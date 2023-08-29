package container

import (
	interfaces "github.com/winc-link/hummingbird/internal/hummingbird/core/interface"
	"github.com/winc-link/hummingbird/internal/pkg/di"
)

var DataDBClientInterfaceName = di.TypeInstanceToName((*interfaces.DataDBClient)(nil))

func DataDBClientFrom(get di.Get) interfaces.DataDBClient {
	return get(DataDBClientInterfaceName).(interfaces.DataDBClient)
}
