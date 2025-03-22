package taia

import (
	"encoding/binary"
	"errors"
	"fmt"
	"go-libtai/lib/tai"
	"math"
	"time"
)

// TAIA is a high-precision time value.
type TAIA struct {
	Sec  tai.TAI // seconds in TAI format
	Nano uint32  // nanoseconds (0...999999999)
	Atto uint32  // attoseconds (0...999999999)
}

var ErrOverflow = errors.New("taia: overflow")
var ErrorUnderflow = errors.New("taia: underflow")
var TAIA_PACK = 16
var TAIA_FMTFRAC = 19

// Now returns the current TAIA time.
func Now() TAIA {
	now := time.Now()
	return TAIA{
		Sec:  tai.TAI{X: tai.TaiUnixOffset + uint64(now.Unix())},
		Nano: uint32(now.Nanosecond()),
		Atto: 0, // Go does not support attoseconds
	}
}

// fraction returns the fractional part of the TAIA time.
func (t TAIA) Frac() float64 {
	return (float64(t.Atto)*float64(1e-09) + float64(t.Nano)) * float64(1e-09)
}

// taiaFmtFrac formats the fractional part of a TAIA timestamp into an 18-character string
func (t *TAIA) FmtFrac() string {
	s := make([]byte, 18)

	x := t.Atto
	s[17] = '0' + byte(x%10)
	x /= 10
	s[16] = '0' + byte(x%10)
	x /= 10
	s[15] = '0' + byte(x%10)
	x /= 10
	s[14] = '0' + byte(x%10)
	x /= 10
	s[13] = '0' + byte(x%10)
	x /= 10
	s[12] = '0' + byte(x%10)
	x /= 10
	s[11] = '0' + byte(x%10)
	x /= 10
	s[10] = '0' + byte(x%10)
	x /= 10
	s[9] = '0' + byte(x%10)

	x = t.Nano
	s[8] = '0' + byte(x%10)
	x /= 10
	s[7] = '0' + byte(x%10)
	x /= 10
	s[6] = '0' + byte(x%10)
	x /= 10
	s[5] = '0' + byte(x%10)
	x /= 10
	s[4] = '0' + byte(x%10)
	x /= 10
	s[3] = '0' + byte(x%10)
	x /= 10
	s[2] = '0' + byte(x%10)
	x /= 10
	s[1] = '0' + byte(x%10)
	x /= 10
	s[0] = '0' + byte(x%10)

	return string(s)
}

// Approx a float64 approximation of the TAIA time.
func (t TAIA) Approx() float64 {
	return float64(t.Sec.X) + t.Frac()
}

// Less returns true if t < u.
func (t TAIA) Less(u TAIA) bool {
	if t.Sec.X < u.Sec.X {
		return true
	}
	if t.Sec.X > u.Sec.X {
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

// Half returns the TAIA time t / 2.
func (t TAIA) Half() TAIA {
	return TAIA{
		Sec:  tai.TAI{X: t.Sec.X / 2},
		Nano: uint32((uint64(t.Sec.X%2)*1_000_000_000 + uint64(t.Nano)) / 2),
		Atto: uint32((uint64(t.Nano%2)*1_000_000_000 + uint64(t.Atto)) / 2),
	}
}

func (t *TAIA) Uint(s uint) {
	t.Sec.X = uint64(s)
	t.Nano = 0
	t.Atto = 0
}

func (t TAIA) Tai() tai.TAI {
	return tai.TAI{X: t.Sec.X}
}

func (t TAIA) String() string {
	return fmt.Sprintf("%18d%.12f", t.Sec.X, t.Frac())
}

func (t TAIA) Pack() [16]byte {
	var s [16]byte
	binary.BigEndian.PutUint64(s[:8], t.Sec.X)
	binary.BigEndian.PutUint32(s[8:], t.Nano)
	binary.BigEndian.PutUint32(s[12:], t.Atto)
	return s
}

func Unpack(s [16]byte) TAIA {
	return TAIA{
		Sec:  tai.TAI{X: binary.BigEndian.Uint64(s[:8])},
		Nano: binary.BigEndian.Uint32(s[8:]),
		Atto: binary.BigEndian.Uint32(s[12:]),
	}
}

// Add returns the TAIA time u + v.
func (u TAIA) Add(v TAIA) (TAIA, error) {
	// Setup the new values

	newSec := uint64(0)
	newNano := u.Nano + v.Nano
	newAtto := u.Atto + v.Atto

	// check for overflows in nano and atto
	if newAtto >= 1_000_000_000 {
		newAtto -= 1_000_000_000
		newNano++
	}

	if newNano >= 1_000_000_000 {
		newNano -= 1_000_000_000
		newSec++
	}

	if u.Sec.X > math.MaxUint64-v.Sec.X+newSec {
		return TAIA{}, ErrOverflow
	}

	return TAIA{
		Sec:  tai.TAI{X: u.Sec.X + v.Sec.X + newSec},
		Nano: newNano,
		Atto: newAtto,
	}, nil

}

// Sub returns the TAIA time u - v.
func (u TAIA) Sub(v TAIA) (TAIA, error) {

	newSec := u.Sec.X
	newNano := u.Nano
	newAtto := u.Atto

	// check for attosecond underflow
	if newAtto < v.Atto { // we need to barrow a nanosecond
		if newNano == 0 { // do we need to borrow a second?
			if newSec == 0 { // out of time
				return TAIA{}, ErrorUnderflow
			}
			// 0 nanoseconds, borow 1 second
			// Set nanoseconds to max
			newSec--
			newNano = 999_999_999
		} else {
			// borrow 1 nanosecond
			newNano--
		}
		newAtto += 1_000_000_000
	}
	newAtto -= v.Atto

	// check for nanosecond underflow
	if newNano < v.Nano { // we need to borrow a second
		if newSec == 0 { // out of time
			return TAIA{}, ErrorUnderflow
		}
		// borrow 1 second and add to nanoseconds
		newSec--
		newNano += 1_000_000_000
	}
	newNano -= v.Nano

	// check for second underflow, i.e: negative time
	if newSec < v.Sec.X {
		return TAIA{}, ErrOverflow
	}
	newSec -= v.Sec.X

	return TAIA{
		Sec:  tai.TAI{X: newSec},
		Nano: newNano,
		Atto: newAtto,
	}, nil

}
