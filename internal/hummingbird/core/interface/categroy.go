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

type CategoryApp interface {
	CategoryTemplateSearch(ctx context.Context, req dtos.CategoryTemplateRequest) ([]dtos.CategoryTemplateResponse, uint32, error)
	Sync(ctx context.Context, versionName string) (int64, error)
}

type UnitApp interface {
	UnitTemplateSearch(ctx context.Context, req dtos.UnitRequest) ([]dtos.UnitResponse, uint32, error)
	Sync(ctx context.Context, versionName string) (int64, error)
}

type DocsApp interface {
	SyncDocs(ctx context.Context, versionName string) (int64, error)
}

type QuickNavigation interface {
	SyncQuickNavigation(ctx context.Context, versionName string) (int64, error)
}

type ThingModelTemplateApp interface {
	ThingModelTemplateSearch(ctx context.Context, req dtos.ThingModelTemplateRequest) ([]dtos.ThingModelTemplateResponse, uint32, error)
	ThingModelTemplateByCategoryKey(ctx context.Context, categoryKey string) (dtos.ThingModelTemplateResponse, error)
	Sync(ctx context.Context, versionName string) (int64, error)
}
