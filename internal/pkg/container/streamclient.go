package container

import (
	"github.com/winc-link/hummingbird/internal/pkg/di"
	"github.com/winc-link/hummingbird/internal/tools/streamclient"
)

var StreamClientName = di.TypeInstanceToName((*streamclient.StreamClient)(nil))

func StreamClientFrom(get di.Get) streamclient.StreamClient {
	return get(StreamClientName).(streamclient.StreamClient)
}
