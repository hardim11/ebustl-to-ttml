package subtitleediting

import (
	"ebustl-to-ttml/internal/filehandler"
	"encoding/json"
	"errors"
)

type Part struct {
	TimecodeStart  string
	TimecodeEnd    string
	OutputFilePath string
}

type SplitJobRequest struct {
	InputFilePath string
	Parts         []Part
}

func SplitJobDeserialiseString(body string) (*SplitJobRequest, error) {
	res := SplitJobRequest{}
	err := json.Unmarshal([]byte(body), &res)
	if err != nil {
		return nil, err
	}

	// sanity check
	if res.InputFilePath == "" {
		return nil, errors.New("job InputFilePath not specified")
	}
	if len(res.Parts) < 1 {
		return nil, errors.New("job at least one part should be specified")
	}
	for _, y := range res.Parts {
		if (y.OutputFilePath == "") || (y.TimecodeEnd == "") || (y.TimecodeStart == "") {
			return nil, errors.New("part information is missing")
		}
	}
	return &res, nil
}

func SplitJobDeserialisefile(filepath string) (*SplitJobRequest, error) {
	byteValue, err := filehandler.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	return SplitJobDeserialiseString(string(*byteValue))
}
