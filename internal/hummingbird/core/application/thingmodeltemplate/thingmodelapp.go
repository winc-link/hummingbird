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
package thingmodeltemplate

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
	"time"
)

type thingModelTemplate struct {
	dic      *di.Container
	dbClient interfaces.DBClient
	lc       logger.LoggingClient
}

func (m thingModelTemplate) ThingModelTemplateSearch(ctx context.Context, req dtos.ThingModelTemplateRequest) ([]dtos.ThingModelTemplateResponse, uint32, error) {
	offset, limit := req.BaseSearchConditionQuery.GetPage()

	thingModelTemplates, total, err := m.dbClient.ThingModelTemplateSearch(offset, limit, req)
	if err != nil {
		m.lc.Errorf("thingModelTemplate Search err %v", err)
		return []dtos.ThingModelTemplateResponse{}, 0, err
	}

	libs := make([]dtos.ThingModelTemplateResponse, len(thingModelTemplates))
	for i, thingModelTemplate := range thingModelTemplates {
		libs[i] = dtos.ThingModelTemplateResponseFromModel(thingModelTemplate)
	}
	return libs, total, nil
}

func (m thingModelTemplate) ThingModelTemplateByCategoryKey(ctx context.Context, categoryKey string) (dtos.ThingModelTemplateResponse, error) {
	thingModelTemplate, err := m.dbClient.ThingModelTemplateByCategoryKey(categoryKey)
	if err != nil {
		m.lc.Errorf("thingModelTemplate Search err %v", err)
		return dtos.ThingModelTemplateResponse{}, err
	}
	var libs dtos.ThingModelTemplateResponse
	libs = dtos.ThingModelTemplateResponseFromModel(thingModelTemplate)
	return libs, nil
}

func (m thingModelTemplate) Sync(ctx context.Context, versionName string) (int64, error) {
	filePath := versionName + "/thing_model_template.json"
	cosApp := resourceContainer.CosAppNameFrom(m.dic.Get)
	bs, err := cosApp.Get(filePath)
	if err != nil {
		m.lc.Errorf(err.Error())
		return 0, err
	}
	var cosThingModelTemplateResp []dtos.CosThingModelTemplateResponse
	err = json.Unmarshal(bs, &cosThingModelTemplateResp)
	if err != nil {
		m.lc.Errorf(err.Error())
		return 0, err
	}

	baseQuery := dtos.BaseSearchConditionQuery{
		IsAll: true,
	}
	dbreq := dtos.ThingModelTemplateRequest{BaseSearchConditionQuery: baseQuery}
	thingModelTemplateResponse, _, err := m.ThingModelTemplateSearch(ctx, dbreq)
	if err != nil {
		return 0, err
	}

	upsertThingModelTemplate := make([]models.ThingModelTemplate, 0)
	for _, cosThingModelTemplate := range cosThingModelTemplateResp {
		var find bool
		for _, localThingModelResponse := range thingModelTemplateResponse {
			if cosThingModelTemplate.CategoryKey == localThingModelResponse.CategoryKey {
				upsertThingModelTemplate = append(upsertThingModelTemplate, models.ThingModelTemplate{
					Id:             localThingModelResponse.Id,
					CategoryName:   cosThingModelTemplate.CategoryName,
					CategoryKey:    cosThingModelTemplate.CategoryKey,
					ThingModelJSON: cosThingModelTemplate.ThingModelJSON,
				})
				find = true
				break
			}
		}
		if !find {
			upsertThingModelTemplate = append(upsertThingModelTemplate, models.ThingModelTemplate{
				Timestamps: models.Timestamps{
					Created: time.Now().Unix(),
				},
				Id:             utils.GenUUID(),
				CategoryName:   cosThingModelTemplate.CategoryName,
				CategoryKey:    cosThingModelTemplate.CategoryKey,
				ThingModelJSON: cosThingModelTemplate.ThingModelJSON,
			})
		}
	}
	rows, err := m.dbClient.BatchUpsertThingModelTemplate(upsertThingModelTemplate)
	m.lc.Infof("upsert thingModelTemplate rows %+v", rows)
	if err != nil {
		return 0, err
	}
	var modelProperty []models.Properties
	var modelEvent []models.Events
	var modelAction []models.Actions

	propertyTemp := map[string]struct{}{}
	eventTemp := map[string]struct{}{}
	actionTemp := map[string]struct{}{}

	for _, cosThingModelTemplate := range cosThingModelTemplateResp {
		p, e, a := dtos.GetModelPropertyEventActionByThingModelTemplate(cosThingModelTemplate.ThingModelJSON)
		for _, properties := range p {
			if _, ok := propertyTemp[properties.Code]; !ok {
				properties.Id = utils.GenUUID()
				properties.System = true
				modelProperty = append(modelProperty, properties)
				propertyTemp[properties.Code] = struct{}{}
			}
		}
		for _, event := range e {
			if _, ok := eventTemp[event.Code]; !ok {
				event.Id = utils.GenUUID()
				event.System = true
				modelEvent = append(modelEvent, event)
				eventTemp[event.Code] = struct{}{}
			}
		}
		for _, action := range a {
			if _, ok := actionTemp[action.Code]; !ok {
				action.Id = utils.GenUUID()
				action.System = true
				modelAction = append(modelAction, action)
				actionTemp[action.Code] = struct{}{}
			}
		}
	}
	m.dbClient.BatchDeleteSystemProperties()
	m.dbClient.BatchDeleteSystemActions()
	m.dbClient.BatchDeleteSystemEvents()

	m.dbClient.BatchInsertSystemProperties(modelProperty)
	m.dbClient.BatchInsertSystemActions(modelAction)
	m.dbClient.BatchInsertSystemEvents(modelEvent)

	return rows, nil
}

func NewThingModelTemplateApp(ctx context.Context, dic *di.Container) interfaces.ThingModelTemplateApp {
	lc := container.LoggingClientFrom(dic.Get)
	dbClient := resourceContainer.DBClientFrom(dic.Get)

	return &thingModelTemplate{
		dic:      dic,
		dbClient: dbClient,
		lc:       lc,
	}
}
