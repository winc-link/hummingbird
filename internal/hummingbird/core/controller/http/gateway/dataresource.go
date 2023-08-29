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

package gateway

import (
	"github.com/gin-gonic/gin"
	"github.com/winc-link/hummingbird/internal/dtos"
	"github.com/winc-link/hummingbird/internal/pkg/errort"
	"github.com/winc-link/hummingbird/internal/pkg/httphelper"
)

// @Tags   资源管理
// @Summary 实例类型
// @Produce json
// @Param   request query    dtos.AddDataResourceReq true "参数"
// @Success 200  {object}  httphelper.CommonResponse
// @Router  /api/v1/typeresource [get]
func (ctl *controller) DataResourceType(c *gin.Context) {
	lc := ctl.lc
	types := ctl.getDataResourceApp().DataResourceType(c)
	httphelper.ResultSuccess(types, c.Writer, lc)
}

// @Tags   资源管理
// @Summary 添加资源管理
// @Produce json
// @Param   request query    dtos.AddDataResourceReq true "参数"
// @Success 200  {object}  httphelper.CommonResponse
// @Router  /api/v1/dataresource [post]
func (ctl *controller) DataResourceAdd(c *gin.Context) {
	lc := ctl.lc
	var req dtos.AddDataResourceReq
	if err := c.ShouldBind(&req); err != nil {
		httphelper.RenderFail(c, errort.NewCommonErr(errort.DefaultReqParamsError, err), c.Writer, lc)
		return
	}
	_, edgeXErr := ctl.getDataResourceApp().AddDataResource(c, req)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(nil, c.Writer, lc)
}

// @Tags   资源管理
// @Summary 添加资源管理
// @Produce json
// @Param   request query    dtos.AddDataResourceReq true "参数"
// @Success 200  {object}  httphelper.CommonResponse
// @Router  /api/v1/dataresource/:resourceId [get]
func (ctl *controller) DataResourceById(c *gin.Context) {
	lc := ctl.lc
	id := c.Param(UrlDataResourceId)
	dataSource, edgeXErr := ctl.getDataResourceApp().DataResourceById(c, id)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(dataSource, c.Writer, lc)
}

// @Tags   资源管理
// @Summary 修改资源管理
// @Produce json
// @Param   request query    dtos.AddDataResourceReq true "参数"
// @Success 200  {object}  httphelper.CommonResponse
// @Router  /api/v1/dataresource [put]
func (ctl *controller) UpdateDataResource(c *gin.Context) {
	lc := ctl.lc
	var req dtos.UpdateDataResource
	if err := c.ShouldBind(&req); err != nil {
		httphelper.RenderFail(c, errort.NewCommonErr(errort.DefaultReqParamsError, err), c.Writer, lc)
		return
	}
	edgeXErr := ctl.getDataResourceApp().UpdateDataResource(c, req)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(nil, c.Writer, lc)
}

// @Tags   资源管理
// @Summary 删除资源管理
// @Produce json
// @Success 200  {object}  httphelper.CommonResponse
// @Router  /api/v1/dataresource/:resourceId [delete]
func (ctl *controller) DataResourceDel(c *gin.Context) {
	lc := ctl.lc
	id := c.Param(UrlDataResourceId)
	edgeXErr := ctl.getDataResourceApp().DelDataResourceById(c, id)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(nil, c.Writer, lc)
}

// @Tags   资源管理
// @Summary 资源管理查询
// @Produce json
// @Success 200  {object}  httphelper.CommonResponse
// @Router  /api/v1/dataresource [get]
func (ctl *controller) DataResourceSearch(c *gin.Context) {
	lc := ctl.lc
	var req dtos.DataResourceSearchQueryRequest
	urlDecodeParam(&req, c.Request, lc)
	dtos.CorrectionPageParam(&req.BaseSearchConditionQuery)
	data, total, edgeXErr := ctl.getDataResourceApp().DataResourceSearch(c, req)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	pageResult := httphelper.NewPageResult(data, total, req.Page, req.PageSize)
	httphelper.ResultSuccess(pageResult, c.Writer, lc)
}

func (ctl *controller) DataResourceHealth(c *gin.Context) {
	lc := ctl.lc
	id := c.Param(UrlDataResourceId)
	edgeXErr := ctl.getDataResourceApp().DataResourceHealth(c, id)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(nil, c.Writer, lc)
}
