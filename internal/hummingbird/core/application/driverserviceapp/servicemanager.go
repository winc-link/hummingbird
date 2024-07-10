package driverserviceapp

import (
	"bytes"
	"context"
	"fmt"
	"github.com/BurntSushi/toml"
	pkgerr "github.com/pkg/errors"
	"github.com/winc-link/hummingbird/internal/dtos"
	"github.com/winc-link/hummingbird/internal/hummingbird/core/container"
	interfaces "github.com/winc-link/hummingbird/internal/hummingbird/core/interface"
	"github.com/winc-link/hummingbird/internal/models"
	"github.com/winc-link/hummingbird/internal/pkg/constants"
	bootstrapContainer "github.com/winc-link/hummingbird/internal/pkg/container"
	"github.com/winc-link/hummingbird/internal/pkg/di"
	"github.com/winc-link/hummingbird/internal/pkg/errort"
	"github.com/winc-link/hummingbird/internal/pkg/logger"
	"github.com/winc-link/hummingbird/internal/pkg/utils"
	"gorm.io/gorm"
	"net"
	"strconv"
	"strings"

	"sync"
)

//  驱动实例管理
func newDriverServiceApp(ctx context.Context, dic *di.Container) *driverServiceAppM {
	dsManager := &driverServiceAppM{
		state:     sync.Map{},
		dic:       dic,
		lc:        bootstrapContainer.LoggingClientFrom(dic.Get),
		dbClient:  container.DBClientFrom(dic.Get),
		ctx:       ctx,
		appModel:  interfaces.DMIFrom(dic.Get),
		dsMonitor: make(map[string]*DeviceServiceMonitor),
	}
	//dsManager.FlushStatsToAgent()
	dsManager.initMonitor()
	return dsManager
}

//
type driverServiceAppM struct {
	state sync.Map
	dic   *di.Container
	lc    logger.LoggingClient
	ctx   context.Context // Bootstrap init 启动传入的， 用来处理done数据

	// interfaces
	dbClient interfaces.DBClient
	appModel interfaces.DMI

	dsMonitor map[string]*DeviceServiceMonitor
}

func (m *driverServiceAppM) getDriverApp() interfaces.DriverLibApp {
	return container.DriverAppFrom(m.dic.Get)
}

func (m *driverServiceAppM) GetState(id string) int {
	state, ok := m.state.Load(id)
	if ok {
		return state.(int)
	}
	m.state.Store(id, constants.RunStatusStopped)

	return constants.RunStatusStopped
}

func (m *driverServiceAppM) SetState(id string, state int) {
	m.state.Store(id, state)
}

func (m *driverServiceAppM) Start(id string) error {
	var err error
	defer func() {
		if err != nil {
			m.SetState(id, constants.RunStatusStopped)
		}
	}()

	if m.InProgress(id) {
		return fmt.Errorf("that id(%s) is staring or stopping, do not to start", id)
	}

	ds, err := m.Get(context.Background(), id)
	if err != nil {
		return err
	}
	dl, err := m.getDriverApp().DriverLibById(ds.DeviceLibraryId)
	if err != nil {
		return err
	}

	driverRunPort, err := utils.GetAvailablePort(ds.GetPort())
	if err != nil {
		return errort.NewCommonErr(errort.CreateConfigFileFail, fmt.Errorf("create cofig file faild %w", err))
	}

	// 获取自身服务运行的ip,并组装运行启动的配置
	runConfig, err := m.buildServiceRunCfg(m.appModel.GetSelfIp(), driverRunPort, ds)
	if err != nil {
		return errort.NewCommonErr(errort.GetAvailablePortFail, fmt.Errorf("get available port fail"))

	}
	dtoDs := dtos.DeviceServiceFromModel(ds)
	dtoRunCfg := dtos.RunServiceCfg{
		ImageRepo:    dl.DockerImageId,
		RunConfig:    runConfig,
		DockerParams: ds.DockerParams,
		DriverName:   dl.Name,
	}
	m.SetState(id, constants.RunStatusStarting)
	_, err = m.appModel.StartInstance(dtoDs, dtoRunCfg)
	if err != nil {
		return err
	}
	m.SetState(id, constants.RunStatusStarted)

	//重新刷新数据
	ds, err = m.Get(context.Background(), id)
	if err != nil {
		return err
	}

	oldBaseAddress := ds.BaseAddress
	// 更新驱动服务数据
	ds.BaseAddress = ds.ContainerName + ":" + strconv.Itoa(driverRunPort)
	// 更新监控ds 如果不更新ping 定时检测会失效
	if oldBaseAddress != ds.BaseAddress {
		err = m.dbClient.UpdateDeviceService(ds)
		if err != nil {
			return err
		}
		if _, ok := m.dsMonitor[ds.Id]; ok {
			m.dsMonitor[ds.Id].ds = dtos.DeviceServiceFromModel(ds)
		}
	}

	return nil
}

