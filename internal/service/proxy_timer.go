package service

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcron"
	"github.com/gogf/gf/v2/os/gctx"
	"golang.org/x/net/context"
)

func Timer() error {
	var (
		err error
		ctx = gctx.New()
	)
	// 天气代理定时器
	_, err = gcron.Add(ctx, "@every 3m", func(ctx context.Context) {
		if err := WeatherProxy(); err != nil {
			g.Log().Error(ctx, err)
		}
	})
	if err != nil {
		return err
	}
	return nil
}
