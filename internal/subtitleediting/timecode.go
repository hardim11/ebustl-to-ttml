package subtitleediting

import (
	"strconv"
	"strings"
)

func SplitTcDelimiters(r rune) bool {
	return r == ':' || r == '.'
}

func TcToFrames(tc string, fps int) int {
	// timecode = 00:12:51:09
	// timecode could also be
	//            00:02:49.760
	// 			  012345678901

	if (len(tc) != 11) && (len(tc) != 12) {
		return -1
	}

	parts := strings.FieldsFunc(tc, SplitTcDelimiters)
	if len(parts) != 4 {
		return -1
	}

	hr, _ := strconv.Atoi(parts[0])
	mn, _ := strconv.Atoi(parts[1])
	sc, _ := strconv.Atoi(parts[2])
	fr, _ := strconv.Atoi(parts[3])

	if len(parts[3]) == 3 {
		deci := int((float32(fr) / float32(1000) * float32(fps)))
		return (hr * 60 * 60 * fps) + (mn * 60 * fps) + (sc * fps) + deci
	} else {
		return (hr * 60 * 60 * fps) + (mn * 60 * fps) + (sc * fps) + fr
	}
}
