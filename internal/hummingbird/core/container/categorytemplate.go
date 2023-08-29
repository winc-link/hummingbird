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
package container

import (
	interfaces "github.com/winc-link/hummingbird/internal/hummingbird/core/interface"
	"github.com/winc-link/hummingbird/internal/pkg/di"
)

var (
	CategoryTemplateAppName = di.TypeInstanceToName((*interfaces.CategoryApp)(nil))
)

func CategoryTemplateAppFrom(get di.Get) interfaces.CategoryApp {
	return get(CategoryTemplateAppName).(interfaces.CategoryApp)
}

var (
	ThingModelTemplateAppName = di.TypeInstanceToName((*interfaces.ThingModelTemplateApp)(nil))
)

func ThingModelTemplateAppFrom(get di.Get) interfaces.ThingModelTemplateApp {
	return get(ThingModelTemplateAppName).(interfaces.ThingModelTemplateApp)
}

var (
	UnitTemplateAppName = di.TypeInstanceToName((*interfaces.UnitApp)(nil))
)

func UnitTemplateAppFrom(get di.Get) interfaces.UnitApp {
	return get(UnitTemplateAppName).(interfaces.UnitApp)
}

var (
	DocsAppName = di.TypeInstanceToName((*interfaces.DocsApp)(nil))
)

func DocsTemplateAppFrom(get di.Get) interfaces.DocsApp {
	return get(DocsAppName).(interfaces.DocsApp)
}

var (
	QuickNavigationAppName = di.TypeInstanceToName((*interfaces.QuickNavigation)(nil))
)

func QuickNavigationAppTemplateAppFrom(get di.Get) interfaces.QuickNavigation {
	return get(QuickNavigationAppName).(interfaces.QuickNavigation)
}
