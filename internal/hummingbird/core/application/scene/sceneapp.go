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

package scene

import (
	"context"
	"errors"
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
	"strconv"
	"strings"
)

type sceneApp struct {
	dic      *di.Container
	dbClient interfaces.DBClient
	lc       logger.LoggingClient
}

func NewSceneApp(ctx context.Context, dic *di.Container) interfaces.SceneApp {
	lc := container.LoggingClientFrom(dic.Get)
	dbClient := resourceContainer.DBClientFrom(dic.Get)

	app := &sceneApp{
		dic:      dic,
		dbClient: dbClient,
		lc:       lc,
	}
	go app.monitor()
	return app
}

func (p sceneApp) AddScene(ctx context.Context, req dtos.SceneAddRequest) (string, error) {
	var scene models.Scene
	scene.Name = req.Name
	scene.Description = req.Description
	resp, err := p.dbClient.AddScene(scene)
	if err != nil {
		return "", err
	}
	return resp.Id, nil
}

func (p sceneApp) UpdateScene(ctx context.Context, req dtos.SceneUpdateRequest) error {
	if req.Id == "" {
		return errort.NewCommonEdgeX(errort.DefaultReqParamsError, "update req id is required", nil)
	}
	scene, edgeXErr := p.dbClient.SceneById(req.Id)
	if edgeXErr != nil {
		return edgeXErr
	}
	if len(req.Conditions) != 1 {
		return errort.NewCommonEdgeX(errort.DefaultReqParamsError, "conditions len not eq 1", nil)
	}
	switch req.Conditions[0].ConditionType {
	case "timer":
		if scene.Status == constants.SceneStart {
			return errort.NewCommonEdgeX(errort.SceneTimerIsStartingNotAllowUpdate, "Please stop this scheduled"+
				" tasks before editing it.", nil)
		}
	case "notify":
		if req.Conditions[0].Option == nil {
			return errort.NewCommonEdgeX(errort.DefaultReqParamsError, "condition option is null", nil)
		}
		actions, sql, err := p.buildEkuiperSqlAndAction(req)
		if err != nil {
			return err
		}
		p.lc.Infof("sql:", sql)

		ekuiperApp := resourceContainer.EkuiperAppFrom(p.dic.Get)

		exist, err := ekuiperApp.RuleExist(ctx, scene.Id)
		if err != nil {
			return err
		}
		if exist {
			err = ekuiperApp.UpdateRule(ctx, actions, scene.Id, sql)
			if err != nil {
				return err
			}
		} else {
			err = ekuiperApp.CreateRule(ctx, actions, scene.Id, sql)
			if err != nil {
				return err
			}
		}
	}
	dtos.ReplaceSceneModelFields(&scene, req)
	edgeXErr = p.dbClient.UpdateScene(scene)
	if edgeXErr != nil {
		return edgeXErr
	}
	return nil
}

func (p sceneApp) SceneById(ctx context.Context, sceneId string) (models.Scene, error) {
	return p.dbClient.SceneById(sceneId)
}

func (p sceneApp) SceneStartById(ctx context.Context, sceneId string) error {
	scene, err := p.dbClient.SceneById(sceneId)
	if err != nil {
		return err
	}
	if len(scene.Conditions) == 0 {
		return errort.NewCommonErr(errort.SceneRuleParamsError, fmt.Errorf("scene id(%s) conditionType param errror", scene.Id))
	}

	switch scene.Conditions[0].ConditionType {
	case "timer":
		tmpJob, errJob := scene.ToRuntimeJob()
		if errJob != nil {
			return errort.NewCommonEdgeX(errort.DefaultSystemError, errJob.Error(), errJob)
		}
		p.lc.Infof("tmpJob: %v", tmpJob)
		conJobApp := resourceContainer.ConJobAppNameFrom(p.dic.Get)
		err = conJobApp.AddJobToRunQueue(tmpJob)
		if err != nil {
			return err
		}
	case "notify":
		if scene.Conditions[0].Option == nil {
			return errort.NewCommonEdgeX(errort.DefaultReqParamsError, "condition option is null", nil)
		}
		if err = p.checkAlertRuleParam(ctx, scene, "start"); err != nil {
			return err
		}
		ekuiperApp := resourceContainer.EkuiperAppFrom(p.dic.Get)
		err = ekuiperApp.StartRule(ctx, sceneId)
		if err != nil {
			return err
		}
	default:
		return errort.NewCommonEdgeX(errort.DefaultReqParamsError, "condition Type value not much", nil)
	}
	return p.dbClient.SceneStart(sceneId)
}

