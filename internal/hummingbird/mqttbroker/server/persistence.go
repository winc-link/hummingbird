package server

import (
	"github.com/winc-link/hummingbird/internal/hummingbird/mqttbroker/config"
	"github.com/winc-link/hummingbird/internal/hummingbird/mqttbroker/persistence/queue"
	"github.com/winc-link/hummingbird/internal/hummingbird/mqttbroker/persistence/session"
	"github.com/winc-link/hummingbird/internal/hummingbird/mqttbroker/persistence/subscription"
	"github.com/winc-link/hummingbird/internal/hummingbird/mqttbroker/persistence/unack"
)

type NewPersistence func(config config.Config) (Persistence, error)

type Persistence interface {
	Open() error
	NewQueueStore(config config.Config, defaultNotifier queue.Notifier, clientID string) (queue.Store, error)
	NewSubscriptionStore(config config.Config) (subscription.Store, error)
	NewSessionStore(config config.Config) (session.Store, error)
	NewUnackStore(config config.Config, clientID string) (unack.Store, error)
	Close() error
}
