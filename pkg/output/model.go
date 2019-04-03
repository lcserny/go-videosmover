package output

const (
	ORIGIN_NAME       = "NAME"
	ORIGIN_DISK       = "DISK"
	ORIGIN_TMDB       = "TMDB"
	ORIGIN_TMDB_CACHE = "TMDB_CACHE"
)

type RequestData struct {
	Name             string `json:"name"`
	Type             string `json:"type"`
	SkipCache        bool   `json:"skipCache"`
	SkipOnlineSearch bool   `json:"skipOnlineSearch"`
	DiskPath         string `json:"diskPath"`
}

type ResponseData struct {
	Names  []string `json:"names"`
	Origin string   `json:"origin"`
}
