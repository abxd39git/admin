package http

import (
	cf "admin/gateway/conf"
	"admin/gateway/http/controller"
	"fmt"

	"github.com/gin-gonic/gin"
)

func initRouter() *gin.Engine {
	r := gin.Default()
	//gin.SetMode(gin.ReleaseMode)

	new(controller.ContentManage).Router(r)
	// new(controller.WalletGroup).Router(r)
	// new(controller.ArticleGroup).Router(r)
	// new(controller.ActionGroup).Router(r)
	// new(controller.CurrencyGroup).Router(r)
	// new(controller.KineGroup).Router(r)
	// new(controller.TokenGroup).Router(r)
	return r
}

func InitHttpServer() {
	port := cf.Cfg.MustInt("http", "port")
	r := initRouter()
	r.Run(fmt.Sprintf(":%d", port))
}
