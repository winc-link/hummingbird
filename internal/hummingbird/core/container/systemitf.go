package container

import (
	//"gitlab.com/tedge/edgex/internal/pkg/di"
	//"gitlab.com/tedge/edgex/internal/tedge/resource/interfaces"
	interfaces "github.com/winc-link/hummingbird/internal/hummingbird/core/interface"
	"github.com/winc-link/hummingbird/internal/pkg/di"
)

// SystemItfName
var SystemItfName = di.TypeInstanceToName((*interfaces.SystemItf)(nil))

// SystemItfFrom
func SystemItfFrom(get di.Get) interfaces.SystemItf {
	return get(SystemItfName).(interfaces.SystemItf)
}
