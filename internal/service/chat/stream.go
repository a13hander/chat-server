package chat

import (
	"context"

	"github.com/a13hander/chat-server/internal/chat"
	desc "github.com/a13hander/chat-server/pkg/chat_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type stream struct {
	grpcStream desc.ChatV1_ConnectChatServer
}

func NewStream(grpcStream desc.ChatV1_ConnectChatServer) *stream {
	return &stream{grpcStream: grpcStream}
}

func (s *stream) Send(message chat.Message) error {
	msg := &desc.Message{
		From:      message.From(),
		Text:      message.Text(),
		CreatedAt: timestamppb.New(message.CreatedAt()),
	}

	return s.grpcStream.Send(msg)
}

func (s *stream) Context() context.Context {
	return s.grpcStream.Context()
}
