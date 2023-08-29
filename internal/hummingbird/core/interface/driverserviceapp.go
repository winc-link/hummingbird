package interfaces

import (
	//"context"

	"context"
	//"gitlab.com/tedge/edgex/internal/dtos"
	//"gitlab.com/tedge/edgex/internal/models"
	"github.com/winc-link/hummingbird/internal/dtos"
	"github.com/winc-link/hummingbird/internal/models"
)

type DriverServiceApp interface {
	GetState(id string) int
	SetState(id string, state int)
	Start(id string) error // 升级
	Stop(id string) error
	ReStart(id string) error
	Add(ctx context.Context, ds models.DeviceService) error
	Update(ctx context.Context, dto dtos.DeviceServiceUpdateRequest) error
	Del(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (models.DeviceService, error)
	Search(ctx context.Context, req dtos.DeviceServiceSearchQueryRequest) ([]models.DeviceService, uint32, error)
	UpdateRunStatus(ctx context.Context, req dtos.UpdateDeviceServiceRunStatusRequest) error
	InProgress(id string) bool
	Upgrade(dl models.DeviceLibrary) error // 升级驱动实例
}
