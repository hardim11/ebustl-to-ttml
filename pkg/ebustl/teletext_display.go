package ebustl

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"

	log "github.com/sirupsen/logrus"
)

// going to try a new approach here
type TeletextPixel struct {
	Character        string
	BackgroundColour string
	ForegroundColour string
	IsPopulated      bool
	IsVisibleChar    bool
	IsNewBackground  bool // this is really only for the white box scenario where no start box is sent
}

func (tp *TeletextPixel) ToString() string {
	return fmt.Sprintf("Pixel: IsPopulated: %t\t IsVisibleChar: %t\t BackgroundColour: %s\t ForegroundColour: %s\t Character: %s",
		tp.IsPopulated,
		tp.IsVisibleChar,
		tp.ForegroundColour,
		tp.BackgroundColour,
		hex.EncodeToString([]byte(tp.Character)),
	)
}

type TeletextLine struct {
	Pixels        []TeletextPixel
	StartBoxIndex int
	EndBoxIndex   int
	ActiveLine    bool
}

type TeletextRaster struct {
	//Pixels [][]TeletextPixel // NOTE this is Y,X - row,column
	FirstActiveLine int
	ActiveLineCount int
	Lines           []TeletextLine
}

func (tr *TeletextRaster) Reset() {
	// NOTE this is Y,X - row,column
	foo := make([]TeletextLine, 24) // <<< FIX
	for row := range foo {
		foo[row].Pixels = make([]TeletextPixel, TELETEXT_WIDTH+5) // added a safety margin as seen some extend too wide
	}
	tr.Lines = foo
}

const COLOUR_BLACK = "#000000"
const COLOUR_RED = "#FF0000"
const COLOUR_GREEN = "#00FF00"
const COLOUR_YELLOW = "#FFFF00"
const COLOUR_BLUE = "#0000FF"
const COLOUR_MAGENTA = "#FF00FF"
const COLOUR_CYAN = "#00FFFF"
const COLOUR_WHITE = "#FFFFFF"
const TELETEXT_WIDTH = 40
const TELETEXT_HEIGHT = 23

// used for ANSI console output only
const AnsiReset = "\033[0m"
const AnsiRed = "\033[31m"
const AnsiGreen = "\033[32m"
const AnsiYellow = "\033[33m"
const AnsiBlue = "\033[34m"
const AnsiMagenta = "\033[35m"
const AnsiCyan = "\033[36m"
const AnsiGray = "\033[37m"
const AnsiWhite = "\033[97m"

func GetAccentedChar(diacritical byte, letter byte, debug bool) (string, error) {
	// this is going to be messy,
	// diacritical
	// https://en.wikipedia.org/wiki/EIA-608
	// using the EBU tech3360 recommendations
	// The accented letters in the Latin-based languages in Teletext are created according to the “floating accent" principle.
	// Column “C0" of the character code table 00 (Latin alphabet) in Annex B contains diacritical marks which are overlaid on
	// another character in the same presentation position. Each single accented character intended for presentation occupies
	// two bytes, and the diacritical mark is sent first (e.g. Ä = C8h 41h, ê = C3h 65h). This is opposite to the order used in
	// Unicode where combining character(s) follow the base character.
	// Ä = C8h 41h -> Ä = 0308h 41h
	// e.g. fmt.Println("A\u0308") > Ä

	if debug {
		log.Debug("GetAccentedChar I've been called, fingers crossed")
	}

	// convert the accent to unicode combining character
	unicode_diacritical := 0
	switch diacritical {
	case 0xc1:
		unicode_diacritical = 0x0300
	case 0xc2:
		unicode_diacritical = 0x0301
	case 0xc3:
		unicode_diacritical = 0x0302
	case 0xc4:
		unicode_diacritical = 0x0303
	case 0xc5:
		unicode_diacritical = 0x0304
	case 0xc6:
		unicode_diacritical = 0x0306
	case 0xc7:
		unicode_diacritical = 0x0307
	case 0xc8:
		unicode_diacritical = 0x0308
	case 0xca:
		unicode_diacritical = 0x030a
	case 0xcb:
		unicode_diacritical = 0x0327
	case 0xcc:
		unicode_diacritical = 0x0332
	case 0xcd:
		unicode_diacritical = 0x030B
	case 0xce:
		unicode_diacritical = 0x0328
	case 0xcf:
		unicode_diacritical = 0x030C
	default:
		if debug {
			log.Debug("GetAccentedChar oopsie")
		}
		return " ", errors.New("Failed to perform diacritical lookup for [" + fmt.Sprintf("%x", diacritical) + "][" + fmt.Sprintf("%x", letter) + "]")
	}

	// this seem wrong but maybe works?
	quote := fmt.Sprintf(`"%c\u%04x"`, letter, unicode_diacritical)
	res, err := strconv.Unquote(quote)
	if err != nil {
		return " ", err
	}

	if debug {
		log.Debug("GetAccentedChar I came up with '" + res + "'")
	}
	return res, nil
}

