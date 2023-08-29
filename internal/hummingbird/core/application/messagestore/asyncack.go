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

import "sync"

type MsgAckChan struct {
	Mu       sync.Mutex
	Id       string
	IsClosed bool
	DataChan chan interface{}
}

func (mac *MsgAckChan) TryCloseChan() {
	mac.Mu.Lock()
	defer mac.Mu.Unlock()
	if !mac.IsClosed {
		close(mac.DataChan)
		mac.IsClosed = true
	}
}

func (mac *MsgAckChan) TrySendDataAndCloseChan(data interface{}) bool {
	mac.Mu.Lock()
	defer mac.Mu.Unlock()
	if !mac.IsClosed {
		mac.DataChan <- data
		close(mac.DataChan)
		mac.IsClosed = true
		return true
	}
	return false
}
