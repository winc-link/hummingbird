package dtos

import (
	"github.com/winc-link/hummingbird/internal/models"
	"github.com/winc-link/hummingbird/internal/pkg/constants"

	//"gitlab.com/tedge/edgex/internal/models"
	//devicelibraryProto "gitlab.com/tedge/edgex/proto/devicelibrary"
	//deviceserviceProto "gitlab.com/tedge/edgex/proto/deviceservice"
	"gopkg.in/yaml.v3"
)

type DeviceService struct {
	Id                 string                 `json:"id,omitempty"`
	Name               string                 `json:"name"`
	Created            int64                  `json:"created,omitempty"`
	Modified           int64                  `json:"modified,omitempty"`
	BaseAddress        string                 `json:"baseAddress"`
	DeviceLibraryId    string                 `json:"deviceLibraryId"`
	Config             map[string]interface{} `json:"config"`
	DockerContainerId  string                 `json:"dockerContainerId"`
	ExpertMode         bool                   `json:"isExpertMode"`
	ExpertModeContent  string                 `json:"expertModeContent"`
	DockerParamsSwitch bool                   `json:"dockerParamsSwitch"`
	DockerParams       string                 `json:"dockerParams"`
	ContainerName      string                 `json:"container_name"`
}

// 启动实例时对应的配置
type RunServiceCfg struct {
	ImageRepo          string
	RunConfig          string
	DockerMountDevices []string
	DockerParams       string
	DriverName         string // 驱动名
}

func DeviceServiceFromModel(ds models.DeviceService) DeviceService {
	var dto DeviceService
	dto.Id = ds.Id
	dto.Name = ds.Name
	dto.BaseAddress = ds.BaseAddress
	dto.DeviceLibraryId = ds.DeviceLibraryId
	dto.Config = ds.Config
	dto.DockerContainerId = ds.DockerContainerId
	dto.ExpertMode = ds.ExpertMode
	dto.ExpertModeContent = ds.ExpertModeContent
	dto.DockerParamsSwitch = ds.DockerParamsSwitch
	dto.DockerParams = ds.DockerParams
	dto.ContainerName = ds.ContainerName
	return dto
}

//func FromDeviceServiceModelToRPC(mds models.DeviceService) *deviceserviceProto.DeviceService {
//	byteConfig, _ := json.Marshal(mds.Config)
//	var ds deviceserviceProto.DeviceService
//	ds.Id = mds.Id
//	ds.Name = mds.Name
//	ds.BaseAddress = mds.BaseAddress
//	ds.DeviceLibraryId = mds.DeviceLibraryId
//	ds.DockerContainerId = mds.DockerContainerId
//	ds.Config = byteConfig
//	ds.ExpertMode = mds.ExpertMode
//	ds.ExpertModeContent = mds.ExpertModeContent
//	ds.DockerParamsSwitch = mds.DockerParamsSwitch
//	ds.DockerParams = mds.DockerParams
//	ds.LogLevel = int64(mds.LogLevel)
//	ds.RunStatus = int32(mds.RunStatus)
//	ds.ImageExist = mds.ImageExist
//	return &ds
//}

//func FromDeviceServiceRpcToModel(ds *deviceserviceProto.DeviceService) models.DeviceService {
//	var config map[string]interface{}
//	if ds.Config != nil {
//		_ = json.Unmarshal(ds.Config, &config)
//	}
//
//	var mds models.DeviceService
//	mds.Id = ds.Id
//	mds.Name = ds.Name
//	mds.BaseAddress = ds.BaseAddress
//	mds.DeviceLibraryId = ds.DeviceLibraryId
//	mds.DockerContainerId = ds.DockerContainerId
//	mds.RunStatus = int(ds.RunStatus)
//	mds.Config = config
//	mds.ExpertMode = ds.ExpertMode
//	mds.ExpertModeContent = ds.ExpertModeContent
//	mds.DockerParamsSwitch = ds.DockerParamsSwitch
//	mds.DockerParams = ds.DockerParams
//	mds.ImageExist = ds.ImageExist
//	mds.DriverType = int(ds.DriverType)
//	return mds
//}

type DeviceServiceAddRequest struct {
	Id                 string                 `json:"id,omitempty" binding:"omitempty,t-special-char"`
	Name               string                 `json:"name"`
	DeviceLibraryId    string                 `json:"deviceLibraryId" binding:"required"`
	Config             map[string]interface{} `json:"config" binding:"required"`
	ExpertMode         bool                   `json:"expertMode"`
	ExpertModeContent  string                 `json:"expertModeContent"`
	DockerParamsSwitch bool                   `json:"dockerParamsSwitch"`
	DockerParams       string                 `json:"dockerParams"`
	DriverType         int                    `json:"driverType" binding:"omitempty,oneof=1 2"` //驱动库类型，1：驱动，2：三方应用
}

