package webhook

import (
	"time"
)

// Event represents the individual events that are available for subscription
type Event struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Label       string `json:"label"`
	Enabled     bool   `json:"enable"`
}

// Webhook represents the interface for the webhook package
type Webhook interface {
	Register(events ...Event) error
	Unregister(name string) error
	Publish(name string, payload interface{}) error
	Subscribe(name string, fn Func) error
	Post(url string, payload []byte) error
	Disco(name string) error
	Fetch(event string) ([]string, error)
	Enable(resource, event, callbackURL string) error
	Disable(resource, event, callbackURL string) error
}

// Payload represents the message that is passed through the queue with additional metadata
type Payload struct {
	Subject   string
	RequestID string
	Body      interface{}
	SentAt    time.Time
}

// Option returns a function to the pointer of webhook
type Option func(*webhook)

// SetName takes a name and overwrites the default name
func SetName(name string) Option {
	return func(w *webhook) {
		w.Name = name
	}
}

// SetDescription takes a description and overwrites the default description
func SetDescription(desc string) Option {
	return func(w *webhook) {
		w.Description = desc
	}
}

// New returns a pointer to the webhook struct
func New(opts ...Option) Webhook {
	// Sane defaults
	wh := webhook{
		Name:        "webhook",
		Description: "a default webhook",
		Version:     "1.0.0", // TODO: Get from git
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Queue:       NewQueue(Nats),
		Store:       NewStore(Consul),
	}

	for _, o := range opts {
		o(&wh)
	}
	return &wh
}

// BasicEvent takes a name and return an Event
func BasicEvent(name string) Event {
	return Event{
		Name:        name,
		Description: "",
		Label:       "",
		Enabled:     true,
	}
}
