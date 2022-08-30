#!/bin/bash

DEF_EXT=""
PASSED_OS="linux"
if [ "$1" ]
  then
    PASSED_OS=$1
    if [ $PASSED_OS = "windows" ]
      then
        DEF_EXT=".exe"
    fi
fi

PASSED_ARCH="amd64"
if [ "$2" ]
  then
    PASSED_ARCH=$2
fi

cd ../..

echo "Running tests"

go test videosmover/pkg/action
go test videosmover/pkg/move
go test videosmover/pkg/delete
go test videosmover/pkg/search
go test videosmover/pkg/output

echo "Building apps in output folder"

cd cmd/commander 
env GOOS=$PASSED_OS GOARCH=$PASSED_ARCH go build -o ../../static/deploy/output/commander$DEF_EXT
cd ../..

cd cmd/proxy_server 
env GOOS=$PASSED_OS GOARCH=$PASSED_ARCH go build -o ../../static/deploy/output/proxy_server$DEF_EXT
cd ../..

cd cmd/webview 
env GOOS=$PASSED_OS GOARCH=$PASSED_ARCH go build -o ../../static/deploy/output/webview$DEF_EXT
cd ../..

cd cmd/cache_server 
env GOOS=$PASSED_OS GOARCH=$PASSED_ARCH go build -o ../../static/deploy/output/cache_server$DEF_EXT
cd ../..

cd cmd/remove_from_qtorrent 
env GOOS=$PASSED_OS GOARCH=$PASSED_ARCH go build -o ../../static/deploy/output/qtorrent_remove$DEF_EXT
cd ../..

echo "Transpiling javascript"

npm run build

echo "Done!"
