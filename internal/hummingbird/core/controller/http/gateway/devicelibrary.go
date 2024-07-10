package gateway

import (
	"github.com/gin-gonic/gin"
	"github.com/winc-link/hummingbird/internal/dtos"
	"github.com/winc-link/hummingbird/internal/pkg/errort"
	"github.com/winc-link/hummingbird/internal/pkg/httphelper"
)

// @Tags    驱动库管理
// @Summary 新增驱动库
// @Produce json
// @Param   request body     dtos.DeviceLibraryAddRequest true "参数"
// @Success 200     {object} httphelper.CommonResponse
// @Router  /api/v1/device-libraries [post]
func (ctl *controller) DeviceLibraryAdd(c *gin.Context) {
	lc := ctl.lc
	var req dtos.DeviceLibraryAddRequest
	if err := c.ShouldBind(&req); err != nil {
		httphelper.RenderFail(c, errort.NewCommonErr(errort.DefaultReqParamsError, err), c.Writer, lc)
		return
	}
	err := ctl.getDriverLibApp().AddDriverLib(c, req)
	if err != nil {
		httphelper.RenderFail(c, err, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(nil, c.Writer, lc)
}

// @Tags    驱动库管理
// @Summary 查询驱动
// @Produce json
// @Param   request query    dtos.DeviceLibrarySearchQueryRequest true "参数"
// @Success 200     {object} httphelper.ResPageResult
// @Router  /api/v1/device-libraries [get]
func (ctl *controller) DeviceLibrariesSearch(c *gin.Context) {
	lc := ctl.lc
	var req dtos.DeviceLibrarySearchQueryRequest
	if err := c.ShouldBind(&req); err != nil {
		httphelper.RenderFail(c, errort.NewCommonErr(errort.DefaultReqParamsError, err), c.Writer, lc)
		return
	}
	dtos.CorrectionPageParam(&req.BaseSearchConditionQuery)

	list, total, edgeXErr := ctl.getDriverLibApp().DeviceLibrariesSearch(c, req)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	libs := make([]dtos.DeviceLibraryResponse, len(list))
	for i, p := range list {
		libs[i] = dtos.DeviceLibraryResponseFromModel(p)
	}
	pageResult := httphelper.NewPageResult(libs, total, req.Page, req.PageSize)

	httphelper.ResultSuccess(pageResult, c.Writer, lc)
}

// @Tags    驱动库管理
// @Summary 删除驱动
// @Produce json
// @Param   deviceLibraryId path     string true "驱动ID"
// @Success 200             {object} httphelper.CommonResponse
// @Router  /api/v1/device-libraries/:deviceLibraryId [delete]
func (ctl *controller) DeviceLibraryDelete(c *gin.Context) {
	lc := ctl.lc
	id := c.Param(UrlParamDeviceLibraryId)
	err := ctl.getDriverLibApp().DeleteDeviceLibraryById(c, id)
	if err != nil {
		httphelper.RenderFail(c, err, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(nil, c.Writer, lc)
}

// @Tags    驱动库管理
// @Summary 获取驱动定义配置信息
// @Produce json
// @Param   request query    dtos.DeviceLibraryConfigRequest true "参数"
// @Success 200     {object} httphelper.CommonResponse
// @Router  /api/v1/device-libraries/config [get]
//func (ctl *controller) DeviceLibraryConfig(c *gin.Context) {
//	lc := ctl.lc
//	var req dtos.DeviceLibraryConfigRequest
//	if err := c.ShouldBindQuery(&req); err != nil {
//		httphelper.RenderFail(c, errort.NewCommonErr(errort.DefaultReqParamsError, err), c.Writer, lc)
//		return
//	}
//	if req.DeviceLibraryId == nil && req.CloudProductId == nil && req.DeviceServiceId == nil && req.DeviceId == nil {
//		err := fmt.Errorf("deviceLibraryConfig req is null")
//		httphelper.RenderFail(c, errort.NewCommonErr(errort.DefaultReqParamsError, err), c.Writer, lc)
//		return
//	}
//
//	dl, err := ctl.getDriverLibApp().GetDeviceLibraryConfig(c, req)
//	//data, edgeXErr := gatewayapp.DeviceLibraryConfig(c, req)
//	if err != nil {
//		httphelper.RenderFail(c, err, c.Writer, lc)
//		return
//	}
//	config, err := dl.GetConfigMap()
//	if err != nil {
//		httphelper.RenderFail(c, err, c.Writer, lc)
//		return
//	}
//
//	httphelper.ResultSuccess(config, c.Writer, lc)
//}

// @Tags    驱动库管理
// @Summary 驱动库升级/下载
// @Produce json
// @Param   deviceLibraryId path     string                           true "驱动ID"
// @Param   request         query    dtos.DeviceLibraryUpgradeRequest true "参数"
// @Success 200             {object} httphelper.CommonResponse
// @Router  /api/v1/device-libraries/:deviceLibraryId/upgrade-download [put]
// Deprecated
//func (ctl *controller) DeviceLibraryUpgrade(c *gin.Context) {
//	lc := ctl.lc
//	var req dtos.DeviceLibraryUpgradeRequest
//	req.Id = c.Param(UrlParamDeviceLibraryId)
//	if err := c.ShouldBind(&req); err != nil {
//		httphelper.RenderFail(c, errort.NewCommonErr(errort.DefaultReqParamsError, err), c.Writer, lc)
//		return
//	}
//
//	err := ctl.getDriverLibApp().UpgradeDeviceLibrary(c, req, true)
//	//edgeXErr := gatewayapp.DeviceLibraryUpgrade(c, req)
//	if err != nil {
//		httphelper.RenderFail(c, err, c.Writer, lc)
//		return
//	}
//
//	httphelper.ResultSuccess(nil, c.Writer, lc)
//}

// @Tags    驱动库管理
// @Summary 上传驱动配置文件
// @Accept  multipart/form-data
// @Produce json
// @Success 200 {object} dtos.DeviceLibraryUploadResponse
// @Router  /api/v1/device-libraries/upload [post]
//func (ctl *controller) DeviceLibraryUpload(c *gin.Context) {
//	lc := ctl.lc
//	file, err := c.FormFile("fileName")
//	if err != nil {
//		httphelper.RenderFail(c, errort.NewCommonErr(errort.DefaultReqParamsError, err), c.Writer, lc)
//		return
//	}
//	var req dtos.DeviceLibraryUploadRequest
//	req.FileName = file.Filename
//	if !utils.CheckFileValid(req.FileName) {
//		err := fmt.Errorf("file name cannot contain special characters: %s", req.FileName)
//		httphelper.RenderFail(c, errort.NewCommonErr(errort.DefaultFileNotSpecialSymbol, err), c.Writer, lc)
//		return
//	}
//
//	if fileSuffix := path.Ext(req.FileName); fileSuffix != ".json" {
//		err := fmt.Errorf("file type not json, filename: %s", req.FileName)
//		httphelper.RenderFail(c, errort.NewCommonErr(errort.DefaultJsonParseError, err), c.Writer, lc)
//		return
//	}
//
//	uploadType, err := strconv.Atoi(c.Request.PostForm["type"][0])
//	if err != nil {
//		httphelper.RenderFail(c, errort.NewCommonErr(errort.DefaultReqParamsError, err), c.Writer, lc)
//		return
//	}
//
//	req.Type = uploadType
//	if req.Type != constants.DeviceLibraryUploadTypeConfig {
//		err := fmt.Errorf("req upload type %d not is %d", req.Type, constants.DeviceLibraryUploadTypeConfig)
//		httphelper.RenderFail(c, errort.NewCommonErr(errort.DefaultReqParamsError, err), c.Writer, lc)
//		return
//	}
//	f, err := file.Open()
//	if err != nil {
//		httphelper.RenderFail(c, errort.NewCommonErr(errort.DefaultReqParamsError, err), c.Writer, lc)
//		return
//	}
//
//	// 将文件流读入请求中,并校验格式
//	req.FileBytes, err = ioutil.ReadAll(f)
//	if err != nil {
//		httphelper.RenderFail(c, errort.NewCommonErr(errort.DefaultReqParamsError, err), c.Writer, lc)
//		return
//	}
//	if !json.Valid(req.FileBytes) {
//		err := fmt.Errorf("config file content must be json")
//		httphelper.RenderFail(c, errort.NewCommonErr(errort.DefaultJsonParseError, err), c.Writer, lc)
//		return
//	}
//
//	fileName, err := ctl.getDriverLibApp().UploadDeviceLibraryConfig(c, req)
//	//resp, edgeXErr := gatewayapp.DeviceLibraryUpload(c, req)
//	if err != nil {
//		httphelper.RenderFail(c, err, c.Writer, lc)
//		return
//	}
//	resp := dtos.DeviceLibraryUploadResponse{
//		FileName: fileName,
//	}
//	httphelper.ResultSuccess(resp, c.Writer, lc)
//}

//@Tags 驱动库管理
//@Summary 驱动库更新
//@Produce json
//@Param  deviceLibraryId   path  string true  "驱动ID"
//@Param  request   query  dtos.UpdateDeviceLibrary true  "参数"
//@Success 200 {object} httphelper.CommonResponse
//@Router /api/v1/device-libraries/:deviceLibraryId [put]
func (ctl *controller) DeviceLibraryUpdate(c *gin.Context) {
	lc := ctl.lc
	var req dtos.UpdateDeviceLibrary
	req.Id = c.Param(UrlParamDeviceLibraryId)
	if err := c.ShouldBind(&req); err != nil {
		httphelper.RenderFail(c, errort.NewCommonErr(errort.DefaultReqParamsError, err), c.Writer, lc)
		return
	}
	err := ctl.getDriverLibApp().UpdateDeviceLibrary(c, req)
	if err != nil {
		httphelper.RenderFail(c, err, c.Writer, lc)
		return
	}
	httphelper.ResultSuccess(nil, c.Writer, lc)
}

// @Tags    驱动库管理
// @Summary 驱动库配置下载
// @Produce application/octet-stream
// @Router  /api/v1/device-libraries/config/download [get]
//func (ctl *controller) DeviceLibraryConfigDownload(c *gin.Context) {
//	//dir, _ := os.Getwd()
//	//if dir == "/" {
//	//	dir = ""
//	//}
//	//filePath := dir + "/template/driver_config_demo.json"
//	//
//	//fileName := path.Base(filePath)
//	//c.Header("Content-Type", "application/octet-stream")
//	//c.Header("Content-Disposition", "attachment; filename="+fileName)
//	//c.File(filePath)
//	cfg := ctl.getDriverLibApp().ConfigDemo()
//	buff := bytes.NewBuffer([]byte(cfg))
//	httphelper.ResultExcelData(c, "driver_config.json", buff)
//}

//@Tags 驱动库分类
//@Summary 驱动库分类
//@Produce json
// @Param request query dtos.DriverClassifyQueryRequest true "参数"
//@Success 200 {object} httphelper.CommonResponse
//@Router /api/v1/device_classify [get]
func (ctl *controller) DeviceClassify(c *gin.Context) {
	lc := ctl.lc
	var req dtos.DriverClassifyQueryRequest
	urlDecodeParam(&req, c.Request, lc)
	dtos.CorrectionPageParam(&req.BaseSearchConditionQuery)
	list, total, edgeXErr := ctl.getDriverLibApp().GetDriverClassify(c, req)
	if edgeXErr != nil {
		httphelper.RenderFail(c, edgeXErr, c.Writer, lc)
		return
	}
	pageResult := httphelper.NewPageResult(list, total, req.Page, req.PageSize)

	httphelper.ResultSuccess(pageResult, c.Writer, lc)
}
