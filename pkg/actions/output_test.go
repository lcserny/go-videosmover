package actions

import (
	. "github.com/lcserny/go-videosmover/pkg/models"
	"testing"
)

func TestOutputAction(t *testing.T) {
	testData := []testActionData{
		{
			request:  OutputRequestData{Name: "", Type: "", SkipOnlineSearch: true},
			response: OutputResponseData{[]string{"", ""}, ORIGIN_TMDB},
		},
	}

	runActionTest(t, testData, OutputAction)
}
