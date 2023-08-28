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
package leveldb

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
	"github.com/winc-link/hummingbird/internal/dtos"
	interfaces "github.com/winc-link/hummingbird/internal/hummingbird/core/interface"
	"github.com/winc-link/hummingbird/internal/models"
	"github.com/winc-link/hummingbird/internal/pkg/constants"
	"github.com/winc-link/hummingbird/internal/pkg/errort"
	"github.com/winc-link/hummingbird/internal/pkg/logger"
	"os"
	"path/filepath"
	"strconv"

	"github.com/syndtr/goleveldb/leveldb"
)

type Client struct {
	client        *LevelDB
	loggingClient logger.LoggingClient
}

func (c *Client) AddDatabaseField(ctx context.Context, tableName string, specsType constants.SpecsType, code string, name string) (err error) {
	return nil
}

func (c *Client) DelDatabaseField(ctx context.Context, tableName, code string) (err error) {
	return nil
}

func (c *Client) ModifyDatabaseField(ctx context.Context, tableName string, specsType constants.SpecsType, code string, name string) (err error) {
	return nil
}

func (c *Client) GetDataDBType() constants.DataType {
	return constants.LevelDB
}

func (c *Client) DropStable(ctx context.Context, table string) (err error) {
	return nil
}

func (c *Client) CreateStable(ctx context.Context, product models.Product) (err error) {
	return nil
}

func (c *Client) DropTable(ctx context.Context, table string) (err error) {
	return nil
}

func (c *Client) CreateTable(ctx context.Context, stable, table string) (err error) {
	return nil
}

func NewClient(config dtos.Configuration, lc logger.LoggingClient) (c interfaces.DataDBClient, errEdgeX error) {
	dataSourceDir := filepath.Dir(config.DataSource)
	_, fileErr := os.Stat(dataSourceDir)
	if fileErr != nil || !os.IsExist(fileErr) {
		_ = os.MkdirAll(dataSourceDir, os.ModePerm)
	}
	var (
		client *leveldb.DB
		err    error
	)
	if client, err = leveldb.OpenFile(dataSourceDir, nil); err != nil {
		if client, err = leveldb.RecoverFile(dataSourceDir, nil); err != nil {
			errEdgeX = errort.NewCommonEdgeX(errort.KindDatabaseError, "database failed to init", err)
			return
		}
	}

	ldb := &LevelDB{
		DB: client,
	}
	c = &Client{
		client:        ldb,
		loggingClient: lc,
	}

	return
}

func (c *Client) CloseSession() {
	c.client.Close()
}

func (c *Client) Insert(ctx context.Context, table string, data map[string]interface{}) (err error) {
	batch := new(leveldb.Batch)
	defer batch.Reset()

	for k, v := range data {
		b, ok := v.([]byte)
		if ok {
			batch.Put([]byte(k), b)
		}
	}
	if err = c.client.Write(batch, &opt.WriteOptions{
		//NoWriteMerge: true,
		//Sync:         true,
	}); err != nil {
		return errort.NewCommonEdgeX(errort.KindDatabaseError, "batch transaction write", err)
	}
	return nil
}

