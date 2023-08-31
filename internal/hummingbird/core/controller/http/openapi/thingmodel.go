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

func (ctl *controller) OpenApiThingModelAddOrUpdate(c *gin.Context) {
	lc := ctl.lc
	var req dtos.OpenApiThingModelAddOrUpdateReq
	if err := c.ShouldBind(&req); err != nil {
		httphelper.RenderFail(c, errort.NewCommonErr(errort.DefaultReqParamsError, err), c.Writer, lc)
		return
	}
	edgeXErr := ctl.getThingModelApp().OpenApiAddThingModel(c, req)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(nil, c.Writer, lc)
}

func (ctl *controller) OpenApiThingModel(c *gin.Context) {
	lc := ctl.lc
	var req dtos.OpenApiQueryThingModelReq
	urlDecodeParam(&req, c.Request, lc)
	data, edgeXErr := ctl.getThingModelApp().OpenApiQueryThingModel(c, req.ProductId)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(data, c.Writer, lc)
}

func (ctl *controller) OpenApiDeleteThingModel(c *gin.Context) {
	lc := ctl.lc
	var req dtos.OpenApiThingModelDeleteReq
	if err := c.ShouldBind(&req); err != nil {
		httphelper.RenderFail(c, errort.NewCommonErr(errort.DefaultReqParamsError, err), c.Writer, lc)
		return
	}
	edgeXErr := ctl.getThingModelApp().OpenApiDeleteThingModel(c, req)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(nil, c.Writer, lc)
}

// OpenApiQueryDeviceEffectivePropertyData 查询设备实时属性
func (ctl *controller) OpenApiQueryDeviceEffectivePropertyData(c *gin.Context) {
	lc := ctl.lc
	var req dtos.DeviceEffectivePropertyDataReq
	urlDecodeParam(&req, c.Request, lc)
	data, err := ctl.getDeviceApp().DeviceEffectivePropertyData(req)
	if err != nil {
		httphelper.RenderFail(c, err, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(data, c.Writer, lc)
}

// OpenApiSetDeviceProperty 设备设备属性
func (ctl *controller) OpenApiSetDeviceProperty(c *gin.Context) {
	lc := ctl.lc
	var req dtos.OpenApiSetDeviceThingModel
	if err := c.ShouldBind(&req); err != nil {
		httphelper.RenderFail(c, errort.NewCommonErr(errort.DefaultReqParamsError, err), c.Writer, lc)
		return
	}
	err := ctl.getDeviceApp().SetDeviceProperty(req)
	if err != nil {
		httphelper.RenderFail(c, err, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(nil, c.Writer, lc)
}

func (ctl *controller) OpenApiInvokeThingService(c *gin.Context) {
	lc := ctl.lc
	var req dtos.InvokeDeviceServiceReq
	if err := c.ShouldBind(&req); err != nil {
		httphelper.RenderFail(c, errort.NewCommonErr(errort.DefaultReqParamsError, err), c.Writer, lc)
		return
	}
	data, err := ctl.getDeviceApp().DeviceInvokeThingService(req)
	if err != nil {
		httphelper.RenderFail(c, err, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(data, c.Writer, lc)
}

func (ctl *controller) OpenApiQueryDevicePropertyData(c *gin.Context) {
	lc := ctl.lc
	var req dtos.ThingModelPropertyDataRequest
	urlDecodeParam(&req, c.Request, lc)
	dtos.CorrectionPageParam(&req.BaseSearchConditionQuery)
	data, total, edgeXErr := ctl.getPersistApp().SearchDeviceThingModelHistoryPropertyData(req)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	pageResult := httphelper.NewPageResult(data, uint32(total), req.Page, req.PageSize)
	httphelper.ResultSuccess(pageResult, c.Writer, lc)
}

func (ctl *controller) OpenApiQueryDeviceEventData(c *gin.Context) {
	lc := ctl.lc
	var req dtos.ThingModelEventDataRequest
	urlDecodeParam(&req, c.Request, lc)
	dtos.CorrectionPageParam(&req.BaseSearchConditionQuery)
	data, total, edgeXErr := ctl.getPersistApp().SearchDeviceThingModelEventData(req)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	pageResult := httphelper.NewPageResult(data, uint32(total), req.Page, req.PageSize)
	httphelper.ResultSuccess(pageResult, c.Writer, lc)
}

func (ctl *controller) OpenApiQueryDeviceServiceData(c *gin.Context) {
	lc := ctl.lc
	var req dtos.ThingModelServiceDataRequest
	urlDecodeParam(&req, c.Request, lc)
	dtos.CorrectionPageParam(&req.BaseSearchConditionQuery)

	data, total, edgeXErr := ctl.getPersistApp().SearchDeviceThingModelServiceData(req)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	pageResult := httphelper.NewPageResult(data, uint32(total), req.Page, req.PageSize)
	httphelper.ResultSuccess(pageResult, c.Writer, lc)
}
