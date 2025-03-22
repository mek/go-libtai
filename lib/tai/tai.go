package tai

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"time"
)

var ErrOverflow = errors.New("tai: overflow")
var ErrUnderflow = errors.New("tai: underflow")
var TAI_PACK = 8

type TAI struct {
	X uint64
}

// TaiUnixOffset is the number of seconds between TAI and Unix time.
const TaiUnixOffset uint64 = 4611686018427387914

// Now returns the current TAI time.
func Now() TAI {
	return TAI{X: TaiUnixOffset + uint64(time.Now().Unix())}
}

// Approx returns a float64 approximation of the TAI time.
func (t TAI) Approx() float64 {
	return float64(t.X)
}

// Is time t less than time u?
func (t TAI) Less(u TAI) bool {
	return t.X < u.X
}

// Pack converts the TAI time to a 8 byte slice.
func (t TAI) Pack() [8]byte {
	var s [8]byte
	binary.BigEndian.PutUint64(s[:], t.X)
	return s
}

// UnPack converts a 8 byte slice to a TAI time.
func UnPack(s [8]byte) TAI {
	return TAI{X: binary.BigEndian.Uint64(s[:])}
}

// Add returns the sum between two TAI times.
func (u TAI) Add(v TAI) (TAI, error) {
	if u.X > math.MaxUint64-v.X {
		return TAI{}, ErrOverflow
	}
	return TAI{X: u.X + v.X}, nil
}

// Sub returns the difference between two TAI times.
func (u TAI) Sub(v TAI) (TAI, error) {
	if u.X < v.X {
		return TAI{}, ErrUnderflow
	}
	return TAI{X: u.X - v.X}, nil
}

func HexDigit(n byte) byte {
	switch {
	case n < 10:
		return n + '0'
	case n < 16:
		return n - 10 + 'a'
	}
	fmt.Printf("HexDigit: invalid digit %d\n", n)
	panic("HexDigit: invalid digit")
}
