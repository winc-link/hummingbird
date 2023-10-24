package docker

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	"github.com/winc-link/hummingbird/internal/pkg/utils"
	"io"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/client"
	flag "github.com/spf13/pflag"
	"github.com/winc-link/hummingbird/internal/dtos"
	"github.com/winc-link/hummingbird/internal/pkg/constants"
	"github.com/winc-link/hummingbird/internal/pkg/di"
	"github.com/winc-link/hummingbird/internal/pkg/errort"
	"github.com/winc-link/hummingbird/internal/pkg/logger"

	dicContainer "github.com/winc-link/hummingbird/internal/pkg/container"
)

const (
	containerInternalMountPath  = "/mnt"
	containerInternalConfigPath = "/etc/driver/res/configuration.toml"
)

type DockerManager struct {
	// 镜像repoTags:imageinfo
	ImageMap map[string]ImageInfo
	// 容器name:ContainerInfo
	ContainerMap      map[string]ContainerInfo
	cli               *client.Client
	ctx               context.Context
	timeout           time.Duration
	dic               *di.Container
	lc                logger.LoggingClient
	authToken         string
	mutex             sync.RWMutex
	dcm               *dtos.DriverConfigManage
	defaultRegistries []string
}

type CustomParams struct {
	user       string
	cpuShares  int64
	memory     int64
	memorySwap int64
	dns        []string
	dnsSearch  []string
	restart    string
	env        []string
	runtime    string
	mnt        []string
	port       []string
}

// 镜像信息
type ImageInfo struct {
	Id       string
	RepoTags []string
}

// 容器信息
type ContainerInfo struct {
	Id    string
	Name  string
	State string
}

// NewDockerManager 创建
func NewDockerManager(ctx context.Context, dic *di.Container, dcm *dtos.DriverConfigManage) (*DockerManager, error) {
	cli, err := client.NewClientWithOpts(client.WithVersion(dcm.DockerApiVersion))
	if err != nil {
		return nil, err
	}
	lc := dicContainer.LoggingClientFrom(dic.Get)
	dockerManager := &DockerManager{
		cli:          cli,
		ImageMap:     make(map[string]ImageInfo),
		ContainerMap: make(map[string]ContainerInfo),
		ctx:          context.Background(),
		timeout:      time.Second * 10,
		dic:          dic,
		lc:           lc,
		dcm:          dcm,
	}
	dockerManager.setDefaultRegistry()

	dockerManager.flushImageMap()
	tickTime := time.Second * 10
	timeTickerChanImage := time.Tick(tickTime)
	go func() {
		for {
			select {
			case <-ctx.Done():
				lc.Info("close to flushImageMap")
				return
			case <-timeTickerChanImage:
				dockerManager.flushImageMap()
			}
		}
	}()

	dockerManager.flushContainerMap()
	timeTickerChanContainer := time.Tick(tickTime)
	go func() {
		for {
			select {
			case <-ctx.Done():
				lc.Info("close to flushContainerMap")
				return
			case <-timeTickerChanContainer:
				dockerManager.flushContainerMap()
			}
		}
	}()

	return dockerManager, nil
}

func (dm *DockerManager) setDefaultRegistry() {
	info, err := dm.cli.Info(dm.ctx)
	if err != nil {
		dm.lc.Errorf("get docker info err:%v", err)
		return
	}
	for _, v := range info.RegistryConfig.IndexConfigs {
		dm.defaultRegistries = append(dm.defaultRegistries, v.Name)
	}
}

// 刷新docker镜像数据至内存中
func (dm *DockerManager) flushImageMap() {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()
	images, err := dm.GetImageList()
	if err == nil {
		dm.ImageMap = make(map[string]ImageInfo)
		for _, image := range images {
			if len(image.RepoTags) > 0 {
				dm.ImageMap[image.RepoTags[0]] = ImageInfo{
					Id:       image.ID,
					RepoTags: image.RepoTags,
				}
			}
		}
	}
}

