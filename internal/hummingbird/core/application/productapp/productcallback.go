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
package productapp

import (
	"context"
	"github.com/winc-link/edge-driver-proto/productcallback"
	"github.com/winc-link/hummingbird/internal/dtos"
	"github.com/winc-link/hummingbird/internal/models"
	"github.com/winc-link/hummingbird/internal/pkg/constants"
	"github.com/winc-link/hummingbird/internal/pkg/di"
	"time"

	//"github.com/winc-link/hummingbird/internal/hummingbird/core/application/driverapp"
	"github.com/winc-link/hummingbird/internal/hummingbird/core/container"
	"github.com/winc-link/hummingbird/internal/tools/rpcclient"
)

func (p *productApp) CreateProductCallBack(productInfo models.Product) {
	deviceServices, total, err := p.dbClient.DeviceServicesSearch(0, -1, dtos.DeviceServiceSearchQueryRequest{})
	if err != nil {
		return
	}
	if total == 0 {
		return
	}
	for _, service := range deviceServices {
		if productInfo.Platform == service.Platform {
			driverService := container.DriverServiceAppFrom(di.GContainer.Get)
			status := driverService.GetState(service.Id)
			if status == constants.RunStatusStarted {
				client, errX := rpcclient.NewDriverRpcClient(service.BaseAddress, false, "", service.Id, p.lc)
				if errX != nil {
					return
				}
				defer client.Close()
				var rpcRequest productcallback.CreateProductCallbackRequest
				rpcRequest.Data = productInfo.TransformToDriverProduct()
				rpcRequest.HappenTime = uint64(time.Now().Unix())
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
				defer cancel()
				_, _ = client.ProductCallBackServiceClient.CreateProductCallback(ctx, &rpcRequest)
			}
		}
	}
}

func (p *productApp) UpdateProductCallBack(productInfo models.Product) {
	//p.lc.Infof("UpdateProductCallBack Platform :%s name :%s id :%s", productInfo.Platform, productInfo.Name, productInfo.Id)
	deviceServices, total, err := p.dbClient.DeviceServicesSearch(0, -1, dtos.DeviceServiceSearchQueryRequest{})
	if err != nil {
		return
	}
	if total == 0 {
		return
	}
	for _, service := range deviceServices {
		if productInfo.Platform == service.Platform {
			driverService := container.DriverServiceAppFrom(di.GContainer.Get)
			status := driverService.GetState(service.Id)
			if status == constants.RunStatusStarted {
				client, errX := rpcclient.NewDriverRpcClient(service.BaseAddress, false, "", service.Id, p.lc)
				if errX != nil {
					return
				}
				defer client.Close()
				var rpcRequest productcallback.UpdateProductCallbackRequest
				rpcRequest.Data = productInfo.TransformToDriverProduct()
				rpcRequest.HappenTime = uint64(time.Now().Second())
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
				defer cancel()
				_, _ = client.ProductCallBackServiceClient.UpdateProductCallback(ctx, &rpcRequest)
			}
		}
	}
}

func (p *productApp) DeleteProductCallBack(productInfo models.Product) {
	//p.lc.Infof("DeleteProductCallBack Platform :%s name :%s id :%s", productInfo.Platform, productInfo.Name, productInfo.Id)
	deviceServices, total, err := p.dbClient.DeviceServicesSearch(0, -1, dtos.DeviceServiceSearchQueryRequest{})
	if total == 0 {
		return
	}
	if err != nil {
		return
	}

	for _, service := range deviceServices {
		if productInfo.Platform == service.Platform {
			driverService := container.DriverServiceAppFrom(di.GContainer.Get)
			status := driverService.GetState(service.Id)
			if status == constants.RunStatusStarted {
				client, errX := rpcclient.NewDriverRpcClient(service.BaseAddress, false, "", service.Id, p.lc)
				if errX != nil {
					return
				}
				defer client.Close()
				var rpcRequest productcallback.DeleteProductCallbackRequest
				rpcRequest.ProductId = productInfo.Id
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
				defer cancel()
				_, _ = client.ProductCallBackServiceClient.DeleteProductCallback(ctx, &rpcRequest)
			}
		}
	}
}
