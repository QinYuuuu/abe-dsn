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

func Enc() {}

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
