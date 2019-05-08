package core

type VideoChecker interface {
	IsVideo(header []byte) bool
	IsSubtitle(path string) bool
}
