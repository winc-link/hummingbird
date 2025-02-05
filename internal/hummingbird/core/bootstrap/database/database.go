//
// Copyright (C) 2020 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package database

import (
	"context"
	"errors"
	"github.com/winc-link/hummingbird/internal/dtos"
	"github.com/winc-link/hummingbird/internal/hummingbird/core/config"
	"github.com/winc-link/hummingbird/internal/hummingbird/core/container"
	"github.com/winc-link/hummingbird/internal/hummingbird/core/infrastructure/mysql"
	"github.com/winc-link/hummingbird/internal/hummingbird/core/infrastructure/sqlite"
	"github.com/winc-link/hummingbird/internal/pkg/constants"
	"github.com/winc-link/hummingbird/internal/tools/datadb/tdengine"

	"github.com/winc-link/hummingbird/internal/pkg/di"
	"github.com/winc-link/hummingbird/internal/pkg/logger"
	"github.com/winc-link/hummingbird/internal/pkg/startup"
	"github.com/winc-link/hummingbird/internal/tools/datadb/leveldb"
	"sync"

	interfaces "github.com/winc-link/hummingbird/internal/hummingbird/core/interface"
	pkgContainer "github.com/winc-link/hummingbird/internal/pkg/container"
)

// Database contains references to dependencies required by the database bootstrap implementation.
type Database struct {
	database *config.ConfigurationStruct
}

// NewDatabase is a factory method that returns an initialized Database receiver struct.
func NewDatabase(database *config.ConfigurationStruct) Database {
	return Database{
		database: database,
	}
}

// Return the dbClient interfaces
func (d Database) newDBClient(
	lc logger.LoggingClient) (interfaces.DBClient, error) {

	databaseInfo := d.database.GetDatabaseInfo()["Primary"]
	switch databaseInfo.Type {
	case string(constants.MySQL):
		return mysql.NewClient(dtos.Configuration{
			Dsn: databaseInfo.Dsn,
		}, lc)
	case string(constants.SQLite):
		return sqlite.NewClient(dtos.Configuration{
			Username:   databaseInfo.Username,
			Password:   databaseInfo.Password,
			DataSource: databaseInfo.DataSource,
		}, lc)
	default:
		panic(errors.New("database configuration error"))
	}
}

func (d Database) newDataDBClient(
	lc logger.LoggingClient) (interfaces.DataDBClient, error) {
	dataDbInfo := d.database.GetDataDatabaseInfo()["Primary"]

	switch dataDbInfo.Type {
	case string(constants.LevelDB):
		return leveldb.NewClient(dtos.Configuration{
			DataSource: dataDbInfo.DataSource,
		}, lc)
	case string(constants.TDengine):
		return tdengine.NewClient(dtos.Configuration{
			Dsn: dataDbInfo.Dsn,
		}, lc)
	default:
		panic(errors.New("database configuration error"))

	}
}

// BootstrapHandler fulfills the BootstrapHandler contract and initializes the database.
func (d Database) BootstrapHandler(
	ctx context.Context,
	wg *sync.WaitGroup,
	startupTimer startup.Timer,
	dic *di.Container) bool {
	lc := pkgContainer.LoggingClientFrom(dic.Get)

	// initialize Metadata db.
	dbClient, err := d.newDBClient(lc)
	if err != nil {
		panic(err)
	}

	dic.Update(di.ServiceConstructorMap{
		container.DBClientInterfaceName: func(get di.Get) interface{} {
			return dbClient
		},
	})

	// initialize Data db.
	dataDbClient, err := d.newDataDBClient(lc)
	if err != nil {
		panic(err)
	}

	dic.Update(di.ServiceConstructorMap{
		container.DataDBClientInterfaceName: func(get di.Get) interface{} {
			return dataDbClient
		},
	})

	lc.Info("DatabaseInfo connected")

	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		case <-ctx.Done():
			interfaces.DMIFrom(di.GContainer.Get).StopAllInstance() //stop all instance
			container.DBClientFrom(di.GContainer.Get).CloseSession()
			container.DataDBClientFrom(di.GContainer.Get).CloseSession()
			lc.Info("DatabaseInfo disconnected")
		}
	}()
	return true
}
