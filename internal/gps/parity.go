package gps

// CheckAndDecodeWord is checking the parity of a 30-bit word and extracting the 24 bytes of data
// d29Star and d30Star are the 29 and 30 bytes of the previous word
func CheckAndDecodeWord(rawWord uint32, d29Star, d30Star uint32) (data uint32, ok bool, newD29, newD30 uint32) {
	// rawWord is a 30-bit word received MSB-first
	// The first bit received (d1) is at bit 29, last parity bit (d30) is at bit 0

	var d [31]uint32 // 1-indexed to match IS-GPS-200 specifications

	// We move the bits in the word into the d array
	for i := 1; i <= 30; i++ {
		d[i] = (rawWord >> (30 - i)) & 1
	}

	// 24 bits of data resolution: Di = di ^ D30*
	// So the data bits are the XOR of the received bits and the last parity bit of the previous word
	var D [31]uint32
	for i := 1; i <= 24; i++ {
		D[i] = d[i] ^ d30Star
	}

	// 6 bits of parity check according to IS-GPS-200 Table 20-XIV
	p1 := d29Star ^ D[1] ^ D[2] ^ D[3] ^ D[5] ^ D[6] ^ D[10] ^ D[11] ^ D[12] ^ D[13] ^ D[14] ^ D[17] ^ D[18] ^ D[20] ^ D[23]
	p2 := d30Star ^ D[2] ^ D[3] ^ D[4] ^ D[6] ^ D[7] ^ D[11] ^ D[12] ^ D[13] ^ D[14] ^ D[15] ^ D[18] ^ D[19] ^ D[21] ^ D[24]
	p3 := d29Star ^ D[1] ^ D[3] ^ D[4] ^ D[5] ^ D[7] ^ D[8] ^ D[12] ^ D[13] ^ D[14] ^ D[15] ^ D[16] ^ D[19] ^ D[20] ^ D[22]
	p4 := d30Star ^ D[2] ^ D[4] ^ D[5] ^ D[6] ^ D[8] ^ D[9] ^ D[13] ^ D[14] ^ D[15] ^ D[16] ^ D[17] ^ D[20] ^ D[21] ^ D[23]
	p5 := d30Star ^ D[1] ^ D[3] ^ D[5] ^ D[6] ^ D[7] ^ D[9] ^ D[10] ^ D[14] ^ D[15] ^ D[16] ^ D[17] ^ D[18] ^ D[21] ^ D[22] ^ D[24]
	p6 := d29Star ^ D[3] ^ D[5] ^ D[6] ^ D[8] ^ D[9] ^ D[10] ^ D[11] ^ D[13] ^ D[15] ^ D[19] ^ D[22] ^ D[23] ^ D[24]

	// Parity check
	valid := (p1 == d[25]) && (p2 == d[26]) && (p3 == d[27]) && (p4 == d[28]) && (p5 == d[29]) && (p6 == d[30])

	// 24 bits of data assembly D1..D24
	var data24 uint32

	for i := 1; i <= 24; i++ {
		data24 = (data24 << 1) | D[i]
	}

	return data24, valid, d[29], d[30]
}
