package locales

import "github.com/nicksnyder/go-i18n/v2/i18n"

func GetZhMessages() []*i18n.Message {
	return []*i18n.Message{
		{
			ID:    "10001",
			Other: `系统错误`,
		},
		{
			ID:    "10002",
			Other: `json 解析失败`,
		},
		{
			ID:    "10003",
			Other: `请求参数错误{{.info | filter | left_space }}`,
		},
		{
			ID:    "10004",
			Other: `请求资源不存在`,
		},
		{
			ID:    "10005",
			Other: `文件名不能包含特殊符号`,
		},
		{
			ID:    "10006",
			Other: `token错误`,
		},
		{
			ID:    "10007",
			Other: `请求参数name重复`,
		},
		{
			ID:    "10008",
			Other: `请求参数name不能含有特殊字符`,
		},
		{
			ID:    "10009",
			Other: `资源已存在，重复操作`,
		},
		{
			ID:    "10010",
			Other: `请求ID为空，资源不存在`,
		},
		{
			ID:    "10011",
			Other: `上传文件错误`,
		},
		{
			ID:    "10012",
			Other: `读取Excel文件错误`,
		},
		{
			ID:    "10013",
			Other: `读取Excel文件错误,必填参数为空`,
		},
		{
			ID:    "10014",
			Other: `数据库操作异常`,
		},
		{
			ID:    "10015",
			Other: ``,
		},
		{
			ID:    "10016",
			Other: `MQTT服务连接失败`,
		},

		{
			ID:    "20101",
			Other: `账号或密码错误`,
		},
		{
			ID:    "20102",
			Other: `系统已经初始化，请登陆`,
		},

		// 驱动库
		{
			ID:    "20201",
			Other: `该驱动库已实例化，请优先删除驱动实例`,
		},
		{
			ID:    "20204",
			Other: `驱动库正在升级，请稍后尝试`,
		},
		{
			ID:    "20205",
			Other: `docker仓库权限校验失败，请检查权限`,
		},
		{
			ID:    "20206",
			Other: `docker镜像版本不存在`,
		},
		// 三方服务
		{
			ID:    "20207",
			Other: `该服务已创建实例，请先删除服务实例`,
		},
		{
			ID:    "20211",
			Other: `驱动不存在`,
		},
		{
			ID:    "20212",
			Other: `驱动配置获取失败`,
		},
		{
			ID:    "20213",
			Other: `驱动镜像下载失败`,
		},
		{
			ID:    "20214",
			Other: `驱动镜像ID不存在`,
		},
		{
			ID:    "20215",
			Other: `内置驱动不允许删除`,
		},
		{
			ID:    "20216",
			Other: `docker镜像仓库不存在`,
		},
		{
			ID:    "20217",
			Other: `响应超时`,
		},

		// 驱动实例
		{
			ID:    "20301",
			Other: `该驱动实例已绑定设备，请优先删除子设备`,
		},
		{
			ID:    "20302",
			Other: `实例正在运行，请先停止运行实例`,
		},
		{
			ID:    "20303",
			Other: `实例正在启动/停止中，请稍后`,
		},
		{
			ID:    "20304",
			Other: `驱动实例自定义配制参数格式不对，请用yaml格式`,
		},
		{
			ID:    "20305",
			Other: `驱动下发指令失败`,
		},
		{
			ID:    "20306",
			Other: `驱动实例未启动，请求无响应`,
		},
		{
			ID:    "20308",
			Other: `驱动实例启动失败`,
		},
		{
			ID:    "20309",
			Other: `实例不存在`,
		},
		{
			ID:    "20310",
			Other: `实例停止异常`,
		},
		{
			ID:    "20311",
			Other: `docker配置参数错误`,
		},
		{
			ID:    "20312",
			Other: `驱动容器名重复`,
		},
		{
			ID:    "20313",
			Other: `获取可用端口失败`,
		},
		{
			ID:    "20314",
			Other: `驱动配置文件创建失败`,
		},
		{
			ID:    "20315",
			Other: `驱动关联的平台非本地`,
		},

		// 设备
		{
			ID:    "20406",
			Other: `该设备产品ID未找到.`,
		},
		{
			ID:    "20407",
			Other: `物模型设备删除不支持`,
		},
		{
			ID:    "20408",
			Other: `子设备绑定/解绑时云端接口返回错误`,
		},
		{
			ID:    "20409",
			Other: `更新云端设备属性时云端接口返回错误`,
		},
		{
			ID:    "20410",
			Other: `设备不存在`,
		},
		{
			ID:    "20411",
			Other: `设备指令不存在`,
		},
		{
			ID:    "20412",
			Other: `本地设备不允许连接云平台`,
		},
		{
			ID:    "20413",
			Other: `设备未与驱动解绑，请先解绑在进行绑定`,
		},
		{
			ID:    "20414",
			Other: `设备平台与驱动平台不一致`,
		},
		{
			ID:    "20414",
			Other: `设备平台与驱动平台不一致`,
		},
		{
			ID:    "20415",
			Other: `该设备已与告警规则绑定，请停止告相关警规则，再进行操作`,
		},
		{
			ID:    "20416",
			Other: `该设备已与场景联动绑定，请停止场景联动规则，再进行操作`,
		},

		// 产品
		{
			ID:    "20601",
			Other: `该产品已绑定功能点，请优先删除子功能点`,
		},
		{
			ID:    "20602",
			Other: `该产品已绑定设备，请优先删除设备`,
		},
		{
			ID:    "20603",
			Other: `请先为该产品添加设备并激活`,
		},
		{
			ID:    "20604",
			Other: `产品不存在`,
		},

		{
			ID:    "20605",
			Other: `同步失败`,
		},

		{
			ID:    "20606",
			Other: `云端产品被删除`,
		},

		{
			ID:    "20607",
			Other: `请先同步产品数据`,
		},
		{
			ID:    "20608",
			Other: `产品属性标识符不存在`,
		},
		{
			ID:    "20609",
			Other: `产品事件标识符不存在`,
		},
		{
			ID:    "20610",
			Other: `产品服务标识符不存在`,
		},
		{
			ID:    "20611",
			Other: `该产品已与告警规则绑定，请停止告相关警规则，再进行操作`,
		},
		{
			ID:    "20612",
			Other: `产品还未发布，请先发布产品，再添加设备`,
		},
		{
			ID:    "20613",
			Other: `请先取消发布产品，在进行操作`,
		},
		{
			ID:    "20614",
			Other: `标识符已重复，请修改标识符`,
		},
		{
			ID:    "20616",
			Other: `数据类型不支持修改，请删除重新建立`,
		},

		// docker config
		{
			ID:    "20701",
			Other: `请先删除绑定的驱动库`,
		},
		{
			ID:    "20702",
			Other: `仓库配置不存在`,
		},

		//专家系统
		{
			ID:    "20802",
			Other: `定时任务中有此场景，请先删除绑定该场景的定时任务`,
		},
		{
			ID:    "20803",
			Other: `智能场景中有内容，请先删除内容`,
		},
		{
			ID:    "20831",
			Other: `联动策略中有内容，请先删除内容`,
		},
		{
			ID:    "20832",
			Other: `联动策略中有内容，请先删除内容`,
		},
		{
			ID:    "20833",
			Other: `已存在执行的定时任务`,
		},
		{
			ID:    "20860",
			Other: `尚不支持该数据类型网关`,
		},
		{
			ID:    "20861",
			Other: `无效的工作类型`,
		},
		{
			ID:    "20862",
			Other: `无效的动作类型`,
		},

		// 物模型
		{
			ID:    "20901",
			Other: "同步物模型失败",
		},
		{
			ID:    "20902",
			Other: "同步物模型需绑定子设备",
		},
		{
			ID:    "20903",
			Other: "参数错误",
		},

		// 功能点
		{
			ID:    "21001",
			Other: `同步功能点，结果返回错误`,
		},
		{
			ID:    "21002",
			Other: `同步功能点，结果返回失败`,
		},
		{
			ID:    "21003",
			Other: `同步功能点，解析功能点失败`,
		},
		{
			ID:    "21004",
			Other: `功能点不存在`,
		},

		// ota
		{
			ID:    "21100",
			Other: "OTA升级错误",
		},
		{
			ID:    "21101",
			Other: "网关升级中",
		},
		{
			ID:    "21102",
			Other: "OTA固件地址错误",
		},
		{
			ID:    "21103",
			Other: "OTA固件下载失败",
		},
		{
			ID:    "21104",
			Other: "OTA固件校验失败",
		},
		{
			ID:    "21105",
			Other: "OTA执行升级命令失败",
		},
		{
			ID:    "21106",
			Other: "OTA固件升级失败",
		},
		{
			ID:    "21107",
			Other: "网关不支持 OTA",
		},
		{
			ID:    "21108",
			Other: "获取OTA升级版本失败",
		},
		{
			ID:    "21109",
			Other: "无可升级版本",
		},
		{
			ID:    "21110",
			Other: "OTA固件版本错误",
		},

		{
			ID:    "22101",
			Other: "连接拒绝，请确保云服务是启动状态",
		},

		//category
		{
			ID:    "21200",
			Other: "产品品类不存在.",
		},
		{
			ID:    "21201",
			Other: "产品物模型不存在.",
		},

		{
			ID:    "21300",
			Other: "告警规则不存在.",
		},
		{
			ID:    "21301",
			Other: "无效的规则.",
		},
		{
			ID:    "21302",
			Other: "规则引擎未找到该规则.",
		},
		{
			ID:    "21303",
			Other: "该规则已启动，请勿重复操作.",
		},
		{
			ID:    "21304",
			Other: "产品或设备已被修改，请从新编辑该规则.",
		},
		{
			ID:    "21305",
			Other: "参数错误，请从新编辑该规则.",
		},
		{
			ID:    "21306",
			Other: "生效时间格式错误，结束时间应该大于开始时间.",
		},
		{
			ID:    "21307",
			Other: "请停止此规则，在进行编辑.",
		},
		//scene
		{
			ID:    "21400",
			Other: "请停止此定时任务，在进行编辑.",
		},
		{
			ID:    "21402",
			Other: "参数错误，请从新编辑该场景.",
		},
		//rule
		{
			ID:    "21500",
			Other: "请停止规则引擎，在进行操作.",
		},
		{
			ID:    "21600",
			Other: "无效的资源配置，请检查资源配置是否争取.",
		},

		/***********************************************************/
		/*********************** string 翻译 ***********************/
		/**********************************************************/

		// agentclient alert
		{
			ID:    "agentclient.alert.notify",
			Other: "通知",
		},
		{
			ID:    "agentclient.alert.warn",
			Other: "警告",
		},
		{
			ID:    "agentclient.alert.error",
			Other: "错误",
		},
		{
			ID:    "agentclient.alert.system_alert",
			Other: "系统告警",
		},
		{
			ID:    "agentclient.alert.driver_alert",
			Other: "驱动告警",
		},

		// import device template
		{
			ID:    "device.import.Cid",
			Other: "设备id",
		},
		{
			ID:    "device.import.Name",
			Other: "设备名",
		},
		{
			ID:    "device.import.Description",
			Other: "描述",
		},
		{
			ID:    "device.import.Ip",
			Other: "ip地址",
		},
		{
			ID:    "device.import.Lat",
			Other: "纬度",
		},
		{
			ID:    "device.import.Lon",
			Other: "经度",
		},
		{
			ID:    "device.import.VendorCode",
			Other: "厂商",
		},
		{
			ID:    "device.import.InstallLocation",
			Other: "安装地址",
		},
		{
			ID:    "device.import.ExtendData",
			Other: "扩展字段",
		},
		{
			ID:    "device.import.IsIPC",
			Other: "是否为IPC设备[false/true]",
		},
		{
			ID:    "20901",
			Other: "同步物模型失败",
		},
		{
			ID:    "20902",
			Other: "同步物模型需绑定子设备",
		},

		/*********************** string 翻译 ***********************/
		{
			ID:    "default.success",
			Other: "成功",
		},
		{
			ID:    "default.fail",
			Other: "失败",
		},
		{
			ID:    "library.upgrade_download",
			Other: "驱动库 {{.name}} 安装/升级{{.status}}",
		},
		{
			ID:    "service.run_status",
			Other: "驱动实例 {{.name}} 运行/停止{{.status}}",
		},
		{
			ID:    "appServiceLibrary.upgrade_download",
			Other: "服务实例 {{.name}} 安装/升级{{.status}}",
		},
		{
			ID:    "appService.run_status",
			Other: "服务 {{.name}} 运行/停止{{.status}}",
		},
		{
			ID:    "appService.delete",
			Other: "服务 {{.name}} 删除{{.status}}",
		},
		{
			ID:    "cloudInstance.log_status",
			Other: "服务 {{.name}} 日志运行/停止{{.status}}",
		},
	}
}
