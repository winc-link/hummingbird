/*******************************************************************************
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

type Doc struct {
	Timestamps `gorm:"embedded"`
	Id         string `gorm:"id;primaryKey;not null;type:string;size:255;comment:主键"`
	Name       string `gorm:"type:string;size:255;comment:名字"`
	Sort       int    `gorm:"type:int;size:8;comment:排序"`
	JumpLink   string `gorm:"type:string;size:255;comment:跳转地址"`
}

func (d *Doc) TableName() string {
	return "doc"
}

func (d *Doc) Get() interface{} {
	return *d
}
