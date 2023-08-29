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
	"github.com/winc-link/edge-driver-proto/cloudinstance"
	"github.com/winc-link/edge-driver-proto/drivercommon"
	"github.com/winc-link/hummingbird/internal/dtos"
	"github.com/winc-link/hummingbird/internal/hummingbird/core/container"
	"github.com/winc-link/hummingbird/internal/pkg/constants"
	"github.com/winc-link/hummingbird/internal/pkg/di"
	"github.com/winc-link/hummingbird/internal/pkg/logger"
	"google.golang.org/grpc"
)

type CloudInstanceServer struct {
	cloudinstance.UnimplementedCloudInstanceServiceServer

	lc  logger.LoggingClient
	dic *di.Container
}

func (s *CloudInstanceServer) DriverReportPlatformInfo(ctx context.Context, request *cloudinstance.DriverReportPlatformInfoRequest) (*cloudinstance.DriverReportPlatformInfoResponse, error) {
	response := new(cloudinstance.DriverReportPlatformInfoResponse)
	response.BaseResponse = new(drivercommon.CommonResponse)
	if request.GetDriverInstanceId() == "" {
		response.BaseResponse.Success = false
		response.BaseResponse.ErrorMessage = "param error"
		return response, nil
	}

	driverService := container.DriverServiceAppFrom(s.dic.Get)
	req := dtos.DeviceServiceUpdateRequest{}
	req.Id = request.GetDriverInstanceId()
	req.Platform = constants.TransformEdgePlatformToDbPlatform(request.GetIotPlatform())
	err := driverService.Update(ctx, req)
	if err != nil {
		response.BaseResponse.ErrorMessage = err.Error()
	}
	response.BaseResponse.Success = true
	return response, nil
}

func (s *CloudInstanceServer) QueryCloudInstanceByPlatform(ctx context.Context, request *cloudinstance.QueryCloudInstanceByPlatformRequest) (*cloudinstance.QueryCloudInstanceByPlatformResponse, error) {
	response := new(cloudinstance.QueryCloudInstanceByPlatformResponse)
	return response, nil

}

func NewCloudInstanceServer(lc logger.LoggingClient, dic *di.Container) *CloudInstanceServer {
	return &CloudInstanceServer{
		lc:  lc,
		dic: dic,
	}
}

func (s *CloudInstanceServer) RegisterServer(server *grpc.Server) {
	cloudinstance.RegisterCloudInstanceServiceServer(server, s)
}
