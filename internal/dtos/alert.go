package dtos

import (
	"github.com/winc-link/hummingbird/internal/models"
	"github.com/winc-link/hummingbird/internal/pkg/constants"
	"github.com/winc-link/hummingbird/internal/pkg/errort"
	"github.com/winc-link/hummingbird/internal/pkg/i18n"
	"time"
)

const (
	SYSTEM_ALERT  models.AlertType = iota + 1 // 系统告警
	DRIVER_ALERT                              // 驱动告警
	LICENSE_ALERT                             //证书过期
)

const (
	ERROR  models.AlertLevel = iota + 1 // 告警级别：错误
	WARN                                // 告警级别：警告
	NOTIFY                              // 告警级别： 通知
)

var (
	AlertTypeTrans = map[models.AlertType]string{
		SYSTEM_ALERT:  i18n.AgentAlertSystem,
		DRIVER_ALERT:  i18n.AgentAlertDriver,
		LICENSE_ALERT: i18n.LicenseAlertExpire,
	}
	AlertLevelTrans = map[models.AlertLevel]string{
		NOTIFY: i18n.AgentAlertNotify,
		WARN:   i18n.AgentAlertWarn,
		ERROR:  i18n.AgentAlertError,
	}
)

// AlertContent 服务和驱动上报告警消息
type (
	ReportAlertsReq struct {
		BaseRequest `json:",inline"`
		ServiceName string            `json:"name"`                        // 服务名
		Type        models.AlertType  `json:"type" binding:"oneof=1 2"`    // 告警类型
		Level       models.AlertLevel `json:"level" binding:"oneof=1 2 3"` // 告警级别
		T           int64             `json:"time"`                        // 告警时间
		Content     string            `json:"content"`
	}

	AlertContentDTO struct {
		ServiceName string            `json:"name"`                                           // 服务名
		Type        models.AlertType  `json:"type" binding:"oneof=1 2" swaggertype:"integer"` // 告警类型
		TypeValue   string            `json:"typeValue"`
		Level       models.AlertLevel `json:"level" binding:"oneof=1 2 3" swaggertype:"integer"` // 告警级别
		LevelValue  string            `json:"levelValue"`
		T           int64             `json:"time"`    // 告警时间
		Content     string            `json:"content"` // 告警内容
	}
)

func NewReportAlertsReq(serviceName string, tp models.AlertType, l models.AlertLevel, t int64, content string) ReportAlertsReq {
	return ReportAlertsReq{
		BaseRequest: NewBaseRequest(),
		ServiceName: serviceName,
		Type:        tp,
		Level:       l,
		T:           t,
		Content:     content,
	}
}

func ToAlertContent(req ReportAlertsReq) models.AlertContent {
	return models.AlertContent{
		ServiceName: req.ServiceName,
		Type:        req.Type,
		Level:       req.Level,
		T:           req.T,
		Content:     req.Content,
	}
}

func AlertContentToDTO(ac models.AlertContent) AlertContentDTO {
	return AlertContentDTO{
		ServiceName: ac.ServiceName,
		Type:        ac.Type,
		Level:       ac.Level,
		T:           ac.T,
		Content:     ac.Content,
	}
}

type ReportAlertRequest struct {
	ServiceName string `json:"serviceName"`
	AlertType   int    `json:"alertType"`  // constants.AlertType_SERVICE
	AlertLevel  int    `json:"alertLevel"` // constants.AlertLevel_ERROR
	AlertTime   int64  `json:"alertTime"`
	Content     string `json:"content"`
}

// GenServerAlert 生成服务警告内容
func GenServerAlert(lvl models.AlertLevel, err error) ReportAlertsReq {
	errw := errort.NewCommonEdgeXWrapper(err)
	return NewReportAlertsReq(
		constants.CoreServiceKey,
		SYSTEM_ALERT,
		lvl,
		time.Now().Unix(),
		i18n.TransCodeDefault(errw.Code(), nil))
}