func (c *Client) GetDeviceService(req dtos.ThingModelServiceDataRequest, device models.Device, product models.Product) ([]dtos.SaveServiceIssueData, int, error) {
	var baseKey string
	var response []dtos.SaveServiceIssueData
	var count int
	if req.DeviceId == "" {
		return response, count, fmt.Errorf("deviceId is nill")
	}
	if req.Code == "" {
		baseKey = req.DeviceId + "-" + constants.Action + "-"
	} else {
		baseKey = req.DeviceId + "-" + constants.Action + "-" + req.Code + "-"
	}
	if len(req.Range) == 2 {
		var firstTime, lastTime string
		if req.Range[0] < req.Range[1] {
			firstTime = strconv.Itoa(int(req.Range[0]))
			lastTime = strconv.Itoa(int(req.Range[1]))
		} else {
			firstTime = strconv.Itoa(int(req.Range[1]))
			lastTime = strconv.Itoa(int(req.Range[0]))
		}
		iter := c.client.NewIterator(&util.Range{Start: []byte(baseKey + firstTime), Limit: []byte(baseKey + lastTime)}, &opt.ReadOptions{
			DontFillCache: true,
		})

		var start int
		var end int

		start = (req.Page-1)*req.PageSize + 1
		end = (req.Page-1)*req.PageSize + req.PageSize

		if iter.Last() {
			count++
			var dbvalue dtos.SaveServiceIssueData
			err := json.Unmarshal(iter.Value(), &dbvalue)
			if err != nil {
				fmt.Println(err)
			}
			response = append(response, dbvalue)
			for iter.Prev() {
				count++
				if count >= start && count <= end {
					var dbvalue dtos.SaveServiceIssueData
					err := json.Unmarshal(iter.Value(), &dbvalue)
					if err != nil {
						fmt.Println(err)
					}
					response = append(response, dbvalue)
				}

			}
		}
		iter.Release()
	}
	return response, count, nil
}

func (c *Client) GetDeviceEventCount(req dtos.ThingModelEventDataRequest) (int, error) {
	var baseKey string
	var count int
	if req.DeviceId == "" {
		return 0, fmt.Errorf("deviceId is nill")
	}
	if req.EventCode == "" {
		baseKey = req.DeviceId + "-" + constants.Event + "-"
	} else {
		baseKey = req.DeviceId + "-" + constants.Event + "-" + req.EventCode + "-"
	}
	if len(req.Range) == 2 {
		var firstTime, lastTime string
		if req.Range[0] < req.Range[1] {
			firstTime = strconv.Itoa(int(req.Range[0]))
			lastTime = strconv.Itoa(int(req.Range[1]))
		} else {
			firstTime = strconv.Itoa(int(req.Range[1]))
			lastTime = strconv.Itoa(int(req.Range[0]))
		}

		iter := c.client.NewIterator(&util.Range{Start: []byte(baseKey + firstTime), Limit: []byte(baseKey + lastTime)}, &opt.ReadOptions{
			DontFillCache: true,
		})
		for iter.Next() {
			count++
		}
		iter.Release()
	}
	return count, nil
}

func (c *Client) GetDeviceEvent(req dtos.ThingModelEventDataRequest, device models.Device, product models.Product) ([]dtos.EventData, int, error) {
	var baseKey string
	var response []dtos.EventData
	var count int

	if req.DeviceId == "" {
		return response, count, fmt.Errorf("deviceId is nill")
	}
	if req.EventCode == "" {
		baseKey = req.DeviceId + "-" + constants.Event + "-"
	} else {
		baseKey = req.DeviceId + "-" + constants.Event + "-" + req.EventCode + "-"
	}
	if len(req.Range) == 2 {
		var firstTime, lastTime string
		if req.Range[0] < req.Range[1] {
			firstTime = strconv.Itoa(int(req.Range[0]))
			lastTime = strconv.Itoa(int(req.Range[1]))
		} else {
			firstTime = strconv.Itoa(int(req.Range[1]))
			lastTime = strconv.Itoa(int(req.Range[0]))
		}

		iter := c.client.NewIterator(&util.Range{Start: []byte(baseKey + firstTime), Limit: []byte(baseKey + lastTime)}, &opt.ReadOptions{
			DontFillCache: true,
		})
		var start int
		var end int

		start = (req.Page-1)*req.PageSize + 1
		end = (req.Page-1)*req.PageSize + req.PageSize

		if iter.Last() {
			count++
			var dbvalue dtos.EventData
			_ = json.Unmarshal(iter.Value(), &dbvalue)
			response = append(response, dbvalue)
			for iter.Prev() {
				count++
				if count >= start && count <= end {
					var dbvalue dtos.EventData
					_ = json.Unmarshal(iter.Value(), &dbvalue)
					response = append(response, dbvalue)
				}

			}
		}

		iter.Release()

	}
	return response, count, nil
}

