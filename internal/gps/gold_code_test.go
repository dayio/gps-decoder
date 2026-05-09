package gps

import (
	"testing"
)

func TestGenerateGoldCode(t *testing.T) {
	prn := 1
	code := GenerateGoldCode(prn)

	// Check the length of the sequence
	if len(code) != 1023 {
		t.Fatalf("Expected length: 1023, got: %d", len(code))
	}

	// Check BPSK signal compliance
	for i, val := range code {
		if val != 1.0 && val != -1.0 {
			t.Errorf("Unexpected BPSK value at index %d : %f. Only 1.0 and -1.0 are allowed", i, val)
		}
	}

	// Verify the first 10 chips (PRN 1 -> Octal 1440 -> Binary 1100100000)
	expectedChips := []float64{
		-1.0, -1.0, 1.0, 1.0, -1.0, 1.0, 1.0, 1.0, 1.0, 1.0,
	}

	for i := range expectedChips {
		if code[i] != expectedChips[i] {
			t.Errorf("Expected chip at index %d: %f, got: %f", i, expectedChips[i], code[i])
		}
	}

	t.Logf("PRN %d validated", prn)
}
