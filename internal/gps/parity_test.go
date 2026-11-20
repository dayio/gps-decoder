package gps

import (
	"testing"
)

// When using the first word TLM we know the last two bits are 0
func TestCheckAndDecodeWord(t *testing.T) {
	// D30* = 0, D29* = 0
	d29Star := uint32(0)
	d30Star := uint32(0)

	// D1..D8 = 1,0,0,0,1,0,1,1
	// D9..D24 = 0 // fake / empty word
	var D [31]uint32
	preambleBits := []uint32{1, 0, 0, 0, 1, 0, 1, 1}
	for i, b := range preambleBits {
		D[i+1] = b
	}

	p1 := d29Star ^ D[1] ^ D[2] ^ D[3] ^ D[5] ^ D[6] ^ D[10] ^ D[11] ^ D[12] ^ D[13] ^ D[14] ^ D[17] ^ D[18] ^ D[20] ^ D[23]
	p2 := d30Star ^ D[2] ^ D[3] ^ D[4] ^ D[6] ^ D[7] ^ D[11] ^ D[12] ^ D[13] ^ D[14] ^ D[15] ^ D[18] ^ D[19] ^ D[21] ^ D[24]
	p3 := d29Star ^ D[1] ^ D[3] ^ D[4] ^ D[5] ^ D[7] ^ D[8] ^ D[12] ^ D[13] ^ D[14] ^ D[15] ^ D[16] ^ D[19] ^ D[20] ^ D[22]
	p4 := d30Star ^ D[2] ^ D[4] ^ D[5] ^ D[6] ^ D[8] ^ D[9] ^ D[13] ^ D[14] ^ D[15] ^ D[16] ^ D[17] ^ D[20] ^ D[21] ^ D[23]
	p5 := d30Star ^ D[1] ^ D[3] ^ D[5] ^ D[6] ^ D[7] ^ D[9] ^ D[10] ^ D[14] ^ D[15] ^ D[16] ^ D[17] ^ D[18] ^ D[21] ^ D[22] ^ D[24]
	p6 := d29Star ^ D[3] ^ D[5] ^ D[6] ^ D[8] ^ D[9] ^ D[10] ^ D[11] ^ D[13] ^ D[15] ^ D[19] ^ D[22] ^ D[23] ^ D[24]

	var rawWord uint32

	for i := 1; i <= 24; i++ {
		rawWord = (rawWord << 1) | (D[i] ^ d30Star)
	}

	rawWord = (rawWord << 1) | p1
	rawWord = (rawWord << 1) | p2
	rawWord = (rawWord << 1) | p3
	rawWord = (rawWord << 1) | p4
	rawWord = (rawWord << 1) | p5
	rawWord = (rawWord << 1) | p6

	data, valid, newD29, newD30 := CheckAndDecodeWord(rawWord, d29Star, d30Star)

	if !valid {
		t.Fatalf("Expected word parity to be valid")
	}

	if (data >> 16) != 0x8B {
		t.Errorf("Expected preamble 0x8B in decoded data, got 0x%X", data>>16)
	}

	if newD29 != p5 || newD30 != p6 {
		t.Errorf("Expected newD29=%d, newD30=%d, got %d, %d", p5, p6, newD29, newD30)
	}
}
