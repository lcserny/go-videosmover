package actions

import (
	"encoding/json"
	"fmt"
	"github.com/Bios-Marcel/wastebasket"
	. "github.com/lcserny/go-videosmover/pkg/models"
	. "github.com/lcserny/goutils"
)

const COULDNT_MOVE_TO_TRASH = "Couldn't move folder '%s' to trash"

func DeleteAction(jsonPayload []byte) (string, error) {
	var requests []DeleteRequestData
	err := json.Unmarshal(jsonPayload, &requests)
	LogError(err)
	if err != nil {
		return "", err
	}

	resultList := make([]DeleteResponseData, 0)
	for _, req := range requests {
		if restricted := pathRemovalIsRestricted(req.Folder); restricted {
			resultList = append(resultList, DeleteResponseData{
				req.Folder,
				[]string{fmt.Sprintf(RESTRICTED_PATH_REASON, req.Folder)},
			})
			continue
		}

		err := wastebasket.Trash(req.Folder)
		if err != nil {
			LogError(err)
			resultList = append(resultList, DeleteResponseData{
				req.Folder,
				[]string{fmt.Sprintf(COULDNT_MOVE_TO_TRASH, req.Folder)},
			})
		}
	}

	return getJSONEncodedString(resultList), nil
}
