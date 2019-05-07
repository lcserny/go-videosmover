package json

import (
	"encoding/json"
	"github.com/lcserny/goutils"
)

func EncodeString(data interface{}) (string, error) {
	bytes, e := EncodeBytes(data)
	return string(bytes), e
}

func EncodeBytes(data interface{}) ([]byte, error) {
	resultBytes, err := json.Marshal(data)
	goutils.LogError(err)
	return resultBytes, err
}

func Decode(body []byte, container interface{}) error {
	return json.Unmarshal(body, container)
}
