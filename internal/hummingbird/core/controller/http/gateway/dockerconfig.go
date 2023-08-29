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

// @Tags    镜像仓库管理
// @Summary 新增镜像
// @Produce json
// @Param   request body     dtos.DockerConfigAddRequest true "参数"
// @Success 200     {object} httphelper.CommonResponse
// @Router  /api/v1/docker-configs [post]
func (ctl *controller) DockerConfigAdd(c *gin.Context) {
	lc := ctl.lc
	var req dtos.DockerConfigAddRequest
	if err := c.ShouldBind(&req); err != nil {
		httphelper.RenderFail(c, errort.NewCommonErr(errort.DefaultReqParamsError, err), c.Writer, lc)
		return
	}
	edgeXErr := ctl.getDriverLibApp().DownConfigAdd(c, req)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(nil, c.Writer, lc)
}

// @Tags    镜像仓库管理
// @Summary 获取镜像列表
// @Produce json
// @Param   request query    dtos.DockerConfigSearchQueryRequest true "参数"
// @Success 200     {object} httphelper.ResPageResult
// @Router  /api/v1/docker-configs [get]
func (ctl *controller) DockerConfigsSearch(c *gin.Context) {
	lc := ctl.lc
	var req dtos.DockerConfigSearchQueryRequest
	urlDecodeParam(&req, c.Request, lc)
	dtos.CorrectionPageParam(&req.BaseSearchConditionQuery)

	list, total, edgeXErr := ctl.getDriverLibApp().DownConfigSearch(c, req)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	dcs := make([]dtos.DockerConfigResponse, len(list))
	for i, v := range list {
		dcs[i] = dtos.DockerConfigResponseFromModel(v)
	}
	pageResult := httphelper.NewPageResult(dcs, total, req.Page, req.PageSize)

	httphelper.ResultSuccess(pageResult, c.Writer, lc)
}

// @Tags    镜像仓库管理
// @Summary 修改仓库信息
// @Produce json
// @Param   request body     dtos.DockerConfigUpdateRequest true "参数"
// @Success 200     {object} httphelper.CommonResponse
// @Router  /api/v1/docker-configs/:dockerConfigId [put]
func (ctl *controller) DockerConfigUpdate(c *gin.Context) {
	lc := ctl.lc
	var req dtos.DockerConfigUpdateRequest
	req.Id = c.Param(UrlParamDockerConfigId)
	if err := c.ShouldBind(&req); err != nil {
		httphelper.RenderFail(c, errort.NewCommonErr(errort.DefaultReqParamsError, err), c.Writer, lc)
		return
	}

	edgeXErr := ctl.getDriverLibApp().DownConfigUpdate(c, req)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(nil, c.Writer, lc)
}

// @Tags    镜像仓库管理
// @Summary 删除仓库信息
// @Produce json
// @Param   dockerConfigId path     string true "镜像ID"
// @Success 200            {object} httphelper.CommonResponse
// @Router  /api/v1/docker-configs/:dockerConfigId [delete]
func (ctl *controller) DockerConfigDelete(c *gin.Context) {
	lc := ctl.lc
	id := c.Param(UrlParamDockerConfigId)
	edgeXErr := ctl.getDriverLibApp().DownConfigDel(c, id)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(nil, c.Writer, lc)
}
