package actions

import (
	"encoding/json"
	"fmt"
	"github.com/Bios-Marcel/wastebasket"
	. "github.com/lcserny/go-videosmover/pkg/models"
	. "github.com/lcserny/goutils"
	"os"
	"path/filepath"
	"strings"
)

const (
	MOVIE_EXISTS_REASON          = "Movie '%s' already exists in '%s'"
	COULDNT_CREATE_FOLDER_REASON = "Couldn't create folder '%s'"
	INPUT_VIDEO_PROBLEM_REASON   = "Video '%s' has problems"
	MOVING_PROBLEM_REASON        = "Problem occurred trying to move '%s'"
	COULDNT_REMOVE_FOLDER_REASON = "Couldn't remove video dir '%s'"
)

type moveExecutor struct {
	resultList   *[]MoveResponseData
	request      *MoveRequestData
	folder, dest string
}

func newMoveExecutor(resultList *[]MoveResponseData, request *MoveRequestData) *moveExecutor {
	videoDir := filepath.Dir(request.Video)
	destination := filepath.Join(request.DiskPath, request.OutName)
	return &moveExecutor{
		resultList: resultList,
		request:    request,
		folder:     videoDir,
		dest:       destination,
	}
}

func appendToUnmovedReasons(resultList *[]MoveResponseData, folder string, reason ...string) {
	*resultList = append(*resultList, MoveResponseData{
		folder,
		reason,
	})
}

func (me *moveExecutor) canProceed(err error, reason string) bool {
	if err != nil {
		LogError(err)
		appendToUnmovedReasons(me.resultList, me.folder, reason)
		return false
	}
	return true
}

func (me *moveExecutor) prepareMove() bool {
	if _, err := os.Stat(me.dest); os.IsNotExist(err) {
		err := os.MkdirAll(me.dest, os.ModePerm)
		if err != nil {
			appendToUnmovedReasons(me.resultList, me.folder, fmt.Sprintf(COULDNT_CREATE_FOLDER_REASON, me.dest))
			return false
		}
	} else {
		if me.request.Type == MOVIE {
			appendToUnmovedReasons(me.resultList, me.folder, fmt.Sprintf(MOVIE_EXISTS_REASON, me.request.OutName, me.request.DiskPath))
			return false
		}
	}
	return true
}

func (me *moveExecutor) moveVideo() bool {
	info, err := os.Stat(me.request.Video)
	if proceed := me.canProceed(err, fmt.Sprintf(INPUT_VIDEO_PROBLEM_REASON, me.request.Video)); !proceed {
		return false
	}

	err = os.Rename(me.request.Video, filepath.Join(me.dest, info.Name()))
	if proceed := me.canProceed(err, fmt.Sprintf(MOVING_PROBLEM_REASON, me.request.Video)); !proceed {
		return false
	}

	return true
}

func (me *moveExecutor) moveSubs() bool {
	unmovedSubs := false
	var unmovedSubsReasons []string

	for _, sub := range me.request.Subs {
		trimmedSub := strings.Replace(sub, me.folder, "", 1)
		subDest := filepath.Join(me.dest, trimmedSub)
		// move sub
		err := os.Rename(sub, subDest)
		if err != nil {
			unmovedSubs = true
			LogError(err)
			unmovedSubsReasons = append(unmovedSubsReasons, fmt.Sprintf(MOVING_PROBLEM_REASON, sub))
			continue
		}
	}

	if unmovedSubs {
		appendToUnmovedReasons(me.resultList, me.folder, unmovedSubsReasons...)
		return false
	}

	return true
}

func cleanFolders(cleaningSet []string, resultList *[]MoveResponseData, restrictedRemovePaths []string) {
	for _, folder := range cleaningSet {
		if restricted := pathRemovalIsRestricted(folder, restrictedRemovePaths); restricted {
			appendToUnmovedReasons(resultList, folder, fmt.Sprintf(RESTRICTED_PATH_REASON, folder))
			continue
		}

		err := wastebasket.Trash(folder)
		if err != nil {
			LogError(err)
			appendToUnmovedReasons(resultList, folder, fmt.Sprintf(COULDNT_REMOVE_FOLDER_REASON, folder))
		}
	}
}

func addToCleanSet(cleaningSet *[]string, folder string) {
	for _, existingFolder := range *cleaningSet {
		if existingFolder == folder {
			return
		}
	}
	*cleaningSet = append(*cleaningSet, folder)
}

func MoveAction(jsonPayload []byte, config *ActionConfig) (string, error) {
	var requests []MoveRequestData
	err := json.Unmarshal(jsonPayload, &requests)
	LogError(err)
	if err != nil {
		return "", err
	}

	resultList := make([]MoveResponseData, 0)
	cleaningSet := make([]string, 0)
	for _, req := range requests {
		moveExecutor := newMoveExecutor(&resultList, &req)
		if proceed := moveExecutor.prepareMove(); !proceed {
			continue
		}
		if proceed := moveExecutor.moveVideo(); !proceed {
			continue
		}
		if proceed := moveExecutor.moveSubs(); !proceed {
			continue
		}
		addToCleanSet(&cleaningSet, moveExecutor.folder)
	}
	cleanFolders(cleaningSet, &resultList, config.RestrictedRemovePaths)

	return getJSONEncodedString(resultList), nil
}
