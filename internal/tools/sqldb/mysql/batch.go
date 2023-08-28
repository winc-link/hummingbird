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

import (
	"fmt"
	"strings"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DialectorType string

const (
	DIALECTOR_MYSQL  DialectorType = "mysql"
	DIALECTOR_PGSQL  DialectorType = "postgres"
	DIALECTOR_SQLITE DialectorType = "sqlite"
)

type batchProvider struct {
	TableName      string   `json:"table_name"`
	Fields         []string `json:"fields"`
	ConflictFields []string `json:"conflict_fields,omitempty"`
	UpdateFields   []string `json:"update_fields"`
	BatchAmount    int      `json:"batch_amount"`
}

func (provider *batchProvider) Update(engine *gorm.DB, records [][]interface{}) error {
	var (
		index = 0
		end   int
		err   error
	)
	for index < len(records) {
		end = index + provider.BatchAmount
		if end > len(records) {
			end = len(records)
		}
		if err = provider.load(engine, records[index:end]); err != nil {
			return err
		}
		index = end
	}
	return err
}

func (provider *batchProvider) engineJudge(engine gorm.DB) DialectorType {
	switch engine.Dialector.Name() {
	//case (&postgres.Dialector{}).Name():
	//	return DIALECTOR_PGSQL
	//case (&mysql.Dialector{}).Name():
	//	return DIALECTOR_MYSQL
	case (&sqlite.Dialector{}).Name():
		return DIALECTOR_SQLITE
	default:
		return ""
	}
}

func (provider *batchProvider) constructSQL(records [][]interface{}, dialectorType DialectorType) (string, error) {
	switch dialectorType {
	case DIALECTOR_PGSQL:
		return provider.constructPGSQL(records), nil
	case DIALECTOR_MYSQL:
		return provider.constructMYSQL(records), nil
	case DIALECTOR_SQLITE:
		return provider.constructSQLite(records), nil
	default:
		return "", fmt.Errorf("dialector type is invalid")
	}
}

func (provider *batchProvider) constructMYSQL(records [][]interface{}) string {
	var (
		valueNames        string
		valuePlaceHolder  string
		valuePlaceHolders string
		sql               string
	)
	valueNames = strings.Join(provider.Fields, ", ")
	valuePlaceHolder = strings.Repeat("?,", len(provider.Fields))
	valuePlaceHolder = "(" + valuePlaceHolder[:len(valuePlaceHolder)-1] + "),"
	valuePlaceHolders = strings.Repeat(valuePlaceHolder, len(records))
	valuePlaceHolders = valuePlaceHolders[:len(valuePlaceHolders)-1]
	sql = "insert into " + provider.TableName + " (" + valueNames + ") values" + valuePlaceHolders
	var onDups []string
	sql += " on duplicate key "
	if len(provider.UpdateFields) > 0 {
		for _, field := range provider.UpdateFields {
			onDups = append(onDups, field+"=values("+field+")")
		}
		sql += "update " + strings.Join(onDups, ", ")
	} else {
		sql += "nothing"
	}
	return sql
}

func (provider *batchProvider) constructPGSQL(records [][]interface{}) string {
	var (
		valueNames        string
		valuePlaceHolder  string
		valuePlaceHolders string
		sql               string
	)
	valueNames = strings.Join(provider.Fields, ", ")
	valuePlaceHolder = strings.Repeat("?,", len(provider.Fields))
	valuePlaceHolder = "(" + valuePlaceHolder[:len(valuePlaceHolder)-1] + "),"
	valuePlaceHolders = strings.Repeat(valuePlaceHolder, len(records))
	valuePlaceHolders = valuePlaceHolders[:len(valuePlaceHolders)-1]
	sql = "insert into " + provider.TableName + " (" + valueNames + ") values" + valuePlaceHolders
	if len(provider.ConflictFields) > 0 {
		var onDups []string
		sql += " on conflict(" + strings.Join(provider.ConflictFields, ", ") + ") do "
		if len(provider.UpdateFields) > 0 {
			for _, field := range provider.UpdateFields {
				onDups = append(onDups, field+"=excluded."+field)
			}
			sql += "update set " + strings.Join(onDups, ", ")
		} else {
			sql += "nothing"
		}
	}
	return sql
}

func (provider *batchProvider) constructSQLite(records [][]interface{}) string {
	return provider.constructPGSQL(records)
}

func (provider *batchProvider) load(engine *gorm.DB, records [][]interface{}) error {
	// 定义变量
	var (
		sql  string
		args []interface{}
		err  error
	)
	// 构造sql
	sql, err = provider.constructSQL(records, provider.engineJudge(*engine))
	if err != nil {
		return err
	}
	// 添加值列表
	for _, record := range records {
		args = append(args, record...)
	}
	return engine.Exec(sql, args...).Error
}
