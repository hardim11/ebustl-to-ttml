package ebustl

import (
	"os"
	"testing"
)

func ReadFile(filepath string) (*[]byte, error) {
	v, err := os.ReadFile(filepath)
	return &v, err
}

func TestGsiBlock(t *testing.T) {

	test_file := `E:\Users\matth\Documents\ITV\fast-files\fast-s3\Subtitles\S3\0t0bmh6_ENG.stl`

	stl_payload, err := ReadFile(test_file)
	if err != nil {
		t.Fatalf("ERROR: " + err.Error())
	}

	stl, err := ReadStlPayload(*stl_payload)
	if err != nil {
		t.Fatalf("ERROR: " + err.Error())
	}

	codepagenumber := "850"
	//fmt.Printf("%+v\n", &(stl).Gsi.CodePageNumber)

	if (stl).Gsi.CodePageNumber != codepagenumber {
		t.Fatalf("ERROR: stl.Gsi.CodePageNumber: expected %+v, got %+v", codepagenumber, stl.Gsi.CodePageNumber)
	}

	DiskFormatCode := "STL25.01"
	if (stl).Gsi.DiskFormatCode != DiskFormatCode {
		t.Fatalf("ERROR: stl.Gsi.DiskFormatCode: expected %+v, got %+v", DiskFormatCode, stl.Gsi.DiskFormatCode)
	}

	DisplayStandardCode := byte(0x31)
	if (stl).Gsi.DisplayStandardCode != DisplayStandardCode {
		t.Fatalf("ERROR: stl.Gsi.DisplayStandardCode: expected %+v, got %+v", DisplayStandardCode, stl.Gsi.DisplayStandardCode)
	}

	CharacterCodeTable := "00"
	if (stl).Gsi.CharacterCodeTable != CharacterCodeTable {
		t.Fatalf("ERROR: stl.Gsi.CharacterCodeTable: expected %+v, got %+v", CharacterCodeTable, stl.Gsi.CharacterCodeTable)
	}

	LanguageCode := "09"
	if (stl).Gsi.LanguageCode != LanguageCode {
		t.Fatalf("ERROR: stl.Gsi.LanguageCode: expected %+v, got %+v", LanguageCode, stl.Gsi.LanguageCode)
	}

	OriginalProgrammeTitle := "Vera S12E2"
	if (stl).Gsi.OriginalProgrammeTitle != OriginalProgrammeTitle {
		t.Fatalf("ERROR: stl.Gsi.OriginalProgrammeTitle: expected %+v, got %+v", OriginalProgrammeTitle, stl.Gsi.OriginalProgrammeTitle)
	}

	OriginalEpisodeTitle := ""
	if (stl).Gsi.OriginalEpisodeTitle != OriginalEpisodeTitle {
		t.Fatalf("ERROR: stl.Gsi.OriginalEpisodeTitle: expected %+v, got %+v", OriginalEpisodeTitle, stl.Gsi.OriginalEpisodeTitle)
	}

	TranslatedProgrammeTitle := ""
	if (stl).Gsi.TranslatedProgrammeTitle != TranslatedProgrammeTitle {
		t.Fatalf("ERROR: stl.Gsi.TranslatedProgrammeTitle: expected %+v, got %+v", TranslatedProgrammeTitle, stl.Gsi.TranslatedProgrammeTitle)
	}

	TranslatorsName := ""
	if (stl).Gsi.TranslatorsName != TranslatorsName {
		t.Fatalf("ERROR: stl.Gsi.TranslatorsName: expected %+v, got %+v", TranslatorsName, stl.Gsi.TranslatorsName)
	}

	TranslatorsContactDetails := ""
	if (stl).Gsi.TranslatorsContactDetails != TranslatorsContactDetails {
		t.Fatalf("ERROR: stl.Gsi.TranslatorsContactDetails: expected %+v, got %+v", TranslatorsContactDetails, stl.Gsi.TranslatorsContactDetails)
	}

	SubtitleListReferenceCode := "1-7314-0052-001"
	if (stl).Gsi.SubtitleListReferenceCode != SubtitleListReferenceCode {
		t.Fatalf("ERROR: stl.Gsi.SubtitleListReferenceCode: expected %+v, got %+v", SubtitleListReferenceCode, stl.Gsi.SubtitleListReferenceCode)
	}

	CreationDate := "230127"
	if (stl).Gsi.CreationDate != CreationDate {
		t.Fatalf("ERROR: stl.Gsi.CreationDate: expected %+v, got %+v", CreationDate, stl.Gsi.CreationDate)
	}

	RevisionDate := "230201"
	if (stl).Gsi.RevisionDate != RevisionDate {
		t.Fatalf("ERROR: stl.Gsi.RevisionDate: expected %+v, got %+v", RevisionDate, stl.Gsi.RevisionDate)
	}

	RevisionNumber := "01"
	if (stl).Gsi.RevisionNumber != RevisionNumber {
		t.Fatalf("ERROR: stl.Gsi.RevisionNumber: expected %+v, got %+v", RevisionNumber, stl.Gsi.RevisionNumber)
	}

	TotalNumberTtiBlocks := "01522"
	if (stl).Gsi.TotalNumberTtiBlocks != TotalNumberTtiBlocks {
		t.Fatalf("ERROR: stl.Gsi.TotalNumberTtiBlocks: expected %+v, got %+v", TotalNumberTtiBlocks, stl.Gsi.TotalNumberTtiBlocks)
	}

	TotalNumberSubtitles := "01522"
	if (stl).Gsi.TotalNumberSubtitles != TotalNumberSubtitles {
		t.Fatalf("ERROR: stl.Gsi.TotalNumberSubtitles: expected %+v, got %+v", TotalNumberSubtitles, stl.Gsi.TotalNumberSubtitles)
	}

	TotalNumberSubtitleGroups := "001"
	if (stl).Gsi.TotalNumberSubtitleGroups != TotalNumberSubtitleGroups {
		t.Fatalf("ERROR: stl.Gsi.TotalNumberSubtitleGroups: expected %+v, got %+v", TotalNumberSubtitleGroups, stl.Gsi.TotalNumberSubtitleGroups)
	}

	MaximumNumberDisplayableCharacters := "40"
	if (stl).Gsi.MaximumNumberDisplayableCharacters != MaximumNumberDisplayableCharacters {
		t.Fatalf("ERROR: stl.Gsi.MaximumNumberDisplayableCharacters: expected %+v, got %+v", MaximumNumberDisplayableCharacters, stl.Gsi.MaximumNumberDisplayableCharacters)
	}

	MaximumNumberDisplayableRows := "23"
	if (stl).Gsi.MaximumNumberDisplayableRows != MaximumNumberDisplayableRows {
		t.Fatalf("ERROR: stl.Gsi.MaximumNumberDisplayableRows: expected %+v, got %+v", MaximumNumberDisplayableRows, stl.Gsi.MaximumNumberDisplayableRows)
	}

	TimeCodeStatus := byte(0x31)
	if (stl).Gsi.TimeCodeStatus != TimeCodeStatus {
		t.Fatalf("ERROR: stl.Gsi.TimeCodeStatus: expected %+v, got %+v", TimeCodeStatus, stl.Gsi.TimeCodeStatus)
	}

	TimeCodeStartProgramme := "10000000"
	if (stl).Gsi.TimeCodeStartProgramme != TimeCodeStartProgramme {
		t.Fatalf("ERROR: stl.Gsi.TimeCodeStartProgramme: expected %+v, got %+v", TimeCodeStartProgramme, stl.Gsi.TimeCodeStartProgramme)
	}

	TimeCodeFirstInCue := "10000318"
	if (stl).Gsi.TimeCodeFirstInCue != TimeCodeFirstInCue {
		t.Fatalf("ERROR: stl.Gsi.TimeCodeFirstInCue: expected %+v, got %+v", TimeCodeFirstInCue, stl.Gsi.TimeCodeFirstInCue)
	}

	TotalNumberDisks := byte('1')
	if (stl).Gsi.TotalNumberDisks != TotalNumberDisks {
		t.Fatalf("ERROR: stl.Gsi.TotalNumberDisks: expected %+v, got %+v", TotalNumberDisks, stl.Gsi.TotalNumberDisks)
	}

	DiskSequenceNumber := byte('1')
	if (stl).Gsi.DiskSequenceNumber != DiskSequenceNumber {
		t.Fatalf("ERROR: stl.Gsi.DiskSequenceNumber: expected %+v, got %+v", DiskSequenceNumber, stl.Gsi.DiskSequenceNumber)
	}

	CountryOrigin := "GBR"
	if (stl).Gsi.CountryOrigin != CountryOrigin {
		t.Fatalf("ERROR: stl.Gsi.CountryOrigin: expected %+v, got %+v", CountryOrigin, stl.Gsi.CountryOrigin)
	}

	Publisher := "accessibility@itv.com"
	if (stl).Gsi.Publisher != Publisher {
		t.Fatalf("ERROR: stl.Gsi.Publisher: expected %+v, got %+v", Publisher, stl.Gsi.Publisher)
	}
	EditorsName := "Jack, Kea, Charlotte W, Kate P,"
	if (stl).Gsi.EditorsName != EditorsName {
		t.Fatalf("ERROR: stl.Gsi.EditorsName: expected %+v, got %+v", EditorsName, stl.Gsi.EditorsName)
	}

	EditorsContactDetails := ""
	if (stl).Gsi.EditorsContactDetails != EditorsContactDetails {
		t.Fatalf("ERROR: stl.Gsi.EditorsContactDetails: expected %+v, got %+v", EditorsContactDetails, stl.Gsi.EditorsContactDetails)
	}

	// rendered values
	if stl.Gsi.TotalNumberTtiBlocksInt != 1522 {
		t.Fatalf("ERROR: stl.Gsi.EditorsContactDetails: expected %+v, got %+v", 1522, stl.Gsi.TotalNumberTtiBlocksInt)
	}

	// 10:00:03:18 == 900000 + 75 + 18 == 900093
	if stl.Ttis[0].TimeCodeInFrames != 900093 {
		t.Fatalf("ERROR: stl.Ttis[0].TimeCodeInFrame: expected %+v, got %+v", 900093, stl.Ttis[0].TimeCodeInFrames)
	}
	// 10:00:05:22 == 900000 + 125 + 22 == 900147
	if stl.Ttis[0].TimeCodeOutFrames != 900147 {
		t.Fatalf("ERROR: stl.Ttis[0].TimeCodeOutFrames: expected %+v, got %+v", 900147, stl.Ttis[0].TimeCodeOutFrames)
	}
	if stl.Ttis[0].SubtitleNumberRendered != 1 {
		t.Fatalf("ERROR: stl.Ttis[0].SubtitleNumberRendered: expected %+v, got %+v", 1, stl.Ttis[0].SubtitleNumberRendered)
	}
	if stl.Ttis[0].VerticalPosition != 22 {
		t.Fatalf("ERROR: stl.Ttis[0].VerticalPosition: expected %+v, got %+v", 22, stl.Ttis[0].VerticalPosition)
	}
	if stl.Ttis[0].JustificationCode != 1 {
		t.Fatalf("ERROR: stl.Ttis[0].JustificationCode: expected %+v, got %+v", 1, stl.Ttis[0].JustificationCode)
	}

	// fmt.Print(stl.Ttis[0].ToString())
	// fmt.Println()

	// fmt.Print(stl.Ttis[1000].ToString())
	// fmt.Println()

	// fmt.Print(stl.Ttis[1001].ToString())
	// fmt.Println()
}
