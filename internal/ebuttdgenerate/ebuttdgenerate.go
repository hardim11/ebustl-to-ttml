package ebuttdgenerate

import (
	"bytes"
	ebustl "ebustl-to-ttml/internal/ebustl"
	"ebustl-to-ttml/internal/ptr"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"fmt"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

const FPS = 25

const COLOUR_BLACK = "#000000"
const COLOUR_RED = "#ff0000"
const COLOUR_GREEN = "#00ff00"
const COLOUR_YELLOW = "#ffff00"
const COLOUR_BLUE = "#0000ff"
const COLOUR_MAGENTA = "#ff00ff"
const COLOUR_CYAN = "#00ffff"
const COLOUR_WHITE = "#ffffff"
const TELETEXT_WIDTH = 40

const TTS_FONTSIZE = "160%"
const TTS_LINEHEIGHT = "120%"

type RegionStruct struct {
	LineNumber int
	RowCount   int
	LeftPad    int
	Width      int
}

type SpanOrLineBreak struct {
	IsLineBreak      bool
	IsPadding        bool
	BackgroundColour string
	ForegroundColour string
	ContentString    string
	SpanStyle        string
	DebugNotes       []string
	//PreserveSpace    bool //TODO should I keep this and use for internal padding?
}

func SpanOrLineBreakNew() SpanOrLineBreak {
	return SpanOrLineBreak{
		IsLineBreak:      false,
		IsPadding:        false,
		BackgroundColour: COLOUR_BLACK,
		ForegroundColour: COLOUR_WHITE,
		ContentString:    "",
	}
}

func (s *SpanOrLineBreak) ToString() string {
	res := fmt.Sprintf("IsLineBreak %t, IsPadding %t, BackgroundColour %s, ForegroundColour %s - \n>%s<", s.IsLineBreak, s.IsPadding, s.BackgroundColour, s.ForegroundColour, s.ContentString)
	return res
}

type SubtitleExtent struct {
	BoxRight  int
	BoxLeft   int
	BoxHeight int
	BoxTop    int
}

func (e *SubtitleExtent) ToString() string {
	return fmt.Sprintf("Extent top=%d, height=%d, left=%d, right=%d", e.BoxTop, e.BoxHeight, e.BoxLeft, e.BoxRight)
}

type NormalisedPara struct {
	SpansAndBreaks []SpanOrLineBreak
	//JustificationCode int
	//Extent SubtitleExtent
}

func (n *NormalisedPara) TrimSpacesWithBr(input string) string {
	res := strings.TrimSpace(input)
	if strings.HasSuffix(res, "<br />") {
		res = strings.TrimSuffix(res, "<br />")
		res = strings.TrimSpace(res)
		res = res + "<br />"
	}

	return res
}

func (n *NormalisedPara) ConvertToXMLEbuTtD() (string, error) {
	res := "\n"
	for _, y := range n.SpansAndBreaks {
		if y.IsLineBreak {
			res = res + "<br />"
			res = res + "\n"
		} else if y.IsPadding {
			res = res + "<span xml:space=\"preserve\">" + XmlEscapeText(y.ContentString) + "</span>"
			res = res + "\n"
		} else {
			if y.BackgroundColour == "" || y.ForegroundColour == "" {
				return "", errors.New("failed to convert Paragraph to XML as foreground / background colours not set - NormalisedPara.ConvertToXML, " + y.ContentString)
			}
			//res = res + "<span tts:backgroundColor=\"" + y.BackgroundColour + "\" tts:color=\"" + y.ForegroundColour + "\">"
			res = res + "<span style=\"" + y.SpanStyle + "\">"
			res = res + XmlEscapeText(y.ContentString)
			res = res + "</span>"
			res = res + "\n"
		}
	}

	return res, nil
}

func DebugPrintByteArrray(byteArry []byte) string {
	return "<!-- " + hex.EncodeToString([]byte(byteArry)) + " -->"
}

func getTtmlDefaultStyle() []EbuTtDOutStyle {
	res := []EbuTtDOutStyle{}
	// defaultStyle
	dStyleAttributes := StyleAttributes{
		TTMLFontFamily:    ptr.String("monospaceSansSerif"), // "proportionalSansSerif, SansSerif, default"
		TTMLFontStyle:     ptr.String("normal"),
		TTMLTextAlign:     ptr.String("left"),
		TTMLWrapOption:    ptr.String("noWrap"),
		EbuTtsLinePadding: ptr.String("0.5c"),
	}
	var ttmlTextStyle = EbuTtDOutStyle{EbuTtDOutHeader: EbuTtDOutHeader{
		ID:                       "defaultStyle",
		EbuTtDOutStyleAttributes: ebuTtDOutStyleAttributesFromStyleAttributes(&dStyleAttributes),
	}}
	res = append(res, ttmlTextStyle)

	return res
}

func getTtmlDefaultRegions(debug_region bool) []EbuTtDOutRegion {
	res := []EbuTtDOutRegion{}

	if debug_region {
		regiondebug := EbuTtDOutRegion{
			EbuTtDOutHeader: EbuTtDOutHeader{
				ID: "debug_region",
				EbuTtDOutStyleAttributes: EbuTtDOutStyleAttributes{
					Extent: ptr.String("80% 16.67%"),
					Origin: ptr.String("10% 79.17%"),
				},
			},
			Comment: "this is only for debugging use",
		}
		res = append(res, regiondebug)
	}

	return res
}

func (r *RegionStruct) decodeRegionNameString(RegionId string) {
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

func getRegionFromNameString(RegionId string) RegionStruct {
	res := RegionStruct{}
	res.decodeRegionNameString(RegionId)
	return res
}

func getEbuTtDRegions(region_id string) EbuTtDOutRegion {
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


		I think this ^ is WRONG for horizontal at least
		I suggest origin is 10% always, extent it 80% always, end of.

		and I think height should be offset too to include a border
		again, I think 10% top and 80% extent * line number


		height
		80 / 24 rows
	*/

	// var regiondata RegionStruct
	// regiondata.DecodeRegionNameString(region_id)
	regiondata := getRegionFromNameString(region_id)

	origin_x := 10.00                                             //float64(regiondata.LeftPad) * 2.5
	origin_y := (float64(regiondata.LineNumber-1) * 3.50) + 10.00 // float64(regiondata.LineNumber-1) * 4.165
	extent_x := 80.00                                             //float64(regiondata.Width) * 2.5
	extent_y := float64(regiondata.RowCount) * 8.33               // float64(regiondata.RowCount) * 8.33

	origin := fmt.Sprintf("%0.2f%% %0.2f%%", origin_x, origin_y)
	extent := fmt.Sprintf("%0.2f%% %0.2f%%", extent_x, extent_y)

	res := EbuTtDOutRegion{
		EbuTtDOutHeader: EbuTtDOutHeader{
			ID: region_id,
			EbuTtDOutStyleAttributes: EbuTtDOutStyleAttributes{
				Origin: &origin,
				Extent: &extent,
			},
		},
	}
	return res
}

func getEbuTtDStles(style_id string) (*EbuTtDOutStyle, error) {
	// lets go with "style[p|s]-{fg}-{bg}({space}-options---)"
	// where stylep == style for Para and imples defaults
	// and styles == style for Span
	res := EbuTtDOutStyle{}
	parts := strings.Split(style_id, "-")
	if len(parts) < 1 {
		return nil, errors.New("no style found")
	}

	// TODO apply <p> / <span> defaults?
	//
	// <P> 		-> <style xml:id="style-3" tts:fontSize="160%" tts:lineHeight="120%" tts:textAlign="center" itts:fillLineGap="false" />
	// <SPAN> 	-> <style xml:id="style-4" tts:backgroundColor="#000000" tts:color="#FFFF00" />
	if len(parts) >= 1 {
		res.ID = style_id
	}
	if strings.HasPrefix(style_id, "styles") {
		// do a span
		if len(parts) >= 2 {
			fgcolor := "#" + parts[1]
			res.Color = &fgcolor
		}
		if len(parts) >= 3 {
			bgcolor := "#" + parts[2]
			res.BackgroundColor = &bgcolor
		}
		if len(parts) >= 4 {
			if parts[3] != "" {
				space := parts[3]
				res.Space = &space
			}
		}
	} else if strings.HasPrefix(style_id, "stylep") {
		// do a paragraph
		res.FontSize = ptr.String(TTS_FONTSIZE)
		res.LineHeight = ptr.String(TTS_LINEHEIGHT)

		// do alignment based on centring
		if len(parts) >= 2 {
			align := parts[1]
			res.TextAlign = &align
		}
	}
	return &res, nil
}

func XmlEscapeText(input string) string {
	// // EscapeText writes to w the properly escaped XML equivalent
	// // of the plain text data s.
	// func EscapeText(w io.Writer, s []byte) error {
	// 	return escapeText(w, s, true)
	// }
	// TODO I don't like mixing serialisation and escaped XML strings.
	//fmt.Printf("XmlEscapeText % x %s\n", input, input)
	buf := new(bytes.Buffer)
	xml.EscapeText(buf, []byte(input))
	return buf.String()
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

func getSubtitlePara(txtraster ebustl.TeletextRaster, fixedalignment bool, configuration EbuTtDConvertConfiguration) (string, int, int, int, []string, error) {
	// returns
	// 		xml for the paragraph
	//		row count
	//		left padding
	//		char count (including code chars)
	//      styles []string
	//		error

	styles := []string{}
	tmpSpans := []SpanOrLineBreak{}

	// loop the active lines
	last_pixel := ebustl.TeletextPixel{
		ForegroundColour: "",
		BackgroundColour: "",
	}
	for row_index, row := range txtraster.Lines {
		if row.ActiveLine {

			current_span := SpanOrLineBreak{}
			if fixedalignment {
				// build in the padding

				// check against negative repeat
				if row.StartBoxIndex > 0 {
					aspan := SpanOrLineBreak{
						IsPadding:     true,
						ContentString: strings.Repeat(" ", row.StartBoxIndex-1),
					}
					tmpSpans = append(tmpSpans, aspan)
				}
			}

			// debug stuff
			if configuration.Debug {
				log.Debugf("ROW %d, row.StartBoxIndex %d, row.EndBoxIndex %d\n", row_index, row.StartBoxIndex, row.EndBoxIndex)
			}
			// try to detect the odd box only ones
			if row.StartBoxIndex < 0 {
				log.Warn("The start end box is not defined")
				// try to find the newbackground ?
			}
			if configuration.Debug {
				log.Debugf("row.EndBoxIndex = %d", row.EndBoxIndex)
			}
			truncated_endboxindex := row.EndBoxIndex
			if truncated_endboxindex > TELETEXT_WIDTH+4 {
				truncated_endboxindex = TELETEXT_WIDTH + 4
				log.Warn("row.EndBoxIndex truncated")
			}
			for pixel_index := row.StartBoxIndex; pixel_index < truncated_endboxindex; pixel_index++ {
				this_pixel := row.Pixels[pixel_index]
				if configuration.Debug {
					log.Debug(this_pixel.ToString())
				}
				if !this_pixel.IsVisibleChar {
					current_span.ContentString = current_span.ContentString + " "
					//continue
				}

				// check for colour change
				if last_pixel.BackgroundColour != this_pixel.BackgroundColour || last_pixel.ForegroundColour != this_pixel.ForegroundColour {
					// save the old one
					if current_span.BackgroundColour != "" && current_span.ForegroundColour != "" {
						tmpSpans = append(tmpSpans, current_span)
					}
					// create a new span
					current_span = SpanOrLineBreak{
						ForegroundColour: this_pixel.ForegroundColour,
						BackgroundColour: this_pixel.BackgroundColour,
					}
				}

				current_span.ContentString = current_span.ContentString + this_pixel.Character

				// reset last pixel
				last_pixel = this_pixel
			}

			if current_span.BackgroundColour != "NONE" && current_span.ForegroundColour != "NONE" {
				tmpSpans = append(tmpSpans, current_span)
			}

			abreak := SpanOrLineBreak{
				IsLineBreak: true,
			}
			tmpSpans = append(tmpSpans, abreak)
		}
	}

	// remove the last <br />
	if len(tmpSpans) > 0 {
		if tmpSpans[len(tmpSpans)-1].IsLineBreak {
			tmpSpans = tmpSpans[:len(tmpSpans)-1]
		}
	}

	res := NormalisedPara{
		SpansAndBreaks: []SpanOrLineBreak{},
	}

	// sanity check that there are no spans with missing foreground / background (caused by trailing new colour)
	for _, span := range tmpSpans {
		if (span.BackgroundColour == "" || span.ForegroundColour == "") && !span.IsLineBreak && !span.IsPadding && strings.TrimSpace(span.ContentString) == "" {
			log.Warn("getSubtitlePara: sanity check removing empty span (trailing new colour?)")
		} else {
			// add styles?
			spanToAdd := span
			// "style[p|s].{fg}.{bg}({space}.options...)"
			if !span.IsLineBreak && !span.IsPadding {
				styleName := fmt.Sprintf(
					"styles-%s-%s",
					strings.ReplaceAll(span.ForegroundColour, "#", ""),
					strings.ReplaceAll(span.BackgroundColour, "#", ""),
				)
				styles = append(styles, styleName)
				spanToAdd.SpanStyle = styleName
			}

			res.SpansAndBreaks = append(res.SpansAndBreaks, spanToAdd)
		}
	}

	// convert spans to XML
	res_string, err := res.ConvertToXMLEbuTtD()
	if err != nil {
		log.Error("getSubtitlePara raised error converting to XML")
		return "", -1, -1, -1, styles, err
	}

	row_count := txtraster.ActiveLineCount
	left_pad := 0                // is this right //TODO
	char_count := TELETEXT_WIDTH // is this right //TODO
	return res_string, row_count, left_pad, char_count, styles, nil
}

func getSubtitle(tti ebustl.Tti, codepage string, configuration EbuTtDConvertConfiguration) (*EbuTtDOutSubtitle, string, []string, error) {
	styles := []string{}
	// convert to a teletext raster
	ttxtraster, err := ebustl.CreateTeletextRasterFromTti(tti, codepage, configuration.TruncateOversizedLines, configuration.Debug, configuration.IgnoreInvalidDiacritical)
	if err != nil {
		return nil, "", styles, err
	}

	if ttxtraster.ActiveLineCount < 0 && ttxtraster.FirstActiveLine < 0 {
		log.Warn("Skipping a TTI as it contains no date")
		return nil, "", styles, nil
	}

	if configuration.Debug {
		log.Debug(ttxtraster.PrintToConsole())
	}

	// convert the raster to paragraphs
	sub_para, rows, left_pad, char_count, styles, err := getSubtitlePara(*ttxtraster, (tti.JustificationCode == 0), configuration)
	if configuration.Debug {
		log.Debugf("sub_para = %#v\n", sub_para)
	}
	if err != nil {
		log.Error("getSubtitle raised error for subtitle cue " + tti.TimeCodeInRendered)
		return nil, "", styles, err
	}
	region := "region." + strconv.Itoa(int(tti.VerticalPosition)) + "." + strconv.Itoa(rows) + "." + strconv.Itoa(left_pad) + "." + strconv.Itoa(char_count)

	ttsTextAlign := "left" // none is covered here (tti.JustificationCode == 0) too
	switch tti.JustificationCode {
	case 1:
		ttsTextAlign = "left"
	case 2:
		ttsTextAlign = "center"
	case 3:
		ttsTextAlign = "right"
	}

	paraStyle := "stylep-" + ttsTextAlign
	res := EbuTtDOutSubtitle{
		Region: region,
		Begin:  FramesToTcMs(int64(tti.TimeCodeInFrames)),
		End:    FramesToTcMs(int64(tti.TimeCodeOutFrames)),
		Text:   sub_para,
		Style:  paraStyle,
	}
	// ID required for EBU
	res.XmlID = fmt.Sprintf("p%d", tti.SubtitleNumberRendered)

	styles = append(styles, paraStyle)
	return &res, region, styles, nil
}

func secondsToTc(seconds int64) string {
	hr, remainder := divmod(seconds, (60 * 60))
	mn, sc := divmod(remainder, (60))
	return fmt.Sprintf("%02d:%02d:%02d", hr, mn, sc)
}

func (t *EbuTtDOut) debugShuffleTimeCodes() {
	shuffle_spacing := 2
	seconds := shuffle_spacing
	for idx := range t.Body.Div.Subtitles {
		t.Body.Div.Subtitles[idx].Begin = secondsToTc(int64(seconds)) + ".000"
		t.Body.Div.Subtitles[idx].End = secondsToTc(int64(seconds+(shuffle_spacing-1))) + ".500"
		seconds = seconds + shuffle_spacing
	}
}

func CreateTtml(stl ebustl.EbuStl, comment string, config *EbuTtDConvertConfiguration) (string, error) {

	//configuration = *config

	// merge TTI's - where a cue is split over multiple TTI's convert to just one
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

	res := EbuTtDOut{}
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
	res.XMLNamespaceEbuTt = "urn:ebu:tt:style"
	res.XMLNamespaceEbuTm = "urn:ebu:tt:metadata"

	res.Lang = getxmlLanguageCode(stl.Gsi.LanguageCode)
	res.CellRsolution = config.CellRsolution
	res.Timebase = "media"

	res.Head.Metadata = &EbuTtDOutMetadata{}
	res.Head.Metadata.Title = stlmerged.Gsi.OriginalProgrammeTitle
	res.Head.Metadata.Description = stlmerged.Gsi.SubtitleListReferenceCode

	fixedStyles := getTtmlDefaultStyle()
	res.Head.Styles = append(res.Head.Styles, fixedStyles[:]...)

	res.Head.Regions = []EbuTtDOutRegion{}
	res.Head.Regions = append(res.Head.Regions, getTtmlDefaultRegions(config.Debug)...)

	res.Body.Style = "defaultStyle"

	// create empty regions list
	regions := map[string]bool{}
	styles := map[string]bool{}

	// deal with the actual subtitles
	for _, aTti := range stlmerged.Ttis {
		if aTti.CommentFlag != 1 {
			// add a para
			para, qregion, qstyles, err := getSubtitle(aTti, stlmerged.Gsi.CharacterCodeTable, *config)
			if err != nil {
				log.Errorf("error in CreateTtml for cue incode %s, error %s", aTti.TimeCodeInRendered, err.Error())
				return "", err
			}
			if para != nil {
				res.Body.Div.Subtitles = append(res.Body.Div.Subtitles, *para)

				// kludge way to have unique regions and styles slice
				regions[qregion] = true
				for _, astyle := range qstyles {
					styles[astyle] = true
				}
			}
		} else {
			log.Infof("INFO: Skipping a comment cue - ID %s \n", aTti.TimeCodeInRendered)
		}
	}

	//create the regions used
	for region_id := range regions {
		res.Head.Regions = append(res.Head.Regions, getEbuTtDRegions(region_id))
	}

	// do styles
	for style_id := range styles {
		// add the style to the head
		style, err := getEbuTtDStles(style_id)
		if err != nil {
			log.Errorf("error in CreateTtml decoding style for style %s, error %s", style_id, err.Error())
			return "", err
		}
		res.Head.Styles = append(res.Head.Styles, *style)
	}

	// debug shuffle timecodes only
	if config.ShuffleTimes {
		res.debugShuffleTimeCodes()
	}

	barrayRes, err := xml.MarshalIndent(res, "", "   ")
	//barrayRes, err := xml.Marshal(res)

	barrayRes = []byte(xml.Header + string(barrayRes))

	//stlmerged.DebugPrint()

	//tmp := string(barrayRes)
	//tmp = strings.TrimSpace(xmlfmt.FormatXML(tmp, "", "  "))

	return string(barrayRes), err
}
