// package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"time"
// )

// type Webhook struct {
// 	Event     string
// 	CreatedAt time.Time
// 	Payload   interface{}
// }

// func main() {
// 	mux := http.NewServeMux()
// 	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 		var webhook Webhook
// 		err := json.NewDecoder(r.Body).Decode(&webhook)
// 		fmt.Printf("received body: %v\n", webhook)
// 		if err != nil {
// 			fmt.Printf("error: %v", err)
// 		}
// 		fmt.Fprintf(w, `{"ok": %v}`, true)
// 	})
// 	fmt.Println("webhook server: listening to port *:4000. press ctrl + c to cancel.")
// 	http.ListenAndServe(":4000", mux)
// }
