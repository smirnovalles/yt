package protocol

import (
	"encoding/binary"
	"fmt"
)

// type mess
type MessageType byte

const (
	TypeHandShake   MessageType = 0x01
	TypeTextMessage MessageType = 0x02
)

var validTypes = map[MessageType]bool{
	TypeHandShake:   true,
	TypeTextMessage: true,
}

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

func Decode(msgType MessageType, payload []byte) (m Message, err error) {

	var msg Message
	switch msgType {
	case TypeHandShake:
		msg = &HandShakeMessage{}
	case TypeTextMessage:
		msg = &TextMessage{}
	default:
		return nil, fmt.Errorf("unknown type: %v", msgType)
	}

	err = msg.SetPayload(payload)

	if err != nil {
		return nil, err
	}

	return msg, err
}

func GetHeader(header []byte) (msgType MessageType, msgLen uint16, err error) {

	if len(header) < MessageHeaderLen {
		return 0, 0, fmt.Errorf("header too short")
	}
	msgType = MessageType(header[0])
	msgLen = binary.BigEndian.Uint16(header[1:3])

	if !validTypes[msgType] {
		return 0, 0, fmt.Errorf("unknown message type: %v", msgType)
	}
	return msgType, msgLen, err
}
