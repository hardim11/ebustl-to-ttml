package ttmlgenerate

import (
	ebustl "ebustl-to-ttml/internal/ebustl"
	"fmt"
	"testing"
)

// func TestGetNextCharFromString(t *testing.T) {

// 	foo := "1234"
// 	res := getNextCharFromString(foo, 0)

// 	// "2"
// 	fmt.Printf("TestgetNextCharFromString got %d\n", res)
// 	if res != 0x32 {
// 		t.Fatalf("ERROR: TestgetNextCharFromString res=%d", res)
// 	}
// }

// func TestGetNextCharFromString2(t *testing.T) {

// 	foo := "1234"
// 	res := getNextCharFromString(foo, 3)

// 	// "2"
// 	fmt.Printf("TestgetNextCharFromString2 got %d\n", res)
// 	if res != 0 {
// 		t.Fatalf("ERROR: TestgetNextCharFromString2 res=%d", res)
// 	}
// }

// func TestTrimTextString(t *testing.T) {
// 	ExtendedTextField := [...]byte{
// 		0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x0D, 0x01, 0x1D, 0x07, 0x0B, 0x0B,
// 		0x4C, 0x41, 0x55, 0x47, 0x48, 0x54, 0x45, 0x52, 0x0A, 0x0A, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20,
// 		0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F,
// 	}

// 	foo := getTextField(ExtendedTextField[:])
// 	input := trimTextString(foo)
// 	fmt.Println(input)
// }

const TRUNCATEOVERSIZEDLINES = true
const DEBUG = false
const IGNOREINVALIDDIACRITICAL = true

func GetConfig() TtmlConvertConfiguration {
	configuration := TtmlConvertConfiguration{
		TruncateOversizedLines: TRUNCATEOVERSIZEDLINES,
		Debug:                  DEBUG,
	}
	return configuration
}

func Test1(t *testing.T) {
	// returns
	// 		xml for the paragraph
	//		row count
	//		left padding
	//		char count (including code chars)
	//		error

	/*
		    	<p begin="485679:35:17.920" end="485679:35:18.240" region="region-50" tts:fontSize="200%" tts:lineHeight="120%" tts:textAlign="center">
					<span tts:backgroundColor="#FF0000" tts:color="#FFFFFF">LAUGHTER</span>
				</p>
			LAUGHTER        (white on red background) 10:38:06:07
	*/

	ExtendedTextField := [...]byte{
		0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x0D, 0x01, 0x1D, 0x07, 0x0B, 0x0B,
		0x4C, 0x41, 0x55, 0x47, 0x48, 0x54, 0x45, 0x52, 0x0A, 0x0A, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20,
		0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F,
	}
	tti := ebustl.Tti{
		ExtendedTextField: ExtendedTextField[:],
		VerticalPosition:  22,
		JustificationCode: 2,
	}
	codepage := "00"
	fixedalignment := (tti.JustificationCode == 0)
	res, _ := ebustl.CreateTeletextRasterFromTti(tti, codepage, TRUNCATEOVERSIZEDLINES, DEBUG, IGNOREINVALIDDIACRITICAL)
	xmlstring, row_count, left_padding, char_count, err := getSubtitlePara(*res, fixedalignment, GetConfig())
	if err != nil {
		t.Fatalf("ERROR: Test1 err=%s", err.Error())
	}
	fmt.Println(xmlstring)
	fmt.Printf("row_count=%d, left_padding=%d, char_count=%d", row_count, left_padding, char_count)
}

