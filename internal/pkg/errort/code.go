package errort

const (
	DefaultSuccess                          uint32 = 00000
	DefaultSystemError                      uint32 = 10001
	DefaultJsonParseError                   uint32 = 10002
	DefaultReqParamsError                   uint32 = 10003
	DefaultResourcesNotFound                uint32 = 10004
	DefaultFileNotSpecialSymbol             uint32 = 10005
	DefaultTokenPermission                  uint32 = 10006
	DefaultNameRepeat                       uint32 = 10007
	DefaultNameSpecialCharacters            uint32 = 10008
	DefaultResourcesRepeat                  uint32 = 10009
	DefaultIdEmpty                          uint32 = 10010
	DefaultUploadFileErrorCode                     = 10011
	DefaultReadExcelErrorCode                      = 10012
	DefaultReadExcelErrorParamsRequiredCode        = 10013
	KindDatabaseError                              = 10014
	KindUnknown                                    = 10015
	MqttConnFail                                   = 10016

	AppPasswordError     uint32 = 20101
	AppSystemInitialized uint32 = 20102

	DeviceLibraryMustDeleteDeviceService uint32 = 20201
	DeviceLibraryUpgradeIng              uint32 = 20204
	DeviceLibraryDockerAuthInvalid       uint32 = 20205
	DeviceLibraryDockerImagesNotFound    uint32 = 20206
	DeviceLibraryNotExist                uint32 = 20211
	DeviceLibraryImageDownloadFail       uint32 = 20213
	DeviceLibraryImageNotFound           uint32 = 20214
	DeviceLibraryNotAllowDelete          uint32 = 20215
	DockerImageRepositoryNotFound        uint32 = 20216
	DeviceLibraryResponseTimeOut         uint32 = 20217

	DeviceServiceMustDeleteDevice     uint32 = 20301
	DeviceServiceMustStopService      uint32 = 20302
	DeviceServiceMustStopDoingService uint32 = 20303
	DeviceServiceSetupYamlFormatError uint32 = 20304
	DeviceServiceSendCommandFail      uint32 = 20305
	DeviceServiceNotStarted           uint32 = 20306
	AppServiceMustStopDoingService           = 20307
	AppServiceMustStopService                = DeviceServiceMustStopService
	ContainerRunFail                         = 20308
	DeviceServiceNotExist                    = 20309
	ContainerStopFail                        = 20310
	DockerParamsParseErr                     = 20311
	DeviceServiceContainerNameRepeat         = 20312
	GetAvailablePortFail                     = 20313
	CreateConfigFileFail                     = 20314
	DeviceServiceMustLocalPlatform           = 20315

	DeviceProductIdNotFound             uint32 = 20406
	DeviceDeleteNotAllowed              uint32 = 20407
	DeviceBindAtopFailed                uint32 = 20408
	DeviceUpdateAtopFailed              uint32 = 20409
	DeviceNotExist                             = 20410
	DeviceCommandNotExist                      = 20411
	DeviceNotAllowConnectPlatform              = 20412
	DeviceNotUnbindDriver                      = 20413
	DeviceAndDriverPlatformNotIdentical        = 20414
	DeviceAssociationAlertRule                 = 20415
	DeviceAssociationSceneRule                 = 20416

	// 产品
	ProductMustDeleteDevice       uint32 = 20602
	ProductNotExist               uint32 = 20604
	ProductPropertyCodeNotExist   uint32 = 20608
	ProductAssociationAlertRule          = 20611
	ProductUnRelease                     = 20612
	ProductRelease                       = 20613
	ThingModelCodeExist                  = 20614
	ThingModeTypeCannotBeModified        = 20616

	// 镜像仓库
	DockerConfigMustDeleteDeviceLibrary uint32 = 20701
	DockerConfigNotExist                uint32 = 20702

	CategoryNotExist   uint32 = 21200
	ThingModelNotExist uint32 = 21201

	AlertRuleNotExist              uint32 = 21300
	InvalidRuleJson                uint32 = 21301
	EkuiperNotFindRule             uint32 = 21302
	AlertRuleStatusStarting        uint32 = 21303
	AlertRuleProductOrDeviceUpdate uint32 = 21304
	AlertRuleParamsError           uint32 = 21305
	EffectTimeParamsError          uint32 = 21306
	StopAlertRule                  uint32 = 21307

	SceneTimerIsStartingNotAllowUpdate = 21400

	SceneRuleParamsError uint32 = 21402

	RuleEngineIsStartingNotAllowUpdate = 21500

	InvalidSource = 21600

	CloudServiceConnectionRefused uint32 = 22101
)

type OpenApiErrorCode uint32

const (
	Success                OpenApiErrorCode = 0
	SystemErrorCode                         = 500
	ParamsError                             = 1104
	FunctionNotSupportCode                  = 2003
	TokenValid                              = 1011
	TokenExpired                            = 1010
	UrlPathIsInvalid                        = 1108
)

type OpenApiErrorMsg string

const (
	SystemErrorMsg        OpenApiErrorMsg = "system error,please contact the admin"
	FunctionNotSupportMsg                 = "function not support"
	ParamsErrorMsg                        = "params incorrect"
	TokenValidMsg                         = "token invalid"
	TokenExpiredMsg                       = "token expired"
	UrlPathIsInvalidMsg                   = "url path is invalid"
)

var OpenApiCodeMsgMap = map[OpenApiErrorCode]OpenApiErrorMsg{
	SystemErrorCode:        SystemErrorMsg,
	FunctionNotSupportCode: FunctionNotSupportMsg,
	ParamsError:            ParamsErrorMsg,
	TokenValid:             TokenValidMsg,
	TokenExpired:           TokenExpiredMsg,
	UrlPathIsInvalid:       UrlPathIsInvalidMsg,
}