func (m *driverServiceAppM) Stop(id string) error {
	ds, err := m.Get(context.Background(), id)
	if err != nil {
		return err
	}

	m.SetState(id, constants.RunStatusStopping)
	stopErr := m.appModel.StopInstance(dtos.DeviceServiceFromModel(ds))
	if stopErr != nil {
		m.SetState(id, constants.RunStatusStopped)
		return errort.NewCommonErr(errort.ContainerStopFail, pkgerr.WithMessage(stopErr, "stop driverService fail"))
	}
	m.SetState(id, constants.RunStatusStopped)
	return nil
}

func (m *driverServiceAppM) ReStart(id string) error {
	err := m.Stop(id)
	if err != nil {
		return fmt.Errorf("dsId(%v), stop err:%v", id, err)
	}
	err = m.Start(id)
	if err != nil {
		return fmt.Errorf("dsId(%v), start err:%v", id, err)
	}
	return nil
}

//
func (m *driverServiceAppM) Add(ctx context.Context, ds models.DeviceService) error {
	if ds.BaseAddress == "" {
		address, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:0", "0.0.0.0"))
		if err != nil {
			return err
		}
		port, _ := utils.AvailablePort(address)
		ds.BaseAddress = ds.ContainerName + ":" + strconv.Itoa(port)
	}

	// 处理专家模式配置
	if ds.ExpertMode && len(ds.ExpertModeContent) > 0 {
		tmpKv, err := dtos.FromYamlStrToMap(ds.ExpertModeContent)
		if err != nil {
			return errort.NewCommonErr(errort.DefaultReqParamsError, fmt.Errorf("parse expertModeContent err:%v", err))
		}
		ds.Config[constants.ConfigKeyDriver] = tmpKv
	}

	ds, err := m.dbClient.AddDeviceService(ds)
	if err != nil {
		return err
	}

	// 添加后台监控
	if _, ok := m.dsMonitor[ds.Id]; ok {
		m.dsMonitor[ds.Id].ds = dtos.DeviceServiceFromModel(ds)
	} else {
		m.dsMonitor[ds.Id] = NewDeviceServiceMonitor(m.ctx, dtos.DeviceServiceFromModel(ds), m.dic)
	}

	//go m.FlushStatsToAgent()
	//go m.autoAddDevice(ds)
	return nil
}

//
func (m *driverServiceAppM) Update(ctx context.Context, dto dtos.DeviceServiceUpdateRequest) error {
	deviceService, edgeXErr := m.Get(ctx, dto.Id)
	if edgeXErr != nil {
		return edgeXErr
	}

	if m.GetState(dto.Id) == constants.RunStatusStarted {
		return errort.NewCommonErr(errort.DeviceServiceMustStopService, fmt.Errorf("service(%v) is running not update", deviceService.Id))
	}
	dtos.ReplaceDeviceServiceModelFieldsWithDTO(&deviceService, dto)
	edgeXErr = m.dbClient.UpdateDeviceService(deviceService)
	if edgeXErr != nil {
		return edgeXErr
	}
	return nil
}

