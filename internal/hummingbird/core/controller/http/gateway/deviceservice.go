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
package gateway

import (
	"github.com/gin-gonic/gin"
	"github.com/winc-link/hummingbird/internal/dtos"
	"github.com/winc-link/hummingbird/internal/pkg/errort"
	"github.com/winc-link/hummingbird/internal/pkg/httphelper"
)

// @Tags    驱动实例管理
// @Summary 查询驱动实例
// @Produce json
// @Param   request query    dtos.DeviceServiceSearchQueryRequest true "参数"
// @Success 200     {object} httphelper.ResPageResult
// @Router  /api/v1/device-servers [get]
func (ctl *controller) DeviceServicesSearch(c *gin.Context) {
	lc := ctl.lc
	var req dtos.DeviceServiceSearchQueryRequest
	urlDecodeParam(&req, c.Request, lc)
	dtos.CorrectionPageParam(&req.BaseSearchConditionQuery)
	if err := c.ShouldBind(&req); err != nil {
		httphelper.RenderFail(c, errort.NewCommonErr(errort.DefaultReqParamsError, err), c.Writer, lc)
		return
	}

	// TODO 驱动实例搜索是否需要查询驱动库
	dss, total, err := ctl.getDriverServiceApp().Search(c, req)
	if err != nil {
		httphelper.RenderFail(c, err, c.Writer, lc)
		return
	}
	data := make([]dtos.DeviceServiceResponse, len(dss))
	for i, ds := range dss {
		dl, err := ctl.getDriverLibApp().DriverLibById(ds.DeviceLibraryId)
		if err != nil {
			httphelper.RenderFail(c, err, c.Writer, lc)
			return
		}
		data[i] = dtos.DeviceServiceResponseFromModel(ds, dl)
	}

	pageResult := httphelper.NewPageResult(data, total, req.Page, req.PageSize)
	httphelper.ResultSuccess(pageResult, c.Writer, lc)
}

// @Tags    驱动实例管理
// @Summary 删除驱动实例（废弃，已经改为websockert形式)
// @Produce json
// @Param   deviceServiceId path     string true "驱动实例 ID"
// @Success 200             {object} httphelper.CommonResponse
// @Router  /api/v1/device_server/:deviceServiceId [delete]
func (ctl *controller) DeviceServiceDelete(c *gin.Context) {
	lc := ctl.lc
	id := c.Param(UrlParamDeviceServiceId)

	err := ctl.getDriverServiceApp().Del(c, id)
	//edgeXErr := gatewayapp.DeviceServiceDelete(c, id)
	if err != nil {
		httphelper.RenderFail(c, err, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(nil, c.Writer, lc)
}

// @Tags    驱动实例管理
// @Summary 编辑驱动实例
// @Produce json
// @Param  request body dtos.DeviceServiceUpdateRequest true "参数"
// @Success 200 {object} httphelper.CommonResponse
// @Router  /api/v1/device-server [put]
func (ctl *controller) DeviceServiceUpdate(c *gin.Context) {
	lc := ctl.lc
	var req dtos.DeviceServiceUpdateRequest
	if err := c.ShouldBind(&req); err != nil {
		httphelper.RenderFail(c, errort.NewCommonErr(errort.DefaultReqParamsError, err), c.Writer, lc)
		return
	}
	err := ctl.getDriverServiceApp().Update(c, req)
	if err != nil {
		httphelper.RenderFail(c, err, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(nil, c.Writer, lc)
}
