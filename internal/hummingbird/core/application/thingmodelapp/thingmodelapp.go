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
package thingmodelapp

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/winc-link/hummingbird/internal/dtos"
	resourceContainer "github.com/winc-link/hummingbird/internal/hummingbird/core/container"
	interfaces "github.com/winc-link/hummingbird/internal/hummingbird/core/interface"
	"github.com/winc-link/hummingbird/internal/models"
	"github.com/winc-link/hummingbird/internal/pkg/constants"
	"github.com/winc-link/hummingbird/internal/pkg/container"
	"github.com/winc-link/hummingbird/internal/pkg/di"
	"github.com/winc-link/hummingbird/internal/pkg/errort"
	"github.com/winc-link/hummingbird/internal/pkg/logger"
	"github.com/winc-link/hummingbird/internal/pkg/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
	"time"
)

const (
	Property = "property"
	Event    = "event"
	Action   = "action"
)

type thingModelApp struct {
	dic      *di.Container
	dbClient interfaces.DBClient
	lc       logger.LoggingClient
}

func (t thingModelApp) AddThingModel(ctx context.Context, req dtos.ThingModelAddOrUpdateReq) (string, error) {

	product, err := t.dbClient.ProductById(req.ProductId)
	if err != nil {
		return "", err
	}
	if product.Status == constants.ProductRelease {
		//产品已发布，不能修改物模型
		return "", errort.NewCommonEdgeX(errort.ProductRelease, "Please cancel publishing the product first before proceeding with the operation", nil)
	}

	if err = validatorReq(req, product); err != nil {
		return "", errort.NewCommonEdgeX(errort.DefaultReqParamsError, "param valida error", err)
	}

	switch req.ThingModelType {
	case Property:
		var property models.Properties
		property.ProductId = req.ProductId
		property.Name = req.Name
		property.Code = req.Code
		property.Description = req.Description
		if req.Property != nil {
			typeSpec, _ := json.Marshal(req.Property.TypeSpec)
			property.AccessMode = req.Property.AccessModel
			property.Require = req.Property.Require
			property.TypeSpec.Type = req.Property.DataType
			property.TypeSpec.Specs = string(typeSpec)
		}
		property.Tag = req.Tag
		err = resourceContainer.DataDBClientFrom(t.dic.Get).AddDatabaseField(ctx, req.ProductId, req.Property.DataType, req.Code, req.Name)
		if err != nil {
			return "", err
		}
		ds, err := t.dbClient.AddThingModelProperty(property)
		if err != nil {
			return "", err
		}
		t.ProductUpdateCallback(req.ProductId)
		return ds.Id, nil
	case Event:
		var event models.Events
		event.ProductId = req.ProductId
		event.Name = req.Name
		event.Code = req.Code
		event.Description = req.Description
		if req.Event != nil {
			var inputOutput []models.InputOutput
			for _, outPutParam := range req.Event.OutPutParam {
				typeSpec, _ := json.Marshal(outPutParam.TypeSpec)
				inputOutput = append(inputOutput, models.InputOutput{
					Code: outPutParam.Code,
					Name: outPutParam.Name,
					TypeSpec: models.TypeSpec{
						Type:  outPutParam.DataType,
						Specs: string(typeSpec),
					},
				})
			}
			event.OutputParams = inputOutput
			event.EventType = req.Event.EventType
		}
		event.Tag = req.Tag
		err = resourceContainer.DataDBClientFrom(t.dic.Get).AddDatabaseField(ctx, req.ProductId, "", req.Code, req.Name)
		if err != nil {
			return "", err
		}
		ds, err := t.dbClient.AddThingModelEvent(event)
		if err != nil {
			return "", err
		}
		t.ProductUpdateCallback(req.ProductId)
		return ds.Id, nil
	case Action:
		var action models.Actions
		action.ProductId = req.ProductId
		action.Name = req.Name
		action.Code = req.Code
		action.Description = req.Description
		if req.Action != nil {
			action.CallType = req.Action.CallType
			var inputOutput []models.InputOutput
			for _, inPutParam := range req.Action.InPutParam {
				typeSpec, _ := json.Marshal(inPutParam.TypeSpec)
				inputOutput = append(inputOutput, models.InputOutput{
					Code: inPutParam.Code,
					Name: inPutParam.Name,
					TypeSpec: models.TypeSpec{
						Type:  inPutParam.DataType,
						Specs: string(typeSpec),
					},
				})
			}
			action.InputParams = inputOutput

			var outOutput []models.InputOutput
			for _, outPutParam := range req.Action.OutPutParam {
				typeSpec, _ := json.Marshal(outPutParam.TypeSpec)
				outOutput = append(outOutput, models.InputOutput{
					Code: outPutParam.Code,
					Name: outPutParam.Name,
					TypeSpec: models.TypeSpec{
						Type:  outPutParam.DataType,
						Specs: string(typeSpec),
					},
				})
			}
			action.OutputParams = outOutput
		}
		action.Tag = req.Tag
		err = resourceContainer.DataDBClientFrom(t.dic.Get).AddDatabaseField(ctx, req.ProductId, "", req.Code, req.Name)
		if err != nil {
			return "", err
		}
		ds, err := t.dbClient.AddThingModelAction(action)
		if err != nil {
			return "", err
		}
		t.ProductUpdateCallback(req.ProductId)
		return ds.Id, nil
	default:
		return "", errort.NewCommonEdgeX(errort.DefaultReqParamsError, "param valida error", fmt.Errorf("req params error"))
	}
}

