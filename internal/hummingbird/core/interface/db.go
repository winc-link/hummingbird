//
// Copyright (C) 2020-2021 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package interfaces

import (
	"github.com/winc-link/hummingbird/internal/dtos"
	"github.com/winc-link/hummingbird/internal/models"
	"gorm.io/gorm"
)

type DBClient interface {
	CloseSession()
	GetDBInstance() *gorm.DB

	QuickNavigationSearch(offset int, limit int, req dtos.QuickNavigationSearchQueryRequest) ([]models.QuickNavigation, uint32, error)
	DocsSearch(offset int, limit int, req dtos.DocsSearchQueryRequest) (docs []models.Doc, total uint32, edgeXErr error)
	BatchUpsertDocsTemplate(d []models.Doc) (int64, error)
	BatchUpsertQuickNavigationTemplate(ds []models.QuickNavigation) (int64, error)
	DeleteQuickNavigation(id string) error

	AddDeviceService(ds models.DeviceService) (models.DeviceService, error)
	DeviceServiceById(id string) (models.DeviceService, error)
	DeleteDeviceServiceById(id string) error
	UpdateDeviceService(ds models.DeviceService) error
	DeviceServicesSearch(offset int, limit int, req dtos.DeviceServiceSearchQueryRequest) ([]models.DeviceService, uint32, error)

	AddProduct(ds models.Product) (models.Product, error)
	ProductsSearch(offset int, limit int, preload bool, req dtos.ProductSearchQueryRequest) ([]models.Product, uint32, error)
	ProductById(id string) (models.Product, error)
	ProductByCloudId(id string) (models.Product, error)
	BatchUpsertProduct(d []models.Product) (int64, error)
	BatchSaveProduct(p []models.Product) error
	BatchDeleteProduct(products []models.Product) error
	BatchDeleteProperties(propertiesId []string) error
	BatchDeleteSystemProperties() error
	BatchInsertSystemProperties(p []models.Properties) (int64, error)
	BatchDeleteSystemEvents() error
	BatchInsertSystemEvents(p []models.Events) (int64, error)
	BatchDeleteEvents(eventId []string) error
	BatchDeleteActions(actionId []string) error
	DeleteProductById(id string) error
	BatchDeleteSystemActions() error
	BatchInsertSystemActions(p []models.Actions) (int64, error)
	DeleteProductObject(d models.Product) error
	UpdateProduct(ds models.Product) error
	AssociationsUpdateProduct(ds models.Product) error
	AssociationsDeleteProductObject(ds models.Product) error

	AddThingModelProperty(ds models.Properties) (models.Properties, error)
	BatchUpsertThingModel(ds interface{}) (int64, error)
	AddThingModelEvent(ds models.Events) (models.Events, error)
	AddThingModelAction(ds models.Actions) (models.Actions, error)
	UpdateThingModelProperty(ds models.Properties) error
	UpdateThingModelEvent(ds models.Events) error
	UpdateThingModelAction(ds models.Actions) error
	ThingModelDeleteProperty(id string) error
	ThingModelDeleteEvent(id string) error
	ThingModelDeleteAction(id string) error
	ThingModelPropertyById(id string) (models.Properties, error)
	ThingModelEventById(id string) (models.Events, error)
	ThingModelActionsById(id string) (models.Actions, error)
	SystemThingModelSearch(modelType, modelName string) (interface{}, error)

	CategoryTemplateSearch(offset int, limit int, req dtos.CategoryTemplateRequest) ([]models.CategoryTemplate, uint32, error)
	CategoryTemplateById(id string) (models.CategoryTemplate, error)
	BatchUpsertCategoryTemplate(d []models.CategoryTemplate) (int64, error)
	ThingModelTemplateSearch(offset int, limit int, req dtos.ThingModelTemplateRequest) ([]models.ThingModelTemplate, uint32, error)
	ThingModelTemplateByCategoryKey(categoryKey string) (models.ThingModelTemplate, error)
	BatchUpsertThingModelTemplate(d []models.ThingModelTemplate) (int64, error)

	UnitSearch(offset int, limit int, req dtos.UnitRequest) ([]models.Unit, uint32, error)
	BatchUpsertUnitTemplate(d []models.Unit) (int64, error)

	AddDevice(d models.Device) (string, error)
	DeviceById(id string) (models.Device, error)
	DeviceOnlineById(id string) (edgeXErr error)
	DeviceOfflineById(id string) (edgeXxErr error)
	DeviceOfflineByCloudInstanceId(id string) (edgeXErr error)
	MsgReportDeviceById(id string) (device models.Device, edgeXErr error)
	DeviceByCloudId(id string) (models.Device, error)
	DevicesSearch(offset int, limit int, req dtos.DeviceSearchQueryRequest) ([]models.Device, uint32, error)
	DeviceMqttAuthInfo(id string) (device models.MqttAuth, edgeXErr error)
	DriverMqttAuthInfo(id string) (device models.MqttAuth, edgeXErr error)
	AddMqttAuthInfo(auth models.MqttAuth) (string, error)
	AddOrUpdateAuth(auth models.MqttAuth) error
	BatchUpsertDevice(d []models.Device) (int64, error)
	BatchDeleteDevice(deviceIds []string) error
	BatchUnBindDevice(ids []string) error
	BatchBindDevice(ids []string, driverInstanceId string) error
	DeleteDeviceById(id string) error
	UpdateDevice(ds models.Device) error
	DeleteDeviceByCloudInstanceId(cloudInstanceId string) error
	AddDeviceLibrary(dl models.DeviceLibrary) (models.DeviceLibrary, error)
	DeviceLibraryById(id string) (models.DeviceLibrary, error)
	DeleteDeviceLibraryById(id string) error
	DeviceLibrariesSearch(offset int, limit int, req dtos.DeviceLibrarySearchQueryRequest) ([]models.DeviceLibrary, uint32, error)
	UpdateDeviceLibrary(dl models.DeviceLibrary) error

	DriverClassifySearch(offset int, limit int, req dtos.DriverClassifyQueryRequest) ([]models.DriverClassify, uint32, error)

	DockerConfigAdd(cfg models.DockerConfig) (models.DockerConfig, error)
	DockerConfigById(id string) (models.DockerConfig, error)
	DockerConfigUpdate(cfg models.DockerConfig) error
	DockerConfigDelete(id string) error
	DockerConfigsSearch(offset int, limit int, req dtos.DockerConfigSearchQueryRequest) ([]models.DockerConfig, uint32, error)

	AbilityByCode(model interface{}, code, productId string) (interface{}, error)

	//// 获取高级配置信息
	GetAdvanceConfig() (models.AdvanceConfig, error)
	// 更新高级配置信息
	UpdateAdvanceConfig(config models.AdvanceConfig) error

	AddMsgGather(msgGather models.MsgGather) error
	MsgGatherSearch(offset int, limit int, req dtos.MsgGatherSearchQueryRequest) (msgGather []models.MsgGather, count uint32, edgeXErr error)

	AddDataResource(dateResource models.DataResource) (string, error)
	UpdateDataResource(dateResource models.DataResource) error
	DelDataResource(id string) error
	//DataResourceById(id string) models.DataResource

	UpdateDataResourceHealth(id string, health bool) error
	SearchDataResource(offset int, limit int, req dtos.DataResourceSearchQueryRequest) (dataResource []models.DataResource, count uint32, edgeXErr error)
	DataResourceById(id string) (models.DataResource, error)
	AddRuleEngine(ruleEngine models.RuleEngine) (string, error)
	UpdateRuleEngine(ruleEngine models.RuleEngine) error
	RuleEngineById(id string) (ruleEngine models.RuleEngine, edgeXErr error)
	RuleEngineSearch(offset int, limit int, req dtos.RuleEngineSearchQueryRequest) (ruleEngine []models.RuleEngine, count uint32, edgeXErr error)
	RuleEngineStart(id string) error
	RuleEngineStop(id string) error
	DeleteRuleEngineById(id string) error

	LanguageSdkByName(name string) (cloudService models.LanguageSdk, edgeXErr error)
	LanguageSearch(offset int, limit int, req dtos.LanguageSDKSearchQueryRequest) (languages []models.LanguageSdk, count uint32, edgeXErr error)
	AddLanguageSdk(cs models.LanguageSdk) (language models.LanguageSdk, edgeXErr error)
	UpdateLanguageSdk(ls models.LanguageSdk) error

	DeviceAlert
	UserDB
	Scene
	SystemMonitor
}

