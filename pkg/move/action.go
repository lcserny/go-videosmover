package move

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"videosmover/pkg"
	"videosmover/pkg/action"

	"github.com/lcserny/goutils"
)

const (
	MOVIE_EXISTS_REASON          = "Movie '%s' already exists in '%s'"
	COULDNT_CREATE_FOLDER_REASON = "Couldn't create folder '%s'"
	INPUT_VIDEO_PROBLEM_REASON   = "Video '%s' has problems"
	MOVING_PROBLEM_REASON        = "Problem occurred trying to move '%s'"
	MAX_PATH_LENGTH_EXCEED       = "Max path length exceeded for '%s'"
	COULDNT_REMOVE_FOLDER_REASON = "Couldn't remove video dir '%s'"
)

func NewAction(cfg *core.ActionConfig, c core.Codec, tm core.TrashMover) core.Action {
	return &moveAction{config: cfg, codec: c, trashMover: tm}
}

type moveAction struct {
	config     *core.ActionConfig
	codec      core.Codec
	trashMover core.TrashMover
}

type moveExecutor struct {
	resultList   *[]ResponseData
	request      *RequestData
	folder, dest string
}

func newMoveExecutor(resultList *[]ResponseData, request *RequestData) *moveExecutor {
	videoDir := filepath.Dir(request.Video)
	destination := filepath.Join(request.DiskPath, request.OutName)
	return &moveExecutor{
		resultList: resultList,
		request:    request,
		folder:     videoDir,
		dest:       destination,
	}
}

func appendToUnmovedReasons(resultList *[]ResponseData, folder string, reason ...string) {
	*resultList = append(*resultList, ResponseData{
		folder,
		reason,
	})
}

func (me *moveExecutor) canProceed(err error, reason string) bool {
	if err != nil {
		goutils.LogError(err)
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
		if me.request.Type == action.MOVIE {
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
		subFolder := filepath.Dir(subDest)

		// move sub
		err := os.MkdirAll(subFolder, os.ModePerm)
		if err != nil {
			unmovedSubs = true
			goutils.LogError(err)
			unmovedSubsReasons = append(unmovedSubsReasons, fmt.Sprintf(MAX_PATH_LENGTH_EXCEED, subFolder))
			continue
		}

		// fix for Windows path length
		if len(sub) >= 255 || len(subDest) >= 255 {
			unmovedSubs = true
			goutils.LogError(fmt.Errorf(fmt.Sprintf(MOVING_PROBLEM_REASON, sub)))
			unmovedSubsReasons = append(unmovedSubsReasons, fmt.Sprintf(MOVING_PROBLEM_REASON, sub))
			continue
		}

		err = os.Rename(sub, subDest)
		if err != nil {
			unmovedSubs = true
			goutils.LogError(err)
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

func cleanFolders(cleaningSet []string, resultList *[]ResponseData, restrictedRemovePaths []string, tm core.TrashMover) {
	for _, folder := range cleaningSet {
		if restricted := action.PathRemovalIsRestricted(folder, restrictedRemovePaths); restricted {
			appendToUnmovedReasons(resultList, folder, fmt.Sprintf(action.RESTRICTED_PATH_REASON, folder))
			continue
		}

		err := tm.MoveToTrash(folder)
		if err != nil {
			goutils.LogError(err)
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

func (ma moveAction) Execute(jsonPayload []byte) (string, error) {
	var requests []RequestData
	if err := ma.codec.Decode(jsonPayload, &requests); err != nil {
		goutils.LogError(err)
		return "", err
	}

	resultList := make([]ResponseData, 0)
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
	cleanFolders(cleaningSet, &resultList, ma.config.RestrictedRemovePaths, ma.trashMover)

	return ma.codec.EncodeString(resultList)
}