func (t thingModelApp) UpdateThingModel(ctx context.Context, req dtos.ThingModelAddOrUpdateReq) error {
	product, err := t.dbClient.ProductById(req.ProductId)
	if err != nil {
		return err
	}
	if product.Status == constants.ProductRelease {
		//产品已发布，不能修改物模型
		return errort.NewCommonEdgeX(errort.ProductRelease, "Please cancel publishing the product first before proceeding with the operation", nil)
	}
	if err := validatorReq(req, product); err != nil {
		return errort.NewCommonEdgeX(errort.DefaultReqParamsError, "param valida error", err)
	}
	switch req.ThingModelType {
	case Property:
		var property models.Properties
		property.Id = req.Id
		property.ProductId = req.ProductId
		property.Name = req.Name
		property.Code = req.Code
		property.Description = req.Description
		if req.Property != nil {
			typeSpec, _ := json.Marshal(req.Property.TypeSpec)
			property.AccessMode = req.Property.AccessModel
			property.Require = req.Property.Require
			property.TypeSpec.Type = req.Property.DataType
			property.TypeSpec.Specs = string(typeSpec)
		}
		property.Tag = req.Tag

		err = resourceContainer.DataDBClientFrom(t.dic.Get).ModifyDatabaseField(ctx, req.ProductId, req.Property.DataType, req.Code, req.Name)
		if err != nil {
			return err
		}

		err = t.dbClient.UpdateThingModelProperty(property)
		if err != nil {
			return err
		}
		t.ProductUpdateCallback(req.ProductId)
		return nil
	case Event:
		var event models.Events
		event.Id = req.Id
		event.ProductId = req.ProductId
		event.Name = req.Name
		event.Code = req.Code
		event.Description = req.Description
		if req.Event != nil {
			var inputOutput []models.InputOutput
			for _, outPutParam := range req.Event.OutPutParam {
				typeSpec, _ := json.Marshal(outPutParam.TypeSpec)
				inputOutput = append(inputOutput, models.InputOutput{
					Code: outPutParam.Code,
					Name: outPutParam.Name,
					TypeSpec: models.TypeSpec{
						Type:  outPutParam.DataType,
						Specs: string(typeSpec),
					},
				})
			}
			event.OutputParams = inputOutput
			event.EventType = req.Event.EventType
		}
		event.Tag = req.Tag
		err = resourceContainer.DataDBClientFrom(t.dic.Get).ModifyDatabaseField(ctx, req.ProductId, "", req.Code, req.Name)
		if err != nil {
			return err
		}
		err = t.dbClient.UpdateThingModelEvent(event)
		if err != nil {
			return err
		}
		t.ProductUpdateCallback(req.ProductId)
		return nil
	case Action:
		var action models.Actions
		action.Id = req.Id
		action.ProductId = req.ProductId
		action.Name = req.Name
		action.Code = req.Code
		action.Description = req.Description
		if req.Action != nil {
			action.CallType = req.Action.CallType
			var inputOutput []models.InputOutput
			for _, inPutParam := range req.Action.InPutParam {
				typeSpec, _ := json.Marshal(inPutParam.TypeSpec)
				inputOutput = append(inputOutput, models.InputOutput{
					Code: inPutParam.Code,
					Name: inPutParam.Name,
					TypeSpec: models.TypeSpec{
						Type:  inPutParam.DataType,
						Specs: string(typeSpec),
					},
				})
			}
			action.InputParams = inputOutput

			var outOutput []models.InputOutput
			for _, outPutParam := range req.Action.OutPutParam {
				typeSpec, _ := json.Marshal(outPutParam.TypeSpec)
				outOutput = append(outOutput, models.InputOutput{
					Code: outPutParam.Code,
					Name: outPutParam.Name,
					TypeSpec: models.TypeSpec{
						Type:  outPutParam.DataType,
						Specs: string(typeSpec),
					},
				})
			}
			action.OutputParams = outOutput
		}
		action.Tag = req.Tag
		err = resourceContainer.DataDBClientFrom(t.dic.Get).ModifyDatabaseField(ctx, req.ProductId, "", req.Code, req.Name)
		if err != nil {
			return err
		}
		err = t.dbClient.UpdateThingModelAction(action)
		if err != nil {
			return err
		}
		t.ProductUpdateCallback(req.ProductId)
		return nil
	default:
		return errort.NewCommonEdgeX(errort.DefaultReqParamsError, "param valida error", fmt.Errorf("req params error"))
	}
}

