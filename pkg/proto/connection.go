package proto

import (
	"bytes"
	"encoding/gob"
)

type ConnRequest struct {
	Initiator string
	Target    string
}

type ConnResponse struct {
	LocalAddress  string
	GlobalAddress string
}

type Punch struct {
	Msg string
}

func (c ConnRequest) Marshal() ([]byte, error) {
	buf := new(bytes.Buffer)

	enc := gob.NewEncoder(buf)
	if err := enc.Encode(c); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (c *ConnRequest) Unmarshal(data []byte) error {
	buf := bytes.NewBuffer(data)

	enc := gob.NewDecoder(buf)
	if err := enc.Decode(c); err != nil {
		return err
	}

	return nil
}

func (c ConnResponse) Marshal() ([]byte, error) {
	buf := new(bytes.Buffer)

	enc := gob.NewEncoder(buf)
	if err := enc.Encode(c); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (c *ConnResponse) Unmarshal(data []byte) error {
	buf := bytes.NewBuffer(data)

	enc := gob.NewDecoder(buf)
	if err := enc.Decode(c); err != nil {
		return err
	}

	return nil
}

func (p Punch) Marshal() ([]byte, error) {
	buf := new(bytes.Buffer)

	enc := gob.NewEncoder(buf)
	if err := enc.Encode(p); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (p *Punch) Unmarshal(data []byte) error {
	buf := bytes.NewBuffer(data)

	enc := gob.NewDecoder(buf)
	if err := enc.Decode(p); err != nil {
		return err
	}

	return nil
}
