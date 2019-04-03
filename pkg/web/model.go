package web

type RequestJsonData struct {
	Action  string      `json:"action"`
	Payload interface{} `json:"payload,string"`
}

type ResponseJsonData struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
	Date  string `json:"date"`
	Body  string `json:"body"`
}
