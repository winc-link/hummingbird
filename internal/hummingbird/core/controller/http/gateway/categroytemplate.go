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
	"github.com/winc-link/hummingbird/internal/pkg/httphelper"
)

// @Tags 标准产品品类
// @Summary 标准产品品类列表
// @Produce json
// @Param request query dtos.CategoryTemplateRequest true "参数"
// @Success 200  {array} dtos.CategoryTemplateResponse
// @Router  /api/v1/category-template [get]
//@Security ApiKeyAuth
func (ctl *controller) CategoryTemplateSearch(c *gin.Context) {
	lc := ctl.lc
	var req dtos.CategoryTemplateRequest
	urlDecodeParam(&req, c.Request, lc)
	dtos.CorrectionPageParam(&req.BaseSearchConditionQuery)
	data, total, edgeXErr := ctl.getCategoryTemplateApp().CategoryTemplateSearch(c, req)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	pageResult := httphelper.NewPageResult(data, total, req.Page, req.PageSize)
	httphelper.ResultSuccess(pageResult, c.Writer, lc)
}

// @Tags 标准产品品类
// @Summary 同步标准产品品类
// @Produce json
// @Param request query dtos.CategoryTemplateRequest true "参数"
// @Router  /api/v1/category-template/sync [post]
//@Security ApiKeyAuth
func (ctl *controller) CategoryTemplateSync(c *gin.Context) {
	lc := ctl.lc
	var req dtos.CategoryTemplateSyncRequest
	urlDecodeParam(&req, c.Request, lc)
	_, edgeXErr := ctl.getCategoryTemplateApp().Sync(c, "Ireland")
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(nil, c.Writer, lc)
}
