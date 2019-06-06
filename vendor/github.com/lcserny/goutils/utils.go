package goutils

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func MakeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func GetRegexSubgroups(exp *regexp.Regexp, text string) map[string]string {
	match := exp.FindStringSubmatch(text)
	resultMap := make(map[string]string)
	for i, name := range exp.SubexpNames() {
		if i != 0 && name != "" {
			resultMap[name] = match[i]
		}
	}
	return resultMap
}

func GetLinesFromString(content string) (lines []string) {
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

func StringsContain(strings []string, match string) bool {
	for _, ele := range strings {
		if ele == match {
			return true
		}
	}
	return false
}

func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func MaxInt(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func GetBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func GetRealPath(path string) (string, os.FileInfo) {
	simLinkInfo, err := os.Lstat(path)
	LogError(err)

	// TODO: don't understand this if, but it works...
	if simLinkInfo.Mode()&os.ModeSymlink == os.ModeSymlink {
		realPath, err := os.Readlink(path)
		LogError(err)
		path = realPath
	}

	fileInfo, err := os.Stat(path)
	LogError(err)

	if path[len(path)-1:] == string(os.PathSeparator) {
		path = path[:len(path)-1]
	}

	return path, fileInfo
}

func GetAbsCurrentPathOf(path string) string {
	currentPath, _ := filepath.Abs(os.Args[0])
	return filepath.Join(filepath.Dir(currentPath), path)
}

func LevenshteinDistance(source, target string) int {
	d := make([][]int, len(source)+1)
	for i := range d {
		d[i] = make([]int, len(target)+1)
	}
	for i := range d {
		d[i][0] = i
	}
	for j := range d[0] {
		d[0][j] = j
	}
	for j := 1; j <= len(target); j++ {
		for i := 1; i <= len(source); i++ {
			if source[i-1] == target[j-1] {
				d[i][j] = d[i-1][j-1]
			} else {
				min := d[i-1][j]
				if d[i][j-1] < min {
					min = d[i][j-1]
				}
				if d[i-1][j-1] < min {
					min = d[i-1][j-1]
				}
				d[i][j] = min + 1
			}
		}

	}
	return d[len(source)][len(target)]
}

func EnvString(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}
