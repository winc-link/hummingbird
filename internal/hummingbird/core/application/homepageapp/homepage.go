/*******************************************************************************
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
package homepageapp

import (
	"context"
	"github.com/winc-link/hummingbird/internal/dtos"
	resourceContainer "github.com/winc-link/hummingbird/internal/hummingbird/core/container"
	interfaces "github.com/winc-link/hummingbird/internal/hummingbird/core/interface"
	"github.com/winc-link/hummingbird/internal/models"
	"github.com/winc-link/hummingbird/internal/pkg/constants"
	"github.com/winc-link/hummingbird/internal/pkg/container"
	"github.com/winc-link/hummingbird/internal/pkg/di"
	"github.com/winc-link/hummingbird/internal/pkg/logger"
	"time"
)

const (
	HummingbridDoc = "https://doc.hummingbird.winc-link.com/"
)

type homePageApp struct {
	dic      *di.Container
	lc       logger.LoggingClient
	dbClient interfaces.DBClient
}

func NewHomePageApp(ctx context.Context, dic *di.Container) interfaces.HomePageItf {
	dbClient := resourceContainer.DBClientFrom(dic.Get)
	return &homePageApp{
		dic:      dic,
		lc:       container.LoggingClientFrom(dic.Get),
		dbClient: dbClient,
	}
}

func (h homePageApp) HomePageInfo(ctx context.Context, req dtos.HomePageRequest) (response dtos.HomePageResponse, err error) {
	var responseResponse dtos.HomePageResponse
	devices, deviceTotal, err := h.dbClient.DevicesSearch(0, -1, dtos.DeviceSearchQueryRequest{})
	var selfDeviceTotal uint32
	for _, device := range devices {
		if device.Platform == constants.IotPlatform_LocalIot {
			selfDeviceTotal++
		}
	}
	responseResponse.PageInfo.Device.Total = deviceTotal
	responseResponse.PageInfo.Device.Self = selfDeviceTotal
	if deviceTotal-selfDeviceTotal < 0 {
		responseResponse.PageInfo.Device.Other = 0
	} else {
		responseResponse.PageInfo.Device.Other = deviceTotal - selfDeviceTotal
	}

	products, productTotal, err := h.dbClient.ProductsSearch(0, -1, false, dtos.ProductSearchQueryRequest{})
	var selfProductTotal uint32
	for _, product := range products {
		if product.Platform == constants.IotPlatform_LocalIot {
			selfProductTotal++
		}
	}
	responseResponse.PageInfo.Product.Total = productTotal
	responseResponse.PageInfo.Product.Self = selfProductTotal
	if productTotal-selfProductTotal < 0 {
		responseResponse.PageInfo.Product.Other = 0
	} else {
		responseResponse.PageInfo.Product.Other = productTotal - selfProductTotal
	}

	responseResponse.PageInfo.CloudInstance.StopCount = responseResponse.PageInfo.CloudInstance.Count - responseResponse.PageInfo.CloudInstance.RunCount
	if responseResponse.PageInfo.CloudInstance.StopCount < 0 {
		responseResponse.PageInfo.CloudInstance.StopCount = 0
	}
	var searchQuickNavigationReq dtos.QuickNavigationSearchQueryRequest
	searchQuickNavigationReq.OrderBy = "sort"
	quickNavigations, _, _ := h.dbClient.QuickNavigationSearch(0, -1, searchQuickNavigationReq)
	navigations := make([]dtos.QuickNavigation, 0)
	for _, navigation := range quickNavigations {
		navigations = append(navigations, dtos.QuickNavigation{
			Id:   navigation.Id,
			Name: navigation.Name,
			Icon: navigation.Icon,
			//JumpLink: navigation.JumpLink,
		})
	}
	responseResponse.QuickNavigation = navigations

	var searchDocsReq dtos.DocsSearchQueryRequest
	searchDocsReq.OrderBy = "sort"
	dbDocs, _, _ := h.dbClient.DocsSearch(0, -1, searchDocsReq)
	docs := make([]dtos.Doc, 0)
	for _, doc := range dbDocs {
		docs = append(docs, dtos.Doc{
			Name:     doc.Name,
			JumpLink: doc.JumpLink,
		})
	}
	responseResponse.Docs.More = HummingbridDoc
	responseResponse.Docs.Doc = docs

	alertRuleApp := resourceContainer.AlertRuleAppNameFrom(h.dic.Get)
	alertResp, _ := alertRuleApp.AlertPlate(ctx, time.Now().AddDate(0, 0, -1).UnixMilli())
	responseResponse.AlertPlate = alertResp

	var alertTotal uint32
	for _, alert := range responseResponse.AlertPlate {
		alertTotal += uint32(alert.Count)
	}
	responseResponse.PageInfo.Alert.Total = alertTotal

	//设备消息总数
	var msgGatherReq dtos.MsgGatherSearchQueryRequest
	msgGatherReq.Date = append(append(append(append(append(append(msgGatherReq.Date,
		time.Now().AddDate(0, 0, -1).Format("2006-01-02")),
		time.Now().AddDate(0, 0, -2).Format("2006-01-02")),
		time.Now().AddDate(0, 0, -3).Format("2006-01-02")),
		time.Now().AddDate(0, 0, -4).Format("2006-01-02")),
		time.Now().AddDate(0, 0, -5).Format("2006-01-02")),
		time.Now().AddDate(0, 0, -6).Format("2006-01-02"))

	msgGather, _, err := h.dbClient.MsgGatherSearch(0, -1, msgGatherReq)
	responseResponse.MsgGather = append(append(append(append(append(append(responseResponse.MsgGather, dtos.MsgGather{
		Date:  time.Now().AddDate(0, 0, -1).Format("2006-01-02"),
		Count: getMsgGatherCountByDate(msgGather, time.Now().AddDate(0, 0, -1).Format("2006-01-02")),
	}), dtos.MsgGather{
		Date:  time.Now().AddDate(0, 0, -2).Format("2006-01-02"),
		Count: getMsgGatherCountByDate(msgGather, time.Now().AddDate(0, 0, -2).Format("2006-01-02")),
	}), dtos.MsgGather{
		Date:  time.Now().AddDate(0, 0, -3).Format("2006-01-02"),
		Count: getMsgGatherCountByDate(msgGather, time.Now().AddDate(0, 0, -3).Format("2006-01-02")),
	}), dtos.MsgGather{
		Date:  time.Now().AddDate(0, 0, -4).Format("2006-01-02"),
		Count: getMsgGatherCountByDate(msgGather, time.Now().AddDate(0, 0, -4).Format("2006-01-02")),
	}), dtos.MsgGather{
		Date:  time.Now().AddDate(0, 0, -5).Format("2006-01-02"),
		Count: getMsgGatherCountByDate(msgGather, time.Now().AddDate(0, 0, -5).Format("2006-01-02")),
	}), dtos.MsgGather{
		Date:  time.Now().AddDate(0, 0, -6).Format("2006-01-02"),
		Count: getMsgGatherCountByDate(msgGather, time.Now().AddDate(0, 0, -6).Format("2006-01-02")),
	})
	return responseResponse, nil
}

func getMsgGatherCountByDate(msgGather []models.MsgGather, data string) int {
	for _, gather := range msgGather {
		if gather.Date == data {
			return gather.Count
		}
	}
	return 0
}
