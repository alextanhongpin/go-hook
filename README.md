# webhook-server
A naive implementation of webhook server


## Naive Version

The naive version will just contain a minimum of two servers - one is the __API server__ that will broadcast a payload upon certain event (e.g. create, update, delete) and another is one the __Webhook server__ that will receive the payload. It is just a simple `POST` request to the webhook server, and the payload can only be sent to a single client.

## Better Version

A better version is to create a separate API that will allow user's to select which events they can subscribe to, and what callback url the payload will be posted to.

The __API Server__ will be sending the message to a queue too, instead of a direct `POST` request. A worker will subscribe to the queue in order to process the payload, query the list of subcribers (callback urls) and sending them.


## Webhook API

The Webhook API allows users to subscribe to events and provide a callback url where the events will be posted.

| Endpoint | Description | 
|--        |--           |
| GET `/webhooks` | Get a list of webhook subscriptions | 
| POST `/webhooks` | Create a webhook with the list of events to be subscribed, and the callback url | 
| DELETE `/webhooks` | Clear all registered webhooks |
| GET `/webhooks/{id}` | Get the info for a specific webhook |
| PUT `/webhooks/{id}` | Update the info for a specific webhook |
| DELETE `/webhooks/{id}` | Delete a webhook by id |

## Internal and External Webhook

If the webhook is open for public to consume (e.g. Slack, Github), then it will require certain authorization. The identity of the creator needs to be embedded during the creation of the webhook too.

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

