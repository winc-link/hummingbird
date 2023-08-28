package container

import (
	"github.com/winc-link/hummingbird/internal/pkg/cache"
	"github.com/winc-link/hummingbird/internal/pkg/di"
)

var CacheFuncName = di.TypeInstanceToName((*cache.Cache)(nil))

func CacheFuncFrom(get di.Get) cache.Cache {
	client, ok := get(CacheFuncName).(cache.Cache)
	if !ok {
		return nil
	}

	return client
}
