# simpleMediaConverter.
Go based Media converter:
Used to convert media files to specified output format. Built-in Parallism for concurrent conversion.
It uses ffmpeg utility for default which avi to mpeg4 format

## How to build

```
go build ./simpleMediaConverter.go ./avi2Mpeg4Conversion.go
```

## How to Run

You need to only a mandatory parameter of path to the input folder where the files are recursively searched through and then
the conversion is run

```
./simpleMediaConverter. -h

Usage of ./simpleMediaConverter:
  -convert string
        Convert input files to output format ... eg: avi2mpeg4 (default "avi2mpeg4")
  -dryrun
        Only list the files to be processed
  -inputdir string
        The input directory where avi files are stored. All files under this folder will be Recursively processed
  -nofiles uint
        no of files to process simulatenosuly (default 20)
  -parallel
        All files to be processed simulatenosuly
  -serial
        All files to be processed serially (default true)

```

## Disclaimer

Created as a personal fun project. Please do not use for commercial purposes.