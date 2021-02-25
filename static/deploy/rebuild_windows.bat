@echo off

cd ..
cd ..

echo Running tests

go test videosmover/pkg/action
go test videosmover/pkg/move
go test videosmover/pkg/delete
go test videosmover/pkg/search
go test videosmover/pkg/output

echo Building apps in output folder

cd cmd/commander
go build -ldflags="-H windowsgui" -o ../../static/deploy/output/commander.exe
cd ..
cd ..

cd cmd/proxy_server
go build -ldflags="-H windowsgui" -o ../../static/deploy/output/proxy_server.exe
cd ..
cd ..

cd cmd/webview
go build -ldflags="-H windowsgui" -o ../../static/deploy/output/webview.exe
cd ..
cd ..

cd cmd/remove_from_qtorrent
go build -ldflags="-H windowsgui" -o ../../static/deploy/output/remove_from_qtorrent.exe
cd ..
cd ..

cd cmd/cache_server
go build -ldflags="-H windowsgui" -o ../../static/deploy/output/cache_server.exe
cd ..
cd ..

cd static/deploy

echo Done!