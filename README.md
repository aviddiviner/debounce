# debounce
`debounce` is a simple command-line utility for handling noisy processes, written in Go.

To install it, just `go get github.com/aviddiviner/debounce`

## Usage

```
Tail lines from STDIN and output the last line to STDOUT, after receiving
nothing for a while.

Usage:
  debounce [options]

Example:
  tail -F * | debounce -d 100
  yes | head -50000 | debounce

Options:
  -d <ms>       Delay in milliseconds [default: 500].
  -h --help     Show this screen.
```