//
// 升级实例： 如果不存在则创建数据、如果存在，但未运行，不做处理、若运行中则重启
func (m *driverServiceAppM) Upgrade(dl models.DeviceLibrary) error {
	dss, _, err := m.Search(m.ctx, dtos.DeviceServiceSearchQueryRequest{DeviceLibraryId: dl.Id})
	// 不存在则创建
	if len(dss) <= 0 {
		version := models.SupportVersion{}
		for _, v := range dl.SupportVersions {
			if v.Version == dl.Version {
				version = v
				break
			}
		}
		err = m.Add(m.ctx, models.DeviceService{
			//Id:                 dsCode,
			Name:               dl.Name,
			DeviceLibraryId:    dl.Id,
			ExpertMode:         version.ExpertMode,
			ExpertModeContent:  version.ExpertModeContent,
			DockerParamsSwitch: version.DockerParamsSwitch,
			DockerParams:       version.DockerParams,
			ContainerName:      dl.ContainerName,
			Config:             make(map[string]interface{}),
		})
		if err != nil {
			m.lc.Errorf("add device service err:%v", err)
			return err
		}
		return nil
	}

	ds := dss[0]
	// 存在则 判断是否更新
	if m.GetState(ds.Id) != constants.RunStatusStarted {
		return nil
	}

	// 重启
	if err = m.Stop(ds.Id); err != nil {
		m.lc.Errorf("stop deviceService(%s) err:%v", ds.Id, err)
		return err
	}
	if err = m.Start(ds.Id); err != nil {
		m.lc.Errorf("start deviceService(%s) err:%v", ds.Id, err)
		return err
	}

	return nil
}

func (m *driverServiceAppM) Search(ctx context.Context, req dtos.DeviceServiceSearchQueryRequest) ([]models.DeviceService, uint32, error) {
	offset, limit := req.BaseSearchConditionQuery.GetPage()

	deviceServices, total, err := m.dbClient.DeviceServicesSearch(offset, limit, req)
	if err != nil {
		return deviceServices, 0, err
	}

	dlIds := make([]string, 0)
	for i, _ := range deviceServices {
		dlIds = append(dlIds, deviceServices[i].DeviceLibraryId)
	}
	dls, _, err := m.getDriverApp().DeviceLibrariesSearch(m.ctx, dtos.DeviceLibrarySearchQueryRequest{
		BaseSearchConditionQuery: dtos.BaseSearchConditionQuery{Ids: dtos.ApiParamsArrayToString(dlIds)},
	})

	if err != nil {
		return deviceServices, 0, err
	}

	dlIdMap := make(map[string]models.DeviceLibrary)
	for i, _ := range dls {
		dlIdMap[dls[i].Id] = dls[i]
	}

	for i, v := range deviceServices {
		deviceServices[i].RunStatus = m.GetState(v.Id)
		if _, ok := dlIdMap[v.DeviceLibraryId]; ok {
			deviceServices[i].ImageExist = dlIdMap[v.DeviceLibraryId].OperateStatus == constants.OperateStatusInstalled
		}
	}

	return deviceServices, total, nil
}

func (m *driverServiceAppM) Del(ctx context.Context, id string) error {
	ds, edgeXErr := m.dbClient.DeviceServiceById(id)
	if edgeXErr != nil {
		return edgeXErr
	}

	// 检查容器是否在运行中
	if m.GetState(id) != constants.RunStatusStopped {
		return errort.NewCommonErr(errort.DeviceServiceMustStopService, fmt.Errorf("must stop service"))
	}
	m.dbClient.GetDBInstance().Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&models.DeviceService{}).Where("id =?", id).Delete(&models.DeviceService{}).Error
		if err != nil {
			return err
		}
		err = tx.Model(&models.Device{}).Where("drive_instance_id = ?", id).Updates(map[string]interface{}{"drive_instance_id": ""}).Error
		if err != nil {
			return err
		}
		return nil
	})

	// 删除容器、监控
	err := m.appModel.DeleteInstance(dtos.DeviceServiceFromModel(ds))
	if err != nil {
		m.lc.Errorf("DeleteInstance err:%v", err)
	}

	// 刷新agent 服务信息
	//go m.FlushStatsToAgent()
	// 删除后台监控
	delete(m.dsMonitor, id)
	m.state.Delete(id)
	return nil
}

//
func (m *driverServiceAppM) Get(ctx context.Context, id string) (models.DeviceService, error) {
	if id == "" {
		return models.DeviceService{}, errort.NewCommonErr(errort.DefaultReqParamsError, fmt.Errorf("id(%s) is empty", id))
	}
	deviceService, err := m.dbClient.DeviceServiceById(id)
	if err != nil {
		return deviceService, err
	}
	deviceService.RunStatus = m.GetState(id)

	dl, _ := m.getDriverApp().DriverLibById(deviceService.DeviceLibraryId)
	deviceService.ImageExist = dl.OperateStatus == constants.OperateStatusInstalled
	return deviceService, nil
}

