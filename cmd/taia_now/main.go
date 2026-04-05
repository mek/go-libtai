// taia_now prints the current TAIA time to standard output.
// It provides high-precision time including nanoseconds.
package main

import (
	"fmt"
	"go-libtai/lib/tai"
	"go-libtai/lib/taia"
)

// Timestamp formats the first 12 bytes of a packed TAIA64 value 
// into a hex string preceded by '@'. This follows the standard 
// format for TAI64N timestamps.
func Timestamp(s [taia.TaiaPack]byte) string {
	timestamp := make([]byte, 25)
	timestamp[0] = '@'
	for i := 0; i < 12; i++ {
		timestamp[i*2+1] = tai.HexDigit((s[i] >> 4) & 0xF)
		timestamp[i*2+2] = tai.HexDigit(s[i] & 0xF)
	}
	return string(timestamp)
}

func main() {
	t := taia.Now()
	s := t.Pack()
	fmt.Println("Current TAIA time:", t.String())
	fmt.Println("Current TAIA timestamp:", Timestamp(s))
}
