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
PEER3ORGBNODE=peer3.orgB.example.com:8081
PEER4ORGBNODE=peer4.orgB.example.com:8091
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
CHANNEL_INSTANTIATE=(
    "AND ('OrgAMSP.peer','OrgBMSP.peer','OrgCMSP.peer','OrgDMSP.peer')"
    "AND ('OrgBMSP.peer','OrgCMSP.peer')"
    "AND ('OrgCMSP.peer','OrgDMSP.peer')"
    "AND ('OrgCMSP.peer','OrgDMSP.peer')"
    )
CHANNELABCD=channelabcd
CHANNELBC=channelbc
CHANNELBCD=channelbcd
CHANNELCD=channelcd


CHANNEL_A=($CHANNELABCD)
CHANNEL_B=($CHANNELABCD $CHANNELBC $CHANNELBCD)
CHANNEL_C=($CHANNELABCD $CHANNELBC $CHANNELBCD $CHANNELCD)
CHANNEL_D=($CHANNELABCD $CHANNELBCD $CHANNELCD)

OrgA_PEERS=(${PEER0ORGANODE})
OrgB_PEERS=(${PEER0ORGBNODE} ${PEER1ORGBNODE} ${PEER2ORGBNODE} ${PEER3ORGBNODE} ${PEER4ORGBNODE})
OrgC_PEERS=(${PEER0ORGCNODE} ${PEER1ORGCNODE} ${PEER2ORGCNODE} ${PEER3ORGCNODE} ${PEER4ORGCNODE})
OrgD_PEERS=(${PEER0ORGDNODE} ${PEER1ORGDNODE} ${PEER2ORGDNODE} ${PEER3ORGDNODE} ${PEER4ORGDNODE} ${PEER5ORGDNODE})

NAME=money_demo
VERSION=1.0

# 切换peer0 orgA
OrgA(){
    CORE_PEER_MSPCONFIGPATH=${PEERROOT}/orgA.example.com/users/Admin@orgA.example.com/msp
    CORE_PEER_ADDRESS=${PEER0ORGANODE}
    CORE_PEER_LOCALMSPID="OrgAMSP"
    echo "org now: orga; node now:peer0"
}

# 切换peer0 orgB
OrgB(){
    CORE_PEER_MSPCONFIGPATH=${PEERROOT}/orgB.example.com/users/Admin@orgB.example.com/msp
    CORE_PEER_ADDRESS=${PEER0ORGBNODE}
    CORE_PEER_LOCALMSPID="OrgBMSP"
    echo "org now: orgb; node now:peer0"
}

# 切换peer0 orgC
OrgC(){
    CORE_PEER_MSPCONFIGPATH=${PEERROOT}/orgC.example.com/users/Admin@orgC.example.com/msp
    CORE_PEER_ADDRESS=${PEER0ORGCNODE}
    CORE_PEER_LOCALMSPID="OrgCMSP"
    echo "org now: orgc; node now:peer0"
}

# 切换peer0 orgD
OrgD(){
    CORE_PEER_MSPCONFIGPATH=${PEERROOT}/orgD.example.com/users/Admin@orgD.example.com/msp
    CORE_PEER_ADDRESS=${PEER0ORGDNODE}
    CORE_PEER_LOCALMSPID="OrgDMSP"
    echo "org now: orgd; node now:peer0"
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
        sleep 1
    done
}

# OrgA加入的channel
JoinChannelA() {
	peer channel join -b ${CHANNELABCD}.block
}

# OrgB加入的channel
JoinChannelB() {
	peer channel join -b ${CHANNELABCD}.block
    peer channel join -b ${CHANNELBC}.block
    peer channel join -b ${CHANNELBCD}.block
}

# OrgC加入的channel
JoinChannelC() {
	peer channel join -b ${CHANNELABCD}.block
    peer channel join -b ${CHANNELBC}.block
    peer channel join -b ${CHANNELBCD}.block
	peer channel join -b ${CHANNELCD}.block
}

# OrgD加入的channel
JoinChannelD() {
	peer channel join -b ${CHANNELABCD}.block
    peer channel join -b ${CHANNELBCD}.block
	peer channel join -b ${CHANNELCD}.block
}

# 加入channel
JoinChannel() {
    OrgA
	for i in ${OrgA_PEERS[@]}; do
        CORE_PEER_ADDRESS=${i}
        JoinChannelA
        echo ${i}" join channel"
	done
    OrgB
    for i in ${OrgB_PEERS[@]}; do
        CORE_PEER_ADDRESS=${i}
        JoinChannelB
        echo ${i}" join channel"
    done
    OrgC
    for i in ${OrgC_PEERS[@]}; do
        CORE_PEER_ADDRESS=${i}
        JoinChannelC
        echo ${i}" join channel"
    done
    OrgD
    for i in ${OrgD_PEERS[@]}; do
        CORE_PEER_ADDRESS=${i}
        JoinChannelD
        echo ${i}" join channel"
    done
}

