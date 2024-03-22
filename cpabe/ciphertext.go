package cpabe

import "github.com/Nik-U/pbc"

type Ciphertext struct {
	C1 *pbc.Element
	C2 *pbc.Element
	D1 map[string]*pbc.Element
	D2 map[string]*pbc.Element
}
