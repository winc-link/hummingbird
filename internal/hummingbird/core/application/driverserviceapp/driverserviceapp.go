package driverserviceapp

import (
	"context"
	interfaces "github.com/winc-link/hummingbird/internal/hummingbird/core/interface"
	"github.com/winc-link/hummingbird/internal/pkg/container"
	"github.com/winc-link/hummingbird/internal/pkg/di"
	"github.com/winc-link/hummingbird/internal/pkg/logger"
)

type driverServiceApp struct {
	dic *di.Container
	lc  logger.LoggingClient

	*driverServiceAppM
}

func NewDriverServiceApp(ctx context.Context, dic *di.Container) interfaces.DriverServiceApp {
	return &driverServiceApp{
		dic: dic,
		lc:  container.LoggingClientFrom(dic.Get),

		driverServiceAppM: newDriverServiceApp(ctx, dic),
	}
}