func validatorReq(req dtos.ThingModelAddOrUpdateReq, product models.Product) error {
	if req.ProductId == "" || req.Name == "" || req.Code == "" {
		return errort.NewCommonEdgeX(errort.DefaultReqParamsError, "params error", nil)
	}
	switch req.ThingModelType {
	case "property":
		for _, property := range product.Properties {
			if strings.ToLower(property.Code) == strings.ToLower(req.Code) {
				return errort.NewCommonEdgeX(errort.ThingModelCodeExist, "code identifier already exists", nil)
			}
		}
	case "action":
		for _, action := range product.Actions {
			if strings.ToLower(action.Code) == strings.ToLower(req.Code) {
				return errort.NewCommonEdgeX(errort.ThingModelCodeExist, "code identifier already exists", nil)
			}
		}
	case "event":
		for _, event := range product.Events {
			if strings.ToLower(event.Code) == strings.ToLower(req.Code) {
				return errort.NewCommonEdgeX(errort.ThingModelCodeExist, "code identifier already exists", nil)
			}
		}
	default:
		return errort.NewCommonEdgeX(errort.DefaultReqParamsError, "params error", nil)
	}
	return nil
}

func NewThingModelApp(ctx context.Context, dic *di.Container) interfaces.ThingModelItf {
	lc := container.LoggingClientFrom(dic.Get)
	dbClient := resourceContainer.DBClientFrom(dic.Get)

	return &thingModelApp{
		dic:      dic,
		dbClient: dbClient,
		lc:       lc,
	}
}

