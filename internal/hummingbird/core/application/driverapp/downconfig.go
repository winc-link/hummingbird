package driverapp

import (
	"context"
	"github.com/google/uuid"
	"github.com/winc-link/hummingbird/internal/dtos"
	"github.com/winc-link/hummingbird/internal/hummingbird/core/container"
	"github.com/winc-link/hummingbird/internal/models"
	"github.com/winc-link/hummingbird/internal/pkg/errort"
	"github.com/winc-link/hummingbird/internal/pkg/utils"
	//"gitlab.com/tedge/edgex/internal/pkg/utils"
	//
	//"github.com/google/uuid"
	//"gitlab.com/tedge/edgex/internal/dtos"
	//"gitlab.com/tedge/edgex/internal/models"
	//"gitlab.com/tedge/edgex/internal/pkg/errort"
	//resourceContainer "gitlab.com/tedge/edgex/internal/tedge/resource/container"
)

// 配置镜像/驱动的账号密码仓库地址， 账号密码可为空
func (app *driverLibApp) DownConfigAdd(ctx context.Context, req dtos.DockerConfigAddRequest) error {
	dc := models.DockerConfig{
		Id:      req.Id,
		Address: req.Address,
	}
	var err error
	// 账号密码为空的情况
	if req.Account == "" || req.Password == "" {
		dc.Account = ""
		dc.Password = ""
	} else {
		dc.Account = req.Account
		dc.SaltKey = generateSaltKey()
		dc.Password, err = utils.EncryptAuthPassword(req.Password, dc.SaltKey)
		if err != nil {
			return err
		}
	}
	return app.DownConfigInternalAdd(dc)
}

func (app *driverLibApp) DownConfigInternalAdd(dc models.DockerConfig) error {
	_, err := app.dbClient.DockerConfigAdd(dc)
	if err != nil {
		return err
	}
	return nil
}

func (app *driverLibApp) DownConfigUpdate(ctx context.Context, req dtos.DockerConfigUpdateRequest) error {
	dbClient := container.DBClientFrom(app.dic.Get)

	if req.Id == "" {
		return errort.NewCommonEdgeX(errort.DefaultReqParamsError, "update req id is required", nil)
	}

	dc, edgeXErr := dbClient.DockerConfigById(req.Id)
	if edgeXErr != nil {
		return edgeXErr
	}

	dtos.ReplaceDockerConfigModelFieldsWithDTO(&dc, req)

	if *req.Password != "" {
		var err error
		dc.SaltKey = generateSaltKey()
		dc.Password, err = utils.EncryptAuthPassword(dc.Password, dc.SaltKey)
		if err != nil {
			return err
		}
	}
	edgeXErr = dbClient.DockerConfigUpdate(dc)
	if edgeXErr != nil {
		return edgeXErr
	}
	return nil
}

func (app *driverLibApp) DownConfigSearch(ctx context.Context, req dtos.DockerConfigSearchQueryRequest) ([]models.DockerConfig, uint32, error) {
	offset, limit := req.BaseSearchConditionQuery.GetPage()

	dcs, total, err := app.dbClient.DockerConfigsSearch(offset, limit, req)
	if err != nil {
		return dcs, 0, err
	}

	return dcs, total, nil
}

func (app *driverLibApp) DownConfigDel(ctx context.Context, id string) error {
	dc, edgeXErr := app.dbClient.DockerConfigById(id)
	if edgeXErr != nil {
		return edgeXErr
	}

	// 判断 此配置是否被 library使用
	_, total, err := app.DeviceLibrariesSearch(ctx, dtos.DeviceLibrarySearchQueryRequest{
		DockerConfigId: id,
	})
	if err != nil {
		return edgeXErr
	}

	if total > 0 {
		return errort.NewCommonEdgeX(errort.DockerConfigMustDeleteDeviceLibrary, "请先删除绑定此配置的驱动", nil)
	}

	err = app.dbClient.DockerConfigDelete(dc.Id)
	if err != nil {
		return err
	}
	return nil
}

// 生成password的salt key
func generateSaltKey() string {
	return uuid.New().String()
}
