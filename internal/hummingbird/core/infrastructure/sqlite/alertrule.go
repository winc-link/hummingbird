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
	"time"
)

func addAlertRule(c *Client, ds models.AlertRule) (alertRule models.AlertRule, edgeXErr error) {
	ts := utils.MakeTimestamp()
	if ds.Created == 0 {
		ds.Created = ts
	}
	ds.Modified = ts

	err := c.client.CreateObject(&ds)
	if err != nil {
		edgeXErr = errort.NewCommonEdgeX(errort.DefaultSystemError, "alert rule creation failed", err)
	}
	return ds, edgeXErr
}

func addAlertList(c *Client, ds models.AlertList) (alertRule models.AlertList, edgeXErr error) {
	ts := utils.MakeTimestamp()
	if ds.Created == 0 {
		ds.Created = ts
	}
	ds.Modified = ts

	err := c.client.CreateObject(&ds)
	if err != nil {
		edgeXErr = errort.NewCommonEdgeX(errort.DefaultSystemError, "alert rule list creation failed", err)
	}
	return ds, edgeXErr
}

func updateAlertRule(c *Client, dl models.AlertRule) error {
	dl.Modified = utils.MakeTimestamp()
	err := c.client.UpdateObject(&dl)
	if err != nil {
		return errort.NewCommonEdgeX(errort.DefaultSystemError, "alert rule update failed", err)
	}
	return nil
}

func alertRuleById(c *Client, id string) (alertRule models.AlertRule, edgeXErr error) {
	if id == "" {
		return alertRule, errort.NewCommonEdgeX(errort.DefaultIdEmpty, "alert rule id is empty", nil)
	}
	err := c.Pool.Table(alertRule.TableName()).First(&alertRule, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return alertRule, errort.NewCommonErr(errort.AlertRuleNotExist, fmt.Errorf("alert rule id(%s) not found", id))
		}
		return alertRule, errort.NewCommonErr(errort.DefaultSystemError, fmt.Errorf("query alert rule fail (Id:%s), %s", alertRule.Id, err))
	}
	return
}

func alertRuleSearch(c *Client, offset int, limit int, req dtos.AlertRuleSearchQueryRequest) (alertRules []models.AlertRule, count uint32, edgeXErr error) {
	dp := models.AlertRule{}
	var total int64
	tx := c.Pool.Table(dp.TableName())
	tx = sqlite.BuildCommonCondition(tx, dp, req.BaseSearchConditionQuery)

	if req.Name != "" {
		tx = tx.Where("`name` = ?", req.Name)
	}
	if req.Status != "" {
		tx = tx.Where("`status` = ?", req.Status)
	}

	err := tx.Count(&total).Error
	if err != nil {
		return []models.AlertRule{}, 0, errort.NewCommonEdgeX(errort.DefaultSystemError, "alert rules failed query from the database", err)
	}

	err = tx.Offset(offset).Limit(limit).Find(&alertRules).Error
	if err != nil {
		return []models.AlertRule{}, 0, errort.NewCommonEdgeX(errort.DefaultSystemError, "alert rules failed query from the database", err)
	}

	return alertRules, uint32(total), nil
}

func alertListLastSend(c *Client, alertRuleId string) (alertList models.AlertList, edgeXErr error) {
	al := models.AlertList{}
	err := c.Pool.Table(al.TableName()).Where("alert_rule_id = ?", alertRuleId).Where("is_send", true).Order("created desc").Last(&alertList).Error
	if err != nil {
		return
	}
	return
}

