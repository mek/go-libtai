//go:build unit

package tai

import (
	"math"
	"testing"
	"time"
)

func TestTAI_Now(t *testing.T) {
	taiNow := Now()
	unixNow := uint64(time.Now().Unix())

	// Expected TAI time should be within 2 seconds of computed value
	expectedTai := TaiUnixOffset + unixNow
	if taiNow.Sec < expectedTai-2 || taiNow.Sec > expectedTai+2 {
		t.Errorf("Now() returned unexpected TAI time: got %d, expected around %d", taiNow.Sec, expectedTai)
	}
}

func TestTAI_PackUnpack(t *testing.T) {
	taiOriginal := TAI{Sec: 4611686018427387914}
	packed := taiOriginal.Pack()
	unpacked := Unpack(packed)

	if taiOriginal.Sec != unpacked.Sec {
		t.Errorf("Pack/Unpack failed: expected %d, got %d", taiOriginal.Sec, unpacked.Sec)
	}
}

func TestTAI_Add(t *testing.T) {
	tests := []struct {
		u, v       TAI
		expected   TAI
		shouldFail bool
	}{
		// Normal addition
		{TAI{Sec: 100}, TAI{Sec: 50}, TAI{Sec: 150}, false},
		// Large addition (no overflow)
		{TAI{Sec: math.MaxUint64 - 10}, TAI{Sec: 5}, TAI{Sec: math.MaxUint64 - 5}, false},
		// Overflow case
		{TAI{Sec: math.MaxUint64 - 1}, TAI{Sec: 2}, TAI{}, true},
		// Boundary Add
		{TAI{Sec: math.MaxUint64 - 1}, TAI{Sec: 1}, TAI{Sec: math.MaxUint64}, false},
		// Boundary Max
		{TAI{Sec: math.MaxUint64}, TAI{Sec: 1}, TAI{}, true},
	}

	for _, tc := range tests {
		result, err := tc.u.Add(tc.v)
		if (err != nil) != tc.shouldFail {
			t.Errorf("Add(%d, %d): unexpected error state: %v", tc.u.Sec, tc.v.Sec, err)
		} else if !tc.shouldFail && (result != tc.expected) {
			t.Errorf("Add(%d, %d): expected %d, got %d", tc.u.Sec, tc.v.Sec, tc.expected.Sec, result.Sec)
		}
	}
}

func TestTAI_Sub(t *testing.T) {
	tests := []struct {
		u, v       TAI
		expected   TAI
		shouldFail bool
	}{
		// Normal subtraction
		{TAI{Sec: 100}, TAI{Sec: 50}, TAI{Sec: 50}, false},
		// Subtracting from itself (zero case)
		{TAI{Sec: 200}, TAI{Sec: 200}, TAI{Sec: 0}, false},
		// Underflow case
		{TAI{Sec: 50}, TAI{Sec: 100}, TAI{}, true},
		// Boundary Sub
		{TAI{Sec: 1}, TAI{Sec: 1}, TAI{Sec: 0}, false},
		// Boundary Sub
		{TAI{Sec: 0}, TAI{Sec: 1}, TAI{}, true},
	}

	for _, tc := range tests {
		result, err := tc.u.Sub(tc.v)
		if (err != nil) != tc.shouldFail {
			t.Errorf("Sub(%d, %d): unexpected error state: %v", tc.u.Sec, tc.v.Sec, err)
		} else if !tc.shouldFail && (result != tc.expected) {
			t.Errorf("Sub(%d, %d): expected %d, got %d", tc.u.Sec, tc.v.Sec, tc.expected.Sec, result.Sec)
		}
	}
}

func TestTAI_Less(t *testing.T) {
	tests := []struct {
		t1, t2   TAI
		expected bool
	}{
		// t1 is less than t2
		{TAI{Sec: 100}, TAI{Sec: 200}, true},
		// t1 is equal to t2
		{TAI{Sec: 100}, TAI{Sec: 100}, false},
		// t1 is greater than t2
		{TAI{Sec: 200}, TAI{Sec: 100}, false},
	}

	for _, tc := range tests {
		result := tc.t1.Less(tc.t2)
		if result != tc.expected {
			t.Errorf("Less(%d, %d): expected %v, got %v", tc.t1.Sec, tc.t2.Sec, tc.expected, result)
		}
	}
}
