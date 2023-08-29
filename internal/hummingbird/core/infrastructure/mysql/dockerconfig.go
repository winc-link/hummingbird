package mysql

import (
	"fmt"
	"github.com/winc-link/hummingbird/internal/dtos"
	"github.com/winc-link/hummingbird/internal/models"
	"github.com/winc-link/hummingbird/internal/pkg/errort"
	"github.com/winc-link/hummingbird/internal/pkg/utils"
	"github.com/winc-link/hummingbird/internal/tools/sqldb/sqlite"

	"gorm.io/gorm"
	//"gitlab.com/tedge/edgex/internal/dtos"
	//"gitlab.com/tedge/edgex/internal/models"
	//"gitlab.com/tedge/edgex/internal/pkg/errort"
	//"gitlab.com/tedge/edgex/internal/pkg/utils"
	//"gitlab.com/tedge/edgex/internal/tools/sqldb/sqlite"
)

func dockerConfigIdExists(c *Client, id string) (bool, error) {
	exists, err := c.client.ExistObject(&models.DockerConfig{
		Id: id,
	})
	if err != nil {
		return false, err
	}
	return exists, nil
}

func dockerConfigAdd(c *Client, dc models.DockerConfig) (models.DockerConfig, error) {
	exists, edgeXErr := dockerConfigIdExists(c, dc.Id)
	if edgeXErr != nil {
		return dc, edgeXErr
	} else if exists {
		return dc, errort.NewCommonEdgeX(errort.DefaultResourcesRepeat, fmt.Sprintf("docker config %s exists", dc.Id), edgeXErr)
	}
	ts := utils.MakeTimestamp()
	if dc.Created == 0 {
		dc.Created = ts
	}
	dc.Modified = ts

	err := c.client.CreateObject(&dc)
	if err != nil {
		edgeXErr = errort.NewCommonEdgeX(errort.DefaultSystemError, "func point creation failed", err)
		return dc, edgeXErr
	}

	return dc, edgeXErr
}

func dockerConfigById(c *Client, id string) (dc models.DockerConfig, edgeXErr error) {
	if id == "" {
		return dc, errort.NewCommonEdgeX(errort.DefaultIdEmpty, "docker config id is empty", nil)
	}
	err := c.client.GetObject(&models.DockerConfig{Id: id}, &dc)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return dc, errort.NewCommonErr(errort.DockerConfigNotExist, fmt.Errorf("docker config id(%s)not found", id))
		}
		return dc, err
	}
	return
}

func dockerConfigDeleteById(c *Client, id string) error {
	if id == "" {
		return errort.NewCommonEdgeX(errort.DefaultIdEmpty, "docker config id is empty", nil)
	}
	rawErr := c.client.DeleteObject(&models.DockerConfig{Id: id})
	if rawErr != nil {
		return errort.NewCommonEdgeX(errort.DefaultSystemError, "docker config deletion failed", rawErr)
	}
	return nil
}

func dockerConfigUpdate(c *Client, dc models.DockerConfig) error {
	dc.Modified = utils.MakeTimestamp()
	err := c.client.UpdateObject(&dc)
	if err != nil {
		return err
	}
	return nil
}

func dockerConfigsSearch(c *Client, offset int, limit int, req dtos.DockerConfigSearchQueryRequest) (dcs []models.DockerConfig, count uint32, edgeXErr error) {
	d := models.DockerConfig{}
	var total int64
	tx := c.Pool.Table(d.TableName())
	tx = sqlite.BuildCommonCondition(tx, d, req.BaseSearchConditionQuery)
	if req.Address != "" {
		tx = tx.Where("`address` = ?", req.Address)
	}
	if req.Account != "" {
		tx = tx.Where("`account` = ?", req.Account)
	}
	err := tx.Count(&total).Error
	if err != nil {
		return []models.DockerConfig{}, 0, errort.NewCommonEdgeX(errort.DefaultSystemError, "device failed query from the database", err)
	}

	err = tx.Offset(offset).Limit(limit).Find(&dcs).Error
	if err != nil {
		return []models.DockerConfig{}, 0, errort.NewCommonEdgeX(errort.DefaultSystemError, "device failed query from the database", err)
	}

	return dcs, uint32(total), nil
}
