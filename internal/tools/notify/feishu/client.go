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

package feishu

import (
	"bytes"
	"github.com/kirinlabs/HttpRequest"
	"github.com/winc-link/hummingbird/internal/pkg/di"
	"github.com/winc-link/hummingbird/internal/pkg/logger"
)

type FeishuClient struct {
	lc logger.LoggingClient
	p  *di.Container
}

func NewFeishuClient(lc logger.LoggingClient, p *di.Container) *FeishuClient {
	return &FeishuClient{
		lc: lc,
		p:  p,
	}
}

func (d *FeishuClient) Send(webhook string, text string) {
	req := HttpRequest.NewRequest()
	req.JSON()
	resp, err := req.Post(webhook, bytes.NewBuffer([]byte(text)))
	if err != nil {
		d.lc.Errorf("feishu send alert message error:", err.Error())
	}
	body, err := resp.Body()
	if err != nil {
		return
	}
	d.lc.Debug("feishu send message", string(body))
}
