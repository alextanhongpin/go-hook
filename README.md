# go-hook

A webhook server that makes it easy to send events from your API Server to clients that subscribes to the topic. The initial goal is to make integration to existing applications easy and to provide a UI/CLI that allows user to subscribe/unsubscribe to specific events.


## Design

## Naive Webhook

![naive_one_to_one](assets/naive_one_to_one.png)

In the naive webhook, the _server_ publishes the events directly to the client. The events can be _create_, _update_, _delete_ or other event sourcing events. The event is published using a `POST` request.

![naive_one_to_many](assets/naive_one_to_many.png)

When the number of clients are increasing, it becomes considerably harder to manage them. Aside from that, the server is knowing too much about the clients. It is best to isolate the logic when you are managing many clients.

## Improved Webhook

![worker_one_to_one](assets/worker_one_to_one.png)

The events will still be sent from the Server to the Client, but through the worker instead. The Worker will receive the events through a queue, which increases reliability of the system. The worker will communicate with the Webhook API to validate the subscribers before sending the payload.

![worker_one_to_many](assets/worker_one_to_many.png)

When the number of clients increases, the Worker can be scaled independently too.


![worker_one_to_many](assets/worker_with_events.png)

In the diagram above, the server registers three different endpoints (`PUT /books`, `POST /books`, `DELETE /books`) that will publish the payload to the Worker. Clients can then choose to subscribe to these events through the UI or CLI by passing them the callback url.
