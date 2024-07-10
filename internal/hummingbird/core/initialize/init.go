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
package initialize

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/winc-link/hummingbird/internal/dtos"
	"github.com/winc-link/hummingbird/internal/hummingbird/core/application"
	"github.com/winc-link/hummingbird/internal/hummingbird/core/application/alertcentreapp"
	"github.com/winc-link/hummingbird/internal/hummingbird/core/application/categorytemplate"
	"github.com/winc-link/hummingbird/internal/hummingbird/core/application/dataresource"
	"github.com/winc-link/hummingbird/internal/hummingbird/core/application/deviceapp"
	"github.com/winc-link/hummingbird/internal/hummingbird/core/application/dmi"
	"github.com/winc-link/hummingbird/internal/hummingbird/core/application/docapp"
	"github.com/winc-link/hummingbird/internal/hummingbird/core/application/driverapp"
	"github.com/winc-link/hummingbird/internal/hummingbird/core/application/driverserviceapp"
	"github.com/winc-link/hummingbird/internal/hummingbird/core/application/homepageapp"
	"github.com/winc-link/hummingbird/internal/hummingbird/core/application/languagesdkapp"
	"github.com/winc-link/hummingbird/internal/hummingbird/core/application/messageapp"
	"github.com/winc-link/hummingbird/internal/hummingbird/core/application/messagestore"
	"github.com/winc-link/hummingbird/internal/hummingbird/core/application/monitor"
	"github.com/winc-link/hummingbird/internal/hummingbird/core/application/persistence"
	"github.com/winc-link/hummingbird/internal/hummingbird/core/application/productapp"
	"github.com/winc-link/hummingbird/internal/hummingbird/core/application/quicknavigationapp"
	"github.com/winc-link/hummingbird/internal/hummingbird/core/application/ruleengine"
	"github.com/winc-link/hummingbird/internal/hummingbird/core/application/scene"
	"github.com/winc-link/hummingbird/internal/hummingbird/core/application/thingmodelapp"
	"github.com/winc-link/hummingbird/internal/hummingbird/core/application/thingmodeltemplate"
	"github.com/winc-link/hummingbird/internal/hummingbird/core/application/timerapp"
	"github.com/winc-link/hummingbird/internal/hummingbird/core/application/unittemplate"
	"github.com/winc-link/hummingbird/internal/hummingbird/core/application/userapp"
	"github.com/winc-link/hummingbird/internal/hummingbird/core/config"
	"github.com/winc-link/hummingbird/internal/hummingbird/core/container"
	"github.com/winc-link/hummingbird/internal/hummingbird/core/controller/rpcserver/driverserver"
	interfaces "github.com/winc-link/hummingbird/internal/hummingbird/core/interface"
	"github.com/winc-link/hummingbird/internal/hummingbird/core/route"
	"github.com/winc-link/hummingbird/internal/pkg/constants"
	pkgContainer "github.com/winc-link/hummingbird/internal/pkg/container"
	"github.com/winc-link/hummingbird/internal/pkg/cos"
	"github.com/winc-link/hummingbird/internal/pkg/crontab"
	"github.com/winc-link/hummingbird/internal/pkg/di"
	"github.com/winc-link/hummingbird/internal/pkg/handlers"
	"github.com/winc-link/hummingbird/internal/pkg/logger"
	"github.com/winc-link/hummingbird/internal/pkg/startup"
	"github.com/winc-link/hummingbird/internal/pkg/timer/jobrunner"
	"github.com/winc-link/hummingbird/internal/tools/ekuiperclient"
	"github.com/winc-link/hummingbird/internal/tools/hpcloudclient"
	"github.com/winc-link/hummingbird/internal/tools/notify/sms"
	"github.com/winc-link/hummingbird/internal/tools/streamclient"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"sync"
)

// Bootstrap contains references to dependencies required by the BootstrapHandler.
type Bootstrap struct {
	router *gin.Engine
}

// NewBootstrap is a factory method that returns an initialized Bootstrap receiver struct.
func NewBootstrap(router *gin.Engine) *Bootstrap {
	return &Bootstrap{
		router: router,
	}
}

