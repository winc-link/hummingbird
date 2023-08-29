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

// @Tags 运维管理
// @Summary 获取系统性能
// @Produce json
// @Param request query dtos.SystemMetricsQuery true "参数"
// @Success 200 {object} dtos.SystemMetricsResponse
// @Router /api/v1/metrics/system [get]
//func (ctl *controller) GetSystemMetricsHandle(c *gin.Context) {
//	ctl.ProxyAgentServer(c)
//}

func (c *controller) SystemMetricsHandler(ctx *gin.Context) {
	var query = dtos.SystemMetricsQuery{}
	if err := ctx.BindQuery(&query); err != nil {
		httphelper.RenderFail(ctx, errort.NewCommonErr(errort.DefaultReqParamsError, err), ctx.Writer, c.lc)
		return
	}
	metrics, err := c.getSystemMonitorApp().GetSystemMetrics(ctx, query)
	if err != nil {
		httphelper.RenderFail(ctx, err, ctx.Writer, c.lc)
		return
	}

	httphelper.ResultSuccess(metrics, ctx.Writer, c.lc)
}

// @Tags 运维管理
// @Summary 操作服务重启
// @Produce json
// @Param  request  body dtos.Operation true "操作"
// @Success 200 {object} httphelper.CommonResponse
// @Router /api/v1/operation [post]
//func (ctl *controller) OperationServiceHandle(c *gin.Context) {
//	ctl.ProxyAgentServer(c)
//}
