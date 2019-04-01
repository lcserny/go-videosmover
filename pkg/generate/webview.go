package generate

import (
	"encoding/json"
	"errors"
	"github.com/lcserny/go-videosmover/pkg/convert"
	"github.com/lcserny/goutils"
	"io/ioutil"
	"path/filepath"
)

func NewWebviewConfig(configsPath, configFile string) *convert.WebviewConfig {
	configBytes, err := ioutil.ReadFile(filepath.Join(configsPath, configFile))
	goutils.LogFatal(err)

	var config convert.WebviewConfig
	err = json.Unmarshal(configBytes, &config)
	goutils.LogFatal(err)

	if config.Host == "" || config.Port == "" {
		goutils.LogFatal(errors.New("No `host` and/or `port` configured"))
	}

	return &config
}
