//
// Copyright (C) 2021 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package middleware

import (
	"context"
	"github.com/winc-link/hummingbird/internal/pkg/constants"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func FromContext(ctx context.Context) string {
	hdr, ok := ctx.Value(constants.CorrelationHeader).(string)
	if !ok {
		hdr = ""
	}
	return hdr
}

func NewId() string {
	return uuid.New().String()
}

func WithCorrelationId(ctx context.Context) context.Context {
	hdr, ok := ctx.Value(constants.CorrelationHeader).(string)
	if ok && hdr != "" {
		return ctx
	}
	hdr = NewId()

	return context.WithValue(ctx, constants.CorrelationHeader, hdr)
}

// CorrelationHeader
func CorrelationHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		hdr := c.GetHeader(constants.CorrelationHeader)
		if hdr == "" {
			hdr = uuid.New().String()
		}
		c.Header(constants.CorrelationHeader, hdr)
		c.Set(constants.CorrelationHeader, hdr)

		c.Next()
	}
}
