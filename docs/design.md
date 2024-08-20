# design

## 协议

### header

#### 定义

header中的key不区分大小写

#### 服务名

ServiceName

#### 函数名

FuncName

#### 序列化方式

Codec

#### 错误

Error

#### metadata

#### in

接受外界的metadata

#### out

发送给外界的matadata

### body

#### req

#### rsp

## 概念

### provider

#### 功能

1. 服务提供方的入口
2. 支持注册结构体
3. 支持注册函数
4. 一个provider可以有多个实例
5. 支持多序列化方式

### consumer

####  功能

1. 服务消费方入口
2. 关联注册中心
3. 

#### 拦截器链

1. 能够在链条上共享数据,不污染用户参数
2. 链条分为系统链和用户链
3. 系统链为框架能够运行的基础
4. 用户链是用户定义在系统链前的拦截器
5. 提供用户替换系统链条的入口
6. 系统链条必须被统一替换,此操作为原子操作
7. 服务过滤,服务负载,服务调用定义为系统拦截器
8. 自己增加功能也是通过

### registry

#### 功能

1. 服务启动后注册
2. 服务关闭后取消注册
3. 服务治理,预留扩展字段
4. 服务健康检查,为每个服务都默认提供一个接口,用于服务检查
5. 服务发现,增加内存缓存/本地缓存/文件配置
6. 消费者注册
7. 服务id最好不变
8. 服务权重
9. 服务需要提供过期时间
10. 服务能够定时注册

#### 服务提供者注册

```
/dodo/serviceName/providers/urlencode(tcp://172.18.1.131:8888/serviceName?version=1.0.0&side=provider&funcs=SayHello,SayWorld&tls=true)
```

#### 服务消费者注册

```
/dodo/serviceName/consumers/urlencode(tcp://172.18.1.131:8888/serviceName?version=1.0.0&side=provider&funcs=SayHello,SayWorld&tls=true)
```

### server

#### 功能

1. 调用幂等
2. 协议解析
3. 支持注册多服务
4. 支持自动获取host
5. 端口不指定默认为17312
6. 支持多种序列化方式
7. 自动生成证书

#### 实现

1. 定义统一接口
2. 实现rpc server

#### 计划

1. 实现rest server对外直接提供http接口

### transport

#### 功能

1. 抽象传输层,对上层提供统一接口
2. transport dail/listen支持tls

#### 实现

1. 实现grpc传输层
2. 实现tcp传输层

### metadata

#### 功能

1. 将内容放入到context中

### log

#### 功能

1. 定义统一日志接口
2. 提供默认实现
3. 支持用户使用自己的日志包

### invoker

#### 功能

1. 支持拦截
2. 封装操作单元

### codec

#### 实现

1. 定义统一接口
2. 实现json序列化

### selector

#### 功能

1. 服务过滤
2. 服务负载
3. 服务缓存

#### 实现

##### 缓存目录

```
# 默认
./.dodo/consumer/selector
```

##### 配置目录

```
# 默认
../config/selector
```

##### 配置优先级

1. 配置目录
2. 注册中心
3. 缓存

### client

#### 实现

1. 定义统一接口
2. 实现rpc client
3. 支持tls

### 功能

#### 自动获取IP地址

1. 可以手动指定服务地址(不推荐)
2. 获取eth0地址
3. 获取非回环地址和非docker地址

client的连接池看是否可以参考http的实现,rpc_pool

```
var DefaultTransport RoundTripper = &Transport{
	Proxy: ProxyFromEnvironment,
	DialContext: (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: true,
	}).DialContext,
	MaxIdleConns:          100,
	IdleConnTimeout:       90 * time.Second,
	TLSHandshakeTimeout:   10 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
}
```

这边一定要能够写一个好的测试程序,尤其是对于模块而言,越不测试就越容易出问题



## 参考

1. [go-micro](https://micro.mu/blog/2016/05/15/resiliency.html)
2. [dubbo](http://dubbo.apache.org/zh-cn/)
3. [rpcx](http://doc.rpcx.site/)
4. [go-chassis](https://go-chassis.readthedocs.io/en/latest/index.html)
5. [ratelimit](github.com/juju/ratelimit)
6. [uber ratelimit](go.uber.org/ratelimit)
7. [metrics](github.com/rcrowley/go-metrics)
8. [breaker](github.com/sony/gobreaker)
9. [fault tolerance](github.com/afex/hystrix-go)
10. [backoff](github.com/cenkalti/backoff)
11. [distribute key/value](github.com/docker/libkv)
12. [Microservices: What's Missing](https://www.slideshare.net/adriancockcroft/microservices-whats-missing-oreilly-software-architecture-new-york#24)
13. [聊聊微服务的服务注册与发现](http://jm.taobao.org/2018/06/26/%E8%81%8A%E8%81%8A%E5%BE%AE%E6%9C%8D%E5%8A%A1%E7%9A%84%E6%9C%8D%E5%8A%A1%E6%B3%A8%E5%86%8C%E4%B8%8E%E5%8F%91%E7%8E%B0/)
14. [How to get the name of a function in Go?](https://stackoverflow.com/questions/7052693/how-to-get-the-name-of-a-function-in-go)
15. [metadata](https://whatis.techtarget.com/definition/metadata)
16. [grpc-metadata](https://github.com/grpc/grpc-go/blob/master/Documentation/grpc-metadata.md)







