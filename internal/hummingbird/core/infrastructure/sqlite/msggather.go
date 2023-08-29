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
	"github.com/winc-link/hummingbird/internal/dtos"
	"github.com/winc-link/hummingbird/internal/models"
	"github.com/winc-link/hummingbird/internal/pkg/errort"
	"github.com/winc-link/hummingbird/internal/pkg/utils"
	"github.com/winc-link/hummingbird/internal/tools/sqldb/sqlite"
)

func addMsgGather(c *Client, msgGather models.MsgGather) error {
	ts := utils.MakeTimestamp()
	if msgGather.Created == 0 {
		msgGather.Created = ts
	}
	msgGather.Modified = ts

	err := c.client.CreateObject(&msgGather)
	if err != nil {
		edgeXErr := errort.NewCommonEdgeX(errort.DefaultSystemError, "add msg gather failed", err)
		return edgeXErr
	}

	return nil
}

func msgGatherSearch(c *Client, offset int, limit int, req dtos.MsgGatherSearchQueryRequest) (dcs []models.MsgGather, count uint32, edgeXErr error) {
	d := models.MsgGather{}
	var total int64
	tx := c.Pool.Table(d.TableName())
	tx = sqlite.BuildCommonCondition(tx, d, req.BaseSearchConditionQuery)

	if len(req.Date) > 0 {
		tx = tx.Where("`date` in (?)", req.Date)
	}
	err := tx.Count(&total).Error
	if err != nil {
		return []models.MsgGather{}, 0, errort.NewCommonEdgeX(errort.DefaultSystemError, "msg gather failed query from the database", err)
	}

	err = tx.Offset(offset).Limit(limit).Find(&dcs).Error
	if err != nil {
		return []models.MsgGather{}, 0, errort.NewCommonEdgeX(errort.DefaultSystemError, "msg gather failed query from the database", err)
	}

	return dcs, uint32(total), nil
}
