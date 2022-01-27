if ! pgrep -x "cache_server" >/dev/null
then
    /home/leonardo/bin/videosmover/cache_server -config=/home/leonardo/bin/videosmover/configs/cache_config.json &
fi

if ! pgrep -x "proxy_server" >/dev/null
then
    /home/leonardo/bin/videosmover/proxy_server -config=/home/leonardo/bin/videosmover/configs/proxy_config.json &
fi

if ! pgrep -x "webview" >/dev/null
then
    /home/leonardo/bin/videosmover/webview -config=/home/leonardo/bin/videosmover/configs/webview_config.json &
fi
