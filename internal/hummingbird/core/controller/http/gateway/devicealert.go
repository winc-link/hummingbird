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
	"time"
)

// @Tags   告警中心
// @Summary 添加告警规则
// @Produce json
// @Param   request query    dtos.RuleAddRequest true "参数"
// @Success 200  {object}  httphelper.CommonResponse
// @Router  /api/v1/alert-rule [post]
func (ctl *controller) AlertRuleAdd(c *gin.Context) {
	lc := ctl.lc
	var req dtos.RuleAddRequest
	if err := c.ShouldBind(&req); err != nil {
		httphelper.RenderFail(c, errort.NewCommonErr(errort.DefaultReqParamsError, err), c.Writer, lc)
		return
	}
	_, edgeXErr := ctl.getAlertRuleApp().AddAlertRule(c, req)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(nil, c.Writer, lc)
}

// @Tags   告警中心
// @Summary 编辑告警规则
// @Produce json
// @Param   request query    dtos.RuleUpdateRequest true "参数"
// @Success 200  {object}  httphelper.CommonResponse
// @Router  /api/v1/alert-rule/:ruleId [put]
func (ctl *controller) AlertRuleUpdate(c *gin.Context) {
	lc := ctl.lc
	var req dtos.RuleUpdateRequest
	if err := c.ShouldBind(&req); err != nil {
		httphelper.RenderFail(c, errort.NewCommonErr(errort.DefaultReqParamsError, err), c.Writer, lc)
		return
	}
	edgeXErr := ctl.getAlertRuleApp().UpdateAlertRule(c, req)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(nil, c.Writer, lc)
}

func (ctl *controller) AlertRuleUpdateField(c *gin.Context) {
	lc := ctl.lc
	var req dtos.RuleFieldUpdate
	if err := c.ShouldBind(&req); err != nil {
		httphelper.RenderFail(c, errort.NewCommonErr(errort.DefaultReqParamsError, err), c.Writer, lc)
		return
	}
	edgeXErr := ctl.getAlertRuleApp().UpdateAlertField(c, req)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(nil, c.Writer, lc)
}

// @Tags    告警中心
// @Summary 告警规则详情
// @Produce json
// @Param   ruleId path   string true "ruleId"
// @Success 200  {object} httphelper.CommonResponse
// @Router /api/v1/alert-rule/:ruleId [get]
func (ctl *controller) AlertRuleById(c *gin.Context) {
	lc := ctl.lc
	id := c.Param(UrlParamRuleId)
	data, edgeXErr := ctl.getAlertRuleApp().AlertRuleById(c, id)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(data, c.Writer, lc)
}

// @Tags    告警中心
// @Summary 告警规则列表
// @Produce json
// @Param   request query   dtos.AlertRuleSearchQueryRequest true "参数"
// @Success 200     {array} []dtos.AlertRuleSearchQueryResponse
// @Router  /api/v1/alert-rule [get]
func (ctl *controller) AlertRuleSearch(c *gin.Context) {
	lc := ctl.lc
	var req dtos.AlertRuleSearchQueryRequest
	urlDecodeParam(&req, c.Request, lc)
	dtos.CorrectionPageParam(&req.BaseSearchConditionQuery)
	data, total, edgeXErr := ctl.getAlertRuleApp().AlertRulesSearch(c, req)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	pageResult := httphelper.NewPageResult(data, total, req.Page, req.PageSize)
	httphelper.ResultSuccess(pageResult, c.Writer, lc)
}

// @Tags    告警中心
// @Summary 告警规则启动
// @Produce json
// @Param   ruleId path   string true "ruleId"
// @Success 200  {object} httphelper.CommonResponse
// @Router /api/v1/alert-rule/:ruleId/start [post]
func (ctl *controller) AlertRuleStart(c *gin.Context) {
	lc := ctl.lc
	id := c.Param(UrlParamRuleId)
	edgeXErr := ctl.getAlertRuleApp().AlertRulesStart(c, id)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(nil, c.Writer, lc)
}

// @Tags    告警中心
// @Summary 告警规则停止
// @Produce json
// @Param   ruleId path   string true "ruleId"
// @Success 200  {object} httphelper.CommonResponse
// @Router /api/v1/alert-rule/:ruleId/stop [post]
func (ctl *controller) AlertRuleStop(c *gin.Context) {
	lc := ctl.lc
	id := c.Param(UrlParamRuleId)
	edgeXErr := ctl.getAlertRuleApp().AlertRulesStop(c, id)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(nil, c.Writer, lc)
}

