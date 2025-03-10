package main

import (
	ebustl "ebustl-to-ttml/internal/ebustl"
	ttmlgenerate "ebustl-to-ttml/internal/ttmlgenerate"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

const STL_EXT = "stl"

func ProcessAfile(sourcefilepath string, outputfilepath string, offsetSeconds int) error {
	fmt.Println("=============================================================")
	fmt.Println("=== Processing " + sourcefilepath + " ===")
	stl, err := ebustl.ReadStlFile(sourcefilepath)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		return err
	}

	// offset to zero based if requested
	if offsetSeconds != 0 {
		stl.OffsetCues(offsetSeconds * stl.Gsi.Fps())
	}

	fmt.Printf("number of ttis = %d\n", len(stl.Ttis))
	fmt.Printf("number of subtitles = %d\n", stl.Gsi.TotalNumberTtiBlocksInt)
	currentTime := time.Now()
	comment := "\nMatt's Golang app, \nsource file=" + sourcefilepath + "\nat " + currentTime.String() + "\nHostname: " + getHostname() + "\n"

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

func createFolder(folder string) error {
	if _, err := os.Stat(folder); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(folder, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	return nil
}

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

func fileexists(path string) bool {
	// https://stackoverflow.com/questions/10510691/how-to-check-whether-a-file-or-directory-exists
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if errors.Is(err, fs.ErrNotExist) {
		return false
	}
	return false
}

func folderScan(in_folder string, out_folder string, move_to string, offsetSeconds int, continueOnError bool) error {
	// run folder scan
	root := os.DirFS(in_folder)
	stlFiles, err := fs.Glob(root, "*."+STL_EXT)
	if err != nil {
		return err
	}

	for _, afile := range stlFiles {
		oldFileExtension := filepath.Ext(afile)
		sourcefilepath := path.Join(in_folder, afile)
		outputfilepath := path.Join(out_folder, strings.Replace(afile, oldFileExtension, ".ttml", 1))
		move_to_filepath := path.Join(move_to, afile)
		fmt.Println(outputfilepath)
		err := ProcessAfile(sourcefilepath, outputfilepath, offsetSeconds)
		if err != nil {
			if continueOnError {
				fmt.Printf("Encountered an error processing %s, error %s\n", sourcefilepath, err.Error())
				fmt.Println(err.Error())
			} else {
				return err
			}
		} else {
			// move file?
			if move_to != "" {
				err := os.Rename(sourcefilepath, move_to_filepath)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}
	return nil
}

func processFolder(in_folder string, out_folder string, move_to string, offsetSeconds int, continuousScanInterval int) error {
	// check source folder exists
	if !fileexists(in_folder) {
		return errors.New("source folder does not exist")
	}

	// check destination folder exists, create if not
	if !fileexists(out_folder) {
		err := createFolder(out_folder)
		if err != nil {
			return err
		}
	}

	// check move to exists
	if move_to != "" {
		if !fileexists(move_to) {
			err := createFolder(move_to)
			if err != nil {
				return err
			}
		}
	}

	if continuousScanInterval < 1 {
		return folderScan(in_folder, out_folder, "", offsetSeconds, false)
	}

	// continuous loop
	for {
		fmt.Printf("Scanning folder %s\n", in_folder)
		err := folderScan(in_folder, out_folder, move_to, offsetSeconds, true)
		if err != nil {
			fmt.Println("processFolder caught error " + err.Error())
		}
		fmt.Printf("Scanning completed, waiting %d seconds to scan again\n", continuousScanInterval)
		// sleep
		time.Sleep(time.Duration(continuousScanInterval) * time.Second)
	}
}

func main() {

	// extract the CLI arguments
	mode := ""
	offsetSeconds := 0
	continuousScanInterval := 0
	flag.StringVar(&mode, "mode", "single", "processing mode; single, folder, edit, split")
	flag.IntVar(&offsetSeconds, "offsetSeconds", 0, "If set, will offset all Cues by this number of seconds, e.g. -36000 would change a Cue at 10:00:00:00 to 00:00:00:00")
	flag.IntVar(&continuousScanInterval, "interval", 0, "If set for folder mode, folder poll with sleep this number of seconds between scans")
	flag.Parse()

	// switch based upon operation mode selected
	switch mode {
	case "single":
		if len(flag.Args()) != 2 {
			fmt.Println("usage: stl-to-ttml -mode=single <input file path> <output file path>")
			os.Exit(1)
		}

		input_file := flag.Args()[0]
		output_file := flag.Args()[1]

		err := ProcessAfile(input_file, output_file, offsetSeconds)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	case "folder":
		move_to := ""
		if continuousScanInterval == 0 {
			if len(flag.Args()) != 2 {
				fmt.Println("usage: stl-to-ttml -mode=folder <input folder path> <output folder path>")
				os.Exit(1)
			}
		} else {
			if len(flag.Args()) != 3 {
				fmt.Println("usage: stl-to-ttml -mode=folder -interval=30 <input folder path> <output folder path> <move to folder>")
				os.Exit(1)
			}
			move_to = flag.Args()[2]
		}

		input_folder := flag.Args()[0]
		output_folder := flag.Args()[1]

		err := processFolder(input_folder, output_folder, move_to, offsetSeconds, continuousScanInterval)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	case "edit":
		fmt.Println("Edit mode - TO DO!")
	case "split":
		fmt.Println("Split mode - TO DO!")
	default:
		fmt.Println("Unknown mode")
	}
}
