//
// Copyright (C) 2021 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package models

import (
	"encoding/json"
)

// 服务类型，服务，驱动
const (
	ServiceTypeEnumService      = "service"
	ServiceTypeEnumDriver       = "driver"
	ServiceTypeEnumThirdPartApp = "appService"
)

// 使用 redis hash 存储，name 做 field

type ServiceStats struct {
	Id          string `gorm:"column:id" json:"id"`        // 服务标识
	Name        string `gorm:"column:name;pk" json:"name"` // 容器名称
	LogPath     string `json:"log_path"`                   // 日志地址, 宿主主机日志路径
	ServiceType string `json:"service_type"`               // 服务类型，服务，驱动，应用
	Healthy     bool   `json:"healthy"`                    // 状态, docker.inspect.State.Running
	Created     string `json:"created"`                    // 创建时间, docker.inspect.Created
	Started     string `json:"started"`                    // 启动时间, docker.inspect.State.StartedAt
}

func (s ServiceStats) ToJSON() string {
	body, _ := json.Marshal(s)
	return string(body)
}

func (s *ServiceStats) TableName() string {
	return "service_stats"
}

func (s *ServiceStats) Get() interface{} {
	return *s
}
