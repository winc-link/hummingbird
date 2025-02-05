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
package mysql

import (
	"fmt"
	"github.com/winc-link/hummingbird/internal/dtos"
	"github.com/winc-link/hummingbird/internal/models"
	"github.com/winc-link/hummingbird/internal/pkg/errort"
	"github.com/winc-link/hummingbird/internal/pkg/logger"
	"github.com/winc-link/hummingbird/internal/pkg/utils"
	"github.com/winc-link/hummingbird/internal/tools/sqldb/mysql"
	"gorm.io/gorm"
	//clientSQLite "github.com/winc-link/hummingbird/internal/tools/sqldb/mysql"
)

type Client struct {
	Pool *gorm.DB
	//cache         interfaces.Cache
	client        mysql.ClientSQLite
	loggingClient logger.LoggingClient
}

func NewClient(config dtos.Configuration, lc logger.LoggingClient) (c *Client, errEdgeX error) {
	client, err := mysql.NewGormClient(config, lc)
	if err != nil {
		errEdgeX = errort.NewCommonErr(errort.DefaultSystemError, fmt.Errorf("database failed to init %w", err))
		return
	}
	// 自动建表
	//if err = client.InitTable(
	//	&models.DeviceLibrary{},
	//	&models.DeviceService{},
	//	&models.Device{},
	//	&models.DockerConfig{},
	//	&models.AdvanceConfig{},
	//	&models.SystemMetrics{},
	//	&models.CategoryTemplate{},
	//	&models.ThingModelTemplate{},
	//	&models.DriverClassify{},
	//	&models.User{},
	//	&models.LanguageSdk{},
	//	&models.Metrics{},
	//	&models.Product{},
	//	&models.Properties{},
	//	&models.Actions{},
	//	&models.Events{},
	//	&models.Unit{},
	//	&models.MqttAuth{},
	//	&models.AlertRule{},
	//	&models.Scene{},
	//	&models.SceneLog{},
	//	&models.AlertList{},
	//	&models.QuickNavigation{},
	//	&models.Doc{},
	//	&models.MsgGather{},
	//	&models.RuleEngine{},
	//	&models.DataResource{},
	//); err != nil {
	//	errEdgeX = errort.NewCommonEdgeX(errort.DefaultSystemError, "database failed to init", err)
	//	return
	//}
	c = &Client{
		client:        client,
		loggingClient: lc,
		Pool:          client.Pool,
	}
	return
}

// CloseSession closes the connections to Redis
func (c *Client) CloseSession() {
	c.client.Close()
}

func (c *Client) GetDBInstance() *gorm.DB {
	return c.Pool
}

func (c *Client) AddDeviceLibrary(dl models.DeviceLibrary) (models.DeviceLibrary, error) {
	if len(dl.Id) == 0 {
		dl.Id = utils.GenUUID()
	}
	return addDeviceLibrary(c, dl)
}

func (c *Client) DockerConfigAdd(dc models.DockerConfig) (models.DockerConfig, error) {
	if len(dc.Id) == 0 {
		dc.Id = utils.GenUUID()
	}
	return dockerConfigAdd(c, dc)
}

func (c *Client) DockerConfigById(id string) (models.DockerConfig, error) {
	return dockerConfigById(c, id)
}

func (c *Client) DockerConfigDelete(id string) error {
	return dockerConfigDeleteById(c, id)
}

func (c *Client) DockerConfigUpdate(dc models.DockerConfig) error {
	return dockerConfigUpdate(c, dc)
}

func (c *Client) DockerConfigsSearch(offset int, limit int, req dtos.DockerConfigSearchQueryRequest) (dcs []models.DockerConfig, total uint32, edgeXErr error) {
	return dockerConfigsSearch(c, offset, limit, req)
}

func (c *Client) DriverClassifySearch(offset int, limit int, req dtos.DriverClassifyQueryRequest) (dcs []models.DriverClassify, total uint32, edgeXErr error) {
	return driverClassifySearch(c, offset, limit, req)
}