func isVivibleChar(input byte) bool {
	if input < 0x20 {
		return false
	}
	if input >= 0x7f && input < 0xa0 {
		return false
	}
	return true
}

func is_tti_empty(myTti Tti) bool {
	// return true if not actual chars found
	for _, chr := range myTti.ExtendedTextField {
		if isVivibleChar(chr) {
			return false
		}
	}
	return true
}

func CreateTeletextRasterFromTti(myTti Tti, codepage string, truncateOversizedLines bool, debug bool, ignoreInvalidDiacritical bool) (*TeletextRaster, error) {
	result := TeletextRaster{}
	result.Reset()

	if is_tti_empty(myTti) {
		log.Warn("This TTI has no visible characters")
		result.ActiveLineCount = -1
		result.FirstActiveLine = -1
		return &result, nil
	}

	var err error

	current_row := int(myTti.VerticalPosition)
	result.FirstActiveLine = current_row
	current_background_colour := COLOUR_BLACK
	current_foreground_colour := COLOUR_WHITE
	current_start_box_index := -1
	current_end_box_index := -1

	column := -1
	safety_we_saw_a_double_height := false
	last_char_was_cr := false
	accented := byte(0)
	for _, chr := range myTti.ExtendedTextField {

		// cope with double height CR's
		if last_char_was_cr && chr == 0x8a {
			last_char_was_cr = false
			continue
		}

		// deal with accented chars, do not increment position as next char is used for actual character
		if accented > 0 {
			if debug {
				log.Println("An accented Character is being processed")
			}

			// but only do for valid accented chars (seen some really flipping odd stuff in the sample files - corrupt I would say but...)

			// get the char using the char + accent
			newchar, err := GetAccentedChar(accented, chr, debug)
			// reset the accent
			accented = 0
			if err != nil && !ignoreInvalidDiacritical {
				return nil, err
			}

			if column >= TELETEXT_WIDTH {
				log.Warnf("the column %d is greater than the width %d, characters should be dropped ", column, TELETEXT_WIDTH)
				if column > (TELETEXT_WIDTH + 5) {
					return nil, errors.New("a line is wider than the teletext width + safety margin, rejecting")
				}
			}

			result.Lines[current_row].ActiveLine = true
			result.Lines[current_row].Pixels[column] = TeletextPixel{
				BackgroundColour: current_background_colour,
				ForegroundColour: current_foreground_colour,
				Character:        string(newchar),
				IsVisibleChar:    isVivibleChar(chr),
				IsPopulated:      true,
			}
			continue
		}

		// move along
		column = column + 1

		// detect accented sequence here
		// if 0xC? then is dual byte accente 0xC0 == 192, 0xCf == 207
		if (chr >= 0xc0) && (chr <= 0xcf) {
			// set accented char
			accented = chr
			// get the next char to get actual char
			// skip forward
			continue
		}

		// if end of stuff, ignore
		if chr == 0x8f {
			continue
		}

		/////////////////////////////
		// deal with control chars //
		/////////////////////////////
		if chr < 32 {
			// confirm double height seen
			if chr == 0x0D {
				safety_we_saw_a_double_height = true
				continue
			}
			// colours
			if chr == 0 {
				current_foreground_colour = COLOUR_BLACK
				continue
			} else if chr == 1 {
				current_foreground_colour = COLOUR_RED
				continue
			} else if chr == 2 {
				current_foreground_colour = COLOUR_GREEN
				continue
			} else if chr == 3 {
				current_foreground_colour = COLOUR_YELLOW
				continue
			} else if chr == 4 {
				current_foreground_colour = COLOUR_BLUE
				continue
			} else if chr == 5 {
				current_foreground_colour = COLOUR_MAGENTA
				continue
			} else if chr == 6 {
				current_foreground_colour = COLOUR_CYAN
				continue
			} else if chr == 7 {
				current_foreground_colour = COLOUR_WHITE
				continue
			}
			// is new background
			if chr == 0x1d {
				// copy foreground to background
				// 	// this is new background,
				// 	// docs say:In practice this control code will typically appear in a sequence of 3 control codes, the first control code (00 – 0F) will
				// 	// set a colour, the second control code (1D) will switch this colour to the background, and the final control code will set a

				current_background_colour = current_foreground_colour
				result.Lines[current_row].Pixels[column] = TeletextPixel{
					BackgroundColour: current_background_colour,
					ForegroundColour: current_foreground_colour,
					Character:        "",
					IsPopulated:      true,
					IsVisibleChar:    false,
					IsNewBackground:  true,
				}
				continue
			}

			// //TODO new page
			// if chr == 0x1d {
			// 	// this is new background,
			// 	// docs say:In practice this control code will typically appear in a sequence of 3 control codes, the first control code (00 – 0F) will
			// 	// set a colour, the second control code (1D) will switch this colour to the background, and the final control code will set a
			// 	// new foreground colour
			// 	if debug {
			// 		log.Debug("DEBUG: <!-- surpessing 0x1d - New background  -->")
			// 	}
			// } else
			if chr == 0x1c {
				// this is Black background (1,2), which I don't believe is relevant here but maybe wrong
				if debug {
					log.Debug("DEBUG: <!-- surpessing 0x1c - Black background (1,2)  -->")
				}
			} else if (chr == 0x09) || (chr == 0x08) {
				log.Warn("WARN: Flash is not supported so ignoring control code [" + hex.EncodeToString([]byte{byte(chr)}) + "] ")
				if debug {
					log.Debug("DEBUG: cue text: " + string(chr))
				}
			} else if chr == 0x0a {
				// update the end box for the line
				if current_end_box_index < 0 {
					current_end_box_index = column
				}
				continue
			} else if chr == 0x0b {
				// update the start box for the line
				current_start_box_index = column
				continue
			} else {
				log.Error("ERROR: unhandled control character  [" + hex.EncodeToString([]byte{byte(chr)}) + "] " + string(chr))
				//return "", 0, 0, 0, errors.New("ttmlgenerate.reformatLine raised unhandled error - unhandled control character  [" + hex.EncodeToString([]byte{byte(chr)}) + "] " + string(chr))
			}
			continue
		}

		// if the character is CR (0x8a) then reset pixel location, reset colours
		// also set endbox if not set already
		if chr == 0x8a {
			//fmt.Printf("CR >>> column = %d, current_end_box_index = %d\n", column, current_end_box_index)
			if current_end_box_index < 0 && column >= TELETEXT_WIDTH {
				current_end_box_index = TELETEXT_WIDTH
			}
			column = -1
			result.Lines[current_row].StartBoxIndex = current_start_box_index
			result.Lines[current_row].EndBoxIndex = current_end_box_index

			current_start_box_index = -1
			current_end_box_index = -1
			current_row = current_row + 1
			current_background_colour = COLOUR_BLACK
			current_foreground_colour = COLOUR_WHITE

			if debug {
				log.Debug("CR")
			}
			last_char_was_cr = true
			continue
		}

		if debug {
			log.Debugf("Setting pixel current_row=%d, column=%d, value=% x\n", current_row, column, chr)
		}

		char := string(chr)

		// finally check for any extended characters here
		if chr > 0x7E {
			if debug {
				log.Println("An extended Character is being processed")
			}
			// get extended characters
			char, err = GetCodePageCharacter(codepage, chr)
			if err != nil {
				return nil, errors.New("unexpected character -  [" + hex.EncodeToString([]byte{chr}) + "] ")
			}
		}

		// set the contents
		if column >= TELETEXT_WIDTH {
			log.Warnf("the column %d is greater than the width %d, characters should be dropped ", column, TELETEXT_WIDTH)
			if column > TELETEXT_WIDTH+4 && !truncateOversizedLines {
				return nil, errors.New("a line is wider than the teletext width + safety margin, rejecting")
			}
		}

		// deal with truncation
		if column < TELETEXT_WIDTH+4 {
			result.Lines[current_row].ActiveLine = true
			result.Lines[current_row].Pixels[column] = TeletextPixel{
				BackgroundColour: current_background_colour,
				ForegroundColour: current_foreground_colour,
				Character:        char,
				IsPopulated:      true,
				IsVisibleChar:    isVivibleChar(chr),
				// was previous char a new
				IsNewBackground: false,
			}
		}
	} // character loop

	// if column == -1 {
	// 	// we have not done anything
	// }

	// cope with missing endbox if too long to add one
	if current_end_box_index < 0 && column >= TELETEXT_WIDTH {
		if debug {
			log.Info("Setting current_end_box_index to TELETEXT_WIDTH")
		}
		current_end_box_index = TELETEXT_WIDTH
	}
	// cope with completely missing endbox because?! seems to be when TTI is full
	if current_end_box_index < 0 && current_start_box_index > 0 {
		log.Warn("CreateTeletextRasterFromTti end box not set but start box is")
		current_end_box_index = column + 1
	}

	result.Lines[current_row].StartBoxIndex = current_start_box_index
	result.Lines[current_row].EndBoxIndex = current_end_box_index
	result.ActiveLineCount = current_row - result.FirstActiveLine + 1
	if !safety_we_saw_a_double_height {
		log.Error("ERROR did not see a double height char")
		return nil, errors.New("did not see a double height char")
	}

	// messy but remove any
	result.ConvertAllNonVisibleToSpaces()

	// cope with no start box
	result.FixupMissingStartBox()

	return &result, nil
}

