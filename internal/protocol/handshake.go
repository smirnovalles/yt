package protocol

type HandShakeMessage struct {
	id string //name node
}

func (h *HandShakeMessage) GetPayload() ([]byte, error) {
	return []byte(h.id), nil
}

func (h *HandShakeMessage) SetPayload(data []byte) error {
	h.id = string(data)
	return nil
}

func (h *HandShakeMessage) Type() MessageType {
	return TypeHandShake
}

func (h *HandShakeMessage) GetID() string { return h.id }