func Test2(t *testing.T) {
	// returns
	// 		xml for the paragraph
	//		row count
	//		left padding
	//		char count (including code chars)
	//		error

	/*
		<span tts:backgroundColor="#000000" tts:color="#FFFFFF">Miss Lemon has yet</span>
		<br />
		<span tts:backgroundColor="#000000" tts:color="#FFFFFF">to see a performance. </span>
		<span tts:backgroundColor="#000000" tts:color="#00FF00">Shh!</span>
	*/

	ExtendedTextField := [...]byte{
		0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x0D, 0x0B, 0x0B, 0x4D, 0x69, 0x73, 0x73, 0x20,
		0x4C, 0x65, 0x6D, 0x6F, 0x6E, 0x20, 0x68, 0x61, 0x73, 0x20, 0x79, 0x65, 0x74, 0x0A, 0x0A, 0x20,
		0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x8A, 0x8A, 0x20, 0x20, 0x20, 0x20, 0x0D, 0x0B,
		0x0B, 0x74, 0x6F, 0x20, 0x73, 0x65, 0x65, 0x20, 0x61, 0x20, 0x70, 0x65, 0x72, 0x66, 0x6F, 0x72,
		0x6D, 0x61, 0x6E, 0x63, 0x65, 0x2E, 0x02, 0x53, 0x68, 0x68, 0x21, 0x0A, 0x0A, 0x20, 0x20, 0x20,
		0x20, 0x20, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F,
		0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F,
	}
	tti := ebustl.Tti{
		ExtendedTextField: ExtendedTextField[:],
		VerticalPosition:  22,
		JustificationCode: 2,
	}
	codepage := "00"
	fixedalignment := (tti.JustificationCode == 0)
	res, _ := ebustl.CreateTeletextRasterFromTti(tti, codepage, TRUNCATEOVERSIZEDLINES, DEBUG, IGNOREINVALIDDIACRITICAL)
	fmt.Println(res.PrintToConsole())
	xmlstring, row_count, left_padding, char_count, err := getSubtitlePara(*res, fixedalignment, GetConfig())
	if err != nil {
		t.Fatalf("ERROR: Test1 err=%s", err.Error())
	}
	fmt.Println(xmlstring)
	fmt.Printf("row_count=%d, left_padding=%d, char_count=%d", row_count, left_padding, char_count)
}

func Test3(t *testing.T) {
	/*
				Matt Simpon's extended Character Code stuff

		Red on Black.Blue on Black.
		Red on Black.Blue on Black.


	*/
	ExtendedTextField := [...]byte{
		0x0D, 0x01, 0x0B, 0x0B, 0x52, 0x65, 0x64, 0x20, 0x6F, 0x6E, 0x20, 0x42, 0x6C, 0x61, 0x63, 0x6B,
		0x2E, 0x04, 0x42, 0x6C, 0x75, 0x65, 0x20, 0x6F, 0x6E, 0x20, 0x42, 0x6C, 0x61, 0x63, 0x6B, 0x2E,
		0x0A, 0x0A, 0x8A, 0x8A, 0x0D, 0x01, 0x0B, 0x0B, 0x52, 0x65, 0x64, 0x20, 0x6F, 0x6E, 0x20, 0x42,
		0x6C, 0x61, 0x63, 0x6B, 0x2E, 0x04, 0x42, 0x6C, 0x75, 0x65, 0x20, 0x6F, 0x6E, 0x20, 0x42, 0x6C,
		0x61, 0x63, 0x6B, 0x2E, 0x0A, 0x0A, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F,
		0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F,
		0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F,
	}
	tti := ebustl.Tti{
		ExtendedTextField: ExtendedTextField[:],
		VerticalPosition:  22,
		JustificationCode: 2,
	}
	codepage := "00"
	fixedalignment := (tti.JustificationCode == 0)
	res, _ := ebustl.CreateTeletextRasterFromTti(tti, codepage, TRUNCATEOVERSIZEDLINES, DEBUG, IGNOREINVALIDDIACRITICAL)
	fmt.Println(res.PrintToConsole())
	xmlstring, row_count, left_padding, char_count, err := getSubtitlePara(*res, fixedalignment, GetConfig())
	if err != nil {
		t.Fatalf("ERROR: Test1 err=%s", err.Error())
	}
	fmt.Println(xmlstring)
	fmt.Printf("row_count=%d, left_padding=%d, char_count=%d", row_count, left_padding, char_count)
}

