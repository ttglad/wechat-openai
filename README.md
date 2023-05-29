# 微信公众号接入GPT自动回复消息

> 项目地址: https://github.com/ttglad/wechat-openai
# 简介
最近GPT火热，所以有了这个想法，利用公众号接入GPT自动回复用户的消息。

# 开始部署

## 一、 环境准备

- 一台 Linux 服务器，建议 国外服务器，或者任何可以长期运行程序的PC设备
- OpenAI 账号 以及生成的 `SECRET KEY` ，本文对账号注册以及 key 生成不做赘述，读者请自行搜索解决方案。
- 一个微信公众号，熟悉微信公众号后台开发配置，项目运行成功后需要在后台配置项目请求的url。

> - 注：OpenAI 的域名 `https://api.openai.com` 在国内由于某种原因可能无法访问，读者需要自己解决 API 访问不通的问题。介绍一种简单的国内代理搭建方式

## 二、 配置

进行配置：

把 `config.yaml.example` 重命名成 `config.yaml`，然后利用文本编辑器修改此文件：

```php
# 监听本地的端口
listen: :8899
# 可选: 代理地址。需要你有本地或远程代理软件，举例: socks5://127.0.0.1:1080
proxy:
# 微信公众号相关配置
officialAccountConfig:
  # 微信公众号后台获取
  appID:
  # 微信公众号后台获取
  appSecret:
  # 微信公众号后台获取
  token:
  # 微信公众号后台获取
  encodingAESKey:
  # 影响滚动返回结果 (5s-13s) 
  timeout: 7

openai:
  # 必填: KEY。 文档: https://platform.openai.com/account/api-keys
  key:
  # 可选: 参数调节
  params:
    # openai的接口地址，放出来是因为有些人做了反向代理，要注意这有安全问题，谨慎使用
    api: https://api.openai.com/v1/chat/completions
    # 暂时请使用 gpt-3.5-turbo
    model: gpt-3.5-turbo
    # 提示。 可以理解为对其身份设定。 文档: https://platform.openai.com/docs/guides/chat/introduction
    # 每个问题都会携带，注意，它也占用token消耗。
    prompt:
    # 影响 问题+回复的长度。  gpt-3.5模型最大4096， 非1个汉字1token
    maxTokens: 1024
    # 温度。 0-2 。较高的值将使输出更加随机，而较低的值将使输出更加集中和确定。
    temperature: 0.8
  # 限制用户问题最大长度。这个以字计算，非token.
  maxQuestionLength: 200
```

## 三、 运行

```bash
go run main.go
```

# 联系作者

- 本项目地址为： https://github.com/ttglad/wechat-openai ，欢迎大家 Star，提交 PR
- 有问题可以在本项目下提 `Issues` 或者发邮件到 `tonneylon@gmail.com`