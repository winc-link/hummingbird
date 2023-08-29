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

func (ctl *controller) SceneAdd(c *gin.Context) {
	lc := ctl.lc
	var req dtos.SceneAddRequest
	if err := c.ShouldBind(&req); err != nil {
		httphelper.RenderFail(c, errort.NewCommonErr(errort.DefaultReqParamsError, err), c.Writer, lc)
		return
	}
	_, edgeXErr := ctl.getSceneApp().AddScene(c, req)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(nil, c.Writer, lc)
}

func (ctl *controller) SceneUpdate(c *gin.Context) {
	lc := ctl.lc
	var req dtos.SceneUpdateRequest
	if err := c.ShouldBind(&req); err != nil {
		httphelper.RenderFail(c, errort.NewCommonErr(errort.DefaultReqParamsError, err), c.Writer, lc)
		return
	}
	edgeXErr := ctl.getSceneApp().UpdateScene(c, req)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(nil, c.Writer, lc)
}

func (ctl *controller) SceneById(c *gin.Context) {
	lc := ctl.lc
	id := c.Param(UrlParamSceneId)
	scene, edgeXErr := ctl.getSceneApp().SceneById(c, id)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(scene, c.Writer, lc)
}

func (ctl *controller) SearchScene(c *gin.Context) {
	lc := ctl.lc
	var req dtos.SceneSearchQueryRequest
	urlDecodeParam(&req, c.Request, lc)
	dtos.CorrectionPageParam(&req.BaseSearchConditionQuery)

	list, total, edgeXErr := ctl.getSceneApp().SceneSearch(c, req)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	pageResult := httphelper.NewPageResult(list, total, req.Page, req.PageSize)
	httphelper.ResultSuccess(pageResult, c.Writer, lc)
}

func (ctl *controller) SceneLogSearch(c *gin.Context) {
	lc := ctl.lc
	sceneId := c.Param(UrlParamSceneId)
	var req dtos.SceneLogSearchQueryRequest
	urlDecodeParam(&req, c.Request, lc)
	dtos.CorrectionPageParam(&req.BaseSearchConditionQuery)
	req.SceneId = sceneId

	ctl.lc.Info("sceneLogSearch log:", req)
	list, total, edgeXErr := ctl.getSceneApp().SceneLogSearch(c, req)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	pageResult := httphelper.NewPageResult(list, total, req.Page, req.PageSize)
	httphelper.ResultSuccess(pageResult, c.Writer, lc)
}

func (ctl *controller) SceneStart(c *gin.Context) {
	lc := ctl.lc
	id := c.Param(UrlParamSceneId)
	edgeXErr := ctl.getSceneApp().SceneStartById(c, id)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(nil, c.Writer, lc)
}

func (ctl *controller) SceneStop(c *gin.Context) {
	lc := ctl.lc
	id := c.Param(UrlParamSceneId)
	edgeXErr := ctl.getSceneApp().SceneStopById(c, id)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(nil, c.Writer, lc)
}

func (ctl *controller) DeleteScene(c *gin.Context) {
	lc := ctl.lc
	id := c.Param(UrlParamSceneId)
	edgeXErr := ctl.getSceneApp().DelSceneById(c, id)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(nil, c.Writer, lc)
}

func (ctl *controller) SceneLog(c *gin.Context) {
	lc := ctl.lc
	id := c.Param(UrlParamSceneId)
	edgeXErr := ctl.getSceneApp().DelSceneById(c, id)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(nil, c.Writer, lc)
}
