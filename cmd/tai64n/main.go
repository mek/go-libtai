// tai64n reads lines from standard input and prepends each line 
// with a TAI64N timestamp. It is similar to DJB's tai64n tool 
// and is useful for high-precision log timestamping.
package main

import (
	"bufio"
	"fmt"
	"go-libtai/lib/tai"
	"go-libtai/lib/taia"
	"io"
	"os"
)

// timestamp returns a formatted TAI64N timestamp for the current time.
func timestamp() string {
	s := make([]byte, 25)
	now := taia.Now()
	nowpack := now.Pack()

	s[0] = '@'
	for i := 0; i < 12; i++ {
		s[i*2+1] = tai.HexDigit((nowpack[i] >> 4) & 0xF)
		s[i*2+2] = tai.HexDigit(nowpack[i] & 0xF)
	}
	return string(s)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	for {
		line, err := in.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				if len(line) > 0 {
					_, _ = out.WriteString(timestamp() + " ")
					_, _ = out.WriteString(line)
					out.Flush()
				}
				break
			} else {
				fmt.Fprintf(os.Stderr, "error reading input: %v\n", err)
				os.Exit(111)
			}
		}

		ts := timestamp()
		_, err = out.WriteString(ts + " ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "error writing output: %v\n", err)
			os.Exit(111)
		}
		_, err = out.WriteString(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error writing output: %v\n", err)
			os.Exit(111)
		}
	}
}
