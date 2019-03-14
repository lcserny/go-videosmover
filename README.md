# Server
#### Description
How to run Server:  
- get https://github.com/gobuffalo/packr  
- make a config from example one for your OS and put it in `cfg/server` folder  
- from root of project `packr build cmd/server.go`  
- launch server

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

v1: 
 - no UI needed, (Ui will be done on requester side), only possibility to listen for requests and execute jar
 (this might come from an Android app, initially)

v2:
 - add Controller/s for UI rendering, use that as entry point
 - add support for html UI (so the requester doesn't need to provide Ui rendering)

# Commander
#### Description
How to run:  
- get https://github.com/gobuffalo/packr    
- from root of project `packr build cmd/commander.go`  
- for the TMDB online API search to work you need to set `TMDB_API_KEY` environment variable  
- execute actions on built app  

# Remove Torrents
#### Description
How to run:  
- get https://github.com/gobuffalo/packr  
- from root of project `packr build cmd/remove_torrents.go`  
- execute actions on built app  
