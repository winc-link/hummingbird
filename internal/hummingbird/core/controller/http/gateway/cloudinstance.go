/*******************************************************************************
 * Copyright 2017 Dell Inc.
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
	"github.com/winc-link/hummingbird/internal/pkg/httphelper"
)

// @Tags 我的云服务实例
// @Summary 云服务实例列表
// @Produce json
// @Param request query dtos.CloudInstanceSearchQueryRequest true "参数"
// @Success 200  {object} dtos.CloudInstanceSearchQueryRequest
// @Router  /api/v1/cloud-instance [get]
//@Security ApiKeyAuth
func (ctl *controller) CloudInstanceSearch(c *gin.Context) {
	lc := ctl.lc
	//var req dtos.CloudInstanceSearchQueryRequest
	//urlDecodeParam(&req, c.Request, lc)
	//dtos.CorrectionPageParam(&req.BaseSearchConditionQuery)
	//data, total := 1,0
	data := make([]string, 0)
	total := 0
	pageResult := httphelper.NewPageResult(data, uint32(total), 1, 10)
	httphelper.ResultSuccess(pageResult, c.Writer, lc)
}
