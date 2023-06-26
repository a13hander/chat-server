package chat

import (
	"time"

	"github.com/a13hander/chat-server/internal/chat"
	desc "github.com/a13hander/chat-server/pkg/chat_v1"
)

type msg struct {
	from    string
	text    string
	created time.Time
}

func ConvertMessageToDomain(m *desc.Message) chat.Message {
	return &msg{
		from:    m.GetFrom(),
		text:    m.GetText(),
		created: m.GetCreatedAt().AsTime(),
	}
}

func (m *msg) From() string {
	return m.from
}

func (m *msg) Text() string {
	return m.text
}

func (m *msg) CreatedAt() time.Time {
	return m.created
}
