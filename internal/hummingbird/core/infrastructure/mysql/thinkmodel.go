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

package mysql

import (
	"fmt"
	"github.com/winc-link/hummingbird/internal/models"
)

func abilityByCode(c *Client, model interface{}, code, productId string) (interface{}, error) {
	var err error
	switch model.(type) {
	case models.Properties:
		ability := models.Properties{}
		err = c.Pool.Model(&models.Properties{}).Where("code = ? and product_id = ?", code, productId).Find(&ability).Error
		if err != nil {
			return nil, err
		} else {
			return ability, nil
		}
	case models.Events:
		ability := models.Events{}
		err = c.Pool.Model(&models.Events{}).Where("code = ? and product_id = ?", code, productId).Find(&ability).Error
		if err != nil {
			return nil, err
		} else {
			return ability, nil
		}
	default:
		return nil, fmt.Errorf("ability type shoud be propery or event")
	}
}
