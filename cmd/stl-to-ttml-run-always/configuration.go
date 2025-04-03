package main

import (
	"encoding/json"
	"os"
)

type serviceConfig struct {
	SourceFolder        string
	TtmlOutputFolder    string
	ProcessedFolder     string
	FailedFolder        string
	ScanIntervalSeconds int
	StopOnError         bool
	Debug               bool
}

func read_config(file_path string) (*serviceConfig, error) {
	bytes, err := os.ReadFile(file_path)
	if err != nil {
		return nil, err
	}

	res := serviceConfig{}
	err = json.Unmarshal(bytes, &res)
	if err != nil {
		return nil, err
	}

	//TODO sanity check here

	//return
	return &res, nil
}
