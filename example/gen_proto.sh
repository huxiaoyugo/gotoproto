#!/usr/bin/env bash

TO_FILE=rpc.proto
SRC_FILE=./proto/proto.go,./proto/proto2.go
PACKAGE_NAME=pack

echo '开始'

./gene_proto -toFileSrc $TO_FILE -src $SRC_FILE -packageName $PACKAGE_NAME

echo '完成'