package actions

import (
	"fmt"
	. "github.com/lcserny/goutils"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

var cachedMP4VideoHeader []byte

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

func addSubtitles(t *testing.T, path string, subs []string) (subtitles []string) {
	for _, sub := range subs {
		subPath := filepath.Join(path, sub)

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
		headerContentPath := "../../../files/mpeg4.mp4"
		bytes, _ := ioutil.ReadFile(headerContentPath)
		cachedMP4VideoHeader = bytes
		fmt.Println("Loaded video header content from", headerContentPath)
	}
	return cachedMP4VideoHeader
}
