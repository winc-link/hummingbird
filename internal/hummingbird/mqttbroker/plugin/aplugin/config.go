package aplugin

// Config is the configuration for the aplugin plugin.
type Config struct {
	// add your config fields
}

// Validate validates the configuration, and return an error if it is invalid.
func (c *Config) Validate() error {
	return nil
}

// DefaultConfig is the default configuration.
var DefaultConfig = Config{}

func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return nil
}
