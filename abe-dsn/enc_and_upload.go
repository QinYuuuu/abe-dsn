package abedsn

import (
	"fmt"
	"math/big"

	"github.com/Nik-U/pbc"
	"github.com/QinYuuuu/abe-dsn/cpabe"
	es "github.com/QinYuuuu/avid-d/erasurecode"
	"github.com/QinYuuuu/avid-d/hasher"

	merkle "github.com/QinYuuuu/avid-d/commit/merklecommitment"
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

func GenerateChunk(symcipher []byte, N, F int) ([]es.ErasureCodeChunk, []merkle.Witness, []byte) {
	escode := es.NewReedSolomonCode(N-2*F, N)
	chunks, err := escode.Encode(symcipher)
	if err != nil {
		fmt.Printf("erasurecode encode wrong: %v\n", err)
		return nil, nil, nil
	}
	dataList := make([][]byte, N)
	for i := 0; i < N; i++ {
		dataList[i] = chunks[i].GetData()
	}
	m, _ := merkle.NewMerkleTree(dataList, hasher.SHA256Hasher)
	witness := make([]merkle.Witness, N)
	for i := 0; i < N; i++ {
		witness[i], err = merkle.CreateWitness(m, i)
		if err != nil {
			fmt.Printf("create witness wrong: %v\n", err)
			return nil, nil, nil
		}
	}
	merklecomm := merkle.Commit(m)
	return chunks, witness, merklecomm
}
