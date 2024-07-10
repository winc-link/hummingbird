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
	"context"
	"fmt"
	"github.com/winc-link/hummingbird/internal/dtos"
	"github.com/winc-link/hummingbird/internal/models"
	"github.com/winc-link/hummingbird/internal/pkg/constants"
	"github.com/winc-link/hummingbird/internal/pkg/errort"
)

func (app *driverLibApp) AddDriverLib(ctx context.Context, dl dtos.DeviceLibraryAddRequest) error {

	dlm := app.addOrUpdateSupportVersionConfig(dtos.FromDeviceLibraryRpcToModel(&dl))
	_, err := app.createDriverLib(dlm)
	if err != nil {
		return err
	}
	return nil
}

func (app *driverLibApp) createDriverLib(dl models.DeviceLibrary) (models.DeviceLibrary, error) {
	return app.dbClient.AddDeviceLibrary(dl)
}

func (app *driverLibApp) addOrUpdateSupportVersionConfig(dl models.DeviceLibrary) models.DeviceLibrary {
	for _, sv := range dl.SupportVersions {
		if sv.Version == dl.Version {
			//dl.SupportVersions[i].ConfigJson = dl.Config
			return dl
		}
	}
	if dl.Version == "" {
		return dl
	}

	// add
	version := models.SupportVersion{
		Version: dl.Version,
	}
	dl.SupportVersions = []models.SupportVersion{version}
	return dl
}

func (app *driverLibApp) DeviceLibrariesSearch(ctx context.Context, req dtos.DeviceLibrarySearchQueryRequest) ([]models.DeviceLibrary, uint32, error) {
	offset, limit := req.BaseSearchConditionQuery.GetPage()

	installingIds := app.manager.FilterState(constants.OperateStatusInstalling)
	if req.DownloadStatus == constants.OperateStatusInstalling && len(installingIds) == 0 {
		app.lc.Infof("deviceLibSearch install status is empty, req: %+v", req)
		return []models.DeviceLibrary{}, 0, nil
	}

	req = app.prepareSearch(req, installingIds)

	deviceLibraries, total, err := app.dbClient.DeviceLibrariesSearch(offset, limit, req)
	if err != nil {
		return deviceLibraries, 0, err
	}

	dlIds := make([]string, 0)
	dlMapDsExist := make(map[string]bool)
	for _, v := range deviceLibraries {
		dlIds = append(dlIds, v.Id)
		dlMapDsExist[v.Id] = false
	}
	dss, _, err := app.dbClient.DeviceServicesSearch(0, -1, dtos.DeviceServiceSearchQueryRequest{DeviceLibraryIds: dtos.ApiParamsArrayToString(dlIds)})
	if err != nil {
		return deviceLibraries, 0, err
	}
	for _, v := range dss {
		dlMapDsExist[v.DeviceLibraryId] = true
	}

	// 设置驱动库状态
	for i, dl := range deviceLibraries {
		stats := app.manager.GetState(dl.Id)
		if stats == constants.OperateStatusInstalling {
			deviceLibraries[i].OperateStatus = constants.OperateStatusInstalling
		} else if app.manager.ExistImage(dl.DockerImageId) && dlMapDsExist[dl.Id] {
			deviceLibraries[i].OperateStatus = constants.OperateStatusInstalled
		} else {
			deviceLibraries[i].OperateStatus = constants.OperateStatusDefault
			deviceLibraries[i].DockerImageId = ""
		}
	}

	return deviceLibraries, total, nil
}

func (app *driverLibApp) prepareSearch(req dtos.DeviceLibrarySearchQueryRequest, installingIds []string) dtos.DeviceLibrarySearchQueryRequest {
	// 处理驱动安装状态
	existImages := app.manager.GetAllImages()
	if req.DownloadStatus == constants.OperateStatusInstalling {
		req.Ids = dtos.ApiParamsArrayToString(installingIds)
	} else if req.DownloadStatus == constants.OperateStatusInstalled {
		req.NoInIds = dtos.ApiParamsArrayToString(installingIds)
		req.ImageIds = dtos.ApiParamsArrayToString(existImages)
	} else if req.DownloadStatus == constants.OperateStatusUninstall || req.DownloadStatus == constants.OperateStatusDefault {
		req.NoInIds = dtos.ApiParamsArrayToString(installingIds)
		req.NoInImageIds = dtos.ApiParamsArrayToString(existImages)
	}

	return req
}

func (app *driverLibApp) DeleteDeviceLibraryById(ctx context.Context, id string) error {
	dl, err := app.dbClient.DeviceLibraryById(id)
	if err != nil {
		return err
	}

	// 内置驱动市场不允许删除
	if dl.IsInternal {
		return errort.NewCommonErr(errort.DeviceLibraryNotAllowDelete, fmt.Errorf("internal library not allow delete"))
	}

	// 删除驱动前需要查看 驱动所属驱动实例是否存在
	_, total, edgeXErr := app.getDriverServiceApp().Search(ctx, dtos.DeviceServiceSearchQueryRequest{DeviceLibraryId: id})
	if edgeXErr != nil {
		return edgeXErr
	}
	if total > 0 {
		return errort.NewCommonErr(errort.DeviceLibraryMustDeleteDeviceService, fmt.Errorf("must delete service"))
	}

	app.manager.Remove(id)

	return nil
}

func (app *driverLibApp) DriverLibById(dlId string) (models.DeviceLibrary, error) {
	dl, err := app.dbClient.DeviceLibraryById(dlId)
	if err != nil {
		app.lc.Errorf("DriverLibById req DeviceLibraryById(%s) err %v", dlId, err)
		return models.DeviceLibrary{}, err
	}
	return dl, nil
}

func (app *driverLibApp) DeviceLibraryById(ctx context.Context, id string) (models.DeviceLibrary, error) {
	dl, err := app.DriverLibById(id)
	if err != nil {
		return models.DeviceLibrary{}, err
	}
	dl.OperateStatus = app.manager.GetState(id)

	return dl, nil
}

// UpgradeDeviceLibrary 升级驱动库版本
func (app *driverLibApp) UpgradeDeviceLibrary(ctx context.Context, req dtos.DeviceLibraryUpgradeRequest) error {
	if app.manager.GetState(req.Id) == constants.OperateStatusInstalling {
		return errort.NewCommonErr(errort.DeviceLibraryUpgradeIng, fmt.Errorf("is upgrading, please wait"))
	}

	//检查驱动是否存在
	_, edgeXErr := app.DriverLibById(req.Id)
	if edgeXErr != nil {
		return edgeXErr
	}

	if err := app.manager.Upgrade(req.Id, req.Version); err != nil {
		return err
	}

	return nil
}

func (app *driverLibApp) UpdateDeviceLibrary(ctx context.Context, update dtos.UpdateDeviceLibrary) error {
	dl, dbErr := app.dbClient.DeviceLibraryById(update.Id)
	if dbErr != nil {
		return dbErr
	}

	dtos.ReplaceDeviceLibraryModelFieldsWithDTO(&dl, update)

	dl = app.addOrUpdateSupportVersionConfig(dl)

	if err := app.dbClient.UpdateDeviceLibrary(dl); err != nil {
		return err
	}

	app.lc.Infof("deviceLibrary(%v) update succ. ", dl.Id)

	return nil
}
