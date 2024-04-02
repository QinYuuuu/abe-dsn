package main

import (
	"fmt"

	abedsn "github.com/QinYuuuu/abe-dsn/abe-dsn"
)

func main() {
	/*
		params := pbc.GenerateA(160, 512)
		pairing := params.NewPairing()
		c := pairing.NewZr().Rand()
		fmt.Println(c)

		fmt.Println(pairing.NewZr().PowZn(c, pairing.NewZr().SetBig(big.NewInt(4))))
		hasher := sha256.New()
		hasher.Write([]byte("test"))
		symkey := hasher.Sum(nil)
		N := 4
		F := 1
		escode := erasurecode.NewReedSolomonCode(N-2*F, N)
		escode.Encode(symkey)*/
	attnum := []int{5, 10, 15, 20, 30, 40, 50, 100}
	nodenum := []int{4, 10, 16, 64, 127, 256}
	tnum := []int{1, 3, 5, 21, 42, 85}
	for i := 0; i < len(attnum); i++ {
		for j := 0; j < len(nodenum); j++ {
			fmt.Printf("\nattnum: %v, nodenum %v\n", attnum[i], nodenum[j])
			abedsn.Test(attnum[i], nodenum[j], tnum[j])
		}
	}

}

//func Test_enc_and_upload(attnum int)
