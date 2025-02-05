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
package deviceapp

import (
	"context"
	"fmt"
	"github.com/winc-link/hummingbird/internal/dtos"
	resourceContainer "github.com/winc-link/hummingbird/internal/hummingbird/core/container"
	interfaces "github.com/winc-link/hummingbird/internal/hummingbird/core/interface"
	"github.com/winc-link/hummingbird/internal/models"
	"github.com/winc-link/hummingbird/internal/pkg/constants"
	"github.com/winc-link/hummingbird/internal/pkg/container"
	"github.com/winc-link/hummingbird/internal/pkg/di"
	"github.com/winc-link/hummingbird/internal/pkg/errort"
	"github.com/winc-link/hummingbird/internal/pkg/logger"
	"github.com/winc-link/hummingbird/internal/pkg/utils"
	"time"
)

type deviceApp struct {
	//*propertyTyApp
	dic      *di.Container
	dbClient interfaces.DBClient
	lc       logger.LoggingClient
}

func NewDeviceApp(ctx context.Context, dic *di.Container) interfaces.DeviceItf {
	lc := container.LoggingClientFrom(dic.Get)
	dbClient := resourceContainer.DBClientFrom(dic.Get)

	return &deviceApp{
		dic:      dic,
		dbClient: dbClient,
		lc:       lc,
	}
}

func (p deviceApp) DeviceById(ctx context.Context, id string) (dtos.DeviceInfoResponse, error) {
	device, err := p.dbClient.DeviceById(id)
	var response dtos.DeviceInfoResponse
	if err != nil {
		return response, err
	}
	var deviceServiceName string
	deviceService, err := p.dbClient.DeviceServiceById(device.DriveInstanceId)
	if err != nil {
		deviceServiceName = deviceService.Name
	}

	response = dtos.DeviceInfoResponseFromModel(device, deviceServiceName)
	return response, nil
}

func (p deviceApp) OpenApiDeviceById(ctx context.Context, id string) (dtos.OpenApiDeviceInfoResponse, error) {
	device, err := p.dbClient.DeviceById(id)
	var response dtos.OpenApiDeviceInfoResponse
	if err != nil {
		return response, err
	}
	response = dtos.OpenApiDeviceInfoResponseFromModel(device)
	return response, nil
}

func (p deviceApp) OpenApiDeviceStatusById(ctx context.Context, id string) (dtos.OpenApiDeviceStatus, error) {
	device, err := p.dbClient.DeviceById(id)
	var response dtos.OpenApiDeviceStatus
	if err != nil {
		return response, err
	}
	response.Status = device.Status
	return response, nil
}

func (p deviceApp) DeviceByCloudId(ctx context.Context, id string) (models.Device, error) {
	return p.dbClient.DeviceByCloudId(id)
}

func (p deviceApp) DeviceModelById(ctx context.Context, id string) (models.Device, error) {
	return p.dbClient.DeviceById(id)
}

func (p *deviceApp) DevicesSearch(ctx context.Context, req dtos.DeviceSearchQueryRequest) ([]dtos.DeviceSearchQueryResponse, uint32, error) {
	offset, limit := req.BaseSearchConditionQuery.GetPage()
	resp, total, err := p.dbClient.DevicesSearch(offset, limit, req)
	if err != nil {
		return []dtos.DeviceSearchQueryResponse{}, 0, err
	}
	devices := make([]dtos.DeviceSearchQueryResponse, len(resp))
	for i, dev := range resp {
		deviceService, _ := p.dbClient.DeviceServiceById(dev.DriveInstanceId)
		devices[i] = dtos.DeviceResponseFromModel(dev, deviceService.Name)
	}
	return devices, total, nil
}

func (p *deviceApp) OpenApiDevicesSearch(ctx context.Context, req dtos.DeviceSearchQueryRequest) ([]dtos.OpenApiDeviceInfoResponse, uint32, error) {
	offset, limit := req.BaseSearchConditionQuery.GetPage()
	resp, total, err := p.dbClient.DevicesSearch(offset, limit, req)
	if err != nil {
		return []dtos.OpenApiDeviceInfoResponse{}, 0, err
	}
	devices := make([]dtos.OpenApiDeviceInfoResponse, len(resp))
	for i, device := range resp {
		devices[i] = dtos.OpenApiDeviceInfoResponseFromModel(device)
	}
	return devices, total, nil
}

