package limit

import (
	"context"
	"fmt"
	"testing"
	"time"
)

var (
	srv    LimitService
	option LimitOption
)

func init() {
	option = LimitOption{
		LimitMillisecond: 100,
		Burst:            1,
	}
	srv = NewLimitService(option)
}

func TestLimitServiceProvider_Consume(t *testing.T) {
	for i := 0; i < 10; i++ {
		if err := srv.Consume(); err != nil {
			fmt.Printf("%v consume failed, error: %v\n", i, err.Error())
		} else {
			fmt.Printf("%v consume success\n", i)
		}
		// sleep time比option.limitMillisecond的值略少一点是为了构造error情况
		time.Sleep(time.Duration(option.LimitMillisecond-1) * time.Millisecond)
	}
}

func TestLimitServiceProvider_ConsumeWithContext(t *testing.T) {
	var blag int64 = 0
	for i := 0; i < 10; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(option.LimitMillisecond-blag)*time.Millisecond)
		defer cancel()
		if err := srv.ConsumeWithContext(ctx); err != nil {
			fmt.Printf("%v consume failed, error: %v\n", i, err.Error())
		} else {
			fmt.Printf("%v consume success\n", i)
		}
		blag = 1 - blag
	}
}
