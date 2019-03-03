package shared

const (
	MOVIE = "movie"
	TV    = "tv"

	ORIGIN_NAME = "NAME"
	ORIGIN_DISK = "DISK"
	ORIGIN_TMDB = "TMDB"
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
	Names  []string `json:"names"`
	Origin string   `json:"origin"`
}

type MoveRequestData struct {
	Video    string   `json:"video"`
	Subs     []string `json:"subs"`
	Type     string   `json:"type"`
	DiskPath string   `json:"diskPath"`
	OutName  string   `json:"outName"`
}

type MoveResponseData struct {
	UnmovedFolder string `json:"unmovedFolder"`
	Reason        string `json:"reason"`
}
