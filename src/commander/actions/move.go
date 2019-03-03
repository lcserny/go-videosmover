package actions

import (
	"encoding/json"
	. "github.com/lcserny/go-videosmover/src/shared"
	. "github.com/lcserny/goutils"
)

func MoveAction(jsonPayload []byte) (string, error) {
	var request []MoveRequestData
	err := json.Unmarshal(jsonPayload, &request)
	LogError(err)
	if err != nil {
		return "", err
	}

	var resultList []MoveResponseData

	

	return getJSONEncodedString(resultList), nil
}
