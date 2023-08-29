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
package mysql

import (
	"fmt"
	"github.com/winc-link/hummingbird/internal/dtos"
	"github.com/winc-link/hummingbird/internal/models"
	"github.com/winc-link/hummingbird/internal/pkg/errort"
	"github.com/winc-link/hummingbird/internal/pkg/utils"
	"github.com/winc-link/hummingbird/internal/tools/sqldb/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func productNameExist(c *Client, name string) (bool, error) {
	exists, err := c.client.ExistObject(&models.Product{Name: name})
	if err != nil {
		return false, err
	}
	return exists, nil
}

func productIdExist(c *Client, id string) (bool, error) {
	exists, err := c.client.ExistObject(&models.Product{Id: id})
	if err != nil {
		return false, err
	}
	return exists, nil
}

func addProduct(c *Client, ds models.Product) (product models.Product, edgeXErr error) {
	exists, edgeXErr := productNameExist(c, ds.Name)
	if edgeXErr != nil {
		return ds, edgeXErr
	} else if exists {
		return ds, errort.NewCommonEdgeX(errort.DefaultResourcesRepeat, fmt.Sprintf("product name %s exists", ds.Id), edgeXErr)
	}

	ts := utils.MakeTimestamp()
	if ds.Created == 0 {
		ds.Created = ts
	}
	ds.Modified = ts

	err := c.client.CreateObject(&ds)
	if err != nil {
		edgeXErr = errort.NewCommonEdgeX(errort.DefaultSystemError, "product creation failed", err)
	}

	return ds, edgeXErr
}

func productById(c *Client, id string) (product models.Product, edgeXErr error) {
	if id == "" {
		return product, errort.NewCommonEdgeX(errort.DefaultIdEmpty, "product id is empty", nil)
	}
	err := c.client.GetPreloadObject(&models.Product{Id: id}, &product)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return product, errort.NewCommonErr(errort.ProductNotExist, fmt.Errorf("product id(%s) not found", id))
		}
		return product, errort.NewCommonErr(errort.DefaultSystemError, fmt.Errorf("query product fail (Id:%s), %s", product.Id, err))
	}
	return
}

func productByCloudId(c *Client, id string) (product models.Product, edgeXErr error) {
	if id == "" {
		return product, errort.NewCommonEdgeX(errort.DefaultIdEmpty, "product id is empty", nil)
	}
	err := c.client.GetPreloadObject(&models.Product{CloudProductId: id}, &product)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return product, errort.NewCommonErr(errort.ProductNotExist, fmt.Errorf("product id(%s) not found", id))
		}
		return product, errort.NewCommonErr(errort.DefaultSystemError, fmt.Errorf("query product fail (Id:%s), %s", product.Id, err))
	}
	//_ = c.cache.SetProduct(product)
	return
}

func productsSearch(c *Client, offset int, limit int, preload bool, req dtos.ProductSearchQueryRequest) (products []models.Product, count uint32, edgeXErr error) {
	dp := models.Product{}
	var total int64
	tx := c.Pool.Table(dp.TableName())
	tx = sqlite.BuildCommonCondition(tx, dp, req.BaseSearchConditionQuery)

	if req.Name != "" {
		tx = tx.Where("`name` LIKE ?", sqlite.MakeLikeParams(req.Name))
	}

	if req.Platform != "" {
		tx = tx.Where("`platform` = ?", req.Platform)
	}

	if req.CloudInstanceId != "" {
		tx = tx.Where("`cloud_instance_id` = ?", req.CloudInstanceId)
	}

	if req.ProductId != "" {
		tx = tx.Where("`product_id` = ?", req.ProductId)
	}

	err := tx.Count(&total).Error
	if err != nil {
		return []models.Product{}, 0, errort.NewCommonEdgeX(errort.DefaultSystemError, "products failed query from the database", err)
	}
	if preload {
		err = tx.Offset(offset).Preload("Properties").Preload("Events").Preload("Actions").Limit(limit).Find(&products).Error
	} else {
		err = tx.Offset(offset).Limit(limit).Find(&products).Error
	}
	if err != nil {
		return []models.Product{}, 0, errort.NewCommonEdgeX(errort.DefaultSystemError, "products failed query from the database", err)
	}

	return products, uint32(total), nil
}

func batchUpsertProduct(c *Client, d []models.Product) (int64, error) {
	if len(d) <= 0 {
		return 0, nil
	}
	tx := c.Pool.Session(&gorm.Session{FullSaveAssociations: true}).Clauses(
		clause.OnConflict{
			UpdateAll: true,
		}).CreateInBatches(d, sqlite.CreateBatchSize)
	num := tx.RowsAffected
	err := tx.Error
	if err != nil {
		return num, err
	}
	return num, nil
}

func batchSaveProduct(c *Client, d []models.Product) error {
	if len(d) <= 0 {
		return nil
	}
	tx := c.Pool.Session(&gorm.Session{FullSaveAssociations: true}).Clauses(
		clause.OnConflict{
			UpdateAll: true,
		}).Save(d)
	//num := tx.RowsAffected
	err := tx.Error
	if err != nil {
		return err
	}
	return nil
}

