package openapi

import (
	"github.com/gorilla/schema"
	"github.com/winc-link/hummingbird/internal/pkg/logger"
	"net/http"
)

const (
	UrlParamSceneId         = "sceneId"
	UrlParamActionId        = "actionId"
	UrlParamStrategyId      = "strategyId"
	UrlParamConditionId     = "conditionId"
	UrlParamJobId           = "jobId"
	UrlParamProductId       = "productId"
	UrlParamCategoryKey     = "categoryKey"
	UrlParamCloudInstanceId = "cloudInstanceId"
	UrlParamDeviceId        = "deviceId"
	UrlParamFuncPointId     = "funcPointId"
	UrlParamDeviceLibraryId = "deviceLibraryId"
	UrlParamDeviceServiceId = "deviceServiceId"
	UrlParamDockerConfigId  = "dockerConfigId"
	UrlParamRuleId          = "ruleId"
	UrlDataResourceId       = "dataResourceId"
	RuleEngineId            = "ruleEngineId"
)

var decoder *schema.Decoder

func init() {
	decoder = schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
}

func urlDecodeParam(obj interface{}, r *http.Request, lc logger.LoggingClient) {
	err := decoder.Decode(obj, r.URL.Query())
	if err != nil {
		lc.Errorf("url decoding err %v", err)
	}
}
