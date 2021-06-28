package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-playground/webhooks/v6/github"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	ConfigPath    = "webhooks-config.json"
	LogLevelError = "ERROR"
	LogLevelINFO  = "INFO"
	LogLevelWarn  = "WARN"
	LogFileName   = "github-webhooks.log"
)

type LogContent struct {
	RepositoryURL string `json:"repository_url"`
	Command       string `json:"command"`
}

type Config struct {
	ListenerPort   string            `json:"listener_port"`
	ListenerRoute  string            `json:"listener_route"`
	WebhooksSecret string            `json:"webhooks_secret"`
	Command        map[string]string `json:"command"`
}

func main() {
	//获取配置
	config, err := GetConfig(ConfigPath)
	if err != nil {
		WriteLog(LogLevelError, "[读取配置] -- "+err.Error(), LogFileName)
		log.Println(err.Error())
		return
	}
	http.HandleFunc(config.ListenerRoute, webHooks)
	err = http.ListenAndServe(":"+config.ListenerPort, nil)
	if err != nil {
		WriteLog(LogLevelError, "[监听失败] -- "+err.Error(), LogFileName)
		log.Fatal("ListenAndServe：", err)
	}
}

func webHooks(w http.ResponseWriter, r *http.Request) {
	//获取配置
	config, err := GetConfig(ConfigPath)
	if err != nil {
		WriteLog(LogLevelError, "[读取配置] -- "+err.Error(), LogFileName)
		log.Println(err.Error())
		return
	}
	//new一个github对象并解析payload
	hook, _ := github.New(github.Options.Secret(config.WebhooksSecret))
	payload, hookErr := hook.Parse(r, github.PushEvent, github.ReleaseEvent, github.PullRequestEvent)
	if hookErr != nil {
		log.Println(hookErr)
		WriteLog(LogLevelError, "[解析错误] -- "+hookErr.Error(), LogFileName)
		return
	}
	switch payload.(type) {
	case github.PushPayload:
		push := payload.(github.PushPayload)
		DealPushEvent(push, *config)
	case github.ReleasePayload:
		release := payload.(github.ReleasePayload)
		log.Println("Release Event:")
		log.Printf("%+v \n", release)
	case github.PullRequestPayload:
		pullRequest := payload.(github.PullRequestPayload)
		log.Println("Pull Request Event:")
		log.Printf("%+v \n", pullRequest)
	}

}

//处理Push事件
func DealPushEvent(formPayload github.PushPayload, config Config) {
	url := formPayload.Repository.URL
	log.Printf("仓库URL：%+v\n", url)
	commandStr := config.Command[url]
	log.Printf("项目命令：%v\n", commandStr)
	LC := LogContent{
		RepositoryURL: url,
		Command:       commandStr,
	}
	logContents, _ := json.Marshal(LC)
	WriteLog(LogLevelINFO, "[命令内容] -- "+string(logContents), LogFileName)
	command := exec.Command("/bin/sh", "-c", commandStr)
	//给标准输入以及标准错误初始化一个buffer，每条命令的输出位置可能是不一样的，
	//比如有的命令会将输出放到stdout，有的放到stderr
	command.Stdout = &bytes.Buffer{}
	command.Stderr = &bytes.Buffer{}
	//执行命令，直到命令结束
	err := command.Run()
	if err != nil {
		//打印程序中的错误以及命令行标准错误中的输出
		WriteLog(LogLevelError, "[执行错误] -- "+err.Error(), LogFileName)
		WriteLog(LogLevelError, "[返回内容] -- "+command.Stderr.(*bytes.Buffer).String(), LogFileName)
		log.Println(err)
		log.Println(command.Stderr.(*bytes.Buffer).String())
		return
	}
	//打印命令行的标准输出
	WriteLog(LogLevelINFO, "[返回内容] -- "+command.Stderr.(*bytes.Buffer).String(), LogFileName)
	log.Println(command.Stdout.(*bytes.Buffer).String())
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

//写日志文件
func WriteLog(level string, content string, fileName string) {
	dateObj := time.Unix(time.Now().Unix(), 0)
	logDate := dateObj.Format("2006-04-02 15:04:05")
	dir, _ := os.Getwd()
	logPath := strings.Replace(dir, "\\", "/", -1) + "/logs/" + dateObj.Format("2006-01-02")
	_, err := PathExistsAndCreate(logPath)
	if err != nil {
		log.Printf("PathExistsAndCreate Error: [%v]\n", err)
		return
	}
	fullName := logPath + "/" + fileName
	logString := fmt.Sprintf("[%v] %v：%v\n", logDate, level, content)
	f, err := os.OpenFile(fullName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	defer func() {
		_ = f.Close()
	}()
	if err != nil {
		log.Printf("open file error![%v]\n", err)
		return
	}
	_, err = f.Write([]byte(logString))
	if err != nil {
		log.Printf("write file error![%v]\n", err)
		return
	}
}

//判断路径存不存在，不存在则创建
func PathExistsAndCreate(path string) (bool, error) {
	_, err := os.Stat(path)
	if !os.IsNotExist(err) {
		return true, nil
	}
	// 创建文件夹
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return false, err
	}
	return true, nil
}
