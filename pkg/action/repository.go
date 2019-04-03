package action

import (
	"errors"
	"github.com/lcserny/go-videosmover/pkg/models"
	"github.com/lcserny/goutils"
	"strings"
	"sync"
)

var (
	actionsMu sync.RWMutex
	actions   = make(map[string]Action)
)

type Action interface {
	Execute([]byte, *models.ActionConfig) (string, error)
}

type unknownAction struct {
}

func (uc *unknownAction) Execute(json []byte, cfg *models.ActionConfig) (string, error) {
	return "", errors.New("Unknown action given")
}

func Register(name string, a Action) {
	actionsMu.Lock()
	defer actionsMu.Unlock()
	if a == nil {
		goutils.LogFatal(errors.New("No Action given to register"))
	}
	actions[strings.ToLower(name)] = a
}

func Get(name string) Action {
	if a, ok := actions[strings.ToLower(name)]; ok {
		return a
	}
	return &unknownAction{}
}