func (c *Client) DeviceLibrariesSearch(offset int, limit int, req dtos.DeviceLibrarySearchQueryRequest) (deviceLibraries []models.DeviceLibrary, total uint32, edgeXErr error) {
	deviceLibraries, total, edgeXErr = deviceLibrariesSearch(c, offset, limit, req)
	if edgeXErr != nil {
		return deviceLibraries, total, edgeXErr
	}
	return deviceLibraries, total, nil
}

func (c *Client) DeviceServicesSearch(offset int, limit int, req dtos.DeviceServiceSearchQueryRequest) (deviceServices []models.DeviceService, total uint32, edgeXErr error) {
	deviceServices, total, edgeXErr = deviceServicesSearch(c, offset, limit, req)
	if edgeXErr != nil {
		return deviceServices, 0, edgeXErr
	}
	return deviceServices, total, nil
}

func (c *Client) DeviceLibraryById(id string) (deviceLibrary models.DeviceLibrary, edgeXErr error) {
	return deviceLibraryById(c, id)
}

func (c *Client) DeleteDeviceLibraryById(id string) error {
	return deleteDeviceLibraryById(c, id)
}

func (c *Client) AddDeviceService(ds models.DeviceService) (models.DeviceService, error) {
	// 驱动实例和驱动id一样，为了防止容器实例名冲突导致数据冲突
	if len(ds.Id) == 0 {
		ds.Id = utils.RandomNum()
	}
	ds.Name = ds.Name + "-" + ds.Id
	return addDeviceService(c, ds)
}

func (c *Client) UpdateDeviceService(ds models.DeviceService) error {
	return updateDeviceService(c, ds)
}

func (c *Client) UpdateDeviceLibrary(dl models.DeviceLibrary) error {
	return updateDeviceLibrary(c, dl)
}

func (c *Client) DeviceServiceById(id string) (deviceService models.DeviceService, edgeXErr error) {
	deviceService, edgeXErr = deviceServiceById(c, id)
	if edgeXErr != nil {
		return deviceService, edgeXErr
	}

	return
}

func (c *Client) DeleteDeviceServiceById(id string) error {
	return deleteDeviceServiceById(c, id)
}

func (c *Client) ProductById(id string) (product models.Product, edgeXErr error) {
	return productById(c, id)
}

func (c *Client) AddProduct(ds models.Product) (product models.Product, edgeXErr error) {
	if len(ds.Id) == 0 {
		ds.Id = utils.RandomNum()
	}
	return addProduct(c, ds)
}

func (c *Client) ProductByCloudId(id string) (product models.Product, edgeXErr error) {
	return productByCloudId(c, id)
}

func (c *Client) BatchUpsertProduct(p []models.Product) (int64, error) {
	return batchUpsertProduct(c, p)
}

func (c *Client) BatchSaveProduct(p []models.Product) error {
	return batchSaveProduct(c, p)
}

func (c *Client) BatchDeleteProduct(products []models.Product) error {
	return batchDeleteProduct(c, products)
}

func (c *Client) BatchDeleteProperties(propertiesIds []string) error {
	return batchDeleteProperties(c, propertiesIds)
}

func (c *Client) BatchDeleteSystemProperties() error {
	return batchDeleteSystemProperties(c)
}

func (c *Client) BatchInsertSystemProperties(p []models.Properties) (int64, error) {
	return batchInsertSystemProperties(c, p)
}

func (c *Client) BatchDeleteEvents(eventIds []string) error {
	return batchDeleteEvents(c, eventIds)
}

func (c *Client) BatchDeleteSystemEvents() error {
	return batchDeleteSystemEvents(c)
}

func (c *Client) BatchInsertSystemEvents(p []models.Events) (int64, error) {
	return batchInsertSystemEvents(c, p)
}

func (c *Client) BatchDeleteActions(actionIds []string) error {
	return batchDeleteActions(c, actionIds)
}
func (c *Client) BatchDeleteSystemActions() error {
	return batchDeleteSystemActions(c)
}

func (c *Client) BatchInsertSystemActions(p []models.Actions) (int64, error) {
	return batchInsertSystemActions(c, p)
}

func (c *Client) DeleteProductById(id string) error {
	return deleteProductById(c, id)
}

func (c *Client) DeleteProductObject(product models.Product) error {
	return deleteProductObject(c, product)
}

