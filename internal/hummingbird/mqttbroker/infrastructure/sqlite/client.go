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
package sqlite

import (
	"github.com/winc-link/hummingbird/internal/dtos"
	"github.com/winc-link/hummingbird/internal/models"
	"github.com/winc-link/hummingbird/internal/pkg/errort"
	clientSQLite "github.com/winc-link/hummingbird/internal/tools/sqldb/sqlite"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Client struct {
	Pool          *gorm.DB
	client        clientSQLite.ClientSQLite
	loggingClient *zap.Logger
}

func NewClient(config dtos.Configuration, lc *zap.Logger) (c *Client, errEdgeX error) {
	client, err := clientSQLite.NewGormClient(config, nil)
	if err != nil {
		errEdgeX = errort.NewCommonEdgeX(errort.DefaultSystemError, "database failed to init", err)
		return
	}
	client.Pool = client.Pool.Debug()
	// 自动建表
	if err = client.InitTable(
		&models.MqttAuth{},
	); err != nil {
		errEdgeX = errort.NewCommonEdgeX(errort.DefaultSystemError, "database failed to init", err)
		return
	}
	c = &Client{
		client:        client,
		loggingClient: lc,
		Pool:          client.Pool,
	}
	return
}

func (client *Client) CloseSession() {
	return
}

func (client *Client) GetMqttAutInfo(clientId string) (models.MqttAuth, error) {
	return getMqttAutInfo(client, clientId)
}
