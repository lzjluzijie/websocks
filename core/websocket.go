package core

import (
	"io"

	"github.com/gorilla/websocket"
)

type Conn struct {
	c      *websocket.Conn
	reader io.Reader
	writer io.Writer
}

func (conn *Conn) Read(p []byte) (n int, err error) {
	reader := conn.reader

	if reader == nil {
		_, reader, err = conn.c.NextReader()
		if err != nil {
			return 0, err
		}
		conn.reader = reader
	}

	return reader.Read(p)
}

func (conn *Conn) Write(p []byte) (n int, err error) {
	writer := conn.writer

	if writer == nil {
		writer, err = conn.c.NextWriter(websocket.BinaryMessage)
		if err != nil {
			return 0, err
		}
		conn.writer = writer
	}

	return writer.Write(p)
}
