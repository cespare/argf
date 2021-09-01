// Package argf provides a simple way of reading line-by-line from either files
// given as command-line arguments or, if none were given, from stdin.
//
// The interface resembles bufio.Scanner.
//
// This package provides a convenient way of handling input for command-line
// utilities. For instance, here is a simple implementation of the Unix utility
// 'cat':
//
//   for argf.Scan() {
//     fmt.Println(argf.String())
//   }
//   if err := argf.Error(); err != nil {
//     fmt.Println(err)
//     os.Exit(1)
//   }
//
// If flags are required, you can call Init(flag.Args()) after flag parsing but
// before any other argf calls to initialize argf with the non-flag arguments
// given in the command-line (presumably filenames).
//
// Multiple goroutines should not call any of the functions in argf
// concurrently.
package argf

import (
	"bufio"
	"bytes"
	"io"
	"os"
)

var (
	initialized bool
	readStdin   bool
	reader      *bufio.Reader
	fileArgs    []string
	file        *os.File
	line        []byte
	curError    error
)

// Init initializes argf's state using some filename arguments. If args is
// empty, argf uses stdin instead of files. Without calling Init(), argf
// initializes itself the first time Scan is called, using os.Args[1:] (ignoring
// the program name).
func Init(args []string) {
	initialized = true
	if len(args) == 0 {
		readStdin = true
		reader = bufio.NewReader(os.Stdin)
	}
	fileArgs = args
}

// Scan reads the next line from either os.Stdin or the current file in os.Args,
// as described in the package documentation. If the current file has been
// exhausted, Scan attempts to open the next file in os.Args, if there is one.
// If there are no more lines to be read from os.Stdin or any files, or if Scan
// encounters an error, false is returned. Otherwise, true is returned and the
// line is available to be accessed by String or Bytes.
func Scan() bool {
	if !initialized {
		args := os.Args
		if len(args) >= 1 {
			args = args[1:] // Shift off the program name
		}
		Init(args)
	}
	if reader == nil {
		if readStdin {
			return false
		}
		if len(fileArgs) == 0 {
			return false
		}
		var err error
		file, err = os.Open(fileArgs[0])
		if err != nil {
			curError = err
			return false
		}
		fileArgs = fileArgs[1:]
		reader = bufio.NewReader(file)
	}
	var err error
	line, err = reader.ReadBytes('\n')
	if err != nil {
		if err != io.EOF {
			curError = err
			return false
		}
		if len(line) == 0 {
			if file != nil {
				file.Close()
			}
			reader = nil
			return Scan()
		}
	}
	line = bytes.TrimRight(line, "\r\n")
	return true
}

// String returns the current line as a string without the trailing newline. It
// panics unless preceeded by a call to Scan that returned true. String may be
// called multiple times consecutively but returns the same line each time.
func String() string {
	if !initialized {
		panic("argf: call to String before Scan.")
	}
	if curError != nil {
		panic("argf: call to String after false Scan()")
	}
	return string(line)
}

// Bytes returns the current line as a []byte without the trailing newline. It
// panics unless preceeded by a call to Scan that returned true. Bytes may be
// called multiple times consecutively but returns the same line each time.
func Bytes() []byte {
	if !initialized {
		panic("argf: call to Bytes before Scan.")
	}
	if curError != nil {
		panic("argf: call to Bytes after false Scan()")
	}
	return line
}

// Error returns the error that caused Scan to return false, unless it was an
// io.EOF, in which case Error returns nil.
func Error() error {
	if !initialized {
		panic("argf: call to Error before Scan.")
	}
	return curError
}
