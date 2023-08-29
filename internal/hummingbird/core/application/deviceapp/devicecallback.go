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
	"github.com/winc-link/edge-driver-proto/devicecallback"
	"github.com/winc-link/hummingbird/internal/hummingbird/core/container"
	"github.com/winc-link/hummingbird/internal/models"
	"github.com/winc-link/hummingbird/internal/pkg/constants"
	"github.com/winc-link/hummingbird/internal/pkg/di"
	"github.com/winc-link/hummingbird/internal/tools/rpcclient"
	"time"
)

func (p *deviceApp) CreateDeviceCallBack(createDevice models.Device) {
	defer func() {
		if err := recover(); err != nil {
			p.lc.Error("CreateDeviceCallBack Panic:", err)
		}
	}()
	deviceService, err := p.dbClient.DeviceServiceById(createDevice.DriveInstanceId)
	if err != nil {
		return
	}

	driverService := container.DriverServiceAppFrom(di.GContainer.Get)
	status := driverService.GetState(deviceService.Id)
	if status == constants.RunStatusStarted {
		client, errX := rpcclient.NewDriverRpcClient(deviceService.BaseAddress, false, "", deviceService.Id, p.lc)
		if errX != nil {
			return
		}
		defer client.Close()
		var rpcRequest devicecallback.CreateDeviceCallbackRequest
		rpcRequest.Data = createDevice.TransformToDriverDevice()
		rpcRequest.HappenTime = uint64(time.Now().Unix())
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		_, _ = client.DeviceCallBackServiceClient.CreateDeviceCallback(ctx, &rpcRequest)
	}
}

func (p *deviceApp) UpdateDeviceCallBack(updateDevice models.Device) {
	defer func() {
		if err := recover(); err != nil {
			p.lc.Error("UpdateDeviceCallBack Panic:", err)
		}
	}()
	deviceService, err := p.dbClient.DeviceServiceById(updateDevice.DriveInstanceId)
	if err != nil {
		return
	}

	driverService := container.DriverServiceAppFrom(di.GContainer.Get)
	status := driverService.GetState(deviceService.Id)

	if status == constants.RunStatusStarted {
		client, errX := rpcclient.NewDriverRpcClient(deviceService.BaseAddress, false, "", deviceService.Id, p.lc)
		if errX != nil {
			return
		}
		defer client.Close()
		var rpcRequest devicecallback.UpdateDeviceCallbackRequest
		rpcRequest.Data = updateDevice.TransformToDriverDevice()
		rpcRequest.HappenTime = uint64(time.Now().Second())
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		_, _ = client.DeviceCallBackServiceClient.UpdateDeviceCallback(ctx, &rpcRequest)
	}
}

func (p *deviceApp) DeleteDeviceCallBack(deleteDevice models.Device) {
	//查出哪些驱动和这个平台相关联，做推送通知。
	defer func() {
		if err := recover(); err != nil {
			p.lc.Error("DeleteDeviceCallBack Panic:", err)
		}
	}()
	deviceService, err := p.dbClient.DeviceServiceById(deleteDevice.DriveInstanceId)
	if err != nil {
		return
	}

	driverService := container.DriverServiceAppFrom(di.GContainer.Get)
	status := driverService.GetState(deviceService.Id)

	if status == constants.RunStatusStarted {
		client, errX := rpcclient.NewDriverRpcClient(deviceService.BaseAddress, false, "", deviceService.Id, p.lc)
		if errX != nil {
			return
		}
		defer client.Close()
		var rpcRequest devicecallback.DeleteDeviceCallbackRequest
		rpcRequest.DeviceId = deleteDevice.Id
		rpcRequest.HappenTime = uint64(time.Now().Second())
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		_, _ = client.DeviceCallBackServiceClient.DeleteDeviceCallback(ctx, &rpcRequest)
	}
}
