//
// Copyright (C) 2020 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package httphelper

import (
	"context"
	"encoding/json"
	"github.com/winc-link/hummingbird/internal/pkg/errort"
	"net/http"
	"net/url"
)

// Helper method to make the get request and return the body
func GetRequest(ctx context.Context, returnValuePointer interface{}, baseUrl string, requestPath string, requestParams url.Values) error {
	req, err := createRequest(ctx, http.MethodGet, baseUrl, requestPath, requestParams)
	if err != nil {
		return err
	}

	res, err := sendRequest(ctx, req)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(res, returnValuePointer); err != nil {
		return errort.NewCommonEdgeX(errort.DefaultJsonParseError, "failed to parse the response body", err)
	}
	return nil
}

// Helper method to make the post JSON request and return the body
func PostRequest(
	ctx context.Context,
	returnValuePointer interface{},
	url string,
	data interface{}) error {

	req, err := createRequestWithRawData(ctx, http.MethodPost, url, data)
	if err != nil {
		return err
	}

	res, err := sendRequest(ctx, req)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(res, returnValuePointer); err != nil {
		return errort.NewCommonEdgeX(errort.DefaultJsonParseError, "failed to parse the response body", err)
	}
	return nil
}

func PutRequest(
	ctx context.Context,
	returnValuePointer interface{},
	url string,
	data interface{}) error {

	req, err := createRequestWithRawData(ctx, http.MethodPut, url, data)
	if err != nil {
		return err
	}

	res, err := sendRequest(ctx, req)
	if err != nil {
		return err
	}
	returnValuePointer = string(res)
	return nil
}

// Helper method to make the delete request and return the body
func DeleteRequest(ctx context.Context, returnValuePointer interface{}, baseUrl string, requestPath string, params url.Values) error {
	req, err := createRequest(ctx, http.MethodDelete, baseUrl, requestPath, params)
	if err != nil {
		return err
	}

	res, err := sendRequest(ctx, req)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(res, returnValuePointer); err != nil {
		return errort.NewCommonEdgeX(errort.DefaultJsonParseError, "failed to parse the response body", err)
	}
	return nil
}
