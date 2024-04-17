package abedsn

import (
	"math/big"

	"github.com/Nik-U/pbc"
	"github.com/QinYuuuu/abe-dsn/cpabe"
)

func VerifyCPABEShare(index, t int, pairing *pbc.Pairing, pk cpabe.ABEpk, d1 map[string]*pbc.Element, d2 map[string]*pbc.Element, comm1 map[string][]*pbc.Element, comm2 map[string][]*pbc.Element) bool {
	for att, d2share := range d2 {
		right := comm2[att][0]
		power := 1
		for j := 1; j < t+1; j++ {
			power := index * power
			tmp := pairing.NewG1().PowBig(comm2[att][j], big.NewInt(int64(power)))
			right = pairing.NewG1().Mul(right, tmp)
		}
		if d2share.Equals(right) {
			continue
		} else {
			return false
		}
	}
	h := pk.Geth()
	g := pk.GetGenerateG()
	ga := pk.Getga()
	for att, d1share := range d1 {
		left := pairing.NewGT().Pair(d1share, g)
		right2 := pairing.NewGT().Pair(d2[att], h[att])
		right1 := comm1[att][0]
		power := 1
		for j := 1; j < t+1; j++ {
			power := index * power
			tmp := pairing.NewG1().PowBig(comm1[att][j], big.NewInt(int64(power)))
			right1 = pairing.NewG1().Mul(right1, tmp)
		}
		right1 = pairing.NewGT().Pair(ga, right1)
		if left.Equals(right1.Mul(right1, right2)) {
			continue
		} else {
			return false
		}
	}
	return true
}

func Aggregate(pairing *pbc.Pairing, r *big.Int, d1s map[int]*pbc.Element, d2s map[int]*pbc.Element, nodelist []int) (*pbc.Element, *pbc.Element) {
	//fmt.Printf("d1s %v\n", d1s)
	indexlist := make([]int, len(nodelist))
	for i := range nodelist {
		indexlist[i] = nodelist[i] + 1
	}
	//fmt.Printf("indexlist %v\n", indexlist)
	d1 := pairing.NewG1().Set1()
	d2 := pairing.NewG1().Set1()
	for _, j := range indexlist {
		tmp := cpabe.GenerateLagrangeCoefficient(indexlist, j, r)
		//fmt.Printf("la %v\n", tmp)
		d1 = pairing.NewG1().Mul(d1, pairing.NewG1().PowBig(d1s[j-1], tmp))
		d2 = pairing.NewG1().Mul(d2, pairing.NewG1().PowBig(d2s[j-1], tmp))
	}
	return d1, d2
}
