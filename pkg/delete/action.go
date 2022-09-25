package delete

import (
	"fmt"
	"github.com/lcserny/goutils"
	"videosmover/pkg"
	"videosmover/pkg/action"
)

const COULDNT_MOVE_TO_TRASH = "Couldn't move folder '%s' to trash"

type deleteAction struct {
	config     *core.ActionConfig
	codec      core.Codec
	trashMover core.TrashMover
}

func NewAction(cfg *core.ActionConfig, c core.Codec, tm core.TrashMover) core.Action {
	return &deleteAction{config: cfg, codec: c, trashMover: tm}
}

func (da deleteAction) Execute(jsonPayload []byte) (string, error) {
	var requests []RequestData
	if err := da.codec.Decode(jsonPayload, &requests); err != nil {
		goutils.LogError(err)
		return "", err
	}

	resultList := make([]ResponseData, 0)
	for _, req := range requests {
		if restricted := action.PathRemovalIsRestricted(req.Folder, da.config.RestrictedRemovePaths); restricted {
			resultList = append(resultList, ResponseData{
				req.Folder,
				[]string{fmt.Sprintf(action.RESTRICTED_PATH_REASON, req.Folder)},
			})
			continue
		}

		if err := da.trashMover.MoveToTrash(req.Folder); err != nil {
			goutils.LogError(err)
			resultList = append(resultList, ResponseData{
				req.Folder,
				[]string{fmt.Sprintf(COULDNT_MOVE_TO_TRASH, req.Folder)},
			})
		}
	}

	return da.codec.EncodeString(resultList)
}
