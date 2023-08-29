//go:build !community
// +build !community

package application

import (
	"context"
	resourceContainer "github.com/winc-link/hummingbird/internal/hummingbird/core/container"
	"github.com/winc-link/hummingbird/internal/pkg/di"
	"github.com/winc-link/hummingbird/internal/pkg/logger"
	"time"

	//"gitlab.com/tedge/edgex/internal/dtos"
	//"gitlab.com/tedge/edgex/internal/pkg/constants"
	//"gitlab.com/tedge/edgex/internal/pkg/container"
	//pkgContainer "github.com/winc-link/hummingbird/internal/pkg/container"
	"github.com/winc-link/hummingbird/internal/pkg/crontab"
	//"gitlab.com/tedge/edgex/internal/pkg/di"
	//"gitlab.com/tedge/edgex/internal/pkg/errort"
	//"gitlab.com/tedge/edgex/internal/pkg/logger"
	//resourceContainer "gitlab.com/tedge/edgex/internal/tedge/resource/container"
	//"gitlab.com/tedge/edgex/internal/tools/atopclient"
)

func InitSchedule(dic *di.Container, lc logger.LoggingClient) {
	lc.Info("init schedule")

	// 每天 1 点
	crontab.Schedule.AddFunc("0 1 * * *", func() {
		lc.Debugf("schedule statistic device msg conut: %v", time.Now().Format("2006-01-02 15:04:05"))
		deviceItf := resourceContainer.DeviceItfFrom(dic.Get)
		err := deviceItf.DevicesReportMsgGather(context.Background())
		if err != nil {
			lc.Error("schedule statistic device err:", err)
		}
	})

	crontab.Start()
}
