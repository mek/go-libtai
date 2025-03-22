package main

import (
	"fmt"
	"go-libtai/lib/tai"
)

func hexDigit(n byte) byte {
	if n < 10 {
		return '0' + n
	}
	return 'a' + (n - 10)
}

func Timestamp(s [8]byte) string {
	timestamp := make([]byte, 25)
	timestamp[0] = '@'
	for i := 0; i < 8; i++ {
		timestamp[i*2+1] = hexDigit((s[i] >> 4) & 0xF)
		timestamp[i*2+2] = hexDigit(s[i] & 0xF)
	}
	return string(timestamp)
}

func main() {
	var s [8]byte

	taiTime := tai.Now()
	s = taiTime.Pack()
	fmt.Println("Current TAI time:", taiTime.X)
	fmt.Println("Current TAI timestamp:", Timestamp(s))
}
