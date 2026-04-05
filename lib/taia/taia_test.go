package taia

import (
	"go-libtai/lib/tai"
	"testing"
	"time"
)

// Test Now()
func TestTAIA_Now(t *testing.T) {
	// Get current TAIA time
	taiNow := Now()

	// Get expected TAI seconds based on Unix time
	unixNow := uint64(time.Now().Unix())

	// Expected TAI time should be within a reasonable range (e.g., ±2 seconds)
	expectedTaiSec := tai.TaiUnixOffset + unixNow

	// Check seconds field
	if taiNow.Sec.Sec < expectedTaiSec-2 || taiNow.Sec.Sec > expectedTaiSec+2 {
		t.Errorf("Now() returned unexpected TAI seconds: got %d, expected around %d", taiNow.Sec.Sec, expectedTaiSec)
	}

	// Check that nanoseconds are within valid range (0 - 999_999_999)
	if taiNow.Nano >= 1_000_000_000 {
		t.Errorf("Now() returned invalid nanoseconds: %d", taiNow.Nano)
	}

	// Check that attoseconds are initialized to zero (since Go does not provide attosecond precision)
	if taiNow.Atto != 0 {
		t.Errorf("Now() should initialize attoseconds to 0, but got %d", taiNow.Atto)
	}
}

// Test TAIA Addition
func TestTAIA_Add(t *testing.T) {
	tests := []struct {
		t1, t2     TAIA
		expected   TAIA
		shouldFail bool
	}{
		// Simple addition (no overflow)
		{
			TAIA{Sec: tai.TAI{Sec: 100}, Nano: 500_000_000, Atto: 100_000_000},
			TAIA{Sec: tai.TAI{Sec: 50}, Nano: 200_000_000, Atto: 300_000_000},
			TAIA{Sec: tai.TAI{Sec: 150}, Nano: 700_000_000, Atto: 400_000_000},
			false,
		},
		// Attosecond overflow into Nanoseconds
		{
			TAIA{Sec: tai.TAI{Sec: 100}, Nano: 500_000_000, Atto: 900_000_000},
			TAIA{Sec: tai.TAI{Sec: 0}, Nano: 100_000_000, Atto: 200_000_000},
			TAIA{Sec: tai.TAI{Sec: 100}, Nano: 600_000_001, Atto: 100_000_000}, // Carry to Nano
			false,
		},
		// Nanosecond overflow into Seconds
		{
			TAIA{Sec: tai.TAI{Sec: 100}, Nano: 900_000_000, Atto: 0},
			TAIA{Sec: tai.TAI{Sec: 1}, Nano: 200_000_000, Atto: 0},
			TAIA{Sec: tai.TAI{Sec: 102}, Nano: 100_000_000, Atto: 0}, // Carry to Sec
			false,
		},
		// Overflow Case
		{
			TAIA{Sec: tai.TAI{Sec: ^uint64(0) - 1}, Nano: 900_000_000, Atto: 0}, // Max uint64 - 1
			TAIA{Sec: tai.TAI{Sec: 10}, Nano: 200_000_000, Atto: 0},
			TAIA{}, // Should fail
			true,
		},
	}

	for _, tc := range tests {
		result, err := tc.t1.Add(tc.t2)
		if (err != nil) != tc.shouldFail {
			t.Errorf("t1 Add(%v, %v, %v)", tc.t1.Sec.Sec, tc.t1.Nano, tc.t1.Atto)
			t.Errorf("t2 Add(%v, %v, %v)", tc.t2.Sec.Sec, tc.t2.Nano, tc.t2.Atto)
			t.Errorf("ex Add(%v, %v, %v)", tc.expected.Sec.Sec, tc.expected.Nano, tc.expected.Atto)
			t.Errorf("re Add(%v, %v, %v)", result.Sec.Sec, result.Nano, result.Atto)
			// t.Errorf("Add(%v, %v): unexpected error state: %v", tc.t1, tc.t2, err)
		} else if !tc.shouldFail && (result != tc.expected) {
			t.Errorf("t1 Add(%v, %v, %v)", tc.t1.Sec.Sec, tc.t1.Nano, tc.t1.Atto)
			t.Errorf("t2 Add(%v, %v, %v)", tc.t2.Sec.Sec, tc.t2.Nano, tc.t2.Atto)
			t.Errorf("ex Add(%v, %v, %v)", tc.expected.Sec.Sec, tc.expected.Nano, tc.expected.Atto)
			t.Errorf("re Add(%v, %v, %v)", result.Sec.Sec, result.Nano, result.Atto)
			// t.Errorf("Add(%v, %v): expected %v, got %v", tc.t1, tc.t2, tc.expected, result)
		}
	}
}

