package cpabe

import "github.com/Nik-U/pbc"

type Message struct {
	mElement *pbc.Element
	mByte    []byte
}

func (m *Message) SetElement(e *pbc.Element) {
	m.mElement = e
}

func (m *Message) GetElement() *pbc.Element {
	return m.mElement
}

func (m *Message) SetByte(b []byte) {
	m.mByte = b
}