// 刷新docker容器数据至内存中
func (dm *DockerManager) flushContainerMap() {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()
	containers, err := dm.GetContainerList()
	if err == nil {
		dm.ContainerMap = make(map[string]ContainerInfo)
		for _, c := range containers {
			if len(c.Names) > 0 {
				dm.ContainerMap[c.Names[0][1:]] = ContainerInfo{
					Id:    c.ID,
					Name:  c.Names[0][1:],
					State: c.State,
				}
			}
		}
	}
}

// 启动容器,这里为强制重启 先删除容器在启动新容器 containerId 可为空
func (dm *DockerManager) ContainerStart(imageRepo string, containerName string, runConfigPath string, mountDevices []string, customParams string, instanceType constants.InstanceType) (ip string, err error) {
	dm.flushImageMap()
	dm.flushContainerMap()
	if !dm.ExistImageById(imageRepo) {
		return "", errort.NewCommonErr(errort.DeviceLibraryDockerImagesNotFound, fmt.Errorf("imageRepo is %s not exist", imageRepo))
	}

	// 1.停止容器,相同名字容器直接删除
	_ = dm.ContainerStop(containerName)
	_ = dm.ContainerRemove(containerName)

	// 2.创建新容器，配置等
	//exposedPorts, portMap := dm.makeExposedPorts(exposePorts)
	//resourceDevices := dm.makeMountDevices(mountDevices)
	binds := make([]string, 0)
	//binds = append(binds, "/etc/localtime:/etc/localtime:ro") // 挂载时区
	var thisRunMode container.NetworkMode
	if instanceType == constants.CloudInstance {

	} else if instanceType == constants.DriverInstance {
		binds = append(binds, runConfigPath+":"+containerInternalConfigPath)                      //挂载启动配置文件
		binds = append(binds, dm.dcm.GetHostMntDir(containerName)+":"+containerInternalMountPath) //挂载日志
		thisRunMode = container.NetworkMode(dm.dcm.NetWorkName)

	}
	dockerCustomParams, err := dm.ParseCustomParams(customParams)
	if err != nil {
		dm.lc.Errorf("dockerCustomParams err: %+v", err)
		return "", err
	}
	binds = append(binds, dockerCustomParams.mnt...)

	dm.lc.Infof("dockerCustomParams: %+v,%+v", dockerCustomParams, customParams)
	dm.lc.Infof("binds: %+v", binds)
	dm.lc.Infof("Image:%+v", dm.ImageMap[imageRepo])
	dm.lc.Infof("thisRunMode:%+v", string(thisRunMode))

	portMap := generateExposedPorts(dockerCustomParams.port)
	restartPolicy := generateRestartPolicy(dockerCustomParams.restart)
	resources := generateResources(dockerCustomParams.cpuShares, dockerCustomParams.memory, dockerCustomParams.memorySwap)

	var _, cErr = dm.cli.ContainerCreate(dm.ctx, &container.Config{
		OpenStdin: true,
		Tty:       true,
		User:      dockerCustomParams.user,
		Image:     imageRepo,
		Env:       dockerCustomParams.env,
	}, &container.HostConfig{
		DNS:           dockerCustomParams.dns,
		DNSSearch:     dockerCustomParams.dnsSearch,
		Resources:     resources,
		Binds:         binds,
		PortBindings:  portMap,
		NetworkMode:   thisRunMode,
		RestartPolicy: restartPolicy,
		Runtime:       dockerCustomParams.runtime,
	}, &network.NetworkingConfig{}, nil, containerName)
	if cErr != nil {
		return "", cErr
	}

	// 3.启动容器
	if err = dm.cli.ContainerStart(dm.ctx, containerName, types.ContainerStartOptions{}); err != nil {
		return "", errort.NewCommonEdgeX(errort.ContainerRunFail, "Start Container Fail", err)
	}

	// 启动后暂停1秒查看状态
	time.Sleep(time.Second * 1)
	dm.flushContainerMap()

	// 4.查看容器信息并返回相应的数据
	status, err := dm.GetContainerRunStatus(containerName)
	dm.lc.Infof("status: %+v", status)

	if err != nil {
		return "", errort.NewCommonEdgeX(errort.DefaultSystemError, "GetContainerRunStatus Fail", err)
	}
	if status != constants.ContainerRunStatusRunning {
		err = errort.NewCommonEdgeX(errort.ContainerRunFail, fmt.Sprintf("%s container status %s please check the log for specific details", containerName, status), nil)
		return
	}
	if thisRunMode.IsHost() {
		ip = constants.HostAddress
	} else {
		ip, err = dm.GetContainerIp(containerName)
	}
	return
}

