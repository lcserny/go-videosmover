package h2non

import (
	"github.com/h2non/filetype"
	"path/filepath"
	core "videosmover/pkg"
)

type videoChecker struct {
	config *core.ActionConfig
}

func NewVideoChecker(cfg *core.ActionConfig) core.VideoChecker {
	return &videoChecker{config: cfg}
}

func (mc videoChecker) IsVideo(header []byte) bool {
	for _, mType := range mc.config.AllowedMIMETypes {
		if filetype.IsMIME(header, mType) {
			return true
		}
	}
	return filetype.IsVideo(header)
}

func (mc videoChecker) IsSubtitle(path string) bool {
	ext := filepath.Ext(path)
	for _, allowedExt := range mc.config.AllowedSubtitleExtensions {
		if ext == allowedExt {
			return true
		}
	}
	return false
}
