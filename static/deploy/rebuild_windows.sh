#!/bin/bash

GPATH=$1;
CDIR=$(pwd);

if [[ $# -eq 0 ]];
then
  echo "Please provide GOPATH path as first arg to script";
  exit 1;
fi

if [[ $CDIR != *static/deploy ]];
then
  echo "Please run from inside #REPO_ROOT/static/deploy folder";
  exit 1;
fi

cd ../..

echo "Running tests ...";

go test videosmover/pkg/action
go test videosmover/pkg/move
go test videosmover/pkg/delete
go test videosmover/pkg/search
go test videosmover/pkg/output

echo "Installing binaries in '$GPATH' ...";

cd cmd/commander && go install -ldflags="-H windowsgui"
cd ../..

cd cmd/proxy_server && go install -ldflags="-H windowsgui"
cd ../..

cd cmd/webview && go install -ldflags="-H windowsgui"
cd ../..

cd cmd/remove_from_qtorrent && go install -ldflags="-H windowsgui"
cd ../..

cd cmd/cache_server && go install -ldflags="-H windowsgui"
cd ../..

cd $CDIR;

echo "Cleaning old output ...";

rm -rf out;

echo "Exporting to 'out' directory ...";

mkdir -p out/bin;
mkdir -p out/cache;
mkdir -p out/configs;
mkdir -p out/html;

cp $GPATH/bin/cache_server.exe out/bin;
cp $GPATH/bin/commander.exe out/bin;
cp $GPATH/bin/proxy_server.exe out/bin;
cp $GPATH/bin/remove_from_qtorrent.exe out/bin;
cp $GPATH/bin/webview.exe out/bin;

cp -R ../html/* out/html;

cp ../configs/cache_config_EXAMPLE.json out/configs/cache_config.json;
cp ../configs/commander_actions_EXAMPLE.json out/configs/commander_actions.json;
cp ../configs/proxy_config_EXAMPLE.json out/configs/proxy_config.json;
cp ../configs/webview_config_EXAMPLE.json out/configs/webview_config.json;

echo "CONSOLESTATE /Hide
start \"\" \"%~dp0bin\cache_server.exe\" -config=\"%~dp0configs\cache_config.json\"
start \"\" \"%~dp0bin\proxy_server.exe\" -config=\"%~dp0configs\proxy_config.json\"" >> out/startup.bat

echo "CONSOLESTATE /Hide
start \"\" \"%~dp0bin\webview.exe\" -config=\"%~dp0configs\webview_config.json\"" >> out/webview.bat

chmod -R +w out/configs;

echo "Exported! Please update out/configs files manually and create shortcuts to bat files provided."
exit 0;