package bootstrap

import (
	"context"
	"github.com/winc-link/hummingbird/internal/pkg/config"
	"github.com/winc-link/hummingbird/internal/pkg/container"
	"github.com/winc-link/hummingbird/internal/pkg/di"
	"github.com/winc-link/hummingbird/internal/pkg/environment"
	"github.com/winc-link/hummingbird/internal/pkg/flags"
	"github.com/winc-link/hummingbird/internal/pkg/handlers"
	"github.com/winc-link/hummingbird/internal/pkg/logger"
	"github.com/winc-link/hummingbird/internal/pkg/processor"
	"github.com/winc-link/hummingbird/internal/pkg/startup"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// Deferred defines the signature of a function returned by RunAndReturnWaitGroup that should be executed via defer.
type Deferred func()

// translateInterruptToCancel spawns a go routine to translate the receipt of a SIGTERM signal to a call to cancel
// the context used by the bootstrap implementation.
func translateInterruptToCancel(ctx context.Context, wg *sync.WaitGroup, cancel context.CancelFunc) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		signalStream := make(chan os.Signal)
		defer func() {
			signal.Stop(signalStream)
			close(signalStream)
		}()
		signal.Notify(signalStream, os.Interrupt, syscall.SIGTERM)
		select {
		case <-signalStream:
			cancel()
			return
		case <-ctx.Done():
			return
		}
	}()
}

// RunAndReturnWaitGroup bootstraps an application.  It loads configuration and calls the provided list of handlers.
// Any long-running process should be spawned as a go routine in a handler.  Handlers are expected to return
// immediately.  Once all of the handlers are called this function will return a sync.WaitGroup reference to the caller.
// It is intended that the caller take whatever additional action makes sense before calling Wait() on the returned
// reference to wait for the application to be signaled to stop (and the corresponding goroutines spawned in the
// various handlers to be stopped cleanly).
func RunAndReturnWaitGroup(
	ctx context.Context,
	cancel context.CancelFunc,
	commonFlags flags.Common,
	serviceKey string,
	configStem string,
	serviceConfig config.Configuration,
	configUpdated processor.UpdatedStream,
	startupTimer startup.Timer,
	dic *di.Container,
	handlers []handlers.BootstrapHandler) (*sync.WaitGroup, Deferred, bool) {

	var err error
	var wg sync.WaitGroup
	deferred := func() {}

	translateInterruptToCancel(ctx, &wg, cancel)

	envVars := environment.NewVariables()
	configProcessor := processor.NewProcessor(commonFlags, envVars, startupTimer, ctx, &wg, configUpdated, dic)
	if err = configProcessor.Process(serviceConfig); err != nil {
		panic(err)
	}

	// Now the the configuration has been processed the logger has been created based on configuration.
	lc := logger.NewClient(serviceKey, serviceConfig.GetLogLevel(), serviceConfig.GetLogPath())
	configProcessor.Logger = lc

	dic.Update(di.ServiceConstructorMap{
		container.ConfigurationInterfaceName: func(get di.Get) interface{} {
			return serviceConfig
		},
		container.LoggingClientInterfaceName: func(get di.Get) interface{} {
			return lc
		},
		container.CancelFuncName: func(get di.Get) interface{} {
			return cancel
		},
	})

	// call individual bootstrap handlers.
	startedSuccessfully := true
	for i := range handlers {
		if handlers[i](ctx, &wg, startupTimer, dic) == false {
			cancel()
			startedSuccessfully = false
			break
		}
	}

	return &wg, deferred, startedSuccessfully
}

// Run bootstraps an application.  It loads configuration and calls the provided list of handlers.  Any long-running
// process should be spawned as a go routine in a handler.  Handlers are expected to return immediately.  Once all of
// the handlers are called this function will wait for any go routines spawned inside the handlers to exit before
// returning to the caller.  It is intended that the caller stop executing on the return of this function.
func Run(
	ctx context.Context,
	cancel context.CancelFunc,
	commonFlags flags.Common,
	serviceKey string,
	configStem string,
	serviceConfig config.Configuration,
	startupTimer startup.Timer,
	dic *di.Container,
	handlers []handlers.BootstrapHandler) {

	wg, deferred, _ := RunAndReturnWaitGroup(
		ctx,
		cancel,
		commonFlags,
		serviceKey,
		configStem,
		serviceConfig,
		nil,
		startupTimer,
		dic,
		handlers,
	)

	defer deferred()

	// wait for go routines to stop executing.
	wg.Wait()
}
