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
	"github.com/winc-link/hummingbird/internal/pkg/di"
	"github.com/winc-link/hummingbird/internal/pkg/logger"
	"google.golang.org/grpc"
)

func RegisterRPCService(lc logger.LoggingClient, dic *di.Container, s *grpc.Server) {
	NewThingModelServer(lc, dic).RegisterServer(s)
	NewDriverDeviceServer(lc, dic).RegisterServer(s)
	NewCloudInstanceServer(lc, dic).RegisterServer(s)
	NewGatewayServer(lc, dic).RegisterServer(s)
	//NewDriverStorageServer(lc, dic).RegisterServer(s)
	NewProductServer(lc, dic).RegisterServer(s)
}
