#!/bin/sh

PEERROOT=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations
ORDEROOT=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations

ORDERERNODE=orderer.yzm.com:7050
PEERORGANODE=peer0.orga.com:7051
PEERORGBNODE=peer0.orgb.com:8051

ORDERERTLS=${ORDEROOT}/yzm.com/orderers/orderer.yzm.com/msp/tlscacerts/tlsca.yzm.com-cert.pem
PEERORGATLS=${PEERROOT}/orga.com/peers/peer0.orga.com/tls/ca.crt
PEERORGBTLS=${PEERROOT}/orgb.com/peers/peer0.orgb.com/tls/ca.crt

CHANNEL_NAME=mychannel

# 切换peer0 orgA
OrgA(){
    CORE_PEER_MSPCONFIGPATH=${PEERROOT}/orga.com/users/Admin@orga.com/msp
    CORE_PEER_ADDRESS=peer0.orga.com:7051
    CORE_PEER_LOCALMSPID="OrgAMSP"
    CORE_PEER_TLS_ROOTCERT_FILE=${PEERROOT}/orga.com/peers/peer0.orga.com/tls/ca.crt
    echo "node now:peer0.orga.com"
}

# 切换peer0 orgB
OrgB(){
    CORE_PEER_MSPCONFIGPATH=${PEERROOT}/orgb.com/users/Admin@orgb.com/msp
    CORE_PEER_ADDRESS=peer0.orgb.com:8051
    CORE_PEER_LOCALMSPID="OrgBMSP"
    CORE_PEER_TLS_ROOTCERT_FILE=${PEERROOT}/orgb.com/peers/peer0.orgb.com/tls/ca.crt
    echo "node now:peer0.orgb.com"
}

# 安装channel
InstallChannel() {
    peer channel create \
        -o ${ORDERERNODE} \
        -c ${CHANNEL_NAME} \
        -f ./channel-artifacts/channel.tx \
        --tls --cafile ${ORDERERTLS}
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
    peer channel update -o ${ORDERERNODE} -c ${CHANNEL_NAME} -f ./channel-artifacts/OrgAMSPanchor.tx --tls --cafile ${ORDERERTLS}
    echo "orga update anchor peer0.orga.com"
    OrgB
    peer channel update -o ${ORDERERNODE} -c ${CHANNEL_NAME} -f ./channel-artifacts/OrgBMSPanchor.tx --tls --cafile ${ORDERERTLS}
    echo "orgb update anchor peer0.orgb.com"
}

# Demo链码打包
PackageDemo() {
    cd /opt/gopath/src/github.com/pipapa/chaincode/demo
    GO111MODULE=on go mod vendor
    cd -
    peer lifecycle chaincode package mycc.tar.gz \
        --path github.com/pipapa/chaincode/demo \
        --lang golang \
        --label mycc_1
}

# 安装链码
InstallChainCode() {
    OrgA
    peer lifecycle chaincode install mycc.tar.gz
    echo "peer0.orga.com install chaincode - mycc"

    OrgB
    peer lifecycle chaincode install mycc.tar.gz
    echo "peer0.orgb.com install chaincode - mycc"
}

# 获取安装的链码
GetChainCodeID() {
    ChainCodeID=$(peer lifecycle chaincode queryinstalled | tail -n 1 | awk -F'Package ID: ' '{print $2}' | awk -F',' '{print $1}')
}

# 安装链码
ApproveChainCode() {
    OrgA
    GetChainCodeID
    peer lifecycle chaincode approveformyorg \
        --channelID ${CHANNEL_NAME} \
        --name mycc \
        --version 1.0 \
        --init-required \
        --package-id ${ChainCodeID} \
        --sequence 1 \
        --tls true \
        --cafile ${ORDERERTLS}
    echo "peer0.orga.com approveformyorg"

    OrgB
    GetChainCodeID
    peer lifecycle chaincode approveformyorg \
        --channelID ${CHANNEL_NAME} \
        --name mycc \
        --version 1.0 \
        --init-required \
        --package-id ${ChainCodeID} \
        --sequence 1 \
        --tls true \
        --cafile ${ORDERERTLS}
    echo "peer0.orgb.com approveformyorg"

    peer lifecycle chaincode checkcommitreadiness \
        --channelID ${CHANNEL_NAME} \
        --name mycc \
        --version 1.0 \
        --init-required \
        --tls true \
        --cafile ${ORDERERTLS} \
        --output json

    peer lifecycle chaincode commit \
        -o ${ORDERERNODE} \
        --channelID ${CHANNEL_NAME} \
        --name mycc \
        --version 1.0 \
        --sequence 1 \
        --init-required \
        --tls true \
        --cafile ${ORDERERTLS} \
        --peerAddresses ${PEERORGANODE} \
        --tlsRootCertFiles ${PEERORGATLS} \
        --peerAddresses ${PEERORGBNODE} \
        --tlsRootCertFiles ${PEERORGBTLS}
    echo "chaincode commit"
}

DemoInit() {
    peer chaincode invoke \
        -o ${ORDERERNODE} \
        --isInit \
        --tls true \
        --cafile ${ORDERERTLS} \
        -C ${CHANNEL_NAME} \
        -n mycc \
        --peerAddresses ${PEERORGANODE} \
        --tlsRootCertFiles ${PEERORGATLS} \
        --peerAddresses ${PEERORGBNODE} \
        --tlsRootCertFiles ${PEERORGBTLS} \
        -c '{"Args":["Init","a","100","b","100"]}' \
        --waitForEvent
}

DemoQuery() {
    peer chaincode query \
        -C ${CHANNEL_NAME} \
        -n mycc \
        -c '{"Args":["query","a"]}'

    peer chaincode query \
        -C ${CHANNEL_NAME} \
        -n mycc \
        -c '{"Args":["query","b"]}'
}

DemoInvoke() {
    peer chaincode invoke \
        -o ${ORDERERNODE} \
        --tls true \
        --cafile ${ORDERERTLS} \
        -C ${CHANNEL_NAME} \
        -n mycc \
        --peerAddresses ${PEERORGANODE} \
        --tlsRootCertFiles ${PEERORGATLS} \
        --peerAddresses ${PEERORGBNODE} \
        --tlsRootCertFiles ${PEERORGBTLS} \
        -c '{"Args":["invoke","a","b","10"]}' --waitForEvent
}

TestDemo() {
    DemoInit
    DemoQuery
    DemoInvoke
    DemoQuery
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
    packagedemo)
        PackageDemo
        ;;
    testdemo)
        TestDemo
        ;;
    installchaincode)
        InstallChainCode
        ApproveChainCode
        ;;
    all)
        InstallChannel
        JoinChannel
        AnchorUpdate
        PackageDemo
        InstallChainCode
        ApproveChainCode
        TestDemo
        ;;
esac
