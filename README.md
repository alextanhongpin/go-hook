# webhook-server
A naive implementation of webhook server


## Naive Version

The naive version will just contain a minimum of two servers - one is the __API server__ that will broadcast a payload upon certain event (e.g. create, update, delete) and another is one the __Webhook server__ that will receive the payload. It is just a simple `POST` request to the webhook server, and the payload can only be sent to a single client.
