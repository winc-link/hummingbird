/*******************************************************************************
 * Copyright 2017 Dell Inc.
 * Copyright (c) 2019 Intel Corporation
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
	"github.com/winc-link/hummingbird/internal/hummingbird/core/controller/http/gateway"
	"github.com/winc-link/hummingbird/internal/hummingbird/core/controller/http/websocket"
	"github.com/winc-link/hummingbird/internal/pkg/di"
	"github.com/winc-link/hummingbird/internal/tools/jwt"

	"github.com/gin-gonic/gin"

	"github.com/swaggo/files"
	gs "github.com/swaggo/gin-swagger"
)

func RegisterGateway(engine *gin.Engine, dic *di.Container) {
	ctl := gateway.New(dic)
	v1 := engine.Group("/api/v1")
	v1.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))

	v1.POST("auth/login", ctl.Login)
	v1.GET("auth/initInfo", ctl.InitInfo)
	v1.POST("auth/init-password", ctl.InitPassword)

	v1.POST("ekuiper/alert", ctl.EkuiperAlert)
	v1.POST("ekuiper/scene", ctl.EkuiperScene) //ekuiper 服务调用
	v1.GET("ws/", websocket.NewServer(dic).Handle)

	v1Auth := v1.Group("", jwt.JWTAuth(false))
	v1Auth.PUT("auth/password", ctl.UpdatePassword)
	/*******首页 *******/
	{
		v1Auth.GET("home-page", ctl.HomePage)
	}
	{
		/******* 运维监控 *******/
		v1Auth.GET("/metrics/system", ctl.SystemMetricsHandler)
	}

	/******* 镜像仓库管理 *******/
	{
		v1Auth.POST("docker-configs", ctl.DockerConfigAdd)
		v1Auth.GET("docker-configs", ctl.DockerConfigsSearch)
		v1Auth.PUT("docker-configs/:dockerConfigId", ctl.DockerConfigUpdate)
		v1Auth.DELETE("docker-configs/:dockerConfigId", ctl.DockerConfigDelete)
	}

	/*******驱动管理 *******/
	{
		v1Auth.POST("device-libraries", ctl.DeviceLibraryAdd)
		v1Auth.GET("device-libraries", ctl.DeviceLibrariesSearch)
		v1Auth.DELETE("device-libraries/:deviceLibraryId", ctl.DeviceLibraryDelete)
		v1Auth.PUT("device-libraries/:deviceLibraryId", ctl.DeviceLibraryUpdate)
		v1Auth.GET("driver-classify", ctl.DeviceClassify)

	}
	/*******驱动实例 *******/
	{
		v1Auth.GET("device-servers", ctl.DeviceServicesSearch)
		v1Auth.PUT("device-server/:deviceServiceId", ctl.DeviceServiceUpdate)
	}

	/*******驱动市场分类 *******/
	{
		v1Auth.GET("device-classify", ctl.DeviceClassify)
	}

	/*******产品管理 *******/
	{
		v1Auth.GET("products", ctl.ProductsSearch)
		v1Auth.GET("product/:productId", ctl.ProductById)
		v1Auth.POST("product", ctl.ProductAdd)
		v1Auth.POST("product-release/:productId", ctl.ProductRelease)
		v1Auth.POST("product-unrelease/:productId", ctl.ProductUnRelease)
		v1Auth.DELETE("product/:productId", ctl.ProductDelete)
		v1Auth.GET("iot-platform", ctl.IotPlatform)
	}
	/*******产品物模型管理 *******/
	{
		v1Auth.GET("thingmodel/system", ctl.SystemThingModelSearch)
		v1Auth.POST("thingmodel", ctl.ThingModelAdd)
		v1Auth.PUT("thingmodel", ctl.ThingModelUpdate)
		v1Auth.DELETE("thingmodel", ctl.ThingModelDelete)
		v1Auth.GET("thingmodel/unit", ctl.ThingModelUnit)
		v1Auth.POST("thingmodel/unit-sync", ctl.ThingModelUnitSync)                       //废弃
		v1Auth.POST("thingmodel/docs-sync", ctl.ThingModelDocsSync)                       //废弃
		v1Auth.POST("thingmodel/quicknavigation-sync", ctl.ThingModelQuickNavigationSync) //废弃

	}

	/*******设备管理 *******/
	{
		v1Auth.POST("device", ctl.DeviceByAdd)
		v1Auth.GET("devices", ctl.DevicesSearch)
		v1Auth.GET("device/:deviceId", ctl.DeviceById)
		v1Auth.DELETE("device/:deviceId", ctl.DeviceDelete)
		v1Auth.DELETE("devices", ctl.DevicesDelete)
		v1Auth.PUT("device/:deviceId", ctl.DeviceUpdate)
		v1Auth.GET("device-mqtt/:deviceId", ctl.DeviceMqttInfoById)
		v1Auth.POST("device-mqtt", ctl.AddMqttAuth)
		v1Auth.GET("device/:deviceId/thing-model/property", ctl.DeviceThingModelPropertyDataSearch)
		v1Auth.GET("device/:deviceId/thing-model/history-property", ctl.DeviceThingModelHistoryPropertyDataSearch)
		v1Auth.GET("device/:deviceId/thing-model/event", ctl.DeviceThingModelEventDataSearch)
		v1Auth.GET("device/:deviceId/thing-model/service", ctl.DeviceThingModelServiceDataSearch)
		v1Auth.GET("device/status-template", ctl.DeviceStatusTemplate)
		v1Auth.GET("devices/import-template", ctl.DeviceImportTemplateDownload)
		v1Auth.POST("devices/import", ctl.DevicesImport)
		v1Auth.POST("device/upload-validated", ctl.UploadValidated)
		v1Auth.PUT("devices/bind-driver", ctl.DevicesBindDriver)
		v1Auth.PUT("devices/unbind-driver", ctl.DevicesUnBindDriver)
		v1Auth.PUT("devices/bind-product", ctl.DevicesBindByProductId)

	}
	/*******品类、物模型同步接口 *******/
	{
		v1Auth.GET("category-template", ctl.CategoryTemplateSearch)
		v1Auth.POST("category-template/sync", ctl.CategoryTemplateSync) //废弃
		v1Auth.GET("thingmodel-template", ctl.ThingModelTemplateSearch)
		v1Auth.GET("thingmodel-template/:categoryKey", ctl.ThingModelTemplateByCategoryKey)
		v1Auth.POST("thingmodel-template/sync", ctl.ThingModelTemplateSync) //废弃
	}

	/*******告警中心接口 *******/
	{
		v1Auth.POST("alert-rule", ctl.AlertRuleAdd)
		v1Auth.PUT("alert-rule/:ruleId", ctl.AlertRuleUpdate)
		v1Auth.PUT("rule-field", ctl.AlertRuleUpdateField)
		v1Auth.GET("alert-rule/:ruleId", ctl.AlertRuleById)
		v1Auth.GET("alert-rule", ctl.AlertRuleSearch)
		v1Auth.DELETE("alert-rule/:ruleId", ctl.AlertRuleDelete)
		v1Auth.POST("alert-rule/:ruleId/start", ctl.AlertRuleStart)
		v1Auth.POST("alert-rule/:ruleId/stop", ctl.AlertRuleStop)
		v1Auth.POST("alert-rule/:ruleId/restart", ctl.AlertRuleRestart)
		v1Auth.GET("alert-list", ctl.AlertSearch)
		v1Auth.GET("alert-plate", ctl.AlertPlate)
		v1Auth.PUT("alert-ignore/:ruleId", ctl.AlertIgnore)
		v1Auth.POST("alert-treated", ctl.AlertTreated)

	}
	/*******规则引擎 *******/
	{
		v1Auth.POST("rule-engine", ctl.RuleEngineAdd)
		v1Auth.PUT("rule-engine", ctl.RuleEngineUpdate)
		v1Auth.GET("rule-engine/:ruleEngineId", ctl.RuleEngineById)
		v1Auth.GET("rule-engine", ctl.RuleEngineSearch)
		v1Auth.POST("rule-engine/:ruleEngineId/start", ctl.RuleEngineStart)
		v1Auth.POST("rule-engine/:ruleEngineId/stop", ctl.RuleEngineStop)
		v1Auth.DELETE("rule-engine/:ruleEngineId/delete", ctl.RuleEngineDelete)
		v1Auth.GET("rule-engine/:ruleEngineId/status", ctl.RuleEngineStatus)

	}
	/*******资源管理 *******/
	{
		v1Auth.GET("typeresource", ctl.DataResourceType)
		v1Auth.PUT("dataresource", ctl.UpdateDataResource)
		v1Auth.POST("dataresource", ctl.DataResourceAdd)
		v1Auth.DELETE("dataresource/:dataResourceId", ctl.DataResourceDel)
		v1Auth.GET("dataresource", ctl.DataResourceSearch)
		v1Auth.GET("dataresource/:dataResourceId", ctl.DataResourceById)
		v1Auth.POST("dataresource/:dataResourceId/health", ctl.DataResourceHealth)

	}
	/*******场景联动 *******/
	{

		v1Auth.POST("scene", ctl.SceneAdd)
		v1Auth.PUT("scene", ctl.SceneUpdate)
		v1Auth.GET("scene/:sceneId", ctl.SceneById)
		v1Auth.GET("scene", ctl.SearchScene)
		v1Auth.POST("scene/:sceneId/start", ctl.SceneStart)
		v1Auth.POST("scene/:sceneId/stop", ctl.SceneStop)
		v1Auth.DELETE("scene/:sceneId", ctl.DeleteScene)
		v1Auth.GET("scene/:sceneId/log", ctl.SceneLogSearch)
	}
	/*******文档中心（sdk） *******/
	{

		v1Auth.GET("language-sdk", ctl.LanguageSdkSearch)
		v1Auth.POST("language-sdk-sync", ctl.LanguageSdkSync) //废弃
	}

}
