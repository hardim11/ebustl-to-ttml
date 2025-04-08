package main

import (
	ebustl "ebustl-to-ttml/internal/ebustl"
	filehandler "ebustl-to-ttml/internal/filehandler"
	subtitleediting "ebustl-to-ttml/internal/subtitleediting"
	ttmlgenerate "ebustl-to-ttml/internal/ttmlgenerate"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const STL_EXT = "stl"

func processAfile(sourcefilepath string, outputfilepath string, offsetSeconds int, debug bool) error {
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
	config.Debug = debug
	ttml, err := ttmlgenerate.CreateTtml(*stl, comment, &config)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	filehandler.WriteFile(outputfilepath, []byte(ttml))
	fmt.Println("=== Outfile " + outputfilepath + " ===")
	fmt.Println("=============================================================")
	return nil
}

func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "<unknown>"
	} else {
		return hostname
	}
}

func folderScan(in_folder string, out_folder string, move_to string, offsetSeconds int, continueOnError bool, debug bool) error {
	// run folder scan
	stlFiles, err := filehandler.FolderScan(in_folder, "*."+STL_EXT)
	if err != nil {
		return err
	}

	for _, afile := range *stlFiles {
		oldFileExtension := filepath.Ext(afile)
		sourcefilepath := filepath.Join(in_folder, afile)
		outputfilepath := filepath.Join(out_folder, strings.Replace(afile, oldFileExtension, ".ttml", 1))
		move_to_filepath := filepath.Join(move_to, afile)
		fmt.Println(outputfilepath)
		err := processAfile(sourcefilepath, outputfilepath, offsetSeconds, debug)
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
				err := filehandler.MoveFile(sourcefilepath, move_to_filepath)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}
	return nil
}

func processFolder(in_folder string, out_folder string, move_to string, offsetSeconds int, continuousScanInterval int, debug bool) error {
	// check source folder exists
	if !filehandler.DoesFileExist(in_folder) {
		return errors.New("source folder does not exist")
	}

	// check destination folder exists, create if not
	if !filehandler.DoesFileExist(out_folder) {
		err := filehandler.CreateFolder(out_folder)
		if err != nil {
			return err
		}
	}

	// check move to exists
	if move_to != "" {
		if !filehandler.DoesFileExist(move_to) {
			err := filehandler.CreateFolder(move_to)
			if err != nil {
				return err
			}
		}
	}

	if continuousScanInterval < 1 {
		return folderScan(in_folder, out_folder, "", offsetSeconds, false, debug)
	}

	// continuous loop
	for {
		fmt.Printf("Scanning folder %s\n", in_folder)
		err := folderScan(in_folder, out_folder, move_to, offsetSeconds, true, debug)
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
	debug := false
	flag.StringVar(&mode, "mode", "single", "processing mode; single, folder, edit, split")
	flag.IntVar(&offsetSeconds, "offsetSeconds", 0, "If set, will offset all Cues by this number of seconds, e.g. -36000 would change a Cue at 10:00:00:00 to 00:00:00:00")
	flag.IntVar(&continuousScanInterval, "interval", 0, "If set for folder mode, folder poll with sleep this number of seconds between scans")
	flag.BoolVar(&debug, "debug", false, "Enables extra debug messages")

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

		err := processAfile(input_file, output_file, offsetSeconds, debug)
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

		err := processFolder(input_folder, output_folder, move_to, offsetSeconds, continuousScanInterval, debug)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	case "conform":
		if len(flag.Args()) != 1 {
			fmt.Println("usage: stl-to-ttml -mode=conform <job file path>")
			os.Exit(1)
		}

		job_file_path := flag.Args()[0]
		err := subtitleediting.ConformJob(job_file_path)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	case "split":
		if len(flag.Args()) != 1 {
			fmt.Println("usage: stl-to-ttml -mode=conform <job file path>")
			os.Exit(1)
		}

		job_file_path := flag.Args()[0]
		err := subtitleediting.PartJob(job_file_path)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	default:
		fmt.Println("Unknown mode")
		os.Exit(1)
	}
	os.Exit(0)
}
