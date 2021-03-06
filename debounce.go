package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/aviddiviner/docopt-go"
)

const usage = `Tail lines from STDIN and output the last line to STDOUT, after receiving
nothing for a while.

Usage:
  debounce [options]

Example:
  tail -F * | debounce -d 100
  yes | head -50000 | debounce

Options:
  -d <ms>       Delay in milliseconds [default: 500].
  -h --help     Show this screen.`

type Options struct {
	Delay int `docopt:"-d"`
}

func writeLine(b []byte, out io.Writer) (err error) {
	if b == nil {
		return
	}
	_, err = out.Write(b)
	if err == nil {
		_, err = out.Write([]byte("\n"))
	}
	return
}

type Quit struct{}

func debounce(src io.Reader, dst io.Writer, errOut io.Writer, delay time.Duration) {
	out := make(chan []byte)
	quit := make(chan Quit)

	// Read from Stdin endlessly and send each line on the `out` channel.
	go func() {
		scanner := bufio.NewScanner(src)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			out <- scanner.Bytes()
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(errOut, "error reading from stdin:", err)
		}
		quit <- Quit{}
	}()

	var b []byte // Last line read.

	// Write the last line to Stdout.
	doWrite := func() {
		err := writeLine(b, dst)
		if err != nil {
			fmt.Fprintln(errOut, "error writing to stdout:", err)
			os.Exit(2)
		}
		b = nil
	}

	// Wait for new lines to arrive, write them after N seconds.
	for {
		select {
		case <-quit:
			doWrite()
			os.Exit(0)
		case b = <-out:
		case <-time.After(delay):
			doWrite()
		}
	}
}

func main() {
	checkErr := func(err error) {
		if err != nil {
			fmt.Fprintln(os.Stderr, "error parsing command line args:", err)
			os.Exit(1)
		}
	}
	var opts Options
	args, err := docopt.ParseDoc(usage)
	checkErr(err)
	checkErr(args.Bind(&opts))

	delay := time.Duration(opts.Delay) * time.Millisecond
	debounce(os.Stdin, os.Stdout, os.Stderr, delay)
}
