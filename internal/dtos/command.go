package dtos

import (
	"encoding/json"
	"fmt"
)

// DebugAssistantReq 调试助手请求参数 TODO 可以直接使用 DpMessage
type DebugAssistantReq struct {
	DeviceId string                 `json:"deviceId,omitempty"`
	OpType   int32                  `json:"opType,omitempty"`
	Data     map[string]interface{} `json:"data" binding:"required"`
	Protocol int32                  `json:"protocol" binding:"required"`
	S        int64                  `json:"s"`
	T        int64                  `json:"t" binding:"required"`
}

func (r DebugAssistantReq) DataString() string {
	body, _ := json.Marshal(r)
	return string(body)
}

// 北向指令
type CmdRequest struct {
	Cid      string
	Protocol int32
	S        int64
	T        int64
	Data     []byte // json encode
}

func (cr CmdRequest) String() string {
	return fmt.Sprintf("cid: %s, protocol: %d, s: %d, t: %d, data: %s", cr.Cid, cr.Protocol, cr.S, cr.T, string(cr.Data))
}

type CommandResponse struct {
	Id        string                 `json:"id"`       // uuid
	Cid       string                 `json:"cid"`      // 设备ID
	Protocol  int32                  `json:"protocol"` // 协议号
	Timestamp int64                  `json:"timestamp"`
	Data      map[string]interface{} `json:"data"` // 序列化后的数据
}

func (cmd CommandResponse) DataJSON() string {
	b, _ := json.Marshal(cmd.Data)
	return string(b)
}

type CommandQueryRequest struct {
	DeviceId string `form:"device_id" binding:"required"`
}

type ListCommandResponse struct {
	List []CommandResponse `json:"list"`
}

func NewListCommandResponse() ListCommandResponse {
	return ListCommandResponse{
		List: []CommandResponse{},
	}
}
