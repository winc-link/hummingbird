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
	"log"
	"sync"

	"github.com/winc-link/hummingbird/internal/dtos"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func BuildCommonCondition(tx *gorm.DB, obj interface{}, req dtos.BaseSearchConditionQuery) *gorm.DB {
	objScheme, err := schema.Parse(obj, &sync.Map{}, schema.NamingStrategy{})
	if err != nil {
		log.Fatal(fmt.Sprintf("schema parse err: %v", err))
		return tx
	}

	if req.Id != "" {
		tx = tx.Where("`id` = ?", req.Id)
	}
	if req.Ids != "" {
		tx = tx.Where("`id` IN ?", dtos.ApiParamsStringToArray(req.Ids))
	}
	if req.LikeId != "" {
		tx = tx.Where("`id` LIKE ?", MakeLikeParams(req.LikeId))
	}

	if req.Name != "" && objScheme.FieldsByName["Name"] != nil {
		tx = tx.Where("`name` = ?", req.Name)
	}
	if req.NameLike != "" && objScheme.FieldsByName["Name"] != nil {
		tx = tx.Where("`name` LIKE ?", MakeLikeParams(req.NameLike))
	}
	if req.OrderBy != "" {
		orderBys := dtos.ApiParamsStringToOrderBy(req.OrderBy)
		for _, v := range orderBys {
			var field *schema.Field
			if objScheme.FieldsByName[v.Key] != nil {
				field = objScheme.FieldsByName[v.Key]
			}
			if objScheme.FieldsByDBName[v.Key] != nil {
				field = objScheme.FieldsByDBName[v.Key]
			}
			if field == nil {
				continue
			}
			do := DescOrder
			if !v.IsDesc {
				do = AscOrder
			}
			tx = tx.Order(fmt.Sprintf("%v %v", field.DBName, do))
		}
	} else {
		tx = tx.Order(fmt.Sprintf("%v %v", OrderFieldCreated, DescOrder))
	}

	return tx
}

func MakeLikeParams(str string) string {
	return "%" + str + "%"
}
