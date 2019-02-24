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
	return "", errors.New("Unknown command given")
}

func NewActionFrom(val string) Action {
	switch strings.ToUpper(val) {
	case SEARCH:
		return new(SearchAction)
	}
	return new(UnknownAction)
}
