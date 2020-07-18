# simpleMediaConverter.

Mass Media converter used to convert media files to specified output format. Built-in Parallism for concurrent conversion and 
support for slack integration for updates

Default: It uses ffmpeg utility by default which converts avi to mpeg4 format. But this can be scaled to by writing any
conversion function : avi to mkv etc

## Pre-requisities
 - Make sure you have go installed
 - The conversion utility like ffmpeg for avi to mpeg conversion

## How to build & Install

### Build
```
make build
```

### Install & Uninstall

It installs the binary in `GOPATH/bin` pointed by go env

```
make install
```

```
make uninstall
```

## How to Run

You need to only a mandatory parameter of path to the input folder where the files are recursively searched through and then
the conversion is run

```
simpleMediaConverter. -h

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

### Dry Run
```
simplemediaconverter -inputdir /PATH/TO/INPUTFOLDER -dryrun
```

### Actual Run with slack notificaitions
```
simplemediaconverter -inputdir /PATH/TO/INPUTFOLDER -nofiles=10 -notify=slack
```

### Actual Run for avi2mpeg(default) in parallel
```
simplemediaconverter -inputdir /PATH/TO/INPUTFOLDER -nofiles=10 -convert=avi2mpeg -parallel -notify=slack
```

You might need to run with sudo user priviliges if you dont have permissions on the mount directory in some cases attaching an external hard disk

## Troubleshooting

If it builds & installs fine and you still cant access the binary from your directory, then it could be because your GOPATH is not in PATH. Make sure to include the below in your `.bashrc`, `.bash_profile` or `.zshrc`

```
export PATH=$PATH:$(go env GOPATH)/bin
```

Fails from notifiers are soft fails and wont break conversion as such

## Suggestions
Suggestions, PRs & Issues are welcome.


## Disclaimer

Created as a personal fun project. Please do not use for commercial purposes.