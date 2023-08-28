package limit

import (
	"context"
	"errors"
	"time"

	"golang.org/x/time/rate"
)

const (
	ConsumeBusyError    = "consume failed because your request is too busy"
	ConsumeTimeoutError = "consume failed because request deadline or waited timeout"
)

type LimitServiceProvider struct {
	limiter *rate.Limiter
}

func NewLimitService(option LimitOption) LimitService {
	return &LimitServiceProvider{
		limiter: rate.NewLimiter(
			rate.Every(time.Duration(option.LimitMillisecond)*time.Millisecond),
			option.Burst,
		),
	}
}

func (provider *LimitServiceProvider) GetOption() LimitOption {
	return LimitOption{
		Burst:            provider.limiter.Burst(),
		LimitMillisecond: int64(1000 / provider.limiter.Limit()),
	}
}

func (provider *LimitServiceProvider) SetOption(option LimitOption) {
	provider.limiter.SetLimit(rate.Every(time.Duration(option.LimitMillisecond) * time.Millisecond))
	provider.limiter.SetBurst(option.Burst)
}

func (provider *LimitServiceProvider) Consume() error {
	if provider.limiter.Allow() {
		return nil
	}
	return errors.New(ConsumeBusyError)
}

func (provider *LimitServiceProvider) ConsumeWithContext(ctx context.Context) error {
	if err := provider.limiter.Wait(ctx); err != nil {
		return errors.New(ConsumeTimeoutError)
	}
	return nil
}

func (provider *LimitServiceProvider) Clone() LimitService {
	return NewLimitService(provider.GetOption())
}
