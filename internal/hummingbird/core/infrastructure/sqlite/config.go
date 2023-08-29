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

package sqlite

import (
	"github.com/winc-link/hummingbird/internal/models"
	"github.com/winc-link/hummingbird/internal/pkg/constants"
	"github.com/winc-link/hummingbird/internal/pkg/errort"
	"gorm.io/gorm"
)

func updateAdvanceConfig(c *Client, config models.AdvanceConfig) error {
	if err := c.client.UpdateObject(&config); err != nil {
		return errort.NewCommonErr(errort.DefaultSystemError, err)
	}
	return nil
}

func getAdvanceConfig(c *Client) (models.AdvanceConfig, error) {
	var config models.AdvanceConfig
	err := c.client.GetObject(&models.AdvanceConfig{ID: constants.DefaultAdvanceConfigID}, &config)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			config.ID = constants.DefaultAdvanceConfigID
			if err = c.client.CreateObject(&config); err != nil {
				return models.AdvanceConfig{}, errort.NewCommonErr(errort.DefaultSystemError, err)
			}
			return config, nil
		}
		return models.AdvanceConfig{}, errort.NewCommonErr(errort.DefaultSystemError, err)
	}
	return config, nil
}
