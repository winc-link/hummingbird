package config

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"gopkg.in/natefinch/lumberjack.v2"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v2"

	pkgconfig "github.com/winc-link/hummingbird/internal/pkg/config"
)

const EdgeMqttBroker = "mqtt-broker"

var (
	defaultPluginConfig = make(map[string]Configuration)
	configFileFullPath  string
	config              Config
)

// Configuration is the interface that enable the implementation to parse config from the global config file.
// Plugin admin and prometheus are two examples.
type Configuration interface {
	// Validate validates the configuration.
	// If returns error, the broker will not start.
	Validate() error
	// Unmarshaler defined how to unmarshal YAML into the config structure.
	yaml.Unmarshaler
}

// RegisterDefaultPluginConfig registers the default configuration for the given plugin.
func RegisterDefaultPluginConfig(name string, config Configuration) {
	if _, ok := defaultPluginConfig[name]; ok {
		panic(fmt.Sprintf("duplicated default config for %s plugin", name))
	}
	defaultPluginConfig[name] = config

}

// DefaultConfig return the default configuration.
// If config file is not provided, mqttd will start with DefaultConfig.
func DefaultConfig() Config {
	c := Config{
		Listeners: DefaultListeners,
		MQTT:      DefaultMQTTConfig,
		API:       DefaultAPI,
		Log: LogConfig{
			Level:    "info",
			FilePath: "/var/tedge/logs/mqtt-broker.log",
		},
		Plugins:           make(pluginConfig),
		PluginOrder:       []string{"aplugin"},
		Persistence:       DefaultPersistenceConfig,
		TopicAliasManager: DefaultTopicAliasManager,
	}

	for name, v := range defaultPluginConfig {
		c.Plugins[name] = v
	}
	return c
}

var DefaultListeners = []*ListenerConfig{
	{
		Address:    "0.0.0.0:58090",
		TLSOptions: nil,
		Websocket:  nil,
	},
	{
		Address: "0.0.0.0:58091",
		Websocket: &WebsocketOptions{
			Path: "/",
		},
	}, {
		Address: "0.0.0.0:21883",
		TLSOptions: &TLSOptions{
			CACert: "/etc/tedge-mqtt-broker/ca.crt",
			Cert:   "/etc/tedge-mqtt-broker/server.pem",
			Key:    "/etc/tedge-mqtt-broker/server.key",
		},
	},
}

// LogConfig is use to configure the log behaviors.
type LogConfig struct {
	// Level is the log level. Possible values: debug, info, warn, error
	Level    string `yaml:"level"`
	FilePath string `yaml:"file_path"`
	// DumpPacket indicates whether to dump MQTT packet in debug level.
	DumpPacket bool `yaml:"dump_packet"`
}

func (l LogConfig) Validate() error {
	level := strings.ToLower(l.Level)
	if level != "debug" && level != "info" && level != "warn" && level != "error" {
		return fmt.Errorf("invalid log level: %s", l.Level)
	}
	return nil
}

// pluginConfig stores the plugin default configuration, key by the plugin name.
// If the plugin has default configuration, it should call RegisterDefaultPluginConfig in it's init function to register.
type pluginConfig map[string]Configuration

func (p pluginConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	for _, v := range p {
		err := unmarshal(v)
		if err != nil {
			return err
		}
	}
	return nil
}

// Config is the configration for mqttd.
type Config struct {
	Listeners []*ListenerConfig `yaml:"listeners"`
	API       API               `yaml:"api"`
	MQTT      MQTT              `yaml:"mqttclient,omitempty"`
	GRPC      GRPC              `yaml:"gRPC"`
	Log       LogConfig         `yaml:"log"`
	PidFile   string            `yaml:"pid_file"`
	ConfigDir string            `yaml:"config_dir"`
	Plugins   pluginConfig      `yaml:"plugins"`
	// PluginOrder is a slice that contains the name of the plugin which will be loaded.
	// Giving a correct order to the slice is significant,
	// because it represents the loading order which affect the behavior of the broker.
	PluginOrder       []string           `yaml:"plugin_order"`
	Persistence       Persistence        `yaml:"persistence"`
	TopicAliasManager TopicAliasManager  `yaml:"topic_alias_manager"`
	Database          pkgconfig.Database `yaml:"data_base"`
}

type GRPC struct {
	Endpoint string `yaml:"endpoint"`
}

