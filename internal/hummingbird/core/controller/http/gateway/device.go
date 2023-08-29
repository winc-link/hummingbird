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
	"github.com/winc-link/hummingbird/internal/pkg/constants"
	"github.com/winc-link/hummingbird/internal/pkg/errort"
	"github.com/winc-link/hummingbird/internal/pkg/httphelper"
)

// @Tags    设备管理
// @Summary 查询设备列表
// @Produce json
// @Param   request query   dtos.DeviceSearchQueryRequest true "参数"
// @Success 200     {array} []dtos.DeviceSearchQueryResponse
// @Router  /api/v1/devices [get]
func (ctl *controller) DevicesSearch(c *gin.Context) {
	lc := ctl.lc
	var req dtos.DeviceSearchQueryRequest
	urlDecodeParam(&req, c.Request, lc)
	dtos.CorrectionPageParam(&req.BaseSearchConditionQuery)
	data, total, edgeXErr := ctl.getDeviceApp().DevicesSearch(c, req)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	pageResult := httphelper.NewPageResult(data, total, req.Page, req.PageSize)
	httphelper.ResultSuccess(pageResult, c.Writer, lc)
}

// @Tags    设备管理
// @Summary 查询详情
// @Produce json
// @Param   deviceId path     string true "pid"
// @Success 200  {object} dtos.DeviceInfoResponse
// @Router  /api/v1/device/:deviceId [get]
func (ctl *controller) DeviceById(c *gin.Context) {
	lc := ctl.lc
	id := c.Param(UrlParamDeviceId)
	data, edgeXErr := ctl.getDeviceApp().DeviceById(c, id)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(data, c.Writer, lc)
}

// @Tags    设备管理
// @Summary 删除设备
// @Produce json
// @Param   deviceId path     string true "pid"
// @Success 200  {object} httphelper.CommonResponse
// @Router  /api/v1/device/:deviceId [delete]
func (ctl *controller) DeviceDelete(c *gin.Context) {
	lc := ctl.lc
	id := c.Param(UrlParamDeviceId)
	edgeXErr := ctl.getDeviceApp().DeleteDeviceById(c, id)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(nil, c.Writer, lc)
}

// @Tags    设备管理
// @Summary 批量删除设备
// @Produce json
// @Param   deviceId path     string true "pid"
// @Success 200  {object} httphelper.CommonResponse
// @Router  /api/v1/devices [delete]
func (ctl *controller) DevicesDelete(c *gin.Context) {
	lc := ctl.lc
	var req dtos.DeviceBatchDelete
	if err := c.ShouldBind(&req); err != nil {
		httphelper.RenderFail(c, errort.NewCommonErr(errort.DefaultReqParamsError, err), c.Writer, lc)
		return
	}
	if len(req.DeviceIds) == 0 {
		httphelper.RenderFail(c, errort.NewCommonErr(errort.DefaultReqParamsError, nil), c.Writer, lc)
		return
	}
	edgeXErr := ctl.getDeviceApp().BatchDeleteDevice(c, req.DeviceIds)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(nil, c.Writer, lc)
}

