package messagebroker

import (
	"encoding/json"
	"os"

	"github.com/inconshreveable/log15"
	"github.com/nats-io/stan.go"

	"github.com/Confialink/wallet-users/internal/config"
)

type Nats struct {
	logger     log15.Logger
	connection stan.Conn
	queueGroup string
}

func NewNats(l log15.Logger, cfg *config.Configuration) MessageBroker {
	logger := l.New("Broker", "NATS")

	connectionLostHandler := stan.SetConnectionLostHandler(func(cn stan.Conn, err error) {
		logger.Error("connection lost", "err", err, "configuration", cfg.MessageBroker)

		// TODO: add reconnect here. The library must provide it.
		// TODO: See the AllowReconnect option and the doReconnect method
	})

	conn, err := stan.Connect(
		cfg.MessageBroker.ClusterID,
		cfg.MessageBroker.ClientID,
		stan.NatsURL(cfg.MessageBroker.URL),
		connectionLostHandler,
	)
	if err != nil {
		logger.Error("cannot connect to NATS", "err", err, "configuration", cfg.MessageBroker)
		os.Exit(1)
	}

	return &Nats{logger: logger, connection: conn, queueGroup: cfg.MessageBroker.QueueGroup}
}

// Publish a message as an asynchronous process.
// It does not return an error if NATS server is offline.
// Use this method if a delivery is not important and the message can be lost.
func (s *Nats) PublishAsync(subject string, data interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		s.logger.New("method", "PublishAsync").Error("cannot marshal data", "err", err, "subject", subject, "data", data)
		return err
	}

	ackHandler := func(_ string, err error) {
		if err != nil {
			s.logger.New("method", "PublishAsync").Error("cannot do Acknowledgements", "err", err, "subject", subject)
		}
	}
	_, err = s.connection.PublishAsync(subject, b, ackHandler)
	if err != nil {
		s.logger.New("method", "PublishAsync").Error("cannot publish data", "err", err, "subject", subject)
	}

	return nil
}

// Publish a message as an synchronous process.
// It returns an error if NATS server does not receive the message.
// Use this method if a delivery is important.
func (s *Nats) Publish(subject string, data interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		s.logger.New("method", "Publish").Error("cannot marshal data", "err", err, "subject", subject, "data", data)
		return err
	}

	err = s.connection.Publish(subject, b)
	if err != nil {
		s.logger.New("method", "Publish").Error("cannot publish data", "err", err, "subject", subject)
		return err
	}

	return nil
}

// Multiple subscriptions using the same channel and queue name are members of the same queue group.
// That means that if a message is published on that channel, only one member of the group receives the message.
// Other subscriptions receive messages independently of the queue groups, that is, a message is delivered to all
// subscriptions and one member of each queue group.
func (s *Nats) QueueSubscribe(subject string, handler MessageHandler) error {
	opt := []stan.SubscriptionOption{
		stan.DurableName(subject),
		stan.DeliverAllAvailable(),
	}

	_, err := s.connection.QueueSubscribe(subject, s.queueGroup, func(m *stan.Msg) {
		handler(string(m.Data))
	}, opt...)

	return err
}
