package tongue

import "fmt"

type TongueMessage struct {
	From    string
	To      string
	Content string
	Encoded bool
}

func NewMessage(from, to, content string, encode bool) TongueMessage {
	msg := TongueMessage{From: from, To: to, Content: content}
	if encode {
		msg.Content = Encode(content)
		msg.Encoded = true
	}
	return msg
}

func (m TongueMessage) String() string {
	return fmt.Sprintf("[%s -> %s]: %s", m.From, m.To, m.Content)
}
