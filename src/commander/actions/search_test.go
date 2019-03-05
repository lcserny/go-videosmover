package actions

import (
	"encoding/json"
	. "github.com/lcserny/go-videosmover/src/shared"
	"testing"
)

// TODO: inMemory FS?
func TestSearchAction(t *testing.T) {
	searches := []struct {
		request       SearchRequestData
		response      []SearchResponseData
		expectedError error
	}{
		{
			request: SearchRequestData{"somePath"},
			response: []SearchResponseData{
				{"somePath/video.mp4", make([]string, 0)},
			},
			expectedError: nil,
		},
	}

	for _, search := range searches {
		reqBytes, err := json.Marshal(search.request)
		if err != nil {
			t.Fatalf("Couldn't decode request: %+v", err)
		}
		resBytes, err := json.Marshal(search.response)
		if err != nil {
			t.Fatalf("Couldn't decode response: %+v", err)
		}

		result, err := SearchAction(reqBytes)
		if err != nil {
			t.Fatalf("Error occurred: %+v", err)
		}
		if result != string(resBytes) {
			t.Errorf("Result mismatch, wanted %s, got: %s", string(resBytes), result)
		}
	}
}