func (p sceneApp) SceneStopById(ctx context.Context, sceneId string) error {
	scene, err := p.dbClient.SceneById(sceneId)
	if err != nil {
		return err
	}
	switch scene.Conditions[0].ConditionType {
	case "timer":
		conJobApp := resourceContainer.ConJobAppNameFrom(p.dic.Get)
		conJobApp.DeleteJob(scene.Id)
	case "notify":
		ekuiperApp := resourceContainer.EkuiperAppFrom(p.dic.Get)
		err = ekuiperApp.StopRule(ctx, sceneId)
		if err != nil {
			return err
		}
	default:
		return errort.NewCommonEdgeX(errort.DefaultReqParamsError, "condition Type value not much", nil)
	}

	return p.dbClient.SceneStop(sceneId)
}

func (p sceneApp) DelSceneById(ctx context.Context, sceneId string) error {
	scene, err := p.dbClient.SceneById(sceneId)
	if err != nil {
		return err
	}

	if len(scene.Conditions) == 0 {
		return p.dbClient.DeleteSceneById(sceneId)
		//return errort.NewCommonEdgeX(errort.DefaultSystemError, "conditions param error", nil)
	}

	switch scene.Conditions[0].ConditionType {
	case "timer":
		conJobApp := resourceContainer.ConJobAppNameFrom(p.dic.Get)
		conJobApp.DeleteJob(scene.Id)
	case "notify":
		ekuiperApp := resourceContainer.EkuiperAppFrom(p.dic.Get)
		err = ekuiperApp.DeleteRule(ctx, sceneId)
		if err != nil {
			return err
		}
	default:
		return errort.NewCommonEdgeX(errort.DefaultReqParamsError, "condition Type value not much", nil)
	}
	return p.dbClient.DeleteSceneById(sceneId)
}

func (p sceneApp) SceneSearch(ctx context.Context, req dtos.SceneSearchQueryRequest) ([]models.Scene, uint32, error) {
	offset, limit := req.BaseSearchConditionQuery.GetPage()
	resp, total, err := p.dbClient.SceneSearch(offset, limit, req)
	if err != nil {
		return []models.Scene{}, 0, err
	}

	return resp, total, nil
}

func (p sceneApp) SceneLogSearch(ctx context.Context, req dtos.SceneLogSearchQueryRequest) ([]models.SceneLog, uint32, error) {
	offset, limit := req.BaseSearchConditionQuery.GetPage()
	resp, total, err := p.dbClient.SceneLogSearch(offset, limit, req)
	if err != nil {
		return []models.SceneLog{}, 0, err
	}
	return resp, total, nil
}