func batchDeleteProduct(c *Client, products []models.Product) error {
	d := models.Product{}
	tx := c.Pool.Table(d.TableName())
	err := tx.Delete(&products).Error
	if err != nil {
		return errort.NewCommonEdgeX(errort.DefaultSystemError, "product batchDeleteProduct failed", err)
	}
	return nil
}

func deleteProductById(c *Client, id string) error {
	if id == "" {
		return errort.NewCommonEdgeX(errort.DefaultIdEmpty, "product id is empty", nil)
	}
	err := c.client.DeleteObject(&models.Product{Id: id})
	if err != nil {
		return errort.NewCommonEdgeX(errort.DefaultSystemError, "product deletion failed", err)
	}
	return nil
}

func deleteProductObject(c *Client, product models.Product) error {
	err := c.client.DeleteObject(&product)
	if err != nil {
		return errort.NewCommonEdgeX(errort.DefaultSystemError, "product deletion failed", err)
	}
	return nil
}

func associationsDeleteProductObject(c *Client, product models.Product) error {
	err := c.client.AssociationsDeleteObject(&product)
	if err != nil {
		return errort.NewCommonEdgeX(errort.DefaultSystemError, "product deletion failed", err)
	}
	return nil
}

func updateProduct(c *Client, dl models.Product) error {
	dl.Modified = utils.MakeTimestamp()
	err := c.client.UpdateObject(&dl)
	if err != nil {
		return errort.NewCommonEdgeX(errort.DefaultSystemError, "product update failed", err)
	}
	return nil
}

func associationsUpdateProduct(c *Client, dl models.Product) error {
	dl.Modified = utils.MakeTimestamp()
	err := c.client.AssociationsUpdateObject(&dl)
	if err != nil {
		return errort.NewCommonEdgeX(errort.DefaultSystemError, "product update failed", err)
	}
	return nil
}

func batchDeleteProperties(c *Client, propertiesIds []string) error {
	d := models.Properties{}
	tx := c.Pool.Table(d.TableName())
	err := tx.Delete(d, propertiesIds).Error
	if err != nil {
		return errort.NewCommonEdgeX(errort.DefaultSystemError, "batchDeleteProperties failed", err)
	}
	return nil
}

func batchDeleteSystemProperties(c *Client) error {
	d := models.Properties{}
	tx := c.Pool.Table(d.TableName())
	err := tx.Where(models.Properties{System: true}).Delete(d).Error
	if err != nil {
		return errort.NewCommonEdgeX(errort.DefaultSystemError, "batchDelete system property failed", err)
	}
	return nil
}

func batchInsertSystemProperties(c *Client, p []models.Properties) (int64, error) {
	if len(p) <= 0 {
		return 0, nil
	}
	tx := c.Pool.CreateInBatches(p, sqlite.CreateBatchSize)
	num := tx.RowsAffected
	err := tx.Error
	if err != nil {
		return num, err
	}
	return num, nil
}

func batchDeleteEvents(c *Client, eventIds []string) error {
	d := models.Events{}
	tx := c.Pool.Table(d.TableName())
	err := tx.Delete(d, eventIds).Error
	if err != nil {
		return errort.NewCommonEdgeX(errort.DefaultSystemError, "batchDeleteEvents failed", err)
	}
	return nil
}

func batchDeleteSystemEvents(c *Client) error {
	d := models.Events{}
	tx := c.Pool.Table(d.TableName())
	err := tx.Where(models.Events{System: true}).Delete(d).Error
	if err != nil {
		return errort.NewCommonEdgeX(errort.DefaultSystemError, "batch system Actions failed", err)
	}
	return nil
}

func batchInsertSystemEvents(c *Client, p []models.Events) (int64, error) {
	if len(p) <= 0 {
		return 0, nil
	}
	tx := c.Pool.CreateInBatches(p, sqlite.CreateBatchSize)
	num := tx.RowsAffected
	err := tx.Error
	if err != nil {
		return num, err
	}
	return num, nil
}

func batchDeleteActions(c *Client, actionIds []string) error {
	d := models.Actions{}
	tx := c.Pool.Table(d.TableName())
	err := tx.Delete(d, actionIds).Error
	if err != nil {
		return errort.NewCommonEdgeX(errort.DefaultSystemError, "batchDeleteActions failed", err)
	}
	return nil
}

func batchDeleteSystemActions(c *Client) error {
	d := models.Actions{}
	tx := c.Pool.Table(d.TableName())
	err := tx.Where(models.Actions{System: true}).Delete(d).Error
	if err != nil {
		return errort.NewCommonEdgeX(errort.DefaultSystemError, "batch Delete system Actions failed", err)
	}
	return nil
}

func batchInsertSystemActions(c *Client, p []models.Actions) (int64, error) {
	if len(p) <= 0 {
		return 0, nil
	}
	tx := c.Pool.CreateInBatches(p, sqlite.CreateBatchSize)
	num := tx.RowsAffected
	err := tx.Error
	if err != nil {
		return num, err
	}
	return num, nil
}
