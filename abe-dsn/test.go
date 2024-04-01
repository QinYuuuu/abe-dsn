package abedsn

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/Nik-U/pbc"
	"github.com/QinYuuuu/abe-dsn/cpabe"
	"github.com/QinYuuuu/abe-dsn/vss"
)

type payload struct {
	c1 *pbc.Element
	d2 map[string]*pbc.Element
}

func Test(attnum, nodenum, t int) {
	if nodenum < 3*t+1 {
		fmt.Printf("node number: %v, faulty node number: %v \n", nodenum, t)
		return
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

	fmt.Printf("access structure: \n%+v\n", ac.GetMatrixAsString())

	pairing, abepk, _ := cpabe.Setup(atts)
	r, _ := new(big.Int).SetString("730750818665451621361119245571504901405976559617", 10)
	param := vss.Setup(pairing, abepk.GetGenerateG(), t, nodenum, r)
	hasher := sha256.New()
	hasher.Write([]byte("test"))
	symkey := hasher.Sum(nil)

	start := time.Now()
	s, _ := vss.RandBigInt(r)
	shares, comm := vss.Share(param, s)

	c1, d2 := GenerateABEciphertext(pairing.NewGT().SetBytes(symkey), pairing, abepk, ac, s)
	p := payload{c1: c1, d2: d2}
	pbytes, _ := json.Marshal(p)
	chunks, merklecomm, root := GenerateChunk(pbytes, nodenum, t)
	end := time.Now()
	fmt.Printf("time: %v", end.Sub(start))
}
