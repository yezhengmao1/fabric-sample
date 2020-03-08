# 环境说明
* Hyperleger/fabric v1.4.4  
* Hyperleger/caliepr-cli v0.3.0  
* node v8.10.0 
* npm v5.6.0 
* docker 19.03.5  
* docker-compose 1.25.4  

# 文件说明

```
chaincode/demo: 测试chaincode
solo-network: solo共识配置
```

## chaincode/demo

* Init : 无
* Invoke : 
  * open - 开户，参数：<"账户名"，"金额">
  * query - 查询，参数：<"账户名">
  * invoke - 转账，参数：<"账户名"，"账户名"，"金额">
  * delete - 销户，参数：<"账户名">

## solo-network

网络：

```
orderer:
	172.22.0.2/orderer.yzm.com:7050
peer:
	172.22.0.3/peer0.orga.com:7051
	172.22.0.4/peer0.orgb.com:8051
```

## multi-channel-network

网络（* - 锚节点）

```
orderer:
	172.22.0.2/orderer.example.com:7050
peer:
	orga:
		172.22.0.3/peer0.orgA.example.com:7051(*)
	orgb:
		172.22.0.4/peer0.orgB.example.com:8051(*)
		172.22.0.5/peer1.orgB.example.com:8061
		172.22.0.6/peer2.orgB.example.com:8071
		172.22.0.7/peer3.orgB.example.com:8081
		172.22.0.8/peer4.orgB.example.com:8091
	orgc:
		172.22.0.9/peer0.orgC.example.com:9051(*)
		172.22.0.10/peer1.orgC.example.com:9061
		172.22.0.11/peer2.orgC.example.com:9071
		172.22.0.12/peer3.orgC.example.com:9081
		172.22.0.13/peer4.orgC.example.com:9091
	orgd:
		172.22.0.14/peer0.orgD.example.com:10051(*)
		172.22.0.15/peer1.orgD.example.com:10061
		172.22.0.16/peer2.orgD.example.com:10071
		172.22.0.17/peer3.orgD.example.com:10081
		172.22.0.18/peer4.orgD.example.com:10091
		172.22.0.19/peer5.orgD.example.com:10101
```

