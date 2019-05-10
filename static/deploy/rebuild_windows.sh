#!/bin/bash

cd ../..

go test videosmover/pkg/action
go test videosmover/pkg/move
go test videosmover/pkg/delete
go test videosmover/pkg/search
go test videosmover/pkg/output

go install -ldflags -H=windowsgui cmd/commander.go
go install -ldflags -H=windowsgui cmd/proxy_server.go
go install -ldflags -H=windowsgui cmd/webview.go
go install -ldflags -H=windowsgui cmd/remove_from_qtorrent.go
