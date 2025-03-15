package subtitleediting

import (
	ebustl "ebustl-to-ttml/internal/ebustl"
	"ebustl-to-ttml/internal/filehandler"
	"ebustl-to-ttml/internal/ttmlgenerate"
	"fmt"
)

// need to support both conform and split operations
func ConformJob(job_file_path string) error {
	// decode job
	job, err := ConformJobRequestDeserialisefile(job_file_path)
	if err != nil {
		return err
	}
	fmt.Println(job)

	// load up subtitle
	stl, err := ebustl.ReadStlFile(job.InputFilePath)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		return err
	}
	// merge TTI's
	stlmerged, err := stl.MergeExtensionBlocksTtis()
	if err != nil {
		return err
	}

	// now perform the conform
	// it is expected that the original has been merged already but what if it hasn't?

	// copy the original
	conformed_stl := ebustl.EbuStl{}
	conformed_stl.Gsi = stlmerged.Gsi // copy the header, need to change the # subs etc later
	// reset cue list
	conformed_stl.Ttis = []ebustl.Tti{}

	// now do the parts
	current_timecode_frame := 0
	for _, aPart := range job.Sources {
		part_start_frames := TcToFrames(aPart.TimecodeStart, conformed_stl.Gsi.Fps())
		part_end_frames := TcToFrames(aPart.TimecodeEnd, conformed_stl.Gsi.Fps())
		duration_frames := part_end_frames - part_start_frames
		if !aPart.Padding {
			// if not padding then insert the cues in the timeframe
			part, err := stlmerged.GetBetweenTimecodes(aPart.TimecodeStart, aPart.TimecodeEnd, current_timecode_frame)
			if err != nil {
				return err
			}

			// copy the cues to the output stl
			conformed_stl.Ttis = append(conformed_stl.Ttis, part.Ttis...)
		}

		current_timecode_frame = current_timecode_frame + duration_frames
	}

	// convert to TTML
	comment := "Conform of " + stl.Gsi.OriginalProgrammeTitle
	config := ttmlgenerate.TtmlConvertConfigurationDefault()
	config.PreserveSpaces = true
	//config.Debug = true
	ttml, err := ttmlgenerate.CreateTtml(conformed_stl, comment, &config)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	filehandler.WriteFile(job.OutputFilePath, []byte(ttml))
	fmt.Println("=== Outfile " + job.OutputFilePath + " ===")
	fmt.Println("=============================================================")
	return nil
}

func createPartFile(stl ebustl.EbuStl, part Part) error {

	fmt.Println("=============================================================")
	fmt.Println("=== Processing incode=" + part.TimecodeStart + ", outcode=" + part.OutputFilePath + " ===")

	// get only the subtitles cues we care about
	partStl, err := stl.GetBetweenTimecodes(part.TimecodeStart, part.TimecodeEnd, 0)
	if err != nil {
		return err
	}

	// sanity check?
	// what if none? That could be valid I guess??

	// convert to TTML
	comment := "Part of " + stl.Gsi.OriginalProgrammeTitle + ", incode=" + part.TimecodeStart + ", outcode=" + part.OutputFilePath
	config := ttmlgenerate.TtmlConvertConfigurationDefault()
	config.PreserveSpaces = true
	//config.Debug = true
	ttml, err := ttmlgenerate.CreateTtml(*partStl, comment, &config)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	outputfilepath := part.OutputFilePath
	filehandler.WriteFile(outputfilepath, []byte(ttml))
	fmt.Println("=== Outfile " + outputfilepath + " ===")
	fmt.Println("=============================================================")
	return nil

}

// need to support both conform and split operations
func PartJob(job_file_path string) error {
	// decode job
	job, err := SplitJobDeserialisefile(job_file_path)
	if err != nil {
		return err
	}
	fmt.Println(job)

	// load up subtitle
	stl, err := ebustl.ReadStlFile(job.InputFilePath)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		return err
	}
	// merge TTI's
	stlmerged, err := stl.MergeExtensionBlocksTtis()
	if err != nil {
		return err
	}

	// for each part
	for _, aPart := range job.Parts {
		err := createPartFile(*stlmerged, aPart)
		if err != nil {
			return err
		}

	}

	return nil
}
