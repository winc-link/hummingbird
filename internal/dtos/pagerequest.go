package dtos

type ReqPageInfo struct {
	Page     int `json:"page" form:"page"`
	PageSize int `json:"page_size" form:"page_size"`
}

type ResPageResult struct {
	List     interface{} `json:"list"`
	Total    uint32      `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
}

func NewPageResult(responses interface{}, total uint32, page int, pageSize int) ResPageResult {
	if responses == nil {
		responses = make([]interface{}, 0)
	}
	return ResPageResult{
		List:     responses,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}
}

type CommonResponse struct {
	Success   bool        `json:"success"`
	ErrorMsg  string      `json:"errorMsg"`
	ErrorCode int         `json:"errorCode"`
	Result    interface{} `json:"result"`
}

const (
	PageDefault        = 1
	PageSizeDefault    = 10
	PageSizeMaxDefault = 1000
)

// 校验并设置PageRequest参数
func CorrectionPageRequest(query *PageRequest) {
	if query.Page <= 0 {
		query.Page = PageDefault
	}

	if query.PageSize >= PageSizeMaxDefault {
		query.PageSize = PageSizeMaxDefault
	} else if query.PageSize <= 0 {
		query.PageSize = PageSizeDefault
	}
}

// 校验并设置page参数
func CorrectionPageParam(query *BaseSearchConditionQuery) {
	if query.Page <= 0 {
		query.Page = PageDefault
	}

	if query.PageSize >= PageSizeMaxDefault {
		query.PageSize = PageSizeMaxDefault
	} else if query.PageSize <= 0 {
		query.PageSize = PageSizeDefault
	}
}
