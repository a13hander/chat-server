package chat

import (
	"errors"
	"sync"

	desc "github.com/a13hander/chat-server/pkg/chat_v1"
	"github.com/google/uuid"
)

type Id string
type UserName string

const magicChannelsLen = 100

var chats = make(map[Id]*Chat)
var mxChat sync.RWMutex

var channels = make(map[Id]chan *desc.Message)
var mxChannel sync.RWMutex

type Chat struct {
	streams map[UserName]desc.ChatV1_ConnectChatServer
	m       sync.RWMutex
}

func Create() Id {
	id := Id(uuid.New().String())

	channels[id] = make(chan *desc.Message, magicChannelsLen)

	return id
}

func Connect(id Id, username UserName, stream desc.ChatV1_ConnectChatServer) error {
	mxChannel.RLock()
	ch, ok := channels[id]
	mxChannel.RUnlock()

	if !ok {
		return errors.New("chat not found")
	}

	mxChat.Lock()
	if _, chatOk := chats[id]; !chatOk {
		chats[id] = &Chat{
			streams: map[UserName]desc.ChatV1_ConnectChatServer{},
		}
	}
	mxChat.Unlock()

	chats[id].m.Lock()
	chats[id].streams[username] = stream
	chats[id].m.Unlock()

	for {
		select {
		case msg, okCh := <-ch:
			if !okCh {
				return nil
			}

			for _, st := range chats[id].streams {
				if err := st.Send(msg); err != nil {
					return err
				}
			}

		case <-stream.Context().Done():
			chats[id].m.Lock()
			delete(chats[id].streams, username)
			chats[id].m.Unlock()
			return nil
		}
	}
}

func SendMessage(id Id, message *desc.Message) error {
	mxChannel.RLock()
	chatChan, ok := channels[id]
	mxChannel.RUnlock()

	if !ok {
		return errors.New("chat not found")
	}

	chatChan <- message

	return nil
}