func (p sceneApp) buildEkuiperSqlAndAction(req dtos.SceneUpdateRequest) (actions []dtos.Actions, sql string, err error) {
	configapp := resourceContainer.ConfigurationFrom(p.dic.Get)
	actions = dtos.GetRuleSceneEkuiperActions(configapp.Service.Url())
	option := req.Conditions[0].Option
	deviceId := option["device_id"]
	deviceName := option["device_name"]
	productId := option["product_id"]
	productName := option["product_name"]
	trigger := option["trigger"]
	code := option["code"]
	if deviceId == "" || deviceName == "" || productId == "" || productName == "" || code == "" || trigger == "" {
		err = errort.NewCommonEdgeX(errort.DefaultReqParamsError, "required parameter missing", nil)
		return
	}
	device, err := p.dbClient.DeviceById(deviceId)
	if err != nil {
		return
	}
	product, err := p.dbClient.ProductById(productId)
	if err != nil {
		return
	}
	if device.ProductId != product.Id {
		err = errort.NewCommonEdgeX(errort.DefaultSystemError, "", nil)
		return
	}

	switch trigger {
	case string(constants.DeviceDataTrigger):
		var codeFind bool
		for _, property := range product.Properties {
			if code == property.Code {
				codeFind = true
				if !property.TypeSpec.Type.AllowSendInEkuiper() {
					err = errort.NewCommonEdgeX(errort.DefaultReqParamsError, "required parameter missing", nil)
				}

				var s int
				switch option["value_cycle"] {
				case "1分钟周期":
					s = 60
				case "5分钟周期":
					s = 60 * 5
				case "15分钟周期":
					s = 60 * 15
				case "30分钟周期":
					s = 60 * 30
				case "60分钟周期":
					s = 60 * 60
				default:
				}

				switch property.TypeSpec.Type {

				case constants.SpecsTypeInt, constants.SpecsTypeFloat:
					valueType := option["value_type"]
					if valueType == "" {
						err = errort.NewCommonEdgeX(errort.DefaultReqParamsError, "required value_type parameter missing", nil)
						return
					}
					switch valueType {
					case constants.Original: //原始值
						decideCondition := option["decide_condition"]
						if decideCondition == "" {
							err = errort.NewCommonEdgeX(errort.DefaultReqParamsError, "required decide_condition parameter missing", nil)
							return
						}
						originalTemp := `SELECT rule_id(),json_path_query(data, "$.%s.time") as report_time ,deviceId FROM mqtt_stream where deviceId = "%s" and messageType = "PROPERTY_REPORT" and  json_path_exists(data, "$.%s") = true and json_path_query(data, "$.%s.value") %s`
						sql = fmt.Sprintf(originalTemp, code, deviceId, code, code, decideCondition)
						return
					case constants.Max:
						sqlTemp := `SELECT window_start(),window_end(),rule_id(),deviceId,max(json_path_query(data, "$.%s.value")) as max_%s FROM mqtt_stream where deviceId = "%s" and messageType = "PROPERTY_REPORT" and json_path_exists(data, "$.%s") = true GROUP BY %s HAVING max_%s %s`
						valueCycle := s
						if valueCycle == 0 {
							err = errort.NewCommonEdgeX(errort.DefaultReqParamsError, "required value_cycle parameter missing", nil)
							return
						}
						decideCondition := option["decide_condition"]
						if decideCondition == "" {
							err = errort.NewCommonEdgeX(errort.DefaultReqParamsError, "required decide_condition parameter missing", nil)
							return
						}
						sql = fmt.Sprintf(sqlTemp, code, code, device.Id, code, fmt.Sprintf("TUMBLINGWINDOW(ss, %d)", valueCycle), code, decideCondition)
						return
					case constants.Min:
						sqlTemp := `SELECT window_start(),window_end(),rule_id(),deviceId,min(json_path_query(data, "$.%s.value")) as min_%s FROM mqtt_stream where deviceId = "%s" and messageType = "PROPERTY_REPORT" and json_path_exists(data, "$.%s") = true GROUP BY %s HAVING min_%s %s`
						valueCycle := s
						if valueCycle == 0 {
							err = errort.NewCommonEdgeX(errort.DefaultReqParamsError, "required value_cycle parameter missing", nil)
							return
						}
						decideCondition := option["decide_condition"]
						if decideCondition == "" {
							err = errort.NewCommonEdgeX(errort.DefaultReqParamsError, "required decide_condition parameter missing", nil)
							return
						}
						sql = fmt.Sprintf(sqlTemp, code, code, device.Id, code, fmt.Sprintf("TUMBLINGWINDOW(ss, %d)", valueCycle), code, decideCondition)
						return
					case constants.Sum:
						sqlTemp := `SELECT window_start(),window_end(),rule_id(),deviceId,sum(json_path_query(data, "$.%s.value")) as sum_%s FROM mqtt_stream where deviceId = "%s" and messageType = "PROPERTY_REPORT" and json_path_exists(data, "$.%s") = true GROUP BY %s HAVING sum_%s %s`
						valueCycle := s
						if valueCycle == 0 {
							err = errort.NewCommonEdgeX(errort.DefaultReqParamsError, "required value_cycle parameter missing", nil)
							return
						}
						decideCondition := option["decide_condition"]
						if decideCondition == "" {
							err = errort.NewCommonEdgeX(errort.DefaultReqParamsError, "required decide_condition parameter missing", nil)
							return
						}
						sql = fmt.Sprintf(sqlTemp, code, code, device.Id, code, fmt.Sprintf("TUMBLINGWINDOW(ss, %d)", valueCycle), code, decideCondition)
						return
					case constants.Avg:
						sqlTemp := `SELECT window_start(),window_end(),rule_id(),deviceId,avg(json_path_query(data, "$.%s.value")) as avg_%s FROM mqtt_stream where deviceId = "%s" and messageType = "PROPERTY_REPORT" and json_path_exists(data, "$.%s") = true GROUP BY %s HAVING avg_%s %s`
						valueCycle := s
						if valueCycle == 0 {
							err = errort.NewCommonEdgeX(errort.DefaultReqParamsError, "required value_cycle parameter missing", nil)
							return
						}
						decideCondition := option["decide_condition"]
						if decideCondition == "" {
							err = errort.NewCommonEdgeX(errort.DefaultReqParamsError, "required decide_condition parameter missing", nil)
							return
						}
						sql = fmt.Sprintf(sqlTemp, code, code, device.Id, code, fmt.Sprintf("TUMBLINGWINDOW(ss, %d)", valueCycle), code, decideCondition)
						return
					}
				case constants.SpecsTypeText:
					decideCondition := option["decide_condition"]
					if decideCondition == "" {
						err = errort.NewCommonEdgeX(errort.DefaultReqParamsError, "required decide_condition parameter missing", nil)
						return
					}
					st := strings.Split(decideCondition, " ")
					if len(st) != 2 {
						return
					}
					sqlTemp := `SELECT rule_id(),json_path_query(data, "$.%s.time") as report_time,deviceId FROM mqtt_stream where deviceId = "%s" and messageType = "PROPERTY_REPORT" and  json_path_exists(data, "$.%s") = true and json_path_query(data, "$.%s.value") = "%s"`
					sql = fmt.Sprintf(sqlTemp, code, deviceId, code, code, st[1])
					return

				case constants.SpecsTypeBool:
					decideCondition := option["decide_condition"]
					if decideCondition == "" {
						err = errort.NewCommonEdgeX(errort.DefaultReqParamsError, "required decide_condition parameter missing", nil)
						return
					}
					st := strings.Split(decideCondition, " ")
					if len(st) != 2 {
						return
					}
					sqlTemp := `SELECT rule_id(),json_path_query(data, "$.%s.time") as report_time,deviceId FROM mqtt_stream where deviceId = "%s" and messageType = "PROPERTY_REPORT" and  json_path_exists(data, "$.%s") = true and json_path_query(data, "$.%s.value") = "%s"`
					if st[1] == "true" {
						sql = fmt.Sprintf(sqlTemp, code, deviceId, code, code, "1")
					} else if st[1] == "false" {
						sql = fmt.Sprintf(sqlTemp, code, deviceId, code, code, "0")
					}
					return
				case constants.SpecsTypeEnum:
					decideCondition := option["decide_condition"]
					if decideCondition == "" {
						err = errort.NewCommonEdgeX(errort.DefaultReqParamsError, "required decide_condition parameter missing", nil)
						return
					}
					st := strings.Split(decideCondition, " ")
					if len(st) != 2 {
						return
					}
					sqlTemp := `SELECT rule_id(),json_path_query(data, "$.%s.time") as report_time,deviceId FROM mqtt_stream where deviceId = "%s" and messageType = "PROPERTY_REPORT" and  json_path_exists(data, "$.%s") = true and json_path_query(data, "$.%s.value") = "%s"`
					sql = fmt.Sprintf(sqlTemp, code, deviceId, code, code, st[1])
					return
				}
			}
		}
		if !codeFind {
			err = errort.NewCommonEdgeX(errort.DefaultReqParamsError, "required code parameter missing", nil)
		}
	case string(constants.DeviceEventTrigger):
		var codeFind bool
		for _, event := range product.Events {
			if code == event.Code {
				codeFind = true
				sqlTemp := `SELECT rule_id(),json_path_query(data, "$.eventTime") as report_time,deviceId FROM mqtt_stream where deviceId = "%s" and messageType = "EVENT_REPORT" and  json_path_exists(data, "$.eventCode") = true and json_path_query(data, "$.eventCode") = "%s"`
				sql = fmt.Sprintf(sqlTemp, device.Id, code)
				return
			}
		}
		if !codeFind {
			err = errort.NewCommonEdgeX(errort.DefaultReqParamsError, "required code parameter missing", nil)
		}
	case string(constants.DeviceStatusTrigger):
		var status string
		deviceStatus := option["status"]
		if deviceStatus == "" {
			err = errort.NewCommonEdgeX(errort.DefaultReqParamsError, "required status parameter missing", nil)
			return
		}
		if deviceStatus == "在线" {
			status = constants.DeviceOnline
		} else if deviceStatus == "离线" {
			status = constants.DeviceOffline
		} else {
			err = errort.NewCommonEdgeX(errort.DefaultReqParamsError, "required status parameter missing", nil)
			return
		}
		sqlTemp := `SELECT rule_id(),json_path_query(data, "$.time") as report_time,deviceId FROM mqtt_stream where deviceId = "%s" and messageType = "DEVICE_STATUS" and  json_path_exists(data, "$.status") = true and json_path_query(data, "$.status") = "%s"`
		sql = fmt.Sprintf(sqlTemp, device.Id, status)
		return
	default:
		err = errort.NewCommonEdgeX(errort.DefaultReqParamsError, "required trigger parameter missing", nil)
		return
	}
	return
}

