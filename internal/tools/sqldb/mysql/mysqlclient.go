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

package mysql

import "gorm.io/gorm"

type DataObject interface {
	TableName() string
	Get() interface{}
}

type ClientSQLite interface {
	CloseSession()

	// 关闭连接
	Close()
	// 初始化建表操作
	InitTable(dos ...DataObject) error
	// 添加数据
	CreateObject(do DataObject) error
	// 判断数据是否存在
	ExistObject(do DataObject) (bool, error)
	// 删除数据
	DeleteObject(do DataObject) error
	// 更新数据
	UpdateObject(do DataObject) error
	// 关联更新
	AssociationsUpdateObject(do DataObject) error
	// 关联删除
	AssociationsDeleteObject(do DataObject) error
	// 查询单个数据
	GetObject(doCond DataObject, do DataObject) error
	// 预加载查询单个数据
	GetPreloadObject(doCond DataObject, do DataObject) error
	// 查询列表数据
	GetObjects(doCond DataObject, do DataObject, likeParam *LikeQueryParam,
		order *OrderQueryParam, offset, limit int) ([]interface{}, int64, error)
	// 批量更新数据
	BatchUpdate(fields, conflictFields, updateFields []string, tableName string, records [][]interface{}) error

	//事物相关
	ExecSqlWithTransaction(funcs ...func(db *gorm.DB) error) error
}
