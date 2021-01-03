// Package pubsub implements an in-memory publish/subscribe mechanism.
//
// Please don't use this in production.
package pubsub

import (
	"context"
	"errors"
	"time"

	"nhooyr.io/websocket"
)

var subscriptions = make(map[uint]*Subscription)

// Subscription represents a set of subscribers on a topic
type Subscription struct {
	// The ID of the subscription
	ID uint
	// The channel that can be used to publish new content to the topic
	C chan []byte
	// A list of websocket connections that are subscribed to the topic. Updates
	// via the channel C will be fanned out to these
	Conns []*websocket.Conn

	listening bool
}

// Listen starts a Subscription listening for updates to publish to subscribers
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

// Subscribe adds the given websocket as a subscriber to the topic denoted by id
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

// Publish publishes content to all subscribers of the topic denoted by id
func Publish(id uint, content []byte) error {
	sub, ok := subscriptions[id]
	if !ok {
		return errors.New("No one's listening")
	}
	sub.C <- content
	return nil
}
