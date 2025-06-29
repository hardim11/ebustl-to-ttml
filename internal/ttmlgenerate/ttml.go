package ttmlgenerate

// copied and adaopted from https://github.com/asticode/go-astisub

import (
	"encoding/xml"
)

// https://www.w3.org/TR/ttaf1-dfxp/
// http://www.skynav.com:8080/ttv/check
// https://www.speechpad.com/captions/ttml

type Justification int

var (
	JustificationUnchanged = Justification(1)
	JustificationLeft      = Justification(2)
	JustificationCentered  = Justification(3)
	JustificationRight     = Justification(4)
)

// StyleAttributes represents style attributes
type StyleAttributes struct {
	TTMLBackgroundColor *string // https://htmlcolorcodes.com/fr/
	TTMLColor           *string
	TTMLDirection       *string
	TTMLDisplay         *string
	TTMLDisplayAlign    *string
	TTMLExtent          *string
	TTMLFontFamily      *string
	TTMLFontSize        *string
	TTMLFontStyle       *string
	TTMLFontWeight      *string
	TTMLLineHeight      *string
	TTMLOpacity         *string
	TTMLOrigin          *string
	TTMLOverflow        *string
	TTMLPadding         *string
	TTMLShowBackground  *string
	TTMLTextAlign       *string
	TTMLTextDecoration  *string
	TTMLTextOutline     *string
	TTMLUnicodeBidi     *string
	TTMLVisibility      *string
	TTMLWrapOption      *string
	TTMLWritingMode     *string
	TTMLZIndex          *int
	TTMLFillLineGap     *string
	TTMLSpace           *string
	TTMLForcedDisplay   *string
}

// TTMLOut represents an output TTML that must be marshaled
// We split it from the input TTML as this time we'll add strict namespaces
type TTMLOut struct {
	Comment string `xml:",comment"`
	Head    Head   `xml:"head"`

	Body                   TTMLOutBody `xml:"body"`
	XMLName                xml.Name    `xml:"http://www.w3.org/ns/ttml tt"`
	XMLNamespaceTTS        string      `xml:"xmlns:tts,attr"`
	XMLNamespaceTTP        string      `xml:"xmlns:ttp,attr"`
	XMLNamespaceTTM        string      `xml:"xmlns:ttm,attr"`
	XMLNamespaceSmpte      string      `xml:"xmlns:smpte,attr,omitempty"`
	XMLNamespaceIMSC       string      `xml:"xmlns:itts,attr,omitempty"`
	XMLNamespaceITTP       string      `xml:"xmlns:ittp,attr,omitempty"`
	XMLNamespaceEbuTt      string      `xml:"xmlns:ebutts,attr,omitempty"`
	Space                  string      `xml:"xml:space,attr,omitempty"`
	Lang                   string      `xml:"xml:lang,attr,omitempty"`
	FrameRate              string      `xml:"ttp:frameRate,attr,omitempty"`
	FrameRateMultiplier    string      `xml:"ttp:frameRateMultiplier,attr,omitempty"`
	Timebase               string      `xml:"ttp:timeBase,attr,omitempty"`
	CellRsolution          string      `xml:"ttp:cellResolution,attr,omitempty"`
	ActiveArea             string      `xml:"ittp:activeArea,attr,omitempty"`
	ProgressivelyDecodable string      `xml:"ittp:progressivelyDecodable,attr,omitempty"`
}

type Head struct {
	//Information HeadSmpte        `xml:"http://www.smpte-ra.org/schemas/2052-1/2010/smpte-tt information,omitempty`
	Metadata *TTMLOutMetadata `xml:"metadata,omitempty"`
	Styles   []TTMLOutStyle   `xml:"styling>style,omitempty"` //!\\ Order is important! Keep Styling above Layout
	Regions  []TTMLOutRegion  `xml:"layout>region,omitempty"`
}

type TTMLOutBody struct {
	Style string     `xml:"style,attr,omitempty"`
	Div   TTMLOutDiv `xml:"div"`
}

type TTMLOutDiv struct {
	Subtitles []TTMLOutSubtitle `xml:"p,omitempty"`
}

// TTMLOutMetadata represents an output TTML Metadata
type TTMLOutMetadata struct {
	Copyright   string `xml:"ttm:copyright,omitempty"`
	Title       string `xml:"ttm:title,omitempty"`
	Description string `xml:"ttm:desc,omitempty"`
}

