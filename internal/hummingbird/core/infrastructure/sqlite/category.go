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
package sqlite

import (
	"fmt"
	"github.com/winc-link/hummingbird/internal/dtos"
	"github.com/winc-link/hummingbird/internal/models"
	"github.com/winc-link/hummingbird/internal/pkg/errort"
	"github.com/winc-link/hummingbird/internal/tools/sqldb/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func categoryTemplateSearch(c *Client, offset int, limit int, req dtos.CategoryTemplateRequest) (categoryTemplates []models.CategoryTemplate, count uint32, edgeXErr error) {
	cs := models.CategoryTemplate{}
	var total int64
	tx := c.Pool.Table(cs.TableName())
	tx = sqlite.BuildCommonCondition(tx, cs, req.BaseSearchConditionQuery)

	if req.CategoryName != "" {
		tx = tx.Where("`category_name` LIKE ?", "%"+req.CategoryName+"%")
	}

	if req.Scene != "" {
		tx = tx.Where("`scene` = ?", req.Scene)
	}

	err := tx.Count(&total).Error
	if err != nil {
		return []models.CategoryTemplate{}, 0, errort.NewCommonEdgeX(errort.DefaultSystemError, "categoryTemplate failed query from the database", err)
	}

	err = tx.Offset(offset).Limit(limit).Find(&categoryTemplates).Error
	if err != nil {
		return []models.CategoryTemplate{}, 0, errort.NewCommonEdgeX(errort.DefaultSystemError, "categoryTemplate failed query from the database", err)
	}

	return categoryTemplates, uint32(total), nil
}

func categoryTemplateById(c *Client, id string) (categoryTemplateInfo models.CategoryTemplate, edgeXErr error) {
	if id == "" {
		return categoryTemplateInfo, errort.NewCommonEdgeX(errort.DefaultIdEmpty, "categoryTemplate id is empty", nil)
	}
	err := c.client.GetObject(&models.CategoryTemplate{Id: id}, &categoryTemplateInfo)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return categoryTemplateInfo, errort.NewCommonErr(errort.CategoryNotExist, fmt.Errorf("categoryTemplate id(%s) not found", id))
		}
		return categoryTemplateInfo, err
	}
	return
}

func batchUpsertCategoryTemplate(c *Client, d []models.CategoryTemplate) (int64, error) {
	if len(d) <= 0 {
		return 0, nil
	}
	tx := c.Pool.Session(&gorm.Session{FullSaveAssociations: true}).Clauses(
		clause.OnConflict{
			UpdateAll: true,
		}).CreateInBatches(d, sqlite.CreateBatchSize)
	num := tx.RowsAffected
	err := tx.Error
	if err != nil {
		return num, err
	}
	return num, nil
}
