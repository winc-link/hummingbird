package gateway

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/schema"
	"github.com/winc-link/hummingbird/internal/pkg/errort"
	"github.com/winc-link/hummingbird/internal/pkg/httphelper"
	"github.com/winc-link/hummingbird/internal/pkg/logger"
	"net/http"
	"net/http/httputil"
	"net/url"
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

func (ctl *controller) ProxyAgentServer(c *gin.Context) {
	proxy := ctl.cfg.Clients["Agent"]
	ctl.lc.Infof("agentProxy: %v", proxy)
	ctl.ServeHTTP(c, fmt.Sprintf("http://%v:%v", proxy.Host, proxy.Port))
}

//func (ctl *controller) ProxySharpServer(c *gin.Context) {
//	proxy := ctl.cfg.Clients["Sharp"]
//	ctl.lc.Infof("sharpProxy: %v", proxy)
//	ctl.ServeHTTP(c, fmt.Sprintf("http://%v:%v", proxy.Host, proxy.Port))
//}

func (ctl *controller) ServeHTTP(c *gin.Context, URL string) {
	parseRootUrl, err := url.Parse(URL)
	if err != nil {
		httphelper.RenderFail(c, errort.NewCommonErr(errort.DefaultReqParamsError, err), c.Writer, ctl.lc)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(parseRootUrl)
	proxy.ServeHTTP(c.Writer, c.Request)
}
