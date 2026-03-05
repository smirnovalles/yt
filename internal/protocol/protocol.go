package protocol

import (
	"encoding/binary"
	"fmt"
)

// type mess
type MessageType byte

const (
	TypeHandShake MessageType = 0x01
	TypeMessage   MessageType = 0x02
)

const (
	MessageTypeLen   = 1
	MessageLenLen    = 2
	MessageHeaderLen = MessageTypeLen + MessageLenLen
	MessageLenMax    = (1 << (MessageLenLen * 8)) - 1
)

type Message interface {
	GetPayload() ([]byte, error)
	SetPayload([]byte) error
	Type() MessageType
}

func setHeader(t MessageType, l int) []byte {
	buf := make([]byte, MessageHeaderLen)
	buf[0] = byte(t)
	binary.BigEndian.PutUint16(buf[1:3], uint16(l))
	return buf

}

func getHeader(data []byte) (t MessageType, payload []byte, err error) {

	len := len(data)

	if len < MessageHeaderLen {
		return t, nil, fmt.Errorf("data head too short")
	}

	t = MessageType(data[0])
	l := binary.BigEndian.Uint16(data[1:3])

	if l > MessageLenMax {
		return t, nil, fmt.Errorf("data len error")
	}

	return t, data[3:], nil

}

func Encode(m Message) (data []byte, err error) {

	payload, err := m.GetPayload()

	if err != nil {
		return nil, err
	}

	mType := m.Type()
	len := len(payload)

	if len > MessageLenMax {
		return nil, fmt.Errorf("message too long: %d", len)
	}

	return append(setHeader(mType, len), payload...), nil

}

func Decode(data []byte) (m Message, err error) {

	t, payload, err := getHeader(data)

	if err != nil {
		return m, err
	}

	var msg Message
	switch t {
	case TypeHandShake:
		msg = &HandShakeMessage{}
	default:
		return nil, fmt.Errorf("unknown type: %v", t)
	}

	err = msg.SetPayload(payload)

	if err != nil {
		return nil, err
	}

	return msg, err
}
