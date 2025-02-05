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

package tdengine

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	interfaces "github.com/winc-link/hummingbird/internal/hummingbird/core/interface"
	"github.com/winc-link/hummingbird/internal/models"
	"github.com/winc-link/hummingbird/internal/pkg/constants"
	"github.com/winc-link/hummingbird/internal/pkg/errort"
	"strconv"
	"strings"
	"time"

	_ "github.com/taosdata/driver-go/v3/taosWS"

	"github.com/gogf/gf/v2/container/gvar"

	"github.com/winc-link/hummingbird/internal/dtos"
	"github.com/winc-link/hummingbird/internal/pkg/logger"

	_ "github.com/taosdata/driver-go/v3/taosRestful"
)

type Client struct {
	client        *sql.DB
	loggingClient logger.LoggingClient
}

func (c *Client) GetDataDBType() constants.DataType {
	return constants.TDengine
}

func (c *Client) CloseSession() {
	c.client.Close()
}

var dbName = "hummingbird"

func NewClient(config dtos.Configuration, lc logger.LoggingClient) (c interfaces.DataDBClient, errEdgeX error) {

	dsn := config.Dsn
	taos, err := sql.Open("taosWS", dsn)

	if err != nil {
		return nil, err
	}
	c = &Client{
		client:        taos,
		loggingClient: lc,
	}
	err = taos.Ping()
	if err != nil {
		return nil, err
	}

	return
}

func (c *Client) Insert(ctx context.Context, table string, data map[string]interface{}) (err error) {
	ts := time.Now().Format("2006-01-02 15:04:05.000")

	var (
		field = []string{"ts"}
		value = []string{"'" + ts + "'"}
	)

	for k, v := range data {
		field = append(field, strings.ToLower(k))
		value = append(value, "'"+gvar.New(v).String()+"'")
	}

	sql := "INSERT INTO ? (?) VALUES (?)"
	_, err = c.client.Exec(sql, table, strings.Join(field, ","), strings.Join(value, ","))
	return

}

func (c *Client) CreateDatabase(ctx context.Context) (err error) {
	var name string
	c.client.QueryRow("SELECT name FROM information_schema.ins_databases WHERE name = '?' LIMIT 1", dbName).Scan(&name)
	if name != "" {
		return
	}
	_, err = c.client.Exec("CREATE DATABASE IF NOT EXISTS " + dbName)

	_, err = c.client.Exec("USE " + dbName)
	return
}

func (c *Client) CreateStable(ctx context.Context, product models.Product) (err error) {
	columns := []string{"ts TIMESTAMP"}

	for _, property := range product.Properties {
		columns = append(columns, c.column(property.TypeSpec.Type, property.Code, property.Name))
	}
	for _, event := range product.Events {
		columns = append(columns, c.column("", event.Code, event.Name))
	}

	for _, action := range product.Actions {
		columns = append(columns, c.column("", action.Code, action.Name))
	}
	sql := fmt.Sprintf("CREATE STABLE IF NOT EXISTS %s.%s (%s) TAGS (device_id NCHAR(255))", dbName, "product_"+product.Id, strings.Join(columns, ","))
	_, err = c.client.Exec(sql)
	return err
}

// CreateTable 创建表
func (c *Client) CreateTable(ctx context.Context, stable, table string) (err error) {
	sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s.device_%s USING %s.%s TAGS('%s')", dbName, table, dbName, "product_"+stable, table)
	fmt.Println(sql)
	_, err = c.client.Exec(sql)
	return
}

func (s *Client) DropStable(ctx context.Context, table string) (err error) {
	sql := fmt.Sprintf("DROP STABLE IF EXISTS %s.%s", dbName, "product_"+table)
	_, err = s.client.Exec(sql)
	return
}

// DropTable 删除子表
func (s *Client) DropTable(ctx context.Context, table string) (err error) {
	sql := fmt.Sprintf("DROP TABLE IF EXISTS %s.%s", dbName, "device_"+table)
	_, err = s.client.Exec(sql)
	return
}

