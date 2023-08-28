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

package ekuiperclient

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/winc-link/hummingbird/internal/dtos"
	"github.com/winc-link/hummingbird/internal/pkg/errort"
	"github.com/winc-link/hummingbird/internal/pkg/logger"
	"strings"

	"github.com/kirinlabs/HttpRequest"
)

const (
	ApiRuleInfoRoute    = "/rules/%s"
	ApiRuleStatusRoute  = "/rules/%s/status"
	ApiRuleCreateRoute  = "/rules"
	ApiRuleStopRoute    = "/rules/%s/stop"
	ApiRuleStartRoute   = "/rules/%s/start"
	ApiRuleRestartRoute = "/rules/%s/restart"
	ApiRuleDeleteRoute  = "/rules/%s"
	ApiRuleUpdateRoute  = "/rules/%s"
)

//Rule rulet2 is not found in registry

type ekuiperClient struct {
	baseUrl string
	lc      logger.LoggingClient
}

func New(baseUrl string, lc logger.LoggingClient) EkuiperClient {
	return &ekuiperClient{
		lc:      lc,
		baseUrl: baseUrl,
	}
}

//GetRule 获取规则
//func (c *ekuiperClient) GetRule(ctx context.Context, ruleId string) (dtos.GetRuleInfoResponse, error) {
//
//	var getRuleResponse dtos.GetRuleInfoResponse
//	err := httphelper.GetRequest(ctx, &getRuleResponse, c.baseUrl, fmt.Sprintf(ApiRuleInfoRoute, ruleId), nil)
//	if err != nil {
//		return dtos.GetRuleInfoResponse{}, err
//	}
//	return getRuleResponse, nil
//}

//RuleExist 规则是否存在
func (c *ekuiperClient) RuleExist(ctx context.Context, ruleId string) (bool, error) {

	req := HttpRequest.NewRequest()
	resp, err := req.Get(c.baseUrl + fmt.Sprintf(ApiRuleInfoRoute, ruleId))
	if err != nil {
		return false, err
	}
	if resp.StatusCode() == 200 {
		//body, err := resp.Body()
		//if err != nil {
		//	return false, err
		//}
		return true, nil
	} else if resp.StatusCode() == 400 {
		body, err := resp.Body()
		if err != nil {
			return false, err
		}
		return false, errort.NewCommonEdgeX(errort.AlertRuleNotExist, string(body), nil)
	} else if resp.StatusCode() == 404 {
		return false, nil
	}

	return false, errort.NewCommonEdgeX(errort.SystemErrorCode, "", nil)
}

//CreateRule 创建规则
func (c *ekuiperClient) CreateRule(ctx context.Context, actions []dtos.Actions, ruleId string, sql string) error {
	req := HttpRequest.NewRequest()
	var createRule dtos.CreateRule
	cr := createRule.BuildCreateRuleParam(actions, ruleId, sql)
	b, _ := json.Marshal(cr)
	resp, err := req.Post(c.baseUrl+ApiRuleCreateRoute, b)
	if err != nil {
		return err
	}
	if resp.StatusCode() == 201 {
		body, err := resp.Body()

		if err != nil {
			return err
		}
		if strings.Contains(string(body), "successfully") {
			return nil
		}
	} else if resp.StatusCode() == 400 {
		body, err := resp.Body()
		if err != nil {
			return err
		}
		return errort.NewCommonEdgeX(errort.InvalidRuleJson, string(body), nil)
	} else if resp.StatusCode() == 404 {
		return errort.NewCommonEdgeX(errort.EkuiperNotFindRule, "Rule engine not found rule", nil)
	}

	return errort.NewCommonEdgeX(errort.SystemErrorCode, "", nil)

}

//UpdateRule  更新规则
func (c *ekuiperClient) UpdateRule(ctx context.Context, actions []dtos.Actions, ruleId string, sql string) error {
	req := HttpRequest.NewRequest()
	var createRule dtos.CreateRule
	cr := createRule.BuildCreateRuleParam(actions, ruleId, sql)
	b, _ := json.Marshal(cr)
	resp, err := req.Put(c.baseUrl+fmt.Sprintf(ApiRuleUpdateRoute, ruleId), b)
	if err != nil {
		return err
	}
	if resp.StatusCode() == 200 {
		body, err := resp.Body()
		if err != nil {
			return err
		}
		if strings.Contains(string(body), "successfully") {
			return nil
		}
	} else if resp.StatusCode() == 400 {
		body, err := resp.Body()

		if err != nil {
			return err
		}
		return errort.NewCommonEdgeX(errort.InvalidRuleJson, string(body), nil)
	} else if resp.StatusCode() == 404 {
		return errort.NewCommonEdgeX(errort.EkuiperNotFindRule, "Rule engine not found rule", nil)
	}
	return errort.NewCommonEdgeX(errort.SystemErrorCode, "", nil)
}

