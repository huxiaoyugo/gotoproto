# gotoproto
Go model switch to proto3 file

#### 使用方法

1、下载安装
go get -v github.com/huxiaoyugo/gotoproto

在$GOPATH/bin/下面会生成gotoproto

2、使用gotoproto进行转化

脚本：

TO_FILE=rpc.proto ## 生成的文件

SRC_FILE=./proto/proto.go,./proto/proto2.go ## go源文件

PACKAGE_NAME=pack ## rpc包名

echo '开始'

gotoproto -toFileSrc $TO_FILE -src $SRC_FILE -packageName $PACKAGE_NAME

echo '完成'