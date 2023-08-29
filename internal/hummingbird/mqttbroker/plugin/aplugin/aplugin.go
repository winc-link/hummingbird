package aplugin

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/winc-link/hummingbird/internal/hummingbird/mqttbroker"

	"github.com/winc-link/hummingbird/internal/hummingbird/mqttbroker/plugin/aplugin/snowflake"

	"github.com/winc-link/hummingbird/internal/hummingbird/mqttbroker/config"
	"github.com/winc-link/hummingbird/internal/hummingbird/mqttbroker/server"
	"go.uber.org/zap"
)

var _ server.Plugin = (*APlugin)(nil)

const Name = "aplugin"

func init() {
	server.RegisterPlugin(Name, New)
	config.RegisterDefaultPluginConfig(Name, &DefaultConfig)
}

func New(config config.Config) (server.Plugin, error) {
	return newAPlugin()
}

type APlugin struct {
	mu            sync.Mutex
	ctx           context.Context
	cancel        context.CancelFunc
	wg            *sync.WaitGroup
	node          *snowflake.Node
	log           *zap.SugaredLogger
	ackMap        sync.Map                 // async ack
	driverClients map[string]*DriverClient // driver map, key is username
	publisher     server.Publisher
	publishChan   chan PublishInfo
}

func newAPlugin() (*APlugin, error) {
	node, err := snowflake.NewNode(1)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &APlugin{
		node:          node,
		ctx:           ctx,
		cancel:        cancel,
		wg:            &sync.WaitGroup{},
		log:           zap.NewNop().Sugar(),
		driverClients: make(map[string]*DriverClient),
		publishChan:   make(chan PublishInfo, 32),
	}, nil
}

func (t *APlugin) Load(service server.Server) error {
	t.log = server.LoggerWithField(zap.String("plugin", Name)).Sugar()
	t.publisher = service.Publisher()
	t.wg.Add(1)
	go func() {
		defer t.wg.Done()
		for {
			select {
			case <-t.ctx.Done():
				t.log.Infof("plugin(%s) exit", t.Name())
				return
			case msg := <-t.publishChan:
				t.log.Infof("publish msg: topic: %s, payload: %s", msg.Topic, string(msg.Payload))
				t.publisher.Publish(&mqttbroker.Message{
					QoS:     1,
					Topic:   msg.Topic,
					Payload: msg.Payload,
				})
			}
		}
	}()
	return nil
}

func (t *APlugin) Unload() error {
	t.cancel()
	t.wg.Wait()
	return nil
}

func (t *APlugin) Name() string {
	return Name
}

func (t *APlugin) genAckChan(id int64) *MsgAckChan {
	ack := &MsgAckChan{
		Id:       id,
		DataChan: make(chan interface{}, 1),
	}
	t.ackMap.Store(id, ack)
	return ack
}

func (t *APlugin) publishWithAckMsg(id int64, topic string, tp int, msg interface{}) (*MsgAckChan, error) {
	payload, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	buff, err := json.Marshal(AsyncMsg{
		Id:   id,
		Type: tp,
		Data: payload,
	})
	if err != nil {
		return nil, err
	}
	ackChan := t.genAckChan(id)
	select {
	case <-time.After(time.Second):
		t.ackMap.Delete(id)
		return nil, errors.New("send auth msg to publish chan timeout")
	case t.publishChan <- PublishInfo{
		Topic:   topic,
		Payload: buff,
	}:
		return ackChan, nil
	}
}

func (t *APlugin) publishNotifyMsg(id int64, topic string, tp int, msg interface{}) error {
	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	buff, err := json.Marshal(AsyncMsg{
		Id:   id,
		Type: tp,
		Data: payload,
	})
	if err != nil {
		return err
	}
	select {
	case <-time.After(time.Second):
		return errors.New("send auth msg to publish chan timeout")
	case t.publishChan <- PublishInfo{
		Topic:   topic,
		Payload: buff,
	}:
		return nil
	}
}

func (t *APlugin) validate(username, password string) error {
	//t.log.Debugf("got clientId: %s, username: %s, password: %s", clientId, username, password)
	passwd, err := md5GenPasswd(username)
	if err != nil {
		return err
	}
	if passwd[8:24] != password {
		return errors.New("auth failure")
	}
	return nil
}

func md5GenPasswd(username string) (string, error) {
	h := md5.New()
	_, err := h.Write([]byte(username))
	if err != nil {
		return "", err
	}
	rs := h.Sum(nil)
	return hex.EncodeToString(rs), nil
}
