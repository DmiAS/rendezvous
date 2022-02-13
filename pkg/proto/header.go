package proto

import (
	"bytes"
	"encoding/binary"
)

const (
	RegisterAction  = iota // sent by user when he wants to register in system
	RegisterApprove        // approve from server about successful registration

	RequestForConnection // sent by user when he wants to connect to some other user in system
	ResponseForInitiator // sent by server to user, who initiate connection
	ResponseForTarget    // sent by server to the user they want to connect to

	PunchInitiatorMessage // sent by user who is initiator
	PunchTargetMessage    // sent by user who is target
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
