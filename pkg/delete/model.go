package delete

type RequestData struct {
	Folder string `json:"folder"`
}

type ResponseData struct {
	UnRemovedFolder string   `json:"unRemovedFolder"`
	Reasons         []string `json:"reasons"`
}
