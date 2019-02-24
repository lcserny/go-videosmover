package videosmover

import (
	"fmt"
	"io/ioutil"
)

type SearchAction struct {
}

func (a *SearchAction) Execute(jsonFile string) (string, error) {
	bytes, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Arg: %s, had this content: %s", jsonFile, string(bytes)), nil
}
