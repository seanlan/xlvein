# XLVein

支持分布式的websocket第三方服务。提供api接口供应用服务调用，实现主动给客户端推送消息。分布式部署可利用nginx或其他负载均衡器实现负载均衡。

## 依赖
- redis【可选】利用pub sub实现多机部署时，各个服务消息交换
- rabbitmq【可选】实现多机部署时，各个服务消息交换

## 配置
配置文件在conf.yml中，配置项如下：
```yaml
host:   :8080 #服务监听地址
# 交换机配置
exchange:
  type: "local" # 消息交换类型 目前支持 rabbitmq,redis,local local方式不支持分布式部署
  rabbitmq: "amqp://xlvein:WFxg5FNMedN@T$z@127.0.0.1:5672/vein"
  exchange_name: "xlvein.im" # rabbitmq交换机名称
  queue_name: "xlvein.im.queue" # 队列名称
  redis: "127.0.0.1:6379" # redis地址，redis作为消息交换试需要
```

## 使用

### 生成客户端websocket连接
```
ws://127.0.0.1:8090/ws/connect?app_id=test&token=TOKEN
```
Token 采用JTW生成，生成规则：
```
var payload={sub:'account'} // account 为用户账号，由应用方自行定义
var appSecret='1jaos8dfasdfjf9sadf'
bindKey = jwt.sign(payload,appSecret,{expiresIn:3600*24*30})
```
token错误会被拒绝链接

### 调用Http Api 推送消息
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

