#!/bin/bash

cd ../..

go test videosmover/pkg/action
go test videosmover/pkg/move
go test videosmover/pkg/delete
go test videosmover/pkg/search
go test videosmover/pkg/output

cd cmd/commander && go install -ldflags="-H windowsgui"
cd ../..

cd cmd/proxy_server && go install -ldflags="-H windowsgui"
cd ../..

cd cmd/webview && go install -ldflags="-H windowsgui"
cd ../..

cd cmd/remove_from_qtorrent && go install -ldflags="-H windowsgui"
cd ../..