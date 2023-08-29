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
package driverapp

import (
	"context"
	"github.com/winc-link/hummingbird/internal/dtos"
)

func (app *driverLibApp) GetDriverClassify(ctx context.Context, req dtos.DriverClassifyQueryRequest) ([]dtos.DriverClassifyResponse, uint32, error) {
	offset, limit := req.BaseSearchConditionQuery.GetPage()
	dcs, total, err := app.dbClient.DriverClassifySearch(offset, limit, req)

	res := make([]dtos.DriverClassifyResponse, len(dcs))
	if err != nil {
		return res, 0, err
	}

	for i, dc := range dcs {
		res[i] = dtos.DriverClassifyResponse{
			Id:   dc.Id,
			Name: dc.Name,
		}
	}
	return res, total, nil
}
