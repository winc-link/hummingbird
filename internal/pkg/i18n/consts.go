package i18n

const (
	AgentAlertNotify = "agentclient.alert.notify"
	AgentAlertWarn   = "agentclient.alert.warn"
	AgentAlertError  = "agentclient.alert.error"
	AgentAlertSystem = "agentclient.alert.system_alert"
	AgentAlertDriver = "agentclient.alert.driver_alert"

	LicenseAlertExpire = "agentclient.alert.license_alert"

	DefaultSuccess                       = "default.success"
	DefaultFail                          = "default.fail"
	LibraryUpgradeDownloadResp           = "library.upgrade_download"           // 驱动库下载结果
	ServiceRunStatusResp                 = "service.run_status"                 // 驱动实例运行结果
	AppServiceLibraryUpgradeDownloadResp = "appServiceLibrary.upgrade_download" //三方服务下载结果
	AppServiceRunStatusResp              = "appService.run_status"              // 三方服务实例运行结果
	AppServiceDeleteResp                 = "appService.delete"                  // 三方服务实例运行结果
	CloudInstanceLogResp                 = "cloudInstance.log_status"           // 云服务日志状态

)

const (
	ContentType     = "Content-Type"
	ContentTypeCBOR = "application/cbor"
	ContentTypeJSON = "application/json"
	ContentTypeYAML = "application/x-yaml"
	ContentTypeText = "text/plain"
	ContentTypeXML  = "application/xml"
	AcceptLanguage  = "Accept-Language"
)
