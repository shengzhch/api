#!/bin/bash

protoc -I/usr/local/include -I. -I$GOPATH/src/ -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
--go_out=plugins=grpc,paths=source_relative:. dto/*.proto

protoc -I/usr/local/include -I. -I$GOPATH/src -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
--swagger_out=logtostderr=true:. dto/*.proto

