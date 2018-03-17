package webhook

import nats "github.com/nats-io/go-nats"

// Func represents a function that takes byte as parameter
type Func func(params []byte)

type natsQueue struct {
	conn *nats.Conn
}

// Publish publishes the payload to a specific topic
func (q *natsQueue) Publish(event string, payload []byte) error {
	return q.conn.Publish(event, payload)
}

// Subscribe subscribes to a topic and handles it with a func
func (q *natsQueue) Subscribe(event string, fn Func) error {
	_, err := q.conn.Subscribe(event, func(m *nats.Msg) {
		fn(m.Data)
	})
	return err
}

// NewNatsQueue returns a new nats queue
func NewNatsQueue() Queue {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		panic(err)
	}
	return &natsQueue{
		conn: nc,
	}
}
