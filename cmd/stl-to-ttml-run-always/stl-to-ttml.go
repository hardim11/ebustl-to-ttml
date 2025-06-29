package main

import (
	ebustl "ebustl-to-ttml/internal/ebustl"
	filehandler "ebustl-to-ttml/internal/filehandler"
	ttmlgenerate "ebustl-to-ttml/internal/ttmlgenerate"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const STL_FILE_MATCH = "*.[sS][tT][lL]" // to ensure case sensitivity for stl
const TTML_EXT = ".ttml"
const CONFIG_FILE = "stl-to-ttml-run-always.json"
const VERSION_INFO = "1.1.0"

func processAfile(sourcefilepath string, outputfilepath string, debug bool) error {
	log.Println("=============================================================")
	log.Println("=== Processing " + sourcefilepath + " ===")
	stl, err := ebustl.ReadStlFile(sourcefilepath)
	if err != nil {
		log.Printf("ERROR: %s\n", err.Error())
		return err
	}

	log.Println("STL Read OK")
	log.Printf("number of ttis = %d\n", len(stl.Ttis))
	log.Printf("number of subtitles = %d\n", stl.Gsi.TotalNumberTtiBlocksInt)
	currentTime := time.Now()

	comment := fmt.Sprintf("Matt's Golang app version %s\nsource file=%s\nat %s\nHostname: %s\n",
		VERSION_INFO,
		sourcefilepath,
		currentTime.String(),
		getHostname(),
	)

	config := ttmlgenerate.TtmlConvertConfigurationDefault()
	config.PreserveSpaces = true
	config.Debug = debug
	ttml, err := ttmlgenerate.CreateTtml(*stl, comment, &config)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	err = filehandler.WriteFile(outputfilepath, []byte(ttml))
	if err != nil {
		log.Println("ERROR failed to write output file")
		return err
	}

	log.Println("=== Outfile " + outputfilepath + " ===")
	log.Println("=============================================================")
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

func folderScan(config serviceConfig) error {
	// run folder scan
	stlFiles, err := filehandler.FolderScan(config.SourceFolder, STL_FILE_MATCH)
	if err != nil {
		return err
	}

	for _, afile := range *stlFiles {
		oldFileExtension := filepath.Ext(afile)
		sourcefilepath := filepath.Join(config.SourceFolder, afile)
		outputfilepath := filepath.Join(config.TtmlOutputFolder, strings.Replace(afile, oldFileExtension, TTML_EXT, 1))
		move_to_filepath := filepath.Join(config.ProcessedFolder, afile)
		failed_filepath := filepath.Join(config.FailedFolder, afile)
		log.Println(outputfilepath)
		err := processAfile(sourcefilepath, outputfilepath, config.Debug)
		if err != nil {
			// oh dear things did not go well
			if config.StopOnError {
				// stop processing
				return err
			} else {
				// move to failed
				log.Printf("Encountered an error processing %s, error %s\n", sourcefilepath, err.Error())
				log.Println(err.Error())
				log.Printf("Moving source file %s to %s\n", sourcefilepath, failed_filepath)
				err := filehandler.MoveFile(sourcefilepath, failed_filepath)
				if err != nil {
					log.Println(err)
				}
			}
		} else {
			// all went well, so move file to the processed folder
			log.Printf("Moving source file %s to %s\n", sourcefilepath, move_to_filepath)
			err := filehandler.MoveFile(sourcefilepath, move_to_filepath)
			if err != nil {
				log.Println(err)
			}
		}
	}
	return nil
}

func checkFolders(config serviceConfig) error {
	// check for empty string
	if config.SourceFolder == "" {
		return errors.New("SourceFolder not set in configuration file")
	}
	if config.TtmlOutputFolder == "" {
		return errors.New("TtmlOutputFolder not set in configuration file")
	}
	if config.ProcessedFolder == "" {
		return errors.New("ProcessedFolder not set in configuration file")
	}
	if config.FailedFolder == "" {
		return errors.New("FailedFolder not set in configuration file")
	}

	// check source folder exists
	if !filehandler.DoesFileExist(config.SourceFolder) {
		return errors.New("SourceFolder \"" + config.SourceFolder + "\" does not exist")
	}

	// check destination folder exists, create if not
	if !filehandler.DoesFileExist(config.TtmlOutputFolder) {
		log.Println("TtmlOutputFolder \"" + config.TtmlOutputFolder + "\" folder not found, attempting to create...")
		err := filehandler.CreateFolder(config.TtmlOutputFolder)
		if err != nil {
			log.Println("Failed to create TtmlOutputFolder \"" + config.TtmlOutputFolder + "\" folder")
			return err
		}
	}

	// check move to exists
	if !filehandler.DoesFileExist(config.ProcessedFolder) {
		log.Println("ProcessedFolder \"" + config.ProcessedFolder + "\" folder not found, attempting to create...")
		err := filehandler.CreateFolder(config.ProcessedFolder)
		if err != nil {
			log.Println("Failed to create ProcessedFolder \"" + config.ProcessedFolder + "\" folder")
			return err
		}
	}

	// check failed exists
	if !filehandler.DoesFileExist(config.FailedFolder) {
		log.Println("FailedFolder \"" + config.FailedFolder + "\" folder not found, attempting to create...")
		err := filehandler.CreateFolder(config.FailedFolder)
		if err != nil {
			log.Println("Failed to create FailedFolder \"" + config.FailedFolder + "\" folder")
			return err
		}
	}

	log.Println("Pre Flight Folder check ok")
	return nil
}

func watchFolder(config serviceConfig) error {
	// check folders exist
	log.Println("EBU STL to TTML converter version " + VERSION_INFO)
	err := checkFolders(config)
	if err != nil {
		log.Println("Failed Folder preflight check, cannot start application")
		return err
	}

	// loop
	for {
		log.Printf("Scanning folder %s\n", config.SourceFolder)
		err := folderScan(config)
		if err != nil {
			log.Println("processFolder caught error " + err.Error())
			if config.StopOnError {
				return err
			}
		}
		log.Printf("Scanning completed, waiting %d seconds to scan again\n", config.ScanIntervalSeconds)
		// sleep
		time.Sleep(time.Duration(config.ScanIntervalSeconds) * time.Second)
	}
}

func main() {
	// read configuration file which should be in the CWD called "stl-to-ttml-run-always.json"
	config, err := read_config(CONFIG_FILE)
	if err != nil {
		// failed to read config
		log.Println(err)
		log.Println("ERROR: unable to start due to inability to read the configuration file \"" + CONFIG_FILE + "\"")
		os.Exit(1)
	}

	err = watchFolder(*config)
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}
