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
package driverapp

import (
	"context"
	"github.com/winc-link/hummingbird/internal/hummingbird/core/container"
	interfaces "github.com/winc-link/hummingbird/internal/hummingbird/core/interface"
	pkgcontainer "github.com/winc-link/hummingbird/internal/pkg/container"
	"github.com/winc-link/hummingbird/internal/pkg/di"
	"github.com/winc-link/hummingbird/internal/pkg/logger"
)

type driverLibApp struct {
	dic      *di.Container
	lc       logger.LoggingClient
	dbClient interfaces.DBClient

	manager DeviceLibraryManager
	market  *driverMarket // Refactor: interface
}

func NewDriverApp(ctx context.Context, dic *di.Container) interfaces.DriverLibApp {
	return newDriverLibApp(dic)
}

func newDriverLibApp(dic *di.Container) *driverLibApp {
	app := &driverLibApp{
		dic:      dic,
		lc:       pkgcontainer.LoggingClientFrom(dic.Get),
		dbClient: container.DBClientFrom(dic.Get),
	}
	app.manager = newDriverLibManager(dic, app)
	app.market = newDriverMarket(dic, app)
	return app
}

func (app *driverLibApp) getDriverServiceApp() interfaces.DriverServiceApp {
	return container.DriverServiceAppFrom(app.dic.Get)
}
