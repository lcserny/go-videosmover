package actions

import (
	. "github.com/lcserny/go-videosmover/pkg/models"
	"testing"
)

func TestOutputAction(t *testing.T) {
	outputs := []testActionData{
		{
			request:  OutputRequestData{Name: "", Type: "", SkipCache: false, DiskPath: "", UseOnlineSearch: false},
			response: OutputResponseData{[]string{"", ""}, ORIGIN_TMDB},
		},
	}

	runActionTest(t, outputs, OutputAction)
}
