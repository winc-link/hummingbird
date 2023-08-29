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

package interfaces

import (
	"context"
	"github.com/winc-link/hummingbird/internal/dtos"
	//"github.com/winc-link/hummingbird/internal/dtos"
)

type AlertRuleApp interface {
	AddAlertRule(ctx context.Context, req dtos.RuleAddRequest) (string, error)
	UpdateAlertRule(ctx context.Context, req dtos.RuleUpdateRequest) error
	UpdateAlertField(ctx context.Context, req dtos.RuleFieldUpdate) error
	AlertRuleById(ctx context.Context, id string) (dtos.RuleResponse, error)
	AlertRulesSearch(ctx context.Context, req dtos.AlertRuleSearchQueryRequest) ([]dtos.AlertRuleSearchQueryResponse, uint32, error)
	AlertRulesDelete(ctx context.Context, id string) error
	AlertRulesStop(ctx context.Context, id string) error
	AlertRulesStart(ctx context.Context, id string) error
	AlertRulesRestart(ctx context.Context, id string) error
	AlertIgnore(ctx context.Context, id string) error
	TreatedIgnore(ctx context.Context, id, message string) error
	AlertPlate(ctx context.Context, beforeTime int64) ([]dtos.AlertPlateQueryResponse, error)
	AlertSearch(ctx context.Context, req dtos.AlertSearchQueryRequest) ([]dtos.AlertSearchQueryResponse, uint32, error)
	AddAlert(ctx context.Context, req map[string]interface{}) error
	CheckRuleByProductId(ctx context.Context, productId string) error
	CheckRuleByDeviceId(ctx context.Context, deviceId string) error
}

type RuleEngineApp interface {
	AddRuleEngine(ctx context.Context, req dtos.RuleEngineRequest) (string, error)
	UpdateRuleEngine(ctx context.Context, req dtos.RuleEngineUpdateRequest) error
	UpdateRuleEngineField(ctx context.Context, req dtos.RuleEngineFieldUpdateRequest) error
	RuleEngineById(ctx context.Context, id string) (dtos.RuleEngineResponse, error)
	RuleEngineSearch(ctx context.Context, req dtos.RuleEngineSearchQueryRequest) ([]dtos.RuleEngineSearchQueryResponse, uint32, error)
	RuleEngineDelete(ctx context.Context, id string) error
	RuleEngineStop(ctx context.Context, id string) error
	RuleEngineStart(ctx context.Context, id string) error
	RuleEngineStatus(ctx context.Context, id string) (map[string]interface{}, error)
}
