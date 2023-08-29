package server

import "github.com/winc-link/hummingbird/internal/hummingbird/mqttbroker"

type publishService struct {
	server *server
}

func (p *publishService) Publish(message *mqttbroker.Message) {
	p.server.mu.Lock()
	p.server.deliverMessage("", message, defaultIterateOptions(message.Topic))
	p.server.mu.Unlock()
}
