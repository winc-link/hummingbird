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
	"github.com/winc-link/edge-driver-proto/gateway"
	"github.com/winc-link/hummingbird/internal/pkg/di"
	"github.com/winc-link/hummingbird/internal/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GatewayServer struct {
	gateway.UnimplementedRpcGatewayServer

	lc  logger.LoggingClient
	dic *di.Container
}

func (s *GatewayServer) GetGatewayInfo(ctx context.Context, empty *emptypb.Empty) (*gateway.GateWayInfoResponse, error) {
	response := new(gateway.GateWayInfoResponse)
	response.Env = "env"
	response.GwId = "gatewayId"
	response.LocalKey = "localKey"

	return response, nil
}

func NewGatewayServer(lc logger.LoggingClient, dic *di.Container) *GatewayServer {
	return &GatewayServer{
		lc:  lc,
		dic: dic,
	}
}

func (s *GatewayServer) RegisterServer(server *grpc.Server) {
	gateway.RegisterRpcGatewayServer(server, s)
}
