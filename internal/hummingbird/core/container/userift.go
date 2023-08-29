package container

import (
	interfaces "github.com/winc-link/hummingbird/internal/hummingbird/core/interface"
	"github.com/winc-link/hummingbird/internal/pkg/di"
)

// UserItfName
var UserItfName = di.TypeInstanceToName((*interfaces.UserItf)(nil))

// UserItfFrom
func UserItfFrom(get di.Get) interfaces.UserItf {
	return get(UserItfName).(interfaces.UserItf)
}
