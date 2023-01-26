#!/bin/sh

/home/leonardo/bin/videosmover/cache_server -config=/home/leonardo/bin/videosmover/configs/cache_config.json &
/home/leonardo/bin/videosmover/proxy_server -config=/home/leonardo/bin/videosmover/configs/proxy_config.json &
/home/leonardo/bin/videosmover/webview -config=/home/leonardo/bin/videosmover/configs/webview_config.json &
