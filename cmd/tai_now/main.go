// tai_now prints the current TAI time to standard output.
// It provides the time as a raw 64-bit integer and as a formatted 
// TAI64 timestamp.
package main

import (
	"fmt"
	"go-libtai/lib/tai"
)

// Timestamp formats a packed TAI64 value into a 
// human-readable hex string preceded by '@'.
func Timestamp(s [tai.TaiPack]byte) string {
	timestamp := make([]byte, 17)
	timestamp[0] = '@'
	for i := 0; i < tai.TaiPack; i++ {
		timestamp[i*2+1] = tai.HexDigit((s[i] >> 4) & 0xF)
		timestamp[i*2+2] = tai.HexDigit(s[i] & 0xF)
	}
	return string(timestamp)
}

func main() {
	taiTime := tai.Now()
	s := taiTime.Pack()
	fmt.Println("Current TAI time:", taiTime.Sec)
	fmt.Println("Current TAI timestamp:", Timestamp(s))
}
