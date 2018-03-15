package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	nats "github.com/nats-io/go-nats"
)

type Webhook struct {
	ID        string      // For tracing purpose
	CreatedAt time.Time   // For tracing purpose
	UserID    string      // Validate that it can only be sent to this user (to avoid spamming)
	URL       string      // The URL pointing to the webhook worker
	Event     string      // The name of the event
	Payload   interface{} // The payload to be posted
}

// Broadcast posts the event and payload to the callback url.
// Potentially change this to send to a message queue for durability.
func (w *Webhook) Broadcast() error {
	payload, err := json.Marshal(w)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", w.URL, bytes.NewBuffer([]byte(payload)))
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

func (w *Webhook) Publish(nc *nats.Conn) error {
	return nc.Publish("foo", []byte("hello world"))
}

func (w *Webhook) Subscribe(nc *nats.Conn) error {
	nc.Subscribe("foo", func(m *nats.Msg) {
		fmt.Println("received a message: %s\n", string(m.Data))
	})
}

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/webhooks", func(w http.ResponseWriter, r *http.Request) {
		wh := makeWebhook("http://localhost:4000/webhooks", "GET_ITEM", `{"hello": "world"}`)
		if err := wh.Broadcast(); err != nil {
			fmt.Printf("Error sending to webhook: %v\n", err)
		}
		fmt.Fprintf(w, `{"hello": "%s"}`, "world")
	})
	fmt.Println("server: listening to port *:8080. press ctrl+c to cancel.")
	http.ListenAndServe(":8080", mux)
}

func makeWebhook(url, event string, payload interface{}) *Webhook {
	return &Webhook{
		Event:     event,
		CreatedAt: time.Now(),
		Payload:   payload,
	}
}
