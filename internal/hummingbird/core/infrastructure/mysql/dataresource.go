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

package mysql

import (
	"fmt"
	"github.com/winc-link/hummingbird/internal/dtos"
	"github.com/winc-link/hummingbird/internal/models"
	"github.com/winc-link/hummingbird/internal/pkg/errort"
	"github.com/winc-link/hummingbird/internal/pkg/utils"
	"github.com/winc-link/hummingbird/internal/tools/sqldb/sqlite"
	"gorm.io/gorm"
)

func dataResourceById(c *Client, id string) (dateResource models.DataResource, edgeXErr error) {
	if id == "" {
		return dateResource, errort.NewCommonEdgeX(errort.DefaultIdEmpty, "dateResource id is empty", nil)
	}
	err := c.Pool.Table(dateResource.TableName()).First(&dateResource, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return dateResource, errort.NewCommonErr(errort.DefaultResourcesNotFound, fmt.Errorf("dateResource id(%s) not found", id))
		}
		return dateResource, errort.NewCommonErr(errort.DefaultSystemError, fmt.Errorf("query dateResource fail (Id:%s), %s", dateResource.Id, err))
	}
	return
}

func addDataResource(c *Client, ds models.DataResource) (id string, edgeXErr error) {
	ts := utils.MakeTimestamp()
	if ds.Created == 0 {
		ds.Created = ts
	}
	ds.Modified = ts
	err := c.client.CreateObject(&ds)
	if err != nil {
		edgeXErr = errort.NewCommonEdgeX(errort.DefaultSystemError, "data resourced creation failed", err)
	}

	return ds.Id, edgeXErr
}

func updateDataResource(c *Client, dl models.DataResource) error {
	dl.Modified = utils.MakeTimestamp()
	err := c.client.UpdateObject(&dl)
	if err != nil {
		return errort.NewCommonEdgeX(errort.DefaultSystemError, "data resource update failed", err)
	}
	return nil
}

func deleteDataResourceById(c *Client, id string) error {
	if id == "" {
		return errort.NewCommonEdgeX(errort.DefaultIdEmpty, "id is empty", nil)
	}
	err := c.client.DeleteObject(&models.DataResource{Id: id})
	if err != nil {
		return errort.NewCommonEdgeX(errort.DefaultSystemError, "data resourced deletion failed", err)
	}
	return nil
}

func updateDataResourceHealth(c *Client, id string, health bool) error {
	if id == "" {
		return errort.NewCommonEdgeX(errort.DefaultIdEmpty, "id is empty", nil)
	}
	d := models.DataResource{}
	tx := c.Pool.Table(d.TableName())
	err := tx.Where("id = ?", id).Updates(map[string]interface{}{"health": health}).Error
	if err != nil {
		return errort.NewCommonEdgeX(errort.DefaultSystemError, "update data resource failed", err)
	}
	return nil
}

func dataResourceSearch(c *Client, offset int, limit int, req dtos.DataResourceSearchQueryRequest) (dataResource []models.DataResource, count uint32, edgeXErr error) {
	dl := models.DataResource{}
	var total int64
	tx := c.Pool.Table(dl.TableName())
	tx = sqlite.BuildCommonCondition(tx, dl, req.BaseSearchConditionQuery)
	// 特殊条件
	if req.Type != "" {
		tx = tx.Where("`type` = ?", req.Type)
	}
	if req.Health != "" {
		isHealth := true
		if req.Health == SearchReqBoolTrue {
			isHealth = true
		} else {
			isHealth = false
		}
		tx = tx.Where("`health` = ?", isHealth)
	}
	err := tx.Count(&total).Error
	if err != nil {
		return []models.DataResource{}, 0, errort.NewCommonEdgeX(errort.DefaultSystemError, "data resource failed query from the database", err)
	}

	err = tx.Offset(offset).Limit(limit).Find(&dataResource).Error
	if err != nil {
		return []models.DataResource{}, 0, errort.NewCommonEdgeX(errort.DefaultSystemError, "data resource failed query from the database", err)
	}

	return dataResource, uint32(total), nil
}
