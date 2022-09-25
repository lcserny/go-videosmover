package json

import (
	"encoding/json"
	"github.com/lcserny/goutils"
	"videosmover/pkg"
)

type jsonCodec struct {
	contentType string
}

func NewJsonCodec() core.Codec {
	return &jsonCodec{contentType: "application/json"}
}

func (jc jsonCodec) EncodeString(data interface{}) (string, error) {
	bytes, e := jc.EncodeBytes(data)
	return string(bytes), e
}

func (jc jsonCodec) EncodeBytes(data interface{}) ([]byte, error) {
	resultBytes, err := json.Marshal(data)
	goutils.LogError(err)
	return resultBytes, err
}

func (jc jsonCodec) Decode(body []byte, container interface{}) error {
	return json.Unmarshal(body, container)
}

func (jc jsonCodec) ContentType() string {
	return jc.contentType
}
