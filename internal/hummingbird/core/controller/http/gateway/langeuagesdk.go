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

func (ctl *controller) LanguageSdkSearch(c *gin.Context) {
	lc := ctl.lc
	var req dtos.LanguageSDKSearchQueryRequest
	if err := c.ShouldBind(&req); err != nil {
		httphelper.RenderFail(c, errort.NewCommonErr(errort.DefaultReqParamsError, err), c.Writer, lc)
		return
	}
	dtos.CorrectionPageParam(&req.BaseSearchConditionQuery)

	list, _, edgeXErr := ctl.getLanguageApp().LanguageSDKSearch(c, req)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	responseData := make(map[string]interface{})
	responseData["doc"] = map[string]interface{}{
		"name": "物联网平台文档",
		"icon": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAOEAAADhCAMAAAAJbSJIAAAAflBMVEX////2iSnRTif5hAj5jjH96tv1gAj94s" +
			"7PSif4jCnSQQHOPQDQSB3129T129XdfGPSQgr2hBjggGf2hiH5iBbROwDegGnTQgrbe2XORB3OQQ/5fQD5hQv+8uf70rX6w5r6uon7y6j838n+9Ov/+v" +
			"X3mk/3kTr4pGL4qm37xqDMD8zoAAAChUlEQVR4nO3ZaVOCUBSA4cykhe5FpNVKW9Tq///BaNHRBC6KzOGced+PIAzPDHdRj466kE+kn6DdfDK8ln6GNnM3" +
			"w8vjC+mnaC83yn3HdoUuufj2mRW6ZPjrMyr0o5XPpNCt+wwKXRKv+8wJ18afSeG/99Oc0N3E2z5Dwq3xZ0xYMP5MCX3R+DMkLJxfDAlX+0+jworxZ0Lok3HIp1" +
			"pYOb8YEJauf0aEtX1KhT48v6gWluw/zQh9aP1TLnTJ7W4+ZUJXZ/1TLAzsP9UL93g/VQmD+2vlwr3GnyJhzf2nWmGD91OFcIf9p0phre9/ioU//282r7PCvdc/JcLG8" +
			"0vHhTV+X1It3HP/qUboRwcafx0VuuT6sL6OCRvtPxUID7T+dVbo7+L7s+rGcWXj4qtiadmyh+w80ONp9R0GjyXX9ev30KKw3wsVhYRR8BbB+ggRIkSIECFChAgRIkSIECFC" +
			"hAgRIkSIECFChAgRIkSIECFChAgRIkSIECFChAgRIkSIECFChAgRIkSIECFChAgRIkSIECFChAgRNhM+TSZPpoXTKO/ZsHD6w9kg2hK+/GmiqVHhZPXp6MWk8HXNEk0MCk/Tjc" +
			"Ov5oSD82zjeLo8bkX4dpL9O5ENTAlnW8Belr0ZEs7mW8CcuJjZEb6nRafS+cyCMM2Xvo9CYH7uPV8k1Qt7vem8BJgTF5+L5kBxYVowBpdlpXhNwvZDiBChfAgRIpQPIUKE8iFEiFA" +
			"+hAgRyocQIUL5ECJEKB9ChAjlQ4gQoXwIESKUDyFChPIhRIhQPoQIEcqHECFC+RAiRCgfQoQI5UOIEKF8CBEilA8hQoTyIUSIUD6ECBHKhxBhVVcnXehqp2f+Ag5ihjFgr47/AAAAAElFTkSuQmCC",
		"addr": "https://doc.hummingbird.winc-link.com/",
	}
	responseData["sdk_language"] = list
	httphelper.ResultSuccess(responseData, c.Writer, lc)
}

func (ctl *controller) LanguageSdkSync(c *gin.Context) {
	lc := ctl.lc
	var req dtos.LanguageSDKSyncRequest
	if err := c.ShouldBind(&req); err != nil {
		httphelper.RenderFail(c, errort.NewCommonErr(errort.DefaultReqParamsError, err), c.Writer, lc)
		return
	}
	edgeXErr := ctl.getLanguageApp().Sync(c, req.VersionName)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(nil, c.Writer, lc)
}
