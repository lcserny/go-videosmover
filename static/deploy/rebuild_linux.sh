#!/bin/bash

cd ../..

go test -v videosmover/pkg/move
go test -v videosmover/pkg/delete
go test -v videosmover/pkg/search
go test -v videosmover/pkg/output

go build -o=../bin/commander cmd/commander.go
go build -o=../bin/proxy_server cmd/proxy_server.go
go build -o=../bin/webview cmd/webview.go
go build -o=../bin/remove_from_qtorrent cmd/remove_from_qtorrent.go
