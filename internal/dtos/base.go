//
// Copyright (C) 2020 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package dtos

import (
	"github.com/google/uuid"
	"github.com/winc-link/hummingbird/internal/pkg/constants"
	//"gitlab.com/tedge/edgex/internal/pkg/constants"
)

// Request defines the base content for request DTOs (data transfer objects).
// This object and its properties correspond to the BaseRequest object in the APIv2 specification:
// https://app.swaggerhub.com/apis-docs/EdgeXFoundry1/core-data/2.x#/BaseRequest
type BaseRequest struct {
	Versionable `json:",inline"`
	RequestId   string `json:"requestId" validate:"len=0|uuid"`
}

func NewBaseRequest() BaseRequest {
	return BaseRequest{
		Versionable: NewVersionable(),
		RequestId:   uuid.NewString(),
	}
}

// BaseResponse defines the base content for response DTOs (data transfer objects).
// This object and its properties correspond to the BaseResponse object in the APIv2 specification:
// https://app.swaggerhub.com/apis-docs/EdgeXFoundry1/core-data/2.x#/BaseResponse
type BaseResponse struct {
	Versionable `json:",inline"`
	RequestId   string      `json:"requestId"`
	Message     interface{} `json:"message,omitempty"`
	StatusCode  int         `json:"statusCode"`
}

// Versionable shows the API version in DTOs
type Versionable struct {
	ApiVersion string `json:"apiVersion"`
}

// BaseWithIdResponse defines the base content for response DTOs (data transfer objects).
// This object and its properties correspond to the BaseWithIdResponse object in the APIv2 specification:
// https://app.swaggerhub.com/apis-docs/EdgeXFoundry1/core-data/2.x#/BaseWithIdResponse
type BaseWithIdResponse struct {
	BaseResponse `json:",inline"`
	Id           string `json:"id"`
}

func NewBaseResponse(requestId string, message string, statusCode int) BaseResponse {
	return BaseResponse{
		Versionable: NewVersionable(),
		RequestId:   requestId,
		Message:     message,
		StatusCode:  statusCode,
	}
}

func NewVersionable() Versionable {
	return Versionable{ApiVersion: constants.ApiVersion}
}

func NewBaseWithIdResponse(requestId string, message string, statusCode int, id string) BaseWithIdResponse {
	return BaseWithIdResponse{
		BaseResponse: NewBaseResponse(requestId, message, statusCode),
		Id:           id,
	}
}
