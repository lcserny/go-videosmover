package videosmover

import (
	"github.com/juju/errors"
	"strings"
)

const (
	SEARCH = "SEARCH"
)

type Action interface {
	Execute(jsonFile string) (string, error)
}

type UnknownAction struct {
}

func (a *UnknownAction) Execute(jsonFile string) (string, error) {
	return "", errors.New("Unknown action given")
}

// TODO: add action to use from qBittorrent when done downloading to add to a db or something,
//  then in Android app on startup it can maybe show you last finished downloading torrents
func NewActionFrom(val string) Action {
	switch strings.ToUpper(val) {
	case SEARCH:
		return new(SearchAction)
	}
	return new(UnknownAction)
}