func (c *Client) column(specsType constants.SpecsType, code string, name string) string {
	column := ""
	tdType := ""
	switch specsType {
	case constants.SpecsTypeInt:
		tdType = "INT"
	case constants.SpecsTypeFloat:
		tdType = "FLOAT"
	case constants.SpecsTypeText:
		tdType = "NCHAR(255)"
	case constants.SpecsTypeDate:
		tdType = "TIMESTAMP"
	case constants.SpecsTypeBool:
		tdType = "BOOL"
	default:
		tdType = "NCHAR(255)"
	}
	column = fmt.Sprintf("%s %s", code, tdType)
	return column
}

func (c *Client) AddDatabaseField(ctx context.Context, stableName string, specsType constants.SpecsType, code string, name string) (err error) {
	sql := fmt.Sprintf("ALTER STABLE %s.%s ADD COLUMN %s", dbName, "product_"+stableName, c.column(specsType, code, name))
	_, err = c.client.Exec(sql)
	return
}

func (c *Client) DelDatabaseField(ctx context.Context, stableName, code string) (err error) {
	sql := fmt.Sprintf("ALTER STABLE %s.%s DROP COLUMN %s", dbName, "product_"+stableName, code)
	_, err = c.client.Exec(sql)
	return
}

func (c *Client) ModifyDatabaseField(ctx context.Context, stableName string, specsType constants.SpecsType, code string, name string) (err error) {
	sql := fmt.Sprintf("ALTER STABLE %s.%s MODIFY COLUMN %s", dbName, "product_"+stableName, c.column(specsType, code, name))
	_, err = c.client.Exec(sql)
	if strings.Contains(err.Error(), "column length could be modified") {
		return errort.NewCommonEdgeX(errort.ThingModeTypeCannotBeModified, "Only varbinary/binary/nchar/geometry column length could be modified, and the length can only be increased, not decreased", nil)
	}
	return
}

func (c *Client) GetDeviceService(req dtos.ThingModelServiceDataRequest, device models.Device, product models.Product) ([]dtos.SaveServiceIssueData, int, error) {
	var response []dtos.SaveServiceIssueData
	var count int
	if len(req.Range) == 2 {
		var firstTime, lastTime string
		if req.Range[0] < req.Range[1] {
			firstTime = strconv.Itoa(int(req.Range[0]))
			lastTime = strconv.Itoa(int(req.Range[1]))
		} else {
			firstTime = strconv.Itoa(int(req.Range[1]))
			lastTime = strconv.Itoa(int(req.Range[0]))
		}
		fi, err := strconv.Atoi(firstTime)
		funix := time.UnixMilli(int64(fi))
		if err != nil {
			return response, count, err
		}
		li, err := strconv.Atoi(lastTime)
		if err != nil {
			return response, count, err
		}
		lunix := time.UnixMilli(int64(li))

		if req.Code != "" {
			err = c.client.QueryRow("select count(*) from ? where ts >= '?' and ts <= '?' and ? is not null", "device_"+device.Id, funix.UTC().Format("2006-01-02 15:04:05.000"), lunix.UTC().Format("2006-01-02 15:04:05.000"), strings.ToLower(req.Code)).Scan(&count)
			if err != nil {
				c.loggingClient.Error("query data:", err)
				return response, count, nil
			}

			sql := fmt.Sprintf("select ts,? from ? where ts >= '?' and ts <= '?' and ? is not null order by ts desc limit %d, %d", (req.Page-1)*req.PageSize, req.PageSize)

			rows, err := c.client.Query(sql, strings.ToLower(req.Code), "device_"+device.Id, funix.UTC().Format("2006-01-02 15:04:05.000"), lunix.UTC().Format("2006-01-02 15:04:05.000"), strings.ToLower(req.Code))

			if err != nil {
				c.loggingClient.Error("query data:", err)
				return response, count, nil
			}
			defer rows.Close()
			for rows.Next() {
				columns, _ := rows.Columns()
				values := make([]any, len(columns))
				var rs gdb.Record
				rs = make(gdb.Record, len(columns))
				for i := range values {
					values[i] = new(any)
				}

				err = rows.Scan(values...)
				if err != nil {
					return nil, count, err
				}
				for i, cs := range columns {
					rs[cs] = gvar.New(values[i])
				}
				var reportData dtos.SaveServiceIssueData
				err = json.Unmarshal([]byte(rs[strings.ToLower(req.Code)].String()), &reportData)
				if err != nil {
					c.loggingClient.Error("err:", err)
					continue
				}
				response = append(response, reportData)
			}
		} else {
			var subSQLs []string
			var code []string
			for _, action := range product.Actions {
				code = append(code, strings.ToLower(action.Code))
				subSQLs = append(subSQLs, strings.ToLower(action.Code)+" is not null")
			}

			if len(subSQLs) == 0 {
				return response, 0, nil
			}
			codes := strings.Join(code, ",")
			res := strings.Join(subSQLs, " or ")
			sql := fmt.Sprintf("select count(*) from ? where ts >= '?' and ts <= '?' and (%s)", res)

			err = c.client.QueryRow(sql, "device_"+device.Id, funix.UTC().Format("2006-01-02 15:04:05.000"), lunix.UTC().Format("2006-01-02 15:04:05.000")).Scan(&count)
			if err != nil {
				c.loggingClient.Error("query data:", err)
				return response, count, nil
			}

			sql2 := fmt.Sprintf("select %s from ? where ts >= '?' and ts <= '?' and (%s) order by ts desc limit %d, %d", codes, res, (req.Page-1)*req.PageSize, req.PageSize)
			rows, err := c.client.Query(sql2, "device_"+device.Id, funix.UTC().Format("2006-01-02 15:04:05.000"), lunix.UTC().Format("2006-01-02 15:04:05.000"))

			if err != nil {
				c.loggingClient.Error("query data:", err)
				return response, count, nil
			}
			defer rows.Close()
			for rows.Next() {
				columns, _ := rows.Columns()
				values := make([]any, len(columns))
				var rs gdb.Record
				rs = make(gdb.Record, len(columns))
				for i := range values {
					values[i] = new(any)
				}

				err = rows.Scan(values...)
				if err != nil {
					return nil, count, err
				}
				for i, cs := range columns {
					rs[cs] = gvar.New(values[i])
				}
				var reportData dtos.SaveServiceIssueData
				for _, value := range rs {
					if value.String() != "" {
						err = json.Unmarshal([]byte(value.String()), &reportData)
						if err != nil {
							c.loggingClient.Error("err:", err)
							continue
						}
						response = append(response, reportData)
					}
				}

			}
		}
	}
	return response, count, nil

}

