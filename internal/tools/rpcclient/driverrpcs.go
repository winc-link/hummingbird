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
package rpcclient

import (
	"errors"
	"github.com/winc-link/edge-driver-proto/productcallback"
	"github.com/winc-link/edge-driver-proto/thingmodel"

	//"github.com/winc-link/edge-driver-proto/productcallback"

	"github.com/winc-link/edge-driver-proto/cloudinstancecallback"
	"github.com/winc-link/edge-driver-proto/devicecallback"

	"github.com/winc-link/hummingbird/internal/pkg/logger"
	"google.golang.org/grpc"
)

type DriverRpcClient struct {
	address string
	Conn    *grpc.ClientConn
	devicecallback.DeviceCallBackServiceClient
	cloudinstancecallback.CloudInstanceCallBackServiceClient
	productcallback.ProductCallBackServiceClient
	thingmodel.ThingModelDownServiceClient
}

func NewDriverRpcClient(address string, useTLS bool, certFile, serverName string, lc logger.LoggingClient) (*DriverRpcClient, error) {
	var (
		err  error
		conn *grpc.ClientConn
	)
	if address == "" {
		return nil, errors.New("required address")
	}
	if conn, err = dialWithLog(address, useTLS, certFile, serverName, lc); err != nil {
		return &DriverRpcClient{}, err
	}
	return &DriverRpcClient{
		address:                            address,
		Conn:                               conn,
		CloudInstanceCallBackServiceClient: cloudinstancecallback.NewCloudInstanceCallBackServiceClient(conn),
		DeviceCallBackServiceClient:        devicecallback.NewDeviceCallBackServiceClient(conn),
		ProductCallBackServiceClient:       productcallback.NewProductCallBackServiceClient(conn),
		ThingModelDownServiceClient:        thingmodel.NewThingModelDownServiceClient(conn),
	}, nil
}

func (d *DriverRpcClient) Close() error {
	return d.Conn.Close()
}
