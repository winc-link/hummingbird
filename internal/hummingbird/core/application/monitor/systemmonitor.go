/*******************************************************************************
 * Copyright 2017.
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

package monitor

import (
	"context"
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
	"github.com/winc-link/hummingbird/internal/dtos"
	"github.com/winc-link/hummingbird/internal/hummingbird/core/container"
	interfaces "github.com/winc-link/hummingbird/internal/hummingbird/core/interface"
	"github.com/winc-link/hummingbird/internal/pkg/di"
	"github.com/winc-link/hummingbird/internal/pkg/logger"
	"github.com/winc-link/hummingbird/internal/pkg/utils"
	"strconv"
	"time"
)

type systemMonitor struct {
	ctx      context.Context
	dic      *di.Container
	dbClient interfaces.DBClient
	lc       logger.LoggingClient
	exitCh   chan struct{}

	ethMap map[string]*dtos.SystemNetwork
	//aa        *AlertApplication
	diskAlert int
	cpuAlert  int
}

func NewSystemMonitor(dic *di.Container, lc logger.LoggingClient) *systemMonitor {
	dbClient := container.DBClientFrom(dic.Get)
	ctx := context.Background()
	m := systemMonitor{
		ctx:      ctx,
		dic:      dic,
		dbClient: dbClient,
		lc:       lc,
		ethMap:   make(map[string]*dtos.SystemNetwork),
		exitCh:   make(chan struct{}),
		//aa:        NewAlertApplication(ctx, dic, lc),
		diskAlert: 3, // 默认告警 3 次
		cpuAlert:  3,
	}

	go m.run()

	return &m
}

func (m *systemMonitor) run() {
	tick := time.Tick(time.Minute)         // 1 minute
	tickClear := time.Tick(24 * time.Hour) // 每24小时删除一次数据
	go func() {
		for {
			select {
			case <-tick:
				metrics := m.collect()
				if err := m.dbClient.UpdateSystemMetrics(metrics); err != nil {
					m.lc.Errorf("failed to UpdateSystemMetrics %v", err)
				}

				//m.reportSystemMetricsAlert(metrics)
			case <-tickClear:
				m.clearMetrics()
			case <-m.exitCh:
				return
			}
		}
	}()
}

func (m *systemMonitor) clearMetrics() {
	min := "0"
	max := strconv.FormatInt(time.Now().Add(-24*time.Hour).UnixMilli(), 10)
	m.lc.Infof("remove system metrics data from %v to %v", min, max)
	if err := m.dbClient.RemoveRangeSystemMetrics(min, max); err != nil {
		m.lc.Error("failed to clearMetrics", err)
	}
}

func (m *systemMonitor) Close() {
	close(m.exitCh)
}

func (m *systemMonitor) collect() dtos.SystemMetrics {
	return dtos.SystemMetrics{
		Timestamp:      time.Now().UnixMilli(),
		CpuUsedPercent: getCpu(),
		CpuAvg:         getCpuLoad(),
		Memory:         getMemory(),
		Network:        getNetwork(m.ethMap),
		Disk:           getDisk(),
		Openfiles:      getOpenfiles(),
	}
}

func getMemory() dtos.SystemMemory {
	v, _ := mem.VirtualMemory()

	return dtos.SystemMemory{
		Total:       v.Total,
		Used:        v.Used,
		UsedPercent: v.UsedPercent,
	}
}

func getCpu() float64 {
	// cpu的使用率
	totalPercent, _ := cpu.Percent(0, false)
	if len(totalPercent) <= 0 {
		return 0
	}
	return totalPercent[0]
}

func getCpuLoad() float64 {
	// cpu的使用率
	avg, _ := load.Avg()
	if avg == nil {
		return 0
	}
	return avg.Load1
}

func getDisk() dtos.SystemDisk {
	// 目录 / 的磁盘使用率
	usage, _ := disk.Usage("/")
	return dtos.SystemDisk{
		Path:        "/",
		Total:       usage.Total,
		Used:        usage.Used,
		UsedPercent: usage.UsedPercent,
	}
}

func getNetwork(ethMap map[string]*dtos.SystemNetwork) map[string]dtos.SystemNetwork {
	stats := make(map[string]dtos.SystemNetwork)

	info, _ := net.IOCounters(true)
	for _, v := range info {
		ethName := v.Name
		if !utils.CheckNetIface(ethName) {
			continue
		}
		if v.BytesSent <= 0 && v.BytesRecv <= 0 {
			continue
		}

		_, ok := ethMap[ethName]
		if !ok {
			ethMap[ethName] = &dtos.SystemNetwork{}
		}
		ethItem := ethMap[ethName]

		var (
			byteRecvPre uint64
			byteSentPre uint64
		)

		now := time.Now().Unix()
		if ethItem.Last == 0 {
			// 第一次采集，没有初始值，不计算
		} else {
			byteRecvPre = v.BytesRecv - ethItem.BytesRecv
			byteSentPre = v.BytesSent - ethItem.BytesSent
		}

		item := dtos.SystemNetwork{
			Name:         ethName,
			BytesSent:    v.BytesSent,
			BytesRecv:    v.BytesRecv,
			BytesRecvPre: byteRecvPre,
			BytesSentPre: byteSentPre,
			Last:         now,
		}
		stats[ethName] = item
		ethMap[ethName] = &item
	}
	return stats
}

func getOpenfiles() int {
	// only linux
	// https://github.com/shirou/gopsutil#process-class
	processes, _ := process.Processes()
	var openfiles int
	for _, pid := range processes {
		files, _ := pid.OpenFiles()
		openfiles += len(files)
	}
	return openfiles
}

func getPlatform() {
	// 查看平台信息
	platform, family, version, _ := host.PlatformInformation()
	fmt.Printf("platform = %v ,family = %v , version = %v \n", platform, family, version)
}
