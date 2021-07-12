package write

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

const (
	LogLevelError = "ERROR"
	LogLevelINFO  = "INFO"
	LogLevelWarn  = "WARN"
	GithubLogName = "github-webhooks.log"
	GitLabLogName = "gitlab-webhooks.log"
)

//写日志文件
func Log(level string, content string, fileName string) {
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
