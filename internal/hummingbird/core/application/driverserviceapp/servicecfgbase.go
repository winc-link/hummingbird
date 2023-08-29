package driverserviceapp

import (
	//"gitlab.com/tedge/edgex/internal/models"
	//"gitlab.com/tedge/edgex/internal/pkg/constants"
	"github.com/winc-link/hummingbird/internal/models"
)

// 驱动运行的配置模版
func getDriverConfigTemplate(ds models.DeviceService) string {

	return getDefaultDriverConfig()
}

func getDefaultDriverConfig() string {
	return `[Logger]
FileName = "/mnt/logs/driver.log"
LogLevel = "INFO" # DEBUG INFO WARN ERROR

[Clients]
[Clients.Core]
Address  = "hummingbird-core:57081"
UseTLS = false
CertFilePath = ""

[Service]
ID = ""
Name = ""
ProductList = []
GwId = ""
LocalKey = ""
Activated = false
[Service.Server]
Address = "0.0.0.0:49991"
UseTLS = false
CertFile = ""
KeyFile = ""

[CustomConfig]`
}
