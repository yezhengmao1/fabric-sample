#!/bin/bash

CPFile() {
    files=$(ls $1)
    echo ${files[0]}
    cd $1
    cp ${files[0]} ./key.pem
    cd -
}

GetAllPkFiles() {
    CPFile crypto-config/peerOrganizations/orga.com/users/Admin@orga.com/msp/keystore
    CPFile crypto-config/peerOrganizations/orga.com/users/User1@orga.com/msp/keystore
    CPFile crypto-config/peerOrganizations/orgb.com/users/Admin@orgb.com/msp/keystore
    CPFile crypto-config/peerOrganizations/orgb.com/users/User1@orgb.com/msp/keystore
}

CleanFiles() {
    rm -rf ./channel-artifacts
    rm -rf ./crypto-config
    rm -rf ./production
}

case $1 in
    pk)
        GetAllPkFiles
        ;;
    clean)
        CleanFiles
        ;;
    up)
        docker-compose -f ./docker-compose-cli.yaml up -d
        docker exec cli /bin/bash -c "scripts/env.sh all"
        ;;
    down)
        docker kill $(docker ps -aq)
        docker system prune
        CleanFiles
        ;;
    networkup)
        GetAllPkFiles
        docker-compose -f ./docker-compose-cli.yaml up -d
        docker exec cli /bin/bash -c "scripts/env.sh all"
        ;;
    networkdown)
        docker kill $(docker ps -qa)
        CleanFiles
        ;;
esac
