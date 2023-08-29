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
	"fmt"
	"github.com/winc-link/hummingbird/internal/dtos"
	"github.com/winc-link/hummingbird/internal/models"
	"github.com/winc-link/hummingbird/internal/pkg/constants"
	"github.com/winc-link/hummingbird/internal/pkg/errort"
	"github.com/winc-link/hummingbird/internal/pkg/utils"
	"github.com/winc-link/hummingbird/internal/tools/sqldb/sqlite"
	"gorm.io/gorm"
)

func deviceServicesSearch(c *Client, offset int, limit int, req dtos.DeviceServiceSearchQueryRequest) (deviceServices []models.DeviceService, count uint32, edgeXErr error) {
	ds := models.DeviceService{}
	var total int64
	tx := c.Pool.Table(ds.TableName())
	tx = sqlite.BuildCommonCondition(tx, ds, req.BaseSearchConditionQuery)

	if req.DeviceLibraryId != "" {
		tx = tx.Where("`device_library_id` = ?", req.DeviceLibraryId)
	}
	if req.DeviceLibraryIds != "" {
		tx = tx.Where("`device_library_id` IN ?", dtos.ApiParamsStringToArray(req.DeviceLibraryIds))
	}
	if req.DriverType != 0 {
		tx = tx.Where("`driver_type` = ?", req.DriverType)
	}
	if req.CloudProductId != "" {
		tx = tx.Where("`cloud_product_id` = ?", req.CloudProductId)
	}
	if req.ProductId != "" {
		tx = tx.Where("`product_id` = ?", req.ProductId)
	}
	if req.Platform != "" {
		tx = tx.Where("`platform` = ?", req.Platform)

	}
	err := tx.Count(&total).Error
	if err != nil {
		return []models.DeviceService{}, 0, errort.NewCommonEdgeX(errort.DefaultSystemError, "deviceServices failed query from the database", err)
	}

	err = tx.Offset(offset).Limit(limit).Find(&deviceServices).Error
	if err != nil {
		return []models.DeviceService{}, 0, errort.NewCommonEdgeX(errort.DefaultSystemError, "deviceServices failed query from the database", err)
	}

	return deviceServices, uint32(total), nil
}

func deviceServiceIdExist(c *Client, id string) (bool, error) {
	exists, err := c.client.ExistObject(&models.DeviceService{Id: id})
	if err != nil {
		return false, err
	}
	return exists, nil
}

func addDeviceService(c *Client, ds models.DeviceService) (addedDeviceService models.DeviceService, edgeXErr error) {
	exists, edgeXErr := deviceServiceIdExist(c, ds.Id)
	if edgeXErr != nil {
		return ds, edgeXErr
	} else if exists {
		return ds, errort.NewCommonEdgeX(errort.DefaultResourcesRepeat, fmt.Sprintf("device service id %s exists", ds.Id), edgeXErr)
	}

	exists, edgeXErr = deviceServiceIdExist(c, ds.Name)
	if edgeXErr != nil {
		return ds, edgeXErr
	} else if exists {
		return ds, errort.NewCommonEdgeX(errort.DefaultResourcesRepeat, fmt.Sprintf("device service name %s exists", ds.Name), edgeXErr)
	}

	ts := utils.MakeTimestamp()
	if ds.Created == 0 {
		ds.Created = ts
	}
	ds.Modified = ts

	err := c.client.CreateObject(&ds)
	if err != nil {
		edgeXErr = errort.NewCommonEdgeX(errort.DefaultSystemError, "device service creation failed", err)
	}

	return ds, edgeXErr
}

func updateDeviceService(c *Client, ds models.DeviceService) error {
	ds.Modified = utils.MakeTimestamp()
	err := c.client.UpdateObject(&ds)
	if err != nil {
		return err
	}
	return nil
}

func deviceServiceById(c *Client, id string) (deviceService models.DeviceService, edgeXErr error) {
	if id == "" {
		return deviceService, errort.NewCommonEdgeX(errort.DefaultIdEmpty, "device service id is empty", nil)
	}
	err := c.client.GetObject(&models.DeviceService{Id: id}, &deviceService)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return deviceService, errort.NewCommonErr(errort.DeviceServiceNotExist, fmt.Errorf("device service id(%s) not found", id))
		}
		return deviceService, err
	}
	return
}

func deleteDeviceServiceById(c *Client, id string) error {
	if id == "" {
		return errort.NewCommonEdgeX(errort.DefaultIdEmpty, "device service id is empty", nil)
	}
	err := c.client.DeleteObject(&models.DeviceService{Id: id})
	if err != nil {
		return errort.NewCommonEdgeX(errort.DefaultSystemError, "device service deletion failed", err)
	}
	return nil
}

func driverMqttAuthInfo(c *Client, id string) (mqttAuth models.MqttAuth, edgeXErr error) {
	if id == "" {
		return mqttAuth, errort.NewCommonEdgeX(errort.DefaultIdEmpty, "device id is empty", nil)
	}
	err := c.client.GetObject(&models.MqttAuth{ResourceId: id, ResourceType: constants.DriverResource}, &mqttAuth)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return mqttAuth, errort.NewCommonErr(errort.DefaultResourcesNotFound, fmt.Errorf("mqtt auth resoure id(%s) not found", id))
		}
		return mqttAuth, errort.NewCommonErr(errort.DefaultSystemError, fmt.Errorf("query mqtt auth fail (resoureId:%s), %s", id, err))
	}
	return
}
