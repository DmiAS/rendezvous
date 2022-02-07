package proto

import "fmt"

type Body interface {
	Marshal() ([]byte, error)
	Unmarshal(data []byte) error
}

type Packet struct {
	Header *Header
	Data   Body
}

func GetHeader(data []byte) (*Header, []byte, error) {
	header := &Header{}
	if err := header.Unmarshal(data); err != nil {
		return nil, nil, fmt.Errorf("failure to unmarshal header: %s", err)
	}
	return header, data[1:], nil
}

func (p Packet) Marshal() ([]byte, error) {
	headerData, err := p.Header.Marshal()
	if err != nil {
		return nil, fmt.Errorf("failure to marshal header: %s", err)
	}

	bodyData, err := p.Data.Marshal()
	if err != nil {
		return nil, fmt.Errorf("failure to marshal data: %s", err)
	}

	return append(headerData, bodyData...), nil
}

func (p *Packet) Unmarshal(data []byte) error {
	if err := p.Header.Unmarshal(data); err != nil {
		return fmt.Errorf("failure to unmarshal header: %s", err)
	}
	// skip header byte
	if err := p.Data.Unmarshal(data[1:]); err != nil {
		return fmt.Errorf("failure to unmarshal data: %s", err)
	}
	return nil
}
