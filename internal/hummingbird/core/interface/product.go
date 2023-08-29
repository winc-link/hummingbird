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
package interfaces

import (
	"context"
	"github.com/winc-link/hummingbird/internal/dtos"
	"github.com/winc-link/hummingbird/internal/models"
)

type ProductItf interface {
	ProductCtlItf
}

type ProductCtlItf interface {
	ProductsSearch(ctx context.Context, req dtos.ProductSearchQueryRequest) ([]dtos.ProductSearchQueryResponse, uint32, error)
	ProductsModelSearch(ctx context.Context, req dtos.ProductSearchQueryRequest) ([]models.Product, uint32, error)
	ProductById(ctx context.Context, id string) (dtos.ProductSearchByIdResponse, error)
	ProductModelById(ctx context.Context, id string) (models.Product, error)
	ProductDelete(ctx context.Context, id string) error
	AddProduct(ctx context.Context, req dtos.ProductAddRequest) (string, error)
	ProductRelease(ctx context.Context, productId string) error
	ProductUnRelease(ctx context.Context, productId string) error
	CreateProductCallBack(productInfo models.Product)
	UpdateProductCallBack(productInfo models.Product)
	DeleteProductCallBack(productInfo models.Product)

	ProductCtlOpenApiItf
}

type ProductCtlOpenApiItf interface {
	OpenApiAddProduct(ctx context.Context, req dtos.OpenApiAddProductRequest) (string, error)
	OpenApiUpdateProduct(ctx context.Context, req dtos.OpenApiUpdateProductRequest) error
	OpenApiProductById(ctx context.Context, id string) (dtos.ProductSearchByIdOpenApiResponse, error)
	OpenApiProductSearch(ctx context.Context, req dtos.ProductSearchQueryRequest) ([]dtos.ProductSearchOpenApiResponse, uint32, error)
}
