// Package taia implements the TAIA64 time format.
// TAIA64 extends TAI64 with nanosecond and attosecond precision, 
// providing a total resolution of one quintillionth of a second.
package taia

import (
	"encoding/binary"
	"errors"
	"fmt"
	"go-libtai/lib/tai"
	"math"
	"strings"
	"time"
)

// TAIA represents a high-precision time value.
// It combines a 64-bit TAI second value with additional precision fields.
type TAIA struct {
	Sec  tai.TAI // TAI seconds
	Nano uint32  // Nanoseconds (0-999,999,999)
	Atto uint32  // Attoseconds (0-999,999,999)
}

// Time conversion constants.
const (
	NanoInSec  = 1_000_000_000
	AttoInNano = 1_000_000_000
)

// Standard error values for TAIA arithmetic.
var ErrOverflow = errors.New("taia: overflow")
var ErrUnderflow = errors.New("taia: underflow")

// TaiaPack is the number of bytes in a packed TAIA64 value.
const TaiaPack = 16

// TaiaFmtFrac is the number of characters used to represent the 
// fractional part in a string.
const TaiaFmtFrac = 19

// Now returns the current time in TAIA format.
// Nanoseconds are retrieved from the system clock; attoseconds are 
// set to zero as they are not supported by the Go standard library.
func Now() TAIA {
	now := time.Now()
	return TAIA{
		Sec:  tai.TAI{Sec: tai.TaiUnixOffset + uint64(now.Unix())},
		Nano: uint32(now.Nanosecond()),
		Atto: 0,
	}
}

// Frac returns the fractional part of the TAIA time as a float64.
func (t TAIA) Frac() float64 {
	return (float64(t.Atto)*float64(1e-09) + float64(t.Nano)) * float64(1e-09)
}

// FmtFrac formats the fractional part into an 18-character string 
// representing nanoseconds and attoseconds.
func (t *TAIA) FmtFrac() string {
	var builder strings.Builder
	fmt.Fprintf(&builder, "%09d%09d", t.Nano, t.Atto)
	return builder.String()
}

// Approx returns a float64 approximation of the TAIA time in seconds.
func (t TAIA) Approx() float64 {
	return float64(t.Sec.Sec) + t.Frac()
}

// Less returns true if time t is earlier than time u.
func (t TAIA) Less(u TAIA) bool {
	if t.Sec.Sec < u.Sec.Sec {
		return true
	}
	if t.Sec.Sec > u.Sec.Sec {
		return false
	}
	if t.Nano < u.Nano {
		return true
	}
	if t.Nano > u.Nano {
		return false
	}
	return t.Atto < u.Atto
}

// Half returns a TAIA value that is exactly half the duration of t.
// It correctly handles carries between the second, nanosecond, and 
// attosecond fields.
func (t TAIA) Half() TAIA {
	return TAIA{
		Sec:  tai.TAI{Sec: t.Sec.Sec / 2},
		Nano: uint32((uint64(t.Sec.Sec%2)*1_000_000_000 + uint64(t.Nano)) / 2),
		Atto: uint32((uint64(t.Nano%2)*1_000_000_000 + uint64(t.Atto)) / 2),
	}
}

// Uint sets the second field to s and resets the fractional precision.
func (t *TAIA) Uint(s uint) {
	t.Sec.Sec = uint64(s)
	t.Nano = 0
	t.Atto = 0
}

// Tai returns the TAI seconds component of the TAIA value.
func (t TAIA) Tai() tai.TAI {
	return tai.TAI{Sec: t.Sec.Sec}
}

// String returns a string representation of the TAIA value.
func (t TAIA) String() string {
	return fmt.Sprintf("%18d%.12f", t.Sec.Sec, t.Frac())
}

// Pack converts the TAIA time into a 16-byte big-endian slice.
// The format is: 8 bytes for seconds, 4 bytes for nanoseconds, 
// and 4 bytes for attoseconds.
func (t TAIA) Pack() [TaiaPack]byte {
	var s [TaiaPack]byte
	binary.BigEndian.PutUint64(s[:8], t.Sec.Sec)
	binary.BigEndian.PutUint32(s[8:], t.Nano)
	binary.BigEndian.PutUint32(s[12:], t.Atto)
	return s
}

// Unpack converts a 16-byte big-endian slice into a TAIA time.
func Unpack(s [TaiaPack]byte) TAIA {
	return TAIA{
		Sec:  tai.TAI{Sec: binary.BigEndian.Uint64(s[:8])},
		Nano: binary.BigEndian.Uint32(s[8:]),
		Atto: binary.BigEndian.Uint32(s[12:]),
	}
}

// Add returns the sum of two TAIA times.
// It handles overflows between attoseconds, nanoseconds, and seconds.
func (u TAIA) Add(v TAIA) (TAIA, error) {
	newSec := uint64(0)
	newNano := u.Nano + v.Nano
	newAtto := u.Atto + v.Atto

	if newAtto >= 1_000_000_000 {
		newAtto -= 1_000_000_000
		newNano++
	}

	if newNano >= 1_000_000_000 {
		newNano -= 1_000_000_000
		newSec++
	}

	if u.Sec.Sec > math.MaxUint64-v.Sec.Sec {
		return TAIA{}, ErrOverflow
	}
	resSec := u.Sec.Sec + v.Sec.Sec
	if resSec > math.MaxUint64-newSec {
		return TAIA{}, ErrOverflow
	}

	return TAIA{
		Sec:  tai.TAI{Sec: resSec + newSec},
		Nano: newNano,
		Atto: newAtto,
	}, nil
}

// Sub returns the difference between two TAIA times.
// It performs a multi-field subtraction with borrows. 
// It returns ErrUnderflow if the result would be negative.
func (u TAIA) Sub(v TAIA) (TAIA, error) {
	newSec := u.Sec.Sec
	newNano := u.Nano
	newAtto := u.Atto

	if newAtto < v.Atto {
		if newNano == 0 {
			if newSec == 0 {
				return TAIA{}, ErrUnderflow
			}
			newSec--
			newNano = 999_999_999
		} else {
			newNano--
		}
		newAtto += 1_000_000_000
	}
	newAtto -= v.Atto

	if newNano < v.Nano {
		if newSec == 0 {
			return TAIA{}, ErrUnderflow
		}
		newSec--
		newNano += 1_000_000_000
	}
	newNano -= v.Nano

	if newSec < v.Sec.Sec {
		return TAIA{}, ErrUnderflow
	}
	newSec -= v.Sec.Sec

	return TAIA{
		Sec:  tai.TAI{Sec: newSec},
		Nano: newNano,
		Atto: newAtto,
	}, nil
}
