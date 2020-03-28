## 一、Fabric 可插拔共识算法 PBFT 开发流程

* `configtxgen`工具源码修改，使其识别`pbft`共识配置。

```go
// common/tools/configtxgen/localconfig/config.go:388
switch ord.OrdererType {
    case 'pbft':
}
// commom/tools/configtxgen/encoder/encoder.go:38
const ConsensusTypePbft = "pbft"
// commom/tools/configtxgen/encoder/encoder.go:215
switch conf.OrdererType {
	case ConsensusTypePbft:
}
```

* 添加共识算法实例

```go
// orderer/common/server/main.go:664
consenters["pbft"] = pbft.New()
```

* 实现共识接口`/orderer/consensus/consensus.go`

```go
// 接口说明 - Consneter 
// 返回 Chain 用于实现处理区块接口
type Consenter interface {
	HandleChain(support ConsenterSupport, metadata *cb.Metadata) (Chain, error)
}
// Chain 处理区块接口
type Chain interface {
   	// 处理 Normal 交易
    Order(env *cb.Envelope, configSeq uint64) error
    // 处理配置交易
    Configure(config *cb.Envelope, configSeq uint64) error
    // 等待接收交易,处理函数交易前
	WaitReady() error
    // 发送错误 chan
    Errored() <-chan struct{}
    // 初始化 Chain 中资源
    Start()
    // 资源释放
    Halt()
}
```

* 编译产生 orderer 镜像（修改`orderer\peer\tools` tag 为 `pbft`）

```
$ make orderer-docker
```

* 编译产生 configtxgen 工具（输出目录：`.build/bin/configtxgen`）

```
$ make configtxgen
```

## 二、网络拓扑

| 类型/组织 |       域名       |    IP/端口/PBFT端口    |   组织名   |
| :-------: | :--------------: | :--------------------: | :--------: |
|  Orderer  | orderer0.yzm.com | 172.22.0.100:6050/6070 | OrdererOrg |
|  Orderer  | orderer1.yzm.com | 172.22.0.101:6051/6071 | OrdererOrg |
|  Orderer  | orderer2.yzm.com | 172.22.0.101:6052/6072 | OrdererOrg |
|  Orderer  | orderer3.yzm.com | 172.22.0.101:6053/6073 | OrdererOrg |
| Peer/OrgA |  peer0.orga.com  |    172.22.0.2:7051     |  OrgAMSP   |

## 三、配置说明

采用环境变量：

* `PBFT_LISTEN_PORT`：PBFT 节点监听端口
* `PBFT_NODE_ID`：PBFT 节点 ID
* `PBFT_NODE_TABLE`：PBFT 网络列表
