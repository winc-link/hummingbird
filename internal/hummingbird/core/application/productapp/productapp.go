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
	"errors"
	"github.com/winc-link/hummingbird/internal/dtos"
	resourceContainer "github.com/winc-link/hummingbird/internal/hummingbird/core/container"
	interfaces "github.com/winc-link/hummingbird/internal/hummingbird/core/interface"
	"github.com/winc-link/hummingbird/internal/models"
	"github.com/winc-link/hummingbird/internal/pkg/constants"
	"github.com/winc-link/hummingbird/internal/pkg/container"
	"github.com/winc-link/hummingbird/internal/pkg/di"
	"github.com/winc-link/hummingbird/internal/pkg/errort"
	"github.com/winc-link/hummingbird/internal/pkg/logger"
	"github.com/winc-link/hummingbird/internal/pkg/utils"
)

type productApp struct {
	//*propertyTyApp
	dic      *di.Container
	dbClient interfaces.DBClient
	lc       logger.LoggingClient
}

func NewProductApp(ctx context.Context, dic *di.Container) interfaces.ProductItf {
	lc := container.LoggingClientFrom(dic.Get)
	dbClient := resourceContainer.DBClientFrom(dic.Get)

	return &productApp{
		dic:      dic,
		dbClient: dbClient,
		lc:       lc,
	}
}

func (p *productApp) ProductsSearch(ctx context.Context, req dtos.ProductSearchQueryRequest) ([]dtos.ProductSearchQueryResponse, uint32, error) {
	offset, limit := req.BaseSearchConditionQuery.GetPage()

	resp, total, err := p.dbClient.ProductsSearch(offset, limit, false, req)
	if err != nil {
		return []dtos.ProductSearchQueryResponse{}, 0, err
	}
	products := make([]dtos.ProductSearchQueryResponse, len(resp))
	for i, p := range resp {
		products[i] = dtos.ProductResponseFromModel(p)
	}
	return products, total, nil
}

func (p *productApp) ProductsModelSearch(ctx context.Context, req dtos.ProductSearchQueryRequest) ([]models.Product, uint32, error) {
	offset, limit := req.BaseSearchConditionQuery.GetPage()

	return p.dbClient.ProductsSearch(offset, limit, true, req)

}

func (p *productApp) ProductById(ctx context.Context, id string) (dtos.ProductSearchByIdResponse, error) {
	resp, err := p.dbClient.ProductById(id)
	if err != nil {
		return dtos.ProductSearchByIdResponse{}, err
	}
	return dtos.ProductSearchByIdFromModel(resp), nil
}

func (p *productApp) OpenApiProductById(ctx context.Context, id string) (dtos.ProductSearchByIdOpenApiResponse, error) {
	resp, err := p.dbClient.ProductById(id)
	if err != nil {
		return dtos.ProductSearchByIdOpenApiResponse{}, err
	}
	return dtos.ProductSearchByIdOpenApiFromModel(resp), nil
}

func (p *productApp) OpenApiProductSearch(ctx context.Context, req dtos.ProductSearchQueryRequest) ([]dtos.ProductSearchOpenApiResponse, uint32, error) {
	offset, limit := req.BaseSearchConditionQuery.GetPage()

	resp, total, err := p.dbClient.ProductsSearch(offset, limit, false, req)
	if err != nil {
		return []dtos.ProductSearchOpenApiResponse{}, 0, err
	}
	products := make([]dtos.ProductSearchOpenApiResponse, len(resp))
	for i, product := range resp {
		products[i] = dtos.ProductSearchOpenApiFromModel(product)
	}
	return products, total, nil
}

func (p *productApp) ProductModelById(ctx context.Context, id string) (models.Product, error) {
	resp, err := p.dbClient.ProductById(id)
	if err != nil {
		return models.Product{}, err
	}
	return resp, nil

}

func (p *productApp) ProductDelete(ctx context.Context, id string) error {
	productInfo, err := p.dbClient.ProductById(id)
	if err != nil {
		return err
	}
	_, total, err := p.dbClient.DevicesSearch(0, -1, dtos.DeviceSearchQueryRequest{ProductId: productInfo.Id})

	if err != nil {
		return err
	}
	if total > 0 {
		return errort.NewCommonEdgeX(errort.ProductMustDeleteDevice, "该产品已绑定子设备，请优先删除子设备", err)
	}
	alertApp := resourceContainer.AlertRuleAppNameFrom(p.dic.Get)
	err = alertApp.CheckRuleByProductId(ctx, id)
	if err != nil {
		return err
	}
	if err = p.dbClient.AssociationsDeleteProductObject(productInfo); err != nil {
		return err
	}
	_ = resourceContainer.DataDBClientFrom(p.dic.Get).DropStable(ctx, productInfo.Id)
	go func() {
		p.DeleteProductCallBack(models.Product{
			Id:       productInfo.Id,
			Platform: productInfo.Platform,
		})
	}()
	return nil
}

