package service

import (
	"github.com/gogf/gf/v2/os/gcron"
	"github.com/gogf/gf/v2/os/gctx"
	"golang.org/x/net/context"
)

func Timer() {
	var (
		err error
		ctx = gctx.New()
	)
	// 天气代理定时器
	_, err = gcron.Add(ctx, "@every 3m", func(ctx context.Context) {
		go WeatherProxy()
	})
	if err != nil {
		panic(err)
	}
}
