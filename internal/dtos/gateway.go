package dtos

type NetIface struct {
	Ifaces []string `json:"ifaces"`
}

type EdgeBaseConfig struct {
}

type EdgeGwConfig struct {
	GwId     string `json:"gwId"`
	SecKey   string `json:"secKey"`
	LocalKey string `json:"localKey"`
	Status   bool   `json:"status"`
}

type EdgeConfig struct {
	//BaseConfig     EdgeBaseConfig `yaml:"baseconfig"`
	//GwConfig       EdgeGwConfig   `yaml:"gwconfig"`
	//SubDeviceLimit int64          `yaml:"subdevicelimit"`
	//ExpiryTime     int64          `yaml:"expiry"`
	//ActiveTime     int64          `yaml:"activeTime"`
	//LastExitTime   int64          `yaml:"lastExitTime"`
	//IsExpired      bool           `yaml:"isExpired"`

	GwId           string `yaml:"gwid"`
	SecKey         string `yaml:"seckey"`
	Status         bool   `yaml:"status"`
	ActiveTime     string `yaml:"activetime"`
	VersionNumber  string `yaml:"versionnumber"`
	SubDeviceLimit int64  `yaml:"subdevicelimit"`
}

func (c EdgeConfig) GetGatewayNumber() string {
	switch c.VersionNumber {
	case "ireland":
		return "Ireland（爱尔兰）"
	case "seattle":
		return "Seattle（西雅图）"
	case "kamakura（镰仓）":
		return "Kamakura"
	default:
		return c.VersionNumber
	}
}

func (c EdgeConfig) IsActivated() bool {
	return c.Status
}

func (c EdgeConfig) CheckThingModelActiveGw() bool {
	return c.GwId != "" && c.SecKey != ""
}
