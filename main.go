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
	attnum := []int{5}
	//attnum := []int{1, 5, 10, 15, 20, 30, 40, 50}
	nodenum := []int{4, 8, 16, 32, 64, 128}
	tnum := []int{1, 2, 5, 10, 21, 42}
	for i := 0; i < len(attnum); i++ {
		for j := 0; j < len(nodenum); j++ {
			fmt.Printf("\nattnum: %v, nodenum %v\n", attnum[i], nodenum[j])
			time, byteAmount := abedsn.Test(attnum[i], nodenum[j], tnum[j])
			fmt.Printf("time: %v\n", time)
			fmt.Printf("communication: %vByte\n", byteAmount)
		}
	}

}

//func Test_enc_and_upload(attnum int)
