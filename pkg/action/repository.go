package action

import (
	"errors"
	"fmt"
	"github.com/lcserny/goutils"
	"strings"
	"sync"
)

var (
	actionsMu sync.RWMutex
	actions   = make(map[string]Action)
)

type Action interface {
	Execute([]byte, *Config) (string, error)
}

type unknownAction struct {
}

func (uc *unknownAction) Execute(json []byte, cfg *Config) (string, error) {
	return "", errors.New("unknown action given")
}

func Register(name string, a Action) {
	actionsMu.Lock()
	defer actionsMu.Unlock()
	if a == nil {
		goutils.LogFatal(errors.New("no Action given to register"))
	}
	if _, dup := actions[name]; dup {
		goutils.LogFatal(errors.New(fmt.Sprintf("action `%s` already defined", name)))
	}
	actions[strings.ToLower(name)] = a
}

func Retrieve(name string) Action {
	if a, ok := actions[strings.ToLower(name)]; ok {
		return a
	}
	return &unknownAction{}
}
