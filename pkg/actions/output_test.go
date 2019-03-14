package actions

import (
	. "github.com/lcserny/go-videosmover/pkg/models"
	"testing"
)

func TestOutputAction(t *testing.T) {
	testData := []testActionData{
		{
			request:  OutputRequestData{Name: "The Lord of the Rings: The Fellowship of <>the Ring (2001)", Type: "movie", SkipOnlineSearch: true},
			response: OutputResponseData{[]string{"The Lord Of The Rings The Fellowship Of The Ring (2001)"}, ORIGIN_NAME},
		},
	}

	runActionTest(t, testData, OutputAction)
}
