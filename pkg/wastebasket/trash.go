package wastebasket

import (
	"github.com/Bios-Marcel/wastebasket"
	"videosmover/pkg"
)

type trashMover struct {
}

func NewTrashMover() core.TrashMover {
	return &trashMover{}
}

func (tm trashMover) MoveToTrash(path string) error {
	return wastebasket.Trash(path)
}