func generateRestartPolicy(restart string) container.RestartPolicy {
	if restart != "" {
		ls := strings.Split(restart, ":")
		if len(ls) == 2 {
			maximumRetryCount, err := strconv.Atoi(ls[1])
			if err != nil {
				maximumRetryCount = 0
			}
			return container.RestartPolicy{
				Name:              ls[0],
				MaximumRetryCount: maximumRetryCount,
			}
		} else if len(ls) == 1 {
			return container.RestartPolicy{
				Name: ls[0],
			}
		}
	}
	return container.RestartPolicy{}
}

func generateResources(cpuShares, memory, memorySwap int64) container.Resources {
	return container.Resources{
		CPUShares:  cpuShares,
		Memory:     memory,
		MemorySwap: memorySwap,
	}
}

// makeExposedPorts pots => [8080:8080/tcp 8090:8090/udp]
func generateExposedPorts(ports []string) nat.PortMap {
	portMap := make(nat.PortMap)
	for _, port := range ports {

		proto := "tcp"

		sp := strings.Split(port, ":")
		if len(sp) != 2 {
			return portMap
		}

		parsePortRange := strings.Split(sp[1], "/")

		if len(parsePortRange) == 2 {
			proto = strings.ToLower(parsePortRange[1])
			if proto != "tcp" && proto != "udp" {
				continue
			}
		}
		tmpPort, _ := nat.NewPort(proto, parsePortRange[0])
		portMap[tmpPort] = []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: sp[0],
			},
		}
	}
	return portMap
}

func (dm *DockerManager) ContainerStop(containerIdOrName string) error {
	defer dm.flushContainerMap()
	if err := dm.cli.ContainerStop(dm.ctx, containerIdOrName, &dm.timeout); err != nil {
		dm.lc.Infof("ContainerStop fail %v", err.Error())
		killErr := dm.cli.ContainerKill(dm.ctx, containerIdOrName, "SIGKILL")
		if killErr != nil {
			dm.lc.Infof("ContainerKill fail %v", err.Error())
		}
		return nil
	}
	return nil
}

// 默认为容器强制删除
func (dm *DockerManager) ContainerRemove(containerIdOrName string) error {
	dm.flushContainerMap()
	if err := dm.cli.ContainerRemove(dm.ctx, containerIdOrName, types.ContainerRemoveOptions{Force: true}); err != nil {
		dm.lc.Infof("ContainerRemove fail containerId: %s, err: %v", containerIdOrName, err.Error())
		// 先不用抛出错误
		return nil
	}
	dm.flushContainerMap()
	return nil
}

func (dm *DockerManager) ImageRemove(imageId string) error {
	if imageId == "" {
		return nil
	}
	dm.lc.Infof("doing remove imageId %s", imageId)
	// 错误只做日志，不做抛出，以免影响后续操作
	dm.flushImageMap()
	if _, ok := dm.ImageMap[imageId]; !ok {
		dm.lc.Infof("remove imageId %s is not exist", imageId)
		return nil
	}
	if _, err := dm.cli.ImageRemove(dm.ctx, imageId, types.ImageRemoveOptions{}); err != nil {
		dm.lc.Infof("ImageRemove imageId %s fail %v", imageId, err.Error())
		return nil
	}
	dm.flushImageMap()
	return nil
}

func (dm *DockerManager) GetImageList() ([]types.ImageSummary, error) {
	return dm.cli.ImageList(dm.ctx, types.ImageListOptions{
		All: true,
	})
}

func (dm *DockerManager) GetContainerList() (containers []types.Container, err error) {
	return dm.cli.ContainerList(dm.ctx, types.ContainerListOptions{
		All: true,
	})
}

// 获取容器运行状态， 目前不做任何错误处理
func (dm *DockerManager) GetContainerRunStatus(containerName string) (status string, err error) {
	if len(containerName) == 0 {
		return constants.ContainerRunStatusExited, nil
	}
	dm.mutex.Lock()
	defer dm.mutex.Unlock()
	if _, ok := dm.ContainerMap[containerName]; !ok {
		return constants.ContainerRunStatusExited, nil
	}
	return dm.ContainerMap[containerName].State, nil
}

