package main

import (
	"log"
	"time"

	"github.com/alextanhongpin/go-hook/webhook"
)

type Payload struct {
	Name string
}

func main() {
	hook := webhook.New(webhook.SetName("book"))

	if err := hook.Register(
		webhook.Event{Name: "book:create"},
		webhook.Event{Name: "book:update"},
		webhook.Event{Name: "book:delete"},
	); err != nil {
		panic(err)
	}
	log.Println("successfully registered events")

	if err := hook.Subscribe("book", func(msg []byte) {
		log.Println("got message:", string(msg))
	}); err != nil {
		log.Println(err)
	}

	if err := hook.Publish("book:create", Payload{Name: "hello world"}); err != nil {
		log.Println(err)
	}

	if err := hook.Publish("book:create", Payload{Name: "hi world"}); err != nil {
		log.Println(err)
	}
	time.Sleep(10 * time.Second)
}
