package cpabe

import (
	"fmt"

	"github.com/Nik-U/pbc"
)

func Setup(atts []string) (ABEpk, ABEmsk) {
	params := pbc.GenerateA(160, 512)
	pairing := params.NewPairing()
	g := pairing.NewG1().Rand()
	alpha := pairing.NewZr().Rand()
	a := pairing.NewZr().Rand()

	msk := ABEmsk{
		galpha: pairing.NewG1().PowZn(g, alpha),
	}
	var h map[string]*pbc.Element
	for _, att := range atts {
		h[att] = pairing.NewG1().Rand()
	}
	pk := ABEpk{
		g:        g,
		eggalpha: pairing.NewGT().Pair(g, msk.galpha),
		ga:       pairing.NewG1().PowZn(g, a),
		h:        h,
	}
	return pk, msk
}

func Enc() {}

func KeyGen(pk ABEpk, atts []string) (PersonalKey, error) {
	for _, att := range atts {
		if hx, ok := pk.h[att]; ok {

		} else {
			return PersonalKey{k: nil, l: nil, kx: nil}, fmt.Errorf("Attribute %s is not valid", att)
		}
	}
}

func Dec() {}
