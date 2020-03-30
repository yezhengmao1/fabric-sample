#!/bin/bash

ORGA=orga
ORGAUSERS=(Admin User1)
VERSION=rbft

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
        CPFile ${PREFIX}/${ORGA}.com/users/${u}@${ORGA}.com/${SUFFIX}
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
    cli)
        env IMAGETAG=${VERSION} docker-compose -f ./docker-compose-cli.yaml up -d
        docker exec -ti cli /bin/bash
        ;;
    up)
        CPAllFiles
        env IMAGETAG=${VERSION} docker-compose -f ./docker-compose-cli.yaml up -d
        sleep 3
        docker exec cli /bin/bash -c "scripts/env.sh all"
        ;;
    down)
        docker kill $(docker ps -qa)
        echo y | docker system prune
        docker rmi $(docker images | grep 'dev-*' | awk '{print $3}')
        echo y | docker system prune
        Clean
        ;;
esac
