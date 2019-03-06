package actions

import (
	"fmt"
	. "github.com/lcserny/goutils"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

var cachedVideoHeader []byte

func setupTmpDir(t *testing.T, prefix string) string {
	dir, err := ioutil.TempDir("", prefix)
	if err != nil {
		t.Fatalf("Couldn't create tmpDir,  %+v", err)
	}

	fmt.Println("Created path", dir)
	return dir
}

func addVideo(t *testing.T, path, video string) string {
	videoPath := filepath.Join(path, video)

	err := os.MkdirAll(filepath.Dir(videoPath), os.ModePerm)
	if err != nil {
		t.Fatalf("Couldn't create folders of file %s, %+v", videoPath, err)
	}

	videoFile, err := os.Create(videoPath)
	if err != nil {
		t.Fatalf("Couldn't create video file %s, %+v", videoPath, err)
	}
	defer CloseFile(videoFile)

	_, err = fmt.Fprint(videoFile, getVideoContent(), os.ModePerm)
	if err != nil {
		t.Fatalf("Couldn't write video content to file %s, %+v", videoPath, err)
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
	if cachedVideoHeader == nil {
		headerContentPath := "../../../files/header.mp4"
		bytes, _ := ioutil.ReadFile(headerContentPath)
		cachedVideoHeader = bytes
		fmt.Println("Loaded video header content from", headerContentPath)
	}
	return cachedVideoHeader
}

func cleanupTmpDir(t *testing.T, path string) {
	fmt.Println("removing path", path)
	err := os.RemoveAll(path)
	if err != nil {
		t.Fatalf("Couldn't remove tmpFiles: %s, %+v", path, err)
	}
}