func (t thingModelApp) ThingModelDelete(ctx context.Context, id string, thingModelType string) error {
	switch thingModelType {
	case Property:
		propertyInfo, err := t.dbClient.ThingModelPropertyById(id)
		if err != nil {
			return err
		}
		err = resourceContainer.DataDBClientFrom(t.dic.Get).DelDatabaseField(ctx, propertyInfo.ProductId, propertyInfo.Code)
		if err != nil {
			return err
		}
		err = t.dbClient.ThingModelDeleteProperty(id)
		if err != nil {
			return err
		}
		t.ProductUpdateCallback(propertyInfo.ProductId)
	case Event:
		eventInfo, err := t.dbClient.ThingModelEventById(id)
		if err != nil {
			return err
		}
		err = resourceContainer.DataDBClientFrom(t.dic.Get).DelDatabaseField(ctx, eventInfo.ProductId, eventInfo.Code)
		if err != nil {
			return err
		}
		err = t.dbClient.ThingModelDeleteEvent(id)
		if err != nil {
			return err
		}
		t.ProductUpdateCallback(eventInfo.ProductId)
	case Action:
		actionInfo, err := t.dbClient.ThingModelActionsById(id)
		if err != nil {
			return err
		}
		err = resourceContainer.DataDBClientFrom(t.dic.Get).DelDatabaseField(ctx, actionInfo.ProductId, actionInfo.Code)
		if err != nil {
			return err
		}
		err = t.dbClient.ThingModelDeleteAction(id)
		if err != nil {
			return err
		}
		t.ProductUpdateCallback(actionInfo.ProductId)
	default:
		return errort.NewCommonEdgeX(errort.DefaultReqParamsError, "param valida error", fmt.Errorf("req params error"))
	}
	return nil
}

func (t thingModelApp) ProductUpdateCallback(productId string) {
	go func() {
		productService := resourceContainer.ProductAppNameFrom(t.dic.Get)
		product, err := productService.ProductModelById(context.Background(), productId)
		if err != nil {
			return
		}
		productService.UpdateProductCallBack(product)
	}()

}

func (t thingModelApp) SystemThingModelSearch(ctx context.Context, req dtos.SystemThingModelSearchReq) (interface{}, error) {
	return t.dbClient.SystemThingModelSearch(req.ThingModelType, req.ModelName)
}

func (t thingModelApp) OpenApiQueryThingModel(ctx context.Context, productId string) (dtos.OpenApiQueryThingModel, error) {
	product, err := t.dbClient.ProductById(productId)
	if err != nil {
		return dtos.OpenApiQueryThingModel{}, err
	}

	var response dtos.OpenApiQueryThingModel

	for _, property := range product.Properties {
		response.Properties = append(response.Properties, dtos.OpenApiThingModelProperties{
			Id:          property.Id,
			Name:        property.Name,
			Code:        property.Code,
			AccessMode:  property.AccessMode,
			Require:     property.Require,
			TypeSpec:    property.TypeSpec,
			Description: property.Description,
		})
	}

	if len(response.Properties) == 0 {
		response.Properties = make([]dtos.OpenApiThingModelProperties, 0)
	}

	for _, event := range product.Events {
		response.Events = append(response.Events, dtos.OpenApiThingModelEvents{
			Id:           event.Id,
			EventType:    event.EventType,
			Name:         event.Name,
			Code:         event.Code,
			Description:  event.Description,
			Require:      event.Require,
			OutputParams: event.OutputParams,
		})
	}

	if len(response.Events) == 0 {
		response.Events = make([]dtos.OpenApiThingModelEvents, 0)
	}

	for _, action := range product.Actions {
		response.Services = append(response.Services, dtos.OpenApiThingModelServices{
			Id:           action.Id,
			Name:         action.Name,
			Code:         action.Code,
			Description:  action.Description,
			Require:      action.Require,
			CallType:     action.CallType,
			InputParams:  action.InputParams,
			OutputParams: action.OutputParams,
		})
	}

	if len(response.Services) == 0 {
		response.Services = make([]dtos.OpenApiThingModelServices, 0)
	}

	return response, nil

}

