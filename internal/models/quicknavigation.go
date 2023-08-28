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

type QuickNavigation struct {
	Timestamps `gorm:"embedded"`
	Id         string `gorm:"id;primaryKey"`
	Name       string `json:"name"`
	Icon       string `json:"icon"`
	Sort       int    `json:"sort"`
	JumpLink   string `json:"jump_link"`
}

func (d *QuickNavigation) TableName() string {
	return "quick_navigation"
}

func (d *QuickNavigation) Get() interface{} {
	return *d
}
