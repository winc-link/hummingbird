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

package models

import "github.com/winc-link/hummingbird/internal/pkg/constants"

type AdvanceConfig struct {
	ID             int                `json:"id" gorm:"id;primaryKey;not null;type:string;size:255;comment:主键"`
	LogLevel       constants.LogLevel `gorm:"default:1;comment:日志等级"` // 日志级别 默认为INFO
	PersistStorage bool               `gorm:"default:0;comment:存储开关"`
	StorageHour    int32              `gorm:"default:24;comment:存储时长"`
}

func (table *AdvanceConfig) TableName() string {
	return "advance_config"
}

func (table *AdvanceConfig) Get() interface{} {
	return *table
}