func (c *Client) AssociationsDeleteProductObject(product models.Product) error {
	return associationsDeleteProductObject(c, product)
}

func (c *Client) UpdateProduct(ds models.Product) error {
	return updateProduct(c, ds)
}

func (c *Client) AssociationsUpdateProduct(ds models.Product) error {
	return associationsUpdateProduct(c, ds)
}

func (c *Client) BatchUpsertDevice(p []models.Device) (int64, error) {
	return batchUpsertDevice(c, p)
}

func (c *Client) ProductsSearch(offset int, limit int, preload bool, req dtos.ProductSearchQueryRequest) (products []models.Product, total uint32, edgeXErr error) {
	products, total, edgeXErr = productsSearch(c, offset, limit, preload, req)
	if edgeXErr != nil {
		return products, 0, edgeXErr
	}
	return products, total, nil
}

func (c *Client) DevicesSearch(offset int, limit int, req dtos.DeviceSearchQueryRequest) (devices []models.Device, total uint32, edgeXErr error) {
	devices, total, edgeXErr = devicesSearch(c, offset, limit, req)
	if edgeXErr != nil {
		return devices, 0, edgeXErr
	}
	return devices, total, nil
}

func (c *Client) DeviceById(id string) (device models.Device, edgeXErr error) {
	return deviceById(c, id)
}

func (c *Client) DeviceOnlineById(id string) (edgeXErr error) {
	return deviceOnlineById(c, id)
}

func (c *Client) DeviceOfflineById(id string) (edgeXErr error) {
	return deviceOfflineById(c, id)
}

func (c *Client) DeviceOfflineByCloudInstanceId(id string) (edgeXErr error) {
	return deviceOfflineByCloudInstanceId(c, id)
}

func (c *Client) MsgReportDeviceById(id string) (device models.Device, edgeXErr error) {
	return msgReportDeviceById(c, id)
}

func (c *Client) DeviceByCloudId(id string) (device models.Device, edgeXErr error) {
	return deviceByCloudId(c, id)
}

func (c *Client) DeviceMqttAuthInfo(id string) (device models.MqttAuth, edgeXErr error) {
	return deviceMqttAuthInfo(c, id)
}

func (c *Client) DriverMqttAuthInfo(id string) (device models.MqttAuth, edgeXErr error) {
	return driverMqttAuthInfo(c, id)
}

func (c *Client) AddDevice(ds models.Device) (deviceId string, edgeXErr error) {
	if len(ds.Id) == 0 {
		ds.Id = utils.RandomNum()
	}
	return addDevice(c, ds)
}

func (c *Client) BatchDeleteDevice(ids []string) error {
	return batchDeleteDevice(c, ids)
}

func (c *Client) BatchUnBindDevice(ids []string) error {
	return batchUnBindDevice(c, ids)
}

func (c *Client) BatchBindDevice(ids []string, driverInstanceId string) error {
	return batchBindDevice(c, ids, driverInstanceId)
}

func (c *Client) DeleteDeviceById(id string) error {
	return deleteDeviceById(c, id)
}

func (c *Client) DeleteDeviceByCloudInstanceId(id string) error {
	return deleteDeviceByCloudInstanceId(c, id)
}

func (c *Client) UpdateDevice(ds models.Device) error {
	return updateDevice(c, ds)
}

func (c *Client) AbilityByCode(model interface{}, code, productId string) (interface{}, error) {
	return abilityByCode(c, model, code, productId)
}

func (c *Client) CategoryTemplateSearch(offset int, limit int, req dtos.CategoryTemplateRequest) ([]models.CategoryTemplate, uint32, error) {
	return categoryTemplateSearch(c, offset, limit, req)
}

func (c *Client) UnitSearch(offset int, limit int, req dtos.UnitRequest) ([]models.Unit, uint32, error) {
	return unitSearch(c, offset, limit, req)
}

func (c *Client) BatchUpsertUnitTemplate(p []models.Unit) (int64, error) {
	return batchUpsertUnitTemplate(c, p)
}

func (c *Client) CategoryTemplateById(id string) (models.CategoryTemplate, error) {
	return categoryTemplateById(c, id)
}

