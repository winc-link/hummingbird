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
	"github.com/winc-link/hummingbird/internal/models"
	"github.com/winc-link/hummingbird/internal/pkg/errort"
	"github.com/winc-link/hummingbird/internal/pkg/utils"
	"github.com/winc-link/hummingbird/internal/tools/sqldb/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func addThingModelProperty(c *Client, ds models.Properties) (models.Properties, error) {
	ts := utils.MakeTimestamp()
	if ds.Created == 0 {
		ds.Created = ts
	}
	ds.Modified = ts

	var edgeXErr error
	err := c.client.CreateObject(&ds)
	if err != nil {
		edgeXErr = errort.NewCommonEdgeX(errort.DefaultSystemError, "thing model property creation failed", err)
	}
	return ds, edgeXErr

}
func batchUpsertThingModel(c *Client, d interface{}) (int64, error) {
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

func addThingModelEvent(c *Client, ds models.Events) (models.Events, error) {
	ts := utils.MakeTimestamp()
	if ds.Created == 0 {
		ds.Created = ts
	}
	ds.Modified = ts

	var edgeXErr error
	err := c.client.CreateObject(&ds)
	if err != nil {
		edgeXErr = errort.NewCommonEdgeX(errort.DefaultSystemError, "thing model event creation failed", err)
	}

	return ds, edgeXErr

}

func addThingModelAction(c *Client, ds models.Actions) (models.Actions, error) {
	ts := utils.MakeTimestamp()
	if ds.Created == 0 {
		ds.Created = ts
	}
	ds.Modified = ts

	var edgeXErr error
	err := c.client.CreateObject(&ds)
	if err != nil {
		edgeXErr = errort.NewCommonEdgeX(errort.DefaultSystemError, "thing model action creation failed", err)
	}

	return ds, edgeXErr
}

func updateThingModelProperty(c *Client, ds models.Properties) error {
	ts := utils.MakeTimestamp()
	if ds.Created == 0 {
		ds.Created = ts
	}
	ds.Modified = ts
	err := c.client.UpdateObject(&ds)
	var edgeXErr error
	if err != nil {
		edgeXErr = errort.NewCommonEdgeX(errort.DefaultSystemError, "thing model property update failed", err)
	}

	return edgeXErr
}

func updateThingModelEvent(c *Client, ds models.Events) error {
	ts := utils.MakeTimestamp()
	if ds.Created == 0 {
		ds.Created = ts
	}
	ds.Modified = ts
	err := c.client.UpdateObject(&ds)
	var edgeXErr error
	if err != nil {
		edgeXErr = errort.NewCommonEdgeX(errort.DefaultSystemError, "thing model events update failed", err)
	}

	return edgeXErr
}

func updateThingModelAction(c *Client, ds models.Actions) error {
	ts := utils.MakeTimestamp()
	if ds.Created == 0 {
		ds.Created = ts
	}
	ds.Modified = ts
	err := c.client.UpdateObject(&ds)
	var edgeXErr error
	if err != nil {
		edgeXErr = errort.NewCommonEdgeX(errort.DefaultSystemError, "thing model action update failed", err)
	}

	return edgeXErr
}

func deleteThingModelPropertyById(c *Client, id string) error {
	if id == "" {
		return errort.NewCommonEdgeX(errort.DefaultIdEmpty, "properties id is empty", nil)
	}
	err := c.client.DeleteObject(&models.Properties{Id: id})
	if err != nil {
		return errort.NewCommonEdgeX(errort.DefaultSystemError, "properties deletion failed", err)
	}
	return nil
}

func deleteThingModelEventById(c *Client, id string) error {
	if id == "" {
		return errort.NewCommonEdgeX(errort.DefaultIdEmpty, "events id is empty", nil)
	}
	err := c.client.DeleteObject(&models.Events{Id: id})
	if err != nil {
		return errort.NewCommonEdgeX(errort.DefaultSystemError, "events deletion failed", err)
	}
	return nil
}

func deleteThingModelActionById(c *Client, id string) error {
	if id == "" {
		return errort.NewCommonEdgeX(errort.DefaultIdEmpty, "actions id is empty", nil)
	}
	err := c.client.DeleteObject(&models.Actions{Id: id})
	if err != nil {
		return errort.NewCommonEdgeX(errort.DefaultSystemError, "actions deletion failed", err)
	}
	return nil
}

func thingModelPropertyById(c *Client, id string) (models.Properties, error) {
	cs := models.Properties{}
	var properties models.Properties
	tx := c.Pool.Table(cs.TableName())
	tx.Where("id = ?", id)
	err := tx.Find(&properties).Error
	return properties, err
}

func thingModelEventById(c *Client, id string) (models.Events, error) {
	cs := models.Events{}
	var event models.Events
	tx := c.Pool.Table(cs.TableName())
	tx.Where("id = ?", id)
	err := tx.Find(&event).Error
	return event, err
}

func thingModeActionById(c *Client, id string) (models.Actions, error) {
	cs := models.Actions{}
	var action models.Actions
	tx := c.Pool.Table(cs.TableName())
	tx.Where("id = ?", id)
	err := tx.Find(&action).Error
	return action, err
}

func systemThingModelSearch(c *Client, modelType, modelName string) (interface{}, error) {
	switch modelType {
	case "property":
		cs := models.Properties{}
		var properties []models.Properties
		tx := c.Pool.Table(cs.TableName())
		if modelName != "" {
			tx.Where("system =1  and `name` LIKE ?", "%"+modelName+"%")
		} else {
			tx.Where("system =1")
		}
		err := tx.Find(&properties).Error
		return properties, err
	case "event":
		cs := models.Events{}
		var events []models.Events
		tx := c.Pool.Table(cs.TableName())
		if modelName != "" {
			tx.Where("system =1 and `name` LIKE ?", "%"+modelName+"%")
		} else {
			tx.Where("system =1")
		}
		err := tx.Find(&events).Error
		return events, err
	case "action":
		cs := models.Actions{}
		var actions []models.Actions
		tx := c.Pool.Table(cs.TableName())
		if modelName != "" {
			tx.Where("system =1 and `name` LIKE ?", "%"+modelName+"%")
		} else {
			tx.Where("system =1")
		}
		err := tx.Find(&actions).Error
		return actions, err

	}
	return nil, nil
}
