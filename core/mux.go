package core

const (
	MessageMethodData = iota
	MessageMethodDial
)

type MuxConn struct {
	ID       int
	DataID   int
	DataChan chan []byte
}

type Message struct {
	Method    byte
	ConnID    int
	MessageID int
	Data      []byte
}
