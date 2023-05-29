package main

import (
	"github.com/gin-gonic/gin"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"wechat-openai/config"
	"wechat-openai/server"
)

func main() {
	//获取wechat实例
	wc := InitWechat()

	cfg := config.GetConfig()

	wechatAccount := server.NewWechatAccount(wc)
	r := gin.Default()
	r.Any("/api/v1/wechat", wechatAccount.Serve)
	//r.GET("/api/v1/wechat/test", wechatAccount.Test)

	r.Run(cfg.Listen)
}

func InitWechat() *wechat.Wechat {
	wc := wechat.NewWechat()
	wc.SetCache(cache.NewMemory())
	return wc
}
