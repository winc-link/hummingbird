/*******************************************************************************
 * Copyright 2017.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 *******************************************************************************/

package messageapp

import (
	"context"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/winc-link/hummingbird/internal/dtos"
	"github.com/winc-link/hummingbird/internal/hummingbird/core/container"
	"github.com/winc-link/hummingbird/internal/pkg/errort"
	pkgMQTT "github.com/winc-link/hummingbird/internal/tools/mqttclient"
)

func (msp *MessageApp) connectMQTT() (mqttClient pkgMQTT.MQTTClient) {
	lc := msp.lc
	var req dtos.NewMQTTClient
	var consumeCallback mqtt.MessageHandler
	var err error
	req, consumeCallback, err = msp.prepareMqttConnectParams()
	if err != nil {
		lc.Errorf("ConnectMQTT failed, err:%v", err)
		return
	}
	connF := func(ctx context.Context) {
		msp.lc.Info("ekuiper mqtt connect")
	}

	disConnF := func(ctx context.Context, msg dtos.CallbackMessage) {
		msp.lc.Info("ekuiper mqtt disconnect")
	}

	mqttClient, err = pkgMQTT.NewMQTTClient(req, lc, consumeCallback, connF, disConnF)
	if err != nil {
		err = errort.NewCommonErr(errort.MqttConnFail, err)
		lc.Errorf("ConnectMQTT failed, err:%v", err)
	}
	return mqttClient
}

func (tmq *MessageApp) prepareMqttConnectParams() (req dtos.NewMQTTClient, consumeCallback mqtt.MessageHandler, err error) {
	config := container.ConfigurationFrom(tmq.dic.Get)

	req = dtos.NewMQTTClient{
		Broker:   config.MessageQueue.URL(),
		ClientId: config.MessageQueue.Optional["ClientId"],
		Username: config.MessageQueue.Optional["Username"],
		Password: config.MessageQueue.Optional["Password"],
	}
	consumeCallback = tmq.ekuiperMsgHandle
	return
}

func (tmq *MessageApp) ekuiperMsgHandle(client mqtt.Client, message mqtt.Message) {

}

func (tmq *MessageApp) pushMsgToMessageBus(msg []byte) {
	config := container.ConfigurationFrom(tmq.dic.Get)
	tmq.ekuiperMqttClient.AsyncPublish(nil, config.MessageQueue.PublishTopicPrefix, msg, false)
}
