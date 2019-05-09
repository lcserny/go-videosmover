package action

import (
	"errors"
	"fmt"
	"github.com/lcserny/goutils"
	"strings"
	"sync"
	core "videosmover/pkg"
)

type actionRepository struct {
	mutex   sync.RWMutex
	actions map[string]core.Action
}

func NewActionRepository() core.ActionRepository {
	return &actionRepository{}
}

type unknownAction struct {
}

func (uc *unknownAction) Execute(json []byte) (string, error) {
	return "", errors.New("unknown action given")
}

func (ar *actionRepository) Register(key string, a core.Action) {
	ar.mutex.Lock()
	defer ar.mutex.Unlock()
	if a == nil {
		goutils.LogFatal(errors.New("no Action given to register"))
	}
	if _, dup := ar.actions[key]; dup {
		goutils.LogFatal(errors.New(fmt.Sprintf("action `%s` already defined", key)))
	}
	ar.actions[strings.ToLower(key)] = a
}

func (ar actionRepository) Retrieve(key string) core.Action {
	if a, ok := ar.actions[strings.ToLower(key)]; ok {
		return a
	}
	return &unknownAction{}
}
