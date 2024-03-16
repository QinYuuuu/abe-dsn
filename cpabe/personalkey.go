package cpabe

import "github.com/Nik-U/pbc"

type PersonalKey struct {
	k  *pbc.Element
	l  *pbc.Element
	kx map[string]*pbc.Element
}
