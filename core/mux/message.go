package mux

import (
	"bytes"
	"encoding/binary"
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

	n, err = m.Read(p)
	return len(p), nil
}

func (m *Message) Write(p []byte) (n int, err error) {
	m.buf = append(m.buf, p...)
	if len(m.buf) >= 13 {
		m.Method = m.buf[0]
		m.ConnID = binary.BigEndian.Uint32(m.buf[1:5])
		m.MessageID = binary.BigEndian.Uint32(m.buf[5:9])
		m.Length = binary.BigEndian.Uint32(m.buf[9:13])
		m.Data = m.buf[13:]
	}
	return len(p), nil
}
