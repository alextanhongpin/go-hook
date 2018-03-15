# webhook-server
A naive implementation of webhook server


## Naive Version

The naive version will just contain a minimum of two servers - one is the __API server__ that will broadcast a payload upon certain event (e.g. create, update, delete) and another is one the __Webhook server__ that will receive the payload. It is just a simple `POST` request to the webhook server, and the payload can only be sent to a single client.

## Better Version

A better version is to create a separate API that will allow user's to select which events they can subscribe to, and what callback url the payload will be posted to.

The __API Server__ will be sending the message to a queue too, instead of a direct `POST` request. A worker will subscribe to the queue in order to process the payload, query the list of subcribers (callback urls) and sending them.
