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
package mqttbroker

import (
	"github.com/winc-link/hummingbird/internal/dtos"
	"github.com/winc-link/hummingbird/internal/hummingbird/mqttbroker/config"
	"github.com/winc-link/hummingbird/internal/hummingbird/mqttbroker/infrastructure/sqlite"
	"github.com/winc-link/hummingbird/internal/hummingbird/mqttbroker/interfaces"
	"go.uber.org/zap"
)

var (
	DbClient interfaces.DBClient
)

func GetDbClient() interfaces.DBClient {
	return DbClient
}

type Database struct {
	conf config.Config
}

// NewDatabase is a factory method that returns an initialized Database receiver struct.
func NewDatabase(conf config.Config) Database {
	return Database{
		conf: conf,
	}
}

// init the dbClient interfaces
func (d Database) InitDBClient(
	lc *zap.Logger) error {
	dbClient, err := sqlite.NewClient(dtos.Configuration{
		Cluster:      d.conf.Database.Cluster,
		Username:     d.conf.Database.Username,
		Password:     d.conf.Database.Password,
		DataSource:   d.conf.Database.DataSource,
		DatabaseName: d.conf.Database.Name,
	}, lc)
	if err != nil {
		return err
	}
	DbClient = dbClient
	return nil
}
