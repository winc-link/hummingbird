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

type (
	SceneLog struct {
		Timestamps `gorm:"embedded"`
		Id         string `gorm:"id;primaryKey;not null;type:string;size:255;comment:主键"`
		SceneId    string `gorm:"index;type:string;size:255;comment:场景ID" json:"scene_id"`
		Name       string `json:"name" gorm:"type:string;size:255;comment:名字"`
		ExecRes    string `json:"exec_res" gorm:"type:text;comment:执行结果"`
	}
)

func (pj *SceneLog) TableName() string {
	return "scene_log"
}

func (pj *SceneLog) Get() interface{} {
	return *pj
}