func (c *Client) GetDevicePropertyCount(req dtos.ThingModelPropertyDataRequest) (int, error) {
	var baseKey string
	var count int
	if req.DeviceId == "" {
		return 0, fmt.Errorf("deviceId is nill")
	}
	if req.Code == "" {
		baseKey = req.DeviceId + "-" + constants.Property + "-"
	} else {
		baseKey = req.DeviceId + "-" + constants.Property + "-" + req.Code + "-"
	}

	if len(req.Range) == 2 {
		var firstTime, lastTime string
		if req.Range[0] < req.Range[1] {
			firstTime = strconv.Itoa(int(req.Range[0]))
			lastTime = strconv.Itoa(int(req.Range[1]))
		} else {
			firstTime = strconv.Itoa(int(req.Range[1]))
			lastTime = strconv.Itoa(int(req.Range[0]))
		}
		iter := c.client.NewIterator(&util.Range{Start: []byte(baseKey + firstTime), Limit: []byte(baseKey + lastTime)}, &opt.ReadOptions{
			DontFillCache: true,
		})
		for iter.Next() {
			count++
		}
		iter.Release()
	}
	return count, nil
}

func (c *Client) GetDeviceProperty(req dtos.ThingModelPropertyDataRequest, device models.Device) ([]dtos.ReportData, int, error) {
	var baseKey string
	var response []dtos.ReportData
	var count int
	if req.DeviceId == "" {
		return response, count, fmt.Errorf("deviceId is nill")
	}
	if req.Code == "" {
		baseKey = req.DeviceId + "-" + constants.Property + "-"
	} else {
		baseKey = req.DeviceId + "-" + constants.Property + "-" + req.Code + "-"
	}

	if len(req.Range) == 2 {
		var firstTime, lastTime string
		if req.Range[0] < req.Range[1] {
			firstTime = strconv.Itoa(int(req.Range[0]))
			lastTime = strconv.Itoa(int(req.Range[1]))
		} else {
			firstTime = strconv.Itoa(int(req.Range[1]))
			lastTime = strconv.Itoa(int(req.Range[0]))
		}

		iter := c.client.NewIterator(&util.Range{Start: []byte(baseKey + firstTime), Limit: []byte(baseKey + lastTime)}, &opt.ReadOptions{
			DontFillCache: true,
		})
		var start int
		var end int

		start = (req.Page-1)*req.PageSize + 1
		end = (req.Page-1)*req.PageSize + req.PageSize

		if iter.Last() {
			count++
			var dbvalue dtos.ReportData
			_ = json.Unmarshal(iter.Value(), &dbvalue)
			response = append(response, dbvalue)
			for iter.Prev() {
				count++
				if count >= start && count <= end {
					var dbvalue dtos.ReportData
					_ = json.Unmarshal(iter.Value(), &dbvalue)
					response = append(response, dbvalue)
				}

			}
		}
		iter.Release()
	} else if req.First {
		iter := c.client.NewIterator(util.BytesPrefix([]byte(baseKey)), &opt.ReadOptions{
			DontFillCache: true,
		})
		if iter.First() {
			var dbvalue dtos.ReportData
			_ = json.Unmarshal(iter.Value(), &dbvalue)
			response = append(response, dbvalue)
		}
		iter.Release()
	} else if req.Last {
		iter := c.client.NewIterator(util.BytesPrefix([]byte(baseKey)), nil)
		if iter.Last() {
			var dbvalue dtos.ReportData
			_ = json.Unmarshal(iter.Value(), &dbvalue)
			response = append(response, dbvalue)
			//break
		}
		iter.Release()
	} else {

	}
	return response, count, nil
}

func (c *Client) GetDeviceMsgCountByGiveTime(deviceId string, startTime, endTime int64) (int, error) {
	//TODO implement me
	panic("implement me")
}
