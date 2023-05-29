package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/officialaccount"
	wechatConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/silenceper/wechat/v2/officialaccount/message"
	log "github.com/sirupsen/logrus"
	"net/http"
	"sync"
	"time"
	"wechat-openai/config"
)

var (
	requests sync.Map // K - 消息ID ， V - chan string
)

type WechatAccount struct {
	wc              *wechat.Wechat
	officialAccount *officialaccount.OfficialAccount
}

func NewWechatAccount(wc *wechat.Wechat) *WechatAccount {
	//init config
	globalCfg := config.GetConfig()
	offCfg := &wechatConfig.Config{
		AppID:          globalCfg.AppID,
		AppSecret:      globalCfg.AppSecret,
		Token:          globalCfg.Token,
		EncodingAESKey: globalCfg.EncodingAESKey,
	}
	log.Debugf("offCfg=%+v", offCfg)
	officialAccount := wc.GetOfficialAccount(offCfg)
	return &WechatAccount{
		wc:              wc,
		officialAccount: officialAccount,
	}
}

//Serve 处理消息
func (ex *WechatAccount) Serve(c *gin.Context) {
	// 传入request和responseWriter
	server := ex.officialAccount.GetServer(c.Request, c.Writer)
	server.SkipValidate(true)
	//设置接收消息的处理方法
	server.SetMessageHandler(func(msg *message.MixMessage) *message.Reply {
		if msg.MsgType == message.MsgTypeText {
			text := message.NewText(msg.Content)
			//if msg.Content == "陶希禾" {
			//	text := message.NewText(msg.Content)
			//	return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
			//}

			var ch chan string
			v, ok := requests.Load(msg.MsgID)
			if !ok {
				ch = make(chan string)
				requests.Store(msg.MsgID, ch)
				ch <- OpenAiQuery(msg.OpenID, msg.Content, time.Second*time.Duration(config.GetConfig().Timeout))
			} else {
				ch = v.(chan string)
			}

			select {
			case result := <-ch:
				//bs := msg.GenerateEchoData(result)
				text = message.NewText(result)
				requests.Delete(msg.MsgID)
			// 超时不要回答，会重试的
			case <-time.After(time.Second * 5):
			}
			return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
		}
		//switch msg.MsgType {
		//case message.MsgTypeText:
		//	text := message.NewText(msg.Content)
		//	return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
		//}
		//回复消息：演示回复用户发送的消息
		//text := message.NewText(msg.Content)
		//return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}

		//article1 := message.NewArticle("测试图文1", "图文描述", "", "")
		//articles := []*message.Article{article1}
		//news := message.NewNews(articles)
		//return &message.Reply{MsgType: message.MsgTypeNews, MsgData: news}

		//voice := message.NewVoice(mediaID)
		//return &message.Reply{MsgType: message.MsgTypeVoice, MsgData: voice}

		//
		//video := message.NewVideo(mediaID, "标题", "描述")
		//return &message.Reply{MsgType: message.MsgTypeVideo, MsgData: video}

		//music := message.NewMusic("标题", "描述", "音乐链接", "HQMusicUrl", "缩略图的媒体id")
		//return &message.Reply{MsgType: message.MsgTypeMusic, MsgData: music}

		//多客服消息转发
		//transferCustomer := message.NewTransferCustomer("")
		//return &message.Reply{MsgType: message.MsgTypeTransfer, MsgData: transferCustomer}
		return nil
	})

	//处理消息接收以及回复
	err := server.Serve()
	if err != nil {
		log.Error("Serve Error, err=%+v", err)
		return
	}
	//发送回复的消息
	err = server.Send()
	if err != nil {
		log.Error("Send Error, err=%+v", err)
		return
	}
}

//Serve 处理消息
func (ex *WechatAccount) Test(c *gin.Context) {
	fmt.Println("request start")
	msg, _ := c.GetQuery("msg")
	replyMsg := OpenAiQuery("0", msg, time.Second*5)
	fmt.Println("request end")
	c.JSON(http.StatusOK, gin.H{
		"reply": replyMsg,
	})
}
