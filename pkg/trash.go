package core

import "os"

type TrashMover interface {
	MoveToTrash(path string) error
}

type trashMover struct {
}

func NewTrashMover() TrashMover {
	return &trashMover{}
}

func (tm trashMover) MoveToTrash(path string) error {
	return os.RemoveAll(path)
}
