package limit

import "context"

type LimitMethodConf interface {
	// 获取需要限流的接口方法名
	GetLimitMethods() map[string]struct{}
}

// 限流服务接口
type LimitService interface {
	// 查询option值
	GetOption() LimitOption
	// 修改option结构体内的参数，支持实时更新
	SetOption(option LimitOption)
	// 直接消费，如果请求过于频繁，则直接拒绝，返回error
	Consume() error
	// 阻塞消费，如果请求过于频繁，通过ctx设置的Deadline或者Timeout来等待一段时间，若资源不够，亦返回error
	ConsumeWithContext(ctx context.Context) error
	// 克隆复制一个新的限流服务接口
	Clone() LimitService
}

type LimitOption struct {
	// 令牌发放的频率，单位毫秒
	LimitMillisecond int64
	// 每次发放的个数
	Burst int
}