//func FromDeviceServiceAddToRpc(req DeviceServiceAddRequest) *deviceserviceProto.DeviceService {
//	byteConfig, _ := json.Marshal(req.Config)
//	return &deviceserviceProto.DeviceService{
//		Id:                 req.Id,
//		Name:               req.Name,
//		DeviceLibraryId:    req.DeviceLibraryId,
//		Config:             byteConfig,
//		ExpertMode:         req.ExpertMode,
//		ExpertModeContent:  req.ExpertModeContent,
//		DockerParamsSwitch: req.DockerParamsSwitch,
//		DockerParams:       req.DockerParams,
//		DriverType:         int32(req.DriverType),
//	}
//}

func DeviceServiceFromDeviceServiceAddRequest(ds DeviceServiceAddRequest) models.DeviceService {
	var mds models.DeviceService
	mds.Id = ds.Id
	mds.Name = ds.Name
	mds.Config = ds.Config
	mds.DeviceLibraryId = ds.DeviceLibraryId
	mds.ExpertMode = ds.ExpertMode
	mds.ExpertModeContent = ds.ExpertModeContent
	mds.DockerParamsSwitch = ds.DockerParamsSwitch
	mds.DockerParams = ds.DockerParams
	mds.DriverType = ds.DriverType
	return mds
}

type DeviceServiceUpdateRequest struct {
	Id                 string                  `json:"id" binding:"required"`
	DeviceLibraryId    *string                 `json:"deviceLibraryId"`
	Name               *string                 `json:"name"`
	Config             *map[string]interface{} `json:"config"`
	ExpertMode         *bool                   `json:"expertMode"`
	ExpertModeContent  *string                 `json:"expertModeContent"`
	DockerParamsSwitch *bool                   `json:"docker_params_switch"`
	DockerParams       *string                 `json:"docker_params"`
	Platform           constants.IotPlatform   `json:"platform"`
	//IsIgnoreRunStatus  bool
}

//func FromRpcToUpdateDeviceService(rpc *deviceserviceProto.UpdateDeviceService) DeviceServiceUpdateRequest {
//	var config map[string]interface{}
//	if rpc.Config != nil {
//		_ = json.Unmarshal(rpc.Config, &config)
//	}
//	return DeviceServiceUpdateRequest{
//		Id:                 rpc.Id,
//		Name:               rpc.Name,
//		DeviceLibraryId:    rpc.DeviceLibraryId,
//		Config:             &config,
//		ExpertMode:         rpc.ExpertMode,
//		ExpertModeContent:  rpc.ExpertModeContent,
//		DockerParamsSwitch: rpc.DockerParamsSwitch,
//		DockerParams:       rpc.DockerParams,
//	}
//}

type UpdateDeviceServiceRunStatusRequest struct {
	Id        string `json:"id"`
	RunStatus int    `json:"run_status"  binding:"required,oneof=1 2"`
}

type DeviceServiceRunLogRequest struct {
	Id      string `json:"id"`
	Operate int    `json:"operate" binding:"required,oneof=1 2"`
}

type DeviceServiceDeleteRequest struct {
	Id string `json:"id" binding:"required"`
}

//func FromUpdateDeviceServiceRunStatusToRpc(req UpdateDeviceServiceRunStatusRequest) *deviceserviceProto.UpdateDeviceServiceRunStatusRequest {
//	return &deviceserviceProto.UpdateDeviceServiceRunStatusRequest{
//		Id:        req.Id,
//		RunStatus: int32(req.RunStatus),
//	}
//}
//
//func FromDeviceServiceSearchQueryRequestToRpc(req DeviceServiceSearchQueryRequest) *deviceserviceProto.DeviceServiceSearchRequest {
//	return &deviceserviceProto.DeviceServiceSearchRequest{
//		BaseSearchConditionQuery: FromBaseSearchConditionQueryToRpc(req.BaseSearchConditionQuery),
//		DeviceLibraryId:          req.DeviceLibraryId,
//		DriverType:               int32(req.DriverType),
//	}
//}

//func FromRpcToUpdateDeviceServiceRunStatus(rpc *deviceserviceProto.UpdateDeviceServiceRunStatusRequest) UpdateDeviceServiceRunStatusRequest {
//	return UpdateDeviceServiceRunStatusRequest{
//		Id:        rpc.Id,
//		RunStatus: int(rpc.RunStatus),
//	}
//}

func ReplaceDeviceServiceModelFieldsWithDTO(ds *models.DeviceService, patch DeviceServiceUpdateRequest) {
	if patch.Config != nil {
		ds.Config = *patch.Config
	}
	if patch.DeviceLibraryId != nil {
		ds.DeviceLibraryId = *patch.DeviceLibraryId
	}
	if patch.ExpertMode != nil {
		ds.ExpertMode = *patch.ExpertMode
	}
	if patch.ExpertModeContent != nil {
		ds.ExpertModeContent = *patch.ExpertModeContent
	}
	if patch.DockerParamsSwitch != nil {
		ds.DockerParamsSwitch = *patch.DockerParamsSwitch
	}
	if patch.DockerParams != nil {
		ds.DockerParams = *patch.DockerParams
	}
	if patch.Platform != "" {
		ds.Platform = patch.Platform
	}
}

