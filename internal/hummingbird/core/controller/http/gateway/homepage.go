/*******************************************************************************
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

// @Tags 首页
// @Summary 首页
// @Produce json
// @Param   request query   dtos.HomePageRequest true "参数"
// @Success 200     {object} httphelper.ResPageResult
// @Router  /api/v1/homepage [get]
// @Security ApiKeyAuth
func (ctl *controller) HomePage(c *gin.Context) {
	lc := ctl.lc
	var req dtos.HomePageRequest
	data, edgeXErr := ctl.getHomePageApp().HomePageInfo(c, req)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(data, c.Writer, lc)
}
