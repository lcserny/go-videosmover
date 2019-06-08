# ![Alt text](static/html/img/application_small.png?raw=true) Videos Mover (v3)
## Description

VideosMover is a lightweight cross-platform desktop-webapp helper which can make it easy to play videos in Kodi media player as it can find your downloaded videos, get video descriptions from *TMDB* (https://www.themoviedb.org/), prepare folders based on type and move them accordingly. *etcd* (https://etcd.io/) client implementation is provided for caching online results. Future versions will have more capabilities to make your video library management even easier, for example, automatically showing you which videos were viewed already and such.  

Light mode screenshot:
![Alt text](static/screens/screen1.jpg?raw=true)

Dark mode screenshot:
![Alt text](static/screens/screen2.jpg?raw=true)

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

Deployment script provided in `static/deploy/rebuild_<OS>.sh`  
Tests provided for each commander action (core of app)  

### Proxy Server
#### How to run Server:  
1. make a config (minimal example provided), options include:  
    * `logFile` path to logfile (default: vm-proxyserver.log)
    * `port` port of server (default: 8077)
    * `bin` list of objects containing uri, path and cfgPath (see example)
2. from root of project `go install cmd/proxy_server.go` NOTE (windows): add `-ldflags="-H windowsgui"`  
3. launch server providing config path  

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
1. make a config (minimal example provided), options include:
    * `logFile` path to logfile (default: vm-webview.log)
    * `port` port of server (default: 8079)
    * `autoDarkModeEnable` enable auto dark theme mode (defaut: true)
    * `autoDarkModeHourStart` hour after which to enable dark theme (default: 18)
    * `autoDarkModeHourEnd` hour until which to enable dark theme (default: 6)
    * `serverPingTimeoutMs` timeout to server close if ping failed (default: 10000 -> 10sec)
    * `videosMoverAPI` address of videos mover API (default: http://localhost:8077/exec-bin/videos-mover)
    * `htmlFilesPath` path to html files and templates
    * `downloadsPath` path to downloaded videos
    * `moviesPath` path to move movies to
    * `tvSeriesPath` path to move tv shows to
    
2. from root of project `go install cmd/webview.go` NOTE (windows): add `-ldflags="-H windowsgui"`  
3. launch server providing config path  

#### TODO:  
- add custom logo  
- improve UI in general (improve html responsiveness, add CSS animations, hover move button and etc.)  

### Commander
#### How to run:      
1. make a config (configured example provided), options include:
    * `logFile` path to logfile (default: vm-commander.log)
    * `cacheAddress` address of cache server (etcd) (default: http://127.0.0.1:2379)
    * `cachePoolSize` size of cache connection pool -> not used by etcd (default: 10)
    * `minimumVideoSize` minimum byte file size to be considered proper video and not trailer or such (default: 52428800 -> 50mb)
    * `similarityPercent` percent used to determine if video exists on disk -> Levenshtein algorithm used (default: 80)
    * `maxOutputWalkDepth` how many child folder levels to scan for videos when executing output action find on disk (default: 2)
    * `maxSearchWalkDepth` how many child folder levels to scan for videos when executing search action find on disk (default: 4)
    * `maxWebSearchResultCount` max results to obtain from online metadata search -> TMDB (default: 10)
    * `headerBytesSize` size of header in bytes, used to determine file types (default: 261)
    * `restrictedRemovePaths` list of paths that are restricted from clean / removal operations
    * `nameTrimRegexes` list of regexes to use when cleaning the name of videos downloaded
    * `searchExcludePaths` list of folder names in downloads path to exclude from searching 
    * `allowedMIMETypes` list of MIME types considered proper videos
    * `allowedSubtitleExtensions` list of subtitle extensions used when deciding what subs to move with the videos
2. from root of project `go install cmd/commander.go` NOTE (windows): add `-ldflags="-H windowsgui"`  
3. *OPTIONAL*: for the TMDB online API search to work you need to set `TMDB_API_KEY` environment variable  
4. execute actions on built app  

### Remove Torrents
#### How to run:  
1. from root of project `go install cmd/remove_from_qtorrent.go` NOTE (windows): add `-ldflags="-H windowsgui"`  
2. execute actions on built app  

#### qBittorrent required settings:    
- enable WebUI, set its port and config it to bypass localhost credentials  
- on torrent download completion use the following command: `/path/to/remove_from_qtorrent -port=portNumber -hash="%I"`  