func (p sceneApp) checkAlertRuleParam(ctx context.Context, scene models.Scene, operate string) error {
	if operate == "start" {
		if scene.Status == constants.SceneStart {
			return errort.NewCommonErr(errort.AlertRuleStatusStarting, fmt.Errorf("scene id(%s) is runing ,not allow start", scene.Id))
		}
	}

	var (
		trigger string
	)

	if len(scene.Conditions) != 1 {
		trigger = scene.Conditions[0].Option["trigger"]

		switch scene.Conditions[0].ConditionType {
		case "timer":
		case "notify":
			ekuiperApp := resourceContainer.EkuiperAppFrom(p.dic.Get)
			exist, err := ekuiperApp.RuleExist(ctx, scene.Id)
			if err != nil {
				return err
			}
			if !exist {

			}

			trigger = scene.Conditions[0].Option["trigger"]
			if trigger != string(constants.DeviceDataTrigger) || trigger != string(constants.DeviceEventTrigger) || trigger != string(constants.DeviceStatusTrigger) {
				return errort.NewCommonErr(errort.SceneRuleParamsError, fmt.Errorf("scene id(%s) trigger param error", scene.Id))
			}

			option := scene.Conditions[0].Option
			deviceId := option["device_id"]
			deviceName := option["device_name"]
			productId := option["product_id"]
			productName := option["product_name"]
			//trigger := option["trigger"]
			code := option["code"]
			if deviceId == "" || deviceName == "" || productId == "" || productName == "" || code == "" || trigger == "" {
				return errort.NewCommonEdgeX(errort.SceneRuleParamsError, "required parameter missing", nil)
			}
			device, err := p.dbClient.DeviceById(deviceId)
			if err != nil {
				return errort.NewCommonErr(errort.SceneRuleParamsError, fmt.Errorf("scene id(%s) device not found", scene.Id))

			}
			product, err := p.dbClient.ProductById(productId)
			if err != nil {
				return errort.NewCommonErr(errort.SceneRuleParamsError, fmt.Errorf("scene id(%s) actions is null", scene.Id))

			}
			if device.ProductId != product.Id {
				return errort.NewCommonErr(errort.SceneRuleParamsError, fmt.Errorf("scene id(%s) actions is null", scene.Id))
			}

		default:
			return errort.NewCommonErr(errort.SceneRuleParamsError, fmt.Errorf("scene id(%s) conditionType param errror", scene.Id))

		}
	}

	//-------------------------

	if len(scene.Actions) == 0 {
		return errort.NewCommonErr(errort.SceneRuleParamsError, fmt.Errorf("scene id(%s) actions is null", scene.Id))
	}

	for _, action := range scene.Actions {
		//检查产品和设备是否存在
		device, err := p.dbClient.DeviceById(action.DeviceID)
		if err != nil {
			return errort.NewCommonErr(errort.SceneRuleParamsError, fmt.Errorf("scene id(%s) device not found", scene.Id))
		}

		product, err := p.dbClient.ProductById(action.ProductID)
		if err != nil {
			return errort.NewCommonErr(errort.SceneRuleParamsError, fmt.Errorf("scene id(%s) product not found", scene.Id))
		}
		if device.ProductId != product.Id {
			return errort.NewCommonErr(errort.SceneRuleParamsError, fmt.Errorf("scene id(%s) actions is null", scene.Id))
		}

		var find bool

		if trigger == string(constants.DeviceDataTrigger) {
			for _, property := range product.Properties {
				if property.Code == action.Code {
					find = true
					break
				}
			}
		}

		if trigger == string(constants.DeviceEventTrigger) {
			for _, event := range product.Events {
				if event.Code == action.Code {
					find = true
					break
				}
			}
		}

		if !find {
			return errort.NewCommonErr(errort.SceneRuleParamsError, fmt.Errorf("scene id(%s) code not found", scene.Id))

		}
		if action.Value == "" {
			return errort.NewCommonErr(errort.SceneRuleParamsError, fmt.Errorf("scene id(%s) value is null", scene.Id))
		}

	}
	return nil
}

