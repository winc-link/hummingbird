package persistence

import (
	"github.com/winc-link/hummingbird/internal/hummingbird/mqttbroker/config"
	"github.com/winc-link/hummingbird/internal/hummingbird/mqttbroker/persistence/queue"
	mem_queue "github.com/winc-link/hummingbird/internal/hummingbird/mqttbroker/persistence/queue/mem"
	"github.com/winc-link/hummingbird/internal/hummingbird/mqttbroker/persistence/session"
	mem_session "github.com/winc-link/hummingbird/internal/hummingbird/mqttbroker/persistence/session/mem"
	"github.com/winc-link/hummingbird/internal/hummingbird/mqttbroker/persistence/subscription"
	mem_sub "github.com/winc-link/hummingbird/internal/hummingbird/mqttbroker/persistence/subscription/mem"
	"github.com/winc-link/hummingbird/internal/hummingbird/mqttbroker/persistence/unack"
	mem_unack "github.com/winc-link/hummingbird/internal/hummingbird/mqttbroker/persistence/unack/mem"
	"github.com/winc-link/hummingbird/internal/hummingbird/mqttbroker/server"
)

func init() {
	server.RegisterPersistenceFactory("memory", NewMemory)
}

func NewMemory(config config.Config) (server.Persistence, error) {
	return &memory{}, nil
}

type memory struct {
}

func (m *memory) NewUnackStore(config config.Config, clientID string) (unack.Store, error) {
	return mem_unack.New(mem_unack.Options{
		ClientID: clientID,
	}), nil
}

func (m *memory) NewSessionStore(config config.Config) (session.Store, error) {
	return mem_session.New(), nil
}

func (m *memory) Open() error {
	return nil
}
func (m *memory) NewQueueStore(config config.Config, defaultNotifier queue.Notifier, clientID string) (queue.Store, error) {
	return mem_queue.New(mem_queue.Options{
		MaxQueuedMsg:    config.MQTT.MaxQueuedMsg,
		InflightExpiry:  config.MQTT.InflightExpiry,
		ClientID:        clientID,
		DefaultNotifier: defaultNotifier,
	})
}

func (m *memory) NewSubscriptionStore(config config.Config) (subscription.Store, error) {
	return mem_sub.NewStore(), nil
}

func (m *memory) Close() error {
	return nil
}
