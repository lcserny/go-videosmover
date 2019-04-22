package delete

import (
	"encoding/json"
	"fmt"
	"github.com/Bios-Marcel/wastebasket"
	"github.com/lcserny/go-videosmover/pkg/action"
	"github.com/lcserny/go-videosmover/pkg/convert"
	utils "github.com/lcserny/goutils"
)

const COULDNT_MOVE_TO_TRASH = "Couldn't move folder '%s' to trash"

func init() {
	action.Register("delete", &deleteAction{})
}

type deleteAction struct {
}

func (da *deleteAction) Execute(jsonPayload []byte, config *action.Config) (string, error) {
	var requests []RequestData
	err := json.Unmarshal(jsonPayload, &requests)
	utils.LogError(err)
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
			utils.LogError(err)
			resultList = append(resultList, ResponseData{
				req.Folder,
				[]string{fmt.Sprintf(COULDNT_MOVE_TO_TRASH, req.Folder)},
			})
		}
	}

	return convert.GetJSONEncodedString(resultList), nil
}
