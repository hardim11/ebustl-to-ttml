package ebustl

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/bamiaux/iobit"
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

func (e *EbuStl) OffsetCues(offset int) {
	// loop the cues and offset
	for idx := range e.Ttis {
		e.Ttis[idx].TimeCodeInFrames = e.Ttis[idx].TimeCodeInFrames + offset
		e.Ttis[idx].TimeCodeOutFrames = e.Ttis[idx].TimeCodeOutFrames + offset

		//TODO sanity check timecodes
		// check not too small or large!

		// do the timecode
		e.Ttis[idx].TimeCodeInRendered = e.FramesToTc(e.Ttis[idx].TimeCodeInFrames)
		e.Ttis[idx].TimeCodeOutRendered = e.FramesToTc(e.Ttis[idx].TimeCodeOutFrames)

		// do the STL timecode
		e.Ttis[idx].TimeCodeIn = e.FramesToStlTimecode(e.Ttis[idx].TimeCodeInFrames)
		e.Ttis[idx].TimeCodeOut = e.FramesToStlTimecode(e.Ttis[idx].TimeCodeOutFrames)
	}
}

func ReadStlFile(filepath string) (*EbuStl, error) {
	v, err := os.ReadFile(filepath) //read the content of file
	if err != nil {
		return nil, err
	}

	return ReadStlPayload(v)
}

func ReadStlPayload(b []byte) (*EbuStl, error) {
	var res EbuStl

	res.Gsi = Gsi{}
	res.Ttis = []Tti{}

	r := iobit.NewReader(b)

	// read GSI
	res.Gsi.Read(&r)

	// skip forward 75 + 576 bytes - Spare Bytes + User Defined Area
	r.Skip((75 + 576) * 8)

	expected_length := 1024 + (res.Gsi.TotalNumberTtiBlocksInt * 128)
	if expected_length > len(b) {
		fmt.Printf("expected_length=%d\n", expected_length)
		fmt.Printf("len(b)=%d\n", len(b))
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
	//fmt.Println(len(e.Ttis))
	var ttiFirst *Tti
	for _, y := range e.Ttis {

		// if we are doing an extension block run, save the first 1 as it holds the header
		if y.ExtensionBlockNumber == 0x0 {
			ttiFirst = &y
		} else if y.ExtensionBlockNumber < 0xff {
			// add them to the first one
			ttiFirst.ExtendedTextField = append(ttiFirst.ExtendedTextField, y.TextField[:]...)
			// fmt.Println(ttiFirst.ExtendedTextField)
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
				// fmt.Print(">>>length == ")
				// fmt.Println(len(ttiFirst.ExtendedTextField))
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
			fmt.Println(aTti.TimeCodeInRendered + "," + strconv.Itoa(int(aTti.VerticalPosition)) + "," + aTti.JustificationCodeRendered + ",\"" + JustLetters(aTti.ExtendedTextField) + "\"")

		} else {
			fmt.Printf("Skipping a comment cue - ID %d \n", aTti.CommentFlag)
		}
	}
}
