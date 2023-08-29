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
	"github.com/winc-link/hummingbird/internal/models"
	"github.com/winc-link/hummingbird/internal/pkg/constants"
)

type DataDBClient interface {
	GetDataDBType() constants.DataType
	CloseSession()

	Insert(ctx context.Context, table string, data map[string]interface{}) (err error)
	GetDeviceProperty(req dtos.ThingModelPropertyDataRequest, device models.Device) ([]dtos.ReportData, int, error)
	GetDeviceService(req dtos.ThingModelServiceDataRequest, device models.Device, product models.Product) ([]dtos.SaveServiceIssueData, int, error)
	GetDeviceEvent(req dtos.ThingModelEventDataRequest, device models.Device, product models.Product) ([]dtos.EventData, int, error)

	CreateTable(ctx context.Context, stable, table string) (err error)
	DropTable(ctx context.Context, table string) (err error)

	CreateStable(ctx context.Context, product models.Product) (err error)
	DropStable(ctx context.Context, table string) (err error)

	AddDatabaseField(ctx context.Context, tableName string, specsType constants.SpecsType, code string, name string) (err error)
	DelDatabaseField(ctx context.Context, tableName, code string) (err error)
	ModifyDatabaseField(ctx context.Context, tableName string, specsType constants.SpecsType, code string, name string) (err error)

	GetDevicePropertyCount(dtos.ThingModelPropertyDataRequest) (int, error)
	GetDeviceEventCount(req dtos.ThingModelEventDataRequest) (int, error)
	GetDeviceMsgCountByGiveTime(deviceId string, startTime, endTime int64) (int, error)
}
