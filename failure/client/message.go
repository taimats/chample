package client

type Message struct {
	From string `json:"from"`
	Text string `json:"text"`
}

func NewMessage(from string, text string) *Message {
	return &Message{
		From: from,
		Text: text,
	}
}
