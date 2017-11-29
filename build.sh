#!/bin/bash

go get ./...
go build -ldflags "-linkmode external -extldflags -static" -o app
docker build -t rem/dataservice:$1 .

