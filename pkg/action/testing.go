package action

import (
	"fmt"
	"github.com/lcserny/goutils"
	"github.com/ryanbradynd05/go-tmdb"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"testing"
	"videosmover/pkg"
)

// TODO: move this whole file better

type TestActionData struct {
	Request  interface{}
	Response interface{}
}

var (
	cachedMP4VideoHeader []byte
	testActionCfg        *Config
)

func getTestActionConfig() *Config {
	if testActionCfg == nil {
		testActionCfg = &Config{
			MinimumVideoSize: 256, SimilarityPercent: 80, MaxOutputWalkDepth: 2, MaxSearchWalkDepth: 4,
			MaxTMDBResultCount: 1, OutTMDBCacheLimit: 100, HeaderBytesSize: 261,
			NameTrimRegexes:       []string{".[sS](\\d{1,2})([-]?[eE](\\d{1,2}))?", "[\\.\\s][sS][0-9]{1,2}[\\.\\s]?", "1080p", "720p"},
			RestrictedRemovePaths: []string{"Downloads"}, SearchExcludePaths: []string{"Programming Stuff"},
			AllowedSubtitleExtensions: []string{".srt", ".sub"},
			AllowedMIMETypes:          []string{"video/x-matroska", "video/x-msvideo", "video/mp4"},
		}
		if key, exists := os.LookupEnv("TMDB_API_KEY"); exists {
			testActionCfg.TmdbAPI = tmdb.Init(tmdb.Config{key, false, nil})
		}
		for _, pat := range testActionCfg.NameTrimRegexes {
			testActionCfg.CompiledNameTrimRegexes = append(testActionCfg.CompiledNameTrimRegexes,
				regexp.MustCompile(fmt.Sprintf("(?i)(-?%s)", pat)))
		}
	}
	return testActionCfg
}

func RunTestAction(t *testing.T, slice []TestActionData, a Action, c core.Codec) {
	config := getTestActionConfig()
	for _, td := range slice {
		reqBytes, err := c.EncodeBytes(td.Request)
		if err != nil {
			t.Fatalf("Couldn't decode request: %+v", err)
		}
		resString, err := c.EncodeString(td.Response)
		if err != nil {
			t.Fatalf("Couldn't decode response: %+v", err)
		}

		result, err := a.Execute(reqBytes, config)
		if err != nil {
			t.Fatalf("Error occurred: %+v", err)
		}

		if result != resString {
			t.Errorf("Result mismatch:\nwanted %s\ngot: %s", resString, result)
		}
	}
}

func SetupTestTmpDir(t *testing.T, prefix string) (string, func()) {
	dir, err := ioutil.TempDir("", prefix)
	if err != nil {
		t.Fatalf("Couldn't create tmpDir,  %+v", err)
	}

	fmt.Println("Created path", dir)
	return dir, func() {
		fmt.Println("removing path", dir)
		err := os.RemoveAll(dir)
		if err != nil {
			t.Fatalf("Couldn't remove tmpFiles: %s, %+v", dir, err)
		}
	}
}

func AddTestVideo(t *testing.T, path, video string) string {
	videoPath := filepath.Join(path, video)

	err := os.MkdirAll(filepath.Dir(videoPath), os.ModePerm)
	if err != nil {
		t.Fatalf("Couldn't create folders of file %s, %+v", videoPath, err)
	}

	err = ioutil.WriteFile(videoPath, getVideoContent(), os.ModePerm)
	if err != nil {
		t.Fatalf("Couldn't create video file %s, %+v", videoPath, err)
	}

	fmt.Println("Created video", videoPath)
	return videoPath
}

func AddTestFile(t *testing.T, path, file string, size int64) string {
	filePath := filepath.Join(path, file)

	err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
	if err != nil {
		t.Fatalf("Couldn't create folders of file %s, %+v", filePath, err)
	}

	newFile, err := os.Create(filePath)
	if err != nil {
		t.Fatalf("Couldn't create file %s, %+v", filePath, err)
	}
	defer goutils.CloseFile(newFile)

	err = os.Truncate(filePath, size)
	if err != nil {
		t.Fatalf("Couldn't write %d zeroes to file %s , %+v", size, filePath, err)
	}

	fmt.Println("Created file", filePath)
	return filePath
}

func AddTestSubtitles(t *testing.T, videoPath string, subs []string) (subtitles []string) {
	subsDir := filepath.Dir(videoPath)
	for _, sub := range subs {
		subPath := filepath.Join(subsDir, sub)

		err := os.MkdirAll(filepath.Dir(subPath), os.ModePerm)
		if err != nil {
			t.Fatalf("Couldn't create folders of file %s, %+v", subPath, err)
		}

		subFile, err := os.Create(subPath)
		if err != nil {
			t.Fatalf("Couldn't create sub file %s, %+v", subPath, err)
		}
		goutils.CloseFile(subFile)

		fmt.Println("Created subtitle", subPath)
		subtitles = append(subtitles, subPath)
	}
	return subtitles
}

func getVideoContent() []byte {
	if cachedMP4VideoHeader == nil {
		headerContentPath := "../../testdata/mpeg4.mp4"
		bytes, _ := ioutil.ReadFile(headerContentPath)
		cachedMP4VideoHeader = bytes
		fmt.Println("Loaded video header content from", headerContentPath)
	}
	return cachedMP4VideoHeader
}
