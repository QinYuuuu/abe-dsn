package cpabe

import "github.com/Nik-U/pbc"

type PersonalKey struct {
	K  *pbc.Element
	L  *pbc.Element
	Kx map[string]*pbc.Element
}
