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

package route

import (
	"github.com/gin-gonic/gin"
	"github.com/winc-link/hummingbird/internal/hummingbird/core/controller/http/openapi"
	"github.com/winc-link/hummingbird/internal/pkg/di"
	"github.com/winc-link/hummingbird/internal/pkg/errort"
	"github.com/winc-link/hummingbird/internal/tools/jwt"
	"github.com/winc-link/hummingbird/internal/tools/openapihelper"
)

func RegisterOpenApi(engine *gin.Engine, dic *di.Container) {
	engine.NoRoute(func(c *gin.Context) {
		openapihelper.ReaderFail(c, errort.UrlPathIsInvalid)
	})

	ctl := openapi.New(dic)

	v := engine.Group("/v1.0/openapi")
	v.POST("/auth/login", ctl.Login)
	v.POST("/token/:refreshToken", ctl.RefreshToken)
	v1 := v.Group("", jwt.JWTAuth(false))

	//产品管理的API
	{
		//创建产品
		v1.POST("/product", ctl.OpenApiCreateProduct)
		//更新产品
		v1.PUT("/product/:productId", ctl.OpenApiUpdateProduct)
		//查询产品列表
		v1.GET("/products", ctl.OpenApiProductSearch)
		//查询产品详细信息。
		v1.GET("/product/:productId", ctl.OpenApiProductById)
		// 发布产品
		v1.GET("/product-release/:productId", ctl.OpenApiProductReleaseById)
		// 取消发布产品
		v1.GET("/product-unrelease/:productId", ctl.OpenApiProductUnReleaseById)
		//删除指定产品。
		v1.DELETE("product/:productId", ctl.OpenApiDeleteProduct)
	}
	//设备管理的API
	{
		//创建设备
		v1.POST("/device", ctl.OpenApiCreateDevice)
		//更新设备
		v1.PUT("/device/:deviceId", ctl.OpenApiUpdateDevice)
		//查询设备列表
		v1.GET("/devices", ctl.OpenApiDeviceSearch)
		//查询设备详细信息。
		v1.GET("/device/:deviceId", ctl.OpenApiDeviceById)
		//删除指定设备。
		v1.DELETE("/device/:deviceId", ctl.OpenApiDeleteDevice)
		//获取设备的运行状态。
		//v1.GET("/deviceStatus/:deviceId", ctl.OpenApiDeviceStatus)

	}
	//物模型管理的API
	{
		//为指定产品的物模型新增功能
		v1.POST("/thingModel", ctl.OpenApiThingModelAddOrUpdate)
		//更新指定产品物模型中的单个功能
		v1.PUT("/thingModel", ctl.OpenApiThingModelAddOrUpdate)
		//查看指定产品的物模型中的功能定义详情
		v1.GET("/thingModel", ctl.OpenApiThingModel)
		//DeleteThingModel
		v1.DELETE("/thingModel", ctl.OpenApiDeleteThingModel)

	}
	//物模型使用的API
	{
		//查询设备实时属性数据。
		v1.GET("/queryDeviceEffectivePropertyData", ctl.OpenApiQueryDeviceEffectivePropertyData)
		//设置设备的属性。
		v1.POST("/setDeviceProperty", ctl.OpenApiSetDeviceProperty)
		//调用设备的服务。
		v1.POST("/invokeThingService", ctl.OpenApiInvokeThingService)
		//查询设备的属性历史数据。
		v1.GET("/queryDevicePropertyData", ctl.OpenApiQueryDevicePropertyData)
		//查询设备的事件历史数据。
		v1.GET("/queryDeviceEventData", ctl.OpenApiQueryDeviceEventData)
		//获取设备的服务记录历史数据。
		v1.GET("/queryDeviceServiceData", ctl.OpenApiQueryDeviceServiceData)
	}
}
