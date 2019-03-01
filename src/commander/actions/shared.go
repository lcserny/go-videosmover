package actions

import (
	"encoding/json"
	. "github.com/lcserny/goutils"
)

func getJSONEncodedString(data interface{}) string {
	resultBytes, err := json.Marshal(data)
	LogError(err)
	return string(resultBytes)
}
