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

package ruleengine

import (
	"context"
	"github.com/winc-link/hummingbird/internal/dtos"
	resourceContainer "github.com/winc-link/hummingbird/internal/hummingbird/core/container"
	"github.com/winc-link/hummingbird/internal/pkg/constants"
	"time"
)

func (p ruleEngineApp) monitor() {
	tickTime := time.Second * 5
	timeTickerChan := time.Tick(tickTime)
	for {
		select {
		case <-timeTickerChan:
			p.checkRuleStatus()
		}
	}
}

func (p ruleEngineApp) checkRuleStatus() {
	ruleEngines, _, err := p.dbClient.RuleEngineSearch(0, -1, dtos.RuleEngineSearchQueryRequest{})
	if err != nil {
		p.lc.Errorf("get engines err:", err)
	}
	ekuiperApp := resourceContainer.EkuiperAppFrom(p.dic.Get)
	for _, ruleEngine := range ruleEngines {
		resp, err := ekuiperApp.GetRuleStats(context.Background(), ruleEngine.Id)
		if err != nil {
			p.lc.Errorf("error:", err)
			continue
		}
		status, ok := resp["status"]
		if ok {
			if status != string(ruleEngine.Status) {
				if status == string(constants.RuleEngineStop) {
					p.dbClient.RuleEngineStop(ruleEngine.Id)
				} else if status == string(constants.RuleEngineStart) {
					p.dbClient.RuleEngineStart(ruleEngine.Id)
				}
			}
		}
	}
}
