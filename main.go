package main

import (
	"fmt"
	"math/rand"
)

func simulate(f int) {
	total := 100
	time := 100
	fail := 0
	for i := 0; i < time; i++ {
		for try := 0; try < f+1; try++ {
			n := rand.Int() % total
			if n > f {
				continue
			} else {
				fail++
				break
			}
		}

	}
	fmt.Printf("success rate of f %v: %v\n", f, time-fail)
}

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
	/*
		attnum := []int{5}
		//attnum := []int{5, 10, 15, 20, 25, 30, 35, 40, 45, 50}
		nodenum := []int{4, 8, 16, 32, 64, 128}
		tnum := []int{1, 2, 5, 10, 21, 42}
		for i := 0; i < len(attnum); i++ {
			for j := 0; j < len(nodenum); j++ {
				fmt.Printf("\nattnum: %v, nodenum %v\n", attnum[i], nodenum[j])
				time, byteAmount := abedsn.Test(attnum[i], nodenum[j], tnum[j])
				fmt.Printf("time: %v\n", time)
				fmt.Printf("communication: %vByte\n", byteAmount)
			}
		}*/
	for i := 0; i < 33; i++ {
		simulate(i)
	}
}

//func Test_enc_and_upload(attnum int)
