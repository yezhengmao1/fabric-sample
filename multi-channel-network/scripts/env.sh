#!/bin/bash

# 证书文件夹
PEERROOT=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations
ORDEROOT=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations

# 节点设置
ORDERERNODE=orderer.example.com:7050
PEER0ORGANODE=peer0.orgA.example.com:7051
PEER0ORGBNODE=peer0.orgB.example.com:8051
PEER1ORGBNODE=peer1.orgB.example.com:8061
PEER2ORGBNODE=peer2.orgB.example.com:8071
PEER3ORGBNODE=peer4.orgB.example.com:8081
PEER4ORGBNODE=peer5.orgB.example.com:8091
PEER0ORGCNODE=peer0.orgC.example.com:9051
PEER1ORGCNODE=peer1.orgC.example.com:9061
PEER2ORGCNODE=peer2.orgC.example.com:9071
PEER3ORGCNODE=peer3.orgC.example.com:9081
PEER4ORGCNODE=peer4.orgC.example.com:9091
PEER0ORGDNODE=peer0.orgD.example.com:10051
PEER1ORGDNODE=peer1.orgD.example.com:10061
PEER2ORGDNODE=peer2.orgD.example.com:10071
PEER3ORGDNODE=peer3.orgD.example.com:10081
PEER4ORGDNODE=peer4.orgD.example.com:10091
PEER5ORGDNODE=peer5.orgD.example.com:10101

CHANNEL_NAME=(channelabcd channelbc channelbcd channelcd)

NAME=money_demo
VERSION=1.0

# 切换peer0 orgA
OrgA(){
    CORE_PEER_MSPCONFIGPATH=${PEERROOT}/orgA.example.com/users/Admin@orgA.example.com/msp
    CORE_PEER_ADDRESS=${PEER0ORGANODE}
    CORE_PEER_LOCALMSPID="OrgAMSP"
    echo "node now:peer0.orga.com"
}

# 切换peer0 orgB
OrgB(){
    CORE_PEER_MSPCONFIGPATH=${PEERROOT}/orgB.example.com/users/Admin@orgB.example.com/msp
    CORE_PEER_ADDRESS=${PEER0ORGBNODE}
    CORE_PEER_LOCALMSPID="OrgBMSP"
    echo "node now:peer0.orgb.com"
}

# 切换peer0 orgC
OrgC(){
    CORE_PEER_MSPCONFIGPATH=${PEERROOT}/orgC.example.com/users/Admin@orgC.example.com/msp
    CORE_PEER_ADDRESS=${PEER0ORGCNODE}
    CORE_PEER_LOCALMSPID="OrgCMSP"
    echo "node now:peer0.orgc.com"
}

# 切换peer0 orgD
OrgC(){
    CORE_PEER_MSPCONFIGPATH=${PEERROOT}/orgD.example.com/users/Admin@orgD.example.com/msp
    CORE_PEER_ADDRESS=${PEER0ORGDNODE}
    CORE_PEER_LOCALMSPID="OrgDMSP"
    echo "node now:peer0.orgd.com"
}

# 安装channel
InstallChannel() {
    # 所有channel包含OrgC - 使用OrgC创建channel
    OrgC
    for i in ${CHANNEL_NAME[@]}; do
        peer channel create \
            -o ${ORDERERNODE} \
            -c ${i} \
            -f ./channel-artifacts/${i}.tx
        echo "install channel " ${i} " done !"
        sleep 3
    done
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
        -n ${NAME} \
        -v ${VERSION} \
        -p github.com/chaincode/demo/
    echo "peer0.orga.com install chaincode - demo"

    OrgB
    peer chaincode install \
        -n ${NAME} \
        -v ${VERSION} \
        -p github.com/chaincode/demo/
    echo "peer0.orgb.com install chaincode - demo"
}

# 实例链码
InstantiateChainCode() {
    peer chaincode instantiate \
        -o ${ORDERERNODE} \
        -C ${CHANNEL_NAME} \
        -n ${NAME} \
        -v ${VERSION} \
        -c '{"Args":["Init"]}' \
        -P "AND ('OrgAMSP.peer','OrgBMSP.peer')"
    sleep 10
    echo "instantiate chaincode"
}

# 链码测试
TestDemo() {
    # 创建账户
    peer chaincode invoke \
        -C ${CHANNEL_NAME} \
        -o ${ORDERERNODE} \
        -n ${NAME} \
        --peerAddresses ${PEERORGANODE} \
        --peerAddresses ${PEERORGBNODE} \
        -c '{"Args":["open","count_a", "100"]}'
    sleep 3
    peer chaincode invoke \
        -C ${CHANNEL_NAME} \
        -o ${ORDERERNODE} \
        -n ${NAME} \
        --peerAddresses ${PEERORGANODE} \
        --peerAddresses ${PEERORGBNODE} \
        -c '{"Args":["open","count_b", "100"]}'
    sleep 3
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
        --peerAddresses ${PEERORGBNODE} \
        -c '{"Args":["invoke","count_a","count_b","50"]}'
    sleep 3
    peer chaincode invoke \
        -C ${CHANNEL_NAME} \
        -o ${ORDERERNODE} \
        -n ${NAME} \
        --peerAddresses ${PEERORGANODE} \
        --peerAddresses ${PEERORGBNODE} \
        -c '{"Args":["open","count_c", "100"]}'
    sleep 3
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
        InstallChannel
        JoinChannel
        AnchorUpdate
        InstallChainCode
        InstantiateChainCode
        OrgA
        TestDemo
        ;;
esac