func (p sceneApp) EkuiperNotify(ctx context.Context, req map[string]interface{}) error {
	sceneId, ok := req["rule_id"]
	if !ok {
		return errort.NewCommonErr(errort.DefaultReqParamsError, errors.New(""))
	}
	var (
		coverSceneId string
	)
	switch sceneId.(type) {
	case string:
		coverSceneId = sceneId.(string)
	case int:
		coverSceneId = strconv.Itoa(sceneId.(int))
	case int64:
		coverSceneId = strconv.Itoa(int(sceneId.(int64)))
	case float64:
		coverSceneId = fmt.Sprintf("%f", sceneId.(float64))
	case float32:
		coverSceneId = fmt.Sprintf("%f", sceneId.(float64))
	}
	if coverSceneId == "" {
		return errort.NewCommonErr(errort.DefaultReqParamsError, errors.New(""))
	}

	scene, err := p.dbClient.SceneById(coverSceneId)
	if err != nil {
		return err
	}

	for _, action := range scene.Actions {
		deviceApp := resourceContainer.DeviceItfFrom(p.dic.Get)
		execRes := deviceApp.DeviceAction(dtos.JobAction{
			ProductId: action.ProductID,
		})
		_, err := p.dbClient.AddSceneLog(models.SceneLog{
			SceneId: scene.Id,
			Name:    scene.Name,
			ExecRes: execRes.ToString(),
		})
		if err != nil {
			p.lc.Errorf("add sceneLog err %v", err.Error())

		}

	}
	return nil
}

func (p sceneApp) CheckSceneByDeviceId(ctx context.Context, deviceId string) error {
	var req dtos.SceneSearchQueryRequest
	req.Status = string(constants.SceneStart)
	scenes, _, err := p.SceneSearch(ctx, req)
	if err != nil {
		return err
	}

	for _, scene := range scenes {
		for _, condition := range scene.Conditions {

			if condition.Option != nil && condition.Option["device_id"] == deviceId {
				return errort.NewCommonEdgeX(errort.DeviceAssociationSceneRule, "This device has been bound to scene rules. Please stop reporting scene rules before proceeding with the operation", nil)
			}
		}
		for _, action := range scene.Actions {
			if action.DeviceID == deviceId {
				return errort.NewCommonEdgeX(errort.DeviceAssociationSceneRule, "This device has been bound to scene rules. Please stop reporting scene rules before proceeding with the operation", nil)
			}
		}
	}

	return nil
}
