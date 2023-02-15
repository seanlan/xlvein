# XLVein

支持分布式的websocket第三方服务，结构简单，依赖少，好部署。提供api接口供应用服务调用，实现主动给客户端推送消息。

单机部署无任何依赖，可直接运行

分布式部署可利用nginx或其他负载均衡器实现负载均衡，利用redis或rabbitmq实现消息交换。只需要修改配置文件即可

## 依赖
- redis【可选】利用pub sub实现多机部署时，各个服务消息交换
- rabbitmq【可选】实现多机部署时，各个服务消息交换

## 使用

启动服务
```bash
go run . start
```
这样服务就启动了，默认监听8090端口

获取websocket连接地址
```bash
go run . test
>>> ws://127.0.0.1:8090/ws/connect?app_id=test&token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NzkwMjQ2ODYsImlzcyI6IlRpbWVUb2tlbiIsInN1YiI6InVzZXJfMSJ9.qANpLFH-Z8z2inP_8NT6FAysR1EWvwPWDfFlNWf0RbY
```
打开 http://www.websocket-test.com/ 输入上面的地址，点击connect，就可以看到连接成功的消息
再次运行
```bash
go run . test
```
通过http调用接口推送消息

可以看到上一个浏览器窗口收到了websocket消息

### 生成websocket连接地址规则
```
ws://HOSTNAME/ws/connect?app_id=test&token=TOKEN
```
Token 采用JTW生成，生成规则：
```javascript
var payload={sub:'account'} // account 为用户账号，由应用方自行定义
var appSecret='1jaos8dfasdfjf9sadf'
bindKey = jwt.sign(payload,appSecret,{expiresIn:3600*24*30})
```
token错误会被拒绝链接


### 调用Http API 推送消息
- 接口地址：/api/v1/im/push
- 参数：
```
{
    "app_key": APPKEY,  // 应用app_key
    "send_to": ACCOUNT, // 接收方账号
    "message": MESSAGE, // json格式 {} ,内容自定义
	"nonce":   NONCE,   // 随机字符串
}
```

## 配置
配置文件在conf.yml中，配置项如下：
```yaml
host:   :8080 #服务监听地址
# 交换机配置
exchange:
  type: "local" # 消息交换类型 目前支持 rabbitmq,redis,local local方式不支持分布式部署
  exchange_name: "xlvein.im"
  queue_name: "xlvein.im.queue"
  rabbitmq: "amqp://xlvein:WFxg5FNMedN@T$z@127.0.0.1:5672/vein" # rabbitmq连接地址[type=rabbitmq 时需要]
  redis: "127.0.0.1:6379" # redis连接地址[type=redis 时需要]
# 应用配置
applications: # 应用列表，可自行添加多个应用
  - {app_key: "test",  app_secret: "j8jasd98efan9sdfj89asjdf"}
  - {app_key: "test2", app_secret: "j8jasd98efan9sdfj89asjdf"}
```
