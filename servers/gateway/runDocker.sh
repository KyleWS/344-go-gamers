#! /usr/bin/env bash
set +e
# create docker network to run containers
#docker rm -f redissvr
#docker rm -f mongosvr
docker rm -f gateway
docker network rm appnet
docker network create appnet

# redis container
docker run -d \
--name redissvr \
--network appnet \
redis

# mongo container
docker run -d \
--name mongosvr \
--network appnet \
-e MONGO_INITDB_DATABASE=users \
mongo

# gateway container
docker pull alexsirr/telestrations-gateway
docker run -d \
-p 443:443 \
--name gateway \
--network appnet \
-v /etc/letsencrypt:/etc/letsencrypt:ro \
-e TLSCERT=/etc/letsencrypt/live/api.telestrations.alexsirr.me/fullchain.pem \
-e TLSKEY=/etc/letsencrypt/live/api.telestrations.alexsirr.me/privkey.pem \
-e REDISADDR=redissvr:6379 \
-e DBADDR=mongosvr:27017 \
-e SESSIONKEY=$(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 128 | head -n 1) \
--restart always \
alexsirr/telestrations-gateway

echo "Docker run complete"

