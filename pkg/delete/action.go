package delete

import (
	"encoding/json"
	"fmt"
	"github.com/Bios-Marcel/wastebasket"
	"github.com/lcserny/go-videosmover/pkg/action"
	"github.com/lcserny/go-videosmover/pkg/convert"
	. "github.com/lcserny/goutils"
)

const COULDNT_MOVE_TO_TRASH = "Couldn't move folder '%s' to trash"

func Action(jsonPayload []byte, config *convert.ActionConfig) (string, error) {
	var requests []convert.DeleteRequestData
	err := json.Unmarshal(jsonPayload, &requests)
	LogError(err)
	if err != nil {
		return "", err
	}

	resultList := make([]convert.DeleteResponseData, 0)
	for _, req := range requests {
		if restricted := action.PathRemovalIsRestricted(req.Folder, config.RestrictedRemovePaths); restricted {
			resultList = append(resultList, convert.DeleteResponseData{
				req.Folder,
				[]string{fmt.Sprintf(action.RESTRICTED_PATH_REASON, req.Folder)},
			})
			continue
		}

		err := wastebasket.Trash(req.Folder)
		if err != nil {
			LogError(err)
			resultList = append(resultList, convert.DeleteResponseData{
				req.Folder,
				[]string{fmt.Sprintf(COULDNT_MOVE_TO_TRASH, req.Folder)},
			})
		}
	}

	return convert.GetJSONEncodedString(resultList), nil
}
