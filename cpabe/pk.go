package cpabe

import "github.com/Nik-U/pbc"

type ABEpk struct {
	g        *pbc.Element
	eggalpha *pbc.Element
	ga       *pbc.Element
	h        map[string]*pbc.Element
}

func (pk ABEpk) GetGenerateG() *pbc.Element {
	return pk.g
}

func (pk ABEpk) Geteggalpha() *pbc.Element {
	return pk.eggalpha
}
