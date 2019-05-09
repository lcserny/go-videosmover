# Videos Mover (v3)
## Description

VideosMover is a helper project which can make it easy to play videos in Kodi media player as it can find your downloaded videos, get video descriptions from TMDB (https://www.themoviedb.org/), prepare folders based on type and move them accordingly. Future versions will have more capabilities to make your video library management even easier, for example, automatically showing you which videos were viewed already and such.  

This project contains 4 applications:
- commander
- proxy_server
- webview
- remove_torrents

#### How to use:  
**Remove_torrents** app is used as an execute action after each torrent download complete (see below for details). It will just remove the finished download from the torrent list so it can be moved.  
The rest are used together, **proxy_server** is meant to run always, so you can add it on startup, it is a small webserver (1mb RAM usage and basically no CPU). It is a proxy from an external tool (can be an Android app or such) to your PC.  
The proxy will execute command-line actions on the **commander** which actually does the jobs, like searching videos, moving them and so on.  
The **webview** is the UI part and basically just opens a tab in your default browser which issues requests on the proxy as needed. Closing the tab will close the webview app as it listens to a "pulse" every second.  

### Proxy Server
#### How to run Server:  
- make a config from example  
- from root of project `go build cmd/proxy_server.go` NOTE (windows): add `-ldflags -H=windowsgui`  
- launch server providing configsFolder path  

JSON requests should be like (regardless of GET or POST usage):  
```
{  
  "action": action to give the executable (like "search")  
  "payload" : json data that will be serialized to tmp file, 
      the path will be sent to jar/bin (file will be deleted after request is done)    
}
```

JSON response should be like:
```
{  
  "code": 200 or 500 or whatever  
  "error" : "" if error not empty (errors might be, jar not found or such)  
  "date": "2019-02-20 20:15:85"  
  "body" : some JSON body response or empty if error  
}
```   

### Webview
#### How to run:
- make a config from example  
- from root of project `go build cmd/webview.go` NOTE (windows): add `-ldflags -H=windowsgui`  
- launch server providing configsFolder path  

#### TODO:  
- improve UI in general (improve html responsiveness, add CSS animations, hover move button and etc.)
- improve tmdb results select UI (add poster, cast etc to dropdown select)
- switch js to `GopherJS`  

### Commander
#### How to run:      
- from root of project `go build cmd/commander.go` NOTE (windows): add `-ldflags -H=windowsgui`  
- for the TMDB online API search to work you need to set `TMDB_API_KEY` environment variable  
- execute actions on built app  

### Remove Torrents
#### How to run:  
- from root of project `go build cmd/remove_from_qtorrent.go` NOTE (windows): add `-ldflags -H=windowsgui`  
- execute actions on built app  

#### qBittorrent required settings:    
- enable WebUI, set its port and config it to bypass localhost credentials  
- on torrent download completion use the following command: `/path/to/remove_from_qtorrent -port=portNumber -hash="%I"`  
