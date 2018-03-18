package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/alextanhongpin/go-hook/webhook"
)

func main() {
	hook := webhook.New(webhook.SetName("book"))
	if err := hook.Disco("book"); err != nil {
		panic(err)
	}

	// Create a new webhook subscription
	// if err := hook.Enable("book", "update", "http://localhost:4000"); err != nil {
	// 	log.Println(err)
	// }

	if err := hook.Subscribe("book.*", func(msg []byte) {
		var p webhook.Payload
		if err := json.Unmarshal(msg, &p); err != nil {
			log.Println(err)
			return
		}
		key := strings.Join(strings.Split(p.Subject, "."), "/")

		subscribers, err := hook.Fetch(key)
		if err != nil {
			panic(err)
		}

		for _, sub := range subscribers {
			if err := hook.Post(sub, msg); err != nil {
				log.Printf("sendError: %s\n", err.Error())
				continue
			}
			// Increment count back to store (or analyse through logs)
			log.Println("sending message to:", sub, p)
		}

	}); err != nil {
		panic(err)
	}

	log.Println("webhook_worker: listening to port *:8080. press ctrl + c to cancel.")
	http.ListenAndServe(":8080", nil)
}
