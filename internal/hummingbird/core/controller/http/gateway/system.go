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
	"github.com/winc-link/hummingbird/internal/pkg/errort"
	"github.com/winc-link/hummingbird/internal/pkg/httphelper"
	"github.com/winc-link/hummingbird/internal/pkg/utils"
	"path"
)

// @Tags 网关管理
// @Summary 网关备份下载
// @Produce json
// @Success 200 {object} httphelper.CommonResponse
// @Router /api/v1/system/backup [get]
func (ctl *controller) SystemBackupHandle(c *gin.Context) {
	lc := ctl.lc
	filePath, edgeXErr := ctl.getSystemApp().SystemBackupFileDownload(c)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}

	fileName := path.Base(filePath)
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.File(filePath)

	// 删除zip文件
	utils.RemoveFileOrDir(filePath)
}

func (ctl *controller) SystemRecoverHandle(c *gin.Context) {
	lc := ctl.lc
	file, err := c.FormFile("fileName")
	if err != nil {
		httphelper.RenderFail(c, errort.NewCommonErr(errort.DefaultReqParamsError, err), c.Writer, lc)
		return
	}
	dist := "/tmp/tedge-recover.zip"
	err = c.SaveUploadedFile(file, dist)
	if err != nil {
		httphelper.RenderFail(c, errort.NewCommonErr(errort.SystemErrorCode, err), c.Writer, lc)
		return
	}
	edgeXErr := ctl.getSystemApp().SystemRecover(c, dist)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}

	httphelper.ResultSuccess(nil, c.Writer, lc)
}
