package streamclient

import "github.com/winc-link/hummingbird/internal/dtos"

type StreamClient interface {
	Send(data dtos.RpcData)
	Recv() <-chan dtos.RpcData
}
