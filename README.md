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

