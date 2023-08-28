package container

import (
	//"gitlab.com/tedge/edgex/internal/pkg/di"
	//"gitlab.com/tedge/edgex/internal/pkg/limit"
	"github.com/winc-link/hummingbird/internal/pkg/di"
	"github.com/winc-link/hummingbird/internal/pkg/limit"
)

// LimitMethodConfName contains the name of the interfaces.LimitMethodConf implementation in the DIC.
var LimitMethodConfName = di.TypeInstanceToName((*limit.LimitMethodConf)(nil))

// LimitMethodConfFrom helper function queries the DIC and returns the interfaces.LimitMethodConf implementation.
func LimitMethodConfFrom(get di.Get) limit.LimitMethodConf {
	conf, ok := get(LimitMethodConfName).(limit.LimitMethodConf)
	if !ok {
		return nil
	}

	return conf
}
