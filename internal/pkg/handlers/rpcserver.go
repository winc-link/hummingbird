package handlers

import (
	"context"
	"errors"

	"net"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"

	"github.com/winc-link/hummingbird/internal/pkg/container"
	"github.com/winc-link/hummingbird/internal/pkg/di"
	"github.com/winc-link/hummingbird/internal/pkg/limit"
	//"gitlab.com/tedge/edgex/internal/pkg/container"
	//"gitlab.com/tedge/edgex/internal/pkg/di"
	//"gitlab.com/tedge/edgex/internal/pkg/limit"
)

var basepath string

func init() {
	_, currentFile, _, _ := runtime.Caller(0)
	basepath = filepath.Dir(currentFile)
}

// Path returns the absolute path the given relative file or directory path,
// relative to the google.golang.org/grpc/examples/data directory in the
// user's GOPATH.  If rel is already absolute, it is returned unmodified.
func path(rel string) string {
	if filepath.IsAbs(rel) {
		return rel
	}

	return filepath.Join(basepath, rel)
}

var kaep = keepalive.EnforcementPolicy{
	MinTime:             5 * time.Second, // If a client pings more than once every 5 seconds, terminate the connection
	PermitWithoutStream: true,            // Allow pings even when there are no active streams
}

type RegisterFunc func(serve *grpc.Server)

type RPCServer struct {
	S         *grpc.Server
	isRunning bool
}

var (
	limitService   limit.LimitService
	methodLimitMap sync.Map
	serverRpcLog   = os.Getenv(SERVER_RPC_LOG_ENABLE)
	serverTimeout  time.Duration
)

func methodLimit(method string, lmc limit.LimitMethodConf) bool {
	_, exist := lmc.GetLimitMethods()[method]
	return exist
}

func requestLimit(ctx context.Context, method string) error {
	value, exist := methodLimitMap.Load(method)
	if !exist {
		value = limitService.Clone()
	}
	var srv = value.(limit.LimitService)
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, serverTimeout)
	defer cancel()
	if err := srv.ConsumeWithContext(ctx); err != nil {
		return err
	}
	methodLimitMap.Store(method, value)
	return nil
}

const (
	// 默认每个rpc接口100ms放入5个请求
	DefaultLimitTime = 100
	DefaultBurst     = 5
	// 默认服务端超时时间为5s
	DefaultServiceTimeout = 5 * time.Second
	SERVER_RPC_LOG_ENABLE = "SERVER_RPC_LOG_ENABLE"
)

func NewRPCServer(ctx context.Context, wg *sync.WaitGroup, dic *di.Container, register RegisterFunc) (*RPCServer, error) {
	var rs RPCServer
	lc := container.LoggingClientFrom(dic.Get)
	rpcConfig := container.ConfigurationFrom(dic.Get).GetBootstrap().RpcServer
	if rpcConfig.Address == "" {
		lc.Error("required rpc address")
		return nil, errors.New("required rpc address")
	}

	lis, err := net.Listen("tcp", rpcConfig.Address)
	if err != nil {
		lc.Errorf("failed to listen: %v", err)
		return nil, err
	}

	var (
	//limitTime int64 = DefaultLimitTime
	//burst           = DefaultBurst
	)
	//if rpcConfig.LimitSetting.LimitTime > 0 {
	//	limitTime = rpcConfig.LimitSetting.LimitTime
	//}
	//if rpcConfig.LimitSetting.Burst > 0 {
	//	burst = rpcConfig.LimitSetting.Burst
	//}
	if rpcConfig.Timeout > 0 {
		serverTimeout = time.Duration(rpcConfig.Timeout) * time.Millisecond
	} else {
		serverTimeout = DefaultServiceTimeout
	}

	limitService = limit.NewLimitService(limit.LimitOption{
		//LimitMillisecond: limitTime,
		//Burst:            burst,
	})
	lmc := container.LimitMethodConfFrom(dic.Get)
	if rpcConfig.UseTLS {
		creds, err := credentials.NewServerTLSFromFile(path(rpcConfig.CertFile), path(rpcConfig.KeyFile))
		if err != nil {
			lc.Errorf("failed to create credentials: %v", err)
			return nil, err
		}
		rs.S = grpc.NewServer(
			grpc.Creds(creds),
			grpc.KeepaliveEnforcementPolicy(kaep),
			withServerInterceptor(lc, lmc, dic),
		)

	} else {
		rs.S = grpc.NewServer(
			grpc.KeepaliveEnforcementPolicy(kaep),
			withServerInterceptor(lc, lmc, dic),
		)
	}
	wg.Add(1)
	go func() {
		defer wg.Done()

		<-ctx.Done()
		lc.Info("RPC server shutting down")
		rs.S.Stop()
		lc.Info("RPC server shut down")
	}()

	// registry Server
	register(rs.S)

	lc.Infof("RPC server starting ( %s )", rpcConfig.Address)

	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
			rs.isRunning = false
		}()
		rs.isRunning = true
		err = rs.S.Serve(lis)
		if err != nil {
			lc.Errorf("RPC server failed: %v", err)
			cancel := container.CancelFuncFrom(dic.Get)
			cancel()
		} else {
			lc.Info("RPC server stopped")
		}
	}()
	return &rs, nil
}