func (c *Client) GetDeviceEvent(req dtos.ThingModelEventDataRequest, device models.Device, product models.Product) ([]dtos.EventData, int, error) {
	var response []dtos.EventData
	var count int

	if len(req.Range) == 2 {
		var firstTime, lastTime string
		if req.Range[0] < req.Range[1] {
			firstTime = strconv.Itoa(int(req.Range[0]))
			lastTime = strconv.Itoa(int(req.Range[1]))
		} else {
			firstTime = strconv.Itoa(int(req.Range[1]))
			lastTime = strconv.Itoa(int(req.Range[0]))
		}
		fi, err := strconv.Atoi(firstTime)
		funix := time.UnixMilli(int64(fi))
		if err != nil {
			return response, count, err
		}
		li, err := strconv.Atoi(lastTime)
		if err != nil {
			return response, count, err
		}
		lunix := time.UnixMilli(int64(li))

		if req.EventCode != "" {
			err = c.client.QueryRow("select count(*) from ? where ts >= '?' and ts <= '?' and ? is not null", "device_"+device.Id, funix.UTC().Format("2006-01-02 15:04:05.000"), lunix.UTC().Format("2006-01-02 15:04:05.000"), strings.ToLower(req.EventCode)).Scan(&count)
			if err != nil {
				c.loggingClient.Error("query data:", err)
				return response, count, nil
			}

			//sql := fmt.Sprintf("select ts,? from ? where ts >= '?' and ts <= '?' and ? is not null order by ts desc limit %d, %d", (req.Page-1)*req.PageSize, req.PageSize)
			sql := fmt.Sprintf("select ts,? from ? where ts >= '?' and ts <= '?' and ? is not null order by ts desc limit %d, %d", (req.Page-1)*req.PageSize, req.PageSize)

			rows, err := c.client.Query(sql, strings.ToLower(req.EventCode), "device_"+device.Id, funix.UTC().Format("2006-01-02 15:04:05.000"), lunix.UTC().Format("2006-01-02 15:04:05.000"), strings.ToLower(req.EventCode))

			if err != nil {
				c.loggingClient.Error("query data:", err)
				return response, count, nil
			}
			defer rows.Close()
			for rows.Next() {
				columns, _ := rows.Columns()
				values := make([]any, len(columns))
				var rs gdb.Record
				rs = make(gdb.Record, len(columns))
				for i := range values {
					values[i] = new(any)
				}

				err = rows.Scan(values...)
				if err != nil {
					return nil, count, err
				}
				for i, cs := range columns {
					rs[cs] = gvar.New(values[i])
				}
				var reportData dtos.EventData
				err = json.Unmarshal([]byte(rs[strings.ToLower(req.EventCode)].String()), &reportData)
				if err != nil {
					c.loggingClient.Error("err:", err)
					continue
				}
				response = append(response, reportData)
			}
		} else {
			var subSQLs []string
			var code []string
			for _, event := range product.Events {
				code = append(code, strings.ToLower(event.Code))
				subSQLs = append(subSQLs, strings.ToLower(event.Code)+" is not null")
			}

			if len(subSQLs) == 0 {
				return response, 0, nil
			}
			codes := strings.Join(code, ",")
			res := strings.Join(subSQLs, " or ")
			sql := fmt.Sprintf("select count(*) from ? where ts >= '?' and ts <= '?' and (%s)", res)

			err = c.client.QueryRow(sql, "device_"+device.Id, funix.UTC().Format("2006-01-02 15:04:05.000"), lunix.UTC().Format("2006-01-02 15:04:05.000")).Scan(&count)
			if err != nil {
				c.loggingClient.Error("query data:", err)
				return response, count, nil
			}

			sql2 := fmt.Sprintf("select %s from ? where ts >= '?' and ts <= '?' and (%s) order by ts desc limit %d, %d", codes, res, (req.Page-1)*req.PageSize, req.PageSize)
			rows, err := c.client.Query(sql2, "device_"+device.Id, funix.UTC().Format("2006-01-02 15:04:05.000"), lunix.UTC().Format("2006-01-02 15:04:05.000"))

			if err != nil {
				c.loggingClient.Error("query data:", err)
				return response, count, nil
			}
			defer rows.Close()
			for rows.Next() {
				columns, _ := rows.Columns()
				values := make([]any, len(columns))
				var rs gdb.Record
				rs = make(gdb.Record, len(columns))
				for i := range values {
					values[i] = new(any)
				}

				err = rows.Scan(values...)
				if err != nil {
					return nil, count, err
				}
				for i, cs := range columns {
					rs[cs] = gvar.New(values[i])
				}
				//fmt.Println()
				var reportData dtos.EventData
				for _, value := range rs {
					if value.String() != "" {
						err = json.Unmarshal([]byte(value.String()), &reportData)
						if err != nil {
							c.loggingClient.Error("err:", err)
							continue
						}
						response = append(response, reportData)
					}
				}

			}
		}

	}
	return response, count, nil
}