func (dm *DockerManager) GetContainerIp(containerId string) (ip string, err error) {
	ip = constants.HostAddress
	res, err := dm.cli.ContainerInspect(dm.ctx, containerId)
	if err != nil {
		return
	}
	dm.lc.Infof("Container Networks %+v", res.NetworkSettings.Networks)
	if _, ok := res.NetworkSettings.Networks["bridge"]; ok {
		return res.NetworkSettings.Networks["bridge"].IPAddress, nil
	}
	//dm.lc.Infof("GetContainerIp fail  networks %v", res.NetworkSettings.Networks)
	return ip, err
}

// 获取镜像id是否存在
func (dm *DockerManager) ExistImageById(imageId string) bool {
	if len(imageId) == 0 {
		return false
	}
	if len(dm.checkGetImageId(imageId)) >= 0 {
		return true
	}
	return false
}

// 获取容器id
func (dm *DockerManager) getContainerIdByName(containerName string) string {
	if len(containerName) == 0 {
		return ""
	}
	for _, v := range dm.ContainerMap {
		if v.Name == "/"+containerName {
			return v.Id
		}
	}
	return ""
}

// 获取容器挂载设备的信息,返回挂载设备路径slice ["/dev/ttyUSB0","/dev/ttyUSB0"]
func (dm *DockerManager) GetContainerMountDevices(containerId string) []string {
	resDevices := make([]string, 0)
	if _, ok := dm.ContainerMap[containerId]; !ok {
		dm.lc.Errorf("containerId is %s not exist", containerId)
		return []string{}
	}
	res, err := dm.cli.ContainerInspect(dm.ctx, containerId)
	if err != nil {
		dm.lc.Errorf("containerInspect err %v", err)
		return []string{}
	}
	for _, v := range res.HostConfig.Devices {
		resDevices = append(resDevices, v.PathOnHost)
	}
	return resDevices
}

func (dm *DockerManager) PullDockerImage(imageUrl string, authToken string) (string, error) {
	var resp io.ReadCloser
	var err error
	var dockerImageId string

	// 每次pull都去重新刷新组装token
	dm.lc.Debugf("authToken len %d", len(authToken))
	resp, err = dm.cli.ImagePull(dm.ctx, imageUrl, types.ImagePullOptions{
		RegistryAuth: authToken,
	})
	if err != nil {
		dm.lc.Errorf("url: %s ImagePull err: %+v", imageUrl, err)
		err = errort.NewCommonErr(errort.DeviceLibraryImageDownloadFail, err)
		return dockerImageId, err
	}

	readResp, err := ioutil.ReadAll(resp)
	if err != nil {
		dm.lc.Errorf("url: %s ImagePull err: %+v", imageUrl, err)
		return dockerImageId, err
	}
	dm.lc.Infof("readResp imageUrl %s, %+v", imageUrl, string(readResp))
	re, err := regexp.Compile(`Digest: (\w+:\w+)`)
	if err != nil {
		dm.lc.Errorf("regexp Compile err %v", err)
		return dockerImageId, err
	}
	strSubMatch := re.FindStringSubmatch(string(readResp))
	if len(strSubMatch) < 2 {
		dm.lc.Errorf("regexp not match imagesId")
		return dockerImageId, errort.NewCommonEdgeX(errort.DeviceLibraryImageDownloadFail, "regexp not match imagesId", nil)
	}

	dockerImageId = dm.checkGetImageId(imageUrl)

	if dockerImageId == "" {
		return "", errort.NewCommonEdgeX(errort.DeviceLibraryImageDownloadFail, "docker images is null", nil)
	}

	dm.lc.Infof("images pull success imageId: %s", dockerImageId)
	return dockerImageId, nil
}

