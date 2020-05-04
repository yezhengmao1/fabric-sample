#!/bin/sh

# 证书文件夹
PEERROOT=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations
ORDEROOT=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations

# 节点设置
ORDERER0NODE=orderer0.yzm.com:6050
ORDERER1NODE=orderer1.yzm.com:6051
ORDERER2NODE=orderer2.yzm.com:6052
ORDERER3NODE=orderer3.yzm.com:6053

ORDERERNODE=${ORDERER1NODE}

PEERORGANODE=peer0.orga.com:7051
CHANNEL_NAME=mychannel

NAME=money_demo
VERSION=1.0

# 切换peer0 orgA
OrgA(){
    CORE_PEER_MSPCONFIGPATH=${PEERROOT}/orga.com/users/Admin@orga.com/msp
    CORE_PEER_ADDRESS=${PEERORGANODE}
    CORE_PEER_LOCALMSPID="OrgAMSP"
    echo "node now:peer0.orga.com"
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
}

# 更新锚节点
AnchorUpdate() {
    OrgA
    peer channel update \
        -o ${ORDERERNODE} \
        -c ${CHANNEL_NAME} \
        -f ./channel-artifacts/OrgAMSPanchor.tx \
    echo "orga update anchor peer0.orga.com"
}

# 安装链码
InstallChainCode() {
    OrgA
    peer chaincode install \
        -n ${NAME} \
        -v ${VERSION} \
        -p github.com/chaincode/demo/
    echo "peer0.orga.com install chaincode - demo"
}

# 实例链码
InstantiateChainCode() {
    peer chaincode instantiate \
        -o ${ORDERERNODE} \
        -C ${CHANNEL_NAME} \
        -n ${NAME} \
        -v ${VERSION} \
        -c '{"Args":["Init"]}' \
        -P "AND ('OrgAMSP.peer')"
    echo "instantiate chaincode"
    sleep 10
}

# 链码测试
TestDemo() {
    # 创建账户
    peer chaincode invoke \
        -C ${CHANNEL_NAME} \
        -o ${ORDERERNODE} \
        -n ${NAME} \
        --peerAddresses ${PEERORGANODE} \
        -c '{"Args":["open","count_a", "100"]}'
    sleep 5
    peer chaincode invoke \
        -C ${CHANNEL_NAME} \
        -o ${ORDERERNODE} \
        -n ${NAME} \
        --peerAddresses ${PEERORGANODE} \
        -c '{"Args":["open","count_b", "100"]}'
    sleep 5
    peer chaincode query \
        -C ${CHANNEL_NAME} \
        -n ${NAME} \
        -c '{"Args":["query","count_a"]}'
    peer chaincode query \
        -C ${CHANNEL_NAME} \
        -n ${NAME} \
        -c '{"Args":["query","count_b"]}'
    peer chaincode invoke \
        -C ${CHANNEL_NAME} \
        -o ${ORDERERNODE} \
        -n ${NAME} \
        --peerAddresses ${PEERORGANODE} \
        -c '{"Args":["invoke","count_a","count_b","50"]}'
    sleep 5
    peer chaincode invoke \
        -C ${CHANNEL_NAME} \
        -o ${ORDERERNODE} \
        -n ${NAME} \
        --peerAddresses ${PEERORGANODE} \
        -c '{"Args":["open","count_c", "100"]}'
    sleep 5
    peer chaincode invoke \
        -C ${CHANNEL_NAME} \
        -o ${ORDERER3NODE} \
        -n ${NAME} \
        --peerAddresses ${PEERORGANODE} \
        -c '{"Args":["invoke","count_a","count_c","10"]}'
    sleep 5
    peer chaincode query \
        -C ${CHANNEL_NAME} \
        -n ${NAME} \
        -c '{"Args":["query","count_a"]}'
    peer chaincode query \
        -C ${CHANNEL_NAME} \
        -n ${NAME} \
        -c '{"Args":["query","count_b"]}'
    peer chaincode query \
        -C ${CHANNEL_NAME} \
        -n ${NAME} \
        -c '{"Args":["query","count_c"]}'
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
        OrgA
        TestDemo
        ;;
    all)
        OrgA
        InstallChannel
        JoinChannel
        AnchorUpdate
        InstallChainCode
        InstantiateChainCode
        TestDemo
        ;;
esac