func Test4(t *testing.T) {
	/*
		Matt Simpon's Character Code stuff

		0 1 2 3 4 5 6 7 8 9 á à â ç é è
		ê í î ñ ó ô ú û ! “ ( )


	*/
	ExtendedTextField := [...]byte{
		0x0D, 0x07, 0x0B, 0x0B, 0x30, 0x20, 0x31, 0x20, 0x32, 0x20, 0x33, 0x20, 0x34, 0x20, 0x35, 0x20,
		0x36, 0x20, 0x37, 0x20, 0x38, 0x20, 0x39, 0x20, 0xC2, 0x61, 0x20, 0xC1, 0x61, 0x20, 0xC3, 0x61,
		0x20, 0xCB, 0x63, 0x20, 0xC2, 0x65, 0x20, 0xC1, 0x65, 0x0A, 0x0A, 0x8A, 0x8A, 0x0D, 0x07, 0x0B,
		0x0B, 0xC3, 0x65, 0x20, 0xC2, 0x69, 0x20, 0xC3, 0x69, 0x20, 0xC4, 0x6E, 0x20, 0xC2, 0x6F, 0x20,
		0xC3, 0x6F, 0x20, 0xC2, 0x75, 0x20, 0xC3, 0x75, 0x20, 0x21, 0x20, 0xAA, 0x20, 0x28, 0x20, 0x29,
		0x0A, 0x0A, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F,
		0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F,
	}
	tti := ebustl.Tti{
		ExtendedTextField: ExtendedTextField[:],
		VerticalPosition:  22,
		JustificationCode: 2,
	}
	codepage := "00"
	fixedalignment := (tti.JustificationCode == 0)
	res, _ := ebustl.CreateTeletextRasterFromTti(tti, codepage, TRUNCATEOVERSIZEDLINES, DEBUG, IGNOREINVALIDDIACRITICAL)
	fmt.Println(res.PrintToConsole())
	xmlstring, row_count, left_padding, char_count, err := getSubtitlePara(*res, fixedalignment, GetConfig())
	if err != nil {
		t.Fatalf("ERROR: Test1 err=%s", err.Error())
	}
	fmt.Println(xmlstring)
	fmt.Printf("row_count=%d, left_padding=%d, char_count=%d", row_count, left_padding, char_count)
}

func Test5(t *testing.T) {

	/*
	   SHEEP BAA
	          Oh, for...

	*/
	ExtendedTextField := [...]byte{
		0x0D, 0x20, 0x20, 0x20, 0x20, 0x20, 0x0B, 0x0B, 0x53, 0x48, 0x45, 0x45, 0x50, 0x20, 0x42, 0x41,
		0x41, 0x0A, 0x0A, 0x8A, 0x8A, 0x0D, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20,
		0x20, 0x20, 0x0B, 0x0B, 0x4F, 0x68, 0x2C, 0x20, 0x66, 0x6F, 0x72, 0x2E, 0x2E, 0x2E, 0x0A, 0x0A,
		0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F,
		0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F,
		0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F,
		0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F,
	}
	tti := ebustl.Tti{
		ExtendedTextField: ExtendedTextField[:],
		VerticalPosition:  1,
		JustificationCode: 0,
	}

	codepage := "00"
	// fixedalignment := (tti.JustificationCode == 0)
	// res, _ := ebustl.CreateTeletextRasterFromTti(tti, codepage, true)
	// 	fmt.Println(res.PrintToConsole())
	ttmlsub, region, err := getSubtitle(tti, codepage, false, GetConfig())
	if err != nil {
		t.Fatalf("ERROR: Test1 err=%s", err.Error())
	}
	fmt.Println(ttmlsub.ToString())
	fmt.Println(region)
	//fmt.Printf("row_count=%d, left_padding=%d, char_count=%d", row_count, left_padding, char_count)
}

