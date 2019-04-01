package convert

import (
	"encoding/json"
	"fmt"
	"github.com/lcserny/goutils"
	"github.com/ryanbradynd05/go-tmdb"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"
)

const TIME_FORMAT = "2006-01-02 15:04:05"

type RequestJsonData struct {
	Action  string      `json:"action"`
	Payload interface{} `json:"payload,string"`
}

type ResponseJsonData struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
	Date  string `json:"date"`
	Body  string `json:"body"`
}

type SearchRequestData struct {
	Path string `json:"path"`
}

type SearchResponseData struct {
	Path      string   `json:"path"`
	Subtitles []string `json:"subtitles"`
}

type OutputRequestData struct {
	Name             string `json:"name"`
	Type             string `json:"type"`
	SkipCache        bool   `json:"skipCache"`
	SkipOnlineSearch bool   `json:"skipOnlineSearch"`
	DiskPath         string `json:"diskPath"`
}

type OutputResponseData struct {
	Names  []string `json:"names"`
	Origin string   `json:"origin"`
}

type MoveRequestData struct {
	Video    string   `json:"video"`
	Subs     []string `json:"subs"`
	Type     string   `json:"type"`
	DiskPath string   `json:"diskPath"`
	OutName  string   `json:"outName"`
}

type MoveResponseData struct {
	UnmovedFolder string   `json:"unmovedFolder"`
	Reasons       []string `json:"reasons"`
}

type DeleteRequestData struct {
	Folder string `json:"folder"`
}

type DeleteResponseData struct {
	UnRemovedFolder string   `json:"unRemovedFolder"`
	Reasons         []string `json:"reasons"`
}

type ActionConfig struct {
	MinimumVideoSize          int64    `json:"minimumVideoSize"`
	SimilarityPercent         int      `json:"similarityPercent"`
	MaxOutputWalkDepth        int      `json:"maxOutputWalkDepth"`
	MaxSearchWalkDepth        int      `json:"maxSearchWalkDepth"`
	MaxTMDBResultCount        int      `json:"maxTMDBResultCount"`
	OutTMDBCacheLimit         int      `json:"outTMDBCacheLimit"`
	HeaderBytesSize           int      `json:"headerBytesSize"`
	RestrictedRemovePaths     []string `json:"restrictedRemovePaths"`
	NameTrimRegexes           []string `json:"nameTrimRegexes"`
	SearchExcludePaths        []string `json:"searchExcludePaths"`
	AllowedMIMETypes          []string `json:"allowedMIMETypes"`
	AllowedSubtitleExtensions []string `json:"allowedSubtitleExtensions"`

	TmdbAPI                 *tmdb.TMDb
	CompiledNameTrimRegexes []*regexp.Regexp
}

type ProxyServerConfig struct {
	Host                       string `json:"host"`
	Port                       string `json:"port"`
	PathVideosMoverJava        string `json:"path.videosMover.java"`
	PathVideosMoverJavaConfigs string `json:"path.videosMover.java.configs"`
	PathVideosMoverBin         string `json:"path.videosMover.bin"`
	PathVideosMoverBinConfigs  string `json:"path.videosMover.bin.configs"`
}

type WebviewConfig struct {
	Host                string `json:"host"`
	Port                string `json:"port"`
	HtmlFilesPath       string `json:"htmlFilesPath"`
	ServerPingTimeoutMs int64  `json:"serverPingTimeoutMs"`
	VideosMoverAPI      string `json:"videosMoverAPI"`
	DownloadsPath       string `json:"downloadsPath"`
	MoviesPath          string `json:"moviesPath"`
	TvSeriesPath        string `json:"tvSeriesPath"`
}

func GetJSONEncodedString(data interface{}) string {
	resultBytes, err := json.Marshal(data)
	goutils.LogError(err)
	return string(resultBytes)
}

func GetJSONBytesForTest(t *testing.T, data interface{}) []byte {
	bytes, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Couldn't decode request: %+v", err)
	}
	return bytes
}

func GetJSONStringForTest(t *testing.T, data interface{}) string {
	bytes, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Couldn't decode response: %+v", err)
	}
	return string(bytes)
}

func GetJsonResponseFromAsBytes(body, err string) []byte {
	if strings.Contains(body, "ERROR") {
		err = body
		body = ""
	}

	code := 200
	if len(err) > 0 {
		code = 500
	}

	responseJsonData := &ResponseJsonData{
		Code:  code,
		Error: err,
		Date:  time.Now().Format(TIME_FORMAT),
		Body:  body,
	}

	jsonBytes, _ := json.Marshal(responseJsonData)
	return jsonBytes
}

func GetErrorJsonResponseAsBytes(err string) []byte {
	responseJsonData := &ResponseJsonData{
		Code:  500,
		Error: err,
		Date:  time.Now().Format(TIME_FORMAT),
		Body:  "",
	}

	jsonBytes, _ := json.Marshal(responseJsonData)
	return jsonBytes
}

func RemoveTmpStoredJsonPayload(tempJsonFile *os.File) {
	err := tempJsonFile.Close()
	goutils.LogErrorWithMessage(fmt.Sprintf("Couldn't close tmpFile: %s", tempJsonFile.Name()), err)

	err = os.Remove(tempJsonFile.Name())
	goutils.LogErrorWithMessage(fmt.Sprintf("Couldn't remove tmpFile: %s", tempJsonFile.Name()), err)
}

func TmpStoreJsonPayload(jsonData interface{}) *os.File {
	tempFile, err := ioutil.TempFile(os.TempDir(), "vms-")
	goutils.LogErrorWithMessage(fmt.Sprintf("Couldn't create tmpFile: %s", tempFile.Name()), err)

	jsonString, err := json.Marshal(jsonData)
	goutils.LogErrorWithMessage("Couldn't convert data to bytes", err)
	if err == nil {
		_, err = tempFile.Write([]byte(jsonString))
	}
	goutils.LogErrorWithMessage(fmt.Sprintf("Couldn't write to tmpFile: %s", tempFile.Name()), err)

	return tempFile
}
