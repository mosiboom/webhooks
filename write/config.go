package write

import (
	"encoding/json"
	"os"
)

type Config struct {
	ListenerPort   string            `json:"listener_port"`
	ListenerRoute  string            `json:"listener_route"`
	WebhooksSecret string            `json:"webhooks_secret"`
	Command        map[string]string `json:"command"`
}

//读取配置文件
func GetConfig(path string) (*Config, error) {
	config := Config{}
	filePtr, err := os.Open(path)
	if err != nil {
		return &config, err
	}
	defer func() {
		_ = filePtr.Close()
	}()
	// 创建json解码器
	decoder := json.NewDecoder(filePtr)
	err = decoder.Decode(&config)
	if err != nil {
		return &config, err
	}
	return &config, nil
}