func (c *Client) BatchUpsertCategoryTemplate(p []models.CategoryTemplate) (int64, error) {
	return batchUpsertCategoryTemplate(c, p)
}

func (c *Client) ThingModelTemplateSearch(offset int, limit int, req dtos.ThingModelTemplateRequest) ([]models.ThingModelTemplate, uint32, error) {
	return thingModelTemplateSearch(c, offset, limit, req)
}
func (c *Client) ThingModelTemplateByCategoryKey(categoryKey string) (models.ThingModelTemplate, error) {
	return thingModelTemplateByCategoryKey(c, categoryKey)
}

func (c *Client) BatchUpsertThingModelTemplate(p []models.ThingModelTemplate) (int64, error) {
	return batchUpsertThingModelTemplate(c, p)
}

func (c *Client) AddThingModelProperty(ds models.Properties) (models.Properties, error) {
	if len(ds.Id) == 0 {
		ds.Id = utils.RandomNum()
	}
	return addThingModelProperty(c, ds)
}

func (c *Client) BatchUpsertThingModel(ds interface{}) (int64, error) {
	return batchUpsertThingModel(c, ds)
}

func (c *Client) AddThingModelEvent(ds models.Events) (models.Events, error) {
	if len(ds.Id) == 0 {
		ds.Id = utils.RandomNum()
	}
	return addThingModelEvent(c, ds)
}

func (c *Client) AddThingModelAction(ds models.Actions) (models.Actions, error) {
	if len(ds.Id) == 0 {
		ds.Id = utils.RandomNum()
	}
	return addThingModelAction(c, ds)
}

func (c *Client) UpdateThingModelProperty(ds models.Properties) error {
	return updateThingModelProperty(c, ds)
}

func (c *Client) UpdateThingModelEvent(ds models.Events) error {
	return updateThingModelEvent(c, ds)
}

func (c *Client) UpdateThingModelAction(ds models.Actions) error {
	return updateThingModelAction(c, ds)
}

func (c *Client) ThingModelDeleteProperty(id string) error {
	return deleteThingModelPropertyById(c, id)
}

func (c *Client) ThingModelDeleteEvent(id string) error {
	return deleteThingModelEventById(c, id)
}

func (c *Client) ThingModelDeleteAction(id string) error {
	return deleteThingModelActionById(c, id)
}

func (c *Client) ThingModelPropertyById(id string) (models.Properties, error) {
	return thingModelPropertyById(c, id)
}

func (c *Client) ThingModelEventById(id string) (models.Events, error) {
	return thingModelEventById(c, id)
}

func (c *Client) ThingModelActionsById(id string) (models.Actions, error) {
	return thingModeActionById(c, id)
}
func (c *Client) SystemThingModelSearch(modelType string, ModelName string) (interface{}, error) {
	return systemThingModelSearch(c, modelType, ModelName)
}

func (c *Client) AddMqttAuthInfo(auth models.MqttAuth) (string, error) {
	if len(auth.Id) == 0 {
		auth.Id = utils.RandomNum()
	}
	return addMqttAuth(c, auth)
}

func (c *Client) AddOrUpdateAuth(auth models.MqttAuth) error {
	if len(auth.Id) == 0 {
		auth.Id = utils.RandomNum()
	}
	return addOrUpdateAuth(c, auth)
}

func (c *Client) AddAlertRule(alertRule models.AlertRule) (models.AlertRule, error) {
	if len(alertRule.Id) == 0 {
		alertRule.Id = utils.RandomNum()
	}
	return addAlertRule(c, alertRule)
}

func (c *Client) AddAlertList(alertRule models.AlertList) (models.AlertList, error) {
	if len(alertRule.Id) == 0 {
		alertRule.Id = utils.RandomNum()
	}
	return addAlertList(c, alertRule)
}

func (c *Client) UpdateAlertRule(rule models.AlertRule) error {
	return updateAlertRule(c, rule)
}

func (c *Client) AlertRuleById(id string) (models.AlertRule, error) {
	return alertRuleById(c, id)
}

