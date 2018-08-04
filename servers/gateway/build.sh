#! /usr/bin/env bash
set -e
CGO_ENABLED=0 GOOS=linux go build -a
docker build -t alexsirr/telestrations-gateway .
go clean