type SystemMonitor interface {
	UpdateSystemMetrics(stats dtos.SystemMetrics) error
	GetSystemMetrics(start, end int64) ([]dtos.SystemMetrics, error)
	RemoveRangeSystemMetrics(min, max string) error
}

type UserDB interface {
	GetUserByUserName(username string) (models.User, error)
	UpdateUser(u models.User) error
	AddUser(u models.User) (models.User, error)
}

type DeviceAlert interface {
	AddAlertRule(rule models.AlertRule) (models.AlertRule, error)
	UpdateAlertRule(rule models.AlertRule) error
	AlertRuleById(id string) (models.AlertRule, error)
	AlertRuleSearch(offset int, limit int, req dtos.AlertRuleSearchQueryRequest) (alertRules []models.AlertRule, total uint32, edgeXErr error)
	DeleteAlertRuleById(id string) error
	AlertRuleStart(id string) error
	AlertRuleStop(id string) error

	AlertListLastSend(alertRuleId string) (alertList models.AlertList, edgeXErr error)
	AddAlertList(alertRule models.AlertList) (models.AlertList, error)
	AlertPlate(beforeTime int64) (plate []dtos.AlertPlateQueryResponse, err error)
	AlertListSearch(offset int, limit int, req dtos.AlertSearchQueryRequest) (alertList []dtos.AlertSearchQueryResponse, total uint32, edgeXErr error)
	AlertIgnore(id string) (edgeXErr error)
	TreatedIgnore(id, message string) (edgeXErr error)
}

type Scene interface {
	AddScene(scene models.Scene) (models.Scene, error)
	SceneById(id string) (models.Scene, error)
	UpdateScene(scene models.Scene) error
	SceneStart(id string) error
	SceneStop(id string) error
	DeleteSceneById(id string) error
	SceneSearch(offset int, limit int, req dtos.SceneSearchQueryRequest) (scenes []models.Scene, total uint32, edgeXErr error)

	AddSceneLog(sceneLog models.SceneLog) (models.SceneLog, error)
	SceneLogSearch(offset int, limit int, req dtos.SceneLogSearchQueryRequest) (sceneLogs []models.SceneLog, total uint32, edgeXErr error)
}
