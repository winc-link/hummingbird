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

package languagesdkapp

import (
	"context"
	"encoding/json"
	"github.com/winc-link/hummingbird/internal/dtos"
	"github.com/winc-link/hummingbird/internal/hummingbird/core/container"
	interfaces "github.com/winc-link/hummingbird/internal/hummingbird/core/interface"
	"github.com/winc-link/hummingbird/internal/models"
	pkgcontainer "github.com/winc-link/hummingbird/internal/pkg/container"
	"github.com/winc-link/hummingbird/internal/pkg/di"
	"github.com/winc-link/hummingbird/internal/pkg/logger"
)

type languageSDKApp struct {
	dic      *di.Container
	lc       logger.LoggingClient
	dbClient interfaces.DBClient
}

func (m languageSDKApp) LanguageSDKSearch(ctx context.Context, req dtos.LanguageSDKSearchQueryRequest) ([]dtos.LanguageSDKSearchResponse, uint32, error) {
	offset, limit := req.BaseSearchConditionQuery.GetPage()
	req.BaseSearchConditionQuery.OrderBy = "sort:asc"
	languages, total, err := m.dbClient.LanguageSearch(offset, limit, req)

	if err != nil {
		return nil, 0, err
	}
	libs := make([]dtos.LanguageSDKSearchResponse, len(languages))
	for i, language := range languages {
		libs[i] = dtos.LanguageSDKSearchResponse{
			Name:        language.Name,
			Icon:        language.Icon,
			Addr:        language.Addr,
			Description: language.Description,
		}
	}
	return libs, total, nil

}

func (m languageSDKApp) Sync(ctx context.Context, versionName string) error {
	filePath := versionName + "/language_sdk.json"
	cosApp := container.CosAppNameFrom(m.dic.Get)
	bs, err := cosApp.Get(filePath)
	if err != nil {
		m.lc.Errorf(err.Error())
	}
	var cosLanguageSdkResp []dtos.LanguageSDK
	err = json.Unmarshal(bs, &cosLanguageSdkResp)
	if err != nil {
		m.lc.Errorf(err.Error())
	}

	for _, sdk := range cosLanguageSdkResp {
		if languageSdk, err := m.dbClient.LanguageSdkByName(sdk.Name); err != nil {
			createModel := models.LanguageSdk{
				Name:        sdk.Name,
				Icon:        sdk.Icon,
				Sort:        sdk.Sort,
				Addr:        sdk.Addr,
				Description: sdk.Description,
			}
			_, err := m.dbClient.AddLanguageSdk(createModel)
			if err != nil {
				return err
			}
		} else {
			updateModel := models.LanguageSdk{
				Id:          languageSdk.Id,
				Name:        sdk.Name,
				Icon:        sdk.Icon,
				Sort:        sdk.Sort,
				Addr:        sdk.Addr,
				Description: sdk.Description,
			}
			err = m.dbClient.UpdateLanguageSdk(updateModel)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func NewLanguageSDKApp(ctx context.Context, dic *di.Container) interfaces.LanguageSDKApp {
	app := &languageSDKApp{
		dic:      dic,
		lc:       pkgcontainer.LoggingClientFrom(dic.Get),
		dbClient: container.DBClientFrom(dic.Get),
	}
	return app
}
