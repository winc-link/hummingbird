package interfaces

import (
	"context"
	"github.com/winc-link/hummingbird/internal/dtos"
)

//先把这些临时放在这里
type SystemItf interface {
	GwConfigItf
	AdvConfigItf
	NetworkItf
	GatewayItf
}

type GwConfigItf interface {
	LoadGatewayConfig() error
	GetGatewayConfig() dtos.EdgeConfig
}

type AdvConfigItf interface {
	GetAdvanceConfig(ctx context.Context) (dtos.AdvanceConfig, error)
	UpdateAdvanceConfig(ctx context.Context, cfg dtos.AdvanceConfig) error
}

type NetworkItf interface {
	GetNetworks(ctx context.Context) (dtos.ConfigNetWorkResponse, dtos.ConfigDnsResponse)
	ConfigNetWork(ctx context.Context, isFlush bool) (resp dtos.ConfigNetWorkResponse, err error)
	ConfigNetWorkUpdate(ctx context.Context, req dtos.ConfigNetworkUpdateRequest) error
	ConfigDns(ctx context.Context) (dtos.ConfigDnsResponse, error)
	ConfigDnsUpdate(ctx context.Context, req dtos.ConfigDnsUpdateRequest) error
}

type GatewayItf interface {
	SystemBackupFileDownload(ctx context.Context) (string, error)
	SystemRecover(ctx context.Context, filepath string) error
}

type Starter interface {
	Conn() error
}
