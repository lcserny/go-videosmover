package actions

import (
	"github.com/pkg/errors"
	"strings"
)

const (
	SEARCH = "SEARCH"
	OUTPUT = "OUTPUT"
)

type RequestSearchData struct {
	Path string `json:"path"`
}

type Action interface {
	Execute(jsonPayload []byte) (string, error)
}

type UnknownAction struct {
}

func (a *UnknownAction) Execute(jsonPayload []byte) (string, error) {
	return "", errors.New("Unknown action given")
}

// TODO: add action to use from qBittorrent when done downloading to add to a db or something,
//  then in Android app on startup it can maybe show you last finished downloading torrents
func NewActionFrom(val string) Action {
	switch strings.ToUpper(val) {
	case SEARCH:
		return new(SearchAction)
	case OUTPUT:
		return new(OutputAction)
	}
	return new(UnknownAction)
}