func alertListSearch(c *Client, offset int, limit int, req dtos.AlertSearchQueryRequest) (alertRules []dtos.AlertSearchQueryResponse, count uint32, edgeXErr error) {
	var total int64
	dp := models.AlertList{}
	tx := c.Pool.Table(dp.TableName()).Select("alert_list.id,alert_list.status," +
		"alert_rule.name,alert_list.alert_result,alert_rule.alert_level,alert_list.trigger_time,alert_list.treated_time,alert_list.message,alert_list.is_send").Joins("left join alert_rule on alert_list.alert_rule_id = alert_rule.id")
	//tx = sqlite.BuildCommonCondition(tx, dp, req.BaseSearchConditionQuery)
	if req.Name != "" {
		tx.Where("alert_rule.name LIKE ?", sqlite.MakeLikeParams(req.Name))
	}
	if req.Status != "" {
		tx.Where("alert_list.status = ?", req.Status)

	}
	if req.AlertLevel != "" {
		tx.Where("alert_rule.alert_level = ?", req.AlertLevel)
	}
	if req.TriggerStartTime > 0 && req.TriggerEndTime > 0 && req.TriggerEndTime-req.TriggerStartTime > 0 {
		tx.Where("alert_list.trigger_time >= ?", req.TriggerStartTime)
		tx.Where("alert_list.trigger_time <= ?", req.TriggerEndTime)
	}
	edgeXErr = tx.Count(&total).Error
	if edgeXErr != nil {
		return []dtos.AlertSearchQueryResponse{}, 0, errort.NewCommonEdgeX(errort.DefaultSystemError, "alert list failed query from the database", edgeXErr)
	}
	tx.Order("alert_list.created desc")
	edgeXErr = tx.Offset(offset).Limit(limit).Scan(&alertRules).Error
	if edgeXErr != nil {
		return []dtos.AlertSearchQueryResponse{}, 0, errort.NewCommonEdgeX(errort.DefaultSystemError, "alert list failed query from the database", edgeXErr)
	}
	return alertRules, uint32(total), nil
}

func deleteAlertRuleById(c *Client, id string) error {
	if id == "" {
		return errort.NewCommonEdgeX(errort.DefaultIdEmpty, "alert rule id is empty", nil)
	}
	err := c.client.DeleteObject(&models.AlertRule{Id: id})
	if err != nil {
		return errort.NewCommonEdgeX(errort.DefaultSystemError, "alert rule deletion failed", err)
	}
	return nil
}

func alertRuleStart(c *Client, id string) error {
	d := models.AlertRule{}
	tx := c.Pool.Table(d.TableName())
	err := tx.Where("id = ?", id).Updates(map[string]interface{}{"status": constants.RuleStart}).Error
	if err != nil {
		return errort.NewCommonEdgeX(errort.DefaultSystemError, "start alert rule failed", err)
	}
	return nil
}

func alertRuleStop(c *Client, id string) error {
	d := models.AlertRule{}
	tx := c.Pool.Table(d.TableName())
	err := tx.Where("id = ?", id).Updates(map[string]interface{}{"status": constants.RuleStop}).Error
	if err != nil {
		return errort.NewCommonEdgeX(errort.DefaultSystemError, "stop alert rule failed", err)
	}
	return nil
}

//subQuery := db.Select("AVG(age)").Where("name LIKE ?", "name%").Table("users")
//db.Select("AVG(age) as avgage").Group("name").Having("AVG(age) > (?)", subQuery).Find(&results)
// SELECT AVG(age) as avgage FROM `users` GROUP BY `name` HAVING AVG(age) > (SELECT AVG(age) FROM `users` WHERE name LIKE "name%")

func alertPlate(c *Client, beforeTime int64) (plate []dtos.AlertPlateQueryResponse, err error) {
	d := models.AlertList{}
	if beforeTime > 0 {
		err = c.Pool.Table(d.TableName()).Raw(
			"SELECT count(alert_list.id) AS count,alert_rule.alert_level FROM alert_list "+
				"JOIN alert_rule on alert_list.alert_rule_id = alert_rule.id and alert_list.created > (?) "+
				"GROUP BY alert_rule.alert_level", beforeTime).Scan(&plate).Error
	} else {
		err = c.Pool.Table(d.TableName()).Raw(
			"SELECT count(alert_list.id) AS count,alert_rule.alert_level FROM alert_list " +
				"JOIN alert_rule on alert_list.alert_rule_id = alert_rule.id" +
				"GROUP BY alert_rule.alert_level").Scan(&plate).Error
	}

	return
}

func alertIgnore(c *Client, id string) (err error) {
	d := models.AlertList{}
	tx := c.Pool.Table(d.TableName())
	err = tx.Where("id = ?", id).Updates(map[string]interface{}{"status": constants.Ignore, "treated_time": time.Now().UnixMilli()}).Error
	if err != nil {
		return errort.NewCommonEdgeX(errort.DefaultSystemError, "alert ignore rule failed", err)
	}
	return nil
}

func treatedIgnore(c *Client, id string, message string) (err error) {
	d := models.AlertList{}
	tx := c.Pool.Table(d.TableName())
	err = tx.Where("id = ?", id).Updates(map[string]interface{}{"status": constants.Treated, "message": message, "treated_time": time.Now().UnixMilli()}).Error
	if err != nil {
		return errort.NewCommonEdgeX(errort.DefaultSystemError, "alert ignore rule failed", err)
	}
	return nil
}
