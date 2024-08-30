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
package driverserver

import (
	"context"
	"fmt"
	"github.com/winc-link/edge-driver-proto/drivercommon"
	device "github.com/winc-link/edge-driver-proto/driverdevice"
	"github.com/winc-link/hummingbird/internal/dtos"
	"github.com/winc-link/hummingbird/internal/hummingbird/core/container"
	"github.com/winc-link/hummingbird/internal/pkg/constants"
	"github.com/winc-link/hummingbird/internal/pkg/di"
	"github.com/winc-link/hummingbird/internal/pkg/errort"
	"github.com/winc-link/hummingbird/internal/pkg/logger"
	"google.golang.org/grpc"
	"strconv"
)

type DriverDeviceServer struct {
	device.UnimplementedRpcDeviceServer
	lc  logger.LoggingClient
	dic *di.Container
}

func (s *DriverDeviceServer) ConnectIotPlatform(ctx context.Context, request *device.ConnectIotPlatformRequest) (*device.ConnectIotPlatformResponse, error) {
	deviceItf := container.DeviceItfFrom(s.dic.Get)
	return deviceItf.ConnectIotPlatform(ctx, request), nil
}

func (s *DriverDeviceServer) DisconnectIotPlatform(ctx context.Context, request *device.DisconnectIotPlatformRequest) (*device.DisconnectIotPlatformResponse, error) {
	deviceItf := container.DeviceItfFrom(s.dic.Get)
	return deviceItf.DisConnectIotPlatform(ctx, request), nil
}

func (s *DriverDeviceServer) GetDeviceConnectStatus(ctx context.Context, request *device.GetDeviceConnectStatusRequest) (*device.GetDeviceConnectStatusResponse, error) {
	deviceItf := container.DeviceItfFrom(s.dic.Get)
	return deviceItf.GetDeviceConnectStatus(ctx, request), nil
}

func (s *DriverDeviceServer) QueryDeviceList(ctx context.Context, request *device.QueryDeviceListRequest) (*device.QueryDeviceListResponse, error) {
	deviceItf := container.DeviceItfFrom(s.dic.Get)

	var platform string
	if request.BaseRequest.UseCloudPlatform {
		platform = string(constants.TransformEdgePlatformToDbPlatform(request.BaseRequest.GetCloudInstanceInfo().GetIotPlatform()))
	} else {
		platform = string(constants.IotPlatform_LocalIot)
	}
	devices, total, err := deviceItf.DevicesModelSearch(ctx, dtos.DeviceSearchQueryRequest{
		DriveInstanceId: request.BaseRequest.DriverInstanceId,
		Platform:        platform,
	})
	response := new(device.QueryDeviceListResponse)
	response.BaseResponse = new(drivercommon.CommonResponse)
	if err != nil {
		response.BaseResponse.Success = false
		response.BaseResponse.ErrorMessage = err.Error()
		return response, nil
	}
	response.BaseResponse.Success = true
	response.Data = new(device.QueryDeviceListResponse_Data)
	response.Data.Total = total
	for _, queryResponse := range devices {
		response.Data.Devices = append(response.Data.Devices, queryResponse.TransformToDriverDevice())
	}
	return response, nil
}

func (s *DriverDeviceServer) QueryDeviceById(ctx context.Context, request *device.QueryDeviceByIdRequest) (*device.QueryDeviceByIdResponse, error) {
	deviceItf := container.DeviceItfFrom(s.dic.Get)

	deviceInfo, err := deviceItf.DeviceModelById(ctx, request.Id)
	response := new(device.QueryDeviceByIdResponse)
	response.BaseResponse = new(drivercommon.CommonResponse)
	if err != nil {
		response.BaseResponse.Success = false
		response.BaseResponse.ErrorMessage = err.Error()
		return response, nil
	}
	response.BaseResponse.Success = true
	response.Data = new(device.QueryDeviceByIdResponse_Data)
	response.Data.Device = deviceInfo.TransformToDriverDevice()
	return response, nil
}

func (s *DriverDeviceServer) CreateDevice(ctx context.Context, request *device.CreateDeviceRequest) (*device.CreateDeviceRequestResponse, error) {
	response := new(device.CreateDeviceRequestResponse)
	response.BaseResponse = new(drivercommon.CommonResponse)
	deviceItf := container.DeviceItfFrom(s.dic.Get)
	productItf := container.ProductAppNameFrom(s.dic.Get)
	productInfo, err := productItf.ProductById(ctx, request.Device.ProductId)
	if err != nil {
		err = errort.NewCommonErr(errort.ProductNotExist, fmt.Errorf(""))
		errWrapper := errort.NewCommonEdgeXWrapper(err)
		response.BaseResponse.Success = false
		response.BaseResponse.Code = strconv.Itoa(int(errWrapper.Code()))
		response.BaseResponse.ErrorMessage = errWrapper.Message()
		return response, nil
	}

	var insertDevice dtos.DeviceAddRequest
	insertDevice.ProductId = productInfo.Id
	insertDevice.Platform = constants.IotPlatform_LocalIot
	insertDevice.Name = request.Device.Name
	insertDevice.DeviceSn = request.Device.DeviceSn
	//insertDevice.d
	insertDevice.DriverInstanceId = request.BaseRequest.GetDriverInstanceId()

	deviceId, err := deviceItf.AddDevice(ctx, insertDevice)
	if err != nil {
		errWrapper := errort.NewCommonEdgeXWrapper(err)
		response.BaseResponse.Success = false
		response.BaseResponse.Code = strconv.Itoa(int(errWrapper.Code()))
		response.BaseResponse.ErrorMessage = errWrapper.Message()
		return response, nil
	}
	deviceInfo, _ := deviceItf.DeviceById(ctx, deviceId)
	response.BaseResponse.Success = true
	response.Data = new(device.CreateDeviceRequestResponse_Data)
	response.Data.Devices = new(device.Device)
	response.Data.Devices.Id = deviceId
	response.Data.Devices.Name = request.Device.Name
	response.Data.Devices.ProductId = request.Device.ProductId
	response.Data.Devices.DeviceSn = request.Device.DeviceSn
	response.Data.Devices.External = request.Device.External
	response.Data.Devices.Secret = deviceInfo.Secret
	response.Data.Devices.Description = request.Device.Description
	response.Data.Devices.Status = device.DeviceStatus_OffLine
	response.Data.Devices.Platform = drivercommon.IotPlatform_LocalIot
	return response, nil
}

func (s *DriverDeviceServer) CreateDeviceAndConnect(ctx context.Context, request *device.CreateDeviceAndConnectRequest) (*device.CreateDeviceAndConnectRequestResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *DriverDeviceServer) DeleteDevice(ctx context.Context, request *device.DeleteDeviceRequest) (*device.DeleteDeviceResponse, error) {
	//TODO implement me
	panic("implement me")
}

var _ device.RpcDeviceServer = (*DriverDeviceServer)(nil)

func NewDriverDeviceServer(lc logger.LoggingClient, dic *di.Container) *DriverDeviceServer {
	return &DriverDeviceServer{
		lc:  lc,
		dic: dic,
	}
}

func (s *DriverDeviceServer) RegisterServer(server *grpc.Server) {
	device.RegisterRpcDeviceServer(server, s)
}
