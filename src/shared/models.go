package shared

const (
	MOVIE = "movie"
	TV    = "tv"
)

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

type MoveRequestData struct {
	Video    string   `json:"video"`
	Subs     []string `json:"subs"`
	DiskPath string   `json:"diskPath"`
	OutName  string   `json:"outName"`
}

type MoveResponseData struct {
	UnmovedVideo string `json:"unmovedVideo"`
}