func Test6(t *testing.T) {

	/*
		OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO
		O                                  O
		O                                  O
		O                                  O
		O                                  O
		O                                  O
		O                                  O
		O                                  O
		O                                  O
		O                                  O
		OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO

	*/
	ExtendedTextField := [...]byte{
		0x0D, 0x07, 0x0B, 0x0B, 0x4F, 0x4F, 0x4F, 0x4F, 0x4F, 0x4F, 0x4F, 0x4F, 0x4F, 0x4F, 0x4F, 0x4F,
		0x4F, 0x4F, 0x4F, 0x4F, 0x4F, 0x4F, 0x4F, 0x4F, 0x4F, 0x4F, 0x4F, 0x4F, 0x4F, 0x4F, 0x4F, 0x4F,
		0x4F, 0x4F, 0x4F, 0x4F, 0x4F, 0x4F, 0x4F, 0x4F, 0x8A, 0x8A, 0x0D, 0x07, 0x0B, 0x0B, 0x4F, 0x20,
		0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20,
		0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20,
		0x20, 0x4F, 0x8A, 0x8A, 0x0D, 0x07, 0x0B, 0x0B, 0x4F, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20,
		0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20,

		0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x4F, 0x8A, 0x8A, 0x0D, 0x07,
		0x0B, 0x0B, 0x4F, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20,
		0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20,
		0x20, 0x20, 0x20, 0x20, 0x20, 0x4F, 0x8A, 0x8A, 0x0D, 0x07, 0x0B, 0x0B, 0x4F, 0x20, 0x20, 0x20,
		0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20,
		0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x4F,
		0x8A, 0x8A, 0x0D, 0x07, 0x0B, 0x0B, 0x4F, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20,

		0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20,
		0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x4F, 0x8A, 0x8A, 0x0D, 0x07, 0x0B, 0x0B,
		0x4F, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20,
		0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20,
		0x20, 0x20, 0x20, 0x4F, 0x8A, 0x8A, 0x0D, 0x07, 0x0B, 0x0B, 0x4F, 0x20, 0x20, 0x20, 0x20, 0x20,
		0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20,
		0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x4F, 0x8A, 0x8A,

		0x0D, 0x07, 0x0B, 0x0B, 0x4F, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20,
		0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20,
		0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x4F, 0x8A, 0x8A, 0x0D, 0x07, 0x0B, 0x0B, 0x4F, 0x20,
		0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20,
		0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20,
		0x20, 0x4F, 0x8A, 0x8A, 0x0D, 0x07, 0x0B, 0x0B, 0x4F, 0x4F, 0x4F, 0x4F, 0x4F, 0x4F, 0x4F, 0x4F,
		0x4F, 0x4F, 0x4F, 0x4F, 0x4F, 0x4F, 0x4F, 0x4F, 0x4F, 0x4F, 0x4F, 0x4F, 0x4F, 0x4F, 0x4F, 0x4F,

		0x4F, 0x4F, 0x4F, 0x4F, 0x4F, 0x4F, 0x4F, 0x4F, 0x4F, 0x4F, 0x4F, 0x4F, 0x8F, 0x8F, 0x8F, 0x8F,
		0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F,
		0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F,
		0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F,
		0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F,
		0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F,
		0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F,
	}
	tti := ebustl.Tti{
		ExtendedTextField: ExtendedTextField[:],
		VerticalPosition:  1,
		JustificationCode: 0,
	}

	codepage := "00"
	// fixedalignment := (tti.JustificationCode == 0)
	// res, _ := ebustl.CreateTeletextRasterFromTti(tti, codepage, true)
	// 	fmt.Println(res.PrintToConsole())
	ttmlsub, region, err := getSubtitle(tti, codepage, false, GetConfig())
	if err != nil {
		t.Fatalf("ERROR: Test1 err=%s", err.Error())
	}
	fmt.Println(ttmlsub.ToString())
	fmt.Println(region)
	//fmt.Printf("row_count=%d, left_padding=%d, char_count=%d", row_count, left_padding, char_count)
}

