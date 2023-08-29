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

// @Tags   规则引擎
// @Summary 添加规则引擎
// @Produce json
// @Param   request query    dtos.RuleEngineRequest true "参数"
// @Success 200  {object}  httphelper.CommonResponse
// @Router  /api/v1/rule-engine [post]
func (ctl *controller) RuleEngineAdd(c *gin.Context) {
	lc := ctl.lc
	var req dtos.RuleEngineRequest
	if err := c.ShouldBind(&req); err != nil {
		httphelper.RenderFail(c, errort.NewCommonErr(errort.DefaultReqParamsError, err), c.Writer, lc)
		return
	}
	_, edgeXErr := ctl.getRuleEngineApp().AddRuleEngine(c, req)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(nil, c.Writer, lc)
}

// @Tags   规则引擎
// @Summary 编辑规则引擎
// @Produce json
// @Param   request query    dtos.RuleEngineUpdateRequest true "参数"
// @Success 200  {object}  httphelper.CommonResponse
// @Router  /api/v1/rule-engine [put]
func (ctl *controller) RuleEngineUpdate(c *gin.Context) {
	lc := ctl.lc
	var req dtos.RuleEngineUpdateRequest
	if err := c.ShouldBind(&req); err != nil {
		httphelper.RenderFail(c, errort.NewCommonErr(errort.DefaultReqParamsError, err), c.Writer, lc)
		return
	}
	edgeXErr := ctl.getRuleEngineApp().UpdateRuleEngine(c, req)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(nil, c.Writer, lc)
}

// @Tags    规则引擎
// @Summary 规则引擎详情
// @Produce json
// @Param   ruleEngineId path   string true "ruleEngineId"
// @Success 200  {object} httphelper.CommonResponse
// @Router /api/v1/rule-engine/:ruleEngineId [get]
func (ctl *controller) RuleEngineById(c *gin.Context) {
	lc := ctl.lc
	id := c.Param(RuleEngineId)
	data, edgeXErr := ctl.getRuleEngineApp().RuleEngineById(c, id)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(data, c.Writer, lc)
}

// @Tags    规则引擎
// @Summary 规则引擎列表
// @Produce json
// @Param   request query   dtos.RuleEngineSearchQueryRequest true "参数"
// @Success 200     {array} []dtos.RuleEngineSearchQueryResponse
// @Router  /api/v1/rule-engine [get]
func (ctl *controller) RuleEngineSearch(c *gin.Context) {
	lc := ctl.lc
	var req dtos.RuleEngineSearchQueryRequest
	urlDecodeParam(&req, c.Request, lc)
	dtos.CorrectionPageParam(&req.BaseSearchConditionQuery)
	data, total, edgeXErr := ctl.getRuleEngineApp().RuleEngineSearch(c, req)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	pageResult := httphelper.NewPageResult(data, total, req.Page, req.PageSize)
	httphelper.ResultSuccess(pageResult, c.Writer, lc)
}

// @Tags    规则引擎
// @Summary 规则引擎启动
// @Produce json
// @Param   ruleEngineId path   string true "ruleEngineId"
// @Success 200  {object} httphelper.CommonResponse
// @Router /api/v1/rule-engine/:ruleEngineId/start [post]
func (ctl *controller) RuleEngineStart(c *gin.Context) {
	lc := ctl.lc
	id := c.Param(RuleEngineId)
	edgeXErr := ctl.getRuleEngineApp().RuleEngineStart(c, id)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(nil, c.Writer, lc)
}

// @Tags    规则引擎
// @Summary 规则引擎停止
// @Produce json
// @Param   ruleEngineId path   string true "ruleEngineId"
// @Success 200  {object} httphelper.CommonResponse
// @Router /api/v1/rule-engine/:ruleEngineId/stop [post]
func (ctl *controller) RuleEngineStop(c *gin.Context) {
	lc := ctl.lc
	id := c.Param(RuleEngineId)
	edgeXErr := ctl.getRuleEngineApp().RuleEngineStop(c, id)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(nil, c.Writer, lc)
}

// @Tags    规则引擎
// @Summary 规则引擎删除
// @Produce json
// @Param   ruleEngineId path   string true "ruleEngineId"
// @Success 200  {object} httphelper.CommonResponse
// @Router  /api/v1/rule-engine/:ruleEngineId/delete [delete]
func (ctl *controller) RuleEngineDelete(c *gin.Context) {
	lc := ctl.lc
	id := c.Param(RuleEngineId)
	edgeXErr := ctl.getRuleEngineApp().RuleEngineDelete(c, id)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(nil, c.Writer, lc)
}

// @Tags    规则引擎
// @Summary 规则引擎状态
// @Produce json
// @Param   ruleEngineId path   string true "ruleEngineId"
// @Success 200  {object} httphelper.CommonResponse
// @Router  /api/v1/rule-engine/:ruleEngineId/status [get]
func (ctl *controller) RuleEngineStatus(c *gin.Context) {
	lc := ctl.lc
	id := c.Param(RuleEngineId)
	r, edgeXErr := ctl.getRuleEngineApp().RuleEngineStatus(c, id)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(r, c.Writer, lc)
}
