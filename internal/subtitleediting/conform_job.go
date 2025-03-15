package subtitleediting

import (
	"ebustl-to-ttml/internal/filehandler"
	"encoding/json"
	"errors"
)

type Source struct {
	TimecodeStart string
	TimecodeEnd   string
	Padding       bool
}

type ConformJobRequest struct {
	OutputFilePath string
	InputFilePath  string
	Sources        []Source
}

func ConformJobRequestDeserialiseString(body string) (*ConformJobRequest, error) {

	res := ConformJobRequest{}
	err := json.Unmarshal([]byte(body), &res)
	if err != nil {
		return nil, err
	}

	// sanity check
	if res.OutputFilePath == "" {
		return nil, errors.New("job OutputFilePath not specified")
	}
	if len(res.Sources) < 1 {
		return nil, errors.New("job at least one source should be specified")
	}

	return &res, nil
}

func ConformJobRequestDeserialisefile(filepath string) (*ConformJobRequest, error) {
	byteValue, err := filehandler.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	return ConformJobRequestDeserialiseString(string(*byteValue))
}
