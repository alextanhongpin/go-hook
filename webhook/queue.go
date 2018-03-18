package webhook

// Func represents a function that takes byte as parameter
// type Func func(params []byte)
type Func func([]byte)

// Queue represents the interface for the queue
type Queue interface {
	Publish(evt string, payload []byte) error
	Subscribe(evt string, fn Func) error
}

const (
	// Nats enum is the queue option
	Nats = 0
)

// NewQueue returns a new queue of choice
func NewQueue(opt int) Queue {
	switch opt {
	case Nats:
		return NewNatsQueue()
	default:
		return NewNatsQueue()
	}
}
