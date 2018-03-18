package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
)

// See more examples here: https://github.com/prometheus/client_golang/blob/master/prometheus/examples_test.go
func prometheusHandler() http.Handler {
	return prometheus.Handler()
}

func main() {
	pushCounter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "service_call_test",
		Help: "number of service call to the endpoint",
	})

	err := prometheus.Register(pushCounter)
	if err != nil {
		panic(err)
	}

	go func() {
		time.Sleep(10 * time.Second)
		for i := 0; i < 10; i++ {
			time.Sleep(1 * time.Second)
			for i := 0; i < rand.Intn(200); i++ {
				pushCounter.Inc()
			}
		}
		log.Println("done")
	}()

	webhookCounterVec := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "webhook_requests_total",
			Help: "How many webhook requests are processed, partitioned by event type and callback_url",
		},
		[]string{"event_type", "callback_url"},
	)
	prometheus.MustRegister(webhookCounterVec)
	webhookCounterVec.WithLabelValues("books:create", "http://localhost:1000").Add(20)
	m := webhookCounterVec.WithLabelValues("books:update", "http://password.com")
	go func() {
		time.Sleep(1 * time.Second)
		for i := 0; i < 100; i++ {
			m.Inc()
		}
		log.Println("done too")
	}()

	r := mux.NewRouter()

	// Can be viewed at http_request_duration_microseconds
	r.Handle("/hello", prometheus.InstrumentHandler(
		"hello_endpoint", hello(),
	))

	r.Handle("/metrics", prometheusHandler())

	log.Println("listening to port *:3000. press ctrl + c to cancel.")
	log.Fatal(http.ListenAndServe(":1234", r))
}

type helloWorld struct{}

func (h *helloWorld) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "hello world")
}

func hello() http.Handler {
	return &helloWorld{}
}
