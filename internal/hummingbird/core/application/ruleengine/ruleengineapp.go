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

package ruleengine

import (
	"context"
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
)

type ruleEngineApp struct {
	dic      *di.Container
	dbClient interfaces.DBClient
	lc       logger.LoggingClient
}

func (p ruleEngineApp) AddRuleEngine(ctx context.Context, req dtos.RuleEngineRequest) (string, error) {
	dataResource, err := p.dbClient.DataResourceById(req.DataResourceId)
	if err != nil {
		return "", err
	}
	randomId := utils.RandomNum()
	ekuiperApp := resourceContainer.EkuiperAppFrom(p.dic.Get)

	sql := req.BuildEkuiperSql()
	var actions []dtos.Actions
	switch dataResource.Type {
	case constants.HttpResource:
		actions = append(actions, dtos.Actions{
			Rest: dataResource.Option,
		})
	case constants.MQTTResource:
		actions = append(actions, dtos.Actions{
			MQTT: dataResource.Option,
		})
	case constants.KafkaResource:
		actions = append(actions, dtos.Actions{
			Kafka: dataResource.Option,
		})
	case constants.InfluxDBResource:
		actions = append(actions, dtos.Actions{
			Influx: dataResource.Option,
		})
	case constants.TDengineResource:
		actions = append(actions, dtos.Actions{
			Tdengine: dataResource.Option,
		})
	default:
		return "", errort.NewCommonErr(errort.DefaultReqParamsError, fmt.Errorf("rule engine action not much"))
	}
	if err = ekuiperApp.CreateRule(ctx, actions, randomId, sql); err != nil {
		return "", err
	}

	var insertRuleEngine models.RuleEngine
	insertRuleEngine.Name = req.Name
	insertRuleEngine.Id = randomId
	insertRuleEngine.Description = req.Description
	insertRuleEngine.Filter = models.Filter(req.Filter)
	insertRuleEngine.DataResourceId = req.DataResourceId
	insertRuleEngine.Status = constants.RuleEngineStop
	id, err := p.dbClient.AddRuleEngine(insertRuleEngine)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (p ruleEngineApp) UpdateRuleEngine(ctx context.Context, req dtos.RuleEngineUpdateRequest) error {
	dataResource, err := p.dbClient.DataResourceById(*req.DataResourceId)
	if err != nil {
		return err
	}
	ruleEngine, err := p.dbClient.RuleEngineById(req.Id)
	if err != nil {
		return err
	}
	sql := req.BuildEkuiperSql()
	var actions []dtos.Actions
	switch dataResource.Type {
	case constants.HttpResource:
		actions = append(actions, dtos.Actions{
			Rest: dataResource.Option,
		})
	case constants.MQTTResource:
		actions = append(actions, dtos.Actions{
			MQTT: dataResource.Option,
		})
	case constants.KafkaResource:
		actions = append(actions, dtos.Actions{
			Kafka: dataResource.Option,
		})
	case constants.InfluxDBResource:
		actions = append(actions, dtos.Actions{
			Influx: dataResource.Option,
		})
	case constants.TDengineResource:
		actions = append(actions, dtos.Actions{
			Tdengine: dataResource.Option,
		})
	default:
		return errort.NewCommonErr(errort.DefaultReqParamsError, fmt.Errorf("rule engine action not much"))
	}
	ekuiperApp := resourceContainer.EkuiperAppFrom(p.dic.Get)
	if err = ekuiperApp.UpdateRule(ctx, actions, req.Id, sql); err != nil {
		return err
	}
	dtos.ReplaceRuleEngineModelFields(&ruleEngine, req)
	err = p.dbClient.UpdateRuleEngine(ruleEngine)
	if err != nil {
		return err
	}
	return nil
}

func (p ruleEngineApp) UpdateRuleEngineField(ctx context.Context, req dtos.RuleEngineFieldUpdateRequest) error {
	//TODO implement me
	panic("implement me")
}

func (p ruleEngineApp) RuleEngineById(ctx context.Context, id string) (dtos.RuleEngineResponse, error) {
	ruleEngine, err := p.dbClient.RuleEngineById(id)
	var ruleEngineResponse dtos.RuleEngineResponse
	if err != nil {
		return ruleEngineResponse, err
	}
	ruleEngineResponse.Id = ruleEngine.Id
	ruleEngineResponse.Name = ruleEngine.Name
	ruleEngineResponse.Description = ruleEngine.Description
	ruleEngineResponse.Created = ruleEngine.Created
	ruleEngineResponse.Filter = dtos.Filter(ruleEngine.Filter)
	ruleEngineResponse.DataResourceId = ruleEngine.DataResourceId
	ruleEngineResponse.DataResource = dtos.DataResourceInfo{
		Name:   ruleEngine.DataResource.Name,
		Type:   string(ruleEngine.DataResource.Type),
		Option: ruleEngine.DataResource.Option,
	}
	return ruleEngineResponse, nil
}

func (p ruleEngineApp) RuleEngineSearch(ctx context.Context, req dtos.RuleEngineSearchQueryRequest) ([]dtos.RuleEngineSearchQueryResponse, uint32, error) {
	offset, limit := req.BaseSearchConditionQuery.GetPage()
	resp, total, err := p.dbClient.RuleEngineSearch(offset, limit, req)
	if err != nil {
		return []dtos.RuleEngineSearchQueryResponse{}, 0, err
	}
	ruleEngines := make([]dtos.RuleEngineSearchQueryResponse, len(resp))
	for i, p := range resp {
		ruleEngines[i] = dtos.RuleEngineSearchQueryResponseFromModel(p)
	}
	return ruleEngines, total, nil
}

func (p ruleEngineApp) RuleEngineDelete(ctx context.Context, id string) error {
	_, err := p.dbClient.RuleEngineById(id)
	if err != nil {
		return err
	}
	ekuiperApp := resourceContainer.EkuiperAppFrom(p.dic.Get)
	err = ekuiperApp.DeleteRule(ctx, id)
	if err != nil {
		return err
	}
	return p.dbClient.DeleteRuleEngineById(id)
}

func (p ruleEngineApp) RuleEngineStop(ctx context.Context, id string) error {
	_, err := p.dbClient.RuleEngineById(id)
	if err != nil {
		return err
	}
	//if alertRule.EkuiperRule() {
	ekuiperApp := resourceContainer.EkuiperAppFrom(p.dic.Get)
	err = ekuiperApp.StopRule(ctx, id)
	if err != nil {
		return err
	}
	return p.dbClient.RuleEngineStop(id)
}

func (p ruleEngineApp) RuleEngineStart(ctx context.Context, id string) error {
	ruleEngine, err := p.dbClient.RuleEngineById(id)
	if err != nil {
		return err
	}
	dataResource, err := p.dbClient.DataResourceById(ruleEngine.DataResourceId)
	if err != nil {
		return err
	}
	if dataResource.Health != true {
		return errort.NewCommonErr(errort.InvalidSource, fmt.Errorf("invalid resource configuration, please check the resource configuration resource id (%s)", dataResource.Id))
	}

	ekuiperApp := resourceContainer.EkuiperAppFrom(p.dic.Get)
	err = ekuiperApp.StartRule(ctx, id)
	if err != nil {
		return err
	}
	//}
	return p.dbClient.RuleEngineStart(id)
}

func (p ruleEngineApp) RuleEngineStatus(ctx context.Context, id string) (map[string]interface{}, error) {
	response := make(map[string]interface{}, 0)
	_, err := p.dbClient.RuleEngineById(id)
	if err != nil {
		return response, err
	}
	ekuiperApp := resourceContainer.EkuiperAppFrom(p.dic.Get)
	return ekuiperApp.GetRuleStats(ctx, id)
}

func NewRuleEngineApp(ctx context.Context, dic *di.Container) interfaces.RuleEngineApp {
	lc := container.LoggingClientFrom(dic.Get)
	dbClient := resourceContainer.DBClientFrom(dic.Get)

	app := &ruleEngineApp{
		dic:      dic,
		dbClient: dbClient,
		lc:       lc,
	}
	go app.monitor()
	return app
}
