package main

import (
	"path/filepath"
	"testing"
)

func TestSingleTtml(t *testing.T) {
	input_file := filepath.Join("..", "..", "testfiles", "0t0bmh6_ENG.stl")
	//output_file_ref := filepath.Join("..", "..", "testfiles", "0t0bmh6_ENG.imsc1.xml")
	output_file_tmp := filepath.Join("..", "..", "testfiles", "TEST_0t0bmh6_ENG.imsc1.xml")
	offsetSeconds := 0
	outputEbuTtD := false
	debug := false
	err := processAfile(input_file, output_file_tmp, offsetSeconds, outputEbuTtD, "test_run", debug)
	if err != nil {
		t.Fatalf("ERROR: Test1 err=%s", err.Error())
		return
	}
	// compare

}

func TestSingleEbuTtD(t *testing.T) {
	input_file := filepath.Join("..", "..", "testfiles", "0t0bmh6_ENG.stl")
	//output_file_ref := filepath.Join("..", "..", "testfiles", "0t0bmh6_ENG.imsc1.xml")
	output_file_tmp := filepath.Join("..", "..", "testfiles", "TEST_0t0bmh6_ENG.ebu.xml")
	offsetSeconds := 0
	outputEbuTtD := true
	debug := false
	err := processAfile(input_file, output_file_tmp, offsetSeconds, outputEbuTtD, "test_run", debug)
	if err != nil {
		t.Fatalf("ERROR: Test1 err=%s", err.Error())
		return
	}
	// compare

}
