package actions

// TODO: add action to use from qBittorrent when done downloading to add to a db or something,
//  then in Android app on startup it can maybe show you last finished downloading torrents
//  For this to work, enable qBittorrent WebUI (bypass localhost login and set port) to issue delete requests from torrent list on download complete

// Example delete torrent
/*POST http://localhost:8078/command/delete
Content-Type: application/x-www-form-urlencoded
hashes=e60bc7149e5c0e9b32f60cefb7c63bad303ceca6*/
func RemoveTorrentAction(jsonPayload []byte, config *ActionConfig) (string, error) {
	return "", nil
}
