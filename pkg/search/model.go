package search

type RequestData struct {
	Path string `json:"path"`
}

type ResponseData struct {
	Path      string   `json:"path"`
	Subtitles []string `json:"subtitles"`
}
