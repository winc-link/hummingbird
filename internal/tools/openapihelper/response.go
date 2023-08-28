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

package openapihelper

import (
	"github.com/gin-gonic/gin"
	"github.com/winc-link/hummingbird/internal/pkg/errort"
	"net/http"
	"time"
)

type CommonOpenResponse struct {
	Code    errort.OpenApiErrorCode `json:"code"`
	Msg     errort.OpenApiErrorMsg  `json:"msg,omitempty"`
	Success bool                    `json:"success"`
	// 返回13位时间戳 time.Now.UnixNano / 1e6
	T      int64       `json:"t"`
	Result interface{} `json:"result,omitempty"`
}

func ReaderSuccess(c *gin.Context, data interface{}) {
	if data == nil {
		data = []interface{}{}
	}
	resp := CommonOpenResponse{
		Result:  data,
		T:       time.Now().UnixNano() / 1e6,
		Success: true,
	}

	c.JSON(http.StatusOK, resp)
}

func ReaderFail(c *gin.Context, code errort.OpenApiErrorCode) {
	resp := CommonOpenResponse{
		Success: false,
		Msg:     errort.OpenApiCodeMsgMap[code],
		Code:    code,
		T:       time.Now().UnixNano() / 1e6,
	}
	c.JSON(http.StatusOK, resp)
}
