package websocket

import (
	"encoding/json"
	"github.com/gin-gonic/gin/binding"
	"github.com/winc-link/hummingbird/internal/dtos"
	"github.com/winc-link/hummingbird/internal/hummingbird/core/container"
	"github.com/winc-link/hummingbird/internal/pkg/errort"
	"github.com/winc-link/hummingbird/internal/pkg/httphelper"
	"github.com/winc-link/hummingbird/internal/pkg/i18n"
	"time"
)

/**
receive: {"code":10002,"data":{"id":"935769","run_status":1}}
*/
func DeviceServiceRunStatus(c *wsClient, data interface{}, code dtos.WsCode) {
	var req dtos.UpdateDeviceServiceRunStatusRequest
	bytes, _ := json.Marshal(data)
	err := binding.JSON.BindBody(bytes, &req)
	if err != nil {
		c.lc.Error(err.Error())
		c.sendData(
			code,
			httphelper.WsResultFail(
				errort.DefaultReqParamsError,
				i18n.TransCode(c.ctx, errort.DefaultReqParamsError, nil),
			),
		)
		return
	}

	driverServerApp := container.DriverServiceAppFrom(c.dic.Get)
	err = driverServerApp.UpdateRunStatus(c.ctx, req)

	errCode := errort.DefaultSuccess
	errMsg := ""
	isSuccess := true
	successMsg := ""
	status := i18n.DefaultSuccess

	if err != nil {
		c.lc.Errorf("DeviceServiceRunStatus err: %+v", err)
		edgeX := errort.NewCommonEdgeXWrapper(err)
		errCode = edgeX.Code()
		errMsg = edgeX.Error()
	}

	// 响应
	dsName := ""
	resp := dtos.UpdateDeviceServiceRunStatusResponse{
		Id: req.Id,
	}
	ds, err := driverServerApp.Get(c.ctx, req.Id)
	if err != nil {
		c.lc.Errorf("get driverServer err:%+v", err)
	} else {
		resp.RunStatus = ds.RunStatus
	}

	if errCode != errort.DefaultSuccess {
		if errCode != errort.ContainerRunFail {
			errMsg = i18n.TransCode(c.ctx, errCode, nil)
		}
		isSuccess = false
		status = i18n.DefaultFail
	}

	msg := i18n.Trans(i18n.GetLang(c.ctx), i18n.ServiceRunStatusResp, map[string]interface{}{
		"name":   dsName,
		"status": i18n.Trans(i18n.GetLang(c.ctx), status, nil),
	})
	if isSuccess {
		successMsg = msg
	}
	time.Sleep(2 * time.Second)
	c.sendData(code, httphelper.WsResult(errCode, resp, errMsg, successMsg))
}