func (p *deviceApp) DevicesModelSearch(ctx context.Context, req dtos.DeviceSearchQueryRequest) ([]models.Device, uint32, error) {
	offset, limit := req.BaseSearchConditionQuery.GetPage()
	return p.dbClient.DevicesSearch(offset, limit, req)
}

func (p *deviceApp) AddDevice(ctx context.Context, req dtos.DeviceAddRequest) (string, error) {
	if req.DriverInstanceId != "" {
		driverInstance, err := p.dbClient.DeviceServiceById(req.DriverInstanceId)
		if err != nil {
			return "", err
		}
		if driverInstance.Platform != "" && driverInstance.Platform != constants.IotPlatform_LocalIot {
			return "", errort.NewCommonErr(errort.DeviceServiceMustLocalPlatform, fmt.Errorf("please sync product data"))
		}
	}

	productInfo, err := p.dbClient.ProductById(req.ProductId)

	if productInfo.Status == constants.ProductUnRelease {
		return "", errort.NewCommonEdgeX(errort.ProductUnRelease, "The product has not been released yet. Please release the product before adding devices", nil)
	}
	deviceId := utils.RandomNum()
	if err != nil {
		return "", err
	}

	err = resourceContainer.DataDBClientFrom(p.dic.Get).CreateTable(ctx, productInfo.Id, deviceId)
	if err != nil {
		return "", err
	}
	var insertDevice models.Device
	insertDevice.Id = deviceId
	insertDevice.Name = req.Name
	insertDevice.DeviceSn = req.DeviceSn
	insertDevice.ProductId = req.ProductId
	insertDevice.Platform = constants.IotPlatform_LocalIot
	insertDevice.DriveInstanceId = req.DriverInstanceId
	insertDevice.Status = constants.DeviceStatusOffline
	insertDevice.Secret = utils.GenerateDeviceSecret(12)
	insertDevice.Description = req.Description
	id, err := p.dbClient.AddDevice(insertDevice)
	if err != nil {
		return "", err
	}
	go func() {
		p.CreateDeviceCallBack(insertDevice)
	}()
	return id, nil
}

func (p *deviceApp) BatchDeleteDevice(ctx context.Context, ids []string) error {
	var searchReq dtos.DeviceSearchQueryRequest
	searchReq.BaseSearchConditionQuery.Ids = dtos.ApiParamsArrayToString(ids)
	devices, _, err := p.dbClient.DevicesSearch(0, -1, searchReq)
	if err != nil {
		return err
	}
	alertApp := resourceContainer.AlertRuleAppNameFrom(p.dic.Get)
	for _, device := range devices {
		edgeXErr := alertApp.CheckRuleByDeviceId(ctx, device.Id)
		if edgeXErr != nil {
			return edgeXErr
		}
	}
	err = p.dbClient.BatchDeleteDevice(ids)
	if err != nil {
		return err
	}
	for _, device := range devices {
		delDevice := device
		go func() {
			p.DeleteDeviceCallBack(delDevice)
		}()
	}
	return nil
}

func (p *deviceApp) DeviceMqttAuthInfo(ctx context.Context, id string) (dtos.DeviceAuthInfoResponse, error) {
	mqttAuth, err := p.dbClient.DeviceMqttAuthInfo(id)
	var response dtos.DeviceAuthInfoResponse
	if err != nil {
		return response, err
	}
	response = dtos.DeviceAuthInfoResponseFromModel(mqttAuth)
	return response, nil
}

func (p *deviceApp) AddMqttAuth(ctx context.Context, req dtos.AddMqttAuthInfoRequest) (string, error) {
	var mqttAuth models.MqttAuth
	mqttAuth.ClientId = req.ClientId
	mqttAuth.UserName = req.UserName
	mqttAuth.Password = req.Password
	mqttAuth.ResourceId = req.ResourceId
	mqttAuth.ResourceType = constants.ResourceType(req.ResourceType)
	return p.dbClient.AddMqttAuthInfo(mqttAuth)
}

