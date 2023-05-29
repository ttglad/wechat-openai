package config

import (
	"flag"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

var (
	cfgFile = flag.String("config", "./config.yaml", "配置文件路径")

	cfg *Config
)

//Config example config
type Config struct {
	Listen                 string `yaml:"listen"`
	Proxy                  string `yaml:"proxy"`
	*OfficialAccountConfig `yaml:"officialAccountConfig"`
	*OpenAI                `yaml:"openai"`
}

//OfficialAccountConfig 公众号相关配置
type OfficialAccountConfig struct {
	AppID          string `yaml:"appID"`
	AppSecret      string `yaml:"appSecret"`
	Token          string `yaml:"token"`
	EncodingAESKey string `yaml:"encodingAESKey"`
	Timeout        int    `yaml:"timeout"`
}

type OpenAI struct {
	Key string `yaml:"key"`

	Params struct {
		Api         string  `yaml:"api"`
		Model       string  `yaml:"model"`
		Prompt      string  `yaml:"prompt"`
		Temperature float32 `yaml:"temperature"`
		MaxTokens   uint16  `yaml:"maxTokens"`
	} `yaml:"params"`

	MaxQuestionLength int `yaml:"maxQuestionLength"`
}

//GetConfig 获取配置
func GetConfig() *Config {
	if cfg != nil {
		return cfg
	}
	bytes, err := ioutil.ReadFile(*cfgFile)
	if err != nil {
		panic(err)
	}

	cfgData := &Config{}
	err = yaml.Unmarshal(bytes, cfgData)
	if err != nil {
		panic(err)
	}
	cfg = cfgData
	return cfg
}
