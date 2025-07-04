# Go Template

![ProxyHub](./logo.png)

基于 [Gin](https://github.com/gin-gonic/gin) 框架以及 [GORM](https://gorm.io/) 构建的后端 API 服务。

开发前，使用 Goland 2021 及以上版本打开 `go.mod`，选中 module 名，右键选择 Refactor 重构为新项目的 module 名。并修改
`internal/app/run.go`
中的常量 `projectName` 为新项目名。

项目结构规范依据 [Standard Go Project Layout](https://github.com/golang-standards/project-layout) 进行约束。

## 主要功能

- 项目整体使用 [golangci-lint](https://github.com/golangci/golangci-lint) 进行较为严格的静态代码检查，其配置文件为
  `.golangci-lint.yml`
- `Gin` 实现的 WEB 服务，路由基于 Action 模式
- 引入了跨域，记录请求和响应，pprof 以及 prometheus metrics 的中间件
- API
  参数校验，依赖 [validator v10](https://github.com/go-playground/validator)
  ，详细用法可查阅[文档](https://pkg.go.dev/github.com/go-playground/validator/v10)及自定义翻译
- 基于 `GORM` 的 `MySQL` 读写分离，默认为单库
- 基于 [jaeger](https://www.jaegertracing.io/) 实现链路追踪
- 基于 [zap](https://github.com/uber-go/zap) 以及 [lumberjack](https://github.com/natefinch/lumberjack/tree/v2.0)
  的日志记录及切割，支持通过 HTTP 请求动态调整日志级别
- 基于 [sarama](https://github.com/Shopify/sarama) 的 `Kafka` 生产者和消费者，可按需使用或删除。Kafka 集群最低支持
  `0.8.2.0` 版本。
- swagger2.0 文档
- 演示的单元测试用例。接口具体实现的函数单元测试，可在 `internal/logic/demo/demo_test.go` 查看

## 注意事项

### 关于配置文件

目录 `configs`，支持使用环境变量覆盖配置文件中的配置项。即，优先从环境变量中获取对应的配置，当环境变量没有时，使用配置文件中的默认配置。

若要查看所有支持的环境变量，可执行 `go run cmd/app/main.go env` 命令

若没有配置如下配置项

- `service_name`，则默认值为 gin-project
- `local_ip`，则默认值为 0.0.0.0
- `api_port`，则默认值为 8000

`run_mode` 为 debug，`Gin` 会运行在 debug 模式，在生产环境中，更换为任意值，则可以运行在 release 模式

模版根据 `configs/dev.yaml` 中的配置，默认启动 jaeger tracer，MySQL client，Kafka client，redis client
以及配置中心服务地址轮训。可根据项目需要删除对应配置屏蔽对应客户端。

项目若需要添加自己的配置文件，可修改 `internal/config/config.go`，添加需要的配置结构体，并在配置文件中默认的 `base `
节点同级新增配置内容。

### 关于 Kafka 消费/生产

考虑到不是每一个项目都需要 Kafka，所以相关初始化代码默认是被注释的。如需启用，可取消 `internal/app/run.go` 中的注释。如果不需要
Kafka，可以删除配置文件中的 Kafka
配置，并删除 `internal/app/run.go` 中注释的代码，以及 `internal/producer` 和 `internal/consumer`

### 关于单元测试建议

所有 db 相关操作，建议用 interface 封装，具体的 interface 与 `models` 下表文件名关联，在 `models/repo` 下新建同名文件，然后在
`models/factory` 下新建同名文件，实现这个
interface，新建一个以 `test_` 为前缀的文件，实现 interface 用于单元测试。

关于接口的单元测试，已经在 `internal/handler/demo_test.go` 中编写了一个示例，可以直接参考。

### 关于编译

`vendor` 目录默认从项目中忽略，为加速 CI 中的编译速度，可以考虑将 vendor 添加到实际项目中

### 关于 RSQL

 操作符         | 操作符      
-------------|----------|
 ==          | 等于       
 !=          | 	不等于     
 =lt= 或 <	   | 小于       
 =lte= 或 <=	 | 小于或等于    
 =gt= 或 >	   | 大于       
 =gte= 或 >=	 | 大于或等于    
 =in=	       | In       
 =out=	      | Not In   
 ;           | AND，和表达式 
 ,           | OR，或表达式  

```bash
# 用户名max，且年龄18
UserName==max;Age=18

# 用户名在max和terry，且年龄18或者性别为男，且住址上开头的
UserName=in=(max,terry);(Age==18,gender=="男");Addr=="上.*"
```

#### 关于 Swagger 文档

```bash
swag init -g ./cmd/app/main.go
```

文档地址： `ip:port/api-docs`
