package route

import (
	"github.com/winc-link/hummingbird/internal/hummingbird/core/container"
	pkgContainer "github.com/winc-link/hummingbird/internal/pkg/container"
	"github.com/winc-link/hummingbird/internal/pkg/di"
	"github.com/winc-link/hummingbird/internal/pkg/utils"
	//"github.com/winc-link/hummingbird/internal/system/monitor/container"

	//"gitlab.com/tedge/edgex/internal/pkg/constants"
	"net/http/httputil"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func LoadWebProxyRoutes(r *gin.Engine, webBuildPath string, dic *di.Container) {
	r.Use(ProxyWeb(r, webBuildPath, dic)).Use(static.ServeRoot("/", webBuildPath))
}

//ProxyServer http proxy
func ProxyServer(c *gin.Context, dic *di.Container) {
	configuration := container.ConfigurationFrom(dic.Get)

	port := strconv.Itoa(configuration.Service.Port)
	addr := configuration.Service.ServerBindAddr + ":" + port

	lc := pkgContainer.LoggingClientFrom(dic.Get)
	parseRootUrl, err := url.Parse("http://" + addr)
	if err != nil {
		lc.Errorf("parse server url err:%v", err)
		c.Data(502, "", []byte("proxy server error"))
	}
	proxy := httputil.NewSingleHostReverseProxy(parseRootUrl)
	proxy.ServeHTTP(c.Writer, c.Request)
}

//ProxyWeb 转发
func ProxyWeb(g *gin.Engine, webBuildPath string, dic *di.Container) gin.HandlerFunc {
	return func(context *gin.Context) {
		ReplaceURLPrefix(context, dic)
		uri := context.Request.URL.Path
		if ok, _ := regexp.MatchString("^(/api/|/v1.0/)", uri); ok {
			ProxyServer(context, dic)
			context.Abort()
			return
		}

		absPath := webBuildPath + context.Request.URL.Path
		if utils.FilePathIsExist(absPath) {
			return
		}
		context.Request.URL.Path = "/"
		// 判断index.html文件是否存在
		indexPath := webBuildPath + "/index.html"
		if !utils.FilePathIsExist(indexPath) {
			context.Data(404, "", []byte("404 not found"))
			context.Abort()
			return
		}
		g.HandleContext(context)
	}
}

func ReplaceURLPrefix(context *gin.Context, dic *di.Container) {
	lc := pkgContainer.LoggingClientFrom(dic.Get)
	//get prefix from env
	prefix := os.Getenv("URLPrefix")
	if prefix == "" {
		return
	}

	prefix = prefix + "/"
	context.Request.URL.Path = strings.ReplaceAll(context.Request.URL.Path, prefix, "")
	lc.Debugf("after replace url path:%s", context.Request.URL.Path)
}