func (p *deviceApp) DeleteDeviceById(ctx context.Context, id string) error {
	deviceInfo, err := p.dbClient.DeviceById(id)
	if err != nil {
		return err
	}
	alertApp := resourceContainer.AlertRuleAppNameFrom(p.dic.Get)
	edgeXErr := alertApp.CheckRuleByDeviceId(ctx, id)
	if edgeXErr != nil {
		return edgeXErr
	}

	sceneApp := resourceContainer.SceneAppNameFrom(p.dic.Get)
	edgeXErr = sceneApp.CheckSceneByDeviceId(ctx, id)
	if edgeXErr != nil {
		return edgeXErr
	}

	err = p.dbClient.DeleteDeviceById(id)
	if err != nil {
		return err
	}
	_ = resourceContainer.DataDBClientFrom(p.dic.Get).DropTable(ctx, id)

	go func() {
		p.DeleteDeviceCallBack(models.Device{
			Id:              deviceInfo.Id,
			DriveInstanceId: deviceInfo.DriveInstanceId,
		})
	}()
	return nil
}

func (p *deviceApp) DeviceUpdate(ctx context.Context, req dtos.DeviceUpdateRequest) error {
	if req.Id == "" {
		return errort.NewCommonEdgeX(errort.DefaultReqParamsError, "update req id is required", nil)

	}
	device, edgeXErr := p.dbClient.DeviceById(req.Id)
	if edgeXErr != nil {
		return edgeXErr
	}

	if device.Platform != constants.IotPlatform_LocalIot {

	}
	alertApp := resourceContainer.AlertRuleAppNameFrom(p.dic.Get)
	edgeXErr = alertApp.CheckRuleByDeviceId(ctx, req.Id)
	if edgeXErr != nil {
		return edgeXErr
	}

	sceneApp := resourceContainer.SceneAppNameFrom(p.dic.Get)
	edgeXErr = sceneApp.CheckSceneByDeviceId(ctx, req.Id)
	if edgeXErr != nil {
		return edgeXErr
	}

	dtos.ReplaceDeviceModelFields(&device, req)
	edgeXErr = p.dbClient.UpdateDevice(device)
	if edgeXErr != nil {
		return edgeXErr
	}
	go func() {
		p.UpdateDeviceCallBack(device)
	}()
	return nil
}

func (p *deviceApp) DevicesUnBindDriver(ctx context.Context, req dtos.DevicesUnBindDriver) error {
	var searchReq dtos.DeviceSearchQueryRequest
	searchReq.BaseSearchConditionQuery.Ids = dtos.ApiParamsArrayToString(req.DeviceIds)
	var devices []models.Device
	var err error
	var total uint32
	devices, total, err = p.dbClient.DevicesSearch(0, -1, searchReq)
	if err != nil {
		return err
	}
	if total == 0 {
		return errort.NewCommonErr(errort.DeviceNotExist, fmt.Errorf("devices not found"))
	}

	err = p.dbClient.BatchUnBindDevice(req.DeviceIds)
	if err != nil {
		return err
	}
	for _, device := range devices {
		callBackDevice := device
		go func() {
			p.DeleteDeviceCallBack(models.Device{
				Id:              callBackDevice.Id,
				DriveInstanceId: callBackDevice.DriveInstanceId,
			})
		}()
	}
	return nil
}

// DevicesBindProductId
func (p *deviceApp) DevicesBindProductId(ctx context.Context, req dtos.DevicesBindProductId) error {
	var searchReq dtos.DeviceSearchQueryRequest
	searchReq.ProductId = req.ProductId
	var devices []models.Device
	var err error
	var total uint32
	devices, total, err = p.dbClient.DevicesSearch(0, -1, searchReq)
	if err != nil {
		return err
	}
	if total == 0 {
		return errort.NewCommonErr(errort.DeviceNotExist, fmt.Errorf("devices not found"))
	}
	var deviceIds []string
	for _, device := range devices {
		deviceIds = append(deviceIds, device.Id)
	}
	err = p.dbClient.BatchBindDevice(deviceIds, req.DriverInstanceId)
	if err != nil {
		return err
	}

	return nil
}

