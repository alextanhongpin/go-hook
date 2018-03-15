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

type WebhookService struct {
	DB *gorm.DB
}

func (svc *WebhookService) Get(id string) Webhook {
	var webhook Webhook
	svc.DB.First(&webhook, "id = ?", id)
	return webhook
}

func (svc *WebhookService) All() []Webhook {
	var webhooks []Webhook
	svc.DB.Find(&webhooks)
	return webhooks
}

func (svc *WebhookService) Migrate() {
	svc.DB.AutoMigrate(&Webhook{})
}

func main() {
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=gorm dbname=webhook password=123456 sslmode=disable")
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	svc := WebhookService{
		DB: db,
	}

	webhook := svc.Get("f7bdc765-2edc-4973-b346-9ec39e9b78ce")
	log.Printf("found one: %+v\n", webhook)
}
