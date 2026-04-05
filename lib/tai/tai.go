// Package tai implements the TAI64 time format.
// TAI64 represents a high-precision, linear time scale that avoids 
// the complexities of leap seconds by counting actual seconds passed.
package tai

import (
	"encoding/binary"
	"errors"
	"math"
	"time"
)

// Standard error values for TAI arithmetic.
var ErrOverflow = errors.New("tai: overflow")
var ErrUnderflow = errors.New("tai: underflow")

// TaiPack is the number of bytes in a packed TAI64 value.
const TaiPack = 8

// TAI represents a Temps Atomique International (International Atomic Time) value.
// It is stored as a 64-bit unsigned integer representing the number of seconds 
// passed since the TAI epoch.
type TAI struct {
	Sec uint64
}

// TaiUnixOffset is the offset used to convert between TAI and Unix time.
// It represents 1970-01-01 00:00:00 TAI, which is 2^62 seconds after the 
// absolute TAI start. The value includes the 10-second difference between 
// TAI and UTC at the start of 1970.
const TaiUnixOffset uint64 = 4611686018427387914

// Now returns the current TAI time based on the system clock.
// It assumes the system clock is synchronized to UTC.
func Now() TAI {
	return TAI{Sec: TaiUnixOffset + uint64(time.Now().Unix())}
}

// Approx returns a float64 approximation of the TAI time.
// This is useful for calculations where some loss of precision is acceptable.
func (t TAI) Approx() float64 {
	return float64(t.Sec)
}

// Less returns true if time t is earlier than time u.
func (t TAI) Less(u TAI) bool {
	return t.Sec < u.Sec
}

// Pack converts the TAI time into an 8-byte big-endian slice.
// This format is compatible with the TAI64 standard.
func (t TAI) Pack() [TaiPack]byte {
	var s [TaiPack]byte
	binary.BigEndian.PutUint64(s[:], t.Sec)
	return s
}

// Unpack converts an 8-byte big-endian slice into a TAI time.
func Unpack(s [TaiPack]byte) TAI {
	return TAI{Sec: binary.BigEndian.Uint64(s[:])}
}

// Add returns the sum of two TAI times. 
// It returns ErrOverflow if the result exceeds the 64-bit range.
func (u TAI) Add(v TAI) (TAI, error) {
	if u.Sec > math.MaxUint64-v.Sec {
		return TAI{}, ErrOverflow
	}
	return TAI{Sec: u.Sec + v.Sec}, nil
}

// Sub returns the difference between two TAI times.
// It returns ErrUnderflow if the result would be negative.
func (u TAI) Sub(v TAI) (TAI, error) {
	if u.Sec < v.Sec {
		return TAI{}, ErrUnderflow
	}
	return TAI{Sec: u.Sec - v.Sec}, nil
}

// HexDigit returns the hexadecimal character for a value n in the range [0, 15].
// If n is outside this range, it returns '0'.
func HexDigit(n byte) byte {
	if n < 10 {
		return '0' + n
	}
	if n < 16 {
		return 'a' + (n - 10)
	}
	return '0'
}
