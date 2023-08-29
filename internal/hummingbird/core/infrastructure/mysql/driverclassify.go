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
package mysql

import (
	"github.com/winc-link/hummingbird/internal/dtos"
	"github.com/winc-link/hummingbird/internal/models"
	"github.com/winc-link/hummingbird/internal/pkg/errort"
	"github.com/winc-link/hummingbird/internal/tools/sqldb/sqlite"
)

func driverClassifySearch(c *Client, offset int, limit int, req dtos.DriverClassifyQueryRequest) (dcs []models.DriverClassify, count uint32, edgeXErr error) {
	d := models.DriverClassify{}
	var total int64
	tx := c.Pool.Table(d.TableName())
	tx = sqlite.BuildCommonCondition(tx, d, req.BaseSearchConditionQuery)
	if req.Name != "" {
		tx = tx.Where("`name` = ?", req.Name)
	}
	err := tx.Count(&total).Error
	if err != nil {
		return []models.DriverClassify{}, 0, errort.NewCommonEdgeX(errort.DefaultSystemError, "device failed query from the database", err)
	}
	err = tx.Offset(offset).Limit(limit).Find(&dcs).Error
	if err != nil {
		return []models.DriverClassify{}, 0, errort.NewCommonEdgeX(errort.DefaultSystemError, "device failed query from the database", err)
	}
	return dcs, uint32(total), nil
}
