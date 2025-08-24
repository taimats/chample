package server

type Message struct {
	Type    int
	Content string
}

func NewMessage(msgType int, content string) *Message {
	return &Message{
		Type:    msgType,
		Content: content,
	}
}
