package generate

import (
	"encoding/json"
	"errors"
	"github.com/lcserny/go-videosmover/pkg/convert"
	"github.com/lcserny/goutils"
	"io/ioutil"
	"path/filepath"
)

func NewServerConfig(configsPath, configFile string) *convert.ProxyServerConfig {
	configBytes, err := ioutil.ReadFile(filepath.Join(configsPath, configFile))
	goutils.LogFatal(err)

	var serverConfig convert.ProxyServerConfig
	err = json.Unmarshal(configBytes, &serverConfig)
	goutils.LogFatal(err)

	if serverConfig.Host == "" || serverConfig.Port == "" {
		goutils.LogFatal(errors.New("No `host` and/or `port` configured"))
	}

	return &serverConfig
}
