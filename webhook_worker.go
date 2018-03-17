package main

import (
	"log"
	"net/http"

	"github.com/alextanhongpin/go-hook/webhook"
)

func main() {
	wh := webhook.New("book/", "")
	// if err := wh.Info("book/"); err != nil {
	// 	panic(err)
	// }

	// if err := wh.Subscribe("book:create", func(msg []byte) {
	// 	log.Println("got message:", string(msg))
	// }); err != nil {
	// 	panic(err)
	// }

	wh.Fetch("book")

	log.Println("webhook_worker: listening to port *:8080. press ctrl + c to cancel.")
	http.ListenAndServe(":8080", nil)
}
