package main

import (
	ebustl "ebustl-to-ttml/internal/ebustl"
	ttmlgenerate "ebustl-to-ttml/internal/ttmlgenerate"
	"fmt"
	"os"
	"time"
)

func ProcessAfile(sourcefilepath string, outputfilepath string, offsetSeconds int) error {
	fmt.Println("=============================================================")
	fmt.Println("=== Processing " + sourcefilepath + " ===")
	stl, err := ebustl.ReadStlFile(sourcefilepath)
	if err != nil {
		fmt.Printf("ERROR: " + err.Error())
		return err
	}

	// offset to zero based if requested
	if offsetSeconds != 0 {
		stl.OffsetCues(offsetSeconds * stl.Gsi.Fps())
	}

	fmt.Printf("number of ttis = %d\n", len(stl.Ttis))
	fmt.Printf("number of subtitles = %d\n", stl.Gsi.TotalNumberTtiBlocksInt)
	currentTime := time.Now()
	comment := "\nMatt's Golang app, \nsource file=" + sourcefilepath + "\n at " + currentTime.String() + "\nHostname: " + getHostname() + "\n"

	config := ttmlgenerate.TtmlConvertConfigurationDefault()
	config.PreserveSpaces = true
	ttml, err := ttmlgenerate.CreateTtml(*stl, comment, &config)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	os.WriteFile(outputfilepath, []byte(ttml), 0644)
	fmt.Println("=== Outfile " + outputfilepath + " ===")
	fmt.Println("=============================================================")
	return nil
}

// func createFolder(folder string) error {
// 	if _, err := os.Stat(folder); errors.Is(err, os.ErrNotExist) {
// 		err := os.Mkdir(folder, os.ModePerm)
// 		if err != nil {
// 			fmt.Println(err)
// 			return err
// 		}
// 	}
// 	return nil
// }

// func fileNameWithoutExtension(fileName string) string {
// 	return strings.TrimSuffix(filepath.Base(fileName), filepath.Ext(fileName))
// }

func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "<unknown>"
	} else {
		return hostname
	}
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("usage: stl-to-ttml <input file path> <output file path>")
		os.Exit(1)
	}

	input_file := os.Args[1]
	output_file := os.Args[2]

	err := ProcessAfile(input_file, output_file, 0)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