func Test7(t *testing.T) {

	/*
		It was all a question of technique.
		Technique?


	*/
	ExtendedTextField := [...]byte{
		0x0D, 0x03, 0x0B, 0x0B, 0x49, 0x74, 0x20, 0x77, 0x61, 0x73, 0x20, 0x61, 0x6C, 0x6C, 0x20, 0x61,
		0x20, 0x71, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6F, 0x6E, 0x20, 0x6F, 0x66, 0x20, 0x74, 0x65, 0x63,
		0x68, 0x6E, 0x69, 0x71, 0x75, 0x65, 0x2E, 0x0A, 0x8A, 0x8A, 0x0D, 0x0B, 0x0B, 0x54, 0x65, 0x63,
		0x68, 0x6E, 0x69, 0x71, 0x75, 0x65, 0x3F, 0x03, 0x0A, 0x0A, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F,
		0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F,
		0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F,
		0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F,
	}
	tti := ebustl.Tti{
		ExtendedTextField: ExtendedTextField[:],
		VerticalPosition:  20,
		JustificationCode: 2,
	}

	codepage := "00"
	// fixedalignment := (tti.JustificationCode == 0)
	// res, _ := ebustl.CreateTeletextRasterFromTti(tti, codepage, true)
	// 	fmt.Println(res.PrintToConsole())
	ttmlsub, region, err := getSubtitle(tti, codepage, false, GetConfig())
	if err != nil {
		t.Fatalf("ERROR: Test1 err=%s", err.Error())
	}
	fmt.Println(ttmlsub.ToString())
	fmt.Println(region)
	//fmt.Printf("row_count=%d, left_padding=%d, char_count=%d", row_count, left_padding, char_count)
}

func Test8(t *testing.T) {

	/*
		And Parliament will be baying for my
		blood if all I've got to show for it <<<< this is > 40 chars
		are two accomplice
	*/

	ExtendedTextField := [...]byte{
		0x0D, 0x0B, 0x0B, 0x03, 0x41, 0x6E, 0x64, 0x20, 0x50, 0x61, 0x72, 0x6C, 0x69, 0x61, 0x6D, 0x65,
		0x6E, 0x74, 0x20, 0x77, 0x69, 0x6C, 0x6C, 0x20, 0x62, 0x65, 0x20, 0x62, 0x61, 0x79, 0x69, 0x6E,
		0x67, 0x20, 0x66, 0x6F, 0x72, 0x20, 0x6D, 0x79, 0x0A, 0x0A, 0x8A, 0x8A, 0x0D, 0x0B, 0x0B, 0x03,
		0x03, 0x62, 0x6C, 0x6F, 0x6F, 0x64, 0x20, 0x69, 0x66, 0x20, 0x61, 0x6C, 0x6C, 0x20, 0x49, 0x27,
		0x76, 0x65, 0x20, 0x67, 0x6F, 0x74, 0x20, 0x74, 0x6F, 0x20, 0x73, 0x68, 0x6F, 0x77, 0x20, 0x66,
		0x6F, 0x72, 0x20, 0x69, 0x74, 0x0A, 0x0A, 0x8A, 0x8A, 0x0D, 0x0B, 0x0B, 0x03, 0x03, 0x61, 0x72,
		0x65, 0x20, 0x74, 0x77, 0x6F, 0x20, 0x61, 0x63, 0x63, 0x6F, 0x6D, 0x70, 0x6C, 0x69, 0x63, 0x65,
	}
	tti := ebustl.Tti{
		ExtendedTextField: ExtendedTextField[:],
		VerticalPosition:  20,
		JustificationCode: 2,
	}

	codepage := "00"
	// fixedalignment := (tti.JustificationCode == 0)
	// res, _ := ebustl.CreateTeletextRasterFromTti(tti, codepage, true)
	// 	fmt.Println(res.PrintToConsole())
	ttmlsub, region, err := getSubtitle(tti, codepage, false, GetConfig())
	if err != nil {
		t.Fatalf("ERROR: Test1 err=%s", err.Error())
	}
	fmt.Println(ttmlsub.ToString())
	fmt.Println(region)
	//fmt.Printf("row_count=%d, left_padding=%d, char_count=%d", row_count, left_padding, char_count)
}

