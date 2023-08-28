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

package agentclient

import (
	"context"
	"github.com/winc-link/hummingbird/internal/dtos"
)

type AgentClient interface {
	GetGps(ctx context.Context) (dtos.Gps, error)
	AddServiceMonitor(ctx context.Context, stats dtos.ServiceStats) error
	DeleteServiceMonitor(ctx context.Context, serviceName string) error
	GetAllDriverMonitor(ctx context.Context) ([]dtos.ServiceStats, error)
	RestartGateway(ctx context.Context) error
	OperationService(ctx context.Context, op dtos.Operation) error
	GetAllAppServiceMonitor(ctx context.Context) ([]dtos.ServiceStats, error)
	GetAllServices(ctx context.Context) (res dtos.ServicesStats, err error)
	Exec(ctx context.Context, req dtos.AgentRequest) (res dtos.AgentResponse, err error)
	ResetGateway(ctx context.Context) error
}
