package mqttclient

import (
	"context"
	"github.com/winc-link/hummingbird/internal/dtos"
)

type MQTTClient interface {
	RegisterConnectCallback(dtos.ConnectHandler)
	RegisterDisconnectCallback(dtos.CallbackHandler)
	AsyncPublish(ctx context.Context, topic string, payload []byte, isSync bool)
	Close()
	GetConnectStatus() bool
}
