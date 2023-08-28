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

package agentclient

import (
	"context"
	"encoding/json"
	"github.com/winc-link/hummingbird/internal/dtos"
	"github.com/winc-link/hummingbird/internal/pkg/constants"
	"github.com/winc-link/hummingbird/internal/pkg/errort"
	"github.com/winc-link/hummingbird/internal/pkg/httphelper"
	"net/url"
)

const (
	ApiTerminalRoute = "/api/v1/terminal"
	ApiServicesRoute = "/api/v1/services"
)

type agentClient struct {
	baseUrl string
}

func New(baseUrl string) AgentClient {
	return &agentClient{
		baseUrl: baseUrl,
	}
}

func (c agentClient) GetGps(ctx context.Context) (dtos.Gps, error) {
	res := dtos.Gps{}
	commonRes := httphelper.CommonResponse{}
	err := httphelper.GetRequest(ctx, &commonRes, c.baseUrl, "/api/v1/gps", nil)
	if err != nil {
		return res, err
	}
	e := httphelper.CommResToSpecial(commonRes.Result, &res)
	if e != nil {
		return res, errort.NewCommonEdgeX(errort.DefaultSystemError, "res type is not gps", nil)
	}

	return res, nil
}

func (c agentClient) AddServiceMonitor(ctx context.Context, stats dtos.ServiceStats) error {
	commonRes := httphelper.CommonResponse{}
	return httphelper.PostRequest(ctx, &commonRes, c.baseUrl+"/api/v1/service/monitor", stats)
}

func (c *agentClient) DeleteServiceMonitor(ctx context.Context, serviceName string) error {
	commonRes := httphelper.CommonResponse{}
	var params = url.Values{}
	params.Set("service_name", serviceName)
	return httphelper.DeleteRequest(ctx, &commonRes, c.baseUrl, "/api/v1/service/monitor", params)
}

func (c *agentClient) GetAllDriverMonitor(ctx context.Context) ([]dtos.ServiceStats, error) {
	var list []dtos.ServiceStats
	commonRes := httphelper.CommonResponse{}
	err := httphelper.GetRequest(ctx, &commonRes, c.baseUrl, "/api/v1/driver/monitor", nil)
	if err != nil {
		return list, err
	}
	e := httphelper.CommResToSpecial(commonRes.Result, &list)
	if e != nil {
		return list, errort.NewCommonEdgeX(errort.DefaultSystemError, "res type is not gps", nil)
	}
	return list, nil
}

func (c *agentClient) RestartGateway(ctx context.Context) error {
	commonRes := httphelper.CommonResponse{}
	err := httphelper.PostRequest(ctx, &commonRes, c.baseUrl+"/api/v1/restart", nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *agentClient) OperationService(ctx context.Context, op dtos.Operation) error {
	commonRes := httphelper.CommonResponse{}
	err := httphelper.PostRequest(ctx, &commonRes, c.baseUrl+"/api/v1/operation", op)
	if err != nil {
		return err
	}

	return nil
}

func (c *agentClient) GetAllAppServiceMonitor(ctx context.Context) ([]dtos.ServiceStats, error) {
	var list []dtos.ServiceStats
	commonRes := httphelper.CommonResponse{}
	err := httphelper.GetRequest(ctx, &commonRes, c.baseUrl, "/api/v1/app_service/monitor", nil)
	if err != nil {
		return list, err
	}
	e := httphelper.CommResToSpecial(commonRes.Result, &list)
	if e != nil {
		return list, errort.NewCommonEdgeX(errort.DefaultSystemError, "res type is not gps", nil)
	}
	return list, nil
}

func (c agentClient) GetAllServices(ctx context.Context) (res dtos.ServicesStats, err error) {
	response := httphelper.CommonResponse{}
	err = httphelper.GetRequest(ctx, &response, c.baseUrl, ApiServicesRoute, nil)
	if err != nil {
		return res, err
	}
	resultByte, _ := json.Marshal(response.Result)
	errJson := json.Unmarshal(resultByte, &res)
	if errJson != nil {
		return res, errJson
	}
	return res, nil
}

func (c agentClient) Exec(ctx context.Context, req dtos.AgentRequest) (res dtos.AgentResponse, err error) {
	if req.TimeoutSeconds == 0 {
		req.TimeoutSeconds = constants.DefaultAgentReqTimeout
	}
	commonRes := httphelper.CommonResponse{}
	err = httphelper.PostRequest(ctx, &commonRes, c.baseUrl+ApiTerminalRoute, req)
	if err != nil {
		return res, err
	}

	resultByte, _ := json.Marshal(commonRes.Result)
	errJson := json.Unmarshal(resultByte, &res)
	if errJson != nil {
		return res, errJson
	}
	return res, nil
}

// 重置网关-并重启
func (c agentClient) ResetGateway(ctx context.Context) error {
	commonRes := httphelper.CommonResponse{}
	err := httphelper.PostRequest(ctx, &commonRes, c.baseUrl+"/api/v1/reset-gateway", nil)
	if err != nil {
		return err
	}
	return nil
}

//TODO: 重新构造一下Response，返回agent的报错信息
//func (c *agentClient) ExecuteOTAUpgrade(ctx context.Context, pid, version string) (httphelper.CommonResponse, error) {
//	upgradeReq := dtos.AgentOTAUpgradeReq{
//		Pid:            pid,
//		Version:        version,
//		TimeoutSeconds: constants.AgentDownloadOTAFirmTimeout,
//	}
//
//	commonRes := httphelper.CommonResponse{}
//	err := httphelper.PostRequest(ctx, &commonRes, c.baseUrl+"/api/v1/ota/upgrade", upgradeReq)
//	if err != nil {
//		return commonRes, err
//	}
//
//	return commonRes, nil
//}
