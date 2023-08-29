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

package dataresource

import (
	"context"
	"database/sql"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/mitchellh/mapstructure"
	"github.com/winc-link/hummingbird/internal/dtos"
	resourceContainer "github.com/winc-link/hummingbird/internal/hummingbird/core/container"
	interfaces "github.com/winc-link/hummingbird/internal/hummingbird/core/interface"
	"github.com/winc-link/hummingbird/internal/models"
	"github.com/winc-link/hummingbird/internal/pkg/constants"
	"github.com/winc-link/hummingbird/internal/pkg/container"
	"github.com/winc-link/hummingbird/internal/pkg/di"
	"github.com/winc-link/hummingbird/internal/pkg/errort"
	"github.com/winc-link/hummingbird/internal/pkg/logger"
	"net/url"
	"time"

	_ "github.com/influxdata/influxdb1-client/v2"
	client "github.com/influxdata/influxdb1-client/v2"
	//_ "github.com/taosdata/driver-go/v2/taosSql"
)

type dataResourceApp struct {
	dic      *di.Container
	dbClient interfaces.DBClient
	lc       logger.LoggingClient
}

func NewDataResourceApp(ctx context.Context, dic *di.Container) interfaces.DataResourceApp {
	lc := container.LoggingClientFrom(dic.Get)
	dbClient := resourceContainer.DBClientFrom(dic.Get)

	app := &dataResourceApp{
		dic:      dic,
		dbClient: dbClient,
		lc:       lc,
	}
	return app
}

