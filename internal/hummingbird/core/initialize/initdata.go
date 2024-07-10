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

package initialize

import (
	"context"
	"encoding/json"
	"github.com/kirinlabs/HttpRequest"
	"github.com/winc-link/hummingbird/internal/hummingbird/core/config"
	pkgContainer "github.com/winc-link/hummingbird/internal/pkg/container"
	"github.com/winc-link/hummingbird/internal/pkg/di"
	"github.com/winc-link/hummingbird/internal/pkg/logger"
	"strings"
	"sync"
	"time"
)

func initApp(ctx context.Context, configuration *config.ConfigurationStruct, dic *di.Container) bool {
	lc := pkgContainer.LoggingClientFrom(dic.Get)
	go initEkuiperStreams(dic, lc, configuration)
	return true
}

func initEkuiperStreams(dic *di.Container, lc logger.LoggingClient, configuration *config.ConfigurationStruct) {
	// time 10s 以保证ekuiper初始化完成
	time.Sleep(10 * time.Second)
	req := HttpRequest.NewRequest()
	r := make(map[string]string)
	r["sql"] = "CREATE STREAM mqtt_stream () WITH (DATASOURCE=\"eventbus/in\", FORMAT=\"JSON\",SHARED = \"true\")"
	b, _ := json.Marshal(r)
	url := configuration.Clients["Ekuiper"].Address() + "/streams"
	resp, err := req.Post(url, b)
	if err != nil {
		lc.Errorf("init ekuiper stream failed error:%+v", err.Error())
		return
	}
	lc.Infof("init ekuiper stream start")
	lc.Info("ekuiper stream start resp code", resp.StatusCode())

	if resp.StatusCode() == 201 {
		body, err := resp.Body()
		if err != nil {
			lc.Errorf("init ekuiper stream failed error:%+v", err.Error())
			return
		}
		if strings.Contains(string(body), "created") {
			lc.Infof("init ekuiper stream success")
			return
		}
	} else if resp.StatusCode() == 400 {
		body, err := resp.Body()
		lc.Infof("init ekuiper stream body", string(body))
		if err != nil {
			lc.Errorf("init ekuiper stream failed error:%+v", err.Error())
			return
		}

		if strings.Contains(string(body), "already exists") {
			lc.Infof("init ekuiper stream plug success")
			return
		}
	} else {
		lc.Errorf("init ekuiper stream failed resp code:%+v", resp.StatusCode())
	}
}

func DownloadEkuiperKafkaPlug(wg *sync.WaitGroup, dic *di.Container, lc logger.LoggingClient) {
	defer wg.Done()
	req := HttpRequest.NewRequest()
	kafkaPlug := make(map[string]string)
	kafkaPlug["name"] = "kafka"
	kafkaPlug["file"] = "https://packages.emqx.net/kuiper-plugins/1.10.0/debian/sinks/kafka_amd64.zip"
	b, _ := json.Marshal(kafkaPlug)
	resp, err := req.Post("http://ekuiper:9081/plugins/sinks", b)
	if err != nil {
		lc.Errorf("down ekuiper kafka plug failed error:%+v", err.Error())
	}
	if resp.StatusCode() == 201 {
		body, err := resp.Body()
		if err != nil {
			lc.Errorf("down ekuiper kafka plug failed error:%+v", err.Error())
			return
		}
		if strings.Contains(string(body), "created") {
			lc.Infof("down ekuiper kafka plug success")
			return
		}
	} else if resp.StatusCode() == 400 {
		body, err := resp.Body()
		if err != nil {
			lc.Errorf("down ekuiper kafka plug failed error:%+v", err.Error())
			return
		}
		if strings.Contains(string(body), "duplicate") {
			lc.Infof("down ekuiper kafka plug success")
			return
		}
	} else {
		lc.Errorf("down ekuiper kafka plug failed resp code:%+v", resp.StatusCode())
	}
}

func DownloadEkuiperTdenginePlug(wg *sync.WaitGroup, dic *di.Container, lc logger.LoggingClient) {
	defer wg.Done()
	req := HttpRequest.NewRequest()
	kafkaPlug := make(map[string]interface{})
	kafkaPlug["name"] = "tdengine"
	kafkaPlug["file"] = "https://packages.emqx.net/kuiper-plugins/1.10.0/debian/sinks/tdengine_amd64.zip"
	kafkaPlug["shellParas"] = []string{"2.4.0.26"}

	b, _ := json.Marshal(kafkaPlug)
	resp, err := req.Post("http://ekuiper:9081/plugins/sinks", b)
	if err != nil {
		lc.Errorf("down ekuiper tdengine plug failed error:%+v", err.Error())
	}
	if resp.StatusCode() == 201 {
		body, err := resp.Body()
		if err != nil {
			lc.Errorf("down ekuiper tdengine plug failed error:%+v", err.Error())
			return
		}
		if strings.Contains(string(body), "created") {
			lc.Infof("down ekuiper tdengine plug success")
			return
		}
	} else if resp.StatusCode() == 400 {
		body, err := resp.Body()
		if err != nil {
			lc.Errorf("down ekuiper tdengine plug failed error:%+v", err.Error())
			return
		}
		if strings.Contains(string(body), "duplicate") {
			lc.Infof("down ekuiper tdengine plug success")
			return
		}
	} else {
		lc.Errorf("down ekuiper tdengine plug failed resp code:%+v", resp.StatusCode())
	}
}
