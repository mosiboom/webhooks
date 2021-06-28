## Golang实现Webhooks服务端的简单自动化部署

- **项目说明：**
  本项目基于Golang，实现简单的自动化部署（目前只实现了github）。可以解决在需要的场景下（例如更新到测试机，自己部署等场景），不用上去服务器也可以自动将更新的项目拉取下来。原理就是使用github仓库的webhook钩子去监听仓库的`push`
  事件，当你将更新的内容推送到远程仓库时，github会向你的接口发起一个请求，将修改的内容和监听的事件发送到该接口
- **注意事项：**
  如果只是想单独用这个服务，可以直接执行我已经编译后的二进制可执行文件（只针对Linux和mac OS x），在`build`文件夹中，并修改配置文件中的参数，然后执行命令`nohub ./github-webhooks &`即可
- **文档参考：**
    - [github官方钩子文档](https://docs.github.com/en/developers/webhooks-and-events/webhooks/about-webhooks)
    - [go-playground/webhooks/v6接口文档](https://pkg.go.dev/github.com/go-playground/webhooks/v6@v6.0.0-beta.3/github#pkg-variables)
    - [go-playground包github仓库](https://github.com/go-playground/webhooks)

1. 配置文件：`webhook-config.json` 详解
    ```json5
    {
        "listener_route": "/webhooks", //服务器路由地址
        "listener_port": "3000",       //服务器监听端口
        "webhooks_secret": "you-secret-key",  //与配置在github webhooks的secret一致
        "command": {  //配置命令 格式：仓库地址:执行命令（Push之后要做些什么）
        "https://github.com/mosiboom/TestRepository": "cd /code/test/TestRepository/ && /usr/bin/git pull"
        }
    }
    ```

2. 修改`github-webhooks.go`中的`WebhooksSecret`的值，该值为加密密钥（需要在仓库配置和创建webhooks时填写的一致）

3. 安装依赖：`go get -u`

4. 运行命令：
   ```go run github-webhooks.go```

5. 构建二进制命令：
   ```go build github-webhooks.go```

6. 可以使用`nohup`后台运行服务器（当然也可以使用其他方式）：
    ```shell
    nohub go run github-webhooks.go &
    #基于二进制构建完毕后的可执行文件
    nohub ./github-webhooks &
    ```

- 声明：该项目是本人学习golang的一个简单练手项目，代码很多的不合理或者写的不对的地方请多多包涵，也希望有人帮我指出我的问题，指导我进步