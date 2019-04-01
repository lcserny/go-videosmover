package testing

import (
	"fmt"
	"github.com/lcserny/go-videosmover/pkg/action"
	"github.com/lcserny/go-videosmover/pkg/convert"
	"github.com/lcserny/go-videosmover/pkg/generate"
	. "github.com/lcserny/goutils"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

type TestActionData struct {
	Request  interface{}
	Response interface{}
}

var cachedMP4VideoHeader []byte

func RunActionTest(t *testing.T, items []TestActionData, action action.Action) {
	config := generate.NewTestActionConfig()
	for _, td := range items {
		reqBytes := convert.GetJSONBytesForTest(t, td.Request)
		resString := convert.GetJSONStringForTest(t, td.Response)

		result, err := action(reqBytes, config)
		if err != nil {
			t.Fatalf("Error occurred: %+v", err)
		}

		if result != resString {
			t.Errorf("Result mismatch:\nwanted %s\ngot: %s", resString, result)
		}
	}
}

func SetupTmpDir(t *testing.T, prefix string) (string, func()) {
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

func AddVideo(t *testing.T, path, video string) string {
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

func AddFile(t *testing.T, path, file string, size int64) string {
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

func AddSubtitles(t *testing.T, videoPath string, subs []string) (subtitles []string) {
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
