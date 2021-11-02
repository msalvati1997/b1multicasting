package basic

import "strconv"

type Message struct {
	MessageHeader map[string]string `json:"MessageHeader"`
	Payload       []byte            `json:"Payload"`
}

func NewMessage(header map[string]string, bytes []byte) Message {
	return Message{
		header,
		bytes,
	}
}

func (m *Message) GetHeaderFieldString(field string) string {
	return m.MessageHeader[field]
}

func (m *Message) GetHeaderFieldInt(field string) int {
	f := m.MessageHeader[field]
	i, _ := strconv.Atoi(f)
	return i
}
