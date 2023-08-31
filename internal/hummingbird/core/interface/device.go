package interfaces

import (
	"context"
	"github.com/winc-link/edge-driver-proto/driverdevice"
	"github.com/winc-link/hummingbird/internal/dtos"
	"github.com/winc-link/hummingbird/internal/models"
)

type DeviceItf interface {
	DeviceCtlItf
	DeviceSyncItf
	OpenApiDeviceItf
}

type DeviceCtlItf interface {
	AddDevice(ctx context.Context, req dtos.DeviceAddRequest) (string, error)

	DevicesSearch(ctx context.Context, req dtos.DeviceSearchQueryRequest) ([]dtos.DeviceSearchQueryResponse, uint32, error)

	DevicesModelSearch(ctx context.Context, req dtos.DeviceSearchQueryRequest) ([]models.Device, uint32, error)

	DeviceById(ctx context.Context, id string) (dtos.DeviceInfoResponse, error)

	DeviceModelById(ctx context.Context, id string) (models.Device, error)

	DeviceByCloudId(ctx context.Context, id string) (models.Device, error)

	DeviceUpdate(ctx context.Context, req dtos.DeviceUpdateRequest) error

	DevicesBindDriver(ctx context.Context, req dtos.DevicesBindDriver) error

	DevicesUnBindDriver(ctx context.Context, req dtos.DevicesUnBindDriver) error

	DevicesBindProductId(ctx context.Context, req dtos.DevicesBindProductId) error

	ConnectIotPlatform(ctx context.Context, request *driverdevice.ConnectIotPlatformRequest) *driverdevice.ConnectIotPlatformResponse

	DisConnectIotPlatform(ctx context.Context, request *driverdevice.DisconnectIotPlatformRequest) *driverdevice.DisconnectIotPlatformResponse

	GetDeviceConnectStatus(ctx context.Context, request *driverdevice.GetDeviceConnectStatusRequest) *driverdevice.GetDeviceConnectStatusResponse

	DeviceMqttAuthInfo(ctx context.Context, id string) (dtos.DeviceAuthInfoResponse, error)

	AddMqttAuth(ctx context.Context, req dtos.AddMqttAuthInfoRequest) (string, error)

	DeleteDeviceById(ctx context.Context, id string) error

	BatchDeleteDevice(ctx context.Context, ids []string) error

	DeviceImportTemplateDownload(ctx context.Context, req dtos.DeviceImportTemplateRequest) (*dtos.ExportFile, error)

	DevicesImport(ctx context.Context, file *dtos.ImportFile, productId, driverInstanceId string) (int64, error)

	UploadValidated(ctx context.Context, file *dtos.ImportFile) error

	DevicesReportMsgGather(ctx context.Context) error

	DeviceAction(jobAction dtos.JobAction) dtos.DeviceExecRes

	DeviceInvokeThingService(invokeDeviceServiceReq dtos.InvokeDeviceServiceReq) (map[string]interface{}, error)

	SetDeviceProperty(req dtos.OpenApiSetDeviceThingModel) error

	DeviceEffectivePropertyData(deviceEffectivePropertyDataReq dtos.DeviceEffectivePropertyDataReq) (dtos.DeviceEffectivePropertyDataResponse, error)
}

type OpenApiDeviceItf interface {
	OpenApiDeviceById(ctx context.Context, id string) (dtos.OpenApiDeviceInfoResponse, error)
	OpenApiDeviceStatusById(ctx context.Context, id string) (dtos.OpenApiDeviceStatus, error)
	OpenApiDevicesSearch(ctx context.Context, req dtos.DeviceSearchQueryRequest) ([]dtos.OpenApiDeviceInfoResponse, uint32, error)
}
type DeviceSyncItf interface {
}
