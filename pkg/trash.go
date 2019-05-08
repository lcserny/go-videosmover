package core

type TrashMover interface {
	MoveToTrash(path string) error
}
