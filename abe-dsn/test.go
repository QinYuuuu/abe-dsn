package abedsn

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"sync"
	"time"

	"github.com/Nik-U/pbc"
	"github.com/QinYuuuu/abe-dsn/cpabe"
	"github.com/QinYuuuu/abe-dsn/vss"
)

type Payload struct {
	C1 *pbc.Element
	C2 *pbc.Element
}

func Test(attnum, nodenum, t int) (time.Duration, int) {
	if nodenum < 3*t+1 {
		fmt.Printf("node number: %v, faulty node number: %v \n", nodenum, t)
		return 0, 0
	}
	atts := make([]string, attnum)
	for i := 0; i < attnum; i++ {
		atts[i] = "att" + strconv.Itoa(i)
	}
	result := ""
	ac := cpabe.AccessStructure{}
	for i := 0; i < attnum; i++ {
		result = result + atts[i] + " and "
	}
	ac.BuildFromPolicy(result[:len(result)-5])

	//fmt.Printf("access structure: \n%+v\n", ac.GetMatrixAsString())

	pairing, abepk, _ := cpabe.Setup(atts)
	r, _ := new(big.Int).SetString("730750818665451621361119245571504901405976559617", 10)
	param := vss.Setup(pairing, abepk.GetGenerateG(), t, nodenum, r)
	hasher := sha256.New()
	hasher.Write([]byte("test"))
	symkey := hasher.Sum(nil)

	lenth := ac.GetL()
	n := ac.GetN()
	start := time.Now()
	s, _ := vss.RandBigInt(r)
	v := []*pbc.Element{pairing.NewZr().SetBig(s)}
	rx := []*pbc.Element{pairing.NewZr().Rand()}
	lambdax := make([]*pbc.Element, n)

	for i := 1; i < lenth; i++ {
		rx = append(rx, pairing.NewZr().Rand())
		v = append(v, pairing.NewZr().Rand())
	}
	for i := 0; i < n; i++ {
		lambdax[i] = cpabe.DotProduct(ac.A[i], v, pairing)
	}
	//
	shares1 := make([][]*big.Int, lenth)
	comm1 := make([][]*pbc.Element, lenth)
	shares2 := make([][]*big.Int, lenth)
	comm2 := make([][]*pbc.Element, lenth)
	var wg sync.WaitGroup
	wg.Add(lenth)
	for i := 0; i < lenth; i++ {
		go func(i int) {
			shares1[i], comm1[i] = vss.Share(param, rx[i].BigInt())
			shares2[i], comm2[i] = vss.Share(param, v[i].BigInt())
			wg.Done()
		}(i)
	}
	c1, c2 := GenerateABEciphertext(pairing.NewGT().SetBytes(symkey), pairing, abepk, ac, s)
	p := Payload{C1: c1, C2: c2}
	pbytes, _ := json.Marshal(p)
	chunks, merklecomm, root := GenerateChunk(pbytes, nodenum, t)
	wg.Wait()
	end := time.Now()
	//fmt.Printf("time: %v\n", end.Sub(start))
	byteAmount := 0
	for j := 0; j < lenth; j++ {
		for i := 0; i < len(shares1[j]); i++ {
			byteAmount += len(shares1[j][i].Bytes())
		}
		for i := 0; i < len(shares2[j]); i++ {
			byteAmount += len(shares2[j][i].Bytes())
		}
		for i := 0; i < len(comm1[j]); i++ {
			byteAmount += len(comm1[j][i].Bytes())
		}
		for i := 0; i < len(comm2[j]); i++ {
			byteAmount += len(comm2[j][i].Bytes())
		}
	}
	for i := 0; i < len(chunks); i++ {
		byteAmount += chunks[i].Size()
		tmp := merklecomm[i].Hash()
		for j := 0; j < len(tmp); j++ {
			byteAmount += len(tmp[j])
		}
	}
	byteAmount += len(root)
	return end.Sub(start), byteAmount
	//fmt.Printf("communication: %vByte\n", byteAmount)
}
