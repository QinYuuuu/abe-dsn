package cpabe_test

import (
	"fmt"
	"testing"

	"github.com/QinYuuuu/abe-dsn/cpabe"
)

func TestSetup(t *testing.T) {
	_, pk, msk := cpabe.Setup([]string{"a", "b"})
	fmt.Printf("pk: %v \n", pk)
	fmt.Printf("pk: %v \n", msk)
}

func TestAccessStruct(t *testing.T) {
	as1 := new(cpabe.AccessStructure)
	as1.BuildFromPolicy("A and B and C and D")
	fmt.Println(as1.A)
}

func TestKeyGen(t *testing.T) {
	pairing, pk, msk := cpabe.Setup([]string{"A", "B"})
	personalkey, err := cpabe.KeyGen(pairing, msk, pk, []string{"A"})
	if err != nil {
		fmt.Printf("KeyGen Wrong: %v", err)
	}
	fmt.Printf("personalkey: %v", personalkey)
}

func TestEncDec(t *testing.T) {
	pairing, pk, msk := cpabe.Setup([]string{"A", "B"})
	var m cpabe.Message
	m.SetElement(pairing.NewGT().Rand())

	as1 := new(cpabe.AccessStructure)
	as1.BuildFromPolicy("A and B")
	fmt.Printf("matrix: %v\n", as1.A)
	personalkey, err := cpabe.KeyGen(pairing, msk, pk, []string{"A", "B"})
	if err != nil {
		fmt.Printf("KeyGen Wrong: %v", err)
	}

	ct, err := cpabe.Enc(pairing, m, *as1, pk)
	if err != nil {
		fmt.Printf("Encrypt Wrong: %v", err)
	}
	fmt.Printf("ciphertext: %v\n", ct)

	m_, err := cpabe.Dec(ct, personalkey)
	if err != nil {
		fmt.Printf("Decrypt Wrong: %v", err)
	}
	fmt.Printf("verify enc and dec: %v", m_.GetElement().Equals(m.GetElement()))
}
