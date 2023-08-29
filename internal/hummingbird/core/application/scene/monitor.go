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

package scene

import (
	"context"
	"github.com/winc-link/hummingbird/internal/dtos"
	resourceContainer "github.com/winc-link/hummingbird/internal/hummingbird/core/container"
	"github.com/winc-link/hummingbird/internal/pkg/constants"
	"time"
)

func (p sceneApp) monitor() {
	tickTime := time.Second * 9
	timeTickerChan := time.Tick(tickTime)
	for {
		select {
		case <-timeTickerChan:
			p.checkSceneRuleStatus()
		}
	}
}

func (p sceneApp) checkSceneRuleStatus() {
	scenes, _, err := p.dbClient.SceneSearch(0, -1, dtos.SceneSearchQueryRequest{})
	if err != nil {
		p.lc.Errorf("get engines err:", err)
	}
	ekuiperApp := resourceContainer.EkuiperAppFrom(p.dic.Get)
	for _, scene := range scenes {
		if len(scene.Conditions) != 1 {
			continue
		}
		if scene.Conditions[0].ConditionType != "notify" {
			continue
		}
		resp, err := ekuiperApp.GetRuleStats(context.Background(), scene.Id)
		if err != nil {
			p.lc.Errorf("error:", err)
			continue
		}
		status, ok := resp["status"]
		if ok {
			if status != string(scene.Status) {
				if status == string(constants.SceneStart) {
					p.dbClient.SceneStart(scene.Id)
				} else if status == string(constants.SceneStop) {
					p.dbClient.SceneStop(scene.Id)
				}
			}
		}
	}
}
