//
// Copyright (C) 2021 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package dtos

import "encoding/json"

/*
 * An Operation for SMA processing.
 *
 *
 * Operation struct
 */
type Operation struct {
	Action  string `json:"action,omitempty" binding:"oneof=start stop restart"` // 动作，重启 restart
	Service string `json:"service,omitempty" binding:"required"`                // 服务名称
}

// String returns a JSON encoded string representation of the model
func (o Operation) String() string {
	out, err := json.Marshal(o)
	if err != nil {
		return err.Error()
	}
	return string(out)
}
