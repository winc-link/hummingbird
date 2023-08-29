/*******************************************************************************
 * Copyright 2017.
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
	driverstorage "github.com/winc-link/edge-driver-proto/driverstorge"
	"github.com/winc-link/hummingbird/internal/hummingbird/core/container"
	"github.com/winc-link/hummingbird/internal/pkg/di"
	"github.com/winc-link/hummingbird/internal/pkg/logger"
	"github.com/winc-link/hummingbird/internal/tools/datadb/leveldb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"os"
	"strings"
	"sync"
)

type DriverStorageServer struct {
	driverstorage.UnimplementedDriverStorageServer
	dMap    map[string]*leveldb.DriverStorageClient
	mu      sync.Mutex
	dirPath string
	lc      logger.LoggingClient
	dic     *di.Container
}

func (s *DriverStorageServer) All(ctx context.Context, req *driverstorage.AllReq) (*driverstorage.KVs, error) {
	id := req.GetDriverServiceId()
	if len(id) <= 0 {
		return nil, status.Error(codes.InvalidArgument, "driver service not set")
	}
	s.lc.Infof("get driver storage, driver service: %s", id)
	client, err := s.getStorageClient(id)
	if err != nil {
		s.lc.Errorf("get leveldb client error: %s,driver service: %s", err, id)
		return nil, status.Error(codes.Internal, err.Error())
	}
	all, err := client.All()
	if err != nil {
		s.lc.Errorf("get all kvs error: %s,driver service: %s", err, id)
		return nil, status.Error(codes.Internal, err.Error())
	}

	var kvs driverstorage.KVs
	for k, v := range all {
		kvs.Kvs = append(kvs.Kvs, &driverstorage.KV{
			Key:   k,
			Value: v,
		})
	}
	return &kvs, nil
}

func (s *DriverStorageServer) Get(ctx context.Context, req *driverstorage.GetReq) (*driverstorage.KVs, error) {
	id := req.GetDriverServiceId()
	if len(id) <= 0 {
		return nil, status.Error(codes.InvalidArgument, "driver service not set")
	}
	client, err := s.getStorageClient(id)
	if err != nil {
		s.lc.Errorf("get leveldb client error: %s,driver service: %s", err, id)
		return nil, status.Error(codes.Internal, err.Error())
	}
	keys := req.GetKeys()
	if len(keys) <= 0 {
		return nil, status.Error(codes.InvalidArgument, "keys length is 0")
	}

	s.lc.Infof("get driver storage, driver service: %s, keys: %+v", id, keys)

	kvs, _ := client.Get(keys)
	// convert
	var resp driverstorage.KVs
	for k, v := range kvs {
		resp.Kvs = append(resp.Kvs, &driverstorage.KV{
			Key:   k,
			Value: v,
		})
	}
	return &resp, nil
}

func (s *DriverStorageServer) Put(ctx context.Context, req *driverstorage.PutReq) (*emptypb.Empty, error) {
	id := req.GetDriverServiceId()
	if len(id) <= 0 {
		return &emptypb.Empty{}, status.Error(codes.InvalidArgument, "driver service not set")
	}
	client, err := s.getStorageClient(id)
	if err != nil {
		s.lc.Errorf("get leveldb client error: %s,driver service: %s", err, id)
		return &emptypb.Empty{}, status.Error(codes.Internal, err.Error())
	}
	data := req.GetData()
	keys := make([]string, 0, len(data))
	kvs := make(map[string][]byte, len(data))
	for _, v := range data {
		keys = append(keys, v.GetKey())
		kvs[v.GetKey()] = v.GetValue()
	}

	s.lc.Infof("put driver storage, driver service: %s, keys: %+v", id, keys)

	if err := client.Put(kvs); err != nil {
		return &emptypb.Empty{}, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (s *DriverStorageServer) Delete(ctx context.Context, req *driverstorage.DeleteReq) (*emptypb.Empty, error) {
	id := req.GetDriverServiceId()
	if len(id) <= 0 {
		return &emptypb.Empty{}, status.Error(codes.InvalidArgument, "driver service not set")
	}
	client, err := s.getStorageClient(id)
	if err != nil {
		s.lc.Errorf("get leveldb client error: %s,driver service: %s", err, id)
		return &emptypb.Empty{}, status.Error(codes.Internal, err.Error())
	}

	keys := req.GetKeys()
	if len(keys) <= 0 {
		return nil, status.Error(codes.InvalidArgument, "keys length is 0")
	}

	s.lc.Infof("delete driver storage, driver service: %s, keys: %+v", id, keys)

	if err := client.Delete(keys); err != nil {
		s.lc.Errorf("driver storage delete keys(%+v) error: %s", keys, err)
		return &emptypb.Empty{}, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func NewDriverStorageServer(lc logger.LoggingClient, dic *di.Container) *DriverStorageServer {
	var dir string
	config := container.ConfigurationFrom(dic.Get)
	if ldb, ok := config.Databases["Data"]; !ok {
		lc.Errorf("leveldb not config")
		os.Exit(-1)
	} else {
		dir = ldb["Primary"].DataSource
		if len(dir) <= 0 {
			lc.Errorf("leveldb not config")
			os.Exit(-1)
		}
	}

	strN := strings.SplitN(dir, "/", 2)
	if len(strN) < 1 {
		lc.Errorf("leveldb config error")
		os.Exit(-1)
	}
	dir = strN[0] + "/"
	return &DriverStorageServer{
		lc:      lc,
		dMap:    make(map[string]*leveldb.DriverStorageClient),
		dirPath: dir,
		dic:     dic,
	}
}

func (s *DriverStorageServer) RegisterServer(server *grpc.Server) {
	driverstorage.RegisterDriverStorageServer(server, s)
}

func (s *DriverStorageServer) getStorageClient(id string) (*leveldb.DriverStorageClient, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var (
		err    error
		ok     bool
		client *leveldb.DriverStorageClient
	)

	if client, ok = s.dMap[id]; !ok {
		if client, err = leveldb.NewDriverStorageClient(s.dirPath, id, s.lc); err != nil {
			return nil, err
		}
		s.dMap[id] = client
	}
	return client, nil
}
