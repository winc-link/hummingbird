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

package handlers

import (
	"context"
	"github.com/winc-link/hummingbird/internal/pkg/di"
	"github.com/winc-link/hummingbird/internal/pkg/errort"
	"github.com/winc-link/hummingbird/internal/pkg/limit"
	"github.com/winc-link/hummingbird/internal/pkg/logger"
	"google.golang.org/grpc"
)

func withServerInterceptor(lc logger.LoggingClient, lmc limit.LimitMethodConf, dic *di.Container) grpc.ServerOption {
	return grpc.UnaryInterceptor(func(ctx context.Context, req interface{},
		info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		//if methodLimit(info.FullMethod, lmc) {
		//if err := requestLimit(ctx, info.FullMethod); err != nil {
		//	lc.Error(err.Error())
		//	return nil, errort.NewRPCStatusErr(err)
		//}
		//}
		//start := time.Now()
		reply, err := handler(ctx, req)
		//var hl = dtos.HandleLog{
		//	Req:      req,
		//	Reply:    reply,
		//	Method:   info.FullMethod,
		//	Duration: time.Since(start).String(),
		//}
		//if err != nil {
		//	hl.Error = err.Error()
		//} else {
		//	hl.Success = true
		//}
		//if serverRpcLog != "" {
		//	lc.Debug(hl.ToString())
		//}
		return reply, errort.NewRPCStatusErr(err)
	})
}
