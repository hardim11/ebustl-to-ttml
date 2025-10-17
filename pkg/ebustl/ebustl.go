package ebustl

import (
	"ebustl-to-ttml/internal/filehandler"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/bamiaux/iobit"
	log "github.com/sirupsen/logrus"
)

type EbuStl struct {
	Gsi  Gsi
	Ttis []Tti
}

func divmod(numerator, denominator int) (quotient, remainder int) {
	// from https://stackoverflow.com/questions/43945675/division-with-returning-quotient-and-remainder
	quotient = numerator / denominator // integer division, decimals are truncated
	remainder = numerator % denominator
	return
}

func (e *EbuStl) FramesToTc(frames int) string {
	fps := e.Gsi.Fps()
	hr, remainder := divmod(frames, (60 * 60 * fps))
	mn, remainder := divmod(remainder, (60 * fps))
	sc, fr := divmod(remainder, fps)

	return fmt.Sprintf("%02d:%02d:%02d:%02d", hr, mn, sc, fr)
}

func (e *EbuStl) FramesToStlTimecode(frames int) [4]byte {
	fps := e.Gsi.Fps()
	hr, remainder := divmod(frames, (60 * 60 * fps))
	mn, remainder := divmod(remainder, (60 * fps))
	sc, fr := divmod(remainder, fps)

	res := [4]byte{}
	res[0] = byte(hr)
	res[1] = byte(mn)
	res[2] = byte(sc)
	res[3] = byte(fr)

	return res
}

func SplitTcDelimiters(r rune) bool {
	return r == ':' || r == '.'
}

func (e *EbuStl) TcToFrames(tc string) (int, error) {
	// timecode = 00:12:51:09
	// timecode could also be
	//            00:02:49.760
	// 			  012345678901

	if (len(tc) != 11) && (len(tc) != 12) {
		return 0, errors.New("timecode string incorrect length")
	}

	parts := strings.FieldsFunc(tc, SplitTcDelimiters)
	if len(parts) != 4 {
		return 0, errors.New("timecode string incorrect sections")
	}

	hr, _ := strconv.Atoi(parts[0])
	mn, _ := strconv.Atoi(parts[1])
	sc, _ := strconv.Atoi(parts[2])
	fr, _ := strconv.Atoi(parts[3])

	fps := e.Gsi.Fps()
	if len(parts[3]) == 3 {
		deci := int((float32(fr) / float32(1000) * float32(fps)))
		return (hr * 60 * 60 * fps) + (mn * 60 * fps) + (sc * fps) + deci, nil
	} else {
		return (hr * 60 * 60 * fps) + (mn * 60 * fps) + (sc * fps) + fr, nil
	}
}

func (e *EbuStl) OffsetCues(offset int) {
	// loop the cues and offset
	for idx := range e.Ttis {
		new_in := e.Ttis[idx].TimeCodeInFrames + offset
		new_out := e.Ttis[idx].TimeCodeOutFrames + offset

		e.Ttis[idx].SetNewtimecodes(new_in, new_out, *e)
		// e.Ttis[idx].TimeCodeInFrames = e.Ttis[idx].TimeCodeInFrames + offset
		// e.Ttis[idx].TimeCodeOutFrames = e.Ttis[idx].TimeCodeOutFrames + offset

		// //TODO sanity check timecodes
		// // check not too small or large!

		// // do the timecode
		// e.Ttis[idx].TimeCodeInRendered = e.FramesToTc(e.Ttis[idx].TimeCodeInFrames)
		// e.Ttis[idx].TimeCodeOutRendered = e.FramesToTc(e.Ttis[idx].TimeCodeOutFrames)

		// // do the STL timecode
		// e.Ttis[idx].TimeCodeIn = e.FramesToStlTimecode(e.Ttis[idx].TimeCodeInFrames)
		// e.Ttis[idx].TimeCodeOut = e.FramesToStlTimecode(e.Ttis[idx].TimeCodeOutFrames)
	}
}

func ReadStlFile(filepath string) (*EbuStl, error) {
	v, err := filehandler.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	return ReadStlPayload(*v)
}

