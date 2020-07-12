# gomediaconverter
Go based Media converter:
Used to convert avi files to mpeg format using ffmpeg utility. Built-in Parallism for concurrent conversion

## How to build

```
go build ./convertAvi2Mpeg.go ./command.go
```

## How to Run

You need to only a mandatory parameter of path to the input folder where the files are recursively searched through and then
the conversion is run

```
./convertAvi2Mpeg -h

Usage of ./convertAvi2Mpeg:
  -dryrun
        Only list the files to be processed
  -inputdir string
        The input directory where avi files are stored. All files under this folder will be Recursively processed
  -nofiles uint
        no of files to process simulatenosuly (default 20)

```

## Disclaimer

Created as a personal fun project. Please do not use for commercial purposes.