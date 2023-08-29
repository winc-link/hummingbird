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

// @Tags 配网助手
// @Summary 获取网卡列表
// @Produce json
// @Success 200 {object} dtos.ConfigNetWorkResponse
// @Router /api/v1/local/config/network [get]
func (ctl *controller) ConfigNetWorkGet(c *gin.Context) {
	lc := ctl.lc
	res, edgeXErr := ctl.getSystemApp().ConfigNetWork(c, false)
	if edgeXErr != nil {
		httphelper.RenderFail(c, errort.NewCommonErr(errort.DefaultReqParamsError, edgeXErr), c.Writer, lc)
		return
	}

	httphelper.ResultSuccess(res, c.Writer, lc)
}

// @Tags 配网助手
// @Summary 修改网卡
// @Produce json
// @Param req body dtos.ConfigNetworkUpdateRequest true "参数"
// @Success 200 {object} dtos.ConfigNetWorkResponse
// @Router /api/v1/local/config/network [put]
func (ctl *controller) ConfigNetWorkUpdate(c *gin.Context) {
	lc := ctl.lc
	var req dtos.ConfigNetworkUpdateRequest
	if err := c.ShouldBind(&req); err != nil {
		httphelper.RenderFail(c, errort.NewCommonErr(errort.DefaultReqParamsError, err), c.Writer, lc)
		return
	}
	edgeXErr := ctl.getSystemApp().ConfigNetWorkUpdate(c, req)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}

	httphelper.ResultSuccess(nil, c.Writer, lc)
}

// @Tags 配网助手
// @Summary 获取dns
// @Produce json
// @Success 200 {object} dtos.ConfigDnsResponse
// @Router /api/v1/local/config/dns [get]
func (ctl *controller) ConfigDnsGet(c *gin.Context) {
	lc := ctl.lc
	resp, edgeXErr := ctl.getSystemApp().ConfigDns(c)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}

	httphelper.ResultSuccess(resp, c.Writer, lc)
}

// @Tags 配网助手
// @Summary 修改dns
// @Produce json
// @Param req body dtos.ConfigDnsUpdateRequest true "参数"
// @Success 200 {object} dtos.ConfigDnsResponse
// @Router /api/v1/local/config/dns [put]
func (ctl *controller) ConfigDnsUpdate(c *gin.Context) {
	lc := ctl.lc
	var req dtos.ConfigDnsUpdateRequest
	if err := c.ShouldBind(&req); err != nil {
		httphelper.RenderFail(c, errort.NewCommonErr(errort.DefaultReqParamsError, err), c.Writer, lc)
		return
	}
	edgeXErr := ctl.getSystemApp().ConfigDnsUpdate(c, req)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}

	httphelper.ResultSuccess(nil, c.Writer, lc)
}
