# task-serve

一个简单的任务服务，包括API服务和算法服务。

## 技术栈

- Golang 1.18
- VsCode
- Docker 
- Redis
- RabbitMQ

## Docker环境配置

- Docker 使用 Docker Desktop
- Redis 安装教程：[链接](https://hub.docker.com/_/redis/)
- RabbitMQ 安装教程：[链接](https://juejin.cn/post/7198430801850105916)

## 环境启动

- 创建Volume

```
docker volume create redis_volume
docker volume create rabbitmq_volume
```

- 启动Redis

```
docker run --name redis_serve -p 6379:6379 -v redis_volume:/data -d redis redis-server --save 60 1 --loglevel warning
```

- 启动RabbitMQ

```
docker run --name rabbitmq_serve -p 15672:15672 -p 5672:5672 -v rabbitmq_volume:/var/lib/rabbitmq -e RABBITMQ_DEFAULT_USER=ciaran -e RABBITMQ_DEFAULT_PASS=123456 -d rabbitmq:management
```

## 代码测试

- 均存放在*_use.go文件中,使用vscode 编辑器的run test通过

## 运行

- 运行算法服务 `go run .\main.go -s algo`
- 创建任务 `go run .\main.go -s api -op create -tname test111 -tcommand none`
- 查询状态 `go run .\main.go -s api -op select -tid task-24`

## 思考

- 当大并发请求时，如何提升服务可用性。我认为
    - 通过建立连接池来保持mq和redis链接，避免频繁链接的时间浪费
    - 通过协程和锁加快api请求的处理速度
    - 通过建立分布式系统来处理大量请求
- 耗时任务如何进行进度上报，取消任务等状态管理。我认为
    - 进度上报可以使用心跳机制，传输任务状态
    - 取消任务也可以通过心跳机制，传输任务操作

## 需要改进的点

- 上述思考的需求
- 这个简单项目，测试交不完备，可能存在许多bug，需要修复
- 项目结构可能需要改进
- 没有协程取消的功能，协程会一直运行到程序结束，当想结束算法服务时，必须结束程序
- 程序并未执行task中的相应命令
- 只是简单使用redis和rabbitmq，并没有完备地构建消费者和生产者。