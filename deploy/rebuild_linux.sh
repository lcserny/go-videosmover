#!/bin/bash

cd ..
go build -o=../bin/commander cmd/commander.go
go build -o=../bin/proxy_server cmd/proxy_server.go
go build -o=../bin/webview cmd/webview.go
go build -o=../bin/remove_from_qtorrent cmd/remove_from_qtorrent.go
