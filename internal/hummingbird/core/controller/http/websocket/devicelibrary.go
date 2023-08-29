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
receive: {"code":10001,"data":{"id":"3208327514","version":"2.0.1"}}
*/
func DeviceLibraryUpgrade(c *wsClient, data interface{}, code dtos.WsCode) {
	var req dtos.DeviceLibraryUpgradeRequest
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

	driverApp := container.DriverAppFrom(c.dic.Get)
	errCode := errort.DefaultSuccess
	err = driverApp.UpgradeDeviceLibrary(c.ctx, req)
	if err != nil {
		c.lc.Errorf("DeviceLibraryUpgrade err: %+v", err)
		errCode = errort.NewCommonEdgeXWrapper(err).Code()
	}

	// 响应
	dlName := ""
	resp := dtos.DeviceLibraryUpgradeResponse{
		Id: req.Id,
	}
	dl, err := driverApp.DeviceLibraryById(c.ctx, req.Id)
	if err != nil {
		c.lc.Errorf("get DeviceLibraryById err:%+v", err)
	} else {
		dlName = dl.Name
		resp.Version = dl.Version
		resp.OperateStatus = dl.OperateStatus
	}

	isSuccess := true
	errMsg := ""
	successMsg := ""
	status := i18n.DefaultSuccess
	if errCode != errort.DefaultSuccess {
		errMsg = i18n.TransCode(c.ctx, errCode, nil)
		isSuccess = false
		status = i18n.DefaultFail
	}

	msg := i18n.Trans(i18n.GetLang(c.ctx), i18n.LibraryUpgradeDownloadResp, map[string]interface{}{
		"name":   dlName,
		"status": i18n.Trans(i18n.GetLang(c.ctx), status, nil),
	})
	if isSuccess {
		successMsg = msg
	} else {
		errMsg = msg + ": " + errMsg
	}
	time.Sleep(2 * time.Second)
	c.sendData(code, httphelper.WsResult(errCode, resp, errMsg, successMsg))
}