func ReadStlPayload(b []byte) (*EbuStl, error) {

	if len(b) < 1024 {
		return nil, errors.New("invalid subtitle input - stl file length < 1204 bytes (the GSI is 1024 so this is not valid)")
	}

	var res EbuStl

	res.Gsi = Gsi{}
	res.Ttis = []Tti{}

	r := iobit.NewReader(b)

	// read GSI
	err := res.Gsi.Read(&r)
	if err != nil {
		log.Error("Failed to read STL GSI block")
		return nil, err
	}

	// skip forward 75 + 576 bytes - Spare Bytes + User Defined Area
	r.Skip((75 + 576) * 8)

	expected_length := 1024 + (res.Gsi.TotalNumberTtiBlocksInt * 128)
	if expected_length > len(b) {
		log.Debugf("expected_length=%d\n", expected_length)
		log.Debugf("len(b)=%d\n", len(b))
		return nil, errors.New("file is shorter than number of cues")
	}

	// read the TTIs next
	for i := 0; i < res.Gsi.TotalNumberTtiBlocksInt; i++ {
		// read
		new_tti, _ := GetTti(&r, res.Gsi)
		res.Ttis = append(res.Ttis, new_tti)
	}

	return &res, nil
}

func (e *EbuStl) MergeExtensionBlocksTtis() (*EbuStl, error) {
	// according to the spec, the 1st TTI caries the tti header,
	// so copy the TTI to a new field
	// edit the Extension block number
	// remove the other TTIs
	res := *e
	// empty the output
	res.Ttis = []Tti{}
	var ttiFirst *Tti
	for _, y := range e.Ttis {

		// if we are doing an extension block run, save the first 1 as it holds the header
		if y.ExtensionBlockNumber == 0x0 {
			ttiFirst = &y
		} else if y.ExtensionBlockNumber < 0xff {
			// add them to the first one
			ttiFirst.ExtendedTextField = append(ttiFirst.ExtendedTextField, y.TextField[:]...)
		}

		if y.ExtensionBlockNumber == 0xff {
			// if we don't have a ttiFirst, then we just copy the tti as is
			if ttiFirst == nil {
				y.ExtendedTextField = y.TextField[:]
				res.Ttis = append(res.Ttis, y)
			} else {
				// this is part of a group of ttis, so
				// save the first one's header
				// append thje text block
				ttiFirst.ExtendedTextField = append(ttiFirst.ExtendedTextField, y.TextField[:]...)
				ttiFirst.ExtensionBlockNumber = 0xff
				foo := *ttiFirst
				res.Ttis = append(res.Ttis, foo)
				ttiFirst = nil
			}
		}

	}

	return &res, nil
}

func JustLetters(input []byte) string {
	res := ""
	for _, abyte := range input {
		if (abyte >= 32) && (abyte <= 127) {
			res = res + string(abyte)
		}
	}
	return res
}

func (e *EbuStl) DebugPrint() {
	// body
	for _, aTti := range e.Ttis {
		if aTti.CommentFlag != 1 {
			log.Debug(aTti.TimeCodeInRendered + "," + strconv.Itoa(int(aTti.VerticalPosition)) + "," + aTti.JustificationCodeRendered + ",\"" + JustLetters(aTti.ExtendedTextField) + "\"\n")

		} else {
			log.Debugf("Skipping a comment cue - ID %d \n", aTti.CommentFlag)
		}
	}
}

func (e *EbuStl) GetBetweenTimecodes(incode string, outcode string, additional_offset int) (*EbuStl, error) {
	// it is expected that the original has been merged already but what if it hasn't?

	// copy the original
	res := EbuStl{}

	res.Gsi = e.Gsi // copy the header, need to change the # subs etc later
	// reset cue list
	res.Ttis = []Tti{}

	// get part in and out frames
	incode_frames, err := res.TcToFrames(incode)
	if err != nil {
		return nil, err
	}
	outcode_frames, err := res.TcToFrames(outcode)
	if err != nil {
		return nil, err
	}

	// loop the source subs and copy in the ones we want to keep
	for _, aCue := range e.Ttis {

		// evaluate if in the range
		if (aCue.TimeCodeInFrames >= incode_frames) && (aCue.TimeCodeInFrames < outcode_frames) {
			// if it is, copy and offset the times to match the part
			newCue := aCue
			// set values
			// adjust timecodes
			new_begin := newCue.TimeCodeInFrames - incode_frames + additional_offset
			new_end := newCue.TimeCodeOutFrames - incode_frames + additional_offset
			newCue.SetNewtimecodes(new_begin, new_end, res)
			// add
			res.Ttis = append(res.Ttis, newCue)
		}
	}

	// TODO update the Gsi
	log.Warn("TODO update the Gsi")

	// return
	return &res, nil
}