func (p *productApp) AddProduct(ctx context.Context, req dtos.ProductAddRequest) (productId string, err error) {
	// 标准品类
	var properties []models.Properties
	var events []models.Events
	var actions []models.Actions

	if req.CategoryTemplateId != "1" {
		categoryTempInfo, err := p.dbClient.CategoryTemplateById(req.CategoryTemplateId)
		if err != nil {
			return "", err
		}
		thingModelTemplateInfo, err := p.dbClient.ThingModelTemplateByCategoryKey(categoryTempInfo.CategoryKey)
		if err != nil {
			return "", err
		}
		if thingModelTemplateInfo.ThingModelJSON != "" {
			properties, events, actions = dtos.GetModelPropertyEventActionByThingModelTemplate(thingModelTemplateInfo.ThingModelJSON)
		}
	}
	secret := utils.GenerateDeviceSecret(15)
	var insertProduct models.Product
	insertProduct.Id = utils.RandomNum()
	insertProduct.Name = req.Name
	insertProduct.CloudProductId = secret
	insertProduct.Platform = constants.IotPlatform_LocalIot
	insertProduct.Protocol = req.Protocol
	insertProduct.NodeType = constants.ProductNodeType(req.NodeType)
	insertProduct.NetType = constants.ProductNetType(req.NetType)
	insertProduct.DataFormat = req.DataFormat
	insertProduct.Factory = req.Factory
	insertProduct.Description = req.Description
	insertProduct.Key = secret
	insertProduct.Status = constants.ProductUnRelease
	insertProduct.Properties = properties
	insertProduct.Events = events
	insertProduct.Actions = actions

	ps, err := p.dbClient.AddProduct(insertProduct)
	if err != nil {
		return "", err
	}
	go func() {
		p.CreateProductCallBack(insertProduct)
	}()
	return ps.Id, nil
}

func (p *productApp) ProductRelease(ctx context.Context, productId string) error {

	var err error
	var productInfo models.Product

	productInfo, err = p.dbClient.ProductById(productId)

	if err != nil {
		return err
	}
	if productInfo.Status == constants.ProductRelease {
		return errors.New("")
	}

	err = resourceContainer.DataDBClientFrom(p.dic.Get).CreateStable(ctx, productInfo)
	if err != nil {
		return err
	}
	productInfo.Status = constants.ProductRelease
	return p.dbClient.UpdateProduct(productInfo)
}

func (p *productApp) ProductUnRelease(ctx context.Context, productId string) error {
	var err error
	var productInfo models.Product

	productInfo, err = p.dbClient.ProductById(productId)

	if err != nil {
		return err
	}
	if productInfo.Status == constants.ProductUnRelease {
		return errors.New("")
	}
	productInfo.Status = constants.ProductUnRelease
	return p.dbClient.UpdateProduct(productInfo)
}

func (p *productApp) OpenApiAddProduct(ctx context.Context, req dtos.OpenApiAddProductRequest) (productId string, err error) {
	//var properties []models.Properties
	//var events []models.Events
	//var actions []models.Actions
	//properties, events, actions = dtos.OpenApiGetModelPropertyEventActionByThingModelTemplate(req)

	//if len(properties) == 0 {
	//	properties = make([]models.Properties, 0)
	//}
	//
	//if len(events) == 0 {
	//	events = make([]models.Events, 0)
	//}
	//
	//if len(actions) == 0 {
	//	actions = make([]models.Actions, 0)
	//}
	var insertProduct models.Product
	insertProduct.Name = req.Name
	insertProduct.CloudProductId = utils.GenerateDeviceSecret(15)
	insertProduct.Platform = constants.IotPlatform_LocalIot
	insertProduct.Protocol = req.Protocol
	insertProduct.NodeType = constants.ProductNodeType(req.NodeType)
	insertProduct.NetType = constants.ProductNetType(req.NetType)
	insertProduct.DataFormat = req.DataFormat
	insertProduct.Factory = req.Factory
	insertProduct.Description = req.Description
	//insertProduct.Properties = properties
	//insertProduct.Events = events
	//insertProduct.Actions = actions

	ps, err := p.dbClient.AddProduct(insertProduct)
	if err != nil {
		return "", err
	}
	go func() {
		p.CreateProductCallBack(insertProduct)
	}()
	return ps.Id, nil
}

func (p *productApp) OpenApiUpdateProduct(ctx context.Context, req dtos.OpenApiUpdateProductRequest) error {

	product, err := p.dbClient.ProductById(req.Id)
	if err != nil {
		return err
	}
	if req.Name != nil {
		product.Name = *req.Name
	}
	product.Platform = constants.IotPlatform_LocalIot
	if req.Protocol != nil {
		product.Protocol = *req.Protocol
	}

	if req.NetType != nil {
		product.NodeType = constants.ProductNodeType(*req.NodeType)
	}

	if req.NetType != nil {
		product.NetType = constants.ProductNetType(*req.NetType)
	}

	if req.DataFormat != nil {
		product.DataFormat = *req.DataFormat
	}

	if req.Factory != nil {
		product.Factory = *req.Factory
	}

	if req.Description != nil {
		product.Description = *req.Description
	}

	err = p.dbClient.UpdateProduct(product)
	if err != nil {
		return err
	}
	go func() {
		p.UpdateProductCallBack(product)
	}()
	return nil
}
