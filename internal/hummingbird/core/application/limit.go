package application

import (
	//"gitlab.com/tedge/edgex/internal/pkg/limit"
	//"gitlab.com/tedge/edgex/internal/tedge/resource/config"
	"github.com/winc-link/hummingbird/internal/hummingbird/core/config"
	"github.com/winc-link/hummingbird/internal/pkg/limit"
)

type LimitMethodConf struct {
	methods map[string]struct{}
}

func NewLimitMethodConf(configuration config.ConfigurationStruct) limit.LimitMethodConf {
	var conf = &LimitMethodConf{methods: make(map[string]struct{})}
	for _, method := range configuration.Writable.LimitMethods {
		conf.methods[method] = struct{}{}
	}
	return conf
}

func (lmc *LimitMethodConf) GetLimitMethods() map[string]struct{} {
	return lmc.methods
}
