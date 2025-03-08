package ebustl

import (
	"fmt"
	"strings"

	"github.com/bamiaux/iobit"
)

// TTI block
// 4.3. Text and Timing Information (TTI) block

type Tti struct {
	SubtitleGroupNumber  byte
	SubtitleNumber       [2]byte
	ExtensionBlockNumber byte
	CumulativeStatus     byte
	TimeCodeIn           [4]byte
	TimeCodeOut          [4]byte
	VerticalPosition     byte
	JustificationCode    byte
	CommentFlag          byte
	TextField            [112]byte

	CumulativeStatusRendered  string
	TimeCodeInRendered        string
	TimeCodeOutRendered       string
	TimeCodeInFrames          int
	TimeCodeOutFrames         int
	JustificationCodeRendered string
	CommentFlagRendered       string
	TextFieldRendered         string
	SubtitleNumberRendered    int
	ExtendedTextField         []byte
}

func (t *Tti) ToString() string {
	res := fmt.Sprintf("Subtitle Group Number: %d\n", t.SubtitleGroupNumber)
	res = res + fmt.Sprintf("Subtitle Number: %d\n", t.SubtitleNumber)
	res = res + fmt.Sprintf("Extension Block Number: %d\n", t.ExtensionBlockNumber)
	res = res + fmt.Sprintf("CumulativeStatus: %d\n", t.CumulativeStatus)
	res = res + fmt.Sprintf("CumulativeStatusRendered: %s\n", t.CumulativeStatusRendered)
	res = res + fmt.Sprintf("TimeCodeInRendered: %s\n", t.TimeCodeInRendered)
	res = res + fmt.Sprintf("TimeCodeOutRendered: %s\n", t.TimeCodeOutRendered)
	res = res + fmt.Sprintf("VerticalPosition: %d\n", t.VerticalPosition)
	res = res + fmt.Sprintf("JustificationCodeRendered: %s\n", t.JustificationCodeRendered)
	res = res + fmt.Sprintf("CommentFlagRendered: %s\n", t.CommentFlagRendered)
	//res = res + fmt.Sprintf("Text = %s", string(t.TextField[:]))
	res = res + fmt.Sprintf("Text = %s", t.TextFieldRendered)

	return res
}

func TimecodeRender(tc [4]byte) string {
	res := fmt.Sprintf("%02d:%02d:%02d:%02d", tc[0], tc[1], tc[2], tc[3])
	return res
}

func (t *Tti) DebugPrintText() string {
	res := strings.Replace(string(t.TextField[:]), "\x8f", "", -1)

	res = strings.Replace(res, "\x80", "[Italics ON]", -1)
	res = strings.Replace(res, "\x81", "[Italics OFF]", -1)
	res = strings.Replace(res, "\x82", "[Underline ON]", -1)
	res = strings.Replace(res, "\x83", "[Underline OFF]", -1)
	res = strings.Replace(res, "\x84", "[Boxing ON]", -1)
	res = strings.Replace(res, "\x85", "[Boxing OFF]", -1)

	return res + "<"
}

func getSubtitleNumber(b [2]byte) int {
	return (int(b[1]) * 256) + int(b[0])
}

func getTimeCodeframes(tc [4]byte, gsi Gsi) int {
	fps := 0
	switch gsi.DiskFormatCode {
	case "STL25.01":
		fps = 25
	case "STL30.01":
		fps = 30
	default:
		// hmmm not good should error
		//TODO actually should fail
		fps = 100
	}

	output := int(tc[0]) * 60 * 60 * fps
	output = output + int(tc[1])*60*fps
	output = output + int(tc[2])*fps
	output = output + int(tc[3])
	return output
}

func GetTti(r *iobit.Reader, gsi Gsi) (Tti, error) {
	res := Tti{}

	res.SubtitleGroupNumber = r.Byte()
	res.SubtitleNumber = [2]byte(r.Bytes(2))
	res.ExtensionBlockNumber = r.Byte()
	res.CumulativeStatus = r.Byte()
	res.TimeCodeIn = [4]byte(r.Bytes(4))
	res.TimeCodeOut = [4]byte(r.Bytes(4))
	res.VerticalPosition = r.Byte()
	res.JustificationCode = r.Byte()
	res.CommentFlag = r.Byte()
	res.TextField = [112]byte(r.Bytes(112))

	res.TimeCodeInRendered = TimecodeRender(res.TimeCodeIn)
	res.TimeCodeOutRendered = TimecodeRender(res.TimeCodeOut)

	// do timecode as frames?
	res.TimeCodeInFrames = getTimeCodeframes(res.TimeCodeIn, gsi)
	res.TimeCodeOutFrames = getTimeCodeframes(res.TimeCodeOut, gsi)

	switch res.CumulativeStatus {
	case 0:
		res.CumulativeStatusRendered = "Subtitle not part of a cumulative set"
	case 1:
		res.CumulativeStatusRendered = "First subtitle of a cumulative set"
	case 2:
		res.CumulativeStatusRendered = "Intermediate subtitle of a cumulative set"
	case 3:
		res.CumulativeStatusRendered = "Last subtitle of a cumulative set"
	default:
		res.CumulativeStatusRendered = "Unknown"
	}

	switch res.JustificationCode {
	case 0:
		res.JustificationCodeRendered = "unchanged presentation"
	case 1:
		res.JustificationCodeRendered = "left-justified text"
	case 2:
		res.JustificationCodeRendered = "centred text"
	case 3:
		res.JustificationCodeRendered = "right-justified text"
	default:
		res.JustificationCodeRendered = "Unknown"
	}

	switch res.CommentFlag {
	case 0:
		res.CommentFlagRendered = "TF contains subtitle data"
	case 1:
		res.CommentFlagRendered = "TF contains comments not intended for transmission"
	default:
		res.JustificationCodeRendered = "Unknown"
	}

	res.TextFieldRendered = res.DebugPrintText()
	res.SubtitleNumberRendered = getSubtitleNumber(res.SubtitleNumber)
	res.ExtendedTextField = res.TextField[:]

	return res, nil
}
