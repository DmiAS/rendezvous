package proto

import (
	"bytes"
	"encoding/gob"
)

type Registration struct {
	User    string
	Address string
}

func (r *Registration) Marshal() ([]byte, error) {
	buf := new(bytes.Buffer)

	enc := gob.NewEncoder(buf)
	if err := enc.Encode(r); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (r *Registration) Unmarshal(data []byte) error {
	buf := bytes.NewBuffer(data)

	enc := gob.NewDecoder(buf)
	if err := enc.Decode(r); err != nil {
		return err
	}

	return nil
}
