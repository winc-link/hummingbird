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
	"github.com/winc-link/hummingbird/internal/pkg/errort"
	"github.com/winc-link/hummingbird/internal/pkg/utils"
	"github.com/winc-link/hummingbird/internal/tools/sqldb/sqlite"
	"gorm.io/gorm"
)

func languageSearch(c *Client, offset int, limit int, req dtos.LanguageSDKSearchQueryRequest) (languages []models.LanguageSdk, count uint32, edgeXErr error) {
	cs := models.LanguageSdk{}
	var total int64
	tx := c.Pool.Table(cs.TableName())
	tx = sqlite.BuildCommonCondition(tx, cs, req.BaseSearchConditionQuery)

	err := tx.Count(&total).Error
	if err != nil {
		return []models.LanguageSdk{}, 0, errort.NewCommonEdgeX(errort.DefaultSystemError, "language sdk failed query from the database", err)
	}

	err = tx.Offset(offset).Limit(limit).Find(&languages).Error
	if err != nil {
		return []models.LanguageSdk{}, 0, errort.NewCommonEdgeX(errort.DefaultSystemError, "language sdk failed query from the database", err)
	}

	return languages, uint32(total), nil
}

func languageByName(c *Client, name string) (language models.LanguageSdk, edgeXErr error) {
	if name == "" {
		return language, errort.NewCommonEdgeX(errort.DefaultIdEmpty, "language sdk name id is empty", nil)
	}
	err := c.client.GetObject(&models.LanguageSdk{Name: name}, &language)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return language, errort.NewCommonErr(errort.DefaultResourcesNotFound, fmt.Errorf("language sdk (%s) not found", name))
		}
		return language, err
	}
	return
}

func addLanguageSdk(c *Client, cs models.LanguageSdk) (language models.LanguageSdk, edgeXErr error) {
	ts := utils.MakeTimestamp()
	if cs.Created == 0 {
		cs.Created = ts
	}
	cs.Modified = ts

	err := c.client.CreateObject(&cs)
	if err != nil {
		edgeXErr = errort.NewCommonEdgeX(errort.DefaultSystemError, "language creation failed", err)
	}

	return cs, edgeXErr
}

func updateLanguageSdk(c *Client, dl models.LanguageSdk) error {
	dl.Modified = utils.MakeTimestamp()
	err := c.client.UpdateObject(&dl)
	if err != nil {
		return err
	}
	return nil
}
