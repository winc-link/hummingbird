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

// @Tags 产品管理
// @Summary 查询产品列表
// @Produce json
// @Param   request query   dtos.ProductSearchQueryRequest true "参数"
// @Success 200     {array} dtos.ProductSearchQueryResponse
// @Router  /api/v1/products [get]
// @Security ApiKeyAuth
func (ctl *controller) ProductsSearch(c *gin.Context) {
	lc := ctl.lc
	var req dtos.ProductSearchQueryRequest
	urlDecodeParam(&req, c.Request, lc)
	dtos.CorrectionPageParam(&req.BaseSearchConditionQuery)

	data, total, edgeXErr := ctl.getProductApp().ProductsSearch(c, req)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	pageResult := httphelper.NewPageResult(data, total, req.Page, req.PageSize)
	httphelper.ResultSuccess(pageResult, c.Writer, lc)
}

// @Tags 产品管理
// @Summary 查询产品详情
// @Produce json
// @Param   productId path     string true "pid"
// @Success 200       {object} dtos.ProductSearchByIdResponse
// @Router  /api/v1/product/:productId [get]
// @Security ApiKeyAuth
func (ctl *controller) ProductById(c *gin.Context) {
	lc := ctl.lc
	id := c.Param(UrlParamProductId)
	data, edgeXErr := ctl.getProductApp().ProductById(c, id)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(data, c.Writer, lc)
}

// @Tags 产品管理
// @Summary 删除产品
// @Produce json
// @Param   productId path     string true "pid"
// @Success 200  {object} httphelper.CommonResponse
// @Router  /api/v1/product/:productId [delete]
// @Security ApiKeyAuth
func (ctl *controller) ProductDelete(c *gin.Context) {
	lc := ctl.lc
	id := c.Param(UrlParamProductId)
	edgeXErr := ctl.getProductApp().ProductDelete(c, id)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(nil, c.Writer, lc)
}

func (ctl *controller) ProductRelease(c *gin.Context) {
	lc := ctl.lc
	productId := c.Param(UrlParamProductId)
	edgeXErr := ctl.getProductApp().ProductRelease(c, productId)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(nil, c.Writer, lc)
}

func (ctl *controller) ProductUnRelease(c *gin.Context) {
	lc := ctl.lc
	productId := c.Param(UrlParamProductId)
	edgeXErr := ctl.getProductApp().ProductUnRelease(c, productId)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(nil, c.Writer, lc)
}

// @Tags 产品管理
// @Summary 添加产品
// @Produce json
// @Param request body dtos.ProductAddRequest true "参数"
// @Success 200  {object} httphelper.CommonResponse
// @Router  /api/v1/product [post]
// @Security ApiKeyAuth
func (ctl *controller) ProductAdd(c *gin.Context) {
	lc := ctl.lc
	var req dtos.ProductAddRequest
	if err := c.ShouldBind(&req); err != nil {
		httphelper.RenderFail(c, errort.NewCommonErr(errort.DefaultReqParamsError, err), c.Writer, lc)
		return
	}
	_, edgeXErr := ctl.getProductApp().AddProduct(c, req)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(nil, c.Writer, lc)
}

// @Tags 产品管理
// @Summary 云平台列表
// @Produce json
// @Success 200  {object} httphelper.CommonResponse
// @Router  /api/v1/iot-platform [get]
// @Security ApiKeyAuth
func (ctl *controller) IotPlatform(c *gin.Context) {
	lc := ctl.lc
	var iotPlatform []constants.IotPlatform
	iotPlatform = append(iotPlatform, constants.IotPlatform_LocalIot)
	httphelper.ResultSuccess(iotPlatform, c.Writer, lc)

}
