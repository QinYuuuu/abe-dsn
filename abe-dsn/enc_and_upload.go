package abedsn

import (
	"math/big"

	"github.com/Nik-U/pbc"
	"github.com/QinYuuuu/abe-dsn/cpabe"
)

func GenerateABEciphertext(key *pbc.Element, pairing *pbc.Pairing, pk cpabe.ABEpk, ac cpabe.AccessStructure, s *big.Int) (*pbc.Element, map[string]*pbc.Element) {
	g := pk.GetGenerateG()
	eggalpha := pk.Geteggalpha()
	eggalphas := pairing.NewGT().PowBig(eggalpha, s)
	c1 := pairing.NewGT().Mul(key, eggalphas)
	r := []*pbc.Element{pairing.NewZr().Rand()}
	lenth := ac.GetL()
	n := ac.GetN()
	for i := 1; i < lenth; i++ {
		r = append(r, pairing.NewZr().Rand())
	}
	d2 := make(map[string]*pbc.Element)
	rho := ac.GetRho()
	for i := 0; i < n; i++ {
		att := rho[i]
		d2[att] = pairing.NewG1().PowZn(g, r[i])
	}
	return c1, d2
}
