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

package sqlite

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

func addRuleEngine(c *Client, ds models.RuleEngine) (id string, edgeXErr error) {
	ts := utils.MakeTimestamp()
	if ds.Created == 0 {
		ds.Created = ts
	}
	ds.Modified = ts
	err := c.client.CreateObject(&ds)
	if err != nil {
		edgeXErr = errort.NewCommonEdgeX(errort.DefaultSystemError, "rule engine creation failed", err)
	}

	return ds.Id, edgeXErr
}

func ruleEngineById(c *Client, id string) (ruleEngine models.RuleEngine, edgeXErr error) {
	if id == "" {
		return ruleEngine, errort.NewCommonEdgeX(errort.DefaultIdEmpty, "rule engine id is empty", nil)
	}
	err := c.Pool.Table(ruleEngine.TableName()).Preload("DataResource").First(&ruleEngine, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return ruleEngine, errort.NewCommonErr(errort.DefaultSystemError, fmt.Errorf("rule engine id id(%s) not found", id))
		}
		return ruleEngine, errort.NewCommonErr(errort.DefaultSystemError, fmt.Errorf("query rule engine id fail (Id:%s), %s", ruleEngine.Id, err))
	}
	return
}

func ruleEngineSearch(c *Client, offset int, limit int, req dtos.RuleEngineSearchQueryRequest) (ruleEngine []models.RuleEngine, count uint32, edgeXErr error) {
	dp := models.RuleEngine{}
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
		return ruleEngine, 0, errort.NewCommonEdgeX(errort.DefaultSystemError, "rules engine failed query from the database", err)
	}

	err = tx.Offset(offset).Limit(limit).Preload("DataResource").Find(&ruleEngine).Error
	if err != nil {
		return ruleEngine, 0, errort.NewCommonEdgeX(errort.DefaultSystemError, "rules engine  failed query from the database", err)
	}

	return ruleEngine, uint32(total), nil
}

func ruleEngineStart(c *Client, id string) error {
	d := models.RuleEngine{}
	tx := c.Pool.Table(d.TableName())
	err := tx.Where("id = ?", id).Updates(map[string]interface{}{"status": constants.RuleStart}).Error
	if err != nil {
		return errort.NewCommonEdgeX(errort.DefaultSystemError, "start alert rule failed", err)
	}
	return nil
}

func ruleEngineStop(c *Client, id string) error {
	d := models.RuleEngine{}
	tx := c.Pool.Table(d.TableName())
	err := tx.Where("id = ?", id).Updates(map[string]interface{}{"status": constants.RuleStop}).Error
	if err != nil {
		return errort.NewCommonEdgeX(errort.DefaultSystemError, "stop alert rule failed", err)
	}
	return nil
}

func deleteRuleEngineById(c *Client, id string) error {
	if id == "" {
		return errort.NewCommonEdgeX(errort.DefaultIdEmpty, "rule engine id is empty", nil)
	}
	err := c.client.DeleteObject(&models.RuleEngine{Id: id})
	if err != nil {
		return errort.NewCommonEdgeX(errort.DefaultSystemError, "rule engine deletion failed", err)
	}
	return nil
}

func updateRuleEngine(c *Client, ruleEngine models.RuleEngine) error {
	ruleEngine.Modified = utils.MakeTimestamp()
	err := c.client.UpdateObject(&ruleEngine)
	if err != nil {
		return err
	}
	return nil
}
