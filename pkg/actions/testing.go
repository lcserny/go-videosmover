package actions

import (
	"encoding/json"
	"fmt"
	. "github.com/lcserny/goutils"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

type testActionData struct {
	request  interface{}
	response interface{}
}

var (
	cachedMP4VideoHeader   []byte
	testSearchActionConfig *ActionConfig
)

func getTestActionConfig() *ActionConfig {
	if testSearchActionConfig == nil {
		testSearchActionConfig = GenerateActionConfig("../../cfg", "commander-actions.test.json")
	}
	return testSearchActionConfig
}

func runActionTest(t *testing.T, items []testActionData, action Action) {
	config := getTestActionConfig()
	for _, td := range items {
		reqBytes := getJSONBytesForTest(t, td.request)
		resString := getJSONStringForTest(t, td.response)

		result, err := action(reqBytes, config)
		if err != nil {
			t.Fatalf("Error occurred: %+v", err)
		}

		if result != resString {
			t.Errorf("Result mismatch:\nwanted %s\ngot: %s", resString, result)
		}
	}
}

func setupTmpDir(t *testing.T, prefix string) (string, func()) {
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

func addVideo(t *testing.T, path, video string) string {
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

func addFile(t *testing.T, path, file string, size int64) string {
	filePath := filepath.Join(path, file)

	err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
	if err != nil {
		t.Fatalf("Couldn't create folders of file %s, %+v", filePath, err)
	}

	newFile, err := os.Create(filePath)
	if err != nil {
		t.Fatalf("Couldn't create file %s, %+v", filePath, err)
	}
	defer CloseFile(newFile)

	err = os.Truncate(filePath, size)
	if err != nil {
		t.Fatalf("Couldn't write %d zeroes to file %s , %+v", size, filePath, err)
	}

	fmt.Println("Created file", filePath)
	return filePath
}

func addSubtitles(t *testing.T, videoPath string, subs []string) (subtitles []string) {
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
		CloseFile(subFile)

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

func getJSONBytesForTest(t *testing.T, data interface{}) []byte {
	bytes, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Couldn't decode request: %+v", err)
	}
	return bytes
}

func getJSONStringForTest(t *testing.T, data interface{}) string {
	bytes, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Couldn't decode response: %+v", err)
	}
	return string(bytes)
}
