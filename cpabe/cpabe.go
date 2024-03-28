package cpabe

import (
	"fmt"

	"github.com/Nik-U/pbc"
)

func Setup(atts []string) (*pbc.Pairing, ABEpk, ABEmsk) {
	params := pbc.GenerateA(160, 512)
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

func Enc(pairing *pbc.Pairing, m *Message, ac *AccessStructure, pk ABEpk) Ciphertext {
	s := pairing.NewZr().Rand()
	c1 := pairing.NewGT().Mul(m.mElement, pairing.NewGT().PowZn(pk.eggalpha, s))
	c2 := pairing.NewG1().PowZn(pk.g, s)

	lenth := ac.GetL()
	n := ac.GetN()

	v := []*pbc.Element{s}
	r := []*pbc.Element{pairing.NewZr().Rand()}
	for i := 1; i < lenth; i++ {
		v = append(v, pairing.NewZr().Rand())
		r = append(r, pairing.NewZr().Rand())
	}
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
	return c
}

func KeyGen(pairing *pbc.Pairing, msk ABEmsk, pk ABEpk, atts []string) (PersonalKey, error) {
	t := pairing.NewZr().Rand()
	kx := make(map[string]*pbc.Element)
	for _, att := range atts {
		if hx, ok := pk.h[att]; ok {
			kx[att] = pairing.NewG1().PowZn(hx, t)
		} else {
			return PersonalKey{k: nil, l: nil, kx: nil}, fmt.Errorf("attribute %s is not valid", att)
		}
	}
	psersonalKey := PersonalKey{
		k:  pairing.NewG1().Mul(msk.galpha, pairing.NewG1().PowZn(pk.ga, t)),
		l:  pairing.NewG1().PowZn(pk.ga, t),
		kx: kx,
	}
	return psersonalKey, nil
}

func Dec(ct Ciphertext, personalkey PersonalKey) (*pbc.Element, error) {

	return
}
