package config

import (
	"github.com/Confialink/wallet-pkg-env_config"
)

type MessageBroker struct {
	ClusterID  string
	ClientID   string
	URL        string
	QueueGroup string
}

func initMessageBrokerConfig() *MessageBroker {
	return &MessageBroker{
		// Cluster ID.
		// This is a static value. It was configured during running NATS server.
		ClusterID: "wallet-nats-streaming",

		// It must be unique for each connection. Only alphanumeric and `-` or `_` characters allowed
		ClientID: "nats-workers-users-1",

		// URL to connect to NATS streaming server.
		URL: env_config.Env("VELMIE_WALLET_USERS_MESSAGE_BROKER_URL", ""),

		// Multiple subscriptions group name.
		// Only one member of the group receives the message for Multiple subscriptions.
		QueueGroup: "users",
	}
}
