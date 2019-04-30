package delete

import (
	"encoding/json"
	"fmt"
	"github.com/Bios-Marcel/wastebasket"
	"github.com/lcserny/goutils"
	"go-videosmover/pkg/action"
	"go-videosmover/pkg/convert"
)

const COULDNT_MOVE_TO_TRASH = "Couldn't move folder '%s' to trash"

func NewAction() action.Action {
	return &deleteAction{}
}

type deleteAction struct {
}

func (da *deleteAction) Execute(jsonPayload []byte, config *action.Config) (string, error) {
	var requests []RequestData
	err := json.Unmarshal(jsonPayload, &requests)
	goutils.LogError(err)
	if err != nil {
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

	return convert.GetJSONEncodedString(resultList), nil
}
