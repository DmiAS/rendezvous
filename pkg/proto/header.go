package proto

import (
	"bytes"
	"encoding/binary"
)

const (
	RegisterAction = iota
	RegisterApprove
	RequestForConnection
	ResponseConnection
	MessageReceive
)

type Header struct {
	Action uint8
}

func (h *Header) Marshal() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, h); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (h *Header) Unmarshal(data []byte) error {
	buf := bytes.NewBuffer(data)
	if err := binary.Read(buf, binary.LittleEndian, h); err != nil {
		return err
	}
	return nil
}
