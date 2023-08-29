package route

import (
	"context"
	"github.com/winc-link/hummingbird/internal/hummingbird/core/container"
	pkgContainer "github.com/winc-link/hummingbird/internal/pkg/container"
	"github.com/winc-link/hummingbird/internal/pkg/di"
	"github.com/winc-link/hummingbird/internal/pkg/startup"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

///var/bin/cmd/hummingbird-ui/build
const WebBuildPath = "./cmd/hummingbird-ui/build"
const WebBuildPath2 = "/var/build"

// WebBootstrap contains references to dependencies required by the BootstrapHandler.
type WebBootstrap struct {
	router  *gin.Engine
	AppMode string
}

// NewWebBootstrap is a factory method that returns an initialized WebBootstrap receiver struct.
func NewWebBootstrap() *WebBootstrap {
	// 不做路由日志输出
	g := gin.New()
	g.Use(gin.Recovery(), gin.Logger())
	return &WebBootstrap{
		router: g,
	}
}

// BootstrapHandler fulfills the BootstrapHandler contract and performs initialization needed by the resource service.
func (b *WebBootstrap) BootstrapHandler(ctx context.Context, wg *sync.WaitGroup, _ startup.Timer, dic *di.Container) bool {
	configuration := container.ConfigurationFrom(dic.Get)
	lc := pkgContainer.LoggingClientFrom(dic.Get)

	lc.Infof("start WebBootstrap BootstrapHandler in...")

	//pwd, _ := os.Getwd()
	LoadWebProxyRoutes(b.router, WebBuildPath, dic)

	if configuration.WebServer.Host == "" || configuration.WebServer.Port == 0 {
		lc.Errorf("WebServer Host is null OR port is 0")
		return false
	}
	port := strconv.Itoa(configuration.WebServer.Port)
	addr := configuration.WebServer.Host + ":" + port
	timeout := time.Second * time.Duration(configuration.WebServer.Timeout)
	server := &http.Server{
		Addr:         addr,
		Handler:      b.router,
		WriteTimeout: timeout,
		ReadTimeout:  timeout,
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		<-ctx.Done()
		lc.Info("WebProxy server shutting down")
		_ = server.Shutdown(context.Background())
		lc.Info("WebProxy server shut down")
	}()

	lc.Info("WebProxy server starting (" + addr + ")")

	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
		}()

		err := server.ListenAndServe()
		if err != nil {
			lc.Errorf("WebProxy server failed: %v", err)
			cancel := pkgContainer.CancelFuncFrom(dic.Get)
			cancel() // this will caused the service to stop
		} else {
			lc.Info("WebProxy server stopped")
		}
	}()

	return true
}
