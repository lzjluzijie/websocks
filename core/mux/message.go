package mux

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

const (
	MessageMethodData = iota
	MessageMethodDial
)

//MessageHeadLength = 13
type Message struct {
	Method    uint8
	ConnID    uint32
	MessageID uint32
	Length    uint32
	Data      []byte

	r   io.Reader
	buf []byte
}

func (m *Message) Read(p []byte) (n int, err error) {
	if m.r == nil {
		h := make([]byte, 13)
		h[0] = m.Method
		binary.BigEndian.PutUint32(h[1:5], m.ConnID)
		binary.BigEndian.PutUint32(h[5:9], m.MessageID)
		binary.BigEndian.PutUint32(h[9:13], m.Length)
		m.r = bytes.NewReader(append(h, m.Data...))
	}

	return m.r.Read(p)
}

func LoadMessage(h []byte) (m *Message) {
	if len(h) != 13 {
		panic(fmt.Sprintf("wrong head length: %d", len(h)))
		return
	}

	m = &Message{
		Method:    h[0],
		ConnID:    binary.BigEndian.Uint32(h[1:5]),
		MessageID: binary.BigEndian.Uint32(h[5:9]),
		Length:    binary.BigEndian.Uint32(h[9:13]),
	}
	return
}
