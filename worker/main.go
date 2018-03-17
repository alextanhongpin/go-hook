package main

import (
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

	if err := hook.Subscribe("book.*", func(msg *webhook.Payload, err error) {
		key := strings.Join(strings.Split(msg.Subject, "."), "/")

		subscribers, err := hook.Fetch(key)
		if err != nil {
			panic(err)
		}

		for _, sub := range subscribers {
			log.Println("sending message to:", sub, msg)
		}

	}); err != nil {
		panic(err)
	}

	log.Println("webhook_worker: listening to port *:8080. press ctrl + c to cancel.")
	http.ListenAndServe(":8080", nil)
}