// @Tags    设备管理
// @Summary 添加设备
// @Produce json
// @Param   request query    dtos.DeviceAddRequest true "参数"
// @Success 200  {object}  httphelper.CommonResponse
// @Router  /api/v1/device [post]
func (ctl *controller) DeviceByAdd(c *gin.Context) {
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

// @Tags    设备管理
// @Summary 查询mqtt连接详情
// @Produce json
// @Param   deviceId path   string true "pid"
// @Success 200  {object} dtos.DeviceAuthInfoResponse
// @Router  /api/v1/device-mqtt/:deviceId [get]
func (ctl *controller) DeviceMqttInfoById(c *gin.Context) {
	lc := ctl.lc
	id := c.Param(UrlParamDeviceId)
	data, edgeXErr := ctl.getDeviceApp().DeviceMqttAuthInfo(c, id)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(data, c.Writer, lc)
}

func (ctl *controller) AddMqttAuth(c *gin.Context) {
	lc := ctl.lc
	var req dtos.AddMqttAuthInfoRequest
	if err := c.ShouldBind(&req); err != nil {
		httphelper.RenderFail(c, errort.NewCommonErr(errort.DefaultReqParamsError, err), c.Writer, lc)
		return
	}
	data, edgeXErr := ctl.getDeviceApp().AddMqttAuth(c, req)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(data, c.Writer, lc)
}

// @Tags 设备管理
// @Summary 设备导入模版下载
// @Produce json
// @Param req query dtos.DeviceImportTemplateRequest true "参数"
// @Success 200 {object} string
// @Router /api/v1/devices/import-template [get]
func (ctl *controller) DeviceImportTemplateDownload(c *gin.Context) {
	lc := ctl.lc
	var req dtos.DeviceImportTemplateRequest

	if err := c.ShouldBind(&req); err != nil {
		httphelper.RenderFail(c, errort.NewCommonErr(errort.DefaultReqParamsError, err), c.Writer, lc)
		return
	}
	file, edgeXErr := ctl.getDeviceApp().DeviceImportTemplateDownload(c, req)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	data, _ := file.Excel.WriteToBuffer()
	httphelper.ResultExcelData(c, file.FileName, data)
}

// @Tags 设备管理
// @Summary 设备导入模版校验
// @Produce json
// @Success 200 {object} string
// @Router /api/v1/device/upload-validated [post]
func (ctl *controller) UploadValidated(c *gin.Context) {
	lc := ctl.lc
	files, _ := c.FormFile("file")
	f, err := files.Open()
	if err != nil {
		err = errort.NewCommonErr(errort.DefaultUploadFileErrorCode, err)
		httphelper.RenderFail(c, err, c.Writer, lc)
		return
	}
	file, edgeXErr := dtos.NewImportFile(f)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	edgeXErr = ctl.getDeviceApp().UploadValidated(c, file)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}

	httphelper.ResultSuccess(nil, c.Writer, lc)
}

// @Tags 设备管理
// @Summary 设备导入
// @Produce json
// @Success 200 {object} httphelper.CommonResponse
// @Router /api/v1/devices/import [post]
func (ctl *controller) DevicesImport(c *gin.Context) {
	lc := ctl.lc
	//productId := c.Param(UrlParamProductId)
	//cloudInstanceId := c.Param(UrlParamCloudInstanceId)
	//var req dtos.ProductSearchQueryRequest
	var req dtos.DevicesImport
	urlDecodeParam(&req, c.Request, lc)

	files, _ := c.FormFile("file")
	f, err := files.Open()
	if err != nil {
		err = errort.NewCommonErr(errort.DefaultUploadFileErrorCode, err)
		httphelper.RenderFail(c, err, c.Writer, lc)
		return
	}
	file, edgeXErr := dtos.NewImportFile(f)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	result, edgeXErr := ctl.getDeviceApp().DevicesImport(c, file, req.ProductId, req.DriverInstanceId)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}

	httphelper.ResultSuccess(result, c.Writer, lc)
}

