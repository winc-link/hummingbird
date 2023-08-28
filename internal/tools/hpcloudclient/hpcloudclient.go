/*******************************************************************************
 * Copyright 2017 Dell Inc.
 * Copyright (c) 2019 Intel Corporation
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
package hpcloudclient

import (
	"github.com/hpcloud/tail"
	"github.com/winc-link/hummingbird/internal/pkg/logger"
	"sync"
)

type hpcloud struct {
	lock       sync.Mutex
	lc         logger.LoggingClient
	hpcloudMap map[string]*hpcloudClient
}

func NewHpcloud(lc logger.LoggingClient) Hpcloud {
	return &hpcloud{
		lc:         lc,
		hpcloudMap: make(map[string]*hpcloudClient),
	}
}

func (h hpcloud) Add(service, filePath string) *hpcloudClient {
	h.lock.Lock()
	defer h.lock.Unlock()
	h.hpcloudMap[service] = newhpcloudClient(filePath)
	return h.hpcloudMap[service]
}

func (h hpcloud) Get(serviceId string) *hpcloudClient {
	h.lock.Lock()
	defer h.lock.Unlock()
	return h.hpcloudMap[serviceId]
}

//-----------------------------------------------------------------

type hpcloudClient struct {
	tails *tail.Tail
}

func (h hpcloudClient) Read() (chan *tail.Line, error) {
	return h.tails.Lines, nil
}

func (h hpcloudClient) Stop() {
	if h.tails != nil {
		_ = h.tails.Stop()
	}
}

func newhpcloudClient(filePath string) *hpcloudClient {
	config := tail.Config{
		ReOpen:    true,                                    // 重新打开
		Follow:    true,                                    // 是否跟随
		Location:  &tail.SeekInfo{Offset: -150, Whence: 2}, // 从文件的哪个地方开始读
		MustExist: false,                                   // 文件不存在不报错
		Poll:      true,
	}
	tails, err := tail.TailFile(filePath, config)
	if err != nil {
		return nil
	}
	return &hpcloudClient{
		tails: tails,
	}
}
