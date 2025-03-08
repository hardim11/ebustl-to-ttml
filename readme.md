# EBU STL to TTML convertor

## What is this?
This is some code to convert EBU STL subtitle files used primarily in Broadcast into [TTML](https://en.wikipedia.org/wiki/Timed_Text_Markup_Language) XML files used also within Broadcast but also in OTT. 

EBU STL files are used to store the raw subtitle lines ("Cues") with a timecode of when to show them. They are generally "played out" alongside the video playback and are generally formatted for presentation on Teletext subtitle systems. This means that the authoring of the EBU STLs files will assume a display area and attributes, such as line and column limits, that the font is monospaced and that subtitles are actually shown using 2 Teletext lines ("Double Height").

## Why?
I was working on a project which required conversion of EBU STL files to TTML for OTT but we were experiencing some complications with the commercial software that we were using. We are still using the commercial software but this project helped me understand the standards better and learn what to look for. As I spent time on the project I thought I would share my awful code for others.

## References
The EBU STL format is documented in [EBU Tech 3264-1991](https://tech.ebu.ch/docs/tech/tech3264.pdf)

The EBU also released a very useful document on suggested approaches to converting EBU STLs to EBU Timed Text in STL [Mapping to EBU-TT (EBU Tech 3360)](https://tech.ebu.ch/docs/tech/tech3360.pdf)

Whilst this code is not outputting EBU Timed Text, it aims to output the [IMSC1](https://www.w3.org/TR/ttml-imsc1.0.1/) standard which, although a superset of EBU-TT, this code should maintain compatibility. [EBU comparison of EBU-TT and IMSC1](https://tech.ebu.ch/docs/events/IBC2015/EBU-TT-D_and_IMSC.pdf)

Within the code, I copied heavily from [Quentin Renard's go-astisub project](https://github.com/asticode/go-astisub) for the TTML code and for ideas. I am sure that that project, and others, would provide the functionality I needed, but the purpose of this project was to learn, not to use!

I hope I have credited all code and standards source but apologies for any I missed. Additionally, I captured a number of public DASH video streams and decoded the subtitles to analise alternative ways of formatting.

## Challenges
The greatest challenge I found was maintaining positioning. STLs offer 3 horizontal positions - left justified, centre justified and right justified - but it also offers unjustified ("unchanged presentation"). This latter mode is often used with space padding between lines to provide more complex positioning. The only method I have found to support this, without breaking up the individual subtitle "cues" into lots of regions, is to use the `xml:space="preserve"` attribute but I am not convinced of how widely this is supported between players.

Additionally, I really struggled with how to parse the control characters, the EBU STL format uses a collection of control characters to format the Cues (e.g. to change the colour) but also to provided accented characters. This, coupled with the "Double Height" of teletext subtitles did mean I used some assumptions, which I hope are valid.

## Code
There are 2 packages;

* ebustl - a very basic STL reader, it does not manipulate the subtitles, rather just reads the file to structures processing later.
* ttmlgenerate - this is where the interpretation of the Cues is performed and where all the styling is applied.

A very basic commanline app is included in ./cmd/stl-to-ttlml which can convert a single suppied file to TTML e.g. (when compiled)

    ./stl-to-ttml subs_test_with_audio_v1.stl output/subs_test_with_audio_v1.ttml

alternatively, you can run it without pre-compiling e.g.

    go run ./stl-to-ttml.go subs_test_with_audio_v1.stl output/subs_test_with_audio_v1.ttml


## Improvements, Errors
As always, please feel free to amend, update, correct!


## Issues
* only 25 fps
* ~~Accented chars~~
* ~~leading and trailing empty box missing~~
* make settings more configurable
* support alternative display approaches.
* create other variants of TTML.
* only a couple of code pages supported in Cues (and only 850 for the header "GSI" section)
