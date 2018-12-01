package mux

const (
	MessageMethodData = iota
	MessageMethodDial
)

type MessageHead struct {
	Method    uint8
	ConnID    uint32
	MessageID uint32
	Length    uint32
}
