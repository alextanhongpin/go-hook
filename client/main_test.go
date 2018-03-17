// package main

// import (
// 	"log"
// 	"time"

// 	"github.com/alextanhongpin/go-hook/webhook"
// )

// type Payload struct {
// 	Name string
// }

// func main() {
// 	wh := webhook.New("book", "book webhook")

// 	if err := wh.Info("book"); err != nil {
// 		log.Println(err)
// 	}
// 	log.Println(wh)

// 	// err := wh.Register([]webhook.WebhookEvent{
// 	// 	webhook.WebhookEvent{Name: "book:create"},
// 	// 	webhook.WebhookEvent{Name: "book:update"},
// 	// 	webhook.WebhookEvent{Name: "book:delete"},
// 	// })
// 	// if err != nil {
// 	// 	panic(err)
// 	// }
// 	// log.Println("successfully registered webhook event")

// 	// err := wh.Subscribe("book:create", func(msg []byte) {
// 	// 	log.Println("got message:", string(msg))
// 	// })
// 	// if err != nil {
// 	// 	log.Println(err)
// 	// }
// 	err = wh.Publish("book:create", Message{Name: "hello world"})
// 	if err != nil {
// 		log.Println(err)
// 	}

// 	err = wh.Publish("book:create", Message{Name: "hi world"})
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	time.Sleep(10 * time.Second)
// 	// mux := http.NewServeMux()
// 	// mux.HandleFunc("/webhooks", func(w http.ResponseWriter, r *http.Request) {
// 	// 	wh := makeWebhook("http://localhost:4000/webhooks", "GET_ITEM", `{"hello": "world"}`)
// 	// 	if err := wh.Broadcast(); err != nil {
// 	// 		fmt.Printf("Error sending to webhook: %v\n", err)
// 	// 	}
// 	// 	fmt.Fprintf(w, `{"hello": "%s"}`, "world")
// 	// })
// 	// fmt.Println("server: listening to port *:8080. press ctrl+c to cancel.")
// 	// http.ListenAndServe(":8080", mux)
// }
