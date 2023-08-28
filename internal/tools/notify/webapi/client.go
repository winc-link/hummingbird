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

package webapi

import (
	"encoding/json"
	"github.com/kirinlabs/HttpRequest"
	"github.com/winc-link/hummingbird/internal/models"
	"github.com/winc-link/hummingbird/internal/pkg/di"
	"github.com/winc-link/hummingbird/internal/pkg/logger"
	"time"
)

type WebApiClient struct {
	lc logger.LoggingClient
	p  *di.Container
}

type WebApiTemplate struct {
	Product struct {
		ProductId   string `json:"product_id"`
		ProductName string `json:"product_name"`
	} `json:"product"`
	Device struct {
		DeviceId   string `json:"device_id"`
		DeviceName string `json:"device_name"`
	} `json:"device"`
	Rule struct {
		RuleId      string `json:"rule_id"`
		RuleName    string `json:"rule_name"`
		AlertLevel  string `json:"alert_level"`
		Trigger     string `json:"trigger,omitempty"`
		TriggerTime int64  `json:"trigger_time,omitempty"`
	} `json:"rule"`
	Message string `json:"message"`
}

func (d *WebApiClient) generateWebApiTemplate(rule models.AlertRule, device models.Device, product models.Product, message map[string]interface{}) WebApiTemplate {
	var temp WebApiTemplate
	msg, _ := json.Marshal(message)
	temp.Message = string(msg)
	temp.Product.ProductId = product.Id
	temp.Product.ProductName = product.Name
	temp.Device.DeviceId = device.Id
	temp.Device.DeviceName = device.Name
	temp.Rule.RuleId = rule.Id
	temp.Rule.RuleName = rule.Name
	temp.Rule.AlertLevel = string(rule.AlertLevel)
	if len(rule.SubRule) > 0 {
		temp.Rule.Trigger = string(rule.SubRule[0].Trigger)
		temp.Rule.TriggerTime = time.Now().UnixMilli()
	}

	return temp
}

func NewWebApiClient(lc logger.LoggingClient, p *di.Container) *WebApiClient {
	return &WebApiClient{
		lc: lc,
		p:  p,
	}
}

func (d *WebApiClient) Send(webhook string, header []map[string]string, rule models.AlertRule, device models.Device, product models.Product, messages map[string]interface{}) {
	if webhook == "" {
		return
	}
	req := HttpRequest.NewRequest()
	req.JSON()
	d.lc.Infof("webapi send header:", header)
	context, _ := json.Marshal(d.generateWebApiTemplate(rule, device, product, messages))
	for _, m := range header {
		req.SetHeaders(m)
	}
	_, err := req.Post(webhook, context)
	if err != nil {
		d.lc.Errorf("webapi send alert message error:", err.Error())
	}
	d.lc.Info("webapi send message")
}
