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
	"github.com/winc-link/hummingbird/internal/dtos"
	"github.com/winc-link/hummingbird/internal/models"
)

func updateSystemMetrics(c *Client, metrics dtos.SystemMetrics) error {
	var m = models.SystemMetrics{
		Data:      metrics.String(),
		Timestamp: metrics.Timestamp,
	}
	return c.client.CreateObject(&m)
}

func getSystemMetrics(c *Client, start, end int64) ([]dtos.SystemMetrics, error) {
	var list []models.SystemMetrics
	if err := c.Pool.Where("timestamp >= ? and timestamp <= ?", start, end).Find(&list).Error; err != nil {
		return nil, err
	}
	var metrics = make([]dtos.SystemMetrics, 0)
	for _, item := range list {
		m, err := dtos.FromModelsSystemMetricsToDTO(item)
		if err != nil {
			return nil, err
		}
		metrics = append(metrics, m)
	}
	return metrics, nil
}

func removeRangeSystemMetrics(c *Client, min, max string) error {
	return c.Pool.Where("timestamp >= ? and timestamp <= ?", min, max).Delete(&models.SystemMetrics{}).Error
}
