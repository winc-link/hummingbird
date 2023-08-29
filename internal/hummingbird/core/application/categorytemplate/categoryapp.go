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
package categorytemplate

import (
	"context"
	"encoding/json"
	"github.com/winc-link/hummingbird/internal/dtos"
	"github.com/winc-link/hummingbird/internal/models"
	"github.com/winc-link/hummingbird/internal/pkg/container"
	"github.com/winc-link/hummingbird/internal/pkg/utils"
	"time"

	resourceContainer "github.com/winc-link/hummingbird/internal/hummingbird/core/container"
	interfaces "github.com/winc-link/hummingbird/internal/hummingbird/core/interface"
	"github.com/winc-link/hummingbird/internal/pkg/di"
	"github.com/winc-link/hummingbird/internal/pkg/logger"
)

type categoryApp struct {
	dic      *di.Container
	dbClient interfaces.DBClient
	lc       logger.LoggingClient
}

func NewCategoryTemplateApp(ctx context.Context, dic *di.Container) interfaces.CategoryApp {
	lc := container.LoggingClientFrom(dic.Get)
	dbClient := resourceContainer.DBClientFrom(dic.Get)

	return &categoryApp{
		dic:      dic,
		dbClient: dbClient,
		lc:       lc,
	}
}

func (m *categoryApp) CategoryTemplateSearch(ctx context.Context, req dtos.CategoryTemplateRequest) ([]dtos.CategoryTemplateResponse, uint32, error) {
	offset, limit := req.BaseSearchConditionQuery.GetPage()

	categoryTemplates, total, err := m.dbClient.CategoryTemplateSearch(offset, limit, req)
	if err != nil {
		m.lc.Errorf("categoryTemplates Search err %v", err)
		return []dtos.CategoryTemplateResponse{}, 0, err
	}

	libs := make([]dtos.CategoryTemplateResponse, len(categoryTemplates))
	for i, categoryTemplate := range categoryTemplates {
		libs[i] = dtos.CategoryTemplateResponseFromModel(categoryTemplate)
	}
	return libs, total, nil
}

func (m *categoryApp) Sync(ctx context.Context, versionName string) (int64, error) {
	filePath := versionName + "/category_template.json"
	cosApp := resourceContainer.CosAppNameFrom(m.dic.Get)
	bs, err := cosApp.Get(filePath)
	if err != nil {
		m.lc.Errorf(err.Error())
		return 0, err
	}
	var cosCategoryTemplateResp []dtos.CosCategoryTemplateResponse
	err = json.Unmarshal(bs, &cosCategoryTemplateResp)
	if err != nil {
		m.lc.Errorf(err.Error())
		return 0, err
	}

	baseQuery := dtos.BaseSearchConditionQuery{
		IsAll: true,
	}
	dbreq := dtos.CategoryTemplateRequest{BaseSearchConditionQuery: baseQuery}
	categoryTemplateResponse, _, err := m.CategoryTemplateSearch(ctx, dbreq)
	if err != nil {
		return 0, err
	}

	upsertCategoryTemplate := make([]models.CategoryTemplate, 0)
	for _, cosCategoryTemplate := range cosCategoryTemplateResp {
		var find bool
		for _, localTemplateResponse := range categoryTemplateResponse {
			if cosCategoryTemplate.CategoryKey == localTemplateResponse.CategoryKey {
				upsertCategoryTemplate = append(upsertCategoryTemplate, models.CategoryTemplate{
					Id:           localTemplateResponse.Id,
					CategoryName: cosCategoryTemplate.CategoryName,
					CategoryKey:  cosCategoryTemplate.CategoryKey,
					Scene:        cosCategoryTemplate.Scene,
				})
				find = true
				break
			}
		}
		if !find {
			upsertCategoryTemplate = append(upsertCategoryTemplate, models.CategoryTemplate{
				Timestamps: models.Timestamps{
					Created: time.Now().Unix(),
				},
				Id:           utils.GenUUID(),
				CategoryName: cosCategoryTemplate.CategoryName,
				CategoryKey:  cosCategoryTemplate.CategoryKey,
				Scene:        cosCategoryTemplate.Scene,
			})
		}
	}
	rows, err := m.dbClient.BatchUpsertCategoryTemplate(upsertCategoryTemplate)
	m.lc.Infof("upsert category template rows %+v", rows)
	if err != nil {
		return 0, err
	}
	return rows, nil
}
