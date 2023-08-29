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

package container

import (
	interfaces "github.com/winc-link/hummingbird/internal/hummingbird/core/interface"
	"github.com/winc-link/hummingbird/internal/pkg/di"
)

var (
	SceneAppName = di.TypeInstanceToName((*interfaces.SceneApp)(nil))
)

func SceneAppNameFrom(get di.Get) interfaces.SceneApp {
	return get(SceneAppName).(interfaces.SceneApp)
}

var (
	ConJobAppName = di.TypeInstanceToName((*interfaces.ConJob)(nil))
)

func ConJobAppNameFrom(get di.Get) interfaces.ConJob {
	return get(ConJobAppName).(interfaces.ConJob)
}
