package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/alextanhongpin/go-hook/webhook"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var payload webhook.Payload
		err := json.NewDecoder(r.Body).Decode(&payload)
		fmt.Printf("received body: %v\n", payload)
		if err != nil {
			fmt.Printf("error: %v", err)
		}
		fmt.Fprintf(w, `{"ok": %v}`, true)
	})

	fmt.Println("webhook server: listening to port *:4000. press ctrl + c to cancel.")
	http.ListenAndServe(":4000", mux)
}
