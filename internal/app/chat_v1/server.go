package chat_v1

import (
	"context"

	"github.com/a13hander/chat-server/internal/chat"
	desc "github.com/a13hander/chat-server/pkg/chat_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Implementation struct {
	desc.UnimplementedChatV1Server
}

func NewImplementation() *Implementation {
	return &Implementation{}
}

func (i *Implementation) CreateChat(ctx context.Context, req *desc.CreateChatRequest) (*desc.CreateChatResponse, error) {
	return &desc.CreateChatResponse{ChatId: string(chat.Create())}, nil
}

func (i *Implementation) ConnectChat(req *desc.ConnectChatRequest, stream desc.ChatV1_ConnectChatServer) error {
	return chat.Connect(
		chat.Id(req.GetChatId()),
		chat.UserName(req.GetUsername()),
		stream,
	)
}

func (i *Implementation) SendMessage(ctx context.Context, req *desc.SendMessageRequest) (*emptypb.Empty, error) {
	err := chat.SendMessage(
		chat.Id(req.GetChatId()),
		req.Message,
	)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
