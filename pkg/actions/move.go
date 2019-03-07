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
	actionConfig *ActionConfig
	resultList   *[]MoveResponseData
	request      *MoveRequestData
	folder, dest string
}

func newMoveExecutor(resultList *[]MoveResponseData, request *MoveRequestData, config *ActionConfig) *moveExecutor {
	videoDir := filepath.Dir(request.Video)
	destination := filepath.Join(request.DiskPath, request.OutName)
	return &moveExecutor{
		actionConfig: config,
		resultList:   resultList,
		request:      request,
		folder:       videoDir,
		dest:         destination,
	}
}

func (me *moveExecutor) appendToUnmovedReasons(reason ...string) {
	*me.resultList = append(*me.resultList, MoveResponseData{
		me.folder,
		reason,
	})
}

func (me *moveExecutor) canProceed(err error, reason string) bool {
	if err != nil {
		LogError(err)
		me.appendToUnmovedReasons(reason)
		return false
	}
	return true
}

func (me *moveExecutor) prepareMove() bool {
	if _, err := os.Stat(me.dest); os.IsNotExist(err) {
		err := os.MkdirAll(me.dest, os.ModePerm)
		if err != nil {
			me.appendToUnmovedReasons(fmt.Sprintf(COULDNT_CREATE_FOLDER_REASON, me.dest))
			return false
		}
	} else {
		if me.request.Type == MOVIE {
			me.appendToUnmovedReasons(fmt.Sprintf(MOVIE_EXISTS_REASON, me.request.OutName, me.request.DiskPath))
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
		me.appendToUnmovedReasons(unmovedSubsReasons...)
		return false
	}

	return true
}

func (me *moveExecutor) cleanIfPossible() {
	if restricted := pathRemovalIsRestricted(me.folder, me.actionConfig.RestrictedRemovePaths); restricted {
		me.appendToUnmovedReasons(fmt.Sprintf(RESTRICTED_PATH_REASON, me.folder))
		return
	}

	err := wastebasket.Trash(me.folder)
	if err != nil {
		LogError(err)
		me.appendToUnmovedReasons(fmt.Sprintf(COULDNT_REMOVE_FOLDER_REASON, me.folder))
	}
}

func MoveAction(jsonPayload []byte, config *ActionConfig) (string, error) {
	var requests []MoveRequestData
	err := json.Unmarshal(jsonPayload, &requests)
	LogError(err)
	if err != nil {
		return "", err
	}

	resultList := make([]MoveResponseData, 0)
	for _, req := range requests {
		moveExecutor := newMoveExecutor(&resultList, &req, config)
		if proceed := moveExecutor.prepareMove(); !proceed {
			continue
		}
		if proceed := moveExecutor.moveVideo(); !proceed {
			continue
		}
		if proceed := moveExecutor.moveSubs(); !proceed {
			continue
		}
		moveExecutor.cleanIfPossible()
	}

	return getJSONEncodedString(resultList), nil
}