func getAnsiColorCode(colour string) string {
	switch colour {
	case COLOUR_BLACK:
		return AnsiGray
	case COLOUR_RED:
		return AnsiRed
	case COLOUR_GREEN:
		return AnsiGreen
	case COLOUR_YELLOW:
		return AnsiYellow
	case COLOUR_BLUE:
		return AnsiBlue
	case COLOUR_MAGENTA:
		return AnsiMagenta
	case COLOUR_CYAN:
		return AnsiCyan
	case COLOUR_WHITE:
		return AnsiWhite
	}
	return AnsiReset

}

func (tr *TeletextRaster) PrintToConsole() string {
	res := ""
	for row := range tr.Lines {
		line := fmt.Sprintf("%02d", row)
		if tr.Lines[row].ActiveLine {
			line = line + " > "
		} else {
			line = line + " X "
		}
		for column := range tr.Lines[row].Pixels {
			if tr.Lines[row].Pixels[column].IsPopulated {
				line = line + getAnsiColorCode(tr.Lines[row].Pixels[column].ForegroundColour) + tr.Lines[row].Pixels[column].Character + getAnsiColorCode("")
			} else {
				line = line + "_"
			}

		}
		line = line + fmt.Sprintf("\t%d\t%d", tr.Lines[row].StartBoxIndex, tr.Lines[row].EndBoxIndex)
		res = res + line + "\n"

		// now do the nochar stuff
		line = "     "
		for column := range tr.Lines[row].Pixels {
			line = line + getAnsiColorCode(tr.Lines[row].Pixels[column].ForegroundColour) + tr.Lines[row].Pixels[column].Character + getAnsiColorCode("")
		}
		res = res + line + "\n"

		// now do the nochar stuff
		line = "     "
		for column := range tr.Lines[row].Pixels {
			if tr.Lines[row].Pixels[column].IsVisibleChar {
				line = line + "+"
			} else {
				line = line + "-"
			}
		}
		res = res + line + "\n"
	}
	res = res + fmt.Sprintf("First Active Line = %d\nActive Line Count = %d\n", tr.FirstActiveLine, tr.ActiveLineCount)
	return res
}

func (tr *TeletextRaster) ConvertAllNonVisibleToSpaces() {
	for line := range tr.Lines {
		for column := range tr.Lines[line].Pixels {
			if !tr.Lines[line].Pixels[column].IsVisibleChar {
				tr.Lines[line].Pixels[column].Character = " "
			}
		}
	}
}

func (tr *TeletextRaster) FixupMissingStartBox() {
	// trying to cope with the single white box at the end scenario
	for lineIdx, line := range tr.Lines {
		if line.ActiveLine {
			if line.StartBoxIndex < 0 {
				// try and find the start box
				for column, pixel := range line.Pixels {
					if pixel.IsNewBackground {
						tr.Lines[lineIdx].StartBoxIndex = column
					}
					// fmt.Printf("%d - Character %s\tBackgroundColour %s\tForegroundColour %s\tIsPopulated %t\tIsVisibleChar %t\tIsNewBackground %t\n",
					// 	column,
					// 	hex.EncodeToString([]byte(pixel.Character)),
					// 	pixel.BackgroundColour,
					// 	pixel.ForegroundColour,
					// 	pixel.IsPopulated,
					// 	pixel.IsVisibleChar,
					// 	pixel.IsNewBackground,
					// )
				}

			}
		}
	}
}
