package ebustl

import (
	"errors"
	"strconv"
	"strings"

	"github.com/bamiaux/iobit"
)

// GSI block
// 4.2. General Subtitle Information (GSI) block

type Gsi struct {
	CodePageNumber                        string //[3]byte  // CPN - 3 bytes
	DiskFormatCode                        string //[8]byte  // DFC - 8 bytes
	DisplayStandardCode                   byte   // DSC - 1 byte
	CharacterCodeTable                    string //[2]byte  // CCT - 2 bytes
	LanguageCode                          string //[2]byte  LC - 2 bytes
	OriginalProgrammeTitle                string //[32]byte // OPT - 32 bytes
	OriginalEpisodeTitle                  string //[32]byte // OET - 32 bytes
	TranslatedProgrammeTitle              string //[32]byte // TPT - 32 bytes
	TranslatedEpisodeTitle                string //[32]byte // TET - 32 bytes
	TranslatorsName                       string //[32]byte // TN - 32 bytes
	TranslatorsContactDetails             string //[32]byte // TCD - 32 bytes
	SubtitleListReferenceCode             string //[16]byte // SLR - 16 bytes
	CreationDate                          string //[6]byte // CD - 6 bytes
	RevisionDate                          string // [6]byte // RD - 6 bytes
	RevisionNumber                        string //[2]byte // RN - 2 bytes (0-99)
	TotalNumberTtiBlocks                  string //[5]byte // TNB - 5 bytes The range is 0-99 999 decimal.
	TotalNumberTtiBlocksInt               int    // the value changed to an int
	TotalNumberSubtitles                  string //[5]byte // TNS - 5 bytes
	TotalNumberSubtitleGroups             string //[3]byte // TNG
	MaximumNumberDisplayableCharacters    string //[2]byte // MNC
	MaximumNumberDisplayableCharactersInt int
	MaximumNumberDisplayableRows          string //[2]byte // MNR
	MaximumNumberDisplayableRowsInt       int
	TimeCodeStatus                        byte      // TCS
	TimeCodeStartProgramme                string    //[8]byte   // TCP
	TimeCodeFirstInCue                    string    //[8]byte   // TCF
	TotalNumberDisks                      byte      // TND
	DiskSequenceNumber                    byte      // DSN
	CountryOrigin                         string    //[3]byte   // CO
	Publisher                             string    //[32]byte  // PUB
	EditorsName                           string    //[32]byte  // EN
	EditorsContactDetails                 string    //[32]byte  // ECD
	SpareBytes                            [75]byte  //
	UserDefinedArea                       [576]byte // UDA

}

func (g *Gsi) Fps() int {
	switch g.DiskFormatCode {
	case "STL25.01":
		return 25
	case "STL30.01":
		return 30
	}
	return 100
}

// func convertToNumber(b []byte) int {
// 	// convert to string,
// 	value := string(b)

// 	// convert string to number
// 	i, _ := strconv.Atoi(value)

//		return i
//	}
func (g *Gsi) Read(r *iobit.Reader) error {
	// read the byte array into the GSI
	g.CodePageNumber = string(r.Bytes(3))
	g.DiskFormatCode = string(r.Bytes(8))
	g.DisplayStandardCode = r.Byte()
	g.CharacterCodeTable = string(r.Bytes(2))
	g.LanguageCode = string(r.Bytes(2))
	g.OriginalProgrammeTitle = strings.TrimSpace(string(r.Bytes(32)))
	g.OriginalEpisodeTitle = strings.TrimSpace(string(r.Bytes(32)))
	g.TranslatedProgrammeTitle = strings.TrimSpace(string(r.Bytes(32)))
	g.TranslatedEpisodeTitle = strings.TrimSpace(string(r.Bytes(32)))
	g.TranslatorsName = strings.TrimSpace(string(r.Bytes(32)))
	g.TranslatorsContactDetails = strings.TrimSpace(string(r.Bytes(32)))
	g.SubtitleListReferenceCode = strings.TrimSpace(string(r.Bytes(16)))
	g.CreationDate = string(r.Bytes(6))
	g.RevisionDate = string(r.Bytes(6))
	g.RevisionNumber = string(r.Bytes(2))
	g.TotalNumberTtiBlocks = string(r.Bytes(5))
	g.TotalNumberSubtitles = string(r.Bytes(5))
	g.TotalNumberSubtitleGroups = string(r.Bytes(3))
	g.MaximumNumberDisplayableCharacters = string(r.Bytes(2))
	g.MaximumNumberDisplayableRows = string(r.Bytes(2))
	g.TimeCodeStatus = r.Byte()
	g.TimeCodeStartProgramme = string(r.Bytes(8))
	g.TimeCodeFirstInCue = string(r.Bytes(8))
	g.TotalNumberDisks = r.Byte()
	g.DiskSequenceNumber = r.Byte()
	g.CountryOrigin = string(r.Bytes(3))
	g.Publisher = strings.TrimSpace(string(r.Bytes(32)))
	g.EditorsName = strings.TrimSpace(string(r.Bytes(32)))
	g.EditorsContactDetails = strings.TrimSpace(string(r.Bytes(32)))

	// render stuff
	g.TotalNumberTtiBlocksInt, _ = strconv.Atoi(g.TotalNumberTtiBlocks[:])

	g.MaximumNumberDisplayableCharactersInt, _ = strconv.Atoi(g.MaximumNumberDisplayableCharacters[:])
	g.MaximumNumberDisplayableRowsInt, _ = strconv.Atoi(g.MaximumNumberDisplayableRows[:])

	// validation?
	if g.CodePageNumber != "437" && g.CodePageNumber != "850" && g.CodePageNumber != "860" && g.CodePageNumber != "863" && g.CodePageNumber != "865" {
		return errors.New("code page number \"" + g.CodePageNumber + "\" is not a valid code page")
	}

	return nil
}