//
func (m *driverServiceAppM) InProgress(id string) bool {
	state, ok := m.state.Load(id)
	if !ok {
		return false
	}
	if state.(int) == constants.RunStatusStarting || state.(int) == constants.RunStatusStopping {
		return true
	}
	return false
}

//
//// 监控驱动运行状态
func (m *driverServiceAppM) initMonitor() {
	dbClient := container.DBClientFrom(m.dic.Get)
	lc := bootstrapContainer.LoggingClientFrom(m.dic.Get)
	ds, _, err := dbClient.DeviceServicesSearch(0, -1, dtos.DeviceServiceSearchQueryRequest{})
	if err != nil {
		lc.Errorf("DeviceServicesSearch err %v", err)
		return
	}
	for _, v := range ds {
		m.dsMonitor[v.Id] = NewDeviceServiceMonitor(m.ctx, dtos.DeviceServiceFromModel(v), m.dic)
	}
}

func (m *driverServiceAppM) UpdateRunStatus(ctx context.Context, req dtos.UpdateDeviceServiceRunStatusRequest) error {
	// 1.正在处理中，返回错误
	if m.InProgress(req.Id) {
		return errort.NewCommonErr(errort.DeviceServiceMustStopDoingService, fmt.Errorf("device service is processing"))
	}

	// 2.请求状态和本地状态一致，无需操作
	if req.RunStatus == m.GetState(req.Id) {
		m.lc.Infof("driverService state is %d", req.RunStatus)
		return nil
	}

	_, err := m.Get(ctx, req.Id)
	if err != nil {
		return err
	}

	if req.RunStatus == constants.RunStatusStopped {
		if err = m.Stop(req.Id); err != nil {
			return err
		}
	} else if req.RunStatus == constants.RunStatusStarted {
		if err = m.Start(req.Id); err != nil {
			return err
		}
	}

	return nil
}

// 将deviceService里的配置转换到配置文件中然后启动服务
func (m *driverServiceAppM) buildServiceRunCfg(serviceIp string, runPort int, ds models.DeviceService) (string, error) {
	if ds.DriverType == constants.DriverLibTypeDefault {
		return m.buildDriverCfg(serviceIp, runPort, ds)
	} else if ds.DriverType == constants.DriverLibTypeAppService {
		//return m.buildAppCfg(serviceIp, runPort, ds)
	}
	return "", nil
}

func (m *driverServiceAppM) buildDriverCfg(localDefaultIp string, runPort int, ds models.DeviceService) (string, error) {
	configuration := &dtos.DriverConfig{}
	sysConfig := container.ConfigurationFrom(m.dic.Get)

	// 读取模版配置
	if _, err := toml.Decode(getDriverConfigTemplate(ds), configuration); err != nil {
		return "", err
	}

	// 修改与核心服务通信的ip
	for k, v := range configuration.Clients {
		if k == "Core" {
			data := v
			data.Address = strings.Replace(data.Address, "127.0.0.1", localDefaultIp, -1)
			data.Address = strings.Split(data.Address, ":")[0] + ":" + strings.Split(sysConfig.RpcServer.Address, ":")[1]
			configuration.Clients[k] = data
		}
	}

	configuration.Service.ID = ds.Id
	configuration.Service.Name = ds.Name
	// 驱动服务只开启rpc服务
	configuration.Service.Server.Address = "0.0.0.0:" + strconv.Itoa(runPort)

	if ds.ExpertMode && ds.ExpertModeContent != "" {
		configuration.CustomParam = string(ds.ExpertModeContent)
	}

	// set log level
	configuration.Logger.LogLevel = constants.LogMap[ds.LogLevel]
	configuration.Logger.FileName = "/mnt/logs/driver.log"

	var buff bytes.Buffer
	e := toml.NewEncoder(&buff)
	err := e.Encode(configuration)
	if err != nil {
		return "", err
	}

	return buff.String(), nil
}