func (p dataResourceApp) AddDataResource(ctx context.Context, req dtos.AddDataResourceReq) (string, error) {
	var insertDataResource models.DataResource
	insertDataResource.Name = req.Name
	insertDataResource.Type = constants.DataResourceType(req.Type)
	insertDataResource.Option = req.Option
	insertDataResource.Option["sendSingle"] = true
	id, err := p.dbClient.AddDataResource(insertDataResource)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (p dataResourceApp) DataResourceById(ctx context.Context, id string) (models.DataResource, error) {
	if id == "" {
		return models.DataResource{}, errort.NewCommonEdgeX(errort.DefaultReqParamsError, "req id is required", nil)

	}
	dataResource, edgeXErr := p.dbClient.DataResourceById(id)
	if edgeXErr != nil {
		return models.DataResource{}, edgeXErr
	}
	return dataResource, nil
}

func (p dataResourceApp) UpdateDataResource(ctx context.Context, req dtos.UpdateDataResource) error {
	if req.Id == "" {
		return errort.NewCommonEdgeX(errort.DefaultReqParamsError, "update req id is required", nil)

	}
	dataResource, edgeXErr := p.dbClient.DataResourceById(req.Id)
	if edgeXErr != nil {
		return edgeXErr
	}
	ruleEngines, _, err := p.dbClient.RuleEngineSearch(0, -1, dtos.RuleEngineSearchQueryRequest{
		Status: string(constants.RuleEngineStart),
	})
	if err != nil {
		return err
	}

	for _, engine := range ruleEngines {
		if engine.DataResourceId == req.Id {
			return errort.NewCommonErr(errort.RuleEngineIsStartingNotAllowUpdate, fmt.Errorf("please stop this rule engine (%s) before editing it", req.Id))
		}
	}

	dtos.ReplaceDataResourceModelFields(&dataResource, req)
	edgeXErr = p.dbClient.UpdateDataResource(dataResource)
	if edgeXErr != nil {
		return edgeXErr
	}
	return nil
}

func (p dataResourceApp) DelDataResourceById(ctx context.Context, id string) error {
	ruleEngines, _, err := p.dbClient.RuleEngineSearch(0, -1, dtos.RuleEngineSearchQueryRequest{
		Status: string(constants.RuleEngineStart),
	})
	if err != nil {
		return err
	}

	for _, engine := range ruleEngines {
		if engine.DataResourceId == id {
			return errort.NewCommonErr(errort.RuleEngineIsStartingNotAllowUpdate, fmt.Errorf("please stop this rule engine (%s) before editing it", id))
		}
	}

	err = p.dbClient.DelDataResource(id)
	if err != nil {
		return err
	}
	return nil
}

func (p dataResourceApp) DataResourceSearch(ctx context.Context, req dtos.DataResourceSearchQueryRequest) ([]models.DataResource, uint32, error) {
	offset, limit := req.BaseSearchConditionQuery.GetPage()
	resp, total, err := p.dbClient.SearchDataResource(offset, limit, req)
	if err != nil {
		return []models.DataResource{}, 0, err
	}
	return resp, total, nil
}

func (p dataResourceApp) DataResourceType(ctx context.Context) []constants.DataResourceType {
	return constants.DataResources
}

func (p dataResourceApp) DataResourceHealth(ctx context.Context, resourceId string) error {
	dataResource, err := p.dbClient.DataResourceById(resourceId)
	if err != nil {
		return err
	}
	//return p.dbClient.UpdateDataResourceHealth(dataResource.Id, true)
	switch dataResource.Type {
	case constants.HttpResource:
		err = p.checkHttpResourceHealth(dataResource)
	case constants.MQTTResource:
		err = p.checkMQTTResourceHealth(dataResource)
	case constants.KafkaResource:
		err = p.checkKafkaResourceHealth(dataResource)
	case constants.InfluxDBResource:
		err = p.checkInfluxDBResourceHealth(dataResource)
	case constants.TDengineResource:
		//err = p.checkTdengineResourceHealth(dataResource)
	default:
		return errort.NewCommonErr(errort.DefaultReqParamsError, fmt.Errorf("resource  type not much"))
	}
	if err != nil {
		return err
	}
	return p.dbClient.UpdateDataResourceHealth(dataResource.Id, true)
}

func (p dataResourceApp) checkHttpResourceHealth(resource models.DataResource) error {
	urlAddr := resource.Option["url"].(string)
	_, err := url.Parse(urlAddr)
	if err != nil {
		return err
	}
	return nil
}

func (p dataResourceApp) checkMQTTResourceHealth(resource models.DataResource) error {
	var (
		server, topic                string
		clientId, username, password string
		//certificationPath,
		//privateKeyPath, rootCaPath, insecureSkipVerify, retained, compression, connectionSelector string
	)
	server = resource.Option["server"].(string)
	topic = resource.Option["topic"].(string)
	clientId = resource.Option["clientId"].(string)
	//protocolVersion = resource.Option["protocolVersion"]
	//qos = resource.Option["qos"].(int)
	username = resource.Option["username"].(string)
	password = resource.Option["password"].(string)
	//certificationPath = resource.Option["certificationPath"]
	//privateKeyPath = resource.Option["privateKeyPath"]
	//rootCaPath = resource.Option["rootCaPath"]
	//insecureSkipVerify = resource.Option["insecureSkipVerify"]
	//retained = resource.Option["retained"]
	//compression = resource.Option["compression"]
	//connectionSelector = resource.Option["connectionSelector"]

	if server == "" || topic == "" || clientId == "" || username == "" || password == "" {

	}
	opts := mqtt.NewClientOptions()
	opts.AddBroker(server)
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.SetClientID(clientId)

	client := mqtt.NewClient(opts)
	token := client.Connect()
	// 如果连接失败，则终止程序
	if token.WaitTimeout(3*time.Second) && token.Error() != nil {
		return token.Error()
	}
	defer client.Disconnect(250)
	return nil
}

func (p dataResourceApp) checkKafkaResourceHealth(resource models.DataResource) error {
	return nil
}

func (p dataResourceApp) checkInfluxDBResourceHealth(resource models.DataResource) error {
	type influxSink struct {
		addr         string
		username     string
		password     string
		measurement  string
		databaseName string
		tagKey       string
		tagValue     string
		fields       string
		cli          client.Client
		fieldMap     map[string]interface{}
		hasTransform bool
	}
	var m influxSink
	if i, ok := resource.Option["addr"]; ok {
		if i, ok := i.(string); ok {
			m.addr = i
		}
	}
	if i, ok := resource.Option["username"]; ok {
		if i, ok := i.(string); ok {
			m.username = i
		}
	}
	if i, ok := resource.Option["password"]; ok {
		if i, ok := i.(string); ok {
			m.password = i
		}
	}
	if i, ok := resource.Option["measurement"]; ok {
		if i, ok := i.(string); ok {
			m.measurement = i
		}
	}
	if i, ok := resource.Option["databasename"]; ok {
		if i, ok := i.(string); ok {
			m.databaseName = i
		}
	}
	if i, ok := resource.Option["tagkey"]; ok {
		if i, ok := i.(string); ok {
			m.tagKey = i
		}
	}
	if i, ok := resource.Option["tagvalue"]; ok {
		if i, ok := i.(string); ok {
			m.tagValue = i
		}
	}
	if i, ok := resource.Option["fields"]; ok {
		if i, ok := i.(string); ok {
			m.fields = i
		}
	}
	if i, ok := resource.Option["dataTemplate"]; ok {
		if i, ok := i.(string); ok && i != "" {
			m.hasTransform = true
		}
	}

	_, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     m.addr,
		Username: m.username,
		Password: m.password,
	})
	if err != nil {
		return err
	}
	return nil
}

