package mysql

import (
	"errors"
	"fmt"
	"github.com/winc-link/hummingbird/internal/dtos"
	"github.com/winc-link/hummingbird/internal/pkg/logger"
	"gorm.io/gorm/clause"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	// 顺序排序
	AscOrder = "asc"
	// 倒序排序
	DescOrder = "desc"

	// 默认排序字段
	DefaultOrderField = "modified"
	OrderFieldCreated = "created"

	// 批量创建的每条sql最大条数
	CreateBatchSize = 500
	// 批量更新的每条sql最大条数
	UpdateBatchSize = 10000
)

type runMode string

const (
	//local  runMode = "local"
	remote = "remote"
)

type GormClient struct {
	Pool *gorm.DB
}

type GormWriter struct {
	lc logger.LoggingClient
}

func (gl *GormWriter) Printf(s string, values ...interface{}) {
	gl.lc.Debugf(s, values...)
}

func (c *GormClient) Close() {
	if c.Pool != nil {
		db, err := c.Pool.DB()
		if err == nil {
			db.Close()
		}
	}
}

func NewGormClient(config dtos.Configuration, lc logger.LoggingClient) (*GormClient, error) {
	dsn := config.Dsn
	pool, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &GormClient{Pool: pool}, nil
}

func (c *GormClient) InitTable(dos ...DataObject) error {
	var tables = make([]interface{}, 0)
	for _, do := range dos {
		tables = append(tables, do)
	}
	return c.Pool.AutoMigrate(tables...)
}

func (c *GormClient) CloseSession() {
	return
}

// CreateObject 添加
func (c *GormClient) CreateObject(object DataObject) (err error) {
	return c.Pool.Table(object.TableName()).Create(object).Error
}

//DeleteObject  删除
func (c *GormClient) DeleteObject(object DataObject) (err error) {
	return c.Pool.Table(object.TableName()).Delete(object).Error
}

//GetObject  获取结果
func (c *GormClient) GetObject(obCond DataObject, object DataObject) (err error) {
	return c.Pool.Table(object.TableName()).Where(obCond).First(object).Error
}

// GetPreloadObject 预加载
func (c *GormClient) GetPreloadObject(obCond DataObject, object DataObject) (err error) {
	return c.Pool.Table(object.TableName()).Preload("Properties").Preload("Events").Preload("Actions").Where(obCond).First(object).Error
}

// UpdateObject 更新数据
func (c *GormClient) UpdateObject(object DataObject) (err error) {
	return c.Pool.Table(object.TableName()).Select("*").Updates(object).Error
}

//AssociationsUpdateObject 关联更新
func (c *GormClient) AssociationsUpdateObject(object DataObject) (err error) {
	return c.Pool.Session(&gorm.Session{FullSaveAssociations: true}).Table(object.TableName()).Select("*").Updates(object).Error
}

//AssociationsDeleteObject 关联删除
func (c *GormClient) AssociationsDeleteObject(object DataObject) (err error) {
	return c.Pool.Session(&gorm.Session{FullSaveAssociations: true}).Table(object.TableName()).Select(clause.Associations).Delete(object).Error
}

// 判断是否存在
func (c *GormClient) ExistObject(do DataObject) (exist bool, err error) {
	var count int64
	template := c.Pool.Table(do.TableName())
	err = template.Where(do).Count(&count).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	if count == 0 {
		return false, nil
	}
	return true, nil
}

// 模糊查询参数
type LikeQueryParam struct {
	Field  string
	Value  string
	Prefix bool
	Suffix bool
}

// 排序查询参数
type OrderQueryParam struct {
	Key  string
	Desc bool
}

// 默认按照最后修改时间进行逆序排序
func NewDefaultOrderQueryParam() *OrderQueryParam {
	return &OrderQueryParam{
		Key:  DefaultOrderField,
		Desc: true,
	}
}

// 条件查询
func (c *GormClient) GetObjects(doCond DataObject, do DataObject, likeParam *LikeQueryParam, order *OrderQueryParam, offset, limit int) (list []interface{}, count int64, err error) {
	template := c.Pool.Table(do.TableName())
	if doCond != nil {
		template = template.Where(doCond)
	}
	if likeParam != nil {
		var value = likeParam.Value
		if likeParam.Prefix {
			value = "%" + value
		}
		if likeParam.Suffix {
			value = value + "%"
		}
		template = template.Where(fmt.Sprintf("%v LIKE ?", likeParam.Field), value)
	}
	if err = template.Count(&count).Error; err != nil {
		return
	}
	if order != nil {
		orderKey := order.Key
		if order.Desc {
			orderKey = fmt.Sprintf("%v %v", orderKey, DescOrder)
		} else {
			orderKey = fmt.Sprintf("%v %v", orderKey, AscOrder)
		}
		template = template.Order(orderKey)
	}
	rows, err := template.Limit(limit).Offset(offset).Rows()
	if err != nil {
		return
	}
	defer rows.Close()

	list = make([]interface{}, 0, count)
	for rows.Next() {
		if err = c.Pool.ScanRows(rows, do); err != nil {
			return
		}
		list = append(list, do.Get())
	}
	return
}

/*
	fields: 			全字段名称列表
	conflictFields: 	唯一字段名称列表
	updateFields: 		更新字段名称列表
	record: 			更新数据列表
*/
func (c *GormClient) BatchUpdate(fields, conflictFields, updateFields []string, tableName string, records [][]interface{}) error {
	if len(records) == 0 {
		return nil
	}
	var provider = &batchProvider{
		TableName:      tableName,
		Fields:         fields,
		ConflictFields: conflictFields,
		UpdateFields:   updateFields,
		BatchAmount:    UpdateBatchSize,
	}
	return provider.Update(c.Pool, records)
}

func (c *GormClient) ExecSqlWithTransaction(funcs ...func(db *gorm.DB) error) error {

	tx := c.Pool.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	for _, f := range funcs {
		err := f(tx)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}
