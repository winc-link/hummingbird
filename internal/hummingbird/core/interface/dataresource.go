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

package interfaces

import (
	"context"
	"github.com/winc-link/hummingbird/internal/dtos"
	"github.com/winc-link/hummingbird/internal/models"
	"github.com/winc-link/hummingbird/internal/pkg/constants"
)

type DataResourceApp interface {
	AddDataResource(ctx context.Context, req dtos.AddDataResourceReq) (string, error)
	DataResourceById(ctx context.Context, id string) (models.DataResource, error)
	UpdateDataResource(ctx context.Context, req dtos.UpdateDataResource) error
	DelDataResourceById(ctx context.Context, id string) error
	DataResourceSearch(ctx context.Context, req dtos.DataResourceSearchQueryRequest) ([]models.DataResource, uint32, error)
	DataResourceType(ctx context.Context) []constants.DataResourceType
	DataResourceHealth(ctx context.Context, resourceId string) error
}
