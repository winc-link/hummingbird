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

//func (ctl *controller) AddJobHandle(c *gin.Context) {
//	ctl.ProxySharpServer(c)
//}
//
//func (ctl *controller) QueryJobHandle(c *gin.Context) {
//	ctl.ProxySharpServer(c)
//}
//
//func (ctl *controller) ChangeJobStatusHandle(c *gin.Context) {
//	ctl.ProxySharpServer(c)
//}
//
//func (ctl *controller) GetJobHandle(c *gin.Context) {
//	ctl.ProxySharpServer(c)
//}
//
//func (ctl *controller) UpdateJobHandle(c *gin.Context) {
//	ctl.ProxySharpServer(c)
//}
//
//func (ctl *controller) DeleteJobHandle(c *gin.Context) {
//	ctl.ProxySharpServer(c)
//}
//
//func (ctl *controller) CheckJobExistByDeviceIdOrSceneIdHandle(c *gin.Context) {
//	ctl.ProxySharpServer(c)
//}
//
//func (ctl *controller) QueryJobLogsHandle(c *gin.Context) {
//	ctl.ProxySharpServer(c)
//}
//
//func (ctl *controller) ExecJobHandle(c *gin.Context) {
//	lc := ctl.lc
//	var req dtos.JobAction
//	if err := c.ShouldBind(&req); err != nil {
//		httphelper.RenderFail(c, errort.NewCommonErr(errort.DefaultReqParamsError, err), c.Writer, lc)
//		return
//	}
//
//	//data, total, edgeXErr := ctl.getDeviceApp().DeviceAction(c, req)
//	//if edgeXErr != nil {
//	//	httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
//	//	return
//	//}
//	//pageResult := httphelper.NewPageResult(data, total, req.Page, req.PageSize)
//	httphelper.ResultSuccess(nil, c.Writer, lc)
//}
