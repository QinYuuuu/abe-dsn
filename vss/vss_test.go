package vss_test

import (
	"fmt"
	"math/big"
	"os"
	"strconv"
	"testing"

	"github.com/Nik-U/pbc"
	"github.com/QinYuuuu/abe-dsn/cpabe"
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

func TestRec(test *testing.T) {
	paramReader, err := os.Open("/home/zhangry2001/abe-dsn/cpabe/a.properties")
	if err != nil {
		fmt.Printf("read a.properties wrong: %v\n", err)
	}
	params, _ := pbc.NewParams(paramReader)
	pairing := params.NewPairing()
	g := pairing.NewG1().Rand()
	r, _ := new(big.Int).SetString("730750818665451621361119245571504901405976559617", 10)
	n := 4
	t := 1
	param := vss.Setup(pairing, g, t, n, r)
	s, _ := new(big.Int).SetString("1", 10)
	shares, _ := vss.Share(param, s)
	fmt.Printf("shares %v\n", shares)
	result := new(big.Int).SetInt64(0)
	list := make([]int, n)
	list2 := make([]*big.Int, n)
	for i := 0; i < n; i++ {
		list[i] = i + 1
		list2[i], _ = new(big.Int).SetString(strconv.Itoa(i+1), 10)
	}
	fmt.Printf("list %v\n", list)
	fx, _ := vss.LagrangeInterpolation(list2, shares, r)
	for i := 0; i < n; i++ {
		l := cpabe.GenerateLagrangeCoefficient(list, i+1, r)
		l2 := vss.Generate0LagrangeCoefficient(list2, i, r)
		fmt.Printf("ls %v %v\n", l, l2)
		result = new(big.Int).Add(result, new(big.Int).Mul(l2, shares[i]))
	}
	fmt.Printf("test rec: %v %v %v\n", s, result.Mod(result, r), fx[0])
}
