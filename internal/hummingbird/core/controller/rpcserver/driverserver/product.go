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
	"github.com/winc-link/edge-driver-proto/driverproduct"
	"github.com/winc-link/hummingbird/internal/dtos"
	"github.com/winc-link/hummingbird/internal/hummingbird/core/container"
	"github.com/winc-link/hummingbird/internal/pkg/constants"
	"github.com/winc-link/hummingbird/internal/pkg/di"
	"github.com/winc-link/hummingbird/internal/pkg/logger"
	"google.golang.org/grpc"
)

type ProductServer struct {
	driverproduct.UnimplementedRpcProductServer
	lc  logger.LoggingClient
	dic *di.Container
}

func (s *ProductServer) QueryProductList(ctx context.Context, request *driverproduct.QueryProductListRequest) (*driverproduct.QueryProductListResponse, error) {
	var productPlatform constants.IotPlatform
	if request.BaseRequest != nil && request.BaseRequest.UseCloudPlatform {
		s.lc.Infof("request.BaseRequest.GetCloudInstanceInfo().IotPlatform:", request.BaseRequest.GetCloudInstanceInfo().IotPlatform)
		switch request.BaseRequest.GetCloudInstanceInfo().IotPlatform {
		case drivercommon.IotPlatform_WinCLinkIot:
			productPlatform = constants.IotPlatform_WinCLinkIot
		case drivercommon.IotPlatform_AliIot:
			productPlatform = constants.IotPlatform_AliIot
		case drivercommon.IotPlatform_HuaweiIot:
			productPlatform = constants.IotPlatform_HuaweiIot
		case drivercommon.IotPlatform_TencentIot:
			productPlatform = constants.IotPlatform_TencentIot
		case drivercommon.IotPlatform_TuyaIot:
			productPlatform = constants.IotPlatform_TuyaIot
		case drivercommon.IotPlatform_OneNetIot:
			productPlatform = constants.IotPlatform_OneNetIot
		default:
			productPlatform = constants.IotPlatform_LocalIot
		}
	} else {
		productPlatform = constants.IotPlatform_LocalIot
	}

	productItf := container.ProductAppNameFrom(s.dic.Get)

	response := new(driverproduct.QueryProductListResponse)
	response.Data = new(driverproduct.QueryProductListResponse_Data)

	response.BaseResponse = new(drivercommon.CommonResponse)

	productsModel, totol, err := productItf.ProductsModelSearch(ctx, dtos.ProductSearchQueryRequest{
		Platform: string(productPlatform),
	})
	if err != nil {
		response.BaseResponse.Success = false
		response.BaseResponse.ErrorMessage = err.Error()
		return response, nil
	}
	response.Data.Total = totol
	var driverProducts []*driverproduct.Product
	for _, productModel := range productsModel {
		driverProducts = append(driverProducts, productModel.TransformToDriverProduct())

	}
	response.BaseResponse.Success = true
	response.Data.Products = driverProducts
	return response, nil
}

func (s *ProductServer) QueryProductById(ctx context.Context, request *driverproduct.QueryProductByIdRequest) (*driverproduct.QueryProductByIdResponse, error) {
	productItf := container.ProductAppNameFrom(s.dic.Get)

	response := new(driverproduct.QueryProductByIdResponse)
	response.BaseResponse = new(drivercommon.CommonResponse)

	productModel, err := productItf.ProductModelById(ctx, request.Id)
	if err != nil {
		response.BaseResponse.Success = false
		response.BaseResponse.ErrorMessage = err.Error()
		return response, nil
	}
	response.Data = new(driverproduct.QueryProductByIdResponse_Data)
	response.Data.Product = productModel.TransformToDriverProduct()

	response.BaseResponse.Success = true
	return response, err

}

func NewProductServer(lc logger.LoggingClient, dic *di.Container) *ProductServer {
	return &ProductServer{
		lc:  lc,
		dic: dic,
	}
}

func (s *ProductServer) RegisterServer(server *grpc.Server) {
	driverproduct.RegisterRpcProductServer(server, s)
}