// Test TAIA Subtraction
func TestTAIA_Sub(t *testing.T) {
	tests := []struct {
		t1, t2     TAIA
		expected   TAIA
		shouldFail bool
	}{
		// Simple subtraction (no underflow)
		{
			TAIA{Sec: tai.TAI{Sec: 150}, Nano: 700_000_000, Atto: 400_000_000},
			TAIA{Sec: tai.TAI{Sec: 50}, Nano: 200_000_000, Atto: 300_000_000},
			TAIA{Sec: tai.TAI{Sec: 100}, Nano: 500_000_000, Atto: 100_000_000},
			false,
		},
		// Attosecond underflow
		{
			TAIA{Sec: tai.TAI{Sec: 100}, Nano: 500_000_000, Atto: 100_000_000},
			TAIA{Sec: tai.TAI{Sec: 0}, Nano: 100_000_000, Atto: 200_000_000},
			TAIA{Sec: tai.TAI{Sec: 100}, Nano: 399_999_999, Atto: 900_000_000}, // Borrow from Nano
			false,
		},
		// Nanosecond underflow
		{
			TAIA{Sec: tai.TAI{Sec: 102}, Nano: 100_000_000, Atto: 0},
			TAIA{Sec: tai.TAI{Sec: 1}, Nano: 200_000_000, Atto: 0},
			TAIA{Sec: tai.TAI{Sec: 100}, Nano: 900_000_000, Atto: 0}, // Borrow from Sec
			false,
		},
		// Underflow Case
		{
			TAIA{Sec: tai.TAI{Sec: 10}, Nano: 0, Atto: 0},
			TAIA{Sec: tai.TAI{Sec: 20}, Nano: 0, Atto: 0}, // Attempting to subtract a larger value
			TAIA{}, // Should fail
			true,
		},
	}

	for _, tc := range tests {
		result, err := tc.t1.Sub(tc.t2)
		if (err != nil) != tc.shouldFail {
			t.Errorf("t1 Sub(%v, %v, %v)", tc.t1.Sec.Sec, tc.t1.Nano, tc.t1.Atto)
			t.Errorf("t2 Sub(%v, %v, %v)", tc.t2.Sec.Sec, tc.t2.Nano, tc.t2.Atto)
			t.Errorf("ex Sub(%v, %v, %v)", tc.expected.Sec.Sec, tc.expected.Nano, tc.expected.Atto)
			t.Errorf("re Sub(%v, %v, %v)", result.Sec.Sec, result.Nano, result.Atto)
			// t.Errorf("Sub(%v, %v): unexpected error state: %v", tc.t1, tc.t2, err)
		} else if !tc.shouldFail && (result != tc.expected) {
			t.Errorf("t1 Sub(%v, %v, %v)", tc.t1.Sec.Sec, tc.t1.Nano, tc.t1.Atto)
			t.Errorf("t2 Sub(%v, %v, %v)", tc.t2.Sec.Sec, tc.t2.Nano, tc.t2.Atto)
			t.Errorf("ex Sub(%v, %v, %v)", tc.expected.Sec.Sec, tc.expected.Nano, tc.expected.Atto)
			t.Errorf("re Sub(%v, %v, %v)", result.Sec.Sec, result.Nano, result.Atto)
			// t.Errorf("Sub(%v, %v): expected %v, got %v", tc.t1, tc.t2, tc.expected, result)
		}
	}
}

