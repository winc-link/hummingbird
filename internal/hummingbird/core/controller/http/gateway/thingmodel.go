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

// @Tags 物模型
// @Summary 查询系统物模型
// @Produce json
// @Param request query dtos.SystemThingModelSearchReq true "参数"
// @Success 200  {object} httphelper.CommonResponse
// @Router  /api/v1/thingmodel/system [get]
// @Security ApiKeyAuth
func (ctl *controller) SystemThingModelSearch(c *gin.Context) {
	lc := ctl.lc
	var req dtos.SystemThingModelSearchReq
	urlDecodeParam(&req, c.Request, lc)
	data, edgeXErr := ctl.getThingModelApp().SystemThingModelSearch(c, req)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(data, c.Writer, lc)
}

// @Tags 物模型
// @Summary 产品添加物模型
// @Produce json
// @Param request body dtos.ThingModelAddOrUpdateReq true "参数"
// @Success 200  {object} httphelper.CommonResponse
// @Router  /api/v1/thingmodel [post]
// @Security ApiKeyAuth
func (ctl *controller) ThingModelAdd(c *gin.Context) {
	lc := ctl.lc
	var req dtos.ThingModelAddOrUpdateReq
	if err := c.ShouldBind(&req); err != nil {
		httphelper.RenderFail(c, errort.NewCommonErr(errort.DefaultReqParamsError, err), c.Writer, lc)
		return
	}
	_, edgeXErr := ctl.getThingModelApp().AddThingModel(c, req)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(nil, c.Writer, lc)
}

// @Tags 物模型
// @Summary 修改产品物模型
// @Produce json
// @Param request body dtos.ThingModelAddOrUpdateReq true "参数"
// @Success 200  {object} httphelper.CommonResponse
// @Router  /api/v1/thingmodel [put]
// @Security ApiKeyAuth
func (ctl *controller) ThingModelUpdate(c *gin.Context) {
	lc := ctl.lc
	var req dtos.ThingModelAddOrUpdateReq
	if err := c.ShouldBind(&req); err != nil {
		httphelper.RenderFail(c, errort.NewCommonErr(errort.DefaultReqParamsError, err), c.Writer, lc)
		return
	}
	edgeXErr := ctl.getThingModelApp().UpdateThingModel(c, req)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(nil, c.Writer, lc)
}

// @Tags 物模型
// @Summary 产品删除物模型
// @Produce json
// @Param request body dtos.ThingModelDeleteReq true "参数"
// @Success 200  {object} httphelper.CommonResponse
// @Router  /api/v1/thingmodel [delete]
// @Security ApiKeyAuth
func (ctl *controller) ThingModelDelete(c *gin.Context) {
	lc := ctl.lc
	var req dtos.ThingModelDeleteReq
	if err := c.ShouldBind(&req); err != nil {
		httphelper.RenderFail(c, errort.NewCommonErr(errort.DefaultReqParamsError, err), c.Writer, lc)
		return
	}
	edgeXErr := ctl.getThingModelApp().ThingModelDelete(c, req.ThingModelId, req.ThingModelType)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(nil, c.Writer, lc)
}

// @Tags 物模型
// @Summary 物模型单位
// @Produce json
// @Param request query dtos.UnitRequest true "参数"
// @Success 200  {array} dtos.UnitResponse
// @Router  /api/v1/thingmodel/unit [get]
// @Security ApiKeyAuth
func (ctl *controller) ThingModelUnit(c *gin.Context) {
	lc := ctl.lc
	var req dtos.UnitRequest
	urlDecodeParam(&req, c.Request, lc)
	dtos.CorrectionPageParam(&req.BaseSearchConditionQuery)
	data, total, edgeXErr := ctl.getUnitModelApp().UnitTemplateSearch(c, req)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	pageResult := httphelper.NewPageResult(data, total, req.Page, req.PageSize)
	httphelper.ResultSuccess(pageResult, c.Writer, lc)
}

// @Tags 物模型
// @Summary 物模型单位同步
// @Produce json
// @Param request query dtos.UnitTemplateSyncRequest true "参数"
// @Router  /api/v1/thingmodel/unit-sync [post]
// @Security ApiKeyAuth
func (ctl *controller) ThingModelUnitSync(c *gin.Context) {
	lc := ctl.lc
	var req dtos.UnitTemplateSyncRequest
	urlDecodeParam(&req, c.Request, lc)
	total, edgeXErr := ctl.getUnitModelApp().Sync(c, "Ireland")
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(total, c.Writer, lc)
}

func (ctl *controller) ThingModelDocsSync(c *gin.Context) {
	lc := ctl.lc
	total, edgeXErr := ctl.getDocsApp().SyncDocs(c, "Ireland")
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(total, c.Writer, lc)
}

func (ctl *controller) ThingModelQuickNavigationSync(c *gin.Context) {
	lc := ctl.lc
	total, edgeXErr := ctl.getQuickNavigationApp().SyncQuickNavigation(c, "Ireland")
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(total, c.Writer, lc)
}

//func (ctl *controller) MsgGather(c *gin.Context) {
//	lc := ctl.lc
//	edgeXErr := ctl.getDeviceApp().DevicesReportMsgGather(context.Background())
//	if edgeXErr != nil {
//		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
//		return
//	}
//	httphelper.ResultSuccess(nil, c.Writer, lc)
//}
