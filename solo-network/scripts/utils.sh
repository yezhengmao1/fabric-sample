#!/bin/bash

GetFile() {
    files=$(ls $1)
    echo ${files[0]}
    cd $1
    cp ${files[0]} ./key.pem
    cd -
}

case $1 in
    pk)
        GetFile crypto-config/peerOrganizations/orga.com/users/Admin@orga.com/msp/keystore
        GetFile crypto-config/peerOrganizations/orga.com/users/User1@orga.com/msp/keystore
        GetFile crypto-config/peerOrganizations/orgb.com/users/Admin@orgb.com/msp/keystore
        GetFile crypto-config/peerOrganizations/orgb.com/users/User1@orgb.com/msp/keystore
        ;;
    clean)
        rm -rf ./channel-artifacts
        rm -rf ./crypto-config
        rm -rf ./production
        ;;
    up)
        docker-compose -f ./docker-compose-cli.yaml up -d
        docker exec -ti cli /bin/bash -c 'scripts/env.sh all'
        ;;
    down)
        docker kill $(docker ps -aq)
        docker system prune
        ;;
esac
