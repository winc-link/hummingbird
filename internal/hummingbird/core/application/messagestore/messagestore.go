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

package messagestore

import (
	"context"
	"github.com/winc-link/hummingbird/internal/pkg/container"
	"github.com/winc-link/hummingbird/internal/pkg/di"
	"github.com/winc-link/hummingbird/internal/pkg/logger"
	"sync"
)

type MessageStores interface {
	StoreMsgId(id string, ch string)
	LoadMsgChan(id string) (interface{}, bool)
	DeleteMsgId(id string)
	GenAckChan(id string) *MsgAckChan
}

type (
	MessageStore struct {
		logger logger.LoggingClient
		ctx    context.Context
		mutex  sync.Mutex
		wg     *sync.WaitGroup
		ackMap sync.Map
	}
)

func NewMessageStore(dic *di.Container) *MessageStore {
	lc := container.LoggingClientFrom(dic.Get)
	return &MessageStore{
		logger: lc,
	}
}

func (wp *MessageStore) StoreRange() {
	wp.ackMap.Range(func(key, value any) bool {
		return true
	})
}

func (wp *MessageStore) StoreMsgId(id string, ch string) {
	wp.ackMap.Store(id, ch)
}

func (wp *MessageStore) DeleteMsgId(id string) {
	wp.ackMap.Delete(id)
}

func (wp *MessageStore) LoadMsgChan(id string) (interface{}, bool) {
	return wp.ackMap.Load(id)
}

func (wp *MessageStore) GenAckChan(id string) *MsgAckChan {
	ack := &MsgAckChan{
		Id:       id,
		DataChan: make(chan interface{}, 1),
	}
	wp.ackMap.Store(id, ack)
	return ack
}