func (p *deviceApp) DevicesBindDriver(ctx context.Context, req dtos.DevicesBindDriver) error {
	var searchReq dtos.DeviceSearchQueryRequest
	searchReq.BaseSearchConditionQuery.Ids = dtos.ApiParamsArrayToString(req.DeviceIds)
	var devices []models.Device
	var err error
	var total uint32
	devices, total, err = p.dbClient.DevicesSearch(0, -1, searchReq)
	if err != nil {
		return err
	}
	if total == 0 {
		return errort.NewCommonErr(errort.DeviceNotExist, fmt.Errorf("devices not found"))
	}
	driverInstance, err := p.dbClient.DeviceServiceById(req.DriverInstanceId)
	if err != nil {
		return err
	}

	for _, device := range devices {
		if device.DriveInstanceId != "" {
			return errort.NewCommonErr(errort.DeviceNotUnbindDriver, fmt.Errorf("please unbind the device with the driver first"))
		}
	}

	for _, device := range devices {
		if driverInstance.Platform != device.Platform && driverInstance.Platform != "" {
			return errort.NewCommonErr(errort.DeviceAndDriverPlatformNotIdentical, fmt.Errorf("the device platform is inconsistent with the drive platform"))
		}
	}

	err = p.dbClient.BatchBindDevice(req.DeviceIds, req.DriverInstanceId)
	if err != nil {
		return err
	}

	for _, device := range devices {
		device.DriveInstanceId = req.DriverInstanceId
		callBackDevice := device
		go func() {
			p.CreateDeviceCallBack(callBackDevice)
		}()
	}
	return nil
}

func (p *deviceApp) DeviceUpdateConnectStatus(id string, status constants.DeviceStatus) error {
	device, edgeXErr := p.dbClient.DeviceById(id)
	if edgeXErr != nil {
		return edgeXErr
	}
	device.Status = status
	if status == constants.DeviceStatusOnline {
		device.LastOnlineTime = utils.MakeTimestamp()
	}
	edgeXErr = p.dbClient.UpdateDevice(device)
	if edgeXErr != nil {
		return edgeXErr
	}
	return nil
}

func setDeviceInfoSheet(file *dtos.ExportFile, req dtos.DeviceImportTemplateRequest) error {
	file.Excel.SetSheetName("Sheet1", dtos.DevicesFilename)

	file.Excel.SetCellStyle(dtos.DevicesFilename, "A1", "A1", file.GetCenterStyle())
	file.Excel.MergeCell(dtos.DevicesFilename, "A1", "B1")
	file.Excel.SetCellStr(dtos.DevicesFilename, "A1", "Device Base Info")

	file.Excel.SetCellStr(dtos.DevicesFilename, "A2", "DeviceName")
	file.Excel.SetCellStr(dtos.DevicesFilename, "B2", "Description")

	return nil
}

func (p *deviceApp) DeviceImportTemplateDownload(ctx context.Context, req dtos.DeviceImportTemplateRequest) (*dtos.ExportFile, error) {
	file, err := dtos.NewExportFile(dtos.DevicesFilename)
	if err != nil {
		return nil, err
	}
	if err := setDeviceInfoSheet(file, req); err != nil {
		p.lc.Error(err.Error())
		return nil, err
	}
	return file, nil
}

func (p *deviceApp) UploadValidated(ctx context.Context, file *dtos.ImportFile) error {
	rows, err := file.Excel.Rows(dtos.DevicesFilename)
	if err != nil {
		return errort.NewCommonEdgeX(errort.DefaultReadExcelErrorCode, "read rows error", err)
	}
	idx := 0
	for rows.Next() {
		idx++
		cols, err := rows.Columns()
		if err != nil {
			return errort.NewCommonEdgeX(errort.DefaultReadExcelErrorCode, "read cols error", err)
		}
		if idx == 1 {
			continue
		}
		if idx == 2 {
			if len(cols) != 2 {
				return errort.NewCommonEdgeX(errort.DefaultReadExcelErrorCode, fmt.Sprintf("read cols error need len %d,but read len %d", 2, len(cols)), err)
			}
			continue
		}

		// 空行过滤
		if len(cols) <= 0 {
			continue
		}

		if cols[0] == "" {
			return errort.NewCommonEdgeX(errort.DefaultReadExcelErrorParamsRequiredCode, fmt.Sprintf("read excel params required %+v", "deviceName"), nil)
		}
	}
	return nil
}

