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

package deviceapp

import (
	"context"
	"encoding/json"
	"github.com/docker/distribution/uuid"
	"github.com/winc-link/edge-driver-proto/thingmodel"
	"github.com/winc-link/hummingbird/internal/dtos"
	"github.com/winc-link/hummingbird/internal/hummingbird/core/container"
	"github.com/winc-link/hummingbird/internal/pkg/constants"
	"github.com/winc-link/hummingbird/internal/pkg/di"
	"github.com/winc-link/hummingbird/internal/tools/rpcclient"
	"time"
)

const (
	ResultSuccess = "success"
	ResultFail    = "fail"
)

// ui控制台、定时任务、场景联动、云平台api 都可以调用
func (p *deviceApp) DeviceAction(jobAction dtos.JobAction) dtos.DeviceExecRes {
	defer func() {
		if err := recover(); err != nil {
			p.lc.Error("CreateDeviceCallBack Panic:", err)
		}
	}()

	device, err := p.dbClient.DeviceById(jobAction.DeviceId)
	if err != nil {
		return dtos.DeviceExecRes{
			Result:  false,
			Message: "device not found",
		}
	}
	deviceService, err := p.dbClient.DeviceServiceById(device.DriveInstanceId)
	if err != nil {
		return dtos.DeviceExecRes{
			Result:  false,
			Message: "driver not found",
		}
	}

	driverService := container.DriverServiceAppFrom(di.GContainer.Get)
	status := driverService.GetState(deviceService.Id)
	if status == constants.RunStatusStarted {
		client, errX := rpcclient.NewDriverRpcClient(deviceService.BaseAddress, false, "", deviceService.Id, p.lc)
		if errX != nil {
			return dtos.DeviceExecRes{
				Result:  false,
				Message: errX.Error(),
			}
		}
		defer client.Close()
		var rpcRequest thingmodel.ThingModelIssueMsg
		rpcRequest.DeviceId = jobAction.DeviceId
		rpcRequest.OperationType = thingmodel.OperationType_PROPERTY_SET
		var data dtos.PropertySet
		data.Version = "v1.0"
		data.MsgId = uuid.Generate().String()
		data.Time = time.Now().UnixMilli()
		param := make(map[string]interface{})
		param[jobAction.Code] = jobAction.Value
		data.Params = param
		rpcRequest.Data = data.ToString()
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		_, err = client.ThingModelDownServiceClient.ThingModelMsgIssue(ctx, &rpcRequest)
		if err == nil {
			return dtos.DeviceExecRes{
				Result:  true,
				Message: ResultSuccess,
			}
		} else {
			return dtos.DeviceExecRes{
				Result:  false,
				Message: err.Error(),
			}
		}
	}
	return dtos.DeviceExecRes{
		Result:  false,
		Message: "driver status stop",
	}
}