// @Tags    告警中心
// @Summary 告警规则重启
// @Produce json
// @Param   ruleId path   string true "ruleId"
// @Success 200  {object} httphelper.CommonResponse
// @Router /api/v1/alert-rule/:ruleId/restart [post]
func (ctl *controller) AlertRuleRestart(c *gin.Context) {
	lc := ctl.lc
	id := c.Param(UrlParamRuleId)
	edgeXErr := ctl.getAlertRuleApp().AlertRulesRestart(c, id)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(nil, c.Writer, lc)
}

// @Tags    告警中心
// @Summary 告警规则删除
// @Produce json
// @Param   ruleId path   string true "ruleId"
// @Success 200  {object} httphelper.CommonResponse
// @Router /api/v1/alert-rule/:ruleId [delete]
func (ctl *controller) AlertRuleDelete(c *gin.Context) {
	lc := ctl.lc
	id := c.Param(UrlParamRuleId)
	edgeXErr := ctl.getAlertRuleApp().AlertRulesDelete(c, id)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(nil, c.Writer, lc)
}

// @Tags    告警中心
// @Summary 告警列表
// @Produce json
// @Param   request query   dtos.AlertSearchQueryRequest true "参数"
// @Success 200     {array} []dtos.AlertSearchQueryResponse
// @Router  /api/v1/alert-list [get]
func (ctl *controller) AlertSearch(c *gin.Context) {
	lc := ctl.lc
	var req dtos.AlertSearchQueryRequest
	urlDecodeParam(&req, c.Request, lc)
	dtos.CorrectionPageParam(&req.BaseSearchConditionQuery)
	data, total, edgeXErr := ctl.getAlertRuleApp().AlertSearch(c, req)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	pageResult := httphelper.NewPageResult(data, total, req.Page, req.PageSize)
	httphelper.ResultSuccess(pageResult, c.Writer, lc)
}

// @Tags    告警中心
// @Summary 告警列表
// @Produce json
// @Param   request query   dtos.AlertSearchQueryRequest true "参数"
// @Success 200     {array} []dtos.AlertSearchQueryResponse
// @Router  /api/v1/alert-plate [get]
func (ctl *controller) AlertPlate(c *gin.Context) {
	lc := ctl.lc
	currentTime := time.Now()
	beforeTime := currentTime.AddDate(0, 0, -7).UnixMilli()
	data, edgeXErr := ctl.getAlertRuleApp().AlertPlate(c, beforeTime)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(data, c.Writer, lc)
}

// @Tags    告警中心
// @Summary 忽略告警
// @Produce json
// @Router  /api/v1/alert-ignore/:ruleId [put]
func (ctl *controller) AlertIgnore(c *gin.Context) {
	lc := ctl.lc
	id := c.Param(UrlParamRuleId)
	edgeXErr := ctl.getAlertRuleApp().AlertIgnore(c, id)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(nil, c.Writer, lc)
}

// @Tags    告警中心
// @Summary 处理告警
// @Produce json
// @Router  /api/v1/alert-treated [post]
func (ctl *controller) AlertTreated(c *gin.Context) {
	lc := ctl.lc
	var req dtos.AlertTreatedRequest
	if err := c.ShouldBind(&req); err != nil {
		httphelper.RenderFail(c, errort.NewCommonErr(errort.DefaultReqParamsError, err), c.Writer, lc)
		return
	}
	edgeXErr := ctl.getAlertRuleApp().TreatedIgnore(c, req.Id, req.Message)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(nil, c.Writer, lc)
}

// @Tags    告警中心
// @Summary 告警列表
// @Produce json
// @Param   request query   dtos.AlertAddRequest true "参数"
// @Success 200  {object}  httphelper.CommonResponse
// @Router  /api/v1/alert [post]
func (ctl *controller) EkuiperAlert(c *gin.Context) {
	lc := ctl.lc

	req := make(map[string]interface{})
	if err := c.ShouldBind(&req); err != nil {
		httphelper.RenderFail(c, errort.NewCommonErr(errort.DefaultReqParamsError, err), c.Writer, lc)
		return
	}
	lc.Info("req....", req)
	edgeXErr := ctl.getAlertRuleApp().AddAlert(c, req)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(nil, c.Writer, lc)
}
