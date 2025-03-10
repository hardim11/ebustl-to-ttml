package ttmlgenerate

import (
	"bytes"
	ebustl "ebustl-to-ttml/internal/ebustl"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const FPS = 25

const COLOUR_BLACK = "#000000"
const COLOUR_RED = "#FF0000"
const COLOUR_GREEN = "#00FF00"
const COLOUR_YELLOW = "#FFFF00"
const COLOUR_BLUE = "#0000FF"
const COLOUR_MAGENTA = "#FF00FF"
const COLOUR_CYAN = "#00FFFF"
const COLOUR_WHITE = "#FFFFFF"
const TELETEXT_WIDTH = 40

var configuration TtmlConvertConfiguration

type RegionStruct struct {
	LineNumber int
	RowCount   int
	LeftPad    int
	Width      int
}

func DebugPrintByteArrray(byteArry []byte) string {
	return "<!-- " + hex.EncodeToString([]byte(byteArry)) + " -->"
}

func GetAccentedChar(diacritical byte, letter byte) (string, error) {
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

	if configuration.Debug {
		fmt.Println("GetAccentedChar I've been called, fingers crossed")
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
		if configuration.Debug {
			fmt.Println("GetAccentedChar oopsie")
		}
		return "", errors.New("Failed to perform diacritical lookup for [" + fmt.Sprintf("%x", diacritical) + "][" + fmt.Sprintf("%x", letter) + "]")
	}

	// this seem wrong but maybe works?
	quote := fmt.Sprintf(`"%c\u%04x"`, letter, unicode_diacritical)
	res, err := strconv.Unquote(quote)
	if err != nil {
		return "", err
	}

	if configuration.Debug {
		fmt.Println("GetAccentedChar I came up with '" + res + "'")
	}
	return res, nil
}

func (r *RegionStruct) DecodeRegionNameString(RegionId string) {
	// we are using the region name to communicate the position using "." to split parameters
	// this decodes that and updates the region struct values

	r.LineNumber = 22
	r.RowCount = 1
	r.LeftPad = 0
	r.Width = 40
	parts := strings.Split(RegionId, ".")
	if parts[0] != "region" {
		return
	}
	if len(parts) != 5 {
		return
	}
	var err error
	r.LineNumber, err = strconv.Atoi(parts[1])
	if err != nil {
		return
	}
	r.RowCount, err = strconv.Atoi(parts[2])
	if err != nil {
		return
	}
	r.LeftPad, err = strconv.Atoi(parts[3])
	if err != nil {
		return
	}
	r.Width, err = strconv.Atoi(parts[4])
	if err != nil {
		return
	}
}

func getTtmlDefaultStyle() []TTMLOutStyle {
	res := []TTMLOutStyle{}

	tTMLFontSize := "18px"
	tTMLFontFamily := "proportionalSansSerif"
	aBackgroundColor := "rgba(0,0,0,0)"
	displayAlign := "center"
	ttsextent := "100% 33%"
	ttsorigin := "0% 66%"
	aTextAlign := "center"

	// backgroundStyle
	aStyleAttributes := StyleAttributes{
		TTMLBackgroundColor: &aBackgroundColor,
		TTMLDisplayAlign:    &displayAlign,
		TTMLExtent:          &ttsextent,
		TTMLFontFamily:      &tTMLFontFamily,
		TTMLFontSize:        &tTMLFontSize,
		TTMLOrigin:          &ttsorigin,
		TTMLTextAlign:       &aTextAlign,
	}
	var ttmlBackgroundStyle = TTMLOutStyle{TTMLOutHeader: TTMLOutHeader{
		ID:                     "backgroundStyle",
		TTMLOutStyleAttributes: ttmlOutStyleAttributesFromStyleAttributes(&aStyleAttributes),
	}}
	res = append(res, ttmlBackgroundStyle)

	// speakerStyle
	bBackgroundColor := "transparent"
	bTtscolor := "white"
	bTtstextOutline := "black 1px"
	bStyleAttributes := StyleAttributes{
		TTMLBackgroundColor: &bBackgroundColor,
		TTMLColor:           &bTtscolor,
		TTMLTextOutline:     &bTtstextOutline,
	}
	var ttmlSpeakerStyle = TTMLOutStyle{TTMLOutHeader: TTMLOutHeader{
		ID:                     "speakerStyle",
		Style:                  "backgroundStyle",
		TTMLOutStyleAttributes: ttmlOutStyleAttributesFromStyleAttributes(&bStyleAttributes),
	}}
	res = append(res, ttmlSpeakerStyle)

	// ttmlStyle
	cTMLFontSize := "80%"
	cTMLFontFamily := "monospaceSansSerif"
	cStyleAttributes := StyleAttributes{
		TTMLFontFamily: &cTMLFontFamily,
		TTMLFontSize:   &cTMLFontSize,
	}
	var ttmlTtmlStyle = TTMLOutStyle{TTMLOutHeader: TTMLOutHeader{
		ID:                     "ttmlStyle",
		TTMLOutStyleAttributes: ttmlOutStyleAttributesFromStyleAttributes(&cStyleAttributes),
	}}
	res = append(res, ttmlTtmlStyle)

	// textStyle
	dBackgroundColor := "black"
	dTtscolor := "white"
	dTtstextOutline := "none"
	dStyleAttributes := StyleAttributes{
		TTMLBackgroundColor: &dBackgroundColor,
		TTMLColor:           &dTtscolor,
		TTMLTextOutline:     &dTtstextOutline,
	}
	var ttmlTextStyle = TTMLOutStyle{TTMLOutHeader: TTMLOutHeader{
		ID:                     "textStyle",
		Style:                  "speakerStyle",
		TTMLOutStyleAttributes: ttmlOutStyleAttributesFromStyleAttributes(&dStyleAttributes),
	}}
	res = append(res, ttmlTextStyle)

	return res
}

func getTtmlDefaultRegions(debug_region bool) []TTMLOutRegion {
	res := []TTMLOutRegion{}

	// <region style="backgroundStyle" xml:id="background"></region>
	region1 := TTMLOutRegion{
		TTMLOutHeader: TTMLOutHeader{
			ID:                     "background",
			Style:                  "backgroundStyle",
			TTMLOutStyleAttributes: TTMLOutStyleAttributes{},
		},
	}
	res = append(res, region1)

	// <region xml:id="full" tts:extent="100% 100%" tts:origin="0% 0%"></region>
	region2extent := "100% 100%"
	region2origin := "0% 0%"
	region2 := TTMLOutRegion{
		TTMLOutHeader: TTMLOutHeader{
			ID: "full",
			TTMLOutStyleAttributes: TTMLOutStyleAttributes{
				Extent: &region2extent,
				Origin: &region2origin,
			},
		},
	}
	res = append(res, region2)

	// <region style="speakerStyle" xml:id="speaker"></region>
	region3 := TTMLOutRegion{
		TTMLOutHeader: TTMLOutHeader{
			ID:                     "speaker",
			Style:                  "speakerStyle",
			TTMLOutStyleAttributes: TTMLOutStyleAttributes{},
		},
	}
	res = append(res, region3)

	if debug_region {
		regionDebugExtent := "80% 16.67%"
		regionDebugOrigin := "10% 79.17%"
		regiondebug := TTMLOutRegion{
			TTMLOutHeader: TTMLOutHeader{
				ID: "debug_region",
				TTMLOutStyleAttributes: TTMLOutStyleAttributes{
					Extent: &regionDebugExtent,
					Origin: &regionDebugOrigin,
				},
			},
			Comment: "this is only for debugging use",
		}
		res = append(res, regiondebug)
	}

	return res
}

func getTtmlRegions(region_id string) TTMLOutRegion {
	/*
		input is like
			region.1.2.3.33
			region.22.1.12.15
			region.22.1.7.26
			region.20.2.10.19
			region.20.2.5.29

		calc:
			Origin X	(padding + control chars) * 2.5
			Origin Y	(vertical line - 1) * 4.165
			Extent X	2.5 * number viewable chars
			Extent Y	lines (double height lines) * 8.33
	*/

	var regiondata RegionStruct
	regiondata.DecodeRegionNameString(region_id)

	origin_x := float64(regiondata.LeftPad) * 2.5
	origin_y := float64(regiondata.LineNumber-1) * 4.165
	extent_x := float64(regiondata.Width) * 2.5
	extent_y := float64(regiondata.RowCount) * 8.33

	origin := fmt.Sprintf("%0.2f%% %0.2f%%", origin_x, origin_y)
	extent := fmt.Sprintf("%0.2f%% %0.2f%%", extent_x, extent_y)

	res := TTMLOutRegion{
		TTMLOutHeader: TTMLOutHeader{
			ID: region_id,
			TTMLOutStyleAttributes: TTMLOutStyleAttributes{
				Origin: &origin,
				Extent: &extent,
			},
		},
	}
	return res
}

func addSpace() string {
	if configuration.PreserveSpaces {
		return " "
	}
	return ""
}

func addPreserveSpace() string {
	if configuration.PreserveSpaces {
		return " xml:space=\"preserve\""
	} else {
		return ""
	}
}

func doSpansXml(spanOn bool, foregroundColour string, backgroundColour string) string {
	// if span already on, we need to close them
	res := ""
	// if there's already a span, close it.
	if spanOn {
		res = res + " </span>"
	}

	res = res + "<span tts:backgroundColor=\"" + backgroundColour + "\" tts:color=\"" + foregroundColour + "\"" + addPreserveSpace() + ">" + addSpace()
	return res
}

func XmlEscapeText(input string) string {
	// // EscapeText writes to w the properly escaped XML equivalent
	// // of the plain text data s.
	// func EscapeText(w io.Writer, s []byte) error {
	// 	return escapeText(w, s, true)
	// }
	// TODO I don't like mixing serialisation and escaped XML strings.
	buf := new(bytes.Buffer)
	xml.EscapeText(buf, []byte(input))
	return buf.String()
}

func getLeadingSpaces(input string) (span string, remaining string) {
	// add in a span
	pos := 0
	for _, ch := range []byte(input) {
		if ch != 0x20 {
			break
		}
		pos = pos + 1
	}
	//create span
	span = "<span tts:backgroundColor=\"transparent\" tts:color=\"transparent\" xml:space=\"preserve\">" + input[:pos] + "</span>"
	//update string
	remaining = input[pos:]

	return span, remaining
}

func getNextChar(inputtrimmed []byte, pos int) byte {
	// creeps forward on the array but protects against running over the end I hope
	if pos >= (len(inputtrimmed) - 1) {
		return 0
	}
	return inputtrimmed[pos+1]
}

func reformatLine(input string, codepage string, fixedalignment bool) (string, error) {
	res := ""
	//res = res + "\n<!-- " + hex.EncodeToString([]byte(input)) + " -->\n"

	backgroundColour := COLOUR_BLACK

	//TODO clean this up to be more efficient, e.g. if first line has a colour, don't add the default etc
	inputtrimmed := strings.ReplaceAll(input, "\b0a", "")       // end box
	inputtrimmed = strings.ReplaceAll(inputtrimmed, "\b0b", "") // start box
	inputtrimmed = strings.ReplaceAll(inputtrimmed, "\b0b", "") // normal height
	inputtrimmed = strings.ReplaceAll(inputtrimmed, "\b0d", "") // double height
	// replace all cr's with <br />
	inputtrimmed = strings.ReplaceAll(inputtrimmed, "\b8a", "<br />") // CR/LF

	//res = res + "<!-- " + hex.EncodeToString([]byte(inputtrimmed)) + " -->"

	if len(inputtrimmed) < 1 {
		return "", nil
	}

	// deal with no justification and leading spaces here?
	// only if the alignment is unset AND there are leading spaces
	prefixspan := ""
	if fixedalignment && inputtrimmed[0] == 0x20 {
		//fmt.Println(hex.EncodeToString([]byte(inputtrimmed)) + inputtrimmed)

		// add in a span with transparent background
		prefixspan, inputtrimmed = getLeadingSpaces(inputtrimmed)
	}

	// loop by character and work out when to turn on and off spans
	spanOn := false
	accented := byte(0)
	for idx, bt := range []byte(inputtrimmed) {
		//bt := byte(ch)
		//res = res + "<!-- " + string(bt) + " > " + hex.EncodeToString([]byte{bt}) + "-->"

		// deal with accented chars
		if accented > 0 {
			// get the char using the char + accent
			newchar, err := GetAccentedChar(accented, bt)
			// reset the accent
			accented = 0
			if err != nil {
				return "", err
			}

			res = res + XmlEscapeText(string(newchar))
			continue
		}

		if bt < 32 {
			// if next character is 0x1d (New Background, then this is actually setting the background)
			//background colours
			if getNextChar([]byte(inputtrimmed), idx) == 0x1d {
				// is background
				if bt == 0 {
					backgroundColour = COLOUR_BLACK
					continue
				} else if bt == 1 {
					backgroundColour = COLOUR_RED
					continue
				} else if bt == 2 {
					backgroundColour = COLOUR_GREEN
					continue
				} else if bt == 3 {
					backgroundColour = COLOUR_YELLOW
					continue
				} else if bt == 4 {
					backgroundColour = COLOUR_BLUE
					continue
				} else if bt == 5 {
					backgroundColour = COLOUR_MAGENTA
					continue
				} else if bt == 6 {
					backgroundColour = COLOUR_CYAN
					continue
				} else if bt == 7 {
					backgroundColour = COLOUR_WHITE
					continue
				}
			}

			// set foreground colour
			if bt == 0 {
				res = res + doSpansXml(spanOn, COLOUR_BLACK, backgroundColour)
				spanOn = true
			} else if bt == 1 {
				res = res + doSpansXml(spanOn, COLOUR_RED, backgroundColour)
				spanOn = true
			} else if bt == 2 {
				res = res + doSpansXml(spanOn, COLOUR_GREEN, backgroundColour)
				spanOn = true
			} else if bt == 3 {
				res = res + doSpansXml(spanOn, COLOUR_YELLOW, backgroundColour)
				spanOn = true
			} else if bt == 4 {
				res = res + doSpansXml(spanOn, COLOUR_BLUE, backgroundColour)
				spanOn = true
			} else if bt == 5 {
				res = res + doSpansXml(spanOn, COLOUR_MAGENTA, backgroundColour)
				spanOn = true
			} else if bt == 6 {
				res = res + doSpansXml(spanOn, COLOUR_CYAN, backgroundColour)
				spanOn = true
			} else if bt == 7 {
				res = res + doSpansXml(spanOn, COLOUR_WHITE, backgroundColour)
				spanOn = true
			} else if bt == 0x1d {
				// this is new page,
				// docs say:In practice this control code will typically appear in a sequence of 3 control codes, the first control code (00 – 0F) will
				// set a colour, the second control code (1D) will switch this colour to the background, and the final control code will set a
				// new foreground colour
				if configuration.Debug {
					res = res + "<!-- surpessing 0x1d - New background  -->"
				}
			} else if bt == 0x1c {
				// this is Black background (1,2), which I don't believe is relevant here but maybe wrong
				if configuration.Debug {
					res = res + "<!-- surpessing 0x1c - Black background (1,2)  -->"
				}
			} else if (bt == 0x09) || (bt == 0x08) {
				fmt.Println("WARN: Flash is not supported so ignoring control code [" + hex.EncodeToString([]byte{bt}) + "] ")
			} else {
				fmt.Println("ERROR: unhandled control character  [" + hex.EncodeToString([]byte{bt}) + "] " + input)
				return "", errors.New("ttmlgenerate.reformatLine raised unhandled error - unhandled control character  [" + hex.EncodeToString([]byte{bt}) + "] " + input)
			}
		} else {
			if (idx == 0) && (!spanOn) {
				res = res + doSpansXml(spanOn, COLOUR_WHITE, COLOUR_BLACK)
				spanOn = true
			}
			if bt > 126 {
				// if 0xC? then is dual byte accente
				if (bt >= 0xc0) && (bt <= 0xcf) {
					// set accented char
					accented = bt
					// get the next char to get actual char
					// skip forward
					continue
				}

				// get extended characters
				char, err := ebustl.GetCodePageCharacter(codepage, bt)
				if err != nil {
					return "", errors.New("unexpected character -  [" + hex.EncodeToString([]byte{bt}) + "] " + input)
				}
				res = res + XmlEscapeText(char)
			} else {
				// encode the char
				res = res + XmlEscapeText(string(bt))
			}
		}
	}
	if spanOn {
		res = res + addSpace() + "</span> "
	}
	return prefixspan + res, nil
}

func getTextField(b []byte) string {
	// get up until the first 0x8f
	tmp := string(b[:])
	i := strings.Index(tmp, "\x8f")
	//TODO fix extenmsion block > 1
	// I think this is fixed now, can remove?
	if i < 0 {
		//fmt.Println(hex.EncodeToString(b))
		fmt.Println("WARN: no 0x8f found so assume end of line")
		return tmp[:]
	}
	return tmp[:i]
}

func isPrintableChar(achar rune) bool {
	// ref tech3264.pdf
	// section 5 character code tables
	if (achar >= 0x20) && (achar <= 0x7f) {
		// valid char
		return true
	}
	if (achar >= 0xa1) && (achar <= 0xff) {
		// valid char
		return true
	}
	return false
}

func getLineLengthAndPad(line string, chars_width *int, chars_pad *int) {
	// calculate the values for this line
	tmp_width := 0
	tmp_path := 40

	tmp_width = len(line)

	for inc, achar := range line {
		if isPrintableChar(achar) {
			tmp_path = inc
			break
		}
	}

	// only update chars_width if > existing value
	if tmp_width > *chars_width {
		*chars_width = tmp_width
	}

	if *chars_pad < tmp_path {
		*chars_pad = tmp_path
	}

}

func getRegionForJustification(justificationCode byte, chars_width int, chars_pad int) (int, int) {
	//calculate the actual pad and width considering the justification option

	if justificationCode == 1 {
		// left
		return chars_pad, chars_width
	} else if justificationCode == 2 {
		// center
		tmp := (TELETEXT_WIDTH - chars_width) / 2
		return tmp, chars_width
	} else if justificationCode == 3 {
		//right
		return TELETEXT_WIDTH - chars_width, chars_width
	}

	// assume unset == space padded
	return chars_pad, chars_width
}

func getSubtitlePara(tti ebustl.Tti, codepage string, fixedalignment bool) (string, int, int, int, error) {
	// returns
	// 		xml for the paragraph
	//		row count
	//		left padding
	//		char count (including code chars)
	//		error
	/*
		the CR/LF indicator, used to initiate the second and subsequent rows of the subtitle display, is conveyed by character code 8Ah;
		the Text Field of the last TTI block of a subtitle must always terminate with code 8Fh;
		unused space in the Text Field will be set to 8Fh.
	*/
	// get raw TextField until 0x8F
	textfield := getTextField(tti.ExtendedTextField[:])

	// deduplicate some codes if double height
	if strings.Contains(textfield, "\x0d") {
		//doubleHeight = true
		// change start box to 1
		textfield = strings.Replace(textfield, "\x0b\x0b", "\x0b", -1)
		// change end box to 1
		textfield = strings.Replace(textfield, "\x0a\x0a", "\x0a", -1)
		// change STL CR to 1
		textfield = strings.Replace(textfield, "\x8a\x8a", "\x8a", -1)
	}

	// we need to split into rows
	lines := strings.Split(textfield, "\x8a")

	// track total char with
	chars_width := 0 // this one we care about the largest
	chars_pad := 0
	res := ""
	for _, y := range lines {
		line := y

		// count the pad chars
		// count the length
		//res = res + "<!-- " + hex.EncodeToString([]byte(line)) + " -->"
		getLineLengthAndPad(line, &chars_width, &chars_pad)

		// ignore the start / end markers / double height
		line = strings.Replace(line, "\x0d", "", -1)
		line = strings.Replace(line, "\x0a", "", -1)
		line = strings.Replace(line, "\x0b", "", -1)
		//getLineLengthAndPad(line, &chars_pad, &chars_width)
		reformatted_line, err := reformatLine(line, codepage, fixedalignment)
		if err != nil {
			return "", 0, 0, 0, err
		}
		res = res + reformatted_line + "<br />"
	}

	// kludge, trim last <br />
	res = strings.TrimSuffix(res, "<br />")

	// calculate the actual pad and width considering the justification option
	chars_pad, chars_width = getRegionForJustification(tti.JustificationCode, chars_width, chars_pad)
	return res, len(lines), chars_pad, chars_width, nil
}

func divmod(numerator, denominator int64) (quotient, remainder int64) {
	// from https://stackoverflow.com/questions/43945675/division-with-returning-quotient-and-remainder
	quotient = numerator / denominator // integer division, decimals are truncated
	remainder = numerator % denominator
	return
}

func FramesToTcMs(frames int64) string {
	hr, remainder := divmod(frames, (60 * 60 * FPS))
	mn, remainder := divmod(remainder, (60 * FPS))
	sc, fr := divmod(remainder, FPS)

	ms := (fr * 1000) / FPS
	return fmt.Sprintf("%02d:%02d:%02d.%03d", hr, mn, sc, ms)
}

func getSubtitle(tti ebustl.Tti, codepage string, addId bool) (*TTMLOutSubtitle, string, error) {
	/*

			target

		      <p begin="00:34:41.920" end="00:34:44.120" region="region-181" tts:fontSize="200%" tts:lineHeight="120%" tts:textAlign="right">
		        <span tts:backgroundColor="#000000" tts:color="#FFFFFF">NIGEL:<br></br></span>
		        <span tts:backgroundColor="#000000" tts:color="#00FF00">&#39;Four? But I only have three
		          chairs.</span>
		      </p>

				Textfield
				0D 07 0B 0B 4E 49 47 45 4C 3A 0A 0A 8A 8A 0D 02
				0B 0B 27 46 6F 75 72 3F 20 42 75 74 20 49 20 6F
				6E 6C 79 20 68 61 76 65 20 74 68 72 65 65 20 63
				68 61 69 72 73 2E 8F 8F 8F 8F 8F 8F 8F 8F 8F 8F
				8F 8F 8F 8F 8F 8F 8F 8F 8F 8F 8F 8F 8F 8F 8F 8F
				8F 8F 8F 8F 8F 8F 8F 8F 8F 8F 8F 8F 8F 8F 8F 8F
				8F 8F 8F 8F 8F 8F 8F 8F 8F 8F 8F 8F 8F 8F 8F 8F

				0D (double height)
				07 (alpha white)  <<< starts the span?
				0B (start box)
				0B (start box)
				4E N
				49 i
				47 g
				45 e
				4C l
				3A :
				0A (end box)
				0A (end box)
				8A (cr)
				8A (cr)
				0D (double height)
				02 (alpha green)    <<< starts the span? end previous span
				0B (start box)
				0B (start box)
				27 '
				46 F
				6F o
				75 u
				72 r
				3F ?
				20 (space)
				42 B
				75 u
				74 t
				20 (space)
				49 I
				20 (space)
				6F o
				6E n
				6C l
				79 y
				20 (space)
				68 h
				61 a
				76 v
				65 e
				20 (space)
				74 t
				68 h
				72 r
				65 e
				65 e
				20 (space)
				63 c
				68 h
				61 a
				69 i
				72 r
				73 s
				2E .
				8F ( the Text Field of the last TTI block of a subtitle must always terminate with code 8Fh )
				8F 8F 8F 8F 8F 8F 8F 8F 8F ( unused space in the Text Field will be set to 8Fh. )

				black background assumed at start of each row

				we then use the justification info + character counts to determine the region for the sub
	*/

	// deal with no justification and leading spaces here?
	sub_para, rows, left_pad, char_count, err := getSubtitlePara(tti, codepage, (tti.JustificationCode == 0))
	if err != nil {
		return nil, "", err
	}
	region := "region." + strconv.Itoa(int(tti.VerticalPosition)) + "." + strconv.Itoa(rows) + "." + strconv.Itoa(left_pad) + "." + strconv.Itoa(char_count)

	ttsTextAlign := "left"
	if tti.JustificationCode == 1 {
		ttsTextAlign = "left"
	} else if tti.JustificationCode == 2 {
		ttsTextAlign = "center"
	} else if tti.JustificationCode == 3 {
		ttsTextAlign = "right"
	}

	ttsFontSize := "200%"
	ttsLineHeight := "120%"
	res := TTMLOutSubtitle{
		Region: region,
		Begin:  FramesToTcMs(int64(tti.TimeCodeInFrames)),
		End:    FramesToTcMs(int64(tti.TimeCodeOutFrames)),
		Text:   sub_para,
		TTMLOutStyleAttributes: TTMLOutStyleAttributes{
			FontSize:   &ttsFontSize,
			LineHeight: &ttsLineHeight,
			TextAlign:  &ttsTextAlign,
		},
	}
	if addId {
		res.ID = strconv.Itoa(tti.SubtitleNumberRendered)
	}

	return &res, region, nil
}

func secondsToTc(seconds int64) string {
	hr, remainder := divmod(seconds, (60 * 60))
	mn, sc := divmod(remainder, (60))
	return fmt.Sprintf("%02d:%02d:%02d", hr, mn, sc)
}

func (t *TTMLOut) debugShuffleTimeCodes() {
	shuffle_spacing := 2
	seconds := shuffle_spacing
	for idx := range t.Body.Div.Subtitles {
		t.Body.Div.Subtitles[idx].Begin = secondsToTc(int64(seconds)) + ".000"
		t.Body.Div.Subtitles[idx].End = secondsToTc(int64(seconds+(shuffle_spacing-1))) + ".500"
		seconds = seconds + shuffle_spacing
	}
}

func CreateTtml(stl ebustl.EbuStl, comment string, config *TtmlConvertConfiguration) (string, error) {

	configuration = *config

	// merge TTI's
	stlmerged, err := stl.MergeExtensionBlocksTtis()
	if err != nil {
		return "", err
	}

	// only allow the stuff we support
	if (stlmerged.Gsi.CharacterCodeTable != "00") && (stlmerged.Gsi.CharacterCodeTable != "01") {
		return "", errors.New("only Latin code table (00), Latin/Cyrillic (01) is supported, code table ID " + stlmerged.Gsi.CharacterCodeTable + " specified")
	}
	if (stlmerged.Gsi.DisplayStandardCode != '1') && (stlmerged.Gsi.DisplayStandardCode != '2') {
		return "", errors.New("only Level 1 and Level 2 Teletext supported, " + string(stlmerged.Gsi.DisplayStandardCode) + " specified")
	}
	if stlmerged.Gsi.DiskFormatCode != "STL25.01" {
		return "", errors.New("only 25fps supported")
	}

	res := TTMLOut{}
	if comment != "" {
		res.Comment = comment
	}

	// root node
	res.XMLNamespaceTTM = "http://www.w3.org/ns/ttml#metadata"
	res.XMLNamespaceTTS = "http://www.w3.org/ns/ttml#styling"
	res.XMLNamespaceTTP = "http://www.w3.org/ns/ttml#parameter"
	res.XMLNamespaceSmpte = "http://www.smpte-ra.org/schemas/2052-1/2010/smpte-tt"
	res.XMLNamespaceIMSC = "http://www.w3.org/ns/ttml/profile/imsc1#styling"
	res.XMLNamespaceITTP = "http://www.w3.org/ns/ttml/profile/imsc1#parameter"
	//res.XMLNamespaceEbuTt = "urn:ebu:tt:style"

	res.Lang = getxmlLanguageCode(stl.Gsi.LanguageCode)
	res.CellRsolution = configuration.CellRsolution

	res.Head.Metadata = &TTMLOutMetadata{}
	res.Head.Metadata.Title = stlmerged.Gsi.OriginalProgrammeTitle
	res.Head.Metadata.Description = stlmerged.Gsi.SubtitleListReferenceCode

	fixedStyles := getTtmlDefaultStyle()
	res.Head.Styles = append(res.Head.Styles, fixedStyles[:]...)

	res.Head.Regions = []TTMLOutRegion{}
	res.Head.Regions = append(res.Head.Regions, getTtmlDefaultRegions(configuration.Debug)...)

	res.Body.Style = "ttmlStyle"

	// regions list
	regions := map[string]bool{}

	// body
	for _, aTti := range stlmerged.Ttis {
		if aTti.CommentFlag != 1 {
			// add a para
			para, region, err := getSubtitle(aTti, stlmerged.Gsi.CharacterCodeTable, configuration.AddId)
			if err != nil {
				return "", err
			}
			res.Body.Div.Subtitles = append(res.Body.Div.Subtitles, *para)

			// kludge way to have unique region slice
			regions[region] = true
		} else {
			fmt.Printf("INFO: Skipping a comment cue - ID %s \n", aTti.TimeCodeInRendered)
		}
	}

	//create the regions used
	for region_id := range regions {
		res.Head.Regions = append(res.Head.Regions, getTtmlRegions(region_id))
	}

	// debug shuffle timecodes only
	if configuration.ShuffleTimes {
		res.debugShuffleTimeCodes()
	}

	barrayRes, err := xml.MarshalIndent(res, "", "   ")
	barrayRes = []byte(xml.Header + string(barrayRes))

	//stlmerged.DebugPrint()

	return string(barrayRes), err
}
