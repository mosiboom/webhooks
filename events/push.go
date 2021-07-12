package events

import (
	"bytes"
	"encoding/json"
	"github.com/go-playground/webhooks/v6/github"
	"github.com/go-playground/webhooks/v6/gitlab"
	"log"
	"os/exec"
	"webhooks/write"
)

type LogContent struct {
	RepositoryURL string `json:"repository_url"`
	Command       string `json:"command"`
}
type Push struct {
}

func (push Push) Github(formPayload github.PushPayload, config write.Config) {
	url := formPayload.Repository.URL
	log.Printf("仓库URL：%+v\n", url)
	commandStr := config.Command[url]
	log.Printf("项目命令：%v\n", commandStr)
	LC := LogContent{
		RepositoryURL: url,
		Command:       commandStr,
	}
	logContents, _ := json.Marshal(LC)
	write.Log(write.LogLevelINFO, "[命令内容] -- "+string(logContents), write.GithubLogName)
	command := exec.Command("/bin/sh", "-c", commandStr)
	//给标准输入以及标准错误初始化一个buffer，每条命令的输出位置可能是不一样的，
	//比如有的命令会将输出放到stdout，有的放到stderr
	command.Stdout = &bytes.Buffer{}
	command.Stderr = &bytes.Buffer{}
	//执行命令，直到命令结束
	err := command.Run()
	if err != nil {
		//打印程序中的错误以及命令行标准错误中的输出
		write.Log(write.LogLevelError, "[执行错误] -- "+err.Error(), write.GithubLogName)
		write.Log(write.LogLevelError, "[返回内容] -- "+command.Stderr.(*bytes.Buffer).String(), write.GithubLogName)
		log.Println(err)
		log.Println(command.Stderr.(*bytes.Buffer).String())
		return
	}
	//打印命令行的标准输出
	write.Log(write.LogLevelINFO, "[返回内容] -- "+command.Stderr.(*bytes.Buffer).String(), write.GithubLogName)
	log.Println(command.Stdout.(*bytes.Buffer).String())
}

func (push Push) Gitlab(formPayload gitlab.PushEventPayload, config write.Config) {
	url := formPayload.Repository.URL
	log.Printf("仓库URL：%+v\n", url)
	commandStr := config.Command[url]
	log.Printf("项目命令：%v\n", commandStr)
	LC := LogContent{
		RepositoryURL: url,
		Command:       commandStr,
	}
	logContents, _ := json.Marshal(LC)
	write.Log(write.LogLevelINFO, "[命令内容] -- "+string(logContents), write.GitLabLogName)
	command := exec.Command("/bin/sh", "-c", commandStr)
	//给标准输入以及标准错误初始化一个buffer，每条命令的输出位置可能是不一样的，
	//比如有的命令会将输出放到stdout，有的放到stderr
	command.Stdout = &bytes.Buffer{}
	command.Stderr = &bytes.Buffer{}
	//执行命令，直到命令结束
	err := command.Run()
	if err != nil {
		//打印程序中的错误以及命令行标准错误中的输出
		write.Log(write.LogLevelError, "[执行错误] -- "+err.Error(), write.GitLabLogName)
		write.Log(write.LogLevelError, "[返回内容] -- "+command.Stderr.(*bytes.Buffer).String(), write.GitLabLogName)
		log.Println(err)
		log.Println(command.Stderr.(*bytes.Buffer).String())
		return
	}
	//打印命令行的标准输出
	write.Log(write.LogLevelINFO, "[返回内容] -- "+command.Stderr.(*bytes.Buffer).String(), write.GitLabLogName)
	log.Println(command.Stdout.(*bytes.Buffer).String())
}
