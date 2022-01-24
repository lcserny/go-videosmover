#!/bin/bash

if ! pgrep -x "cache_server" 
then
    $(pwd)/cache_server -config=configs/cache_config.json &
    disown
fi

if ! pgrep -x "proxy_server" 
then
    $(pwd)/proxy_server -config=configs/proxy_config.json &
    disown
fi

if ! pgrep -x "webview" 
then
    $(pwd)/webview -config=configs/webview_config.json &
    disown
fi
