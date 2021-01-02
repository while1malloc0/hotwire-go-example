package pubsub

import (
	"context"
	"errors"
	"time"

	"nhooyr.io/websocket"
)

var subscriptions = make(map[uint]*Subscription)

type Subscription struct {
	ID    uint
	C     chan []byte
	Conns []*websocket.Conn

	listening bool
}

func (s *Subscription) Listen() {
	if s.listening {
		return
	}

	go func(s *Subscription) {
		for {
			select {
			case msg := <-s.C:
				for i := range s.Conns {
					ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
					s.Conns[i].Write(ctx, websocket.MessageText, msg)
					cancel()
				}
			}
		}
	}(s)
}

func Subscribe(id uint, conn *websocket.Conn) {
	sub, ok := subscriptions[id]
	if !ok {
		sub = &Subscription{ID: id, C: make(chan []byte), Conns: []*websocket.Conn{conn}}
		sub.Listen()
		subscriptions[id] = sub
		return
	}
	sub.Conns = append(sub.Conns, conn)
}

func Broadcast(id uint, content []byte) error {
	sub, ok := subscriptions[id]
	if !ok {
		return errors.New("No one's listening")
	}
	sub.C <- content
	return nil
}
