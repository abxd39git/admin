package main

import (
	"fmt"
	"os"

	"admin/app"
	"admin/middleware"
	"admin/session"
	"admin/utils"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"admin/cron"
)

func main() {
	if os.Getenv("ADMIN_API_ENV") == "" {
		panic("环境变量ADMIN_API_ENV未设置")
	}

	// 定时任务
	cron.InitCron()

	//启动定时器
	//go new(models.TokenFeeDailySheet).BoottimeTimingSettlement()
	//go new(models.WalletInoutDailySheet).BoottimeTimingSettlement()
	//go new(models.CurencyFeeDailySheet).BoottimeTimingSettlement()

	// 配置gin
	r := gin.Default()

	// session
	r.Use(sessions.Sessions("mysession", session.Store))

	// custom middleware
	r.Use(middleware.JsCors())
	r.Use(middleware.CheckLogin())
	app.Router(r)

	r.Run(fmt.Sprintf(":%d", utils.Cfg.MustInt("http", "port")))
}
