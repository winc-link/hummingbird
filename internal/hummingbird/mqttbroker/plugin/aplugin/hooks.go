package aplugin

import (
	"context"
	"github.com/winc-link/hummingbird/internal/hummingbird/mqttbroker/server"
)

func (t *APlugin) HookWrapper() server.HookWrapper {
	return server.HookWrapper{}
}

func (t *APlugin) OnBasicAuthWrapper(pre server.OnBasicAuth) server.OnBasicAuth {
	return func(ctx context.Context, client server.Client, req *server.ConnectRequest) (err error) {
		return nil
	}
}

func (t *APlugin) OnSubscribeWrapper(pre server.OnSubscribe) server.OnSubscribe {
	return func(ctx context.Context, client server.Client, req *server.SubscribeRequest) error {
		return nil
	}
}

func (t *APlugin) OnUnsubscribeWrapper(pre server.OnUnsubscribe) server.OnUnsubscribe {
	return func(ctx context.Context, client server.Client, req *server.UnsubscribeRequest) error {
		return nil
	}
}

func (t *APlugin) OnMsgArrivedWrapper(pre server.OnMsgArrived) server.OnMsgArrived {
	return func(ctx context.Context, client server.Client, req *server.MsgArrivedRequest) error {
		return nil
	}
}

func (t *APlugin) OnConnectedWrapper(pre server.OnConnected) server.OnConnected {
	return func(ctx context.Context, client server.Client) {
	}
}

func (t *APlugin) OnClosedWrapper(pre server.OnClosed) server.OnClosed {
	return func(ctx context.Context, client server.Client, err error) {

	}
}
