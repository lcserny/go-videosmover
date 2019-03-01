package shared

type SearchRequestData struct {
	Path string `json:"path"`
}

type SearchResponseData struct {
	Path      string   `json:"path"`
	Subtitles []string `json:"subtitles"`
}

type OutputRequestData struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	MoviesPath string `json:"moviesPath"`
	TvPath     string `json:"tvPath"`
}

type OutputResponseData struct {
	Name string `json:"name"`
}
