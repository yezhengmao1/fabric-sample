#!/bin/bash

ORGA=orgA.example.com
ORGB=orgB.example.com
ORGC=orgC.example.com
ORGD=orgD.example.com
ORGAUSERS=(Admin)
ORGBUSERS=(Admin User1 User2 User3 User4)
ORGCUSERS=(Admin User1 User2 User3 User4)
ORGDUSERS=(Admin User1 User2 User3 User4 User5)
VERSION=1.4.4

# 复制keystore
CPFile() {
    files=$(ls $1)
    echo ${files[0]}
    cd $1
    cp ${files[0]} ./key.pem
    cd -
}

# 复制所有文件keystore
CPAllFiles() {
    PREFIX=crypto-config/peerOrganizations
    SUFFIX=msp/keystore
    for u in ${ORGAUSERS[@]}; do
        CPFile ${PREFIX}/${ORGA}/users/${u}@${ORGA}/${SUFFIX}
    done
    for u in ${ORGBUSERS[@]}; do
        CPFile ${PREFIX}/${ORGB}/users/${u}@${ORGB}/${SUFFIX}
    done
    for u in ${ORGCUSERS[@]}; do
        CPFile ${PREFIX}/${ORGC}/users/${u}@${ORGC}/${SUFFIX}
    done
    for u in ${ORGDUSERS[@]}; do
        CPFile ${PREFIX}/${ORGD}/users/${u}@${ORGD}/${SUFFIX}
    done
}

# 清理缓存文件
Clean() {
    rm -rf ./channel-artifacts
    rm -rf ./crypto-config
    rm -rf ./production
    rm -rf /tmp/crypto
}

case $1 in
    # 压力测试启动/关闭
    up)
        CPAllFiles
        env IMAGETAG=${VERSION} docker-compose -f ./docker-compose-cli.yaml up -d
        docker exec cli /bin/bash -c "scripts/env.sh all"
        ;;
    down)
        docker kill $(docker ps -qa)
        echo y | docker system prune
        Clean
        ;;
esac
