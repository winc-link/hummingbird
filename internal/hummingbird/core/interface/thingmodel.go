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
)

type ThingModelItf interface {
	ThingModelCtlItf
}

type ThingModelCtlItf interface {
	ThingModelDelete(ctx context.Context, id string, thingModelType string) error
	AddThingModel(ctx context.Context, req dtos.ThingModelAddOrUpdateReq) (string, error)
	UpdateThingModel(ctx context.Context, req dtos.ThingModelAddOrUpdateReq) error
	SystemThingModelSearch(ctx context.Context, req dtos.SystemThingModelSearchReq) (interface{}, error)
	OpenApiAddThingModel(ctx context.Context, req dtos.OpenApiThingModelAddOrUpdateReq) error
	OpenApiQueryThingModel(ctx context.Context, productId string) (dtos.OpenApiQueryThingModel, error)
	OpenApiDeleteThingModel(ctx context.Context, req dtos.OpenApiThingModelDeleteReq) error
}
