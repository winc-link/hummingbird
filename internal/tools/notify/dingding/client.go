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

package dingding

import (
	"bytes"
	"fmt"
	"github.com/kirinlabs/HttpRequest"
	"github.com/winc-link/hummingbird/internal/models"
	"github.com/winc-link/hummingbird/internal/pkg/constants"
	"github.com/winc-link/hummingbird/internal/pkg/logger"
	"strings"
)

type DingDingClient struct {
	lc logger.LoggingClient
}

func NewDingDingClient(lc logger.LoggingClient) *DingDingClient {
	return &DingDingClient{
		lc: lc,
	}
}

func dingDingTemplate(trigger constants.Trigger, deviceName string, ruleName string) string {
	// 字符串拼接
	var rst strings.Builder
	rst.WriteString("# 设备触发告警 \n\n")
	rst.WriteString(fmt.Sprintf("设备名称：%s \n\n", deviceName))
	rst.WriteString(fmt.Sprintf("告警规则名称：%s \n\n", ruleName))
	rst.WriteString(fmt.Sprintf("触发方式：%s \n\n", string(trigger)))
	return rst.String()
}

func (d *DingDingClient) Send(webhook string, rule models.Rule, device models.Device, product models.Product, ruleName string) {
	req := HttpRequest.NewRequest()
	req.JSON()
	resp, err := req.Post(webhook, bytes.NewBuffer([]byte(dingDingTemplate(rule.Trigger, device.Name, ruleName))))
	if err != nil {
		d.lc.Errorf("dingding send alert message error:", err.Error())
	}
	body, err := resp.Body()
	if err != nil {
		return
	}
	d.lc.Debug("dingding send message", string(body))
}
