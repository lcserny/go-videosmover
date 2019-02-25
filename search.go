package videosmover

import (
	"encoding/json"
	. "github.com/lcserny/goutils"
	"io/ioutil"
	"os"
	"path/filepath"
)

type SearchAction struct {
}

// TODO: put these in a shared go project `go-videosmover-shared`
type RequestSearchData struct {
	Path string `json:"path"`
}

type ResponseSearchData struct {
	Path string `json:"path"`
}

func (a *SearchAction) Execute(jsonFile string) (string, error) {
	jsonRequestBytes, err := ioutil.ReadFile(jsonFile)
	LogError(err)
	if err != nil {
		return "", err
	}

	var request RequestSearchData
	err = json.Unmarshal(jsonRequestBytes, &request)
	LogError(err)
	if err != nil {
		return "", err
	}

	var resultList []ResponseSearchData
	err = filepath.Walk(request.Path, func(path string, info os.FileInfo, err error) error {
		LogError(err)
		if !info.IsDir() {
			resultList = append(resultList, ResponseSearchData{path})
		}
		return nil
	})
	LogError(err)
	if err != nil {
		return "", err
	}

	resultBytes, err := json.Marshal(resultList)
	LogError(err)
	if err != nil {
		return "", err
	}

	return string(resultBytes), nil
}