func Test9(t *testing.T) {

	/*
		missing double height
	*/

	ExtendedTextField := [...]byte{
		0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F,
		0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F,
		0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F,
		0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F,
		0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F,
		0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F,
		0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F,
	}
	tti := ebustl.Tti{
		ExtendedTextField: ExtendedTextField[:],
		VerticalPosition:  20,
		JustificationCode: 2,
	}

	codepage := "00"
	// fixedalignment := (tti.JustificationCode == 0)
	// res, _ := ebustl.CreateTeletextRasterFromTti(tti, codepage, true)
	// 	fmt.Println(res.PrintToConsole())
	ttmlsub, region, err := getSubtitle(tti, codepage, false, GetConfig())
	if err != nil {
		t.Fatalf("ERROR: Test1 err=%s", err.Error())
	}
	fmt.Println(ttmlsub.ToString())
	fmt.Println(region)
	//fmt.Printf("row_count=%d, left_padding=%d, char_count=%d", row_count, left_padding, char_count)
}

func Test10(t *testing.T) {

	/*
		no start box (just a single box at the right)
	*/

	ExtendedTextField := [...]byte{
		0x20, 0x20, 0x20, 0x20, 0x0D, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20,
		0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20,
		0x20, 0x20, 0x20, 0x20, 0x20, 0x07, 0x1D, 0x07, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F,
		0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F,
		0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F,
		0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F,
		0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F,
	}
	tti := ebustl.Tti{
		ExtendedTextField: ExtendedTextField[:],
		VerticalPosition:  20,
		JustificationCode: 0,
	}

	codepage := "00"
	// fixedalignment := (tti.JustificationCode == 0)
	// res, _ := ebustl.CreateTeletextRasterFromTti(tti, codepage, true)
	// 	fmt.Println(res.PrintToConsole())
	ttmlsub, region, err := getSubtitle(tti, codepage, false, GetConfig())
	if err != nil {
		t.Fatalf("ERROR: Test1 err=%s", err.Error())
	}
	fmt.Println(ttmlsub.ToString())
	fmt.Println(region)
	//fmt.Printf("row_count=%d, left_padding=%d, char_count=%d", row_count, left_padding, char_count)
}

func Test11(t *testing.T) {

	/*
		really long tti

	*/

	ExtendedTextField := [...]byte{
		0x0D, 0x0B, 0x0B, 0x54, 0x48, 0x45, 0x20, 0x41, 0x44, 0x56, 0x45, 0x4E, 0x54, 0x55, 0x52, 0x45,
		0x53, 0x20, 0x4F, 0x46, 0x20, 0x53, 0x48, 0x45, 0x52, 0x4C, 0x4F, 0x43, 0x4B, 0x20, 0x48, 0x4F,
		0x4C, 0x4D, 0x45, 0x53, 0x20, 0x20, 0x23, 0x20, 0x43, 0x6F, 0x70, 0x70, 0x65, 0x72, 0x20, 0x42,
		0x65, 0x65, 0x63, 0x68, 0x65, 0x73, 0x20, 0x54, 0x58, 0x3A, 0x20, 0x32, 0x31, 0x2F, 0x30, 0x32,
		0x20, 0x20, 0x44, 0x55, 0x45, 0x3A, 0x20, 0x31, 0x38, 0x2F, 0x30, 0x32, 0x20, 0x30, 0x39, 0x2F,
		0x34, 0x39, 0x2F, 0x35, 0x35, 0x2F, 0x30, 0x30, 0x20, 0x43, 0x4F, 0x4E, 0x54, 0x49, 0x4E, 0x55,
		0x4F, 0x55, 0x53, 0x0A, 0x0A, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F,
	}
	tti := ebustl.Tti{
		ExtendedTextField: ExtendedTextField[:],
		VerticalPosition:  20,
		JustificationCode: 0,
	}

	codepage := "00"
	// fixedalignment := (tti.JustificationCode == 0)
	// res, _ := ebustl.CreateTeletextRasterFromTti(tti, codepage, true)
	// 	fmt.Println(res.PrintToConsole())
	cfg := GetConfig()
	ttmlsub, region, err := getSubtitle(tti, codepage, false, cfg)
	if err != nil {
		t.Fatalf("ERROR: Test11 err=%s", err.Error())
	}
	fmt.Println(ttmlsub.ToString())
	fmt.Println(region)
	//fmt.Printf("row_count=%d, left_padding=%d, char_count=%d", row_count, left_padding, char_count)
}

