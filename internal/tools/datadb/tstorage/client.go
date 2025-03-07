package tstorage

import (
	"context"
	tstorage "github.com/nakabonne/tstorage"
	"github.com/winc-link/hummingbird/internal/dtos"
	interfaces "github.com/winc-link/hummingbird/internal/hummingbird/core/interface"
	"github.com/winc-link/hummingbird/internal/models"
	"github.com/winc-link/hummingbird/internal/pkg/constants"
	"github.com/winc-link/hummingbird/internal/pkg/logger"
	"github.com/winc-link/hummingbird/internal/pkg/utils"
	"os"
	"path/filepath"
	"time"
)

type Client struct {
	client        tstorage.Storage
	loggingClient logger.LoggingClient
}

func (c *Client) GetDataDBType() constants.DataType {
	return constants.Tstorage
}

func (c *Client) CloseSession() {
	c.client.Close()
}

func (c *Client) Insert(ctx context.Context, table string, data map[string]interface{}) (err error) {
	metric := table
	var rows []tstorage.Row

	timestamp := time.Now().UnixMilli()

	for code, value := range data {
		var labels []tstorage.Label
		labels = append(labels, tstorage.Label{
			Name: "code", Value: code,
		})

		rows = append(rows, tstorage.Row{
			Metric: metric,
			Labels: labels,
			DataPoint: tstorage.DataPoint{
				Timestamp: timestamp,
				Value:     utils.ConvertToFloat64(value),
			},
		})
	}
	return c.client.InsertRows(rows)
}

func paginate(arr []*tstorage.DataPoint, page, pageSize int) []*tstorage.DataPoint {
	arr = reverseArray(arr)
	// 计算起始索引
	start := (page - 1) * pageSize
	end := start + pageSize

	// 边界检查
	if start >= len(arr) {
		return []*tstorage.DataPoint{}
	}
	if end > len(arr) {
		end = len(arr) // 不能超出数组范围
	}

	return arr[start:end]
}

// reverseArray 翻转数组
func reverseArray(arr []*tstorage.DataPoint) []*tstorage.DataPoint {
	n := len(arr)
	reversed := make([]*tstorage.DataPoint, n)
	for i, v := range arr {
		reversed[n-1-i] = v
	}
	return reversed
}

func (c *Client) GetDeviceProperty(req dtos.ThingModelPropertyDataRequest, device models.Device) ([]dtos.ReportData, int, error) {
	var response []dtos.ReportData
	var count int
	if len(req.Range) == 2 {
		var startTime, endTime int64
		if req.Range[0] < req.Range[1] {
			startTime = req.Range[0]
			endTime = req.Range[1]
		} else {
			startTime = req.Range[1]
			endTime = req.Range[0]
		}
		var labels []tstorage.Label

		labels = append(labels, tstorage.Label{
			Name: "code", Value: req.Code,
		})
		points, err := c.client.Select(constants.DB_PREFIX+device.Id, labels, startTime, endTime)
		if err != nil {
			c.loggingClient.Error("tstorage query data:", err)
			return []dtos.ReportData{}, count, err
		}
		count = len(points)

		paginateRes := paginate(points, req.Page, req.PageSize)
		for _, re := range paginateRes {
			response = append(response, dtos.ReportData{
				Value: re.Value,
				Time:  re.Timestamp,
			})
		}

	} else if req.Last {
		var labels []tstorage.Label
		labels = append(labels, tstorage.Label{
			Name: "code", Value: req.Code,
		})
		var startTime, endTime int64
		// 获取当前时间
		now := time.Now()
		// 计算半小时前的时间
		past := now.Add(-30 * time.Minute)
		// 转换为毫秒时间戳
		endTime = now.UnixMilli()
		startTime = past.UnixMilli()

		points, err := c.client.Select(constants.DB_PREFIX+device.Id, labels, startTime, endTime)
		if err != nil {
			c.loggingClient.Error("tstorage query data:", err)
			return []dtos.ReportData{}, count, nil
		}

		dataPoint := points[len(points)-1]
		var reportData dtos.ReportData
		reportData.Time = dataPoint.Timestamp
		reportData.Value = dataPoint.Value
		response = append(response, reportData)

	}
	return response, count, nil
}

func (c *Client) GetDeviceService(req dtos.ThingModelServiceDataRequest, device models.Device, product models.Product) ([]dtos.SaveServiceIssueData, int, error) {
	// tstorage 不支持服务查询
	var response []dtos.SaveServiceIssueData
	var count int
	return response, count, nil
}

func (c *Client) GetDeviceEvent(req dtos.ThingModelEventDataRequest, device models.Device, product models.Product) ([]dtos.EventData, int, error) {
	// tstorage 不支持事件查询
	var response []dtos.EventData
	var count int
	return response, count, nil
}

func (c *Client) CreateTable(ctx context.Context, stable, table string) (err error) {
	return nil
}

func (c *Client) DropTable(ctx context.Context, table string) (err error) {
	return nil
}

func (c *Client) CreateStable(ctx context.Context, product models.Product) (err error) {
	return nil
}

func (c *Client) DropStable(ctx context.Context, table string) (err error) {
	return nil
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

func (c *Client) GetDevicePropertyCount(request dtos.ThingModelPropertyDataRequest) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Client) GetDeviceEventCount(req dtos.ThingModelEventDataRequest) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Client) GetDeviceMsgCountByGiveTime(deviceId string, startTime, endTime int64) (int, error) {
	var labels []tstorage.Label
	var count int
	labels = append(labels, tstorage.Label{})
	points, err := c.client.Select(constants.DB_PREFIX+deviceId, labels, startTime, endTime)
	if err != nil {
		c.loggingClient.Error("tstorage query data:", err)
		return count, err
	}
	count = len(points)
	return count, nil
}

func NewClient(config dtos.Configuration, lc logger.LoggingClient) (c interfaces.DataDBClient, errEdgeX error) {

	dataSourceDir := filepath.Dir(config.DataSource)
	_, fileErr := os.Stat(dataSourceDir)
	if fileErr != nil || !os.IsExist(fileErr) {
		_ = os.MkdirAll(dataSourceDir, os.ModePerm)
	}
	storage, err := tstorage.NewStorage(
		tstorage.WithTimestampPrecision(tstorage.Milliseconds),
		tstorage.WithDataPath(dataSourceDir),
		//tstorage.WithRetention(365*24*time.Hour),
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		client:        storage,
		loggingClient: lc,
	}, nil
}
