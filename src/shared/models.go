package shared

type SearchRequestData struct {
	Path string `json:"path"`
}

type SearchResponseData struct {
	Path      string   `json:"path"`
	Subtitles []string `json:"subtitles"`
}

type OutputRequestData struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	DiskPath string `json:"diskPath"`
}

type OutputResponseData struct {
	Name string `json:"name"`
}
