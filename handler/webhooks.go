package handler

import (
	"github.com/go-playground/webhooks/v6/github"
	"github.com/go-playground/webhooks/v6/gitlab"
	"log"
	"net/http"
	"webhooks/events"
	"webhooks/write"
)

const (
	GitHubConfigPath = "github-config.json"
	GitLabConfigPath = "gitlab-config.json"
)

//GitHubHanlder:github相关的处理
func GitHubHanlder(w http.ResponseWriter, r *http.Request) {
	//获取配置
	config, err := write.GetConfig(GitHubConfigPath)
	if err != nil {
		write.Log(write.LogLevelError, "[读取配置] -- "+err.Error(), write.GithubLogName)
		log.Println(err.Error())
		return
	}
	//new一个github对象并解析payload
	hook, _ := github.New(github.Options.Secret(config.WebhooksSecret))
	payload, hookErr := hook.Parse(r, github.PushEvent, github.ReleaseEvent, github.PullRequestEvent)
	if hookErr != nil {
		log.Println(hookErr)
		write.Log(write.LogLevelError, "[解析错误] -- "+hookErr.Error(), write.GithubLogName)
		return
	}
	switch payload.(type) {
	case github.PushPayload:
		push := payload.(github.PushPayload)
		events.Push{}.Github(push, *config)
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

func GitLabHanlder(w http.ResponseWriter, r *http.Request) {
	//获取配置
	config, err := write.GetConfig(GitLabConfigPath)
	if err != nil {
		write.Log(write.LogLevelError, "[读取配置] -- "+err.Error(), write.GitLabLogName)
		log.Println(err.Error())
		return
	}
	//new一个github对象并解析payload
	hook, _ := gitlab.New(gitlab.Options.Secret(config.WebhooksSecret))
	payload, hookErr := hook.Parse(r, gitlab.PushEvents)
	if hookErr != nil {
		log.Println(hookErr)
		write.Log(write.LogLevelError, "[解析错误] -- "+hookErr.Error(), write.GitLabLogName)
		return
	}
	switch payload.(type) {
	case gitlab.PushEventPayload:
		push := payload.(gitlab.PushEventPayload)
		events.Push{}.Gitlab(push, *config)
	}
}
