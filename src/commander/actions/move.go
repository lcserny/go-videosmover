package actions

import (
	"encoding/json"
	"fmt"
	. "github.com/lcserny/go-videosmover/src/shared"
	. "github.com/lcserny/goutils"
	"os"
	"path/filepath"
	"strings"
)

const (
	MOVIE_EXISTS_REASON          = "Movie %s already exists in %s"
	COULDNT_CREATE_FOLDER_REASON = "Couldn't create folder %s"
	INPUT_VIDEO_PROBLEM_REASON   = "Video %s has problems"
	MOVING_PROBLEM_REASON        = "Problem occurred trying to move %s"
	RESTRICTED_PATH_REASON       = "Video dir %s is in restricted path"
	COULDNT_REMOVE_FOLDER_REASON = "Couldn't remove video dir %s"

	RESTRICTED_REMOVE_PATHS_FILE_KEY = "restricted_remove_paths"
)

var restrictedRemovePaths []string

func init() {
	restrictedRemovePathsContent, err := configFolder.FindString(RESTRICTED_REMOVE_PATHS_FILE_KEY)
	LogError(err)
	restrictedRemovePaths = GetLinesFromString(restrictedRemovePathsContent)
}

// TODO: test and refactor
func MoveAction(jsonPayload []byte) (string, error) {
	var request []MoveRequestData
	err := json.Unmarshal(jsonPayload, &request)
	LogError(err)
	if err != nil {
		return "", err
	}

	var resultList []MoveResponseData
	for _, req := range request {
		videoDir := filepath.Dir(req.Video)
		dest := filepath.Join(req.DiskPath, req.OutName)

		if _, err := os.Stat(dest); os.IsNotExist(err) {
			err := os.MkdirAll(dest, os.ModePerm)
			if err != nil {
				resultList = append(resultList, MoveResponseData{videoDir,
					[]string{fmt.Sprintf(COULDNT_CREATE_FOLDER_REASON, dest)}})
				continue
			}
		} else {
			if req.Type == MOVIE {
				resultList = append(resultList, MoveResponseData{videoDir,
					[]string{fmt.Sprintf(MOVIE_EXISTS_REASON, dest, req.DiskPath)}})
				continue
			}
		}

		info, err := os.Stat(req.Video)
		if err != nil {
			LogError(err)
			resultList = append(resultList, MoveResponseData{videoDir,
				[]string{fmt.Sprintf(INPUT_VIDEO_PROBLEM_REASON, req.Video)}})
			continue
		}

		// move video
		err = os.Rename(req.Video, filepath.Join(dest, info.Name()))
		if err != nil {
			LogError(err)
			resultList = append(resultList, MoveResponseData{videoDir,
				[]string{fmt.Sprintf(MOVING_PROBLEM_REASON, req.Video)}})
			continue
		}

		// move subs, forEach
		unmovedSubs := false
		var unmovedSubsReasons []string
		for _, sub := range req.Subs {
			trimmedSub := strings.Replace(sub, videoDir, "", 1)
			subDest := filepath.Join(dest, trimmedSub)
			err = os.Rename(sub, subDest)
			if err != nil {
				unmovedSubs = true
				LogError(err)
				unmovedSubsReasons = append(unmovedSubsReasons, fmt.Sprintf(MOVING_PROBLEM_REASON, sub))
				continue
			}
		}
		if unmovedSubs {
			resultList = append(resultList, MoveResponseData{videoDir, unmovedSubsReasons})
			continue
		}

		// clean folder if possible (not restricted path or such)
		restricted := false
		for _, restrictedFolder := range restrictedRemovePaths {
			if strings.Index(videoDir, restrictedFolder) != -1 {
				restricted = true
				resultList = append(resultList, MoveResponseData{videoDir,
					[]string{fmt.Sprintf(RESTRICTED_PATH_REASON, videoDir)}})
				break
			}
		}
		if !restricted {
			err := os.RemoveAll(videoDir)
			if err != nil {
				LogError(err)
				resultList = append(resultList, MoveResponseData{videoDir,
					[]string{fmt.Sprintf(COULDNT_REMOVE_FOLDER_REASON, videoDir)}})
			}
		}
	}

	return getJSONEncodedString(resultList), nil
}
