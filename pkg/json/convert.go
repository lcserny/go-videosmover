package json

import (
	"encoding/json"
	"github.com/lcserny/goutils"
)

func EncodeString(data interface{}) string {
	return string(EncodeBytes(data))
}

func EncodeBytes(data interface{}) []byte {
	resultBytes, err := json.Marshal(data)
	goutils.LogError(err)
	return resultBytes
}
