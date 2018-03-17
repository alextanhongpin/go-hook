package main

import (
	"html/template"
	"log"
	"net/http"
)

// webui for the webhook

var (
	templates map[string]*template.Template
)

func init() {
	templates = setupTemplates(map[string]string{
		"index": "public/index.html",
	})
}

func main() {
	mux := http.NewServeMux()

	// Serve static files
	mux.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	// Index
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		templates["index"].Execute(w, nil)
	})

	log.Println("listening to port *:9090. press ctrl + c to cancel.")
	log.Fatal(http.ListenAndServe(":9090", mux))
}

// Setup the templates if it doesn't exist
func setupTemplates(mapping map[string]string) map[string]*template.Template {
	templates := make(map[string]*template.Template)

	for k, v := range mapping {
		t, err := template.ParseFiles(v)
		if err != nil {
			log.Println(err)
		}
		templates[k] = t
	}

	return templates
}
