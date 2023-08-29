/*******************************************************************************
 * Copyright 2023 Winc link Inc.
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
package main

import (
	"context"
	"github.com/winc-link/hummingbird/internal/hummingbird/core"

	"github.com/gin-gonic/gin"
)

// @title 赢创万联（蜂鸟） API
// @version 1.0
// @description Swagger API for Golang Project hummingbird.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email email@winc-link.com

// @BasePath /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name x-token
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	gin.SetMode(gin.ReleaseMode)
	core.Main(ctx, cancel, gin.Default())
}
