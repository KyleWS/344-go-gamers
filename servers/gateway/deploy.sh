#! /usr/bin/env bash
./build.sh
docker push alexsirr/telestrations-gateway
ssh root@192.241.152.251 'bash -s' < runDocker.sh
