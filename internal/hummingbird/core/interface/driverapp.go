package interfaces

import (
	"context"
	"github.com/winc-link/hummingbird/internal/dtos"
	"github.com/winc-link/hummingbird/internal/models"
)

type DriverLibApp interface {
	AddDriverLib(ctx context.Context, dl dtos.DeviceLibraryAddRequest) error
	DeleteDeviceLibraryById(ctx context.Context, id string) error
	DeviceLibraryById(ctx context.Context, id string) (models.DeviceLibrary, error)
	DeviceLibrariesSearch(ctx context.Context, req dtos.DeviceLibrarySearchQueryRequest) ([]models.DeviceLibrary, uint32, error)
	UpdateDeviceLibrary(ctx context.Context, update dtos.UpdateDeviceLibrary) error
	UpgradeDeviceLibrary(ctx context.Context, req dtos.DeviceLibraryUpgradeRequest) error
	DriverLibById(dlId string) (models.DeviceLibrary, error)
	GetDriverClassify(ctx context.Context, req dtos.DriverClassifyQueryRequest) ([]dtos.DriverClassifyResponse, uint32, error)
	GetDeviceLibraryAndMirrorConfig(dlId string) (dl models.DeviceLibrary, dc models.DockerConfig, err error)
	DriverDownConfigItf
}
type DriverDownConfigItf interface {
	DownConfigAdd(ctx context.Context, req dtos.DockerConfigAddRequest) error
	DownConfigUpdate(ctx context.Context, req dtos.DockerConfigUpdateRequest) error
	DownConfigSearch(ctx context.Context, req dtos.DockerConfigSearchQueryRequest) ([]models.DockerConfig, uint32, error)
	DownConfigDel(ctx context.Context, id string) error
}
