/*******************************************************************************
 * Copyright 2017 Dell Inc.
 * Copyright (c) 2019 Intel Corporation
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
package websocket

import (
	"encoding/json"
	"github.com/gin-gonic/gin/binding"
	"github.com/winc-link/hummingbird/internal/dtos"
	"github.com/winc-link/hummingbird/internal/hummingbird/core/container"
	interfaces "github.com/winc-link/hummingbird/internal/hummingbird/core/interface"
	"github.com/winc-link/hummingbird/internal/pkg/constants"
	"github.com/winc-link/hummingbird/internal/pkg/errort"
	"github.com/winc-link/hummingbird/internal/pkg/httphelper"
	"github.com/winc-link/hummingbird/internal/pkg/i18n"
	"time"
)

/**
receive: {"code":10004,"data":{"id":"123291"}}
*/
func DeviceServiceLog(c *wsClient, data interface{}, code dtos.WsCode) {
	var req dtos.DeviceServiceRunLogRequest
	bytes, _ := json.Marshal(data)
	err := binding.JSON.BindBody(bytes, &req)
	if err != nil {
		c.lc.Error(err.Error())
		c.sendData(
			code,
			httphelper.WsResultFail(
				errort.DefaultReqParamsError,
				i18n.TransCode(c.ctx, errort.DefaultReqParamsError, nil),
			),
		)
		return
	}

	driverServerApp := container.DriverServiceAppFrom(c.dic.Get)

	ds, err := driverServerApp.Get(c.ctx, req.Id)
	if err != nil {
		c.lc.Error(err.Error())
		c.sendData(
			code,
			httphelper.WsResultFail(
				errort.DefaultReqParamsError,
				i18n.TransCode(c.ctx, errort.DefaultReqParamsError, nil),
			),
		)
		return
	}

	driverLibApp := container.DriverAppFrom(c.dic.Get)
	dl, err := driverLibApp.DriverLibById(ds.DeviceLibraryId)
	if err != nil {
		c.lc.Error(err.Error())
		c.sendData(
			code,
			httphelper.WsResultFail(
				errort.DefaultReqParamsError,
				i18n.TransCode(c.ctx, errort.DefaultReqParamsError, nil),
			),
		)
		return
	}
	logfilePath := interfaces.DMIFrom(c.dic.Get).GetDriverInstanceLogPath(dl.ContainerName)

	if req.Operate == constants.StatusRead {
		//读取日志
		StartReadServiceLog(c, req.Id, logfilePath, code)

	} else if req.Operate == constants.StatusStop {
		//停止
		StopReadServiceLog(c, req.Id, code)
		errCode := errort.DefaultSuccess
		errMsg := ""
		successMsg := ""
		status := i18n.DefaultSuccess

		msg := i18n.Trans(i18n.GetLang(c.ctx), i18n.CloudInstanceLogResp, map[string]interface{}{
			"name":   ds.Name,
			"status": i18n.Trans(i18n.GetLang(c.ctx), status, nil),
		})
		successMsg = msg
		c.sendData(code, httphelper.WsResult(errCode, "", errMsg, successMsg))
	}

}

func StartReadServiceLog(c *wsClient, serviceId, logFilePath string, code dtos.WsCode) {
	driverServerLogApp := container.HpcServiceAppFrom(c.dic.Get)
	hpc := driverServerLogApp.Add(serviceId, logFilePath)

	tails, err := hpc.Read()
	if err != nil {
		//报错
	}
	for {
		line, ok := <-tails //遍历chan，读取日志内容
		if !ok {
			c.lc.Info("stop")
			return
		}
		c.lc.Infof("msg %+v", line.Text)
		c.sendData(code, httphelper.WsResult(errort.DefaultSuccess, line.Text, "", ""))
		//return
	}
}

func StopReadServiceLog(c *wsClient, serviceId string, code dtos.WsCode) {
	driverServerLogApp := container.HpcServiceAppFrom(c.dic.Get)
	hpc := driverServerLogApp.Get(serviceId)
	if hpc == nil {
		//报错
		c.sendData(
			code,
			httphelper.WsResultFail(
				errort.DefaultReqParamsError,
				i18n.TransCode(c.ctx, errort.DefaultReqParamsError, nil),
			),
		)
		return
	}
	hpc.Stop()
}

/**
receive: {"code":10003,"data":{"id":"3208327514"}}
*/
func DeviceLibraryDelete(c *wsClient, data interface{}, code dtos.WsCode) {
	var req dtos.DeviceServiceDeleteRequest
	bytes, _ := json.Marshal(data)
	err := binding.JSON.BindBody(bytes, &req)
	if err != nil {
		c.lc.Error(err.Error())
		c.sendData(
			code,
			httphelper.WsResultFail(
				errort.DefaultReqParamsError,
				i18n.TransCode(c.ctx, errort.DefaultReqParamsError, nil),
			),
		)
		return
	}
	driverServiceApp := container.DriverServiceAppFrom(c.dic.Get)
	// 响应
	dsName := ""

	ds, err := driverServiceApp.Get(c.ctx, req.Id)
	if err != nil {
		c.sendData(
			code,
			httphelper.WsResultFail(
				errort.DefaultReqParamsError,
				i18n.TransCode(c.ctx, errort.DefaultReqParamsError, nil),
			),
		)
		return
	} else {
		dsName = ds.Name
	}

	errCode := errort.DefaultSuccess
	err = driverServiceApp.Del(c.ctx, req.Id)
	if err != nil {
		c.lc.Errorf("del cloud service err: %+v", err)
		errCode = errort.NewCommonEdgeXWrapper(err).Code()
	}

	isSuccess := true
	errMsg := ""
	successMsg := ""
	status := i18n.DefaultSuccess

	if errCode != errort.DefaultSuccess {
		errMsg = i18n.TransCode(c.ctx, errCode, nil)
		isSuccess = false
		status = i18n.DefaultFail
	}
	msg := i18n.Trans(i18n.GetLang(c.ctx), i18n.AppServiceDeleteResp, map[string]interface{}{
		"name":   dsName,
		"status": i18n.Trans(i18n.GetLang(c.ctx), status, nil),
	})
	if isSuccess {
		successMsg = msg
	} else {
		errMsg = msg + ": " + errMsg
	}
	time.Sleep(2 * time.Second)
	c.sendData(code, httphelper.WsResult(errCode, "", errMsg, successMsg))
}
