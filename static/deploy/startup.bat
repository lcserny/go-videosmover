CONSOLESTATE /Hide
start "" "%~dp0bin\cache_server.exe" -config="%~dp0configs\cache_config.json"
start "" "%~dp0bin\proxy_server.exe" -config="%~dp0configs\proxy_config.json"
