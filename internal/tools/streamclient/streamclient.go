package streamclient

import (
	"github.com/winc-link/hummingbird/internal/dtos"
	"github.com/winc-link/hummingbird/internal/pkg/logger"
	"time"
)

const (
	pubTimeout = time.Millisecond * 10
)

type streamClient struct {
	msgCh chan dtos.RpcData
	lc    logger.LoggingClient
}

func (c *streamClient) Send(data dtos.RpcData) {
	select {
	case c.msgCh <- data:
	case <-time.After(pubTimeout):
		c.lc.Warnf("send stream message timeout, data: %+v", data)
	}
}

func (c *streamClient) Recv() <-chan dtos.RpcData {
	return c.msgCh
}

func NewStreamClient(lc logger.LoggingClient) *streamClient {
	return &streamClient{
		msgCh: make(chan dtos.RpcData),
		lc:    lc,
	}
}
