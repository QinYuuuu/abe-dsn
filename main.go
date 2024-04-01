package main

import (
	"crypto/sha256"
	"fmt"
	"math/big"

	"github.com/Nik-U/pbc"
	"github.com/QinYuuuu/avid-d/erasurecode"
)

func main() {
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
	escode.Encode(symkey)
}

//func Test_enc_and_upload(attnum int)
