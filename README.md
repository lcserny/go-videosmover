# Proxy Server
#### Description
How to run Server:  
- make a config from example one for your OS  
- from root of project `go build cmd/proxy_server.go` NOTE (windows): add `-ldflags -H=windowsgui`  
- launch server providing configsFolder path  

JSON requests should be like (regardless of GET or POST usage):  
{  
"action": action to give the executable (like `search`)  
"payload" : json data that will be serialized to tmp file, the path will be sent to jar/bin (file will be deleted after request is done)    
}  

JSON response should be like:  
{  
"code": 200 or 500 or whatever  
"error" : "" if error not empty (errors might be, jar not found or such)  
"date": "2019-02-20 20:15:85"  
"body" : some JSON body response or empty if error  
}  

# Webview
#### Description
How to run:
- make a config from example one for your OS  
- from root of project `go build cmd/webview.go` NOTE (windows): add `-ldflags -H=windowsgui`  
- launch server providing configsFolder path  

TODO:  
- improve tmdb results select UI (add poster, cast etc to dropdown select)
- improve UI in general (add CSS animations and etc.)
- separate statics better (move from base to search.js and such)

# Commander
#### Description
How to run:      
- from root of project `go build cmd/commander.go` NOTE (windows): add `-ldflags -H=windowsgui`  
- for the TMDB online API search to work you need to set `TMDB_API_KEY` environment variable  
- execute actions on built app  

# Remove Torrents
#### Description
How to run:  
- from root of project `go build cmd/remove_from_qtorrent.go` NOTE (windows): add `-ldflags -H=windowsgui`  
- execute actions on built app  

#### qBittorrent notes:    
- enable WebUI, set its port and config it to bypass localhost credentials  
- on torrent completion use following command: `/path/to/remove_from_qtorrent -port=portNumber -hash="%I"`  