func (p *deviceApp) DevicesImport(ctx context.Context, file *dtos.ImportFile, productId, driverInstanceId string) (int64, error) {
	productService := resourceContainer.ProductAppNameFrom(p.dic.Get)
	productInfo, err := productService.ProductById(ctx, productId)
	if err != nil {
		return 0, err
	}

	if productInfo.Status == string(constants.ProductUnRelease) {
		return 0, errort.NewCommonEdgeX(errort.ProductUnRelease, "The product has not been released yet. Please release the product before adding devices", nil)
	}

	if driverInstanceId != "" {
		driverService := resourceContainer.DriverServiceAppFrom(p.dic.Get)
		driverInfo, err := driverService.Get(ctx, driverInstanceId)
		if err != nil {
			return 0, err
		}
		if driverInfo.Platform != "" && driverInfo.Platform != constants.IotPlatform_LocalIot {
			return 0, errort.NewCommonEdgeX(errort.DeviceServiceMustLocalPlatform, "driver service must local platform", err)
		}
	}

	rows, err := file.Excel.Rows(dtos.DevicesFilename)
	if err != nil {
		return 0, errort.NewCommonEdgeX(errort.DefaultReadExcelErrorCode, "read rows error", err)
	}
	devices := make([]models.Device, 0)
	idx := 0
	for rows.Next() {
		idx++
		deviceAddRequest := models.Device{
			ProductId:       productId,
			DriveInstanceId: driverInstanceId,
		}
		cols, err := rows.Columns()
		if err != nil {
			return 0, errort.NewCommonEdgeX(errort.DefaultReadExcelErrorCode, "read cols error", err)
		}
		if idx == 1 {
			continue
		}
		if idx == 2 {
			if len(cols) != 2 {
				return 0, errort.NewCommonEdgeX(errort.DefaultReadExcelErrorCode, fmt.Sprintf("read cols error need len %d,but read len %d", 2, len(cols)), err)
			}
			continue
		}

		// 空行过滤
		if len(cols) <= 0 {
			continue
		}

		deviceAddRequest.Id = utils.RandomNum()
		deviceAddRequest.Name = cols[0]
		if len(cols) >= 2 {
			deviceAddRequest.Description = cols[1]
		}
		deviceAddRequest.Status = constants.DeviceStatusOffline
		deviceAddRequest.Platform = constants.IotPlatform_LocalIot
		deviceAddRequest.Created = utils.MakeTimestamp()
		deviceAddRequest.Secret = utils.GenerateDeviceSecret(12)
		if deviceAddRequest.Name == "" {
			return 0, errort.NewCommonEdgeX(errort.DefaultReadExcelErrorParamsRequiredCode, fmt.Sprintf("read excel params required %+v", deviceAddRequest), nil)
		}
		devices = append(devices, deviceAddRequest)
	}

	for _, device := range devices {
		err = resourceContainer.DataDBClientFrom(p.dic.Get).CreateTable(ctx, productInfo.Id, device.Id)
		if err != nil {
			return 0, err
		}
	}

	total, err := p.dbClient.BatchUpsertDevice(devices)
	if err != nil {
		return 0, err
	}

	for _, device := range devices {
		addDevice := device
		go func() {
			p.CreateDeviceCallBack(addDevice)
		}()
	}
	return total, nil
}

func (p *deviceApp) DevicesReportMsgGather(ctx context.Context) error {
	var count int
	var err error
	startTime, endTime := GetYesterdayStartTimeAndEndTime()
	persistApp := resourceContainer.PersistItfFrom(p.dic.Get)
	count, err = persistApp.SearchDeviceMsgCount(startTime, endTime)
	if err != nil {

	}
	var msgGather models.MsgGather
	msgGather.Count = count
	msgGather.Date = time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	return p.dbClient.AddMsgGather(msgGather)
}

func GetYesterdayStartTimeAndEndTime() (int64, int64) {
	NowTime := time.Now()
	var startTime time.Time
	if NowTime.Hour() == 0 && NowTime.Minute() == 0 && NowTime.Second() == 0 {
		startTime = time.Unix(NowTime.Unix()-86399, 0) //当天的最后一秒
	} else {
		startTime = time.Unix(NowTime.Unix()-86400, 0)
	}
	currentYear := startTime.Year()
	currentMonth := startTime.Month()
	currentDay := startTime.Day()
	yesterdayStartTime := time.Date(currentYear, currentMonth, currentDay, 0, 0, 0, 0, time.Local).UnixMilli()
	yesterdayEndTime := time.Date(currentYear, currentMonth, currentDay, 23, 59, 59, 0, time.Local).UnixMilli()
	return yesterdayStartTime, yesterdayEndTime
}
