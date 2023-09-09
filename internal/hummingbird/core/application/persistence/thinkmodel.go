/*******************************************************************************
 * Copyright 2017 Dell Inc.
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

package persistence

import (
	"context"
	"encoding/json"
	"github.com/winc-link/edge-driver-proto/thingmodel"
	"github.com/winc-link/hummingbird/internal/dtos"
	"github.com/winc-link/hummingbird/internal/hummingbird/core/application/messagestore"
	resourceContainer "github.com/winc-link/hummingbird/internal/hummingbird/core/container"
	interfaces "github.com/winc-link/hummingbird/internal/hummingbird/core/interface"
	"github.com/winc-link/hummingbird/internal/models"
	"github.com/winc-link/hummingbird/internal/pkg/constants"
	"github.com/winc-link/hummingbird/internal/pkg/container"
	"github.com/winc-link/hummingbird/internal/pkg/di"
	"github.com/winc-link/hummingbird/internal/pkg/logger"
	"github.com/winc-link/hummingbird/internal/pkg/utils"
	"strconv"
)

type persistApp struct {
	dic          *di.Container
	lc           logger.LoggingClient
	dbClient     interfaces.DBClient
	dataDbClient interfaces.DataDBClient
}

func NewPersistApp(dic *di.Container) *persistApp {
	lc := container.LoggingClientFrom(dic.Get)
	dbClient := resourceContainer.DBClientFrom(dic.Get)
	dataDbClient := resourceContainer.DataDBClientFrom(dic.Get)
	pstApp := &persistApp{
		lc:           lc,
		dic:          dic,
		dbClient:     dbClient,
		dataDbClient: dataDbClient,
	}

	return pstApp
}

func (pst *persistApp) SaveDeviceThingModelData(req dtos.ThingModelMessage) error {
	switch pst.dataDbClient.GetDataDBType() {
	case constants.LevelDB:
		return pst.saveDeviceThingModelToLevelDB(req)
	case constants.TDengine:
		return pst.saveDeviceThingModelToTdengine(req)
	default:
		return nil
	}
}

func (pst *persistApp) saveDeviceThingModelToLevelDB(req dtos.ThingModelMessage) error {
	switch req.GetOpType() {
	case thingmodel.OperationType_PROPERTY_REPORT:
		propertyMsg, err := req.TransformMessageDataByProperty()
		if err != nil {
			return err
		}
		kvs := make(map[string]interface{})
		for s, data := range propertyMsg.Data {
			key := generatePropertyLeveldbKey(req.Cid, s, data.Time)
			value, err := data.Marshal()
			if err != nil {
				continue
			}
			kvs[key] = value
		}
		//批量写。
		err = pst.dataDbClient.Insert(context.Background(), "", kvs)
		if err != nil {
			return err
		}
	case thingmodel.OperationType_EVENT_REPORT:
		eventMsg, err := req.TransformMessageDataByEvent()
		if err != nil {
			return err
		}
		kvs := make(map[string]interface{})
		var key string
		key = generateEventLeveldbKey(req.Cid, eventMsg.Data.EventCode, eventMsg.Data.EventTime)
		value, _ := eventMsg.Data.Marshal()
		kvs[key] = value
		//批量写。
		err = pst.dataDbClient.Insert(context.Background(), "", kvs)
		if err != nil {
			return err
		}
	case thingmodel.OperationType_SERVICE_EXECUTE:
		serviceMsg, err := req.TransformMessageDataByService()
		kvs := make(map[string]interface{})
		var key string
		key = generateActionLeveldbKey(req.Cid, serviceMsg.Code, serviceMsg.Time)
		value, _ := serviceMsg.Marshal()
		kvs[key] = value
		err = pst.dataDbClient.Insert(context.Background(), "", kvs)

		if err != nil {
			return err
		}

	case thingmodel.OperationType_PROPERTY_GET_RESPONSE:
		msg, err := req.TransformMessageDataByGetProperty()
		if err != nil {
			return err
		}
		messageStore := resourceContainer.MessageStoreItfFrom(pst.dic.Get)
		ack, ok := messageStore.LoadMsgChan(msg.MsgId)
		if !ok {
			//超时了。
			return nil
		}
		if v, ok := ack.(*messagestore.MsgAckChan); ok {
			v.TrySendDataAndCloseChan(msg.Data)
			messageStore.DeleteMsgId(msg.MsgId)
		}

	case thingmodel.OperationType_PROPERTY_SET_RESPONSE:
		msg, err := req.TransformMessageDataBySetProperty()
		if err != nil {
			return err
		}
		messageStore := resourceContainer.MessageStoreItfFrom(pst.dic.Get)
		ack, ok := messageStore.LoadMsgChan(msg.MsgId)
		if !ok {
			//超时了。
			return nil
		}
		if v, ok := ack.(*messagestore.MsgAckChan); ok {
			v.TrySendDataAndCloseChan(msg.Data)
			messageStore.DeleteMsgId(msg.MsgId)
		}
	case thingmodel.OperationType_SERVICE_EXECUTE_RESPONSE:
		serviceMsg, err := req.TransformMessageDataByServiceExec()
		if err != nil {
			return err
		}
		//var find bool
		//var callType constants.CallType
		//
		//for _, action := range product.Actions {
		//	if action.Code == serviceMsg.Code {
		//		find = true
		//		callType = action.CallType
		//		break
		//	}
		//}
		//
		//if !find {
		//	return errors.New("")
		//}
		messageStore := resourceContainer.MessageStoreItfFrom(pst.dic.Get)
		ack, ok := messageStore.LoadMsgChan(serviceMsg.MsgId)
		if !ok {
			//可能是超时了。
			return nil
		}

		if v, ok := ack.(*messagestore.MsgAckChan); ok {
			v.TrySendDataAndCloseChan(serviceMsg.OutputParams)
			messageStore.DeleteMsgId(serviceMsg.MsgId)
		}
		//if callType == constants.CallTypeSync {
		//	messageStore := resourceContainer.MessageStoreItfFrom(pst.dic.Get)
		//	ack, ok := messageStore.LoadMsgChan(serviceMsg.MsgId)
		//	if !ok {
		//		//可能是超时了。
		//		return nil
		//	}
		//
		//	if v, ok := ack.(*messagestore.MsgAckChan); ok {
		//		v.TrySendDataAndCloseChan(serviceMsg.OutputParams)
		//		messageStore.DeleteMsgId(serviceMsg.MsgId)
		//	}
		//
		//} else if callType == constants.CallTypeAsync {
		//	kvs := make(map[string]interface{})
		//	var key string
		//	key = generateActionLeveldbKey(req.Cid, serviceMsg.Code, serviceMsg.Time)
		//	value, _ := serviceMsg.Marshal()
		//	kvs[key] = value
		//	err = pst.dataDbClient.Insert(context.Background(), "", kvs)
		//
		//	if err != nil {
		//		return err
		//	}
		//}
	case thingmodel.OperationType_DATA_BATCH_REPORT:
		msg, err := req.TransformMessageDataByBatchReport()
		if err != nil {
			return err
		}
		t := msg.Time
		kvs := make(map[string]interface{})

		for code, property := range msg.Data.Properties {
			var data dtos.ReportData
			data.Value = property.Value
			data.Time = t
			key := generatePropertyLeveldbKey(req.Cid, code, t)
			value, err := data.Marshal()
			if err != nil {
				continue
			}
			kvs[key] = value
		}

		for code, event := range msg.Data.Events {
			var data dtos.EventData
			data.OutputParams = event.OutputParams
			data.EventTime = t
			data.EventCode = code
			key := generateEventLeveldbKey(req.Cid, code, t)
			value, _ := data.Marshal()
			kvs[key] = value

		}
		//批量写。
		err = pst.dataDbClient.Insert(context.Background(), "", kvs)

		if err != nil {
			return err
		}
		return nil
	}
	return nil
}

func (pst *persistApp) saveDeviceThingModelToTdengine(req dtos.ThingModelMessage) error {
	switch req.GetOpType() {
	case thingmodel.OperationType_PROPERTY_REPORT:
		propertyMsg, err := req.TransformMessageDataByProperty()
		if err != nil {
			return err
		}
		data := make(map[string]interface{})
		for s, reportData := range propertyMsg.Data {
			data[s] = reportData.Value
		}
		err = pst.dataDbClient.Insert(context.Background(), constants.DB_PREFIX+req.Cid, data)
		if err != nil {
			return err
		}

	case thingmodel.OperationType_EVENT_REPORT:
		eventMsg, err := req.TransformMessageDataByEvent()
		if err != nil {
			return err
		}
		data := make(map[string]interface{})
		data[eventMsg.Data.EventCode] = eventMsg.Data
		err = pst.dataDbClient.Insert(context.Background(), constants.DB_PREFIX+req.Cid, data)
		if err != nil {
			return err
		}

	case thingmodel.OperationType_SERVICE_EXECUTE:
		serviceMsg, err := req.TransformMessageDataByService()
		if err != nil {
			return err
		}
		v, _ := serviceMsg.Marshal()
		data := make(map[string]interface{})
		data[serviceMsg.Code] = string(v)
		err = pst.dataDbClient.Insert(context.Background(), constants.DB_PREFIX+req.Cid, data)
		if err != nil {
			return err
		}

	case thingmodel.OperationType_PROPERTY_GET_RESPONSE:
		msg, err := req.TransformMessageDataByGetProperty()
		if err != nil {
			return err
		}
		messageStore := resourceContainer.MessageStoreItfFrom(pst.dic.Get)
		ack, ok := messageStore.LoadMsgChan(msg.MsgId)
		if !ok {
			//超时了。
			return nil
		}
		if v, ok := ack.(*messagestore.MsgAckChan); ok {
			v.TrySendDataAndCloseChan(msg.Data)
			messageStore.DeleteMsgId(msg.MsgId)
		}
	case thingmodel.OperationType_PROPERTY_SET_RESPONSE:
		msg, err := req.TransformMessageDataBySetProperty()
		if err != nil {
			return err
		}
		messageStore := resourceContainer.MessageStoreItfFrom(pst.dic.Get)
		ack, ok := messageStore.LoadMsgChan(msg.MsgId)
		if !ok {
			//超时了。
			return nil
		}
		if v, ok := ack.(*messagestore.MsgAckChan); ok {
			v.TrySendDataAndCloseChan(msg.Data)
			messageStore.DeleteMsgId(msg.MsgId)
		}
	case thingmodel.OperationType_DATA_BATCH_REPORT:

		msg, err := req.TransformMessageDataByBatchReport()
		if err != nil {
			return err
		}
		data := make(map[string]interface{})

		for code, property := range msg.Data.Properties {
			data[code] = property.Value
		}
		for code, event := range msg.Data.Events {
			var eventData dtos.EventData
			eventData.OutputParams = event.OutputParams
			eventData.EventCode = code
			eventData.EventTime = msg.Time
			data[code] = eventData

		}
		//批量写。
		err = pst.dataDbClient.Insert(context.Background(), constants.DB_PREFIX+req.Cid, data)

		if err != nil {
			return err
		}
		return nil
	case thingmodel.OperationType_SERVICE_EXECUTE_RESPONSE:
		serviceMsg, err := req.TransformMessageDataByServiceExec()
		if err != nil {
			return err
		}

		device, err := pst.dbClient.DeviceById(req.Cid)
		if err != nil {
			return err
		}

		_, err = pst.dbClient.ProductById(device.ProductId)
		if err != nil {
			return err
		}

		//var find bool
		//var callType constants.CallType
		//
		//for _, action := range product.Actions {
		//	if action.Code == serviceMsg.Code {
		//		find = true
		//		callType = action.CallType
		//		break
		//	}
		//}
		//
		//if !find {
		//	return errors.New("")
		//}

		messageStore := resourceContainer.MessageStoreItfFrom(pst.dic.Get)
		ack, ok := messageStore.LoadMsgChan(serviceMsg.MsgId)
		if !ok {
			//可能是超时了。
			return nil
		}

		if v, ok := ack.(*messagestore.MsgAckChan); ok {
			v.TrySendDataAndCloseChan(serviceMsg.OutputParams)
			messageStore.DeleteMsgId(serviceMsg.MsgId)
		}

		//if callType == constants.CallTypeSync {
		//	messageStore := resourceContainer.MessageStoreItfFrom(pst.dic.Get)
		//	ack, ok := messageStore.LoadMsgChan(serviceMsg.MsgId)
		//	if !ok {
		//		//可能是超时了。
		//		return nil
		//	}
		//
		//	if v, ok := ack.(*messagestore.MsgAckChan); ok {
		//		v.TrySendDataAndCloseChan(serviceMsg.OutputParams)
		//		messageStore.DeleteMsgId(serviceMsg.MsgId)
		//	}
		//
		//} else if callType == constants.CallTypeAsync {
		//	v, _ := serviceMsg.Marshal()
		//	data := make(map[string]interface{})
		//	data[serviceMsg.Code] = string(v)
		//	err = pst.dataDbClient.Insert(context.Background(), constants.DB_PREFIX+req.Cid, data)
		//	if err != nil {
		//		return err
		//	}
		//}

	}
	return nil
}

func generatePropertyLeveldbKey(cid, code string, reportTime int64) string {
	return cid + "-" + constants.Property + "-" + code + "-" + strconv.Itoa(int(reportTime))
}

func generateOncePropertyLeveldbKey(cid, code string) string {
	return cid + "-" + constants.Property + "-" + code
}

func generateEventLeveldbKey(cid, code string, reportTime int64) string {
	return cid + "-" + constants.Event + "-" + code + "-" + strconv.Itoa(int(reportTime))
}

func generateOnceEventLeveldbKey(cid, code string) string {
	return cid + "-" + constants.Event + "-" + code
}

func generateActionLeveldbKey(cid, code string, reportTime int64) string {
	return cid + "-" + constants.Action + "-" + code + "-" + strconv.Itoa(int(reportTime))
}

func generateOnceActionLeveldbKey(cid, code string) string {
	return cid + "-" + constants.Action + "-" + code
}

func (pst *persistApp) searchDeviceThingModelPropertyDataFromLevelDB(req dtos.ThingModelPropertyDataRequest) (interface{}, error) {
	deviceInfo, err := pst.dbClient.DeviceById(req.DeviceId)
	if err != nil {
		return nil, err
	}
	var productInfo models.Product
	response := make([]dtos.ThingModelDataResponse, 0)
	productInfo, err = pst.dbClient.ProductById(deviceInfo.ProductId)
	if err != nil {
		return nil, err
	}
	if req.Code == "" {
		for _, property := range productInfo.Properties {
			req.Code = property.Code
			ksv, _, err := pst.dataDbClient.GetDeviceProperty(req, deviceInfo)
			if err != nil {
				pst.lc.Errorf("GetDeviceProperty error %+v", err)
				continue
			}
			var reportData dtos.ReportData
			if len(ksv) > 0 {
				reportData = ksv[0]
			}
			var unit string
			if property.TypeSpec.Type == constants.SpecsTypeInt || property.TypeSpec.Type == constants.SpecsTypeFloat {
				var typeSpecIntOrFloat models.TypeSpecIntOrFloat
				_ = json.Unmarshal([]byte(property.TypeSpec.Specs), &typeSpecIntOrFloat)
				unit = typeSpecIntOrFloat.Unit
			} else if property.TypeSpec.Type == constants.SpecsTypeEnum {
				//enum 的单位需要特殊处理一下
				enumTypeSpec := make(map[string]string)
				_ = json.Unmarshal([]byte(property.TypeSpec.Specs), &enumTypeSpec)

				for key, value := range enumTypeSpec {
					s := utils.InterfaceToString(reportData.Value)
					if key == s {
						unit = value
					}
				}
			}
			response = append(response, dtos.ThingModelDataResponse{
				ReportData: reportData,
				Code:       property.Code,
				DataType:   string(property.TypeSpec.Type),
				Name:       property.Name,
				Unit:       unit,
				AccessMode: property.AccessMode,
			})
		}
	}
	return response, nil
}

func (pst *persistApp) searchDeviceThingModelHistoryPropertyDataFromTDengine(req dtos.ThingModelPropertyDataRequest) (interface{}, int, error) {
	var count int
	deviceInfo, err := pst.dbClient.DeviceById(req.DeviceId)
	if err != nil {
		return nil, count, err
	}
	var productInfo models.Product
	productInfo, err = pst.dbClient.ProductById(deviceInfo.ProductId)
	if err != nil {
		return nil, count, err
	}
	var response []dtos.ReportData

	for _, property := range productInfo.Properties {
		if property.Code == req.Code {
			req.Code = property.Code
			response, count, err = pst.dataDbClient.GetDeviceProperty(req, deviceInfo)
			if err != nil {
				pst.lc.Errorf("GetDeviceProperty error %+v", err)
			}
			var typeSpecIntOrFloat models.TypeSpecIntOrFloat
			if property.TypeSpec.Type == constants.SpecsTypeInt || property.TypeSpec.Type == constants.SpecsTypeFloat {
				_ = json.Unmarshal([]byte(property.TypeSpec.Specs), &typeSpecIntOrFloat)
			}

			if typeSpecIntOrFloat.Unit == "" {
				typeSpecIntOrFloat.Unit = "-"
			}
			break
		}
	}
	return response, count, nil
}

func (pst *persistApp) searchDeviceThingModelPropertyDataFromTDengine(req dtos.ThingModelPropertyDataRequest) (interface{}, error) {
	deviceInfo, err := pst.dbClient.DeviceById(req.DeviceId)
	if err != nil {
		return nil, err
	}
	var productInfo models.Product
	response := make([]dtos.ThingModelDataResponse, 0)
	productInfo, err = pst.dbClient.ProductById(deviceInfo.ProductId)
	if err != nil {
		return nil, err
	}
	if req.Code == "" {
		for _, property := range productInfo.Properties {
			req.Code = property.Code
			ksv, _, err := pst.dataDbClient.GetDeviceProperty(req, deviceInfo)
			if err != nil {
				pst.lc.Errorf("GetDeviceProperty error %+v", err)
				continue
			}
			var reportData dtos.ReportData
			if len(ksv) > 0 {
				reportData = ksv[0]
			}
			var unit string
			if property.TypeSpec.Type == constants.SpecsTypeInt || property.TypeSpec.Type == constants.SpecsTypeFloat {
				var typeSpecIntOrFloat models.TypeSpecIntOrFloat
				_ = json.Unmarshal([]byte(property.TypeSpec.Specs), &typeSpecIntOrFloat)
				unit = typeSpecIntOrFloat.Unit
			} else if property.TypeSpec.Type == constants.SpecsTypeEnum {
				//enum 的单位需要特殊处理一下
				enumTypeSpec := make(map[string]string)
				_ = json.Unmarshal([]byte(property.TypeSpec.Specs), &enumTypeSpec)
				//pst.lc.Info("reportDataType enumTypeSpec", enumTypeSpec)

				for key, value := range enumTypeSpec {
					s := utils.InterfaceToString(reportData.Value)
					if key == s {
						unit = value
					}
				}
			}

			if unit == "" {
				unit = "-"
			}
			response = append(response, dtos.ThingModelDataResponse{
				ReportData: reportData,
				Code:       property.Code,
				DataType:   string(property.TypeSpec.Type),
				Name:       property.Name,
				Unit:       unit,
				AccessMode: property.AccessMode,
			})
		}
	}
	return response, nil
}

func (pst *persistApp) SearchDeviceThingModelPropertyData(req dtos.ThingModelPropertyDataRequest) (interface{}, error) {

	switch pst.dataDbClient.GetDataDBType() {
	case constants.LevelDB:
		return pst.searchDeviceThingModelPropertyDataFromLevelDB(req)

	case constants.TDengine:
		return pst.searchDeviceThingModelPropertyDataFromTDengine(req)

	default:
		return make([]interface{}, 0), nil
	}
}

func (pst *persistApp) searchDeviceThingModelHistoryPropertyDataFromLevelDB(req dtos.ThingModelPropertyDataRequest) (interface{}, int, error) {
	deviceInfo, err := pst.dbClient.DeviceById(req.DeviceId)
	if err != nil {
		return nil, 0, err
	}
	var count int
	var productInfo models.Product
	productInfo, err = pst.dbClient.ProductById(deviceInfo.ProductId)
	if err != nil {
		return nil, 0, err
	}
	var response []dtos.ReportData
	for _, property := range productInfo.Properties {
		if property.Code == req.Code {
			response, count, err = pst.dataDbClient.GetDeviceProperty(req, deviceInfo)
			if err != nil {
				pst.lc.Errorf("GetHistoryDeviceProperty error %+v", err)
			}
			break
		}
	}
	return response, count, nil
}

func (pst *persistApp) SearchDeviceThingModelHistoryPropertyData(req dtos.ThingModelPropertyDataRequest) (interface{}, int, error) {
	switch pst.dataDbClient.GetDataDBType() {
	case constants.LevelDB:
		return pst.searchDeviceThingModelHistoryPropertyDataFromLevelDB(req)
	case constants.TDengine:
		return pst.searchDeviceThingModelHistoryPropertyDataFromTDengine(req)
	}
	response := make([]interface{}, 0)
	return response, 0, nil
}

func (pst *persistApp) searchDeviceThingModelServiceDataFromLevelDB(req dtos.ThingModelServiceDataRequest) ([]dtos.ThingModelServiceDataResponse, int, error) {
	var count int
	deviceInfo, err := pst.dbClient.DeviceById(req.DeviceId)
	if err != nil {
		return nil, count, err
	}
	var productInfo models.Product
	productInfo, err = pst.dbClient.ProductById(deviceInfo.ProductId)
	if err != nil {
		return nil, count, err
	}
	var response dtos.ThingModelServiceDataResponseArray
	if req.Code == "" {
		for _, action := range productInfo.Actions {
			req.Code = action.Code
			var ksv []dtos.SaveServiceIssueData
			ksv, count, err = pst.dataDbClient.GetDeviceService(req, deviceInfo, productInfo)
			if err != nil {
				continue
			}
			for _, data := range ksv {
				response = append(response, dtos.ThingModelServiceDataResponse{
					ServiceName: action.Name,
					Code:        data.Code,
					InputData:   data.InputParams,
					OutputData:  data.OutputParams,
					ReportTime:  data.Time,
				})
			}
		}
	} else {
		var ksv []dtos.SaveServiceIssueData
		ksv, count, err = pst.dataDbClient.GetDeviceService(req, deviceInfo, productInfo)

		if err != nil {
			return nil, count, err
		}
		for _, data := range ksv {
			name := getServiceName(productInfo.Actions, data.Code)
			response = append(response, dtos.ThingModelServiceDataResponse{
				ServiceName: name,
				Code:        data.Code,
				InputData:   data.InputParams,
				OutputData:  data.OutputParams,
				ReportTime:  data.Time,
			})
		}
	}
	return response, count, nil
}

func (pst *persistApp) searchDeviceThingModelServiceDataFromTDengine(req dtos.ThingModelServiceDataRequest) ([]dtos.ThingModelServiceDataResponse, int, error) {

	var response dtos.ThingModelServiceDataResponseArray
	var count int
	deviceInfo, err := pst.dbClient.DeviceById(req.DeviceId)
	if err != nil {
		return response, count, err
	}
	var productInfo models.Product
	productInfo, err = pst.dbClient.ProductById(deviceInfo.ProductId)
	if err != nil {
		return response, count, err

	}
	serviceData, count, err := pst.dataDbClient.GetDeviceService(req, deviceInfo, productInfo)

	if err != nil {
		return response, count, err
	}

	for _, data := range serviceData {
		name := getServiceName(productInfo.Actions, data.Code)
		response = append(response, dtos.ThingModelServiceDataResponse{
			InputData:   data.InputParams,
			OutputData:  data.OutputParams,
			Code:        data.Code,
			ReportTime:  data.Time,
			ServiceName: name,
		})
	}

	return response, count, nil

}

func (pst *persistApp) SearchDeviceThingModelServiceData(req dtos.ThingModelServiceDataRequest) ([]dtos.ThingModelServiceDataResponse, int, error) {
	var response dtos.ThingModelServiceDataResponseArray
	switch pst.dataDbClient.GetDataDBType() {
	case constants.LevelDB:
		return pst.searchDeviceThingModelServiceDataFromLevelDB(req)
	case constants.TDengine:
		return pst.searchDeviceThingModelServiceDataFromTDengine(req)
	}
	return response, 0, nil
}

func (pst *persistApp) searchDeviceThingModelEventDataFromLevelDB(req dtos.ThingModelEventDataRequest) (dtos.ThingModelEventDataResponseArray, int, error) {
	var count int

	deviceInfo, err := pst.dbClient.DeviceById(req.DeviceId)
	if err != nil {
		return nil, count, err
	}
	var response dtos.ThingModelEventDataResponseArray
	var productInfo models.Product
	productInfo, err = pst.dbClient.ProductById(deviceInfo.ProductId)
	if err != nil {
		return nil, count, err
	}
	var ksv []dtos.EventData
	ksv, count, err = pst.dataDbClient.GetDeviceEvent(req, deviceInfo, productInfo)
	if err != nil {
		return nil, count, err
	}
	for _, data := range ksv {
		var eventType string
		var name string
		eventType, name = getEventTypeAndName(productInfo.Events, data.EventCode)
		response = append(response, dtos.ThingModelEventDataResponse{
			EventCode:  data.EventCode,
			EventType:  eventType,
			OutputData: data.OutputParams,
			ReportTime: data.EventTime,
			Name:       name,
		})
	}
	return response, count, nil
}

func (pst *persistApp) searchDeviceThingModelEventDataFromTDengine(req dtos.ThingModelEventDataRequest) (dtos.ThingModelEventDataResponseArray, int, error) {
	var response dtos.ThingModelEventDataResponseArray
	var count int
	deviceInfo, err := pst.dbClient.DeviceById(req.DeviceId)
	if err != nil {
		return response, count, err
	}
	var productInfo models.Product
	productInfo, err = pst.dbClient.ProductById(deviceInfo.ProductId)
	if err != nil {
		return response, count, err

	}
	eventData, count, err := pst.dataDbClient.GetDeviceEvent(req, deviceInfo, productInfo)

	if err != nil {
		return response, count, err
	}

	for _, data := range eventData {
		var eventType string
		var name string
		eventType, name = getEventTypeAndName(productInfo.Events, data.EventCode)
		response = append(response, dtos.ThingModelEventDataResponse{
			EventCode:  data.EventCode,
			EventType:  eventType,
			Name:       name,
			OutputData: data.OutputParams,
			ReportTime: data.EventTime,
		})

	}
	return response, count, nil
}

func (pst *persistApp) SearchDeviceThingModelEventData(req dtos.ThingModelEventDataRequest) ([]dtos.ThingModelEventDataResponse, int, error) {
	var response dtos.ThingModelEventDataResponseArray
	switch pst.dataDbClient.GetDataDBType() {
	case constants.LevelDB:
		return pst.searchDeviceThingModelEventDataFromLevelDB(req)
	case constants.TDengine:
		return pst.searchDeviceThingModelEventDataFromTDengine(req)
	}
	return response, 0, nil

}

func (pst *persistApp) searchDeviceMsgCountFromLevelDB(startTime, endTime int64) (int, error) {
	var (
		count int
		err   error
	)

	devices, _, err := pst.dbClient.DevicesSearch(0, -1, dtos.DeviceSearchQueryRequest{})
	if err != nil {
		return 0, err
	}

	for _, device := range devices {
		product, err := pst.dbClient.ProductById(device.ProductId)
		if err != nil {
			pst.lc.Errorf("search product:", err)
		}
		for _, property := range product.Properties {
			var req dtos.ThingModelPropertyDataRequest
			req.DeviceId = device.Id
			req.Code = property.Code
			req.Range = append(req.Range, startTime, endTime)
			propertyCount, err := pst.dataDbClient.GetDevicePropertyCount(req)
			if err != nil {
				return 0, err
			}
			count += propertyCount
		}

		for _, event := range product.Events {
			var req dtos.ThingModelEventDataRequest
			req.DeviceId = device.Id
			req.EventCode = event.Code
			req.Range = append(req.Range, startTime, endTime)
			eventCount, err := pst.dataDbClient.GetDeviceEventCount(req)
			if err != nil {
				return 0, err
			}
			count += eventCount
		}
	}

	return count, err

}

func (pst *persistApp) searchDeviceMsgCountFromTDengine(startTime, endTime int64) (int, error) {
	var (
		count int
		err   error
	)

	devices, _, err := pst.dbClient.DevicesSearch(0, -1, dtos.DeviceSearchQueryRequest{})
	if err != nil {
		return 0, err
	}

	for _, device := range devices {
		msgCount, err := pst.dataDbClient.GetDeviceMsgCountByGiveTime(device.Id, startTime, endTime)
		if err != nil {
			return 0, err
		}
		count += msgCount
	}
	return 0, nil
}

// SearchDeviceMsgCount 统计设备的消息总数（属性、事件都算在内）
func (pst *persistApp) SearchDeviceMsgCount(startTime, endTime int64) (int, error) {

	switch pst.dataDbClient.GetDataDBType() {
	case constants.LevelDB:
		return pst.searchDeviceMsgCountFromLevelDB(startTime, endTime)
	case constants.TDengine:
		return pst.searchDeviceMsgCountFromTDengine(startTime, endTime)
	}

	return 0, nil
}

func getEventTypeAndName(events []models.Events, code string) (string, string) {
	for _, event := range events {
		if event.Code == code {
			return event.EventType, event.Name
		}
	}
	return "", ""
}

func getServiceName(events []models.Actions, code string) string {
	for _, event := range events {
		if event.Code == code {
			return event.Name
		}
	}
	return ""
}
