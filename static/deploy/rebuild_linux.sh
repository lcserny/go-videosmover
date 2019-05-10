#!/bin/bash

cd ../..

go test -v videosmover/pkg/action
go test -v videosmover/pkg/move
go test -v videosmover/pkg/delete
go test -v videosmover/pkg/search
go test -v videosmover/pkg/output

go install cmd/commander.go
go install cmd/proxy_server.go
go install cmd/webview.go
go install cmd/remove_from_qtorrent.go
