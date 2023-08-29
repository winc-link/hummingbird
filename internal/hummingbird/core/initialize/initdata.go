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

package initialize

import (
	"context"
	"github.com/winc-link/hummingbird/internal/hummingbird/core/container"
	pkgContainer "github.com/winc-link/hummingbird/internal/pkg/container"
	"github.com/winc-link/hummingbird/internal/pkg/di"
	"github.com/winc-link/hummingbird/internal/pkg/errort"
	"github.com/winc-link/hummingbird/internal/pkg/logger"
	"sync"
)

func initApp(ctx context.Context, dic *di.Container) bool {
	lc := pkgContainer.LoggingClientFrom(dic.Get)
	dbClient := container.DBClientFrom(dic.Get)
	_, edgeXErr := dbClient.GetUserByUserName("admin")
	if edgeXErr != nil {
		if errort.Is(errort.AppPasswordError, edgeXErr) {
			var wg sync.WaitGroup
			wg.Add(6)
			go syncQuickNavigation(&wg, dic, lc)
			go syncDocTemplate(&wg, dic, lc)
			go syncUnitTemplate(&wg, dic, lc)
			go syncCategory(&wg, dic, lc)
			go syncThingModel(&wg, dic, lc)
			go syncDocuments(&wg, dic, lc)
			//go initEkuiperStreams(&wg, dic, lc)
			wg.Wait()
			lc.Infof("initApp end...")
		}

	}
	return true
}

func syncCategory(wg *sync.WaitGroup, dic *di.Container, lc logger.LoggingClient) {
	defer wg.Done()
	categoryApp := container.CategoryTemplateAppFrom(dic.Get)

	_, err := categoryApp.Sync(context.Background(), "Ireland")
	lc.Infof("sync category start...")
	if err != nil {
		lc.Errorf("sync category fail...")
	}
	lc.Infof("sync category success")

}

func syncUnitTemplate(wg *sync.WaitGroup, dic *di.Container, lc logger.LoggingClient) {
	defer wg.Done()
	unitTempApp := container.UnitTemplateAppFrom(dic.Get)

	_, err := unitTempApp.Sync(context.Background(), "Ireland")
	lc.Infof("sync unit start")

	if err != nil {
		lc.Errorf("sync unit fail")
	}
	lc.Infof("sync unit success")

}

func syncDocTemplate(wg *sync.WaitGroup, dic *di.Container, lc logger.LoggingClient) {
	defer wg.Done()
	docApp := container.DocsTemplateAppFrom(dic.Get)

	_, err := docApp.SyncDocs(context.Background(), "Ireland")
	lc.Infof("sync doc start")

	if err != nil {
		lc.Errorf("sync doc fail")
	}
	lc.Infof("sync doc success")

}

func syncQuickNavigation(wg *sync.WaitGroup, dic *di.Container, lc logger.LoggingClient) {
	defer wg.Done()
	quickApp := container.QuickNavigationAppTemplateAppFrom(dic.Get)
	_, err := quickApp.SyncQuickNavigation(context.Background(), "Ireland")
	lc.Infof("sync quickNavigation start")
	if err != nil {
		lc.Errorf("sync quickNavigation fail...", err.Error())
	}
	lc.Infof("sync quickNavigation success")
}

func syncThingModel(wg *sync.WaitGroup, dic *di.Container, lc logger.LoggingClient) {
	defer wg.Done()
	thingModelApp := container.ThingModelTemplateAppFrom(dic.Get)
	_, err := thingModelApp.Sync(context.Background(), "Ireland")
	lc.Infof("sync thingModel start")
	if err != nil {
		lc.Errorf("sync thingModel fail...", err.Error())
	}
	lc.Infof("sync thingModel success")
}

func syncDocuments(wg *sync.WaitGroup, dic *di.Container, lc logger.LoggingClient) {
	defer wg.Done()
	languageApp := container.LanguageAppNameFrom(dic.Get)
	err := languageApp.Sync(context.Background(), "Ireland")
	lc.Infof("sync language start")
	if err != nil {
		lc.Errorf("sync language fail...", err.Error())
	}
	lc.Infof("sync language success")
}
