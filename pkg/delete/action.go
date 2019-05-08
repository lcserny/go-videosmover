package delete

import (
	"fmt"
	"github.com/Bios-Marcel/wastebasket"
	"github.com/lcserny/goutils"
	"videosmover/pkg"
	"videosmover/pkg/action"
)

const COULDNT_MOVE_TO_TRASH = "Couldn't move folder '%s' to trash"

type deleteAction struct {
	codec core.Codec
}

func NewAction(c core.Codec) action.Action {
	return &deleteAction{codec: c}
}

func (da *deleteAction) Execute(jsonPayload []byte, config *action.Config) (string, error) {
	var requests []RequestData
	if err := da.codec.Decode(jsonPayload, &requests); err != nil {
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

	return da.codec.EncodeString(resultList)
}
