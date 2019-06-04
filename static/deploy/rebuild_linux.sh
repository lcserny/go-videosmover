#!/bin/bash

cd ../..

go test videosmover/pkg/action
go test videosmover/pkg/move
go test videosmover/pkg/delete
go test videosmover/pkg/search
go test videosmover/pkg/output

go install cmd/cachestore.go
go install cmd/commander.go
go install cmd/proxy_server.go
go install cmd/webview.go
go install cmd/remove_from_qtorrent.go
