package vss_test

import (
	"fmt"
	"math/big"
	"os"
	"testing"

	"github.com/Nik-U/pbc"
	"github.com/QinYuuuu/abe-dsn/vss"
)

func TestShare(test *testing.T) {
	paramReader, err := os.Open("/home/zhangry2001/abe-dsn/cpabe/a.properties")
	if err != nil {
		fmt.Printf("read a.properties wrong: %v\n", err)
	}
	params, _ := pbc.NewParams(paramReader)
	pairing := params.NewPairing()
	g := pairing.NewG1().Rand()
	r, _ := new(big.Int).SetString("730750818665451621361119245571504901405976559617", 10)
	n := 3
	t := 1
	param := vss.Setup(pairing, g, t, n, r)
	s, _ := new(big.Int).SetString("1", 10)
	shares, commitment := vss.Share(param, s)
	fmt.Printf("shares: %v\n", shares)
	fmt.Printf("commitment: %v\n", commitment)
	for i := 0; i < n; i++ {
		result := vss.Verify(param, commitment, shares[i], i)
		fmt.Printf("verify shares: %v %v\n", i, result)
	}
}
