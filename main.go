package main

import (
	"log"
	"net/http"
	"webhooks/handler"
	"webhooks/write"
)

const (
	ConfigPath = "webhooks-config.json"
)

func main() {
	http.HandleFunc("/github", handler.GitHubHanlder)
	http.HandleFunc("/gitlab", handler.GitLabHanlder)
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		write.Log(write.LogLevelError, "[监听失败] -- "+err.Error(), write.GithubLogName)
		log.Fatal("ListenAndServe：", err)
	}
}
