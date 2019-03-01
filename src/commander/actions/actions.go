package actions

import (
	"github.com/pkg/errors"
	"strings"
)

// TODO: add action to use from qBittorrent when done downloading to add to a db or something,
//  then in Android app on startup it can maybe show you last finished downloading torrents
var actionsMap = map[string]Action{
	"SEARCH": SearchAction,
	"OUTPUT": OutputAction,
}

type Action func(jsonPayload []byte) (string, error)

func UnknownAction(jsonPayload []byte) (string, error) {
	return "", errors.New("Unknown action given")
}

func NewActionFrom(val string) Action {
	if action, ok := actionsMap[strings.ToUpper(val)]; ok {
		return action
	}
	return UnknownAction
}
