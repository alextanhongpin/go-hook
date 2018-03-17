# Webhook Server

As a developer, I want to send subscribe to events from an API, in order to publish it to a webhook server.


## Naive Version

```
[ API Server ] POST /callback_urls --> [ Webhook Server ]
```

The naive version will just contain a minimum of two servers:

- __API server__ is responsible for publishing the event (e.g. create, update, delete) and associated payload to the webhook server. For simplicity, the event is published using a `POST` request, and only one recipient can be registered at a time.
- __Webhook server__ is the end client. The __API Server__ will post to this endpoint.

## Better Version

```bash
# User subscribes to the webhook
[ User ] POST /webhooks?callback_urls=http://example.com --> [ Webhook API ]

# To handle the load, the payload is sent to the queue first
[ API Server ] SendToQueue --> [ Webhook Worker ] POST /callback_urls --> [ Webhook Server ]
                                        |
				        | Check subscribers and get their callback urls
					|
			         [ Webhook API ]
```

A better version is to create a separate API that will allow user's to select which events they can subscribe to, and what callback url the payload will be posted to.

The __API Server__ will be sending the message to a queue too, instead of a direct `POST` request. A worker will subscribe to the queue in order to process the payload, query the list of subcribers (callback urls) and sending them.

- webhook_api: allows users to create new webhook subscriptions
- webhook_server: the webhook will post to this server
- webhook_worker: the api server will send the payload through a queue to the webhook worker which will validate the payload through the webhook api and send the message to the webhook server

## Webhook API

The Webhook API allows users to subscribe to events and provide a callback url where the events will be posted.

| Method | Endpoint | Description | 
|--      |--        |--           |
| GET    | `/webhooks` | Get a list of webhook subscriptions | 
| POST   | `/webhooks` | Create a webhook with the list of events to be subscribed, and the callback url | 
| DELETE | `/webhooks` | Clear all registered webhooks |
| GET    | `/webhooks/{id}` | Get the info for a specific webhook |
| PUT    | `/webhooks/{id}` | Update the info for a specific webhook |
| DELETE | `/webhooks/{id}` | Delete a webhook by id |

## Webhook API Model

Webhook:

```js
{
	"id": "1",
	"created_at": "",
	"updated_at": "",
	"user_id": "",
	"is_verified": false, 
	"status": "active|error|stop",
	"callback_urls": [],
	"events": ["books:get", "user:create"], // Do we allow user's to subscribe to different resource topics?
	"invocation_count": 0,
	"version": "0.0.1"
}
```

Webhook Events:

```js
{
	"id": "",
	"name": "books:get",
	"service": "", // The service triggering this
	"count": 0,
	"created_at": "",
	"updated_at": "",
	"callback_url": "",
	"batch_size": 10, // Number of items per batch
	"error_count": 0, // 
	"retry_policy": {},
	"version": "git version"
}
```

Webhook API:

```js
{
	"name": "books api",
	"description": "books that serves api",
	// "events": [
	// 	"books:create",
	// 	"books:update",
	// 	"books:delete"
	// ],
	"events": [
		{
			"name": "books:create",
			"description": "",
			"created_at": "",
			"updated_at": "",
			"enabled": true,
			"payload": {},
			"metadata": {}
		}
	],
	"created_at": "",
	"updated_at": "",
	"version": ""
}
```

## Homogenous/Heteregenous Events

Homogenous events can be different events for the same resource, e.g. `books:get`, `books:create`.

Heteregenous events can be different events and different resources, e.g. `books:create`, `users:create`.

It is preferable to store each event and each resource in a new row to simplify query. Caching can be done through Redis too to reduce calls to the database, and each worker can just point to a redis cluster.

## Internal and External Webhook

If the webhook is open for public to consume (e.g. Slack, Github), then it will require certain authorization. The identity of the creator needs to be embedded during the creation of the webhook too.

## UI

The UI for selecting the topics to subscribe should contain the bare minimum:

```
Name: Book Webhook
Description: Contains book events that user can subscribe to
Created At: 1s ago
Updated At: 1s ago

Callback Urls (comma-separated): _____________________ [ VERIFY ]

Events:
[x] book:create
[x] book:update
[-] book:delete

[ OK ]
```

## Security

Some thoughts and scenarios that could happen:

- can I register any URL?
- do I need to confirm the URL I am posting at?
- can I post to other user's URL (DDOS)?
- how can I verify once-only delivery?
- what if the server to be posted is down?
- can I subscribe multiple urls for the same topic?
- are there any retry policy?
- batching requests?
- how do I deal with changes to the event name

## Interface

The webhook package should contain the following interface

| Method | Description |
|--      |--           |
| register | Register a new service and events that users can subscribe to. If the service already exists, it will be updated. |
| deregister | Remove an existing service from the store |
| update | Update an existing service, if the service does not exist, throws an error |
| publish | Publish the payload to a topic |
| subscribe | Subscribes to the topic and receives the payload |
| fetch | Get all the topics subscribed and the equivalent callback urls to be stored in memory as a dictionary |

## Dependencies

- Nats as the messaging queue
- Consul as the key-value store

```bash
# Go client
go get github.com/nats-io/go-nats

# Server
go get github.com/nats-io/gnatsd
```
