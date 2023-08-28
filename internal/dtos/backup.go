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

package dtos

const (
	BackupFileTypeDbResource = "db_resource"
	BackupFileTypeDbExpert   = "db_expert"
	BackupFileTypeDbGateway  = "db_gateway"
	BackupFileTypeCheck      = "check.json"
	BackupUnZipDir           = "/tmp/edge-recover"
)

// 备份/恢复时的校验文件
type BackupFileCheck struct {
	GatewayId string `json:"gateway_id"`
	Version   string `json:"version"`
}

type BackupCommand struct {
	BackupType int
}
