package main

import (
	"fmt"
	"go-libtai/lib/taia"
)

func hexDigit(n byte) byte {
	if n < 10 {
		return '0' + n
	}
	return 'a' + (n - 10)
}

func Timestamp(s [16]byte) string {
	timestamp := make([]byte, 25)
	timestamp[0] = '@'
	for i := 0; i < 12; i++ {
		timestamp[i*2+1] = hexDigit((s[i] >> 4) & 0xF)
		timestamp[i*2+2] = hexDigit(s[i] & 0xF)
	}
	return string(timestamp)
}

func main() {
	var s [16]byte

	t := taia.Now()
	s = t.Pack()
	fmt.Println("Current TAIA time:", t.String())
	fmt.Println("Current TAIA timestamp:", Timestamp(s))
}
