package logging

import (
	"io"
	"io/ioutil"
	"log"
)

// LogFlags = log.Ldate | log.Ltime | log.Lshortfile
const LogFlags = log.Ldate | log.Ltime | log.Lshortfile

// debug makes a logger that prints to w if verbose is true, prefixed with the LogFlags and error, and discards otherwise
func Debug(w io.Writer, verbose bool) *log.Logger {
	if verbose {
		return log.New(w, "DEBUG ", log.Ldate|log.Ltime|log.Lshortfile)
	}
	return log.New(ioutil.Discard, "", 0)
}

// debug makes a logger that prints to w if verbose is true, prefixed with the LogFlags and error
func Error(w io.Writer) *log.Logger {
	return log.New(w, "ERROR ", log.Ldate|log.Ltime|log.Lshortfile)
}
