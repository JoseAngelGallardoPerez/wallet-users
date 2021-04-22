package messagebroker

// Function to handle received messages from a Broker
type MessageHandler func(jsonData string)

type MessageBroker interface {
	// Publish a message as an asynchronous process.
	// Use this method if a delivery is not important and the message can be lost.
	PublishAsync(subject string, data interface{}) error

	// Publish a message as an synchronous process.
	// Use this method if a delivery is important.
	Publish(subject string, data interface{}) error

	// Multiple subscriptions using the same channel and queue name are members of the same queue group.
	// That means that if a message is published on that channel, only one member of the group receives the message.
	// Other subscriptions receive messages independently of the queue groups, that is, a message is delivered to all
	// subscriptions and one member of each queue group.
	QueueSubscribe(subject string, handler MessageHandler) error
}