// TestTAIA_PackUnpack tests the Pack and Unpack functions.
func TestTAIA_PackUnpack(t *testing.T) {
	// Define a TAIA struct with specific values
	original := TAIA{
		Sec:  tai.TAI{Sec: 4611686018427387914},
		Nano: 500_000_000,
		Atto: 250_000_000,
	}

	// Pack the TAIA struct into bytes
	packed := original.Pack()

	// Unpack it back into a TAIA struct
	unpacked := Unpack(packed)

	// Ensure the unpacked values match the original values
	if original.Sec.Sec != unpacked.Sec.Sec {
		t.Errorf("Pack/Unpack mismatch for Sec: expected %d, got %d", original.Sec.Sec, unpacked.Sec.Sec)
	}
	if original.Nano != unpacked.Nano {
		t.Errorf("Pack/Unpack mismatch for Nano: expected %d, got %d", original.Nano, unpacked.Nano)
	}
	if original.Atto != unpacked.Atto {
		t.Errorf("Pack/Unpack mismatch for Atto: expected %d, got %d", original.Atto, unpacked.Atto)
	}
}

func TestTAIA_Less(t *testing.T) {
	tests := []struct {
		u, v     TAIA
		expected bool
	}{
		// u is less than v
		{TAIA{Sec: tai.TAI{Sec: 100}, Nano: 500_000_000, Atto: 0},
			TAIA{Sec: tai.TAI{Sec: 200}, Nano: 0, Atto: 0}, true},

		// u is greater than v
		{TAIA{Sec: tai.TAI{Sec: 300}, Nano: 0, Atto: 0},
			TAIA{Sec: tai.TAI{Sec: 200}, Nano: 500_000_000, Atto: 0}, false},

		// u and v are equal
		{TAIA{Sec: tai.TAI{Sec: 100}, Nano: 500_000_000, Atto: 123_456_789},
			TAIA{Sec: tai.TAI{Sec: 100}, Nano: 500_000_000, Atto: 123_456_789}, false},
	}

	for _, tc := range tests {
		result := tc.u.Less(tc.v)
		if result != tc.expected {
			t.Errorf("u Less(%v, %v, %v)", tc.u.Sec.Sec, tc.u.Nano, tc.u.Atto)
			t.Errorf("u Less(%v, %v, %v)", tc.v.Sec.Sec, tc.v.Nano, tc.v.Atto)
			t.Errorf("Less(%v, %v): expected %v, got %v", tc.u, tc.v, tc.expected, result)
		}
	}
}

func TestTAIA_Half(t *testing.T) {
	tests := []struct {
		input    TAIA
		expected TAIA
	}{
		// Even values
		{
			TAIA{Sec: tai.TAI{Sec: 100}, Nano: 500_000_000, Atto: 200_000_000},
			TAIA{Sec: tai.TAI{Sec: 50}, Nano: 250_000_000, Atto: 100_000_000},
		},
		// Odd seconds, carry-over
		{
			TAIA{Sec: tai.TAI{Sec: 101}, Nano: 0, Atto: 0},
			TAIA{Sec: tai.TAI{Sec: 50}, Nano: 500_000_000, Atto: 0},
		},
		// Odd nanoseconds, carry-over
		{
			TAIA{Sec: tai.TAI{Sec: 100}, Nano: 999_999_999, Atto: 0},
			TAIA{Sec: tai.TAI{Sec: 50}, Nano: 499_999_999, Atto: 500_000_000},
		},
	}

	for _, tc := range tests {
		result := tc.input.Half()
		if result != tc.expected {
			t.Errorf("Half(input): %v %v %v", tc.input.Sec.Sec, tc.input.Nano, tc.input.Atto)
			t.Errorf("Half(expected): %v %v %v", tc.expected.Sec.Sec, tc.expected.Nano, tc.expected.Atto)
			t.Errorf("Half(result): %v %v %v", result.Sec.Sec, result.Nano, result.Atto)
		}
	}
}
