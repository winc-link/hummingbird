package processor

import (
	"context"
	"fmt"
	"github.com/pelletier/go-toml/v2"
	"github.com/winc-link/hummingbird/internal/pkg/config"
	"github.com/winc-link/hummingbird/internal/pkg/di"
	"github.com/winc-link/hummingbird/internal/pkg/environment"
	"github.com/winc-link/hummingbird/internal/pkg/flags"
	"github.com/winc-link/hummingbird/internal/pkg/logger"
	"github.com/winc-link/hummingbird/internal/pkg/startup"
	"io/ioutil"
	"sync"
)

// UpdatedStream defines the stream type that is notified by ListenForChanges when a configuration update is received.
type UpdatedStream chan struct{}

type Processor struct {
	Logger          logger.LoggingClient
	flags           flags.Common
	envVars         *environment.Variables
	startupTimer    startup.Timer
	ctx             context.Context
	wg              *sync.WaitGroup
	configUpdated   UpdatedStream
	dic             *di.Container
	overwriteConfig bool
}

// NewProcessor creates a new configuration Processor
func NewProcessor(
	flags flags.Common,
	envVars *environment.Variables,
	startupTimer startup.Timer,
	ctx context.Context,
	wg *sync.WaitGroup,
	configUpdated UpdatedStream,
	dic *di.Container,
) *Processor {
	return &Processor{
		flags:         flags,
		envVars:       envVars,
		startupTimer:  startupTimer,
		ctx:           ctx,
		wg:            wg,
		configUpdated: configUpdated,
		dic:           dic,
	}
}

func (cp *Processor) Process(serviceConfig config.Configuration) error {
	// Create some shorthand for frequently used items
	//envVars := cp.envVars

	cp.overwriteConfig = cp.flags.OverwriteConfig()

	// Local configuration must be loaded first in case need registry processor info and/or
	// need to push it to the Configuration Provider.
	if err := cp.loadFromFile(serviceConfig); err != nil {
		return err
	}

	return nil
}

// LoadFromFile attempts to read and unmarshal toml-based configuration into a configuration struct.
func (cp *Processor) loadFromFile(config config.Configuration) error {
	configDir := environment.GetConfDir(cp.Logger, cp.flags.ConfigDirectory())
	contents, err := ioutil.ReadFile(configDir)
	if err != nil {
		return fmt.Errorf("could not load configuration file (%s): %s", configDir, err.Error())
	}
	if err = toml.Unmarshal(contents, config); err != nil {
		return fmt.Errorf("could not load configuration file (%s): %s", configDir, err.Error())
	}

	fmt.Println(fmt.Sprintf("Loaded configuration from %s", configDir))

	return nil
}

// logConfigInfo logs the processor info message with number over overrides that occurred.
func (cp *Processor) logConfigInfo(message string, overrideCount int) {
	if cp.Logger == nil {
		fmt.Println(fmt.Sprintf("%s (%d envVars overrides applied)", message, overrideCount))
		return
	}
	cp.Logger.Info(fmt.Sprintf("%s (%d envVars overrides applied)", message, overrideCount))
}
