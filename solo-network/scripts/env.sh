#!/bin/sh

# 证书文件夹
PEERROOT=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations
ORDEROOT=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations

# 节点设置
ORDERERNODE=orderer.yzm.com:7050
PEERORGANODE=peer0.orga.com:7051
PEERORGBNODE=peer0.orgb.com:8051
CHANNEL_NAME=mychannel

# 切换peer0 orgA
OrgA(){
    CORE_PEER_MSPCONFIGPATH=${PEERROOT}/orga.com/users/Admin@orga.com/msp
    CORE_PEER_ADDRESS=${PEERORGANODE}
    CORE_PEER_LOCALMSPID="OrgAMSP"
    echo "node now:peer0.orga.com"
}

# 切换peer0 orgB
OrgB(){
    CORE_PEER_MSPCONFIGPATH=${PEERROOT}/orgb.com/users/Admin@orgb.com/msp
    CORE_PEER_ADDRESS=${PEERORGBNODE}
    CORE_PEER_LOCALMSPID="OrgBMSP"
    echo "node now:peer0.orgb.com"
}

# 安装channel
InstallChannel() {
    peer channel create \
        -o ${ORDERERNODE} \
        -c ${CHANNEL_NAME} \
        -f ./channel-artifacts/channel.tx \
    echo "install channel"
}

# 加入channel
JoinChannel() {
    OrgA
    peer channel join -b ${CHANNEL_NAME}.block
    echo "peer0.orga.com join channel" 
    OrgB
    peer channel join -b ${CHANNEL_NAME}.block
    echo "peer0.orgb.com join channel"
}

# 更新锚节点
AnchorUpdate() {
    OrgA
    peer channel update \
        -o ${ORDERERNODE} \
        -c ${CHANNEL_NAME} \
        -f ./channel-artifacts/OrgAMSPanchor.tx \
    echo "orga update anchor peer0.orga.com"
    OrgB
    peer channel update \
        -o ${ORDERERNODE} \
        -c ${CHANNEL_NAME} \
        -f ./channel-artifacts/OrgBMSPanchor.tx \
    echo "orgb update anchor peer0.orgb.com"
}

# 安装链码
InstallChainCode() {
    OrgA
    peer chaincode install \
        -n demo \
        -v 1.0 \
        -p github.com/chaincode/demo/
    echo "peer0.orga.com install chaincode - demo"

    OrgB
    peer chaincode install \
        -n demo \
        -v 1.0 \
        -p github.com/chaincode/demo/
    echo "peer0.orgb.com install chaincode - demo"
}

# 实例链码
InstantiateChainCode() {
    peer chaincode instantiate \
        -o ${ORDERERNODE} \
        -C ${CHANNEL_NAME} \
        -n demo \
        -v 1.0 \
        -c '{"Args":["Init","a","100","b","100"]}' \
        -P "AND ('OrgAMSP.peer','OrgBMSP.peer')"
    sleep 3
    echo "instantiate chaincode"
}

# 链码测试
TestDemo() {
    peer chaincode query \
        -C ${CHANNEL_NAME} \
        -n demo \
        -c '{"Args":["query","a"]}'
    peer chaincode query \
        -C ${CHANNEL_NAME} \
        -n demo \
        -c '{"Args":["query","b"]}'
    peer chaincode invoke \
        -C ${CHANNEL_NAME} \
        -o ${ORDERERNODE} \
        -n demo \
        --peerAddresses ${PEERORGANODE} \
        --peerAddresses ${PEERORGBNODE} \
        -c '{"Args":["invoke","a","b","1"]}'
    # 等待共识完成
    sleep 3
    peer chaincode query \
        -C ${CHANNEL_NAME} \
        -n demo \
        -c '{"Args":["query","a"]}'
    peer chaincode query \
        -C ${CHANNEL_NAME} \
        -n demo \
        -c '{"Args":["query","b"]}'
}

case $1 in
    installchannel)
        InstallChannel
        ;;
    joinchannel)
        JoinChannel
        ;;
    anchorupdate)
        AnchorUpdate
        ;;
    installchaincode)
        InstallChainCode
        ;;
    instantiatechaincode)
        InstantiateChainCode
        ;;
    testdemo)
        TestDemo
        ;;
    all)
        InstallChannel
        JoinChannel
        AnchorUpdate
        InstallChainCode
        InstantiateChainCode
        TestDemo
        ;;
esac
