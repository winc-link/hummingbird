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

package openapi

import (
	"github.com/gin-gonic/gin"
	"github.com/winc-link/hummingbird/internal/dtos"
	"github.com/winc-link/hummingbird/internal/pkg/errort"
	"github.com/winc-link/hummingbird/internal/pkg/httphelper"
)

func (ctl *controller) OpenApiCreateDevice(c *gin.Context) {
	lc := ctl.lc
	var req dtos.DeviceAddRequest
	if err := c.ShouldBind(&req); err != nil {
		httphelper.RenderFail(c, errort.NewCommonErr(errort.DefaultReqParamsError, err), c.Writer, lc)
		return
	}
	_, edgeXErr := ctl.getDeviceApp().AddDevice(c, req)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(nil, c.Writer, lc)
}

func (ctl *controller) OpenApiUpdateDevice(c *gin.Context) {
	lc := ctl.lc
	id := c.Param(UrlParamDeviceId)
	var req dtos.DeviceUpdateRequest
	req.Id = id
	if err := c.ShouldBind(&req); err != nil {
		httphelper.RenderFail(c, errort.NewCommonErr(errort.DefaultReqParamsError, err), c.Writer, lc)
		return
	}
	err := ctl.getDeviceApp().DeviceUpdate(c, req)
	if err != nil {
		httphelper.RenderFail(c, err, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(nil, c.Writer, lc)
}

func (ctl *controller) OpenApiDeviceSearch(c *gin.Context) {
	lc := ctl.lc
	var req dtos.DeviceSearchQueryRequest
	urlDecodeParam(&req, c.Request, lc)
	dtos.CorrectionPageParam(&req.BaseSearchConditionQuery)
	data, total, edgeXErr := ctl.getDeviceApp().OpenApiDevicesSearch(c, req)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	pageResult := httphelper.NewPageResult(data, total, req.Page, req.PageSize)
	httphelper.ResultSuccess(pageResult, c.Writer, lc)
}

func (ctl *controller) OpenApiDeviceById(c *gin.Context) {
	lc := ctl.lc
	id := c.Param(UrlParamDeviceId)
	data, edgeXErr := ctl.getDeviceApp().OpenApiDeviceById(c, id)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(data, c.Writer, lc)
}

func (ctl *controller) OpenApiDeleteDevice(c *gin.Context) {
	lc := ctl.lc
	id := c.Param(UrlParamDeviceId)
	edgeXErr := ctl.getDeviceApp().DeleteDeviceById(c, id)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(nil, c.Writer, lc)
}

func (ctl *controller) OpenApiDeviceStatus(c *gin.Context) {
	lc := ctl.lc
	id := c.Param(UrlParamDeviceId)
	data, edgeXErr := ctl.getDeviceApp().OpenApiDeviceStatusById(c, id)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(data, c.Writer, lc)
}
