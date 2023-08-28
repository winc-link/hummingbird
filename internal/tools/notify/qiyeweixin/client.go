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

package yiqiweixin

import (
	"encoding/json"
	"github.com/kirinlabs/HttpRequest"
	"github.com/winc-link/hummingbird/internal/pkg/di"
	"github.com/winc-link/hummingbird/internal/pkg/logger"
)

type WeixinClient struct {
	lc logger.LoggingClient
	p  *di.Container
}

func NewWeiXinClient(lc logger.LoggingClient, p *di.Container) *WeixinClient {
	return &WeixinClient{
		lc: lc,
		p:  p,
	}
}

type QiYeWeiXinTemplate struct {
	Msgtype  string   `json:"msgtype"`
	Markdown Markdown `json:"markdown"`
}

type Markdown struct {
	Content string `json:"content"`
}

func (d *WeixinClient) Send(webhook string, text string) {
	if webhook == "" {
		return
	}
	req := HttpRequest.NewRequest()
	req.JSON()
	context, _ := json.Marshal(text)
	resp, err := req.Post(webhook, context)
	if err != nil {
		d.lc.Errorf("weixin send alert message error:", err.Error())
	}
	body, err := resp.Body()
	if err != nil {
		return
	}
	d.lc.Debug("weixin send message", string(body))
}