func (c *Client) AlertRuleSearch(offset int, limit int, req dtos.AlertRuleSearchQueryRequest) (alertRules []models.AlertRule, total uint32, edgeXErr error) {
	alertRules, total, edgeXErr = alertRuleSearch(c, offset, limit, req)
	if edgeXErr != nil {
		return alertRules, 0, edgeXErr
	}
	return alertRules, total, nil
}

func (c *Client) AlertListSearch(offset int, limit int, req dtos.AlertSearchQueryRequest) (alertList []dtos.AlertSearchQueryResponse, total uint32, edgeXErr error) {
	alertList, total, edgeXErr = alertListSearch(c, offset, limit, req)
	if edgeXErr != nil {
		return alertList, 0, edgeXErr
	}
	return alertList, total, nil
}

func (c *Client) AlertIgnore(id string) (edgeXErr error) {
	return alertIgnore(c, id)
}

func (c *Client) TreatedIgnore(id, message string) (edgeXErr error) {
	return treatedIgnore(c, id, message)
}

func (c *Client) AlertListLastSend(alertRuleId string) (alertList models.AlertList, edgeXErr error) {
	return alertListLastSend(c, alertRuleId)
}

func (c *Client) DeleteAlertRuleById(id string) error {
	return deleteAlertRuleById(c, id)
}

func (c *Client) AlertRuleStart(id string) error {
	return alertRuleStart(c, id)
}

func (c *Client) AlertRuleStop(id string) error {
	return alertRuleStop(c, id)
}

func (c *Client) AlertPlate(beforeTime int64) (plate []dtos.AlertPlateQueryResponse, err error) {
	return alertPlate(c, beforeTime)
}

func (c *Client) QuickNavigationSearch(offset int, limit int, req dtos.QuickNavigationSearchQueryRequest) (quickNavigations []models.QuickNavigation, total uint32, edgeXErr error) {
	quickNavigations, total, edgeXErr = quickNavigationSearch(c, offset, limit, req)
	if edgeXErr != nil {
		return quickNavigations, 0, edgeXErr
	}
	return quickNavigations, total, nil
}

func (c *Client) DocsSearch(offset int, limit int, req dtos.DocsSearchQueryRequest) (docs []models.Doc, total uint32, edgeXErr error) {
	docs, total, edgeXErr = docsSearch(c, offset, limit, req)
	if edgeXErr != nil {
		return docs, 0, edgeXErr
	}
	return docs, total, nil
}

func (c *Client) BatchUpsertDocsTemplate(ds []models.Doc) (int64, error) {
	return batchUpsertDocsTemplate(c, ds)
}

func (c *Client) BatchUpsertQuickNavigationTemplate(ds []models.QuickNavigation) (int64, error) {
	return batchUpsertQuickNavigationTemplate(c, ds)
}

func (c *Client) DeleteQuickNavigation(id string) error {
	return deleteQuickNavigation(c, id)
}

func (c *Client) GetAdvanceConfig() (models.AdvanceConfig, error) {
	return getAdvanceConfig(c)
}

func (c *Client) UpdateAdvanceConfig(config models.AdvanceConfig) error {
	return updateAdvanceConfig(c, config)
}

func (c *Client) AddMsgGather(msgGather models.MsgGather) error {
	if len(msgGather.Id) == 0 {
		msgGather.Id = utils.RandomNum()
	}
	return addMsgGather(c, msgGather)
}

func (c *Client) MsgGatherSearch(offset int, limit int, req dtos.MsgGatherSearchQueryRequest) (dcs []models.MsgGather, count uint32, edgeXErr error) {
	return msgGatherSearch(c, offset, limit, req)
}

func (c *Client) AddDataResource(dateResource models.DataResource) (string, error) {
	if len(dateResource.Id) == 0 {
		dateResource.Id = utils.RandomNum()
	}
	return addDataResource(c, dateResource)
}

func (c *Client) UpdateDataResource(dateResource models.DataResource) error {
	return updateDataResource(c, dateResource)
}

func (c *Client) DelDataResource(id string) error {
	return deleteDataResourceById(c, id)
}

func (c *Client) UpdateDataResourceHealth(id string, health bool) error {
	return updateDataResourceHealth(c, id, health)
}

