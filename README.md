# 说明

* [PBFT 共识实现与网络搭建方法](https://www.yezhem.com/index.php/archives/52/)    
* [SOLO 共识网络搭建方法](https://www.yezhem.com/index.php/archives/39/)

# 环境
* Hyperleger/fabric v1.4.4  
* Hyperleger/caliepr-cli v0.3.0  
* node v8.10.0 
* npm v5.6.0 
* docker 19.03.5  
* docker-compose 1.25.4  

# 文件

```
chaincode/demo:       测试 chaincode
chaincode/callback:   hyperleger/caliper测试用例

pbft:                 可插拔 PBFT 共识算法简单实现
rbft:                 可插拔 RBFT 共识算法简单实现

solo-network:          solo共识配置
pbft-network:          pbft共识配置 
rbft-network:          rbft共识配置
multi-channel-network: solo多链配置
```

# 链码

| 函数 |       功能       |    参数    |
| :-------: | :--------------: | :--------------------: |
|  open  | 开户 | 账户名, 金额 |
|  query  | 查询 | 账户名 |
|  invoke  | 转账 | 账户名, 账户名, 金额 |
|  delete  | 销户 | 账户名 |

# 编译

[编译 pbft 说明文件](https://github.com/yezhem/fabric-sample/blob/master/pbft/doc.md)：`./pbft/doc.md` 

[编译 rbft 说明文件](https://github.com/yezhem/fabric-sample/blob/master/rbft/doc.md)：`./rbft/doc.md`

# 测试

```
$ npx caliper launch master --caliper-workspace <pbft或rbft-network> --caliper-benchconfig benchmarks/config.yaml --caliper-networkconfig benchmarks/network.yaml
```



