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

	pairing, abepk, abemsk := cpabe.Setup(atts)
	r, _ := new(big.Int).SetString("730750818665451621361119245571504901405976559617", 10)
	param := vss.Setup(pairing, abepk.GetGenerateG(), t, nodenum, r)
	hasher := sha256.New()
	hasher.Write([]byte("test"))
	symkey := hasher.Sum(nil)

	lenth := ac.GetL()
	n := ac.GetN()
	//start := time.Now()
	s, _ := vss.RandBigInt(r)
	v := []*pbc.Element{pairing.NewZr().SetBig(s)}

	//fmt.Printf("eggalphas %v\n", pairing.NewGT().PowBig(abepk.Geteggalpha(), s))

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
			shares2[i], comm2[i] = vss.Share(param, rx[i].BigInt())
			shares1[i], comm1[i] = vss.Share(param, lambdax[i].BigInt())
			wg.Done()
		}(i)
	}

	c1, c2 := GenerateABEciphertext(pairing.NewGT().SetBytes(symkey), pairing, abepk, ac, s)
	p := Payload{C1: c1, C2: c2}
	pbytes, _ := json.Marshal(p)
	chunks, merklecomm, root := GenerateChunk(pbytes, nodenum, t)
	wg.Wait()
	//end := time.Now()
	//fmt.Printf("time: %v\n", end.Sub(start))
	/*
		for i := 0; i < lenth; i++ {
			fmt.Printf("vi %v\n", v[i])
			fmt.Printf("shares1 : %v\n", shares1[i])
		}*/
	rho := ac.GetRho()
	ga := abepk.Getga()
	g := abepk.GetGenerateG()
	h := abepk.Geth()
	/*
		d10x := make(map[string]*pbc.Element)
		d20x := make(map[string]*pbc.Element)
		for x := 0; x < lenth; x++ {
			att := rho[x]
			tmp1 := pairing.NewG1().PowZn(ga, lambdax[x])
			tmp2 := pairing.NewG1().PowZn(h[att], pairing.NewZr().Neg(rx[x]))
			d10x[att] = pairing.NewG1().Mul(tmp1, tmp2)
			d20x[att] = pairing.NewG1().PowZn(g, rx[x])
		}
	*/
	nodelist := make([]int, nodenum)
	d1xi := make(map[string]map[int]*pbc.Element)
	d2xi := make(map[string]map[int]*pbc.Element)
	for i := 0; i < nodenum; i++ {
		nodelist[i] = i
	}

	for x := 0; x < lenth; x++ {
		att := rho[x]
		d1xi[att] = make(map[int]*pbc.Element)
		d2xi[att] = make(map[int]*pbc.Element)
		for i := 0; i < nodenum; i++ {
			tmp1 := pairing.NewG1().PowBig(ga, shares1[x][i])
			tmp2 := pairing.NewG1().PowZn(h[att], pairing.NewZr().Neg(pairing.NewZr().SetBig(shares2[x][i])))
			d1xi[att][i] = pairing.NewG1().Mul(tmp1, tmp2)
			d2xi[att][i] = pairing.NewG1().PowBig(g, shares2[x][i])
		}
	}

	start := time.Now()
	d1x := make(map[string]*pbc.Element)
	d2x := make(map[string]*pbc.Element)

	for x := 0; x < lenth; x++ {
		att := rho[x]
		d1x[att], d2x[att] = Aggregate(pairing, r, d1xi[att], d2xi[att], nodelist)
		//fmt.Printf("d1x ? %v\n", d10x[att].Equals(d1x[att]))
		//fmt.Printf("d2x ? %v\n", d20x[att].Equals(d2x[att]))
	}

	ct := cpabe.Ciphertext{C1: c1, C2: c2, D1: d1x, D2: d2x, AccessStructure: ac, Pairing: pairing}
	psk, _ := cpabe.KeyGen(pairing, abemsk, abepk, atts)
	cpabe.Dec(ct, psk)
	end := time.Now()
	//fmt.Printf("plaintext: %v\n", pairing.NewGT().SetBytes(symkey))
	//fmt.Printf("decrypt: %v\n", m.GetElement())
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
