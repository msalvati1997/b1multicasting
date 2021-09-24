package basic

type Message struct {
	MessageHeader map[string]string
	Payload       []byte
}

func NewMessage(bytes []byte) Message {
	return Message{
		make(map[string]string),
		bytes,
	}
}
