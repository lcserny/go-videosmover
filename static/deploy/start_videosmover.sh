#!/bin/sh

cd ~/bin/videosmover

./cache_server -config=configs/cache_config.json &
./proxy_server -config=configs/proxy_config.json &
./webview -config=configs/webview_config.json &

