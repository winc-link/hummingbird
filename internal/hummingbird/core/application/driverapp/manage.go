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
package driverapp

import (
	"fmt"
	"github.com/winc-link/hummingbird/internal/dtos"
	"github.com/winc-link/hummingbird/internal/hummingbird/core/container"
	resourceContainer "github.com/winc-link/hummingbird/internal/hummingbird/core/container"
	interfaces "github.com/winc-link/hummingbird/internal/hummingbird/core/interface"
	"github.com/winc-link/hummingbird/internal/models"
	"github.com/winc-link/hummingbird/internal/pkg/constants"
	pkgcontainer "github.com/winc-link/hummingbird/internal/pkg/container"
	"github.com/winc-link/hummingbird/internal/pkg/di"
	"github.com/winc-link/hummingbird/internal/pkg/errort"
	"github.com/winc-link/hummingbird/internal/pkg/logger"
	"sync"
)

type DeviceLibraryManager interface {
	GetState(dlId string) string
	SetState(dlId, state string)
	FilterState(state string) []string
	Upgrade(dlId, version string) error
	Remove(dlId string) error
	//
	GetAllImages() []string // 获取所有安装镜像ID
	ExistImage(dockerImageId string) bool
}

type deviceLibraryManager struct {
	libs sync.Map

	dic *di.Container
	lc  logger.LoggingClient

	driverApp interfaces.DriverLibApp
	appModel  interfaces.DMI
}

func newDriverLibManager(dic *di.Container, app interfaces.DriverLibApp) *deviceLibraryManager {
	return &deviceLibraryManager{
		libs:      sync.Map{},
		dic:       dic,
		lc:        pkgcontainer.LoggingClientFrom(dic.Get),
		appModel:  interfaces.DMIFrom(dic.Get),
		driverApp: app,
	}
}

func (m *deviceLibraryManager) GetState(dlId string) string {
	state, ok := m.libs.Load(dlId)
	if ok {
		return state.(string)
	}
	m.libs.Store(dlId, constants.OperateStatusDefault)
	return constants.OperateStatusDefault
}

func (m *deviceLibraryManager) SetState(dlId, state string) {
	m.libs.Store(dlId, state)
}

func (m *deviceLibraryManager) FilterState(state string) []string {
	var list []string
	m.libs.Range(func(key, value interface{}) bool {
		if value.(string) == state {
			list = append(list, key.(string))
		}
		return true
	})

	return list
}

func (m *deviceLibraryManager) Remove(dlId string) error {
	dbClient := container.DBClientFrom(m.dic.Get)
	dl, err := dbClient.DeviceLibraryById(dlId)
	if err != nil {
		return err
	}

	// 删除自定义驱动
	if err := dbClient.DeleteDeviceLibraryById(dlId); err != nil {
		return err
	}

	m.asyncRemoveImage(dl)
	return nil
}

func (m *deviceLibraryManager) GetAllImages() []string {
	return m.appModel.GetAllApp()
}

func (m *deviceLibraryManager) ExistImage(dockerImageId string) bool {
	return m.appModel.StateApp(dockerImageId)
}

// updateDL 下载新版本镜像，并更新驱动库版本信息
func (m *deviceLibraryManager) updateDL(dlId, updateVersion string) error {
	// 获取驱动库信息 Refactor:
	dl, dc, err := m.driverApp.GetDeviceLibraryAndMirrorConfig(dlId)
	if err != nil {
		m.SetState(dlId, constants.OperateStatusDefault)
		return err
	}

	imageId, err := m.downloadVersion(dl, dc, updateVersion)
	if err != nil {
		return err
	}

	old := dl
	// 添加驱动库版本
	dl.DockerImageId = imageId
	newVersionInfo, isExistVersion := getNewSupportVersion(dl.SupportVersions, dl.Version, updateVersion)
	if !isExistVersion {
		dl.SupportVersions = append(dl.SupportVersions, newVersionInfo)
	}
	dl = m.updateDLDefaultVersion(dl, newVersionInfo)

	dbClient := resourceContainer.DBClientFrom(m.dic.Get)
	if err := dbClient.UpdateDeviceLibrary(dl); err != nil {
		m.lc.Errorf("updateDeviceLibrary %s fail %+v", dl.Id, err)
		return err
	}

	m.cleanOldVersion(old, dl)

	return nil
}

// cleanOldVersion 异步清理旧版本镜像
func (m *deviceLibraryManager) cleanOldVersion(oldDL, dl models.DeviceLibrary) {
	if oldDL.DockerImageId == dl.DockerImageId {
		return
	}
	m.asyncRemoveImage(oldDL)
}

func (m *deviceLibraryManager) asyncRemoveImage(dl models.DeviceLibrary) {
	// 镜像删除
	go m.appModel.RemoveApp(dtos.DeviceLibraryFromModel(dl))
}

func getNewSupportVersion(versions models.SupportVersions, curVersion, newVersion string) (models.SupportVersion, bool) {
	newVersionInfo := models.SupportVersion{}
	for _, v := range versions {
		if v.Version == newVersion {
			return v, true
		}
		// 将老版本的配置复制到新版本中
		if v.Version == curVersion {
			newVersionInfo = v
		}
	}

	newVersionInfo.Version = newVersion
	return newVersionInfo, false
}

func (m *deviceLibraryManager) updateDLDefaultVersion(dl models.DeviceLibrary, newVersion models.SupportVersion) models.DeviceLibrary {
	dl.Version = newVersion.Version
	return dl
}

func (m *deviceLibraryManager) downloadVersion(dl models.DeviceLibrary, dc models.DockerConfig, version string) (string, error) {
	// 3.下载应用
	var cfg dtos.DockerConfig
	if dl.IsInternal {
		cfg = dtos.WincLinkDockerConfig()
	} else {
		cfg = dtos.DockerConfigFromModel(dc)
	}
	imageId, err := m.appModel.DownApp(cfg, dtos.DeviceLibraryFromModel(dl), version)
	if err != nil {
		return "", err
	}
	return imageId, nil
}

func (m *deviceLibraryManager) Upgrade(dlId, updateVersion string) error {
	if m.GetState(dlId) == constants.OperateStatusInstalling {
		return errort.NewCommonErr(errort.DeviceLibraryUpgradeIng, fmt.Errorf("device library upgradeing"))
	}

	// 1. 设置为升级中
	m.SetState(dlId, constants.OperateStatusInstalling)
	m.lc.Infof("1.start updateDeviceLibrary %v to version %v", dlId, updateVersion)

	// 下载新版本镜像，并更新驱动库版本信息
	if err := m.updateDL(dlId, updateVersion); err != nil {
		m.SetState(dlId, constants.OperateStatusDefault)
		m.lc.Errorf("updateDeviceLibrary version fail", err)
		return err
	}

	m.SetState(dlId, constants.OperateStatusInstalled)
	//m.lc.Infof("2.updateDeviceLibrary %v version %v", dlId, updateVersion)

	// 3. 升级驱动实例
	if err := m.upgradeDeviceService(dlId); err != nil {
		m.lc.Errorf("3.upgradeDeviceService error %v", err)
		return err
	}

	return nil
}

func (m *deviceLibraryManager) upgradeDeviceService(dlId string) error {
	dl, err := m.driverApp.DriverLibById(dlId)
	if err != nil {
		return err
	}

	return resourceContainer.DriverServiceAppFrom(m.dic.Get).Upgrade(dl)
}
