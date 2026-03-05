package protocol

type TextMessage struct {
	text string
}

func (t *TextMessage) GetPayload() ([]byte, error) {
	return []byte(t.text), nil
}

func (t *TextMessage) SetPayload(data []byte) error {
	t.text = string(data)
	return nil
}

func (t *TextMessage) Type() MessageType {
	return TypeTextMessage
}

func (t *TextMessage) GetText() string {
	return t.text
}
