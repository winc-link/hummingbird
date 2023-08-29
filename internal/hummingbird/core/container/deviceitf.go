package container

import (
	interfaces "github.com/winc-link/hummingbird/internal/hummingbird/core/interface"
	"github.com/winc-link/hummingbird/internal/pkg/di"
)

// DeviceItfName
var DeviceItfName = di.TypeInstanceToName((*interfaces.DeviceItf)(nil))

// DeviceItfFrom
func DeviceItfFrom(get di.Get) interfaces.DeviceItf {
	return get(DeviceItfName).(interfaces.DeviceItf)
}
