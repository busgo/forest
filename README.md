# forest

>   分布式任务调度系统

### 前言
在企业系统开发过程中难免少不了一些定时任务来进行定时触发执行任务，对于非分布式环境系统中，我们只需要在对应系统中内部集成一些调度库进行配置定时触发即可。
比如：使用Spring框架集成quartz,只需要进行一些简单的配置就能定时执行任务了。但是随着企业的系统越来越多、逐步从单一应用慢慢演变为微服务集群。
在分布式集群系统中主要面临出如：任务的重复执行、没有统一定时任务配置、任务节点故障转移、任务监控&报警等一些列的功能都是要在分布式系统中进行解决。
### 系统架构


//TODO 

### 选举
由于一个任务调度集群有多台提供服务，我们在可以从集群节点中选举出一台领导节点来进行发号师令，比较成熟的选举算法(Paxos、Raft 等)这里不做讨论。这里使用etcd中的租约机制来实现选举功能。

当一个调度服务节点启动的时候首先尝试发起选举请求(PUT 节点 /forest/server/leader/),如果执行成功则选举成功。如果判断已经有其他调度服务节点已经选举成功过则放弃选举请求并进行监听(/forest/server/leader/)选举节点变化。如果有领导下线通知则立即发起选举。


### 快速开始

####    先决条件
   *    golang(>=1.11)
   *    git 
   
####    源代码安装

```shell
    git clone https://github.com/busgo/forest.git
    cd forest/forest
    go build forest.go
```
等待自动下载依赖库文件

```go


appledeMacBook-Pro:forest apple$ go build forest.go 
go: finding github.com/xiang90/probing v0.0.0-20190116061207-43a291ad63a2
go: finding github.com/dgrijalva/jwt-go v3.2.0+incompatible
go: finding github.com/golang/groupcache v0.0.0-20190702054246-869f871628b6
go: finding github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f
go: finding github.com/coreos/go-systemd v0.0.0-20190620071333-e64a0ec8b42a
go: finding github.com/prometheus/client_golang v1.0.0
go: finding github.com/coreos/bbolt v1.3.3
go: finding github.com/prometheus/common v0.4.1
go: finding github.com/grpc-ecosystem/go-grpc-middleware v1.0.0
go: finding github.com/coreos/etcd v3.3.13+incompatible
go: finding github.com/tmc/grpc-websocket-proxy v0.0.0-20190109142713-0ad062ec5ee5
go: finding github.com/grpc-ecosystem/grpc-gateway v1.9.4
go: finding github.com/gogo/protobuf v1.1.1
go: finding github.com/alecthomas/template v0.0.0-20160405071501-a0175ee3bccc
go: finding github.com/prometheus/procfs v0.0.0-20181005140218-185b4288413d
go: finding gopkg.in/yaml.v2 v2.2.1
...
```