type DeviceServiceSearchQueryRequest struct {
	BaseSearchConditionQuery `schema:",inline"`
	ProductId                string `form:"productId"`
	CloudProductId           string `form:"cloudProductId"`
	DeviceLibraryId          string `form:"deviceLibraryId"`  // 驱动库ID
	DeviceLibraryIds         string `form:"deviceLibraryIds"` // 驱动库IDs
	Platform                 string `form:"platform"`
	DriverType               int    `form:"driver_type" binding:"omitempty,oneof=1 2"` //驱动库类型，1：驱动，2：三方应用
}

//func ToDeviceServiceSearchQueryRequestDTO(req *deviceserviceProto.DeviceServiceSearchRequest) DeviceServiceSearchQueryRequest {
//
//	if req.BaseSearchConditionQuery == nil {
//		return DeviceServiceSearchQueryRequest{
//			DeviceLibraryId: req.DeviceLibraryId,
//			DriverType:      int(req.DriverType),
//		}
//	} else {
//		return DeviceServiceSearchQueryRequest{
//			BaseSearchConditionQuery: ToBaseSearchConditionQueryDTO(req.BaseSearchConditionQuery),
//			DeviceLibraryId:          req.DeviceLibraryId,
//			DriverType:               int(req.DriverType),
//		}
//	}
//}

//func FromDeviceServiceUpdateToRpc(req DeviceServiceUpdateRequest) *deviceserviceProto.UpdateDeviceService {
//	var byteConfig []byte
//	if req.Config != nil {
//		byteConfig, _ = json.Marshal(&req.Config)
//	} else {
//		byteConfig = nil
//	}
//
//	return &deviceserviceProto.UpdateDeviceService{
//		Id:                 req.Id,
//		Name:               req.Name,
//		DeviceLibraryId:    req.DeviceLibraryId,
//		Config:             byteConfig,
//		ExpertMode:         req.ExpertMode,
//		ExpertModeContent:  req.ExpertModeContent,
//		DockerParamsSwitch: req.DockerParamsSwitch,
//		DockerParams:       req.DockerParams,
//	}
//}

/************** Response **************/

type DeviceServiceResponse struct {
	Id            string                `json:"id"`
	Name          string                `json:"name"`
	DeviceLibrary DeviceLibraryResponse `json:"deviceLibrary"`
	//Version       string                `json:"version"`
	RunStatus          int         `json:"runStatus"`
	Config             interface{} `json:"config"`
	ExpertMode         bool        `json:"expertMode"`
	ExpertModeContent  string      `json:"expertModeContent"`
	DockerParamsSwitch bool        `json:"dockerParamsSwitch"`
	DockerParams       string      `json:"dockerParams"`
	CreateAt           int64       `json:"create_at"`
	ImageExist         bool        `json:"imageExist"`
	Platform           string      `json:"platform"`
}

func DeviceServiceResponseFromModel(ds models.DeviceService, dl models.DeviceLibrary) DeviceServiceResponse {
	return DeviceServiceResponse{
		Id:   ds.Id,
		Name: ds.Name,
		//Version: DeviceLibraryResponseFromModel(dl).Version,
		DeviceLibrary:      DeviceLibraryResponseFromModel(dl),
		RunStatus:          ds.RunStatus,
		Config:             ds.Config,
		ExpertMode:         ds.ExpertMode,
		ExpertModeContent:  ds.ExpertModeContent,
		DockerParamsSwitch: ds.DockerParamsSwitch,
		DockerParams:       ds.DockerParams,
		ImageExist:         ds.ImageExist,
		CreateAt:           ds.Created,
		Platform:           string(ds.Platform),
	}
}

//func FromDeviceServiceRpcToResponse(ds *deviceserviceProto.DeviceService, dl *devicelibraryProto.DeviceLibrary) DeviceServiceResponse {
//	var cfg map[string]interface{}
//	_ = json.Unmarshal(ds.Config, &cfg)
//	return DeviceServiceResponse{
//		Id:                 ds.Id,
//		Name:               ds.Name,
//		RunStatus:          int(ds.RunStatus),
//		DeviceLibrary:      FromDeviceLibraryRpcToResponse(dl),
//		Config:             cfg,
//		ExpertMode:         ds.ExpertMode,
//		ExpertModeContent:  ds.ExpertModeContent,
//		DockerParamsSwitch: ds.DockerParamsSwitch,
//		DockerParams:       ds.DockerParams,
//		ImageExist:         ds.ImageExist,
//	}
//}

func FromYamlStrToMap(yamlStr string) (m map[string]interface{}, err error) {
	err = yaml.Unmarshal([]byte(yamlStr), &m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

type UpdateDeviceServiceRunStatusResponse struct {
	Id        string `json:"id"`
	RunStatus int    `json:"run_status"`
}

type UpdateServiceLogLevelConfigRequest struct {
	Id       string `json:"id"` // 驱动或应用ID
	LogLevel int64  `json:"logLevel"`
}
