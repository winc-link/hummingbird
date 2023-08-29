package aplugin

import (
	"encoding/json"
	"sync"

	"github.com/winc-link/hummingbird/internal/hummingbird/mqttbroker/server"
)

const (
	Auth = iota + 1 // 连接鉴权
	Sub             // 设备订阅校验
	Pub             // 设备发布校验
	UnSub
	Connected
	Closed
)

type PublishInfo struct {
	Topic   string
	Payload []byte
}

type MsgAckChan struct {
	Mu       sync.Mutex
	Id       int64
	IsClosed bool
	DataChan chan interface{} // auth ack, device sub ack, device pub ack
}

func (mac *MsgAckChan) tryCloseChan() {
	mac.Mu.Lock()
	defer mac.Mu.Unlock()
	if !mac.IsClosed {
		close(mac.DataChan)
		mac.IsClosed = true
	}
}

func (mac *MsgAckChan) trySendDataAndCloseChan(data interface{}) bool {
	mac.Mu.Lock()
	defer mac.Mu.Unlock()
	if !mac.IsClosed {
		mac.DataChan <- data
		close(mac.DataChan)
		mac.IsClosed = true
		return true
	}
	return false
}

type (
	// AsyncMsg 异步消息统一收发
	AsyncMsg struct {
		Id   int64
		Type int             // 1:连接鉴权,2:设备订阅校验,3:设备发布校验,4:unsub,6:closed
		Data json.RawMessage // auth ack sub ack pub ack
	}
)

type (
	AuthCheck struct {
		ClientId string
		Username string
		Password string
		Pass     bool
		Msg      string
	}
)

type PubTopic struct {
	ClientId string
	Username string
	Topic    string
	QoS      byte
	Retained bool
	Pass     bool
	Msg      string
}

// SubTopic 三方设备或服务订阅topic校验
type (
	SubTopic struct {
		Topic string
		QoS   byte
		Pass  bool
		Msg   string
	}
	SubTopics struct {
		ClientId string
		Username string
		Topics   []SubTopic
	}
)

type (
	// ConnectedNotify 三方设备或服务连接成功后通知对应驱动
	ConnectedNotify struct {
		ClientId string
		Username string
		IP       string
		Port     string
	}

	// ClosedNotify 三方设备或服务断开连接后通知对应驱动
	ClosedNotify struct {
		ClientId string
		Username string
	}

	UnSubNotify struct {
		ClientId string
		Username string
		Topics   []string
	}
)

type DriverClient struct {
	mu        sync.RWMutex
	ClientId  string
	Username  string
	PubTopic  string
	SubTopic  string
	ClientMap map[string]*ThirdClient // key is clientId
}

func (dc *DriverClient) AddThirdClient(client server.Client) {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	dc.ClientMap[client.ClientOptions().ClientID] = newThirdClient(client)
}

func (dc *DriverClient) DeleteThirdClient(clientId string) {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	delete(dc.ClientMap, clientId)
}

type ThirdClient struct {
	mu     sync.RWMutex
	client server.Client
	subs   map[string]struct{}
	pubs   map[string]struct{}
}

func (tc *ThirdClient) AddTopics(topics []string, t int) {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	if t == Sub {
		for i := range topics {
			tc.subs[topics[i]] = struct{}{}
		}
	} else if t == Pub {
		for i := range topics {
			tc.pubs[topics[i]] = struct{}{}
		}
	}
}

func (tc *ThirdClient) DeleteTopics(topics []string, t int) {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	if t == UnSub {
		for i := range topics {
			delete(tc.subs, topics[i])
		}
	}
}

func (tc *ThirdClient) CheckTopic(topic string, t int) bool {
	tc.mu.RLock()
	defer tc.mu.RUnlock()

	if t == Sub {
		_, ok := tc.subs[topic]
		return ok
	} else if t == Pub {
		_, ok := tc.pubs[topic]
		return ok
	}
	return false
}

func newThirdClient(c server.Client) *ThirdClient {
	return &ThirdClient{
		client: c,
		subs:   make(map[string]struct{}),
		pubs:   make(map[string]struct{}),
	}
}
