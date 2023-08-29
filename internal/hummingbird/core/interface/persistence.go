package interfaces

import (
	"github.com/winc-link/hummingbird/internal/dtos"
)

type PersistItf interface {
	PersistDeviceItf
}

type PersistDeviceItf interface {
	SaveDeviceThingModelData(req dtos.ThingModelMessage) error
	SearchDeviceThingModelPropertyData(req dtos.ThingModelPropertyDataRequest) (interface{}, error)
	SearchDeviceThingModelHistoryPropertyData(req dtos.ThingModelPropertyDataRequest) (interface{}, int, error)
	SearchDeviceThingModelEventData(req dtos.ThingModelEventDataRequest) ([]dtos.ThingModelEventDataResponse, int, error)
	SearchDeviceThingModelServiceData(req dtos.ThingModelServiceDataRequest) ([]dtos.ThingModelServiceDataResponse, int, error)
	SearchDeviceMsgCount(startTime, endTime int64) (int, error)
}
