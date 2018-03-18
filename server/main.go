package main

import (
	"log"
	"time"

	"github.com/alextanhongpin/go-hook/webhook"
)

// Payload represents the payload to be passed to the queues
type Payload struct {
	Name string
}

func main() {
	hook := webhook.New(
		webhook.SetName("book"),
		webhook.SetDescription("subscribe to book events for changes"),
	)

	if err := hook.Register(
		webhook.BasicEvent("book.create"),
		webhook.BasicEvent("book.update"),
		webhook.BasicEvent("book.delete"),
	); err != nil {
		panic(err)
	}

	if err := hook.Publish("book.create", Payload{Name: "hello world"}); err != nil {
		log.Println(err)
	}

	if err := hook.Publish("book.update", Payload{Name: "hi world"}); err != nil {
		log.Println(err)
	}

	time.Sleep(1 * time.Second)
}