AnchorUpdateA() {
    for i in ${CHANNEL_A[@]}; do
        peer channel update \
            -o ${ORDERERNODE} \
            -c ${i} \
            -f ./channel-artifacts/OrgAMSPanchor_${i}.tx
    done
}

AnchorUpdateB() {
    for i in ${CHANNEL_B[@]}; do
         peer channel update \
            -o ${ORDERERNODE} \
            -c ${i} \
            -f ./channel-artifacts/OrgBMSPanchor_${i}.tx
    done
}

AnchorUpdateC() {
    for i in ${CHANNEL_C[@]}; do
        peer channel update \
            -o ${ORDERERNODE} \
            -c ${i} \
            -f ./channel-artifacts/OrgCMSPanchor_${i}.tx
    done
}

AnchorUpdateD() {
    for i in ${CHANNEL_D[@]}; do
        peer channel update \
            -o ${ORDERERNODE} \
            -c ${i} \
            -f ./channel-artifacts/OrgDMSPanchor_${i}.tx
    done
}

# 更新锚节点
AnchorUpdate() {
    OrgA
    AnchorUpdateA
    OrgB
    AnchorUpdateB
    OrgC
    AnchorUpdateC
    OrgD
    AnchorUpdateD
}

InstallChainCodeFunc() {
    peer chaincode install \
        -n ${NAME} \
        -v ${VERSION} \
        -p github.com/chaincode/demo/
}

# 安装链码
InstallChainCode() {
    OrgA
    for i in ${OrgA_PEERS[@]}; do
        CORE_PEER_ADDRESS=${i}
        InstallChainCodeFunc
        echo ${i}
    done
    OrgB
    for i in ${OrgB_PEERS[@]}; do
        CORE_PEER_ADDRESS=${i}
        InstallChainCodeFunc
        echo ${i}
    done
    OrgC
    for i in ${OrgC_PEERS[@]}; do
        CORE_PEER_ADDRESS=${i}
        InstallChainCodeFunc
        echo ${i}
    done
    OrgD
    for i in ${OrgD_PEERS[@]}; do
        CORE_PEER_ADDRESS=${i}
        InstallChainCodeFunc
        echo ${i}
    done
}

# 实例链码
InstantiateChainCode() {
    OrgC
    for i in ${!CHANNEL_NAME[@]}; do
        peer chaincode instantiate \
            -o ${ORDERERNODE} \
            -C ${CHANNEL_NAME[i]} \
            -n ${NAME} \
            -v ${VERSION} \
            -c '{"Args":["Init"]}' \
            -P "${CHANNEL_INSTANTIATE[i]}"
        sleep 1
    done
    for i in ${CHANNEL_NAME[@]}; do
        peer chaincode list --instantiated -C ${i}
    done
}

# 链码测试
TestDemo() {
    OrgC
    # 创建账户
    peer chaincode invoke \
        -C ${CHANNELABCD} \
        -o ${ORDERERNODE} \
        -n ${NAME} \
        --peerAddresses ${PEER0ORGANODE} \
        --peerAddresses ${PEER0ORGBNODE} \
        --peerAddresses ${PEER0ORGCNODE} \
        --peerAddresses ${PEER0ORGDNODE} \
        -c '{"Args":["open","count_a", "100"]}'
    peer chaincode invoke \
        -C ${CHANNELABCD} \
        -o ${ORDERERNODE} \
        -n ${NAME} \
        --peerAddresses ${PEER0ORGANODE} \
        --peerAddresses ${PEER0ORGBNODE} \
        --peerAddresses ${PEER0ORGCNODE} \
        --peerAddresses ${PEER0ORGDNODE} \
        -c '{"Args":["open","count_b", "100"]}'
    peer chaincode invoke \
        -C ${CHANNELABCD} \
        -o ${ORDERERNODE} \
        -n ${NAME} \
        --peerAddresses ${PEER0ORGANODE} \
        --peerAddresses ${PEER0ORGBNODE} \
        --peerAddresses ${PEER0ORGCNODE} \
        --peerAddresses ${PEER0ORGDNODE} \
        -c '{"Args":["invoke","count_a", "count_b","1"]}'
    peer chaincode query \
        -C ${CHANNELABCD} \
        -n ${NAME} \
        -c '{"Args":["query","count_a"]}'
    peer chaincode query \
        -C ${CHANNELABCD} \
        -n ${NAME} \
        -c '{"Args":["query","count_b"]}'
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
