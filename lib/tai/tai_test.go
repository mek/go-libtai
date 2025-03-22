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
	if taiNow.X < expectedTai-2 || taiNow.X > expectedTai+2 {
		t.Errorf("Now() returned unexpected TAI time: got %d, expected around %d", taiNow.X, expectedTai)
	}
}

func TestTAI_PackUnpack(t *testing.T) {
	taiOriginal := TAI{X: 4611686018427387914}
	packed := taiOriginal.Pack()
	unpacked := UnPack(packed)

	if taiOriginal.X != unpacked.X {
		t.Errorf("Pack/Unpack failed: expected %d, got %d", taiOriginal.X, unpacked.X)
	}
}

func TestTAI_Add(t *testing.T) {
	tests := []struct {
		u, v       TAI
		expected   TAI
		shouldFail bool
	}{
		// Normal addition
		{TAI{X: 100}, TAI{X: 50}, TAI{X: 150}, false},
		// Large addition (no overflow)
		{TAI{X: math.MaxUint64 - 10}, TAI{X: 5}, TAI{X: math.MaxUint64 - 5}, false},
		// Overflow case
		{TAI{X: math.MaxUint64 - 1}, TAI{X: 2}, TAI{}, true},
		// Boundary Add
		{TAI{X: math.MaxUint64 - 1}, TAI{X: 1}, TAI{X: math.MaxUint64}, false},
		// Boundary Max
		{TAI{X: math.MaxUint64}, TAI{X: 1}, TAI{}, true},
	}

	for _, tc := range tests {
		result, err := tc.u.Add(tc.v)
		if (err != nil) != tc.shouldFail {
			t.Errorf("Add(%d, %d): unexpected error state: %v", tc.u.X, tc.v.X, err)
		} else if !tc.shouldFail && (result != tc.expected) {
			t.Errorf("Add(%d, %d): expected %d, got %d", tc.u.X, tc.v.X, tc.expected.X, result.X)
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
		{TAI{X: 100}, TAI{X: 50}, TAI{X: 50}, false},
		// Subtracting from itself (zero case)
		{TAI{X: 200}, TAI{X: 200}, TAI{X: 0}, false},
		// Underflow case
		{TAI{X: 50}, TAI{X: 100}, TAI{}, true},
		// Boundary Sub
		{TAI{X: 1}, TAI{X: 1}, TAI{X: 0}, false},
		// Boundary Sub
		{TAI{X: 0}, TAI{X: 1}, TAI{}, true},
	}

	for _, tc := range tests {
		result, err := tc.u.Sub(tc.v)
		if (err != nil) != tc.shouldFail {
			t.Errorf("Sub(%d, %d): unexpected error state: %v", tc.u.X, tc.v.X, err)
		} else if !tc.shouldFail && (result != tc.expected) {
			t.Errorf("Sub(%d, %d): expected %d, got %d", tc.u.X, tc.v.X, tc.expected.X, result.X)
		}
	}
}

func TestTAI_Less(t *testing.T) {
	tests := []struct {
		t1, t2   TAI
		expected bool
	}{
		// t1 is less than t2
		{TAI{X: 100}, TAI{X: 200}, true},
		// t1 is equal to t2
		{TAI{X: 100}, TAI{X: 100}, false},
		// t1 is greater than t2
		{TAI{X: 200}, TAI{X: 100}, false},
	}

	for _, tc := range tests {
		result := tc.t1.Less(tc.t2)
		if result != tc.expected {
			t.Errorf("Less(%d, %d): expected %v, got %v", tc.t1.X, tc.t2.X, tc.expected, result)
		}
	}
}
