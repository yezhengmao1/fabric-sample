#!/bin/bash

GENESIS_PROFILE=Genesis
CHANNEL_PROFILE=Channel
SYS_CHANNEL=sys-channel
CHANNEL_NAME=mychannel
VERSION=1.4.4

FABRIC_CFG_PATH=$PWD

ORG_NAMES=(OrgAMSP OrgBMSP)

# 检测cryptogen和版本
if ! [ -x "$(command -v cryptogen)" ] ; then
    echo -e "\033[31m no cryptogen\033[0m"
    exit 1
fi
if [ ${VERSION} != "$(cryptogen version | grep Version | awk -F ': ' '{print $2}')" ] ; then
    echo -e "\033[31m cryptogen need version \033[0m"${VERSION}
    exit 1
fi
# 检测configtxgen和版本
if ! [ -x "$(command -v configtxgen)" ] ; then 
    echo -e "\033[31m no configtxgen\033[0m"
    exit 1
fi
if [ ${VERSION} != "$(configtxgen --version | grep Version | awk -F ': ' '{print $2}')" ] ; then
    echo -e "\033[31m configtxgen need version \033[0m"${VERSION}
    exit 1
fi
# 生成证书文件
echo -e "\033[31m clear crypto files\033[0m"
rm -rf crypto-config
echo -e "\033[31m generate crypto files\033[0m"
cryptogen generate --config ./crypto-config.yaml
# 清理多余文件
echo -e "\033[31m clear block files\033[0m"
rm -rf ./channel-artifacts
mkdir ./channel-artifacts
# 生成创世块
echo -e "\033[31m generate genesis block\033[0m"
configtxgen \
    -profile ${GENESIS_PROFILE} \
    -channelID ${SYS_CHANNEL} \
    -outputBlock ./channel-artifacts/genesis.block \
# 生成通道交易
echo -e "\033[31m generate channel transcation\033[0m"
configtxgen \
    -profile ${CHANNEL_PROFILE} \
    -channelID ${CHANNEL_NAME} \
    -outputCreateChannelTx ./channel-artifacts/channel.tx
# 生成铆节点配置
echo -e "\033[31m generate anchor transcation\033[0m"
for i in ${ORG_NAMES[@]}; do
    configtxgen \
        -profile ${CHANNEL_PROFILE} \
        -channelID ${CHANNEL_NAME} \
        -outputAnchorPeersUpdate ./channel-artifacts/${i}anchor.tx \
        -asOrg ${i}
done