func Test12(t *testing.T) {

	/*
		1g3fgkq_ENG.stl
		really long tti

	*/

	ExtendedTextField := [...]byte{
		0x0D, 0x0B, 0x0B, 0x03, 0x48, 0x6F, 0x77, 0x20, 0x77, 0x6F, 0x75, 0x6C, 0x64, 0x20, 0x79, 0x6F,
		0x75, 0x20, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x20, 0x6D, 0x79, 0x20, 0x6C, 0x6F,
		0x6F, 0x6B, 0x73, 0x3F, 0x20, 0x07, 0x20, 0x0A, 0x0A, 0x8A, 0x8A, 0x0D, 0x0B, 0x0B, 0x07, 0x07,
		0x07, 0x07, 0x03, 0x07, 0x4D, 0x79, 0x20, 0x66, 0x72, 0x69, 0x65, 0x6E, 0x64, 0x20, 0x69, 0x73,
		0x20, 0x63, 0x68, 0x75, 0x62, 0x62, 0x79, 0x2E, 0x20, 0x03, 0x4A, 0x6F, 0x65, 0x79, 0x20, 0x2D,
		0x07, 0x59, 0x6F, 0x75, 0x20, 0x61, 0x72, 0x65, 0x21, 0x0A, 0x0A, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F,
		0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F,
	}
	tti := ebustl.Tti{
		ExtendedTextField: ExtendedTextField[:],
		VerticalPosition:  20,
		JustificationCode: 0,
	}

	codepage := "00"
	// fixedalignment := (tti.JustificationCode == 0)
	// res, _ := ebustl.CreateTeletextRasterFromTti(tti, codepage, true)
	// 	fmt.Println(res.PrintToConsole())
	cfg := GetConfig()
	ttmlsub, region, err := getSubtitle(tti, codepage, false, cfg)
	if err != nil {
		t.Fatalf("ERROR: Test12 err=%s", err.Error())
	}
	fmt.Println(ttmlsub.ToString())
	fmt.Println(region)
	//fmt.Printf("row_count=%d, left_padding=%d, char_count=%d", row_count, left_padding, char_count)
}

func Test13(t *testing.T) {

	/*
		1g3fgkq_ENG.stl
		really long tti

	*/

	ExtendedTextField := [...]byte{
		0x0D, 0x03, 0x0B, 0x0B, 0x42, 0x75, 0x74, 0x20, 0x49, 0x20, 0x73, 0x65, 0x63, 0x6F, 0x6E, 0x64,
		0x2D, 0x67, 0x75, 0x65, 0x73, 0x73, 0x20, 0x6D, 0x79, 0x73, 0x65, 0x6C, 0x66, 0xC1, 0x0A, 0x0A,
		0x8A, 0x8A, 0x0D, 0x03, 0x0B, 0x0B, 0x61, 0x6E, 0x64, 0x20, 0x74, 0x68, 0x61, 0x74, 0x27, 0x73,
		0x20, 0x6D, 0x79, 0x20, 0x62, 0x69, 0x67, 0x67, 0x65, 0x73, 0x74, 0x20, 0x77, 0x65, 0x61, 0x6B,
		0x6E, 0x65, 0x73, 0x73, 0x2E, 0x0A, 0x0A, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F,
		0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F,
		0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F, 0x8F,
	}
	tti := ebustl.Tti{
		ExtendedTextField: ExtendedTextField[:],
		VerticalPosition:  20,
		JustificationCode: 0,
	}

	codepage := "00"
	// fixedalignment := (tti.JustificationCode == 0)
	// res, _ := ebustl.CreateTeletextRasterFromTti(tti, codepage, true)
	// 	fmt.Println(res.PrintToConsole())
	cfg := GetConfig()
	cfg.IgnoreInvalidDiacritical = IGNOREINVALIDDIACRITICAL
	cfg.Debug = true
	ttmlsub, region, err := getSubtitle(tti, codepage, false, cfg)
	if err != nil {
		t.Fatalf("ERROR: Test13 err=%s", err.Error())
	}
	fmt.Println(ttmlsub.ToString())
	fmt.Println(region)
	//fmt.Printf("row_count=%d, left_padding=%d, char_count=%d", row_count, left_padding, char_count)
}
