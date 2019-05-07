package delete

import (
	"fmt"
	"github.com/Bios-Marcel/wastebasket"
	"github.com/lcserny/goutils"
	"videosmover/pkg/action"
	"videosmover/pkg/json"
)

const COULDNT_MOVE_TO_TRASH = "Couldn't move folder '%s' to trash"

func NewAction() action.Action {
	return &deleteAction{}
}

type deleteAction struct {
}

func (da *deleteAction) Execute(jsonPayload []byte, config *action.Config) (string, error) {
	var requests []RequestData
	if err := json.Decode(jsonPayload, &requests); err != nil {
		goutils.LogError(err)
		return "", err
	}

	resultList := make([]ResponseData, 0)
	for _, req := range requests {
		if restricted := action.PathRemovalIsRestricted(req.Folder, config.RestrictedRemovePaths); restricted {
			resultList = append(resultList, ResponseData{
				req.Folder,
				[]string{fmt.Sprintf(action.RESTRICTED_PATH_REASON, req.Folder)},
			})
			continue
		}

		err := wastebasket.Trash(req.Folder)
		if err != nil {
			goutils.LogError(err)
			resultList = append(resultList, ResponseData{
				req.Folder,
				[]string{fmt.Sprintf(COULDNT_MOVE_TO_TRASH, req.Folder)},
			})
		}
	}

	return json.EncodeString(resultList)
}
