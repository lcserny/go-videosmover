#!/bin/bash

cd ../..

go test videosmover/pkg/action
go test videosmover/pkg/move
go test videosmover/pkg/delete
go test videosmover/pkg/search
go test videosmover/pkg/output

cd cmd/commander && go install
cd ../..

cd cmd/proxy_server && go install
cd ../..

cd cmd/webview && go install
cd ../..

cd cmd/remove_from_qtorrent && go install
cd ../..