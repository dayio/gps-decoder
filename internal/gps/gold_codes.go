package gps

func GenerateGoldCode(prn int) []float64 {

	taps := map[int][2]int{
		1: {2, 6}, 2: {3, 7}, 3: {4, 8}, 4: {5, 9},
		5: {1, 9}, 6: {2, 10}, 7: {1, 8}, 8: {2, 9},
		9: {3, 10}, 10: {2, 3}, 11: {3, 4}, 12: {5, 6},
		13: {6, 7}, 14: {7, 8}, 15: {8, 9}, 16: {9, 10},
		17: {1, 4}, 18: {2, 5}, 19: {3, 6}, 20: {4, 7},
		21: {5, 8}, 22: {6, 9}, 23: {1, 3}, 24: {4, 6},
		25: {5, 7}, 26: {6, 8}, 27: {7, 9}, 28: {8, 10},
		29: {1, 6}, 30: {2, 7}, 31: {3, 8}, 32: {4, 9},
	}

	code := make([]float64, 1023)
	prnTaps := taps[prn]

	register1 := make([]int, 10)
	register2 := make([]int, 10)

	for i := range 10 {
		register1[i] = 1
		register2[i] = 1
	}

	for i := range 1023 {
		outputRegister1 := register1[9]
		outputRegister2 := register2[prnTaps[0]-1] ^ register2[prnTaps[1]-1]

		if outputRegister1^outputRegister2 == 0 {
			code[i] = 1.0
		} else {
			code[i] = -1.0
		}

		polynomial1 := register1[2] ^ register1[9]
		polynomial2 := register2[1] ^ register2[2] ^ register2[5] ^ register2[7] ^ register2[8] ^ register2[9]

		copy(register1[1:], register1[:9])
		copy(register2[1:], register2[:9])

		register1[0] = polynomial1
		register2[0] = polynomial2
	}

	return code
}
