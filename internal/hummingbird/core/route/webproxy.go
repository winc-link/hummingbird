package route

import (
	"context"
	"fmt"
	"github.com/winc-link/hummingbird/internal/hummingbird/core/container"
	"github.com/winc-link/hummingbird/internal/pkg/color"
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
	fmt.Println(color.Red(string(color.LogoContent)))

	tip()
	fmt.Println(color.Green("Server run at:"))
	fmt.Printf("-  Web:   http://localhost:%d/ \r\n", configuration.WebServer.Port)
	fmt.Println(color.Green("Swagger run at:"))
	fmt.Printf("-  Local:   http://localhost:%d/api/v1/swagger/index.html \r\n", configuration.Service.Port)
	fmt.Printf("%s Enter Control + C Shutdown Server \r\n", time.Now().Format("2006-01-02 15:04:05"))

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

func tip() {
	usageStr := `欢迎使用 ` + color.Green(`Hummingbird物联网平台（社区版） `+"版本：v1.0")
	fmt.Printf("%s \n\n", usageStr)
}
