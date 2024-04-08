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

func (pk ABEpk) Getga() *pbc.Element {
	return pk.ga
}

func (pk ABEpk) Geteggalpha() *pbc.Element {
	return pk.eggalpha
}
func (pk ABEpk) Geth() map[string]*pbc.Element {
	return pk.h
}