//GetRuleStats 获取规则状态
func (c *ekuiperClient) GetRuleStats(ctx context.Context, ruleId string) (map[string]interface{}, error) {
	req := HttpRequest.NewRequest()
	resp, err := req.Get(c.baseUrl + fmt.Sprintf(ApiRuleStatusRoute, ruleId))
	ruleStatusResponse := make(map[string]interface{})
	if err != nil {
		return ruleStatusResponse, err
	}
	if resp.StatusCode() == 200 {
		body, err := resp.Body()
		if err != nil {
			return ruleStatusResponse, err
		}
		err = json.Unmarshal(body, &ruleStatusResponse)
		if err != nil {
			return ruleStatusResponse, errort.NewCommonEdgeX(errort.InvalidRuleJson, string(body), nil)
		}
		return ruleStatusResponse, nil
	}
	return ruleStatusResponse, errort.NewCommonEdgeX(errort.SystemErrorCode, "", nil)
}

//StartRule 启动规则
func (c *ekuiperClient) StartRule(ctx context.Context, ruleId string) error {
	req := HttpRequest.NewRequest()
	resp, err := req.Post(c.baseUrl + fmt.Sprintf(ApiRuleStartRoute, ruleId))
	if err != nil {
		return err
	}
	if resp.StatusCode() == 200 {
		body, err := resp.Body()
		if err != nil {
			return err
		}
		if strings.Contains(string(body), "started") {
			return nil
		}
	} else if resp.StatusCode() == 400 {
		body, err := resp.Body()
		if err != nil {
			return err
		}
		return errort.NewCommonEdgeX(errort.InvalidRuleJson, string(body), nil)
	} else if resp.StatusCode() == 404 {
		return errort.NewCommonEdgeX(errort.EkuiperNotFindRule, "Rule engine not found rule", nil)
	}
	return errort.NewCommonEdgeX(errort.SystemErrorCode, "", nil)
}

//StopRule 启动规则
func (c *ekuiperClient) StopRule(ctx context.Context, ruleId string) error {
	req := HttpRequest.NewRequest()
	resp, err := req.Post(c.baseUrl + fmt.Sprintf(ApiRuleStopRoute, ruleId))
	if err != nil {
		return err
	}
	if resp.StatusCode() == 200 {
		body, err := resp.Body()

		if err != nil {
			return err
		}
		if strings.Contains(string(body), "stopped") {
			return nil
		}
	} else if resp.StatusCode() == 400 {
		body, err := resp.Body()
		if err != nil {
			return err
		}
		return errort.NewCommonEdgeX(errort.InvalidRuleJson, string(body), nil)
	} else if resp.StatusCode() == 404 {
		return errort.NewCommonEdgeX(errort.EkuiperNotFindRule, "Rule engine not found rule", nil)
	}
	return errort.NewCommonEdgeX(errort.SystemErrorCode, "", nil)
}

//RestartRule 重启规则
func (c *ekuiperClient) RestartRule(ctx context.Context, ruleId string) error {
	req := HttpRequest.NewRequest()
	resp, err := req.Post(c.baseUrl + fmt.Sprintf(ApiRuleRestartRoute, ruleId))
	if err != nil {
		return err
	}
	if resp.StatusCode() == 200 {
		body, err := resp.Body()

		if err != nil {
			return err
		}
		if strings.Contains(string(body), "restarted") {
			return nil
		}
	} else if resp.StatusCode() == 400 {
		body, err := resp.Body()
		if err != nil {
			return err
		}
		return errort.NewCommonEdgeX(errort.InvalidRuleJson, string(body), nil)
	} else if resp.StatusCode() == 404 {
		return errort.NewCommonEdgeX(errort.EkuiperNotFindRule, "Rule engine not found rule", nil)
	}
	return errort.NewCommonEdgeX(errort.SystemErrorCode, "", nil)
}

//DeleteRule 删除规则
func (c *ekuiperClient) DeleteRule(ctx context.Context, ruleId string) error {
	req := HttpRequest.NewRequest()
	resp, err := req.Delete(c.baseUrl + fmt.Sprintf(ApiRuleDeleteRoute, ruleId))
	if err != nil {
		return err
	}
	c.lc.Infof("resp status code", resp.StatusCode())
	if resp.StatusCode() == 200 {
		body, err := resp.Body()

		if err != nil {
			return err
		}
		if strings.Contains(string(body), "dropped") {
			return nil
		}
	} else if resp.StatusCode() == 400 {
		body, err := resp.Body()
		if err != nil {
			return err
		}
		return errort.NewCommonEdgeX(errort.InvalidRuleJson, string(body), nil)
	} else if resp.StatusCode() == 404 {
		return nil
	}
	return errort.NewCommonEdgeX(errort.SystemErrorCode, "", nil)
}