func (p *deviceApp) DeviceInvokeThingService(invokeDeviceServiceReq dtos.InvokeDeviceServiceReq) dtos.DeviceExecRes {
	defer func() {
		if err := recover(); err != nil {
			p.lc.Error("CreateDeviceCallBack Panic:", err)
		}
	}()

	device, err := p.dbClient.DeviceById(invokeDeviceServiceReq.DeviceId)
	if err != nil {
		return dtos.DeviceExecRes{
			Result:  false,
			Message: "device not found",
		}
	}

	product, err := p.dbClient.ProductById(device.ProductId)
	if err != nil {
		return dtos.DeviceExecRes{
			Result:  false,
			Message: "product not found",
		}
	}
	var find bool
	var callType constants.CallType
	for _, action := range product.Actions {
		if action.Code == invokeDeviceServiceReq.Code {
			find = true
			callType = action.CallType
		}
	}

	if !find {
		return dtos.DeviceExecRes{
			Result:  false,
			Message: "code not found",
		}
	}

	deviceService, err := p.dbClient.DeviceServiceById(device.DriveInstanceId)
	if err != nil {
		return dtos.DeviceExecRes{
			Result:  false,
			Message: "driver not found",
		}
	}

	driverService := container.DriverServiceAppFrom(di.GContainer.Get)
	status := driverService.GetState(deviceService.Id)
	if status == constants.RunStatusStarted {
		client, errX := rpcclient.NewDriverRpcClient(deviceService.BaseAddress, false, "", deviceService.Id, p.lc)
		if errX != nil {
			return dtos.DeviceExecRes{
				Result:  false,
				Message: errX.Error(),
			}
		}
		defer client.Close()
		var rpcRequest thingmodel.ThingModelIssueMsg
		rpcRequest.DeviceId = invokeDeviceServiceReq.DeviceId
		rpcRequest.OperationType = thingmodel.OperationType_SERVICE_EXECUTE
		var data dtos.InvokeDeviceService
		data.Version = "v1.0"
		data.MsgId = uuid.Generate().String()
		data.Time = time.Now().UnixMilli()
		data.Data.Code = invokeDeviceServiceReq.Code
		data.Data.InputParams = invokeDeviceServiceReq.Items
		rpcRequest.Data = data.ToString()

		if callType == constants.CallTypeAsync {
			//saveServiceInfo := genSaveServiceInfo(data.MsgId, data.Time, invokeDeviceServiceReq)
			var saveServiceInfo dtos.ThingModelMessage
			saveServiceInfo.OpType = int32(thingmodel.OperationType_SERVICE_EXECUTE)
			saveServiceInfo.Cid = device.Id
			var saveData dtos.SaveServiceIssueData
			saveData.MsgId = data.MsgId
			saveData.Code = invokeDeviceServiceReq.Code
			saveData.Time = data.Time
			saveData.InputParams = invokeDeviceServiceReq.Items
			saveData.OutputParams = map[string]interface{}{
				"result":  true,
				"message": "success",
			}
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()
			_, err = client.ThingModelDownServiceClient.ThingModelMsgIssue(ctx, &rpcRequest)
			if err == nil {
				persistItf := container.PersistItfFrom(p.dic.Get)
				_ = persistItf.SaveDeviceThingModelData(saveServiceInfo)
				return dtos.DeviceExecRes{
					Result:  true,
					Message: ResultSuccess,
				}
			} else {
				return dtos.DeviceExecRes{
					Result:  false,
					Message: err.Error(),
				}
			}
		} else if callType == constants.CallTypeSync {
			var saveServiceInfo dtos.ThingModelMessage
			saveServiceInfo.OpType = int32(thingmodel.OperationType_SERVICE_EXECUTE)
			saveServiceInfo.Cid = device.Id
			var saveData dtos.SaveServiceIssueData
			saveData.MsgId = data.MsgId
			saveData.Code = invokeDeviceServiceReq.Code
			saveData.Time = data.Time
			saveData.InputParams = invokeDeviceServiceReq.Items
			persistItf := container.PersistItfFrom(p.dic.Get)

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()
			_, err = client.ThingModelDownServiceClient.ThingModelMsgIssue(ctx, &rpcRequest)

			if err != nil {
				return dtos.DeviceExecRes{
					Result:  false,
					Message: err.Error(),
				}
			}

			messageStore := container.MessageStoreItfFrom(p.dic.Get)
			ch := messageStore.GenAckChan(data.MsgId)

			select {
			case <-time.After(5 * time.Second):
				ch.TryCloseChan()
				saveData.OutputParams = map[string]interface{}{
					"result":  false,
					"message": "wait response timeout",
				}
				s, _ := json.Marshal(saveData)
				saveServiceInfo.Data = string(s)
				_ = persistItf.SaveDeviceThingModelData(saveServiceInfo)
				return dtos.DeviceExecRes{
					Result:  false,
					Message: "wait response timeout",
				}
			case <-ctx.Done():
				saveData.OutputParams = map[string]interface{}{
					"result":  false,
					"message": "wait response timeout",
				}
				s, _ := json.Marshal(saveData)
				saveServiceInfo.Data = string(s)
				_ = persistItf.SaveDeviceThingModelData(saveServiceInfo)
				return dtos.DeviceExecRes{
					Result:  false,
					Message: "wait response timeout",
				}
			case resp := <-ch.DataChan:
				if v, ok := resp.(map[string]interface{}); ok {
					saveData.OutputParams = v
					s, _ := json.Marshal(saveData)
					saveServiceInfo.Data = string(s)
					_ = persistItf.SaveDeviceThingModelData(saveServiceInfo)
					message, _ := json.Marshal(v)
					return dtos.DeviceExecRes{
						Result:  true,
						Message: string(message),
					}
				}
			}

		}
	}
	return dtos.DeviceExecRes{
		Result:  false,
		Message: "driver status stop",
	}
}
