package cpabe

import "github.com/Nik-U/pbc"

type Ciphertext struct {
	C1              *pbc.Element
	C2              *pbc.Element
	D1              map[string]*pbc.Element
	D2              map[string]*pbc.Element
	AccessStructure AccessStructure
	Pairing         *pbc.Pairing
}

func NewCiphertext(pairing *pbc.Pairing, ac AccessStructure) *Ciphertext {
	c := new(Ciphertext)
	c.Pairing = pairing
	c.AccessStructure = ac
	c.D1 = make(map[string]*pbc.Element)
	c.D2 = make(map[string]*pbc.Element)
	return c
}