type TLSOptions struct {
	// CACert is the trust CA certificate file.
	CACert string `yaml:"cacert"`
	// Cert is the path to certificate file.
	Cert string `yaml:"cert"`
	// Key is the path to key file.
	Key string `yaml:"key"`
	// Verify indicates whether to verify client cert.
	Verify bool `yaml:"verify"`
}

type ListenerConfig struct {
	Address     string `yaml:"address"`
	*TLSOptions `yaml:"tls"`
	Websocket   *WebsocketOptions `yaml:"websocket"`
}

type WebsocketOptions struct {
	Path string `yaml:"path"`
}

func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type config Config
	raw := config(DefaultConfig())
	if err := unmarshal(&raw); err != nil {
		return err
	}
	emptyMQTT := MQTT{}
	if raw.MQTT == emptyMQTT {
		raw.MQTT = DefaultMQTTConfig
	}
	if len(raw.Plugins) == 0 {
		raw.Plugins = make(pluginConfig)
		for name, v := range defaultPluginConfig {
			raw.Plugins[name] = v
		}
	} else {
		for name, v := range raw.Plugins {
			if v == nil {
				raw.Plugins[name] = defaultPluginConfig[name]
			}
		}
	}
	*c = Config(raw)
	return nil
}

func (c Config) Validate() (err error) {
	err = c.Log.Validate()
	if err != nil {
		return err
	}
	err = c.API.Validate()
	if err != nil {
		return err
	}
	err = c.MQTT.Validate()
	if err != nil {
		return err
	}
	err = c.Persistence.Validate()
	if err != nil {
		return err
	}
	for _, conf := range c.Plugins {
		err := conf.Validate()
		if err != nil {
			return err
		}
	}
	return nil
}

func ParseConfig(filePath string) (Config, error) {
	if filePath == "" {
		return DefaultConfig(), nil
	}
	if _, err := os.Stat(filePath); err != nil {
		fmt.Println("unspecificed configuration file, use default config")
		return DefaultConfig(), nil
	}
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return config, err
	}
	config = DefaultConfig()
	err = yaml.Unmarshal(b, &config)
	if err != nil {
		return config, err
	}
	config.ConfigDir = path.Dir(filePath)
	err = config.Validate()
	if err != nil {
		return Config{}, err
	}
	configFileFullPath = filePath
	return config, err
}

func UpdateLogLevel(level string) {
	config.Log.Level = level
}

func GetLogLevel() string {
	return config.Log.Level
}

func WriteToFile() error {
	return config.writeToFile()
}

func (c Config) writeToFile() error {
	var (
		err  error
		buff bytes.Buffer
	)
	e := yaml.NewEncoder(&buff)
	if err = e.Encode(c); err != nil {
		return err
	}
	if err = ioutil.WriteFile(configFileFullPath+".tmp", buff.Bytes(), 0644); err != nil {
		return err
	}
	os.Remove(configFileFullPath)
	return os.Rename(configFileFullPath+".tmp", configFileFullPath)
}

func (c Config) GetLogger(config LogConfig) (*zap.AtomicLevel, *zap.Logger, error) {
	var logLevel zapcore.Level
	err := logLevel.UnmarshalText([]byte(config.Level))
	if err != nil {
		return nil, nil, err
	}
	var level = zap.NewAtomicLevelAt(logLevel)
	if config.FilePath == "" {
		cfg := zap.NewDevelopmentConfig()
		cfg.Level = level
		cfg.EncoderConfig.ConsoleSeparator = " "
		cfg.EncoderConfig.LineEnding = zapcore.DefaultLineEnding
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000")
		cfg.EncoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
		cfg.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		logger, err := cfg.Build(zap.AddStacktrace(zapcore.PanicLevel))
		if err != nil {
			return nil, nil, err
		}
		return &level, logger.Named(EdgeMqttBroker), nil
	}

	writeSyncer := getLogWriter(config)
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, level.Level())
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.PanicLevel))

	return &level, logger.Named(EdgeMqttBroker), nil
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.LineEnding = zapcore.DefaultLineEnding
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000")
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	encoderConfig.ConsoleSeparator = " "
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter(cfg LogConfig) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   cfg.FilePath,
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     7,
		Compress:   false,
	}
	return zapcore.AddSync(lumberJackLogger)
}
