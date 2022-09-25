package core

type Codec interface {
	EncodeString(data interface{}) (string, error)
	EncodeBytes(data interface{}) ([]byte, error)
	Decode(body []byte, container interface{}) error
	ContentType() string
}
