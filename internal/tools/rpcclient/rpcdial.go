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
	"github.com/winc-link/hummingbird/internal/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/protobuf/types/known/emptypb"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var Empty = &emptypb.Empty{}

// grpc内部对client的ping server的周期限制为最小10s
var kacp = keepalive.ClientParameters{
	Time:                10 * time.Second, // send pings every 10 seconds if there is no activity
	Timeout:             3 * time.Second,  // wait 3 second for ping ack before considering the connection dead
	PermitWithoutStream: true,             // send pings even without active streams
}

var ConnParams = grpc.ConnectParams{
	Backoff: backoff.Config{
		BaseDelay:  time.Second * 1.0,
		Multiplier: 1.0,
		Jitter:     0,
		MaxDelay:   1.0 * time.Second,
	},
	MinConnectTimeout: time.Second * 1,
}

var basepath string

func init() {
	_, currentFile, _, _ := runtime.Caller(0)
	basepath = filepath.Dir(currentFile)
}

// Path returns the absolute path the given relative file or directory path,
// relative to the google.golang.org/grpc/examples/data directory in the
// user's GOPATH.  If rel is already absolute, it is returned unmodified.
func path(rel string) string {
	if filepath.IsAbs(rel) {
		return rel
	}

	return filepath.Join(basepath, rel)
}

const (
	CLIENT_RPC_LOG_ENABLE = "CLIENT_RPC_LOG_ENABLE"
)

func dialWithLog(address string, withTLS bool, certFile, serverName string, lc logger.LoggingClient) (*grpc.ClientConn, error) {
	var (
		err   error
		creds credentials.TransportCredentials
		conn  *grpc.ClientConn
	)
	var clientRpcLog = os.Getenv(CLIENT_RPC_LOG_ENABLE)
	if withTLS {
		if creds, err = credentials.NewClientTLSFromFile(path(certFile), serverName); err != nil {
			return nil, err
		}
		if clientRpcLog != "" {
			conn, err = grpc.Dial(address,
				grpc.WithTransportCredentials(creds),
				grpc.WithBlock(),
				grpc.WithKeepaliveParams(kacp),
				withClientLoggerInterceptor(lc),
				grpc.WithConnectParams(ConnParams),
				withEdgeErrorInterceptor(),
			)
		} else {
			conn, err = grpc.Dial(address,
				grpc.WithTransportCredentials(creds),
				grpc.WithBlock(),
				grpc.WithKeepaliveParams(kacp),
				grpc.WithConnectParams(ConnParams),
				withEdgeErrorInterceptor(),
			)
		}
		if err != nil {
			return nil, err
		}
	} else {
		if clientRpcLog != "" {
			conn, err = grpc.Dial(address,
				grpc.WithInsecure(), /*grpc.WithBlock(),*/
				grpc.WithKeepaliveParams(kacp),
				withClientLoggerInterceptor(lc),
				grpc.WithConnectParams(ConnParams),
				withEdgeErrorInterceptor(),
			)
		} else {
			conn, err = grpc.Dial(address,
				grpc.WithInsecure(), /*grpc.WithBlock(),*/
				grpc.WithKeepaliveParams(kacp),
				grpc.WithConnectParams(ConnParams),
			)
		}
		if err != nil {
			return nil, err
		}
	}
	return conn, nil
}
