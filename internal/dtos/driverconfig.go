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
package dtos

import (
	"github.com/winc-link/hummingbird/internal/pkg/constants"
	"github.com/winc-link/hummingbird/internal/pkg/utils"
	"sync"
)

var dcMutex sync.Mutex

type DriverConfigManage struct {
	HostRootDir string // 宿主主机dir
	NetWorkName string
	DockerManageConfig
}

type DockerManageConfig struct {
	ContainerConfigPath string // 容器内部驱动运行启动的配置文件路径
	DockerApiVersion    string // docker 版本号
	DockerRunMode       string // 默认运行模式 默认是host
	DockerSelfName      string // edge启动的容器名
	Privileged          bool   // 容器特权
}

func (m *DriverConfigManage) SetHostRootDir(dir string) {
	dcMutex.Lock()
	defer dcMutex.Unlock()
	m.HostRootDir = dir
}

func (m *DriverConfigManage) SetNetworkName(networkName string) {
	dcMutex.Lock()
	defer dcMutex.Unlock()
	m.NetWorkName = networkName
}

// 存储驱动上传配置定义文件目录  /var/tedge/edgex-driver-data/driver-library/
func (m *DriverConfigManage) GetLibraryDir() string {
	return utils.GetPwdDir() + "/" + constants.DriverLibraryDir + "/"
}

// 驱动二进制文件路径 /var/tedge/edgex-driver-data/bin/modbus-1234
func (m *DriverConfigManage) GetBinPath(serverName string) string {
	return utils.GetPwdDir() + "/" + constants.DriverBinDir + "/" + serverName
}

// 驱动启动的配置文件路径 /var/edge/run-config/modbus-1234.toml
func (m *DriverConfigManage) GetRunConfigPath(serviceName string) string {
	return constants.DockerHummingbirdRootDir + "/" + constants.DriverRunConfigDir + "/" + serviceName + constants.ConfigSuffix
}

// docker挂载
func (m *DriverConfigManage) GetHostRunConfigPath(serviceName string) string {
	return m.HostRootDir + "/" + constants.DriverRunConfigDir + "/" + serviceName + constants.ConfigSuffix
}

// 二进制版本路径
// 驱动启动的配置文件路径 /var/edge/mnt/modbus-1234.toml
func (m *DriverConfigManage) GetMntDir(serviceName string) string {
	return constants.DockerHummingbirdRootDir + "/" + constants.DriverMntDir + "/" + serviceName
}

// docker挂载 的日志：只针对docker版本，二进制版本需要改动日志存储地址 /var/edge/mnt/modbus-1234.toml
func (m *DriverConfigManage) GetHostMntDir(serviceName string) string {
	return m.HostRootDir + "/" + constants.DriverMntDir + "/" + serviceName
}

// 二进制版本 驱动运行日志文件 /var/tedge/mnt/modbus-1234/logs/driver.log
func (m *DriverConfigManage) GetLogFilePath(serviceName string) string {
	return utils.GetPwdDir() + "/" + constants.DriverMntDir + "/" + serviceName + "/" + constants.DriverDefaultLogPath
}

// docker挂载
//logfilePath = "/var/edge/edge-driver-data/mnt/aliyun-iot/edgex-aliyun-cloud.log"

///var/edge/edge-driver-data/mnt/aliyun-iot
func (m *DriverConfigManage) GetHostLogFilePath(serviceName string) string {
	return constants.DockerHummingbirdRootDir + "/" + constants.DriverMntDir + "/" + serviceName + "/" + constants.DriverDefaultLogPath
}
