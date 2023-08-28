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

type LanguageSdk struct {
	Timestamps  `gorm:"embedded"`
	Id          string `gorm:"id;primaryKey;not null;type:string;size:255;comment:主键"`
	Name        string `json:"name" gorm:"type:string;size:255;comment:名字"`
	Icon        string `json:"icon" gorm:"type:text;comment:图标"`
	Addr        string `json:"addr" gorm:"type:string;size:255;comment:地址"`
	Description string `json:"description" gorm:"type:text;comment:描述"`
	Sort        int    `gorm:"type:int;size:8;comment:排序"`
}

func (d *LanguageSdk) TableName() string {
	return "language_sdk"
}

func (d *LanguageSdk) Get() interface{} {
	return *d
}
