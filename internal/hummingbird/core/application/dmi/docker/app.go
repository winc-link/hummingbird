package docker

import (
	"context"
	"fmt"
	"github.com/docker/docker/errdefs"
	"github.com/winc-link/hummingbird/internal/dtos"
	resourceContainer "github.com/winc-link/hummingbird/internal/hummingbird/core/container"
	interfaces "github.com/winc-link/hummingbird/internal/hummingbird/core/interface"
	"github.com/winc-link/hummingbird/internal/pkg/constants"
	"github.com/winc-link/hummingbird/internal/pkg/container"
	"github.com/winc-link/hummingbird/internal/pkg/di"
	"github.com/winc-link/hummingbird/internal/pkg/errort"
	"github.com/winc-link/hummingbird/internal/pkg/logger"
	"github.com/winc-link/hummingbird/internal/pkg/utils"
	"io/ioutil"
	"path"
	"strings"
	"sync"
)

var _ interfaces.DMI = new(dockerImpl)

type dockerImpl struct {
	dic *di.Container
	lc  logger.LoggingClient
	dm  *DockerManager
	dcm *dtos.DriverConfigManage
}

func (d *dockerImpl) GetDriverInstanceLogPath(serviceName string) string {
	return d.dcm.GetHostLogFilePath(serviceName)
}
func New(dic *di.Container, ctx context.Context, wg *sync.WaitGroup, dcm dtos.DriverConfigManage) (*dockerImpl, error) {
	dm, err := NewDockerManager(ctx, dic, &dcm)
	if err != nil {
		return nil, err
	}
	dil := &dockerImpl{
		dic: dic,
		lc:  container.LoggingClientFrom(dic.Get),
		dm:  dm,
		dcm: &dcm,
	}
	dil.setDcmRootDir()
	return dil, nil
}
func (d *dockerImpl) StopAllInstance() {
	dbClient := resourceContainer.DBClientFrom(d.dic.Get)

	deviceService, _, err := dbClient.DeviceServicesSearch(0, -1, dtos.DeviceServiceSearchQueryRequest{})
	if err != nil {
		d.lc.Errorf("search service error :", err.Error())
	}
	for _, service := range deviceService {
		d.lc.Info(fmt.Sprintf("stop docker instance[%s]", service.ContainerName))
		err := d.StopInstance(dtos.DeviceService{ContainerName: service.ContainerName})
		if err != nil {
			d.lc.Error(fmt.Sprintf("stop docker instance[%s] error:", err.Error()))
		}
	}
}

// DownApp 下载应用
func (d *dockerImpl) DownApp(cfg dtos.DockerConfig, app dtos.DeviceLibrary, toVersion string) (string, error) {
	authToken, err := d.getAuthToken(cfg.Address, cfg.Account, cfg.Password, cfg.SaltKey)
	if err != nil {
		return "", err
	}

	// 2. pull images
	return d.getApp(authToken, d.getImageUrl(cfg.Address, app.DockerRepoName, toVersion))
}

// getAuthToken 获取 docker 认证 token
func (d *dockerImpl) getAuthToken(address, account, pass, salt string) (string, error) {
	if account == "" || pass == "" {
		return "", nil
	}

	// 处理docker密码
	rawPassword, err := utils.DecryptAuthPassword(pass, salt)
	if err != nil {
		d.lc.Errorf("3.getAuthToken docker id:%s, account:%s, password err:%n", address, account, err)
		return "", err
	}
	return d.dm.GetAuthToken(account, rawPassword, address), nil
}

// getApp 下载镜像
func (d *dockerImpl) getApp(token, imageUrl string) (string, error) {
	dockerImageId, dockerErr := d.dm.PullDockerImage(imageUrl, token)
	if dockerErr != nil {
		code := errort.DefaultSystemError
		err := fmt.Errorf("driver library download error")
		if errdefs.IsUnauthorized(dockerErr) || errdefs.IsForbidden(dockerErr) ||
			strings.Contains(dockerErr.Error(), "denied") ||
			strings.Contains(dockerErr.Error(), "unauthorized") {
			code = errort.DeviceLibraryDockerAuthInvalid
			err = fmt.Errorf("docker auth invalid")
		} else if errdefs.IsNotFound(dockerErr) || strings.Contains(dockerErr.Error(), "not found") {
			code = errort.DeviceLibraryDockerImagesNotFound
			err = fmt.Errorf("docker images not found")
		} else if strings.Contains(dockerErr.Error(), "invalid reference format") {
			code = errort.DeviceLibraryDockerImagesNotFound
			err = fmt.Errorf("docker images not found, url invalid")
		}
		d.lc.Errorf("4.getApp imageUrl %s, PullDockerImage err:%v", imageUrl, dockerErr)

		return "", errort.NewCommonErr(code, err)
	}

	return dockerImageId, nil
}

func (d *dockerImpl) getImageUrl(address, repoName, version string) string {
	return path.Clean(address+"/"+repoName) + ":" + version
}

// StateApp 驱动软件下载情况
func (d *dockerImpl) StateApp(dockerImageId string) bool {
	return d.dm.ExistImageById(dockerImageId)
}

// GetAllApp 获取所有镜像信息
func (d *dockerImpl) GetAllApp() []string {
	return d.dm.GetAllImagesIds()
}

func (d *dockerImpl) setDcmRootDir() {
	res, err := d.dm.GetContainerInspect(d.dcm.DockerSelfName)
	if err != nil {
		d.lc.Errorf("GetContainerInspect:%v", err)
		return
	}

	for _, v := range res.Mounts {
		if v.Destination == constants.DockerHummingbirdRootDir {
			d.dcm.SetHostRootDir(v.Source)
			break
		}
	}
	networkName := res.HostConfig.NetworkMode.UserDefined()
	d.dcm.SetNetworkName(networkName)
}

func (d *dockerImpl) genRunServiceConfig(name string, cfgContent string, instanceType constants.InstanceType) (string, error) {
	var err error
	var filePath string
	if instanceType == constants.CloudInstance {

	} else if instanceType == constants.DriverInstance {
		filePath = d.dcm.GetRunConfigPath(name)
		err = utils.CreateDirIfNotExist(constants.DockerHummingbirdRootDir + "/" + constants.DriverRunConfigDir)
		if err != nil {
			return "", err
		}
	}
	err = ioutil.WriteFile(filePath, []byte(cfgContent), 0644)

	if err != nil {
		d.lc.Error("err:", err)
		return "", err
	}
	if instanceType == constants.CloudInstance {

	} else if instanceType == constants.DriverInstance {
		return d.dcm.GetHostRunConfigPath(name), nil
	}
	return "", nil
}
