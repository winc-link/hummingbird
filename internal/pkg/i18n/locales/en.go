package locales

import "github.com/nicksnyder/go-i18n/v2/i18n"

func GetEnMessages() []*i18n.Message {
	return []*i18n.Message{
		{
			ID:    "10001",
			Other: `A system error has occurred while processing your request.`,
		},
		{
			ID:    "10002",
			Other: `Failed to parse JSON data.`,
		},
		{
			ID:    "10003",
			Other: `An error has occurred while processing your request parameters.{{.info | filter | left_space }}`,
		},
		{
			ID:    "10004",
			Other: `The requested resource does not exist.`,
		},
		{
			ID:    "10005",
			Other: `The file name cannot contain special characters.`,
		},
		{
			ID:    "10006",
			Other: `Request token is invalid.`,
		},
		{
			ID:    "10007",
			Other: `Request parameter name is duplicate`,
		},
		{
			ID:    "10008",
			Other: `Request parameter name cannot contain special characters`,
		},
		{
			ID:    "10009",
			Other: `Resource already exists, repeat operation`,
		},
		{
			ID:    "10010",
			Other: `The request ID is empty and the resource does not exist`,
		},
		{
			ID:    "10011",
			Other: `Upload file error`,
		},
		{
			ID:    "10012",
			Other: `Error reading Excel file`,
		},
		{
			ID:    "10013",
			Other: `Error reading Excel file, The required parameter is empty`,
		},
		{
			ID:    "10014",
			Other: `Database operation exception`,
		},
		{
			ID:    "10015",
			Other: ``,
		},
		{
			ID:    "10016",
			Other: `MQTT service connection failed`,
		},

		{
			ID:    "20101",
			Other: `The specified password is incorrect.`,
		},
		{
			ID:    "20102",
			Other: `The system has been initialized. Please log in to the system.`,
		},

		// 驱动库
		{
			ID:    "20201",
			Other: `The driver library has been instantiated. Delete the driver instance first.`,
		},
		{
			ID:    "20204",
			Other: `The driver library is being upgraded, please try later.`,
		},
		{
			ID:    "20205",
			Other: `Docker warehouse permission verification failed, please check the permissions.`,
		},
		{
			ID:    "20206",
			Other: `The docker image version does not exist.`,
		},
		//三方服务
		{
			ID:    "20207",
			Other: `The app service has been instantiated. please Delete the app service instance first.`,
		},
		{
			ID:    "20211",
			Other: `Driver library not exist.`,
		},
		{
			ID:    "20212",
			Other: `Driver library failed to get default configuration.`,
		},
		{
			ID:    "20213",
			Other: `Driver image download failed.`,
		},
		{
			ID:    "20214",
			Other: `Driver image ID not found.`,
		},
		{
			ID:    "20215",
			Other: `The built-in driver does not allow deletion.`,
		},
		{
			ID:    "20216",
			Other: `The docker image repository does not exist.`,
		},
		{
			ID:    "20217",
			Other: `Time out.`,
		},

		// 驱动实例
		{
			ID:    "20301",
			Other: `The specified driver instance has been associated with devices. Delete the sub-devices first.`,
		},
		{
			ID:    "20302",
			Other: `The specified instance is running. Stop the instance first.`,
		},
		{
			ID:    "20303",
			Other: `The specified instance is running or stopping. Please wait.`,
		},
		{
			ID:    "20304",
			Other: `The specified instance setup is not yaml format`,
		},
		{
			ID:    "20305",
			Other: `The driver failed to issue commands`,
		},
		{
			ID:    "20306",
			Other: `The driver instance is not started, and the request is not responding`,
		},
		{
			ID:    "20307",
			Other: `The specified instance is running or stopping. Please wait.`,
		},
		{
			ID:    "20308",
			Other: `Failed to start the driver instance.`,
		},
		{
			ID:    "20309",
			Other: `The driver instance not exist.`,
		},
		{
			ID:    "20310",
			Other: `The driver instance stop fail.`,
		},
		{
			ID:    "20311",
			Other: `The driver instance docker params parse err.`,
		},
		{
			ID:    "20312",
			Other: `The driver instance container name repeat.`,
		},
		{
			ID:    "20313",
			Other: `Get available port fail.`,
		},
		{
			ID:    "20314",
			Other: `Failed to create config file.`,
		},
		{
			ID:    "20315",
			Other: `The driver instance not local platform.`,
		},

		// 设备

		{
			ID:    "20406",
			Other: `Product id not found when adding device.`,
		},
		{
			ID:    "20407",
			Other: `Thing model device deletion is not supported.`,
		},
		{
			ID:    "20408",
			Other: `Response of cloud API error when bindinging sub device`,
		},
		{
			ID:    "20409",
			Other: `Response of cloud API error when updating device attributes`,
		},
		{
			ID:    "20410",
			Other: `Device not exist`,
		},
		{
			ID:    "20411",
			Other: `Device command not exist`,
		},
		{
			ID:    "20412",
			Other: `Local devices are not allowed to connect to the platform`,
		},
		{
			ID:    "20413",
			Other: `Please unbind the device with the driver first`,
		},
		{
			ID:    "20414",
			Other: `The device platform is inconsistent with the drive platform`,
		},
		{
			ID:    "20415",
			Other: `This device has been bound to alarm rules. Please stop reporting relevant alarm rules before proceeding with the operation`,
		},
		{
			ID:    "20416",
			Other: `This device has been bound to scene rules. Please stop reporting scene rules before proceeding with the operation`,
		},

		// 产品
		{
			ID:    "20601",
			Other: `The specified product has been associated with DPs. Delete the DPs first.`,
		},
		{
			ID:    "20602",
			Other: `The specified product has been associated with sub-devices. Delete the sub-devices first.`,
		},
		{
			ID:    "20603",
			Other: `No active device found with the product, please add device and activate it.`,
		},
		{
			ID:    "20604",
			Other: `Product not exist.`,
		},

		{
			ID:    "20605",
			Other: `Sync failed.`,
		},

		{
			ID:    "20606",
			Other: `Product resource not exist.`,
		},

		{
			ID:    "20607",
			Other: `Please sync product data.`,
		},
		{
			ID:    "20608",
			Other: `Product property code exist.`,
		},
		{
			ID:    "20609",
			Other: `Product event code exist.`,
		},
		{
			ID:    "20610",
			Other: `Product server code exist.`,
		},
		{
			ID:    "20611",
			Other: `This product has been bound to alarm rules. Please stop reporting relevant alarm rules before proceeding with the operation.`,
		},
		{
			ID:    "20612",
			Other: `The product has not been released yet. Please release the product before adding devices`,
		},
		{
			ID:    "20613",
			Other: `Please cancel publishing the product first before proceeding with the operation`,
		},
		{
			ID:    "20614",
			Other: `Code identifier already exists`,
		},
		{
			ID:    "20616",
			Other: `The data type does not support modification. Please delete and recreate`,
		},

		// docker config
		{
			ID:    "20701",
			Other: `Please delete the bound driver library first`,
		},
		{
			ID:    "20702",
			Other: `Driver repository configuration does not exist`,
		},

		//专家系统
		{
			ID:    "20802",
			Other: `There is this scene in the timed task, please delete the timed task bound to this scene first`,
		},
		{
			ID:    "20803",
			Other: `There is content in the smart scene, please delete the content first`,
		},
		{
			ID:    "20831",
			Other: `There is content in the linkage strategy, please delete the content first`,
		},
		{
			ID:    "20832",
			Other: `There is content in the linkage strategy, please delete the content first`,
		},
		{
			ID:    "20833",
			Other: `Scheduled tasks that already exist`,
		},
		{
			ID:    "20860",
			Other: `The data type gateway is not yet supported`,
		},
		{
			ID:    "20861",
			Other: `Invalid job type`,
		},
		{
			ID:    "20862",
			Other: `Invalid action type`,
		},

		// 物模型
		{
			ID:    "20901",
			Other: "sync thing model failed",
		},
		{
			ID:    "20902",
			Other: "subDevice needed to be bond when sync thing model",
		},
		{
			ID:    "20903",
			Other: "parameters error or incomplete",
		},

		// 功能点
		{
			ID:    "21001",
			Other: `Sync function point，result return error.`,
		},
		{
			ID:    "21002",
			Other: `Sync function point，result return fail.`,
		},
		{
			ID:    "21003",
			Other: `Sync function point，parse fail.`,
		},
		{
			ID:    "21004",
			Other: `Function point not exist.`,
		},

		{
			ID:    "20901",
			Other: "sync thing model failed",
		},
		{
			ID:    "20902",
			Other: "subDevice needed to be bond when sync thing model",
		},
		{
			ID:    "20903",
			Other: "parameters error or incomplete",
		},
		// ota
		{
			ID:    "21100",
			Other: "OTA upgrade error",
		},
		{
			ID:    "21101",
			Other: "Upgrading the gateway",
		},
		{
			ID:    "21102",
			Other: "OTA firmware download url error",
		},
		{
			ID:    "21103",
			Other: "OTA firmware download error",
		},
		{
			ID:    "21104",
			Other: "OTA verify firmware error",
		},
		{
			ID:    "21105",
			Other: "OTA execute upgrade command error",
		},
		{
			ID:    "21106",
			Other: "OTA upgrade firmware error",
		},
		{
			ID:    "21107",
			Other: "gateway does not support OTA",
		},
		{
			ID:    "21108",
			Other: "OTA get upgrade version error",
		},
		{
			ID:    "21109",
			Other: "OTA not upgrade version",
		},
		{
			ID:    "21110",
			Other: "OTA upgrade version wrong",
		},
		//category
		{
			ID:    "21200",
			Other: "The category does not exist.",
		},
		{
			ID:    "21201",
			Other: "The thingModel does not exist.",
		},

		//alert rule
		{
			ID:    "21300",
			Other: "Rule alert does not exist.",
		},
		{
			ID:    "21301",
			Other: "Invalid rule json.",
		},
		{
			ID:    "21302",
			Other: "Rule engine not found rule.",
		},
		{
			ID:    "21303",
			Other: "Rule is starting.",
		},
		{
			ID:    "21304",
			Other: "The product or device has been modified. Please edit the rule again.",
		},
		{
			ID:    "21305",
			Other: "Parameter error, please edit the rule again.",
		},
		{
			ID:    "21306",
			Other: "The format of the effective time is incorrect. The end time should be greater than the start time.",
		},
		{
			ID:    "21307",
			Other: "Please stop this rule before editing it.",
		},
		//scene
		{
			ID:    "21400",
			Other: "Please stop this scheduled tasks before editing it.",
		},
		{
			ID:    "21402",
			Other: "Parameter error, please edit the scene again.",
		},
		//rule
		{
			ID:    "21500",
			Other: "Please stop this rule engine before editing it.",
		},
		{
			ID:    "21600",
			Other: "Invalid resource configuration, please check the resource configuration.",
		},
		/***********************************************************/
		/*********************** string 翻译 ***********************/
		/**********************************************************/
		{
			ID:    "agentclient.alert.notify",
			Other: "notify",
		},
		{
			ID:    "agentclient.alert.warn",
			Other: "warn",
		},
		{
			ID:    "agentclient.alert.error",
			Other: "error",
		},
		{
			ID:    "agentclient.alert.system_alert",
			Other: "system_alert",
		},
		{
			ID:    "agentclient.alert.driver_alert",
			Other: "driver_alert",
		},
		{
			ID:    "20901",
			Other: "sync thing model failed",
		},
		{
			ID:    "20902",
			Other: "subDevice needed to be bond when sync thing model",
		},

		{
			ID:    "default.success",
			Other: "Success",
		},
		{
			ID:    "default.fail",
			Other: "Fail",
		},
		{
			ID:    "library.upgrade_download",
			Other: "Driver Library {{.name}} Install/Upgrade {{.status}}",
		},
		{
			ID:    "service.run_status",
			Other: "Driver Service {{.name}} Run/Stop {{.status}}",
		},
		{
			ID:    "appService.run_status",
			Other: "service {{.name}} Install/Upgrade{{.status}}",
		},
		{
			ID:    "appServiceLibrary.upgrade_download",
			Other: "service instance {{.name}} Run/Stop{{.status}}",
		},
		{
			ID:    "appService.delete",
			Other: "service instance {{.name}} Delete{{.status}}",
		},
		{
			ID:    "cloudInstance.log_status",
			Other: "service instance {{.name}} log Run/Stop{{.status}}",
		},
	}
}
