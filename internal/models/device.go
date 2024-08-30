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
	"crypto/hmac"
	"crypto/sha1"
	"database/sql/driver"
	"encoding/hex"
	"github.com/winc-link/edge-driver-proto/driverdevice"
	"github.com/winc-link/hummingbird/internal/pkg/constants"
	"strings"
)

type Device struct {
	Timestamps      `gorm:"embedded"`
	Id              string                 `gorm:"id;primaryKey;not null;type:string;size:255;comment:主键"`
	CloudDeviceId   string                 `gorm:"type:string;size:255;comment:云设备ID"`
	CloudProductId  string                 `gorm:"type:string;size:255;comment:云产品ID"`
	CloudInstanceId string                 `gorm:"index;type:string;size:255;comment:云实例ID"`
	DriveInstanceId string                 `gorm:"index;type:string;size:255;comment:驱动实例ID"`
	Name            string                 `gorm:"type:string;size:255;comment:名字"`
	DeviceSn        string                 `gorm:"type:string;size:255;comment:设备唯一编码"`
	Status          constants.DeviceStatus `gorm:"type:string;size:50;comment:设备状态"`
	Description     string                 `gorm:"type:text;comment:描述"`
	ProductId       string                 `gorm:"type:string;size:255;comment:产品ID"`
	Secret          string                 `gorm:"type:string;size:255;comment:密钥"`
	Platform        constants.IotPlatform  `gorm:"type:string;size:50;comment:平台名称"`
	InstallLocation string                 `gorm:"type:string;size:255;comment:安装地址"`
	LastSyncTime    int64                  `gorm:"comment:最后一次同步时间"`
	LastOnlineTime  int64                  `gorm:"comment:最后一次在线时间"`
	Product         Product                `gorm:"foreignKey:ProductId"`
}

// ProtocolProperties contains the device connection information in key/value pair
type ProtocolProperties map[string]string

type Protocols map[string]ProtocolProperties

func (c Protocols) Value() (driver.Value, error) {
	return GormValueWrap(c)
}

func (c *Protocols) Scan(value interface{}) error {
	return GormScanWrap(value, c)
}

func (table *Device) TableName() string {
	return "device"
}

func (table *Device) Get() interface{} {
	return *table
}

func (table *Device) TransformToDriverDevice() *driverdevice.Device {
	driverDevice := new(driverdevice.Device)
	driverDevice.Id = table.Id
	driverDevice.Name = table.Name
	driverDevice.ProductId = table.ProductId
	driverDevice.Description = table.Description
	driverDevice.Status = table.Status.TransformToDriverDeviceStatus()
	driverDevice.ProductId = table.ProductId
	driverDevice.Secret = table.Secret
	driverDevice.Platform = table.Platform.TransformToDriverDevicePlatform()
	return driverDevice
}

func HmacSha1(keyStr, value string) string {
	key := []byte(keyStr)
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(value))
	res := hex.EncodeToString(mac.Sum(nil))
	return strings.ToUpper(res)
}
