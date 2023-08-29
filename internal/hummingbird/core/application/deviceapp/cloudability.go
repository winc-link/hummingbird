/*******************************************************************************
 * Copyright 2017 Dell Inc.
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
	"github.com/docker/distribution/uuid"
	"github.com/winc-link/hummingbird/internal/hummingbird/core/container"
	"github.com/winc-link/hummingbird/internal/pkg/constants"
	"github.com/winc-link/hummingbird/internal/pkg/errort"
	"strconv"

	"github.com/winc-link/edge-driver-proto/drivercommon"
	"github.com/winc-link/edge-driver-proto/driverdevice"
)

func (p deviceApp) ConnectIotPlatform(ctx context.Context, request *driverdevice.ConnectIotPlatformRequest) *driverdevice.ConnectIotPlatformResponse {
	response := new(driverdevice.ConnectIotPlatformResponse)
	baseResponse := new(drivercommon.CommonResponse)
	baseResponse.Success = false
	deviceInfo, err := p.DeviceById(ctx, request.DeviceId)

	if err != nil {
		errWrapper := errort.NewCommonEdgeXWrapper(err)
		baseResponse.Code = strconv.Itoa(int(errWrapper.Code()))
		baseResponse.ErrorMessage = errWrapper.Message()
		response.BaseResponse = baseResponse
		return response
	}
	//把消息投体进入消息总线
	messageApp := container.MessageItfFrom(p.dic.Get)
	messageApp.DeviceStatusToMessageBus(ctx, deviceInfo.Id, constants.DeviceOnline)

	err = p.dbClient.DeviceOnlineById(request.DeviceId)
	if err != nil {
		errWrapper := errort.NewCommonEdgeXWrapper(err)
		baseResponse.Code = strconv.Itoa(int(errWrapper.Code()))
		baseResponse.ErrorMessage = errWrapper.Message()
		response.BaseResponse = baseResponse
		return response
	} else {
		baseResponse.Success = true
		baseResponse.RequestId = uuid.Generate().String()
		response.Data = new(driverdevice.ConnectIotPlatformResponse_Data)
		response.Data.Status = driverdevice.ConnectStatus_ONLINE
		response.BaseResponse = baseResponse
		return response
	}
}

func (p deviceApp) DisConnectIotPlatform(ctx context.Context, request *driverdevice.DisconnectIotPlatformRequest) *driverdevice.DisconnectIotPlatformResponse {
	deviceInfo, err := p.DeviceById(ctx, request.DeviceId)
	response := new(driverdevice.DisconnectIotPlatformResponse)
	baseResponse := new(drivercommon.CommonResponse)
	baseResponse.Success = false
	if err != nil {
		errWrapper := errort.NewCommonEdgeXWrapper(err)
		baseResponse.Code = strconv.Itoa(int(errWrapper.Code()))
		baseResponse.ErrorMessage = errWrapper.Message()
		response.BaseResponse = baseResponse
		return response
	}
	messageApp := container.MessageItfFrom(p.dic.Get)
	messageApp.DeviceStatusToMessageBus(ctx, deviceInfo.Id, constants.DeviceOffline)

	err = p.dbClient.DeviceOfflineById(request.DeviceId)
	if err != nil {
		errWrapper := errort.NewCommonEdgeXWrapper(err)
		baseResponse.Code = strconv.Itoa(int(errWrapper.Code()))
		baseResponse.ErrorMessage = errWrapper.Message()
		response.BaseResponse = baseResponse
		return response
	}
	baseResponse.Success = true
	baseResponse.RequestId = uuid.Generate().String()
	response.Data = new(driverdevice.DisconnectIotPlatformResponse_Data)
	response.Data.Status = driverdevice.ConnectStatus_OFFLINE
	response.BaseResponse = baseResponse
	return response
}

func (p deviceApp) GetDeviceConnectStatus(ctx context.Context, request *driverdevice.GetDeviceConnectStatusRequest) *driverdevice.GetDeviceConnectStatusResponse {
	deviceInfo, err := p.DeviceById(ctx, request.DeviceId)
	response := new(driverdevice.GetDeviceConnectStatusResponse)
	baseResponse := new(drivercommon.CommonResponse)
	baseResponse.Success = false
	if err != nil {
		errWrapper := errort.NewCommonEdgeXWrapper(err)
		baseResponse.Code = strconv.Itoa(int(errWrapper.Code()))
		baseResponse.ErrorMessage = errWrapper.Message()
		response.BaseResponse = baseResponse
		return response
	}

	baseResponse.Success = true
	baseResponse.RequestId = uuid.Generate().String()
	response.Data = new(driverdevice.GetDeviceConnectStatusResponse_Data)
	if deviceInfo.Status == constants.DeviceStatusOnline {
		response.Data.Status = driverdevice.ConnectStatus_ONLINE
	} else {
		response.Data.Status = driverdevice.ConnectStatus_OFFLINE
	}
	response.BaseResponse = baseResponse
	return response

}
