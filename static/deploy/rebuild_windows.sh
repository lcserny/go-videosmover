#!/bin/bash

cd ../..

go test -v videosmover/pkg/action
go test -v videosmover/pkg/move
go test -v videosmover/pkg/delete
go test -v videosmover/pkg/search
go test -v videosmover/pkg/output

go install -ldflags -H=windowsgui cmd/commander.go
go install -ldflags -H=windowsgui cmd/proxy_server.go
go install -ldflags -H=windowsgui cmd/webview.go
go install -ldflags -H=windowsgui cmd/remove_from_qtorrent.go
