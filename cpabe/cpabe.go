package cpabe

import (
	"fmt"
	"os"

	"github.com/Nik-U/pbc"
)

func Setup(atts []string) (*pbc.Pairing, ABEpk, ABEmsk) {
	paramReader, err := os.Open("a.properties")
	if err != nil {
		fmt.Printf("read a.properties wrong: %v\n", err)
	}
	params, _ := pbc.NewParams(paramReader)
	pairing := params.NewPairing()
	g := pairing.NewG1().Rand()
	alpha := pairing.NewZr().Rand()
	a := pairing.NewZr().Rand()

	msk := ABEmsk{
		galpha: pairing.NewG1().PowZn(g, alpha),
	}
	h := make(map[string]*pbc.Element)
	for _, att := range atts {
		h[att] = pairing.NewG1().Rand()
	}
	pk := ABEpk{
		g:        g,
		eggalpha: pairing.NewGT().Pair(g, msk.galpha),
		ga:       pairing.NewG1().PowZn(g, a),
		h:        h,
	}
	return pairing, pk, msk
}

func Enc(pairing *pbc.Pairing, m Message, ac AccessStructure, pk ABEpk) (Ciphertext, error) {
	s := pairing.NewZr().Rand()
	c1 := pairing.NewGT().Mul(m.mElement, pairing.NewGT().PowZn(pk.eggalpha, s))
	c2 := pairing.NewG1().PowZn(pk.g, s)
	fmt.Printf("s: %v\n", s)
	fmt.Printf("eggalphas: %v\n", pairing.NewGT().PowZn(pk.eggalpha, s))
	lenth := ac.GetL()
	n := ac.GetN()

	v := []*pbc.Element{s}
	r := []*pbc.Element{pairing.NewZr().Rand()}
	for i := 1; i < lenth; i++ {
		v = append(v, pairing.NewZr().Rand())
		r = append(r, pairing.NewZr().Rand())
	}
	fmt.Printf("v: %v\n", v)
	d1 := make(map[string]*pbc.Element)
	d2 := make(map[string]*pbc.Element)
	for i := 0; i < n; i++ {
		lambdax := DotProduct(ac.A[i], v, pairing)
		att := ac.rho[i]
		tmp := pairing.NewG1().PowZn(pk.ga, lambdax)
		d1[att] = pairing.NewG1().Mul(tmp, pairing.NewG1().PowZn(pk.h[att], pairing.NewZr().Neg(r[i])))
		d2[att] = pairing.NewG1().PowZn(pk.g, r[i])
	}

	c := Ciphertext{
		C1:              c1,
		C2:              c2,
		D1:              d1,
		D2:              d2,
		AccessStructure: ac,
		Pairing:         pairing,
	}
	return c, nil
}

func KeyGen(pairing *pbc.Pairing, msk ABEmsk, pk ABEpk, atts []string) (PersonalKey, error) {
	t := pairing.NewZr().Rand()
	kx := make(map[string]*pbc.Element)
	for _, att := range atts {
		if hx, ok := pk.h[att]; ok {
			kx[att] = pairing.NewG1().PowZn(hx, t)
		} else {
			return PersonalKey{}, fmt.Errorf("attribute %s is not valid", att)
		}
	}
	psersonalKey := PersonalKey{
		K:  pairing.NewG1().Mul(msk.galpha, pairing.NewG1().PowZn(pk.ga, t)),
		L:  pairing.NewG1().PowZn(pk.g, t),
		Kx: kx,
	}
	return psersonalKey, nil
}

func Dec(ct Ciphertext, personalkey PersonalKey) (Message, error) {
	m := Message{}
	pairing := ct.Pairing
	ac := ct.AccessStructure

	atts := make([]string, 0, len(personalkey.Kx))
	for att := range personalkey.Kx {
		atts = append(atts, att)
	}

	toUse := make([]int, 0, len(personalkey.Kx))
	for i := 0; i < len(ac.rho); i++ {

		att := ac.rho[i]
		//fmt.Println(i, att)
		if Contains(atts, att) {
			toUse = append(toUse, i)
		}
	}

	// submatrix of A corresponding to user attribute
	subA := make([][]int, 0, len(personalkey.Kx))
	for _, index := range toUse {
		subA = append(subA, ac.A[index])
	}
	// matrix transpose
	subA_T := make([][]*pbc.Element, len(subA[0]))
	for i := range subA_T {
		subA_T[i] = make([]*pbc.Element, len(subA))
	}

	for i := 0; i < len(subA); i++ {
		for j := 0; j < len(subA[0]); j++ {
			subA_T[j][i] = pairing.NewZr()
			if subA[i][j] == 1 {
				subA_T[j][i].Set1()
			} else if subA[i][j] == 0 {
				subA_T[j][i].Set0()
			} else if subA[i][j] == -1 {
				subA_T[j][i].Neg(pairing.NewZr().Set1())
			}
		}
	}
	fmt.Printf("subA_T: %v\n", subA_T)
	b := []*pbc.Element{pairing.NewZr().Set1()}
	for i := 1; i < len(toUse); i++ {
		b = append(b, pairing.NewZr().Set0())
	}
	cx := GaussianElimination(subA_T, b, pairing)
	fmt.Printf("cx: %v\n", cx)
	tmp := pairing.NewGT().Set1()
	for i, index := range toUse {
		att := ac.rho[index]
		tmp1 := pairing.NewGT().Pair(ct.D1[att], personalkey.L)
		tmp2 := pairing.NewGT().Pair(ct.D2[att], personalkey.Kx[att])
		tmp.Mul(tmp, pairing.NewGT().PowZn(pairing.NewGT().Mul(tmp1, tmp2), cx[i]))
	}
	tmp = pairing.NewGT().Div(pairing.NewGT().Pair(ct.C2, personalkey.K), tmp)
	fmt.Printf("eggalphas: %v\n", tmp)
	m.mElement = pairing.NewGT().Div(ct.C1, tmp)

	return m, nil
}
