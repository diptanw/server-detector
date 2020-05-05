package events

import (
	"github.com/nats-io/nats.go"
)

// MsgHandler is a callback function that processes messages delivered to
// asynchronous subscribers
type MsgHandler func(msg *Message)

// Message defines the structure of a message to be received/published
type Message struct {
	Data []byte
}

// Connect creates and connects a new instance of Messenger
func Connect(serverURI string) (*Messenger, error) {
	conn, err := nats.Connect(serverURI)
	if err != nil {
		return nil, err
	}

	return &Messenger{
		conn: conn,
	}, nil
}

// Messenger is wrapper for nats.Client to abstract the implementation
type Messenger struct {
	conn *nats.Conn
	subs []*nats.Subscription
}

func subscribeCallback(cb MsgHandler) nats.MsgHandler {
	return func(m *nats.Msg) {
		msg := &Message{m.Data}
		cb(msg)
	}
}

// Subscribe add a new subscription to message queue
func (rc *Messenger) Subscribe(subject string, cb MsgHandler) error {
	sub, err := rc.conn.Subscribe(subject, subscribeCallback(cb))
	rc.subs = append(rc.subs, sub)

	return err
}

// Publish publishes a new message to message queue
func (rc *Messenger) Publish(subject string, msg Message) error {
	return rc.conn.Publish(subject, msg.Data)
}

func (rc *Messenger) Close() {
	for _, s := range rc.subs {
		s.Unsubscribe() // nolint[errcheck]
		s.Drain()       // nolint[errcheck]
	}

	rc.subs = nil

	rc.conn.Close()
}
