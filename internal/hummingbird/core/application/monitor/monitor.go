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
	"github.com/winc-link/hummingbird/internal/dtos"
	"github.com/winc-link/hummingbird/internal/hummingbird/core/container"
	"github.com/winc-link/hummingbird/internal/pkg/constants"
	pkgcontainer "github.com/winc-link/hummingbird/internal/pkg/container"
	"github.com/winc-link/hummingbird/internal/pkg/di"
	"github.com/winc-link/hummingbird/internal/pkg/logger"
	"sync"
	"time"
)

type monitor struct {
	dic               *di.Container
	lc                logger.LoggingClient
	serviceMonitorMap sync.Map
	ctx               context.Context
	exitCh            chan struct{}

	systemMonitor *systemMonitor
}

func NewMonitor(ctx context.Context, dic *di.Container) *monitor {
	lc := pkgcontainer.LoggingClientFrom(dic.Get)

	m := monitor{
		dic:               dic,
		lc:                lc,
		serviceMonitorMap: sync.Map{},
		ctx:               context.Background(),
		exitCh:            make(chan struct{}),
		systemMonitor:     NewSystemMonitor(dic, lc),
	}

	//go m.run()

	return &m
}

func systemMetricsTypeToTime(t string) (time.Time, time.Time) {
	switch t {
	case constants.HourMetricsType:
		end := time.Now()
		start := time.Now().Add(-1 * time.Hour)
		return start, end
	case constants.HalfDayMetricsType:
		end := time.Now()
		start := time.Now().Add(-12 * time.Hour)
		return start, end
	case constants.DayMetricsType:
		end := time.Now()
		start := time.Now().Add(-24 * time.Hour)
		return start, end
	default:
		end := time.Now()
		start := time.Now().Add(-1 * time.Hour)
		return start, end
	}
}

func (m *monitor) GetSystemMetrics(ctx context.Context, query dtos.SystemMetricsQuery) (dtos.SystemMetricsResponse, error) {
	dbClient := container.DBClientFrom(m.dic.Get)

	start, end := systemMetricsTypeToTime(query.MetricsType)
	metrics, err := dbClient.GetSystemMetrics(start.UnixMilli(), end.UnixMilli())
	if err != nil {
		return dtos.SystemMetricsResponse{}, err
	}

	resp := dtos.SystemMetricsResponse{
		Metrics: make([]dtos.SystemStatResponse, 0),
	}
	step := 1
	if query.MetricsType == constants.HalfDayMetricsType || query.MetricsType == constants.DayMetricsType {
		step = 5
	}
	for i := 0; i < len(metrics); i = i + step {
		metric := metrics[i]

		item := dtos.SystemStatResponse{
			Timestamp:         metric.Timestamp,
			CpuUsedPercent:    metric.CpuUsedPercent,
			MemoryTotal:       metric.Memory.Total,
			MemoryUsed:        metric.Memory.Used,
			MemoryUsedPercent: metric.Memory.UsedPercent,
			DiskTotal:         metric.Disk.Total,
			DiskUsed:          metric.Disk.Used,
			DiskUsedPercent:   metric.Disk.UsedPercent,
			Openfiles:         metric.Openfiles,
		}

		if iface, ok := metric.Network[query.Iface]; ok {
			item.NetSentBytes = iface.BytesSentPre
			item.NetRecvBytes = iface.BytesRecvPre
		}

		resp.Metrics = append(resp.Metrics, item)
	}
	resp.Total = len(resp.Metrics)

	return resp, nil
}
