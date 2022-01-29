#!/bin/bash

cd ../..

echo "Running tests"

go test videosmover/pkg/action
go test videosmover/pkg/move
go test videosmover/pkg/delete
go test videosmover/pkg/search
go test videosmover/pkg/output

echo "Building apps in output folder"

cd cmd/commander && go build -o ../../static/deploy/output/commander
cd ../..

cd cmd/proxy_server && go build -o ../../static/deploy/output/proxy_server
cd ../..

cd cmd/webview && go build -o ../../static/deploy/output/webview
cd ../..

cd cmd/cache_server && go build -o ../../static/deploy/output/cache_server
cd ../..

cd cmd/remove_from_qtorrent 
env GOOS=windows GOARCH=amd64 go build -o ../../static/deploy/output/qtorrent_remove.exe .
cd ../..

echo "Transpiling javascript"

npm run build

echo "Done!"