func (t thingModelApp) OpenApiAddThingModel(ctx context.Context, req dtos.OpenApiThingModelAddOrUpdateReq) error {

	_, err := t.dbClient.ProductById(req.ProductId)
	if err != nil {
		return err
	}
	var properties []models.Properties
	var events []models.Events
	var action []models.Actions
	for _, property := range req.Properties {
		propertyId := property.Id
		if propertyId == "" {
			propertyId = utils.RandomNum()
		}
		properties = append(properties, models.Properties{
			Id:          propertyId,
			ProductId:   req.ProductId,
			Name:        property.Name,
			Code:        property.Code,
			AccessMode:  property.AccessMode,
			Require:     property.Require,
			TypeSpec:    property.TypeSpec,
			Description: property.Description,
			Timestamps: models.Timestamps{
				Created: time.Now().UnixMilli(),
			},
		})
	}

	for _, event := range req.Events {
		eventId := event.Id
		if eventId == "" {
			eventId = utils.RandomNum()
		}
		events = append(events, models.Events{
			Id:           eventId,
			ProductId:    req.ProductId,
			Name:         event.Name,
			EventType:    event.EventType,
			Code:         event.Code,
			Description:  event.Description,
			Require:      event.Require,
			OutputParams: event.OutputParams,
			Timestamps: models.Timestamps{
				Created: time.Now().UnixMilli(),
			},
		})
	}

	for _, service := range req.Services {
		serviceId := service.Id
		if serviceId == "" {
			serviceId = utils.RandomNum()
		}
		action = append(action, models.Actions{
			Id:           serviceId,
			ProductId:    req.ProductId,
			Name:         service.Name,
			Code:         service.Code,
			Description:  service.Description,
			Require:      service.Require,
			CallType:     service.CallType,
			InputParams:  service.InputParams,
			OutputParams: service.OutputParams,
			Timestamps: models.Timestamps{
				Created: time.Now().UnixMilli(),
			},
		})
	}

	//t.dbClient.ThingModelActionsById()

	var shouldCallBack bool
	updateFunc := func(source interface{}, db *gorm.DB) error {
		tx := db.Session(&gorm.Session{FullSaveAssociations: true}).Clauses(clause.OnConflict{UpdateAll: true}).Save(source)
		return tx.Error
		//tx := db.Save(&source).Error
		//return tx
	}

	db := t.dbClient.GetDBInstance()
	//db.Begin()
	if len(properties) > 0 {
		shouldCallBack = true
		err := updateFunc(properties, db)
		if err != nil {
			//db.Rollback()
			return err
		}
	}

	if len(events) > 0 {
		shouldCallBack = true
		err := updateFunc(events, db)
		if err != nil {
			//db.Rollback()
			return err
		}
	}

	if len(action) > 0 {
		shouldCallBack = true
		err := updateFunc(action, db)
		if err != nil {
			//db.Rollback()
			return err
		}
	}

	//db.Commit()
	if shouldCallBack {
		t.ProductUpdateCallback(req.ProductId)
	}

	return nil
}

func (t thingModelApp) OpenApiDeleteThingModel(ctx context.Context, req dtos.OpenApiThingModelDeleteReq) error {
	product, err := t.dbClient.ProductById(req.ProductId)
	if err != nil {
		return err
	}

	var productPropertyIds []string
	var productEventIds []string
	var productService []string

	for _, property := range product.Properties {
		productPropertyIds = append(productPropertyIds, property.Id)
	}

	for _, event := range product.Events {
		productEventIds = append(productEventIds, event.Id)
	}

	for _, action := range product.Actions {
		productService = append(productService, action.Id)
	}

	for _, id := range req.PropertyIds {
		if utils.InStringSlice(id, productPropertyIds) {
			if err = t.dbClient.ThingModelDeleteProperty(id); err != nil {
				return err
			}
		}
	}

	for _, id := range req.EventIds {
		if utils.InStringSlice(id, productEventIds) {
			if err = t.dbClient.ThingModelDeleteEvent(id); err != nil {
				return err
			}
		}
	}

	for _, id := range req.ServiceIds {
		if utils.InStringSlice(id, productService) {
			if err = t.dbClient.ThingModelDeleteAction(id); err != nil {
				return err
			}
		}
	}

	return nil

}
