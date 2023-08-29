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
	"github.com/winc-link/edge-driver-proto/drivercommon"
	"github.com/winc-link/edge-driver-proto/thingmodel"
	"github.com/winc-link/hummingbird/internal/dtos"
	coreContainer "github.com/winc-link/hummingbird/internal/hummingbird/core/container"
	"github.com/winc-link/hummingbird/internal/pkg/di"
	"github.com/winc-link/hummingbird/internal/pkg/logger"
	"google.golang.org/grpc"
)

type ThingModelServer struct {
	thingmodel.UnimplementedThingModelUpServiceServer

	lc  logger.LoggingClient
	dic *di.Container
}

func (s *ThingModelServer) ThingModelMsgReport(ctx context.Context, msg *thingmodel.ThingModelMsg) (response *drivercommon.CommonResponse, err error) {
	message := dtos.ThingModelMessageFromThingModelMsg(msg)
	messageItf := coreContainer.MessageItfFrom(s.dic.Get)
	return messageItf.ThingModelMsgReport(ctx, message)
}

var _ thingmodel.ThingModelUpServiceServer = (*ThingModelServer)(nil)

func NewThingModelServer(lc logger.LoggingClient, dic *di.Container) *ThingModelServer {
	return &ThingModelServer{
		lc:  lc,
		dic: dic,
	}
}

func (s *ThingModelServer) RegisterServer(server *grpc.Server) {
	thingmodel.RegisterThingModelUpServiceServer(server, s)
}