// @Tags    设备管理
// @Summary 更新设备
// @Produce json
// @Param   deviceId path   string true "pid"
// @Success 200  {object} httphelper.CommonResponse
// @Router /api/v1/device/:deviceId [put]
func (ctl *controller) DeviceUpdate(c *gin.Context) {
	lc := ctl.lc
	var req dtos.DeviceUpdateRequest
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

// @Tags    设备管理
// @Summary 设备批量绑定驱动
// @Produce json
// @Param   request query    dtos.DevicesBindDriver true "参数"
// @Success 200  {object} httphelper.CommonResponse
// @Router /api/v1/devices/bind-driver [put]
func (ctl *controller) DevicesBindDriver(c *gin.Context) {
	lc := ctl.lc
	var req dtos.DevicesBindDriver
	if err := c.ShouldBind(&req); err != nil {
		httphelper.RenderFail(c, errort.NewCommonErr(errort.DefaultReqParamsError, err), c.Writer, lc)
		return
	}
	err := ctl.getDeviceApp().DevicesBindDriver(c, req)
	if err != nil {
		httphelper.RenderFail(c, err, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(nil, c.Writer, lc)
}

// @Tags    设备管理
// @Summary 设备批量与驱动解绑
// @Produce json
// @Param   request query    dtos.DevicesBindDriver true "参数"
// @Success 200  {object} httphelper.CommonResponse
// @Router /api/v1/unbind-driver [put]
func (ctl *controller) DevicesUnBindDriver(c *gin.Context) {
	lc := ctl.lc
	var req dtos.DevicesUnBindDriver
	if err := c.ShouldBind(&req); err != nil {
		httphelper.RenderFail(c, errort.NewCommonErr(errort.DefaultReqParamsError, err), c.Writer, lc)
		return
	}
	err := ctl.getDeviceApp().DevicesUnBindDriver(c, req)
	if err != nil {
		httphelper.RenderFail(c, err, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(nil, c.Writer, lc)
}

func (ctl *controller) DevicesBindByProductId(c *gin.Context) {
	lc := ctl.lc
	var req dtos.DevicesBindProductId
	if err := c.ShouldBind(&req); err != nil {
		httphelper.RenderFail(c, errort.NewCommonErr(errort.DefaultReqParamsError, err), c.Writer, lc)
		return
	}
	err := ctl.getDeviceApp().DevicesBindProductId(c, req)
	if err != nil {
		httphelper.RenderFail(c, err, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(nil, c.Writer, lc)
}

// @Tags    设备管理
// @Summary 查看设备属性数据
// @Produce json
// @Param   request query   dtos.ThingModelPropertyDataRequest true "参数"
// @Success 200     {array} []dtos.ThingModelDataResponse
// @Router /api/v1/device/:deviceId/thing-model/property [get]
func (ctl *controller) DeviceThingModelPropertyDataSearch(c *gin.Context) {
	lc := ctl.lc
	var req dtos.ThingModelPropertyDataRequest
	urlDecodeParam(&req, c.Request, lc)
	deviceId := c.Param(UrlParamDeviceId)
	req.DeviceId = deviceId
	data, edgeXErr := ctl.getPersistApp().SearchDeviceThingModelPropertyData(req)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(data, c.Writer, lc)
}

func (ctl *controller) DeviceThingModelHistoryPropertyDataSearch(c *gin.Context) {
	lc := ctl.lc
	var req dtos.ThingModelPropertyDataRequest
	urlDecodeParam(&req, c.Request, lc)
	dtos.CorrectionPageParam(&req.BaseSearchConditionQuery)
	deviceId := c.Param(UrlParamDeviceId)
	req.DeviceId = deviceId
	data, total, edgeXErr := ctl.getPersistApp().SearchDeviceThingModelHistoryPropertyData(req)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	pageResult := httphelper.NewPageResult(data, uint32(total), req.Page, req.PageSize)
	httphelper.ResultSuccess(pageResult, c.Writer, lc)
}

// @Tags    设备管理
// @Summary 查看设备事件数据
// @Produce json
// @Param   request query   dtos.ThingModelPropertyDataRequest true "参数"
// @Success 200     {array} []dtos.ThingModelEventDataResponse
// @Router /api/v1/device/:deviceId/thing-model/event [get]
func (ctl *controller) DeviceThingModelEventDataSearch(c *gin.Context) {
	lc := ctl.lc
	var req dtos.ThingModelEventDataRequest
	urlDecodeParam(&req, c.Request, lc)
	deviceId := c.Param(UrlParamDeviceId)
	dtos.CorrectionPageParam(&req.BaseSearchConditionQuery)
	req.DeviceId = deviceId
	data, total, edgeXErr := ctl.getPersistApp().SearchDeviceThingModelEventData(req)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	pageResult := httphelper.NewPageResult(data, uint32(total), req.Page, req.PageSize)
	httphelper.ResultSuccess(pageResult, c.Writer, lc)
}

// @Tags    设备管理
// @Summary 查看设备服务调用数据
// @Produce json
// @Param   request query   dtos.ThingModelServiceDataRequest true "参数"
// @Success 200     {array} []dtos.ThingModelServiceDataResponse
// @Router /api/v1/device/:deviceId/thing-model/service [get]
func (ctl *controller) DeviceThingModelServiceDataSearch(c *gin.Context) {
	lc := ctl.lc
	var req dtos.ThingModelServiceDataRequest
	urlDecodeParam(&req, c.Request, lc)
	deviceId := c.Param(UrlParamDeviceId)
	dtos.CorrectionPageParam(&req.BaseSearchConditionQuery)
	req.DeviceId = deviceId
	data, total, edgeXErr := ctl.getPersistApp().SearchDeviceThingModelServiceData(req)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	pageResult := httphelper.NewPageResult(data, uint32(total), req.Page, req.PageSize)
	httphelper.ResultSuccess(pageResult, c.Writer, lc)
}

func (ctl *controller) DeviceStatusTemplate(c *gin.Context) {
	lc := ctl.lc
	var deviceStatus []constants.DeviceStatus
	deviceStatus = append(append(append(deviceStatus),
		constants.DeviceStatusOnline),
		constants.DeviceStatusOffline)
	httphelper.ResultSuccess(deviceStatus, c.Writer, lc)
}
