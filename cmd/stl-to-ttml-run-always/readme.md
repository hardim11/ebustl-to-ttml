# Folder Scanner 
This version of the app will only provide the folder scan functionality. It should run and periodically scan the nominated directlty for *.stl files, for any file found, it should convert them to TTML and move the file to a processed folder.

# Configuration
The application requires a JSON file called "stl-to-ttml-run-always.json" located in the working directory (probably the same folder as the executable). Now command line arguments are required.

## Example Configuration file
```
{
	"SourceFolder": "\\\\pc7\\video_share\\ITV\\subtitles\\in",
	"TtmlOutputFolder": "\\\\pc7\\video_share\\ITV\\subtitles\\output",
	"ProcessedFolder": "\\\\pc7\\video_share\\ITV\\subtitles\\processed",
	"FailedFolder": "\\\\pc7\\video_share\\ITV\\subtitles\\failed",
	"ScanIntervalSeconds": 10,
	"StopOnError": false,
	"Debug": false
}
```
Hopefully the fields are self-explanatory. Normally "StopOnError" and "Debug" would be set to false. Beware that windows paths will include the "\\" charcter which will need to be escaped with another "\\" e.g. "\\\\".

## Compiling the executable
With the Repository downloaded to your computer and with Golang installed (this was written against Go 1.22), in a console window, change directory to the cmd\stl-to-ttml-run-always folder and execute following command

    go build .

On windows, an executable called "stl-to-ttml-run-always.exe" should be created (on Linux / Mac it will just be called "stl-to-ttml-run-always")

## Running
No command line arguments are required.

With the configuration file created, open a command prompt in the folder with the execuatable and just run it

	cd \myfolder\
	stl-to-ttml-run-always.exe

The application should start polling and log to the console output


