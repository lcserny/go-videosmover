#!/bin/bash

cd ../..
go build -ldflags -H=windowsgui -o=../bin/commander.exe cmd/commander.go
go build -ldflags -H=windowsgui -o=../bin/proxy_server.exe cmd/proxy_server.go
go build -ldflags -H=windowsgui -o=../bin/webview.exe cmd/webview.go
go build -ldflags -H=windowsgui -o=../bin/remove_from_qtorrent.exe cmd/remove_from_qtorrent.go
