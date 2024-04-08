package vss

import (
	"math/big"
	"strconv"

	"github.com/Nik-U/pbc"
)

type Param struct {
	pairing *pbc.Pairing
	g       *pbc.Element
	r       *big.Int
	n       int
	t       int
}

func (p Param) GetT() int {
	return p.t
}

func Setup(pairing *pbc.Pairing, g *pbc.Element, t, n int, r *big.Int) Param {
	param := Param{
		pairing: pairing,
		g:       g,
		r:       r,
		n:       n,
		t:       t,
	}
	return param
}

func Share(param Param, s *big.Int) ([]*big.Int, []*pbc.Element) {
	f := make([]*big.Int, param.t+1)
	c := make([]*pbc.Element, param.t+1)
	f[0] = s
	c[0] = param.pairing.NewG1().PowBig(param.g, s)
	for i := 1; i < param.t+1; i++ {
		f[i], _ = RandBigInt(param.r)
		c[i] = param.pairing.NewG1().PowBig(param.g, f[i])
	}
	//fmt.Printf("f(x): %v\n", f)
	shares := make([]*big.Int, param.n)
	for j := 0; j < param.n; j++ {
		bigj, _ := new(big.Int).SetString(strconv.Itoa(j), 10)
		shares[j] = polynomialEval(f, bigj, param.r)
	}

	return shares, c
}

func Verify(param Param, c []*pbc.Element, share *big.Int, index int) bool {
	left := param.pairing.NewG1().PowBig(param.g, share)
	right := c[0]
	power := 1
	for j := 1; j < param.t+1; j++ {
		power := index * power
		tmp := param.pairing.NewG1().PowBig(c[j], big.NewInt(int64(power)))
		right = param.pairing.NewG1().Mul(right, tmp)
	}
	return left.Equals(right)
}
