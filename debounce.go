package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"time"
)

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
	delay := 500 * time.Millisecond
	debounce(os.Stdin, os.Stdout, os.Stderr, delay)
}
