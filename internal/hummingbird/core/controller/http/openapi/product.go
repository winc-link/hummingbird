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

func (ctl *controller) OpenApiCreateProduct(c *gin.Context) {
	lc := ctl.lc
	var req dtos.OpenApiAddProductRequest
	if err := c.ShouldBind(&req); err != nil {
		httphelper.RenderFail(c, errort.NewCommonErr(errort.DefaultReqParamsError, err), c.Writer, lc)
		return
	}
	_, edgeXErr := ctl.getProductApp().OpenApiAddProduct(c, req)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(nil, c.Writer, lc)
}

func (ctl *controller) OpenApiProductById(c *gin.Context) {
	lc := ctl.lc
	id := c.Param(UrlParamProductId)
	data, edgeXErr := ctl.getProductApp().OpenApiProductById(c, id)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(data, c.Writer, lc)
}

func (ctl *controller) OpenApiProductReleaseById(c *gin.Context) {
	lc := ctl.lc
	id := c.Param(UrlParamProductId)
	edgeXErr := ctl.getProductApp().ProductRelease(c, id)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(nil, c.Writer, lc)
}

func (ctl *controller) OpenApiProductUnReleaseById(c *gin.Context) {
	lc := ctl.lc
	id := c.Param(UrlParamProductId)
	edgeXErr := ctl.getProductApp().ProductUnRelease(c, id)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(nil, c.Writer, lc)
}

func (ctl *controller) OpenApiProductSearch(c *gin.Context) {
	lc := ctl.lc
	var req dtos.ProductSearchQueryRequest
	urlDecodeParam(&req, c.Request, lc)
	dtos.CorrectionPageParam(&req.BaseSearchConditionQuery)

	data, total, edgeXErr := ctl.getProductApp().OpenApiProductSearch(c, req)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	pageResult := httphelper.NewPageResult(data, total, req.Page, req.PageSize)
	httphelper.ResultSuccess(pageResult, c.Writer, lc)
}

func (ctl *controller) OpenApiDeleteProduct(c *gin.Context) {
	lc := ctl.lc
	id := c.Param(UrlParamProductId)
	edgeXErr := ctl.getProductApp().ProductDelete(c, id)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(nil, c.Writer, lc)
}

func (ctl *controller) OpenApiUpdateProduct(c *gin.Context) {
	lc := ctl.lc
	id := c.Param(UrlParamProductId)
	var req dtos.OpenApiUpdateProductRequest
	if err := c.ShouldBind(&req); err != nil {
		httphelper.RenderFail(c, errort.NewCommonErr(errort.DefaultReqParamsError, err), c.Writer, lc)
		return
	}
	req.Id = id
	edgeXErr := ctl.getProductApp().OpenApiUpdateProduct(c, req)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(nil, c.Writer, lc)

}
