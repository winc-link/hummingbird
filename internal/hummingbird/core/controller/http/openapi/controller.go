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

package openapi

import (
	resourceContainer "github.com/winc-link/hummingbird/internal/hummingbird/core/container"
	interfaces "github.com/winc-link/hummingbird/internal/hummingbird/core/interface"
	"github.com/winc-link/hummingbird/internal/pkg/container"
	"github.com/winc-link/hummingbird/internal/pkg/di"
	"github.com/winc-link/hummingbird/internal/pkg/logger"
)

type controller struct {
	lc  logger.LoggingClient
	dic *di.Container
	//cfg   *config.ConfigurationStruct
	//gwApp *gatewayapp.GatewayApp
}

func New(dic *di.Container) *controller {
	return &controller{
		lc:  container.LoggingClientFrom(dic.Get),
		dic: dic,
		//cfg:   resourceContainer.ConfigurationFrom(dic.Get),
		//gwApp: gatewayapp.NewGatewayApp(dic),
	}
}

func (ctl *controller) getUserApp() interfaces.UserItf {
	return resourceContainer.UserItfFrom(ctl.dic.Get)
}

func (ctl *controller) getProductApp() interfaces.ProductItf {
	return resourceContainer.ProductAppNameFrom(ctl.dic.Get)
}

func (ctl *controller) getDeviceApp() interfaces.DeviceItf {
	return resourceContainer.DeviceItfFrom(ctl.dic.Get)
}

func (ctl *controller) getThingModelApp() interfaces.ThingModelCtlItf {
	return resourceContainer.ThingModelAppNameFrom(ctl.dic.Get)
}

func (ctl *controller) getPersistApp() interfaces.PersistItf {
	return resourceContainer.PersistItfFrom(ctl.dic.Get)
}