func (c *Client) SearchDataResource(offset int, limit int, req dtos.DataResourceSearchQueryRequest) (dataResource []models.DataResource, count uint32, edgeXErr error) {
	return dataResourceSearch(c, offset, limit, req)
}

func (c *Client) DataResourceById(id string) (models.DataResource, error) {
	return dataResourceById(c, id)
}

func (c *Client) AddRuleEngine(ruleEngine models.RuleEngine) (string, error) {
	if len(ruleEngine.Id) == 0 {
		ruleEngine.Id = utils.RandomNum()
	}
	return addRuleEngine(c, ruleEngine)
}

func (c *Client) UpdateRuleEngine(ruleEngine models.RuleEngine) error {
	return updateRuleEngine(c, ruleEngine)
}

func (c *Client) RuleEngineById(id string) (ruleEngine models.RuleEngine, edgeXErr error) {
	return ruleEngineById(c, id)
}

func (c *Client) RuleEngineSearch(offset int, limit int, req dtos.RuleEngineSearchQueryRequest) (ruleEngine []models.RuleEngine, count uint32, edgeXErr error) {
	return ruleEngineSearch(c, offset, limit, req)
}

func (c *Client) RuleEngineStart(id string) error {
	return ruleEngineStart(c, id)
}

func (c *Client) RuleEngineStop(id string) error {
	return ruleEngineStop(c, id)
}

func (c *Client) DeleteRuleEngineById(id string) error {
	return deleteRuleEngineById(c, id)
}

func (c *Client) AddScene(scene models.Scene) (models.Scene, error) {
	if len(scene.Id) == 0 {
		scene.Id = utils.RandomNum()
	}
	return addScene(c, scene)
}

func (c *Client) UpdateScene(scene models.Scene) error {
	if len(scene.Id) == 0 {
		scene.Id = utils.RandomNum()
	}
	return updateScene(c, scene)
}
func (c *Client) SceneById(id string) (models.Scene, error) {
	return sceneById(c, id)
}

func (c *Client) SceneStart(id string) error {
	return sceneStart(c, id)
}

func (c *Client) SceneStop(id string) error {
	return sceneStop(c, id)
}

func (c *Client) DeleteSceneById(id string) error {
	return deleteSceneById(c, id)
}

func (c *Client) SceneSearch(offset int, limit int, req dtos.SceneSearchQueryRequest) (scenes []models.Scene, total uint32, edgeXErr error) {
	return sceneSearch(c, offset, limit, req)
}

func (c *Client) AddSceneLog(sceneLog models.SceneLog) (models.SceneLog, error) {
	if len(sceneLog.Id) == 0 {
		sceneLog.Id = utils.RandomNum()
	}
	return addSceneLog(c, sceneLog)
}

func (c *Client) SceneLogSearch(offset int, limit int, req dtos.SceneLogSearchQueryRequest) (sceneLogs []models.SceneLog, total uint32, edgeXErr error) {
	return sceneLogSearch(c, offset, limit, req)
}

func (c *Client) LanguageSdkByName(name string) (cloudService models.LanguageSdk, edgeXErr error) {
	return languageByName(c, name)
}

func (c *Client) LanguageSearch(offset int, limit int, req dtos.LanguageSDKSearchQueryRequest) (languages []models.LanguageSdk, count uint32, edgeXErr error) {
	return languageSearch(c, offset, limit, req)
}

func (c *Client) AddLanguageSdk(ls models.LanguageSdk) (language models.LanguageSdk, edgeXErr error) {
	if len(ls.Id) == 0 {
		ls.Id = utils.RandomNum()
	}
	return addLanguageSdk(c, ls)
}

func (c *Client) UpdateLanguageSdk(ls models.LanguageSdk) error {
	return updateLanguageSdk(c, ls)
}

func (c *Client) UpdateSystemMetrics(metrics dtos.SystemMetrics) error {
	return updateSystemMetrics(c, metrics)
}

func (c *Client) GetSystemMetrics(start, end int64) ([]dtos.SystemMetrics, error) {
	return getSystemMetrics(c, start, end)
}

func (c *Client) RemoveRangeSystemMetrics(min, max string) error {
	return removeRangeSystemMetrics(c, min, max)
}
