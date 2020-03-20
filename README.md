# 环境说明
* Hyperleger/fabric v1.4.4  
* Hyperleger/caliepr-cli v0.3.0  
* node v8.10.0 
* npm v5.6.0 
* docker 19.03.5  
* docker-compose 1.25.4  

# 文件说明

```
chaincode/demo:       测试 chaincode
chaincode/callback:   hyperleger/caliper测试用例

pbft:                 可插拔 PBFT 共识算法简单实现

solo-network:          solo共识配置
pbft-network:          pbft共识配置 
multi-channel-network: solo多链配置
```

## chaincode 说明

| 函数 |       功能       |    参数    |
| :-------: | :--------------: | :--------------------: |
|  open  | 开户 | 账户名, 金额 |
|  query  | 查询 | 账户名 |
|  invoke  | 转账 | 账户名, 账户名, 金额 |
|  delete  | 销户 | 账户名 |


## 网络拓扑 solo-network

| 类型/组织 |      域名       |     IP/端口     |   组织名   |
| :-------: | :-------------: | :-------------: | :--------: |
|  Orderer  | orderer.yzm.com | 172.22.0.2:7050 | OrdererOrg |
| Peer/OrgA | peer0.orga.com  | 172.22.0.3:7051 |  OrgAMSP   |
| Peer/OrgB | peer0.orgb.com  | 172.22.0.4:8051 |  OrgBMSP   |

## 网络拓扑 pbft-network

| 类型/组织 |       域名       |    IP/端口/PBFT端口    |   组织名   |
| :-------: | :--------------: | :--------------------: | :--------: |
|  Orderer  | orderer0.yzm.com | 172.22.0.100:6050/6070 | OrdererOrg |
|  Orderer  | orderer1.yzm.com | 172.22.0.101:6051/6071 | OrdererOrg |
|  Orderer  | orderer2.yzm.com | 172.22.0.101:6052/6072 | OrdererOrg |
|  Orderer  | orderer3.yzm.com | 172.22.0.101:6053/6073 | OrdererOrg |
| Peer/OrgA |  peer0.orga.com  |    172.22.0.2:7051     |  OrgAMSP   |
| Peer/OrgB |  peer0.orgb.com  |    172.22.0.3:8051     |  OrgBMSP   |


## 网络拓扑 multi-channel-network

| 类型/组织 |          域名           | IP/端口/PBFT端口 |   组织名   |
| :-------: | :---------------------: | :--------------: | :--------: |
|  Orderer  |   orderer.example.com   | 172.22.0.2:7050  | OrdererOrg |
| Peer/OrgA | peer0.orgA.example.com | 172.22.0.3:7051  |  OrgAMSP   |

| 类型/组织 |          域名           | IP/端口/PBFT端口 |   组织名   |
| :-------: | :---------------------: | :--------------: | :--------: |
| Peer/OrgB | peer0.orgB.example.com | 172.22.0.4:8051  |  OrgBMSP   |
| Peer/OrgB | peer1.orgB.example.com | 172.22.0.5:8061  |  OrgBMSP   |
| Peer/OrgB | peer2.orgB.example.com | 172.22.0.6:8071  |  OrgBMSP   |
| Peer/OrgB | peer3.orgB.example.com | 172.22.0.7:8081  |  OrgBMSP   |
| Peer/OrgB | peer4.orgB.example.com | 172.22.0.8:8091  |  OrgBMSP   |

| 类型/组织 |          域名           | IP/端口/PBFT端口 |   组织名   |
| :-------: | :---------------------: | :--------------: | :--------: |
| Peer/OrgC | peer0.orgC.example.com | 172.22.0.9:9051  |  OrgCMSP   |
| Peer/OrgC | peer1.orgC.example.com | 172.22.0.10:9061  |  OrgCMSP   |
| Peer/OrgC | peer2.orgC.example.com | 172.22.0.11:9071  |  OrgCMSP   |
| Peer/OrgC | peer3.orgC.example.com | 172.22.0.12:9081  |  OrgCMSP   |
| Peer/OrgC | peer4.orgC.example.com | 172.22.0.13:9091  |  OrgCMSP   |

| 类型/组织 |          域名           | IP/端口/PBFT端口 |   组织名   |
| :-------: | :---------------------: | :--------------: | :--------: |
| Peer/OrgD | peer0.orgD.example.com | 172.22.0.14:10051  |  OrgDMSP   |
| Peer/OrgD | peer1.orgD.example.com | 172.22.0.15:10061  |  OrgDMSP   |
| Peer/OrgD | peer2.orgD.example.com | 172.22.0.16:10071  |  OrgDMSP   |
| Peer/OrgD | peer3.orgD.example.com | 172.22.0.17:10081  |  OrgDMSP   |
| Peer/OrgD | peer4.orgD.example.com | 172.22.0.18:10091  |  OrgDMSP   |
| Peer/OrgD | peer5.orgD.example.com | 172.22.0.19:10101  |  OrgDMSP   |