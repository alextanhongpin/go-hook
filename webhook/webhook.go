package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
)

type webhook struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Version     string    `json:"version"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Events      []Event   `json:"events"`
	Queue       Queue     `json:"-"`
	Store       Store     `json:"-"`
}

// Register will create a new entry for the webhook into the store
func (w *webhook) Register(events ...Event) error {
	w.Events = events

	value, err := json.Marshal(w)
	if err != nil {
		return err
	}

	return w.Store.Put(w.Name, value)
}

// Enable will add the event and the callback url to the store
func (w *webhook) Enable(resource, event, callbackURL string) error {
	keys := []string{resource, event, uuid.Must(uuid.NewV4()).String()}
	return w.Store.Put(strings.Join(keys, "/"), []byte(callbackURL))
}

// Disable will remove the event and the callback url from the store
func (w *webhook) Disable(resource, event, requestID string) error {
	keys := []string{resource, event, requestID}
	return w.Store.Delete(strings.Join(keys, "/"))
}

// Fetch will fetch all events with the prefix
func (w *webhook) Fetch(prefix string) ([]string, error) {
	return w.Store.List(prefix)
}

// Disco returns the info of the webhook based on the name
func (w *webhook) Disco(eventType string) error {
	val, err := w.Store.Get(eventType)
	if err != nil {
		return err
	}

	var wh webhook
	if err := json.Unmarshal(val, &wh); err != nil {
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
func (w *webhook) Publish(eventType string, body interface{}) error {
	if err := w.checkExist(eventType); err != nil {
		return err
	}
	payload, err := json.Marshal(Payload{
		RequestID: uuid.Must(uuid.NewV4()).String(),
		SentAt:    time.Now(),
		Body:      body,
		Subject:   eventType,
	})
	if err != nil {
		return err
	}
	return w.Queue.Publish(eventType, payload)
}

// Subscribe will listen to the registered topic for the payload
func (w *webhook) Subscribe(eventType string, fn Func) error {
	// if err := w.checkExist(eventType); err != nil {
	// 	return err
	// }
	return w.Queue.Subscribe(eventType, fn)
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
