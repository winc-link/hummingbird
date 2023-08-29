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
package unittemplate

import (
	"context"
	"encoding/json"
	"github.com/winc-link/hummingbird/internal/dtos"
	resourceContainer "github.com/winc-link/hummingbird/internal/hummingbird/core/container"
	interfaces "github.com/winc-link/hummingbird/internal/hummingbird/core/interface"
	"github.com/winc-link/hummingbird/internal/models"
	"github.com/winc-link/hummingbird/internal/pkg/container"
	"github.com/winc-link/hummingbird/internal/pkg/di"
	"github.com/winc-link/hummingbird/internal/pkg/logger"
	"github.com/winc-link/hummingbird/internal/pkg/utils"
)

type unitApp struct {
	dic      *di.Container
	dbClient interfaces.DBClient
	lc       logger.LoggingClient
}

func NewUnitTemplateApp(ctx context.Context, dic *di.Container) interfaces.UnitApp {
	lc := container.LoggingClientFrom(dic.Get)
	dbClient := resourceContainer.DBClientFrom(dic.Get)

	return &unitApp{
		dic:      dic,
		dbClient: dbClient,
		lc:       lc,
	}
}

func (m *unitApp) UnitTemplateSearch(ctx context.Context, req dtos.UnitRequest) ([]dtos.UnitResponse, uint32, error) {
	offset, limit := req.BaseSearchConditionQuery.GetPage()

	unitTemplates, total, err := m.dbClient.UnitSearch(offset, limit, req)
	if err != nil {
		m.lc.Errorf("unit Templates Search err %v", err)
		return []dtos.UnitResponse{}, 0, err
	}

	libs := make([]dtos.UnitResponse, len(unitTemplates))
	for i, unitTemplate := range unitTemplates {
		libs[i] = dtos.UnitTemplateResponseFromModel(unitTemplate)
	}
	return libs, total, nil
}

func (m *unitApp) Sync(ctx context.Context, versionName string) (int64, error) {
	filePath := versionName + "/unit_template.json"
	cosApp := resourceContainer.CosAppNameFrom(m.dic.Get)
	bs, err := cosApp.Get(filePath)
	if err != nil {
		m.lc.Errorf(err.Error())
		return 0, err
	}
	var cosUnitTemplateResp []dtos.CosUnitTemplateResponse
	err = json.Unmarshal(bs, &cosUnitTemplateResp)
	if err != nil {
		m.lc.Errorf(err.Error())
		return 0, err
	}

	baseQuery := dtos.BaseSearchConditionQuery{
		IsAll: true,
	}
	dbreq := dtos.UnitRequest{BaseSearchConditionQuery: baseQuery}
	unitTemplateResponse, _, err := m.UnitTemplateSearch(ctx, dbreq)
	if err != nil {
		return 0, err
	}

	upsertUnitTemplate := make([]models.Unit, 0)
	for _, cosUnitTemplate := range cosUnitTemplateResp {
		var find bool
		for _, localTemplateResponse := range unitTemplateResponse {
			if cosUnitTemplate.UnitName == localTemplateResponse.UnitName {
				upsertUnitTemplate = append(upsertUnitTemplate, models.Unit{
					Id:       localTemplateResponse.Id,
					UnitName: cosUnitTemplate.UnitName,
					Symbol:   cosUnitTemplate.Symbol,
				})
				find = true
				break
			}
		}
		if !find {
			upsertUnitTemplate = append(upsertUnitTemplate, models.Unit{
				Id:       utils.GenUUID(),
				UnitName: cosUnitTemplate.UnitName,
				Symbol:   cosUnitTemplate.Symbol,
			})
		}
	}
	rows, err := m.dbClient.BatchUpsertUnitTemplate(upsertUnitTemplate)
	if err != nil {
		return 0, err
	}
	return rows, nil
}
