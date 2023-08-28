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

package ekuiperclient

import (
	"context"
	"github.com/winc-link/hummingbird/internal/dtos"
)

type EkuiperClient interface {
	RuleExist(ctx context.Context, ruleId string) (bool, error)
	CreateRule(ctx context.Context, actions []dtos.Actions, ruleId string, sql string) error
	UpdateRule(ctx context.Context, actions []dtos.Actions, ruleId string, sql string) error
	GetRuleStats(ctx context.Context, ruleId string) (map[string]interface{}, error)
	StartRule(ctx context.Context, ruleId string) error
	StopRule(ctx context.Context, ruleId string) error
	RestartRule(ctx context.Context, ruleId string) error
	DeleteRule(ctx context.Context, ruleId string) error
}
