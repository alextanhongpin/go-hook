package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Webhook struct {
	Event     string
	CreatedAt time.Time
	Payload   interface{}
}

func (w *Webhook) Broadcast(url string) error {
	payload, err := json.Marshal(w)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payload)))
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

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/webhooks", func(w http.ResponseWriter, r *http.Request) {
		wh := makeWebhook("get_item", `{"hello": "world"}`)
		if err := wh.Broadcast("http://localhost:4000/webhooks"); err != nil {
			fmt.Printf("Error sending to webhook: %v\n", err)
		}
		fmt.Fprintf(w, `{"hello": "%s"}`, "world")
	})
	fmt.Println("server: listening to port *:8080. press ctrl+c to cancel.")
	http.ListenAndServe(":8080", mux)
}

func makeWebhook(event string, payload interface{}) *Webhook {
	return &Webhook{
		Event:     event,
		CreatedAt: time.Now(),
		Payload:   payload,
	}
}
