package actions

import (
	"github.com/pkg/errors"
	"strings"
)

const (
	SEARCH = "SEARCH"
	OUTPUT = "OUTPUT"
)

type Action func(jsonPayload []byte) (string, error)

func UnknownAction(jsonPayload []byte) (string, error) {
	return "", errors.New("Unknown action given")
}

// TODO: add action to use from qBittorrent when done downloading to add to a db or something,
//  then in Android app on startup it can maybe show you last finished downloading torrents
func NewActionFrom(val string) Action {
	switch strings.ToUpper(val) {
	case SEARCH:
		return SearchAction
	case OUTPUT:
		return OutputAction
	}
	return UnknownAction
}