func (p dataResourceApp) checkTdengineResourceHealth(resource models.DataResource) error {
	type taosConfig struct {
		ProvideTs      bool     `json:"provideTs"`
		Port           int      `json:"port"`
		Ip             string   `json:"ip"` // To be deprecated
		Host           string   `json:"host"`
		User           string   `json:"user"`
		Password       string   `json:"password"`
		Database       string   `json:"database"`
		Table          string   `json:"table"`
		TsFieldName    string   `json:"tsFieldName"`
		Fields         []string `json:"fields"`
		STable         string   `json:"sTable"`
		TagFields      []string `json:"tagFields"`
		DataTemplate   string   `json:"dataTemplate"`
		TableDataField string   `json:"tableDataField"`
	}
	cfg := &taosConfig{
		User:     "root",
		Password: "taosdata",
	}
	err := MapToStruct(resource.Option, cfg)
	if err != nil {
		return fmt.Errorf("read properties %v fail with error: %v", resource.Option, err)
	}
	if cfg.Ip != "" {
		fmt.Errorf("Deprecated: Tdengine sink ip property is deprecated, use host instead.")
		if cfg.Host == "" {
			cfg.Host = cfg.Ip
		}
	}
	if cfg.Host == "" {
		cfg.Host = "localhost"
	}
	if cfg.User == "" {
		return fmt.Errorf("propert user is required.")
	}
	if cfg.Password == "" {
		return fmt.Errorf("propert password is required.")
	}
	if cfg.Database == "" {
		return fmt.Errorf("property database is required")
	}
	if cfg.Table == "" {
		return fmt.Errorf("property table is required")
	}
	if cfg.TsFieldName == "" {
		return fmt.Errorf("property TsFieldName is required")
	}
	if cfg.STable != "" && len(cfg.TagFields) == 0 {
		return fmt.Errorf("property tagFields is required when sTable is set")
	}
	url := fmt.Sprintf(`%s:%s@tcp(%s:%d)/%s`, cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
	//m.conf = cfg
	_, err = sql.Open("taosSql", url)
	if err != nil {
		return err
	}
	return nil
}

func MapToStruct(input, output interface{}) error {
	config := &mapstructure.DecoderConfig{
		TagName: "json",
		Result:  output,
	}
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	return decoder.Decode(input)
}
