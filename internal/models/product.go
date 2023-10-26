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
package models

import (
	"github.com/winc-link/edge-driver-proto/driverproduct"
	"github.com/winc-link/hummingbird/internal/pkg/constants"
)

type Product struct {
	Timestamps      `gorm:"embedded"`
	Id              string                    `gorm:"id;primaryKey;not null;type:string;size:255;comment:主键"`
	Name            string                    `gorm:"type:string;size:255;comment:名字"`
	Key             string                    `gorm:"type:string;size:255;comment:产品标识"`
	CloudProductId  string                    `gorm:"type:string;size:255;comment:云产品ID"`
	CloudInstanceId string                    `gorm:"index;type:string;size:255;comment:云实例ID"`
	Platform        constants.IotPlatform     `gorm:"type:string;size:255;comment:平台"`
	Protocol        string                    `gorm:"type:string;size:255;comment:协议"`
	NodeType        constants.ProductNodeType `gorm:"type:string;size:255;comment:节点类型"`
	NetType         constants.ProductNetType  `gorm:"type:string;size:255;comment:网络类型"`
	DataFormat      string                    `gorm:"type:string;size:255;comment:数据类型"`
	LastSyncTime    int64                     `gorm:"comment:最后一次同步时间"`
	Factory         string                    `gorm:"type:string;size:255;comment:工厂名称"`
	Description     string                    `gorm:"type:text;comment:描述"`
	Status          constants.ProductStatus   `gorm:"type:string;size:255;comment:产品状态"`
	Extra           MapStringString           `gorm:"type:string;size:255;comment:扩展字段"`
	Properties      []Properties              `gorm:"foreignKey:ProductId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"` // 物模型的属性列表
	Events          []Events                  `gorm:"foreignKey:ProductId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"` // 物模型的事件列表
	Actions         []Actions                 `gorm:"foreignKey:ProductId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"` // 物模型的动作列表
}

func (d *Product) TableName() string {
	return "product"
}

func (d *Product) Get() interface{} {
	return *d
}

func (d *Product) TransformToDriverProduct() *driverproduct.Product {
	driverProduct := new(driverproduct.Product)
	driverProduct.Id = d.Id
	driverProduct.Name = d.Name
	driverProduct.Description = d.Description
	driverProduct.NodeType = d.NodeType.TransformToDriverProductNodeType()
	driverProduct.DataFormat = d.DataFormat
	driverProduct.Platform = d.Platform.TransformToDriverDevicePlatform()
	driverProduct.NetType = d.NetType.TransformToDriverProductNetType()
	driverProduct.ProtocolType = d.Protocol
	driverProduct.Key = d.Key
	driverProduct.CreateAt = uint64(d.Created)
	var p []*driverproduct.Properties
	var e []*driverproduct.Events
	var a []*driverproduct.Actions

	for _, property := range d.Properties {
		driverProperty := new(driverproduct.Properties)
		driverProperty.Name = property.Name
		driverProperty.ProductId = property.ProductId
		driverProperty.Code = property.Code
		driverProperty.Description = property.Description
		driverProperty.Required = property.Require
		driverProperty.AccessMode = property.AccessMode
		driverProperty.TypeSpec = new(driverproduct.TypeSpec)
		driverProperty.TypeSpec.Type = string(property.TypeSpec.Type)
		driverProperty.TypeSpec.Specs = property.TypeSpec.Specs
		p = append(p, driverProperty)
	}

	for _, event := range d.Events {
		driverEvent := new(driverproduct.Events)
		driverEvent.Name = event.Name
		driverEvent.ProductId = event.ProductId
		driverEvent.Code = event.Code
		driverEvent.Description = event.Description
		driverEvent.Required = event.Require
		driverEvent.Type = event.EventType
		var driverOutParams []*driverproduct.OutputParams
		for _, outparam := range event.OutputParams {
			driverOutParam := new(driverproduct.OutputParams)
			driverOutParam.Code = outparam.Code
			driverOutParam.Name = outparam.Name
			driverOutParam.TypeSpec = new(driverproduct.TypeSpec)
			driverOutParam.TypeSpec.Type = string(outparam.TypeSpec.Type)
			driverOutParam.TypeSpec.Specs = outparam.TypeSpec.Specs
			driverOutParams = append(driverOutParams, driverOutParam)
		}
		driverEvent.OutputParams = driverOutParams
		e = append(e, driverEvent)
	}

	for _, action := range d.Actions {
		driverAction := new(driverproduct.Actions)
		driverAction.Name = action.Name
		driverAction.ProductId = action.ProductId
		driverAction.Code = action.Code
		driverAction.Description = action.Description
		driverAction.Required = action.Require
		driverAction.CallType = string(action.CallType)

		var driverInParams []*driverproduct.InputParams
		for _, inparam := range action.InputParams {
			driverInParam := new(driverproduct.InputParams)
			driverInParam.Code = inparam.Code
			driverInParam.Name = inparam.Name
			driverInParam.TypeSpec = new(driverproduct.TypeSpec)
			driverInParam.TypeSpec.Type = string(inparam.TypeSpec.Type)
			driverInParam.TypeSpec.Specs = inparam.TypeSpec.Specs
			driverInParams = append(driverInParams, driverInParam)
		}
		driverAction.InputParams = driverInParams

		var driverOutParams []*driverproduct.OutputParams
		for _, outparam := range action.OutputParams {
			driverOutParam := new(driverproduct.OutputParams)
			driverOutParam.Code = outparam.Code
			driverOutParam.Name = outparam.Name
			driverOutParam.TypeSpec = new(driverproduct.TypeSpec)
			driverOutParam.TypeSpec.Type = string(outparam.TypeSpec.Type)
			driverOutParam.TypeSpec.Specs = outparam.TypeSpec.Specs
			driverOutParams = append(driverOutParams, driverOutParam)
		}
		driverAction.OutputParams = driverOutParams

		a = append(a, driverAction)
	}
	driverProduct.Properties = p
	driverProduct.Events = e
	driverProduct.Actions = a

	return driverProduct
}