func (b *Bootstrap) BootstrapHandler(ctx context.Context, wg *sync.WaitGroup, _ startup.Timer, dic *di.Container) bool {

	configuration := container.ConfigurationFrom(dic.Get)
	lc := pkgContainer.LoggingClientFrom(dic.Get)

	if !b.initClient(ctx, wg, dic, configuration, lc) {
		return false
	}

	if !initApp(ctx, configuration, dic) {
		return false
	}

	// rpc 服务
	if ok := initRPCServer(ctx, wg, dic); !ok {
		return false
	}
	lc.Infof("init rpc server")

	// http 路由
	route.LoadRestRoutes(b.router, dic)

	// 业务逻辑
	application.InitSchedule(dic, lc)

	wg.Add(1)
	go func() {
		defer wg.Done()

		<-ctx.Done()
		crontab.Stop()
	}()

	return true
}

func (b *Bootstrap) initClient(ctx context.Context, wg *sync.WaitGroup, dic *di.Container, configuration *config.ConfigurationStruct, lc logger.LoggingClient) bool {

	appMode, err := dmi.New(dic, ctx, wg, dtos.DriverConfigManage{
		DockerManageConfig: dtos.DockerManageConfig{
			ContainerConfigPath: configuration.DockerManage.ContainerConfigPath,
			DockerApiVersion:    configuration.DockerManage.DockerApiVersion,
			DockerRunMode:       constants.NetworkModeHost,
			DockerSelfName:      constants.CoreServiceName,
			Privileged:          configuration.DockerManage.Privileged,
		},
	})
	if err != nil {
		lc.Error("create driver model interface error %v", err)
		return false
	}

	dic.Update(di.ServiceConstructorMap{
		interfaces.DriverModelInterfaceName: func(get di.Get) interface{} {
			return appMode
		},
	})
	homePageApp := homepageapp.NewHomePageApp(ctx, dic)
	dic.Update(di.ServiceConstructorMap{
		container.HomePageAppName: func(get di.Get) interface{} {
			return homePageApp
		},
	})

	languageApp := languagesdkapp.NewLanguageSDKApp(ctx, dic)
	dic.Update(di.ServiceConstructorMap{
		container.LanguageSDKAppName: func(get di.Get) interface{} {
			return languageApp
		},
	})

	monitorApp := monitor.NewMonitor(ctx, dic)
	dic.Update(di.ServiceConstructorMap{
		container.MonitorAppName: func(get di.Get) interface{} {
			return monitorApp
		},
	})

	streamClient := streamclient.NewStreamClient(lc)
	dic.Update(di.ServiceConstructorMap{
		pkgContainer.StreamClientName: func(get di.Get) interface{} {
			return streamClient
		},
	})

	driverApp := driverapp.NewDriverApp(ctx, dic)
	dic.Update(di.ServiceConstructorMap{
		container.DriverAppName: func(get di.Get) interface{} {
			return driverApp
		},
	})

	driverServiceApp := driverserviceapp.NewDriverServiceApp(ctx, dic)
	dic.Update(di.ServiceConstructorMap{
		container.DriverServiceAppName: func(get di.Get) interface{} {
			return driverServiceApp
		},
	})

	productApp := productapp.NewProductApp(ctx, dic)
	dic.Update(di.ServiceConstructorMap{
		container.ProductAppName: func(get di.Get) interface{} {
			return productApp
		},
	})

	thingModelApp := thingmodelapp.NewThingModelApp(ctx, dic)
	dic.Update(di.ServiceConstructorMap{
		container.ThingModelAppName: func(get di.Get) interface{} {
			return thingModelApp
		},
	})

	deviceApp := deviceapp.NewDeviceApp(ctx, dic)
	dic.Update(di.ServiceConstructorMap{
		container.DeviceItfName: func(get di.Get) interface{} {
			return deviceApp
		},
	})

	alertCentreApp := alertcentreapp.NewAlertCentreApp(ctx, dic)
	dic.Update(di.ServiceConstructorMap{
		container.AlertRuleAppName: func(get di.Get) interface{} {
			return alertCentreApp
		},
	})

	ruleEngineApp := ruleengine.NewRuleEngineApp(ctx, dic)
	dic.Update(di.ServiceConstructorMap{
		container.RuleEngineAppName: func(get di.Get) interface{} {
			return ruleEngineApp
		},
	})

	sceneApp := scene.NewSceneApp(ctx, dic)
	dic.Update(di.ServiceConstructorMap{
		container.SceneAppName: func(get di.Get) interface{} {
			return sceneApp
		},
	})

	conJobApp := timerapp.NewCronTimer(ctx, jobrunner.NewJobRunFunc(dic), dic)
	dic.Update(di.ServiceConstructorMap{
		container.ConJobAppName: func(get di.Get) interface{} {
			return conJobApp
		},
	})

	dataResourceApp := dataresource.NewDataResourceApp(ctx, dic)
	dic.Update(di.ServiceConstructorMap{
		container.DataResourceName: func(get di.Get) interface{} {
			return dataResourceApp
		},
	})

	cosApp := cos.NewCos("", "", "")
	dic.Update(di.ServiceConstructorMap{
		container.CosAppName: func(get di.Get) interface{} {
			return cosApp
		},
	})

	categoryTemplateApp := categorytemplate.NewCategoryTemplateApp(ctx, dic)
	dic.Update(di.ServiceConstructorMap{
		container.CategoryTemplateAppName: func(get di.Get) interface{} {
			return categoryTemplateApp
		},
	})

	unitTemplateApp := unittemplate.NewUnitTemplateApp(ctx, dic)
	dic.Update(di.ServiceConstructorMap{
		container.UnitTemplateAppName: func(get di.Get) interface{} {
			return unitTemplateApp
		},
	})

	docsApp := docapp.NewDocsApp(ctx, dic)
	dic.Update(di.ServiceConstructorMap{
		container.DocsAppName: func(get di.Get) interface{} {
			return docsApp
		},
	})

	quickNavigationApp := quicknavigationapp.NewQuickNavigationApp(ctx, dic)
	dic.Update(di.ServiceConstructorMap{
		container.QuickNavigationAppName: func(get di.Get) interface{} {
			return quickNavigationApp
		},
	})

	thingModelTemplateApp := thingmodeltemplate.NewThingModelTemplateApp(ctx, dic)
	dic.Update(di.ServiceConstructorMap{
		container.ThingModelTemplateAppName: func(get di.Get) interface{} {
			return thingModelTemplateApp
		},
	})

	hpcloudServiceApp := hpcloudclient.NewHpcloud(lc)
	dic.Update(di.ServiceConstructorMap{
		container.HpcServiceAppName: func(get di.Get) interface{} {
			return hpcloudServiceApp
		},
	})

	smsServiceApp := sms.NewSmsClient(lc, "",
		"", "")
	dic.Update(di.ServiceConstructorMap{
		container.SmsServiceAppName: func(get di.Get) interface{} {
			return smsServiceApp
		},
	})

	limitMethodApp := application.NewLimitMethodConf(*configuration)
	dic.Update(di.ServiceConstructorMap{
		pkgContainer.LimitMethodConfName: func(get di.Get) interface{} {
			return limitMethodApp
		},
	})

	ekuiperApp := ekuiperclient.New(configuration.Clients["Ekuiper"].Address(), lc)
	dic.Update(di.ServiceConstructorMap{
		container.EkuiperAppName: func(get di.Get) interface{} {
			return ekuiperApp
		},
	})

	persistItf := persistence.NewPersistApp(dic)
	dic.Update(di.ServiceConstructorMap{
		container.PersistItfName: func(get di.Get) interface{} {
			return persistItf
		},
	})

	userItf := userapp.New(dic)
	dic.Update(di.ServiceConstructorMap{
		container.UserItfName: func(get di.Get) interface{} {
			return userItf
		},
	})

	messageItf := messageapp.NewMessageApp(dic, configuration.Clients["Ekuiper"].Address())
	dic.Update(di.ServiceConstructorMap{
		container.MessageItfName: func(get di.Get) interface{} {
			return messageItf
		},
	})

	messageStoreItf := messagestore.NewMessageStore(dic)
	dic.Update(di.ServiceConstructorMap{
		container.MessageStoreItfName: func(get di.Get) interface{} {
			return messageStoreItf
		},
	})
	return true
}

func initRPCServer(ctx context.Context, wg *sync.WaitGroup, dic *di.Container) bool {
	lc := pkgContainer.LoggingClientFrom(dic.Get)
	_, err := handlers.NewRPCServer(ctx, wg, dic, func(serve *grpc.Server) {
		driverserver.RegisterRPCService(lc, dic, serve)
		reflection.Register(serve)
	})
	if err != nil {
		lc.Errorf("initRPCServer err:%v", err)
		return false
	}
	return true
}
