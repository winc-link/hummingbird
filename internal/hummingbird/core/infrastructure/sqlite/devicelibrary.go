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
	"github.com/winc-link/hummingbird/internal/pkg/utils"
	"github.com/winc-link/hummingbird/internal/tools/sqldb/sqlite"
	"gorm.io/gorm"
)

func deviceLibraryById(c *Client, id string) (deviceLibrary models.DeviceLibrary, edgeXErr error) {
	if id == "" {
		return deviceLibrary, errort.NewCommonEdgeX(errort.DefaultIdEmpty, "deviceLibrary id is empty", nil)
	}
	err := c.client.GetObject(&models.DeviceLibrary{Id: id}, &deviceLibrary)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return deviceLibrary, errort.NewCommonErr(errort.DeviceLibraryNotExist, fmt.Errorf("device library id(%s) not found", id))
		}
		return deviceLibrary, err
	}
	return
}

func deleteDeviceLibraryById(c *Client, id string) error {
	if id == "" {
		return errort.NewCommonEdgeX(errort.DefaultIdEmpty, "id is empty", nil)
	}
	err := c.client.DeleteObject(&models.DeviceLibrary{Id: id})
	if err != nil {
		return errort.NewCommonEdgeX(errort.DefaultSystemError, "device library deletion failed", err)
	}
	return nil
}

func deviceLibraryIdExists(c *Client, id string) (bool, error) {
	exists, err := c.client.ExistObject(&models.DeviceLibrary{Id: id})
	if err != nil {
		return false, err
	}
	return exists, nil
}

func addDeviceLibrary(c *Client, dl models.DeviceLibrary) (models.DeviceLibrary, error) {
	// query device library name and id to avoid the conflict
	exists, edgeXErr := deviceLibraryIdExists(c, dl.Id)
	if edgeXErr != nil {
		return dl, edgeXErr
	} else if exists {
		return dl, errort.NewCommonEdgeX(errort.DefaultResourcesRepeat, fmt.Sprintf("device library id %s exists", dl.Id), edgeXErr)
	}

	// check docker config id exists
	exists, edgeXErr = dockerConfigIdExists(c, dl.DockerConfigId)
	if edgeXErr != nil {
		return dl, edgeXErr
	} else if !exists {
		return dl, errort.NewCommonEdgeX(errort.DockerImageRepositoryNotFound, fmt.Sprintf("docker config id %s not exists", dl.Id), edgeXErr)
	}

	ts := utils.MakeTimestamp()
	if dl.Created == 0 {
		dl.Created = ts
	}
	dl.Modified = ts

	err := c.client.CreateObject(&dl)
	if err != nil {
		edgeXErr = errort.NewCommonEdgeX(errort.DefaultSystemError, "device library creation failed", err)
	}

	return dl, edgeXErr
}

const (
	// url中请求参数有判断真假的
	SearchReqBoolTrue  = "true"
	SearchReqBoolFalse = "false"
)

func deviceLibrariesSearch(c *Client, offset int, limit int, req dtos.DeviceLibrarySearchQueryRequest) (deviceLibraries []models.DeviceLibrary, count uint32, edgeXErr error) {
	dl := models.DeviceLibrary{}
	var total int64
	tx := c.Pool.Table(dl.TableName())
	tx = sqlite.BuildCommonCondition(tx, dl, req.BaseSearchConditionQuery)
	// 特殊条件
	if req.DockerConfigId != "" {
		tx = tx.Where("`docker_config_id` = ?", req.DockerConfigId)
	}
	if req.IsInternal != "" {
		isInternal := true
		if req.IsInternal == SearchReqBoolTrue {
			isInternal = true
		} else {
			isInternal = false
		}
		tx = tx.Where("`is_internal` = ?", isInternal)
	}
	if req.DockerRepoName != "" {
		tx = tx.Where("`docker_repo_name` = ?", req.DockerRepoName)
	}
	if req.NameAliasLike != "" {
		tx = tx.Where("`name` LIKE ? OR `alias` LIKE ? OR `description` LIKE ?", sqlite.MakeLikeParams(req.NameAliasLike), sqlite.MakeLikeParams(req.NameAliasLike), sqlite.MakeLikeParams(req.NameAliasLike))
	}
	if req.NoInIds != "" {
		tx = tx.Where("`id` NOT IN ?", dtos.ApiParamsStringToArray(req.NoInIds))
	}
	if req.ImageIds != "" {
		tx = tx.Where("`docker_image_id` IN ?", dtos.ApiParamsStringToArray(req.ImageIds))
	}
	if req.NoInImageIds != "" {
		tx = tx.Where("`docker_image_id` NOT IN ?", dtos.ApiParamsStringToArray(req.NoInImageIds))
	}
	if req.DriverType != 0 {
		tx = tx.Where("`driver_type` = ?", req.DriverType)
	}
	if req.ClassifyId != 0 {
		tx = tx.Where("`classify_id` = ?", req.ClassifyId)
	}

	err := tx.Count(&total).Error
	if err != nil {
		return []models.DeviceLibrary{}, 0, errort.NewCommonEdgeX(errort.DefaultSystemError, "deviceLibraries failed query from the database", err)
	}

	err = tx.Offset(offset).Limit(limit).Order("id desc").Find(&deviceLibraries).Error
	if err != nil {
		return []models.DeviceLibrary{}, 0, errort.NewCommonEdgeX(errort.DefaultSystemError, "deviceLibraries failed query from the database", err)
	}

	return deviceLibraries, uint32(total), nil
}

func updateDeviceLibrary(c *Client, dl models.DeviceLibrary) error {
	dl.Modified = utils.MakeTimestamp()
	err := c.client.UpdateObject(&dl)
	if err != nil {
		return err
	}
	return nil
}
