package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	consul "github.com/hashicorp/consul/api"
	nats "github.com/nats-io/go-nats"
)

// Func represents a function that takes byte as parameter
type Func func(params []byte)

// Event represents the individual events that are available for subscription
type Event struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Label       string `json:"label"`
	Enabled     bool   `json:"enable"`
}

// Webhook represents the interface for the webhook package
type Webhook interface {
	Register(events []Event) error
	Unregister(name string) error
	Publish(name string, payload interface{}) error
	Subscribe(name string, fn Func) error
	Post(url, name string, payload interface{}) error
	Info(name string) error
	Fetch(event string) ([]string, error)
	Enable(resource, event, callbackURL string) error
	Disable(resource, event, callbackURL string) error
}

type webhook struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Version     string     `json:"version"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Events      []Event    `json:"events"`
	Queue       *nats.Conn `json:"-"`
	Store       *consul.KV `json:"-"`
}

type Payload struct {
	RequestID string
	Body      interface{}
	CreatedAt time.Time
}

// New returns a pointer to the webhook struct
func New(name, desc string) Webhook {
	client, err := consul.NewClient(consul.DefaultConfig())
	if err != nil {
		panic(err)
	}

	kv := client.KV()

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		panic(err)
	}

	// POST to the webhook api
	return &webhook{
		Name:        name,
		Description: desc,
		Version:     "1.0.0", // TODO: Get from git
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Queue:       nc,
		Store:       kv,
	}
}

// Register will create a new entry for the webhook into the store
func (w *webhook) Register(events []Event) error {
	w.Events = events
	value, err := json.Marshal(w)
	if err != nil {
		return err
	}

	_, err = w.Store.Put(&consul.KVPair{
		Key:   w.Name,
		Value: value,
	}, nil)

	return err
}

// Enable will add the event and the callback url to the store
func (w *webhook) Enable(resource, event, callbackURL string) error {
	keys := []string{resource, event, callbackURL}
	_, err := w.Store.Put(&consul.KVPair{
		Key:   strings.Join(keys, "/"),
		Value: []byte(callbackURL),
	}, nil)
	return err
}

// Disable will remove the event and the callback url from the store
func (w *webhook) Disable(resource, event, callbackURL string) error {
	keys := []string{resource, event, callbackURL}
	_, err := w.Store.Delete(strings.Join(keys, "/"), nil)
	if err != nil {
		return err
	}
	return nil
}

// Fetch will fetch all events with the prefix
func (w *webhook) Fetch(prefix string) ([]string, error) {
	kvPairs, meta, err := w.Store.List(prefix, nil)
	if err != nil {
		return nil, err
	}
	for _, kv := range kvPairs {
		log.Println(kv.Key, string(kv.Value))
	}
	log.Println(meta)
	return nil, nil
}

// Info returns the info of the webhook based on the name
func (w *webhook) Info(eventType string) error {
	kv, meta, err := w.Store.Get(eventType, nil)
	_ = meta

	var wh webhook
	if err := json.Unmarshal(kv.Value, &wh); err != nil {
		return err
	}
	w.Name = wh.Name
	w.Description = wh.Description
	w.Version = wh.Version
	w.CreatedAt = wh.CreatedAt
	w.UpdatedAt = wh.UpdatedAt
	w.Events = wh.Events
	return err
}

// Unregister will remove the webhook from the store
func (w *webhook) Unregister(name string) error {
	return nil
}

// Publish will send the payload to the registered topic
func (w *webhook) Publish(eventType string, payload interface{}) error {
	if err := w.checkExist(eventType); err != nil {
		return err
	}
	// Content enricher - enrich the content with useful information such as UUID and also timestamp
	// publish
	bytePayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	log.Println("sending to", eventType)
	return w.Queue.Publish(eventType, bytePayload)
}

// Subscribe will listen to the registered topic for the payload
func (w *webhook) Subscribe(eventType string, fn Func) error {
	if err := w.checkExist(eventType); err != nil {
		return err
	}
	_, err := w.Queue.Subscribe(eventType, func(m *nats.Msg) {
		fn(m.Data)
	})
	return err
}

// Post will send a POST request to the targetted endpoint
func (w *webhook) Post(url, eventType string, payload interface{}) error {
	// payload, err := json.Marshal(payload)
	// if err != nil {
	// 	return err
	// }
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte("payload")))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println("response:", string(body))
	return nil
}

func (w *webhook) checkExist(eventType string) error {
	var found bool
	for i := 0; i < len(w.Events); i++ {
		evt := w.Events[i]
		if evt.Name == eventType {
			found = true
		}
	}
	if !found {
		return fmt.Errorf("publishError: event with the name %s is not registered", eventType)
	}
	return nil
}
