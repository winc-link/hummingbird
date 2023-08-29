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
	"github.com/winc-link/hummingbird/internal/pkg/constants"
	"github.com/winc-link/hummingbird/internal/pkg/errort"
	"github.com/winc-link/hummingbird/internal/pkg/utils"
	"github.com/winc-link/hummingbird/internal/tools/sqldb/sqlite"
	"gorm.io/gorm"
)

func addScene(c *Client, ds models.Scene) (scene models.Scene, edgeXErr error) {
	ts := utils.MakeTimestamp()
	if ds.Created == 0 {
		ds.Created = ts
	}
	ds.Modified = ts

	err := c.client.CreateObject(&ds)
	if err != nil {
		edgeXErr = errort.NewCommonEdgeX(errort.DefaultSystemError, "scene creation failed", err)
	}
	return ds, edgeXErr
}

func updateScene(c *Client, dl models.Scene) error {
	dl.Modified = utils.MakeTimestamp()
	err := c.client.UpdateObject(&dl)
	if err != nil {
		return errort.NewCommonEdgeX(errort.DefaultSystemError, "scene update failed", err)
	}
	return nil
}

func sceneById(c *Client, id string) (scene models.Scene, err error) {
	if id == "" {
		return scene, errort.NewCommonEdgeX(errort.DefaultIdEmpty, "scene id is empty", nil)
	}
	err = c.client.GetObject(&models.Scene{Id: id}, &scene)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return scene, errort.NewCommonErr(errort.DefaultResourcesNotFound, fmt.Errorf("scene id(%s) not found", id))
		}
		return scene, err
	}
	return
}

func sceneStart(c *Client, id string) error {
	d := models.Scene{}
	tx := c.Pool.Table(d.TableName())
	err := tx.Where("id = ?", id).Updates(map[string]interface{}{"status": constants.SceneStart}).Error
	if err != nil {
		return errort.NewCommonEdgeX(errort.DefaultSystemError, "start scene rule failed", err)
	}
	return nil
}

func sceneStop(c *Client, id string) error {
	d := models.Scene{}
	tx := c.Pool.Table(d.TableName())
	err := tx.Where("id = ?", id).Updates(map[string]interface{}{"status": constants.SceneStop}).Error
	if err != nil {
		return errort.NewCommonEdgeX(errort.DefaultSystemError, "start scene rule failed", err)
	}
	return nil
}

func deleteSceneById(c *Client, id string) error {
	if id == "" {
		return errort.NewCommonEdgeX(errort.DefaultIdEmpty, "del scene id is empty", nil)
	}
	err := c.client.DeleteObject(&models.Scene{Id: id})
	if err != nil {
		return errort.NewCommonEdgeX(errort.DefaultSystemError, "del scene deletion failed", err)
	}
	return nil
}

func sceneSearch(c *Client, offset int, limit int, req dtos.SceneSearchQueryRequest) (scene []models.Scene, count uint32, edgeXErr error) {
	dp := models.Scene{}
	var total int64
	tx := c.Pool.Table(dp.TableName())
	tx = sqlite.BuildCommonCondition(tx, dp, req.BaseSearchConditionQuery)

	if req.Name != "" {
		tx = tx.Where("`name` LIKE ?", sqlite.MakeLikeParams(req.Name))
	}
	if req.Status != "" {
		tx = tx.Where("`status` = ?", req.Status)
	}
	err := tx.Count(&total).Error
	if err != nil {
		return scene, 0, errort.NewCommonEdgeX(errort.DefaultSystemError, "scene search failed query from the database", err)
	}

	err = tx.Offset(offset).Limit(limit).Find(&scene).Error
	if err != nil {
		return scene, 0, errort.NewCommonEdgeX(errort.DefaultSystemError, "scene search failed query from the database", err)
	}

	return scene, uint32(total), nil
}

func addSceneLog(c *Client, ds models.SceneLog) (sceneLog models.SceneLog, edgeXErr error) {
	ts := utils.MakeTimestamp()
	if ds.Created == 0 {
		ds.Created = ts
	}
	ds.Modified = ts

	err := c.client.CreateObject(&ds)
	if err != nil {
		edgeXErr = errort.NewCommonEdgeX(errort.DefaultSystemError, "scene log creation failed", err)
	}
	return ds, edgeXErr
}

func sceneLogSearch(c *Client, offset int, limit int, req dtos.SceneLogSearchQueryRequest) (sceneLogs []models.SceneLog, count uint32, edgeXErr error) {
	dp := models.SceneLog{}
	var total int64
	tx := c.Pool.Table(dp.TableName())
	tx = sqlite.BuildCommonCondition(tx, dp, req.BaseSearchConditionQuery)

	if req.StartAt > 0 && req.EndAt > 0 && req.EndAt-req.StartAt > 0 {
		tx.Where("created > ?", req.StartAt).Where("created < ?", req.EndAt)
	}
	if req.SceneId != "" {
		tx = tx.Where("`scene_id` = ?", req.SceneId)
	}

	err := tx.Count(&total).Error
	if err != nil {
		return sceneLogs, 0, errort.NewCommonEdgeX(errort.DefaultSystemError, "scene log search failed query from the database", err)
	}

	err = tx.Offset(offset).Limit(limit).Find(&sceneLogs).Error
	if err != nil {
		return sceneLogs, 0, errort.NewCommonEdgeX(errort.DefaultSystemError, "scene log search failed query from the database", err)
	}

	return sceneLogs, uint32(total), nil
}
