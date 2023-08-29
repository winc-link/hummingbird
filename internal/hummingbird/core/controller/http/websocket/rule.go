package websocket

import (
	"github.com/winc-link/hummingbird/internal/dtos"
	"github.com/winc-link/hummingbird/internal/pkg/httphelper"
)

type WsData struct {
	Code dtos.WsCode `json:"code"`
	Data interface{} `json:"data"`
}

type WsResponse struct {
	Code dtos.WsCode               `json:"code"`
	Data httphelper.CommonResponse `json:"data"`
}

// 前端websockets
// 前端请求处理
type wsFunc func(*wsClient, interface{}, dtos.WsCode)

var wsFuncMap = map[dtos.WsCode]wsFunc{
	//驱动相关
	dtos.WsCodeDeviceLibraryUpgrade:   DeviceLibraryUpgrade,
	dtos.WsCodeDeviceServiceRunStatus: DeviceServiceRunStatus,
	dtos.WsCodeDeviceServiceLog:       DeviceServiceLog,
	dtos.WsCodeDeviceLibraryDelete:    DeviceLibraryDelete,

	//多语言
	dtos.WsCodeCheckLang: CheckLang,
}