func (dm *DockerManager) checkGetImageId(imageUrl string) string {
	dm.flushImageMap()
	for imageId, v := range dm.ImageMap {
		repoTags := v.RepoTags
		if utils.InStringSlice(imageUrl, repoTags) {
			return imageId
		}
		// 补充默认docker url前缀 	如：默认dockerhub的nginx下载下来是image是 nginx:1.12.0 那么补充默认后变成 docker.io/nginx:1.12.0
		for i, _ := range dm.defaultRegistries {
			for j, _ := range repoTags {
				repoTags[j] = dm.defaultRegistries[i] + "/" + repoTags[j]
			}
			if utils.InStringSlice(imageUrl, repoTags) {
				return imageId
			}
		}
	}
	return ""
}

func (dm *DockerManager) ContainerIsExist(containerIdOrName string) bool {
	_, err := dm.cli.ContainerInspect(dm.ctx, containerIdOrName)
	if err != nil {
		return false
	}
	return true
}
func (dm *DockerManager) ContainerRename(originNameOrId string, nowName string) bool {
	err := dm.cli.ContainerRename(dm.ctx, originNameOrId, nowName)
	if err != nil {
		dm.lc.Errorf("ContainerRename %v", err)
		return false
	}
	return true
}

// Deprecated: docker api 的flush接口返回的值无效
func (dm *DockerManager) FlushAuthToken(username string, password string, serverAddress string) {
	dm.lc.Debugf("docker token from serverAddress %s, username %s", serverAddress, username)
	// docker api登陆结果的IdentityToken为空，这里的token为明文组装后base64的值 详见https://github.com/moby/moby/issues/38830
	dm.authToken = base64.StdEncoding.EncodeToString([]byte(`{"username":"` + username + `", "password": "` + password + `", "serveraddress": "` + serverAddress + `"}`))
}

func (dm *DockerManager) GetAuthToken(username string, password string, serverAddress string) string {
	token := base64.StdEncoding.EncodeToString([]byte(`{"username":"` + username + `", "password": "` + password + `", "serveraddress": "` + serverAddress + `"}`))
	dm.lc.Debugf("docker token from serverAddress %s, username %s token %s", serverAddress, username, token)
	return token
}

// 自定义docker启动参数解析
func (dm *DockerManager) ParseCustomParams(cmd string) (CustomParams, error) {
	params := CustomParams{
		user:      "",
		cpuShares: 0,
		memory:    0,
		dns:       []string{},
		dnsSearch: []string{},
		restart:   "",
		env:       []string{},
		runtime:   "",
		port:      []string{},
		mnt:       []string{},
	}
	if cmd == "" {
		return params, nil
	}
	cmd = strings.Replace(cmd, "\n", "", -1)
	strArr := strings.Split(cmd, "\\")
	args := make([]string, 0)
	for _, v := range strArr {
		args = append(args, strings.Split(v, " ")...)
	}
	f := flag.NewFlagSet("", flag.ContinueOnError)
	f.StringVarP(&params.user, "user", "u", "", "")
	f.Int64VarP(&params.cpuShares, "cpu-shares", "c", 0, "")
	f.Int64VarP(&params.memory, "memory", "m", 0, "")
	f.Int64VarP(&params.memorySwap, "memory-swap", "", -1, "")
	f.StringArrayVarP(&params.dns, "dns", "", []string{}, "")
	f.StringArrayVarP(&params.dnsSearch, "dns-search", "", []string{}, "")
	f.StringVarP(&params.restart, "restart", "", "on-failure:10", "")
	f.StringArrayVarP(&params.env, "env", "e", []string{}, "")
	f.StringVarP(&params.runtime, "runtime", "", "", "")
	f.StringArrayVarP(&params.port, "publish", "p", []string{}, "")
	f.StringArrayVarP(&params.mnt, "volume", "v", []string{}, "")

	err := f.Parse(args)
	if err != nil {
		return CustomParams{}, errort.NewCommonErr(errort.DockerParamsParseErr, fmt.Errorf("parse docker params err:%v", err))
	}
	return params, nil

}

func (dm *DockerManager) GetAllImagesIds() []string {
	ids := make([]string, 0)
	tmpMap := dm.ImageMap
	for id, _ := range tmpMap {
		ids = append(ids, id)
	}
	return ids
}

func (dm *DockerManager) GetContainerInspect(containerName string) (types.ContainerJSON, error) {
	return dm.cli.ContainerInspect(dm.ctx, containerName)
}