func (c *Client) GetDeviceMsgCountByGiveTime(deviceId string, startTime, endTime int64) (int, error) {
	var count int
	err := c.client.QueryRow("select count(*) from ? where ts >= '?' and ts <= '?'", "hummingbird_"+deviceId, time.Unix(startTime, 0).Format("2006-01-02 15:04:05.000"), time.Unix(endTime, 0).Format("2006-01-02 15:04:05.000")).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (c *Client) GetDeviceProperty(req dtos.ThingModelPropertyDataRequest, device models.Device) ([]dtos.ReportData, int, error) {
	var response []dtos.ReportData
	var count int
	if len(req.Range) == 2 {
		var firstTime, lastTime string
		if req.Range[0] < req.Range[1] {
			firstTime = strconv.Itoa(int(req.Range[0]))
			lastTime = strconv.Itoa(int(req.Range[1]))
		} else {
			firstTime = strconv.Itoa(int(req.Range[1]))
			lastTime = strconv.Itoa(int(req.Range[0]))
		}
		fi, err := strconv.Atoi(firstTime)
		funix := time.UnixMilli(int64(fi))
		if err != nil {
			return response, count, err
		}
		li, err := strconv.Atoi(lastTime)
		if err != nil {
			return response, count, err
		}
		lunix := time.UnixMilli(int64(li))

		err = c.client.QueryRow("select count(*) from ? where ts >= '?' and ts <= '?' and ? is not null", "device_"+device.Id, funix.UTC().Format("2006-01-02 15:04:05.000"), lunix.UTC().Format("2006-01-02 15:04:05.000"), strings.ToLower(req.Code)).Scan(&count)
		if err != nil {
			c.loggingClient.Error("query data:", err)
			return []dtos.ReportData{}, count, nil
		}
		var sql string
		if req.IsAll {
			sql = fmt.Sprintf("select ts,? from ? where ts >= '?' and ts <= '?' and ? is not null order by ts desc")
		} else {
			sql = fmt.Sprintf("select ts,? from ? where ts >= '?' and ts <= '?' and ? is not null order by ts desc limit %d, %d", (req.Page-1)*req.PageSize, req.PageSize)
		}

		rows, err := c.client.Query(sql, strings.ToLower(req.Code), "device_"+device.Id, funix.UTC().Format("2006-01-02 15:04:05.000"), lunix.UTC().Format("2006-01-02 15:04:05.000"), strings.ToLower(req.Code))

		if err != nil {
			c.loggingClient.Error("query data:", err)
			return []dtos.ReportData{}, count, nil
		}
		defer rows.Close()
		for rows.Next() {
			columns, _ := rows.Columns()
			values := make([]any, len(columns))
			var rs gdb.Record
			rs = make(gdb.Record, len(columns))
			for i := range values {
				values[i] = new(any)
			}

			err = rows.Scan(values...)
			if err != nil {
				return nil, count, err
			}
			for i, cs := range columns {
				rs[cs] = gvar.New(values[i])
			}
			var reportData dtos.ReportData
			reportData.Time = rs["ts"].Time().UnixMilli()
			if reportData.Time < 0 {
				reportData.Time = 0
			}
			reportData.Value = rs[strings.ToLower(req.Code)].String()
			response = append(response, reportData)
		}
	} else if req.Last {
		sql := "select ts,last(?) as ? from hummingbird.?"
		rows, err := c.client.Query(sql, strings.ToLower(req.Code), strings.ToLower(req.Code), "device_"+device.Id)
		if err != nil {
			return []dtos.ReportData{}, count, nil
		}
		defer rows.Close()
		columns, _ := rows.Columns()
		values := make([]any, len(columns))
		var rs gdb.Record
		rs = make(gdb.Record, len(columns))
		for i := range values {
			values[i] = new(any)
		}

		for rows.Next() {
			err = rows.Scan(values...)
			if err != nil {
				return nil, count, err
			}

			for i, cs := range columns {
				rs[cs] = gvar.New(values[i])
			}
		}
		var reportData dtos.ReportData
		reportData.Time = rs["ts"].Time().UnixMilli()
		if reportData.Time < 0 {
			reportData.Time = 0
		}
		reportData.Value = rs[strings.ToLower(req.Code)].String()
		response = append(response, reportData)
	}

	return response, count, nil
}

func (c *Client) GetOne(ctx context.Context, sql string, args ...any) (rs gdb.Record, err error) {

	rows, err := c.client.Query(sql, args...)
	if err != nil {
		g.Log().Error(ctx, err, sql, args)
		return nil, err
	}
	defer rows.Close()

	columns, _ := rows.Columns()
	values := make([]any, len(columns))
	rs = make(gdb.Record, len(columns))
	for i := range values {
		values[i] = new(any)
	}

	for rows.Next() {
		err = rows.Scan(values...)
		if err != nil {
			return nil, err
		}

		for i, cs := range columns {
			rs[cs] = gvar.New(values[i])
		}

		rows.Close()
	}

	return
}

func (c *Client) GetDevicePropertyCount(request dtos.ThingModelPropertyDataRequest) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Client) GetDeviceEventCount(req dtos.ThingModelEventDataRequest) (int, error) {
	//TODO implement me
	panic("implement me")
}
