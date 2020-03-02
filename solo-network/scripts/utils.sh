#!/bin/bash
case $1 in
    clean)
        rm -rf ./channel-artifacts
        rm -rf ./crypto-config
        rm -rf ./production
        ;;
    up)
        docker-compose -f ./docker-compose-cli.yaml up -d
        docker exec -ti cli /bin/bash
        ;;
    down)
        docker kill $(docker ps -aq)
        docker system prune
        ;;
esac
