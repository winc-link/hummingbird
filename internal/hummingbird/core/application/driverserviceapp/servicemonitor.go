package driverserviceapp

import (
	"context"
	"github.com/winc-link/hummingbird/internal/dtos"
	interfaces "github.com/winc-link/hummingbird/internal/hummingbird/core/interface"
	"github.com/winc-link/hummingbird/internal/pkg/constants"
	"github.com/winc-link/hummingbird/internal/pkg/di"
	"github.com/winc-link/hummingbird/internal/pkg/logger"
	"github.com/winc-link/hummingbird/internal/pkg/middleware"
	"sync"
	"time"

	//"gitlab.com/tedge/edgex/internal/tedge/resource/interfaces"
	//
	//"gitlab.com/tedge/edgex/internal/dtos"
	//"gitlab.com/tedge/edgex/internal/pkg/constants"
	pkgContainer "github.com/winc-link/hummingbird/internal/pkg/container"
	//"gitlab.com/tedge/edgex/internal/pkg/di"
	//"gitlab.com/tedge/edgex/internal/pkg/logger"
	//"gitlab.com/tedge/edgex/internal/pkg/middleware"
	resourceContainer "github.com/winc-link/hummingbird/internal/hummingbird/core/container"
)

/**
检测驱动实例运行状态，通过rpc
*/
type DeviceServiceMonitor struct {
	ds        dtos.DeviceService
	isRunning bool
	ctx       context.Context
	dic       *di.Container
	lc        logger.LoggingClient
	exitChan  chan struct{}
	mutex     sync.RWMutex
}

func NewDeviceServiceMonitor(ctx context.Context, ds dtos.DeviceService, dic *di.Container) *DeviceServiceMonitor {
	dsm := &DeviceServiceMonitor{
		ds:        ds,
		isRunning: false,
		ctx:       ctx,
		dic:       dic,
		lc:        pkgContainer.LoggingClientFrom(dic.Get),
		exitChan:  make(chan struct{}),
	}
	go dsm.monitor()
	return dsm
}

func (dsm *DeviceServiceMonitor) monitor() {
	// 监控间隔
	tickTime := time.Second * 5
	timeTickerChan := time.Tick(tickTime)
	for {
		select {
		case <-dsm.ctx.Done():
			dsm.lc.Infof("close to DeviceServiceMonitor dsId: %s", dsm.ds.Id)
			return
		case <-dsm.exitChan:
			dsm.lc.Infof("close to DeviceServiceMonitor dsId: %s", dsm.ds.Id)
			return
		case <-timeTickerChan:
			dsm.CheckServiceAvailable()
		}
	}
}

func (dsm *DeviceServiceMonitor) CheckServiceAvailable() {
	dsm.mutex.Lock()
	defer dsm.mutex.Unlock()
	ctx := middleware.WithCorrelationId(context.Background())
	dsApp := resourceContainer.DriverServiceAppFrom(dsm.dic.Get)
	_, err := dsApp.Get(dsm.ctx, dsm.ds.Id)
	if err != nil {
		dsm.lc.Infof("monitor get driver instance err %+v", err.Error())
		dsm.exitChan <- struct{}{}
		return
	}

	isRunning := interfaces.DMIFrom(dsm.dic.Get).InstanceState(dsm.ds)

	// 驱动运行状态更改
	if dsm.isRunning != isRunning {
		dsm.lc.Debugf("id [%s] before status [%v] current status [%v]: %v", dsm.ds.Id, dsm.isRunning, isRunning, middleware.FromContext(ctx))
		// 状态更改上报
	}

	// 更新管理驱动实例状态
	if !dsApp.InProgress(dsm.ds.Id) {
		if isRunning {
			dsApp.SetState(dsm.ds.Id, constants.RunStatusStarted)
		} else {
			dsApp.SetState(dsm.ds.Id, constants.RunStatusStopped)
		}
	}
	if dsm.isRunning && !isRunning {
		//OfflineDevicesByServiceId(ctx, dsm.dic, dsm.ds.Id)
	}
	dsm.isRunning = isRunning
}

func (dsm *DeviceServiceMonitor) Stop() {
	close(dsm.exitChan)
}
