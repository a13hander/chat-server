package chat

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Id string
type UserName string

const magicChannelsLen = 100

var chats = make(map[Id]*Chat)
var mxChat sync.RWMutex

var channels = make(map[Id]chan Message)
var mxChannel sync.RWMutex

type Message interface {
	From() string
	Text() string
	CreatedAt() time.Time
}

type Stream interface {
	Send(message Message) error
	Context() context.Context
}

type Chat struct {
	streams map[UserName]Stream
	m       sync.RWMutex
}

func Create() Id {
	id := Id(uuid.New().String())

	channels[id] = make(chan Message, magicChannelsLen)

	return id
}

func Connect(id Id, username UserName, stream Stream) error {
	mxChannel.RLock()
	ch, ok := channels[id]
	mxChannel.RUnlock()

	if !ok {
		return errors.New("chat not found")
	}

	mxChat.Lock()
	if _, chatOk := chats[id]; !chatOk {
		chats[id] = &Chat{
			streams: map[UserName]Stream{},
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

func SendMessage(id Id, message Message) error {
	mxChannel.RLock()
	chatChan, ok := channels[id]
	mxChannel.RUnlock()

	if !ok {
		return errors.New("chat not found")
	}

	chatChan <- message

	return nil
}
