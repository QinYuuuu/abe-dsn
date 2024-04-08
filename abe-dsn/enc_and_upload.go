package abedsn

import (
	"fmt"
	"math/big"
	"sync"

	"github.com/Nik-U/pbc"
	"github.com/QinYuuuu/abe-dsn/cpabe"
	es "github.com/QinYuuuu/avid-d/erasurecode"
	"github.com/QinYuuuu/avid-d/hasher"

	merkle "github.com/QinYuuuu/avid-d/commit/merklecommitment"
)

func GenerateABEciphertext(key *pbc.Element, pairing *pbc.Pairing, pk cpabe.ABEpk, ac cpabe.AccessStructure, s *big.Int) (*pbc.Element, *pbc.Element) {
	g := pk.GetGenerateG()
	eggalpha := pk.Geteggalpha()
	/*
		ga := pk.Getga()
		h := pk.Geth()
	*/
	eggalphas := pairing.NewGT().PowBig(eggalpha, s)
	c1 := pairing.NewGT().Mul(key, eggalphas)
	c2 := pairing.NewG1().PowBig(g, s)
	return c1, c2
}

func GenerateChunk(symcipher []byte, N, F int) ([]es.ErasureCodeChunk, []merkle.Witness, []byte) {
	escode := es.NewReedSolomonCode(N-2*F, N)
	chunks, err := escode.Encode(symcipher)
	if err != nil {
		fmt.Printf("erasurecode encode wrong: %v\n", err)
		return nil, nil, nil
	}
	dataList := make([][]byte, N)
	var wg sync.WaitGroup
	wg.Add(N)
	for i := 0; i < N; i++ {
		go func(i int) {
			dataList[i] = chunks[i].GetData()
			wg.Done()
		}(i)
	}
	wg.Wait()
	m, _ := merkle.NewMerkleTree(dataList, hasher.SHA256Hasher)
	witness := make([]merkle.Witness, N)
	wg.Add(N)
	for i := 0; i < N; i++ {
		go func(i int) {
			witness[i], _ = merkle.CreateWitness(m, i)
			wg.Done()
		}(i)
	}
	wg.Wait()
	merklecomm := merkle.Commit(m)
	return chunks, witness, merklecomm
}
