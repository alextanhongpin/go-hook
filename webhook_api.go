package main

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/lib/pq"
	"github.com/satori/go.uuid"
)

type Webhook struct {
	// gorm.Model
	ID              string `sql:"type:uuid;primary_key;"`
	Name            string
	IsVerified      bool
	Status          string
	CallbackURLs    pq.StringArray `gorm:"type:varchar(100)[]"`
	Events          pq.StringArray `gorm:"type:varchar(100)[]"`
	InvocationCount int64
	Version         string
}

func (w *Webhook) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("ID", uuid.NewV4().String())
}

func main() {
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=gorm dbname=webhook password=123456 sslmode=disable")
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&Webhook{})

	// createWebhook(db)
	getWebhook(db)
}

func createWebhook(db *gorm.DB) {
	// Create a new webhook
	db.Create(&Webhook{
		IsVerified:   false,
		Name:         "haha webhook",
		Status:       "haha",
		CallbackURLs: pq.StringArray{"http:localhost:4000/webhooks"},
		Events:       pq.StringArray{"haha:create", "haha:update", "haha:delete"},
		Version:      "0.0.1",
	})
}

func getWebhook(db *gorm.DB) {
	// Get the webhooks
	var webhooks []Webhook
	db.Find(&webhooks)
	log.Printf("found many: %#v\n", webhooks)

	// Get a webhook by id
	var webhook Webhook
	db.First(&webhook, "id = ?", "f7bdc765-2edc-4973-b346-9ec39e9b78ce")
	log.Printf("found one: %+v\n", webhook)
}
