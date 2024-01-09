package mqttd

var (
	DefaultConfigDir = "./res/"
)

func GetDefaultConfigDir() (string, error) {
	return DefaultConfigDir, nil
}
