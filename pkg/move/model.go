package move

type RequestData struct {
	Video    string   `json:"video"`
	Subs     []string `json:"subs"`
	Type     string   `json:"type"`
	DiskPath string   `json:"diskPath"`
	OutName  string   `json:"outName"`
}

type ResponseData struct {
	UnmovedFolder string   `json:"unmovedFolder"`
	Reasons       []string `json:"reasons"`
}