// TTMLOutStyleAttributes represents output TTML style attributes
type TTMLOutStyleAttributes struct {
	Display         *string `xml:"tts:display,attr,omitempty"`
	DisplayAlign    *string `xml:"tts:displayAlign,attr,omitempty"`
	BackgroundColor *string `xml:"tts:backgroundColor,attr,omitempty"`
	Color           *string `xml:"tts:color,attr,omitempty"`
	Direction       *string `xml:"tts:direction,attr,omitempty"`
	Origin          *string `xml:"tts:origin,attr,omitempty"`
	Extent          *string `xml:"tts:extent,attr,omitempty"`
	FontFamily      *string `xml:"tts:fontFamily,attr,omitempty"`
	FontSize        *string `xml:"tts:fontSize,attr,omitempty"`
	FontStyle       *string `xml:"tts:fontStyle,attr,omitempty"`
	FontWeight      *string `xml:"tts:fontWeight,attr,omitempty"`
	LineHeight      *string `xml:"tts:lineHeight,attr,omitempty"`
	Opacity         *string `xml:"tts:opacity,attr,omitempty"`

	Overflow       *string `xml:"tts:overflow,attr,omitempty"`
	Padding        *string `xml:"tts:padding,attr,omitempty"`
	ShowBackground *string `xml:"tts:showBackground,attr,omitempty"`
	TextAlign      *string `xml:"tts:textAlign,attr,omitempty"`
	TextDecoration *string `xml:"tts:textDecoration,attr,omitempty"`
	TextOutline    *string `xml:"tts:textOutline,attr,omitempty"`
	UnicodeBidi    *string `xml:"tts:unicodeBidi,attr,omitempty"`
	Visibility     *string `xml:"tts:visibility,attr,omitempty"`
	WrapOption     *string `xml:"tts:wrapOption,attr,omitempty"`
	WritingMode    *string `xml:"tts:writingMode,attr,omitempty"`
	ZIndex         *int    `xml:"tts:zIndex,attr,omitempty"`
	FillLineGap    *string `xml:"itts:fillLineGap,attr,omitempty"`
	Space          *string `xml:"tts:space,attr,omitempty"`
	ForcedDisplay  *string `xml:"itts:forcedDisplay,attr,omitempty"`
	PreserveSpace  *string `xml:"xml:space,attr,omitempty"`
}

// ttmlOutStyleAttributesFromStyleAttributes converts StyleAttributes into a TTMLOutStyleAttributes
func ttmlOutStyleAttributesFromStyleAttributes(s *StyleAttributes) TTMLOutStyleAttributes {
	if s == nil {
		return TTMLOutStyleAttributes{}
	}
	return TTMLOutStyleAttributes{
		BackgroundColor: s.TTMLBackgroundColor,
		Color:           s.TTMLColor,
		Direction:       s.TTMLDirection,
		Display:         s.TTMLDisplay,
		DisplayAlign:    s.TTMLDisplayAlign,
		Extent:          s.TTMLExtent,
		FontFamily:      s.TTMLFontFamily,
		FontSize:        s.TTMLFontSize,
		FontStyle:       s.TTMLFontStyle,
		FontWeight:      s.TTMLFontWeight,
		LineHeight:      s.TTMLLineHeight,
		Opacity:         s.TTMLOpacity,
		Origin:          s.TTMLOrigin,
		Overflow:        s.TTMLOverflow,
		Padding:         s.TTMLPadding,
		ShowBackground:  s.TTMLShowBackground,
		TextAlign:       s.TTMLTextAlign,
		TextDecoration:  s.TTMLTextDecoration,
		TextOutline:     s.TTMLTextOutline,
		UnicodeBidi:     s.TTMLUnicodeBidi,
		Visibility:      s.TTMLVisibility,
		WrapOption:      s.TTMLWrapOption,
		WritingMode:     s.TTMLWritingMode,
		ZIndex:          s.TTMLZIndex,
		FillLineGap:     s.TTMLFillLineGap,
		Space:           s.TTMLSpace,
		ForcedDisplay:   s.TTMLForcedDisplay,
	}
}

// TTMLOutHeader represents an output TTML header
type TTMLOutHeader struct {
	ID    string `xml:"xml:id,attr,omitempty"`
	Style string `xml:"style,attr,omitempty"`
	TTMLOutStyleAttributes
}

// TTMLOutRegion represents an output TTML region
type TTMLOutRegion struct {
	Comment string `xml:",comment"`
	TTMLOutHeader
	XMLName xml.Name `xml:"region"`
}

// TTMLOutStyle represents an output TTML style
type TTMLOutStyle struct {
	TTMLOutHeader
	XMLName xml.Name `xml:"style"`
}

// TTMLOutSubtitle represents an output TTML subtitle
type TTMLOutSubtitle struct {
	Begin  string `xml:"begin,attr"` // was TTMLOutDuration
	End    string `xml:"end,attr"`   // was TTMLOutDuration
	Region string `xml:"region,attr,omitempty"`
	Style  string `xml:"style,attr,omitempty"`
	ID     string `xml:"id,attr,omitempty"`
	Items  []TTMLOutItem

	Text string `xml:",innerxml"` // `xml:",chardata"`
	TTMLOutStyleAttributes
}

func (s *TTMLOutSubtitle) ToString() string {
	res, _ := xml.MarshalIndent(s, "", "   ")
	return string(res)
}

// TTMLOutItem represents an output TTML Item
type TTMLOutItem struct {
	Style string `xml:"style,attr,omitempty"`
	Text  string `xml:",chardata"`
	TTMLOutStyleAttributes
	XMLName xml.Name
}

// WriteToTTMLOptions represents TTML write options.
type WriteToTTMLOptions struct {
	Indent string // Default is 4 spaces.
}

// WriteToTTMLOption represents a WriteToTTML option.
type WriteToTTMLOption func(o *WriteToTTMLOptions)

// WriteToTTMLWithIndentOption sets the indent option.
func WriteToTTMLWithIndentOption(indent string) WriteToTTMLOption {
	return func(o *WriteToTTMLOptions) {
		o.Indent = indent
	}
}
