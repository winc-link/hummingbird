/*******************************************************************************
 * Copyright 2017 Dell Inc.
 * Copyright (c) 2019 Intel Corporation
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
package interfaces

import (
	"github.com/winc-link/hummingbird/internal/dtos"
	"github.com/winc-link/hummingbird/internal/pkg/di"
)

var (
	DriverModelInterfaceName = di.TypeInstanceToName((*DMI)(nil))
)

func DMIFrom(get di.Get) DMI {
	return get(DriverModelInterfaceName).(DMI)
}

type DMI interface {
	DriverInInstanceDMI
	StopAllInstance()
}

//驱动相关接口
type DriverInInstanceDMI interface {
	// DownApp 下载驱动
	DownApp(cfg dtos.DockerConfig, app dtos.DeviceLibrary, toVersion string) (string, error)

	RemoveApp(app dtos.DeviceLibrary) error
	GetAllApp() []string
	// 检查驱动软件情况
	StateApp(dockerImageId string) bool

	InstanceState(ins dtos.DeviceService) bool
	// StartInstance 启动实例
	StartInstance(ins dtos.DeviceService, cfg dtos.RunServiceCfg) (string, error) // 返回服务所在的ip
	// StopInstance 停止实例
	StopInstance(ins dtos.DeviceService) error
	// DeleteInstance 删除实例
	DeleteInstance(ins dtos.DeviceService) error

	GetDriverInstanceLogPath(serviceName string) string
	// GetSelfIp 获取当前服务运行的内网ip
	GetSelfIp() string
}
