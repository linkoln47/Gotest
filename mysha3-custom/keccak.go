package main

func padMsg(input []byte, rateBytes int, domain byte) []byte {
	paddedLen := ((len(input) / rateBytes) + 1) * rateBytes

	padded := make([]byte, paddedLen)
	copy(padded, input)

	padded[len(input)] ^= domain
	padded[paddedLen-1] ^= 0x80

	return padded
}

func absorbBlock(state *[25]uint64, block []byte) {
	for i := 0; i < len(block)/8; i++ {
		j := i * 8

		lane := uint64(block[j]) |
			uint64(block[j+1])<<8 |
			uint64(block[j+2])<<16 |
			uint64(block[j+3])<<24 |
			uint64(block[j+4])<<32 |
			uint64(block[j+5])<<40 |
			uint64(block[j+6])<<48 |
			uint64(block[j+7])<<56

		state[i] ^= lane
	}
}

func squeeze(state *[25]uint64, rateBytes, outLenBytes int) []byte {
	digest := make([]byte, 0, outLenBytes)

	for len(digest) < outLenBytes {
		need := outLenBytes - len(digest)
		if need > rateBytes {
			need = rateBytes
		}

		for i := 0; i < need/8; i++ {
			lane := state[i]
			digest = append(digest,
				byte(lane),
				byte(lane>>8),
				byte(lane>>16),
				byte(lane>>24),
				byte(lane>>32),
				byte(lane>>40),
				byte(lane>>48),
				byte(lane>>56),
			)
		}

		if len(digest) < outLenBytes {
			keccakF1600(state)
		}
	}
	return digest
}

func keccakXOF(rateBytes int, input []byte, domain byte, outLenBytes int) []byte {
	var state [25]uint64

	padded := padMsg(input, rateBytes, domain)

	for offset := 0; offset < len(padded); offset += rateBytes {
		block := padded[offset : offset+rateBytes]
		absorbBlock(&state, block)
		keccakF1600(&state)
	}
	return squeeze(&state, rateBytes, outLenBytes)
}

func rotl(x uint64, n uint) uint64 {
	return (x << n) | (x >> (64 - n))
}

func theta(state *[25]uint64) {
	var c [5]uint64
	var d [5]uint64

	for x := 0; x < 5; x++ {
		c[x] = state[x] ^
			state[x+5] ^
			state[x+10] ^
			state[x+15] ^
			state[x+20]
	}

	for x := 0; x < 5; x++ {
		d[x] = c[(x+4)%5] ^ rotl(c[(x+1)%5], 1)
	}

	for y := 0; y < 5; y++ {
		for x := 0; x < 5; x++ {
			state[x+5*y] ^= d[x]
		}
	}
}

func rho(state *[25]uint64) {
	rotationOffsets := [25]uint{
		0, 1, 62, 28, 27,
		36, 44, 6, 55, 20,
		3, 10, 43, 25, 39,
		41, 45, 15, 21, 8,
		18, 2, 61, 56, 14,
	}

	for i := 0; i < 25; i++ {
		state[i] = rotl(state[i], rotationOffsets[i])
	}
}

func pi(state *[25]uint64) {
	var newState [25]uint64

	for y := 0; y < 5; y++ {
		for x := 0; x < 5; x++ {
			newX := y
			newY := (2*x + 3*y) % 5

			newState[newX+5*newY] = state[x+5*y]
		}
	}

	*state = newState
}

func chi(state *[25]uint64) {
	var newState [25]uint64

	for y := 0; y < 5; y++ {
		for x := 0; x < 5; x++ {
			current := state[x+5*y]
			next := state[((x+1)%5)+5*y]
			nextNext := state[((x+2)%5)+5*y]

			newState[x+5*y] = current ^ ((^next) & nextNext)
		}
	}

	*state = newState
}

func iota(state *[25]uint64, round int) {
	roundConstants := [24]uint64{
		0x0000000000000001, 0x0000000000008082, 0x800000000000808a, 0x8000000080008000,
		0x000000000000808b, 0x0000000080000001, 0x8000000080008081, 0x8000000000008009,
		0x000000000000008a, 0x0000000000000088, 0x0000000080008009, 0x000000008000000a,
		0x000000008000808b, 0x800000000000008b, 0x8000000000008089, 0x8000000000008003,
		0x8000000000008002, 0x8000000000000080, 0x000000000000800a, 0x800000008000000a,
		0x8000000080008081, 0x8000000000008080, 0x0000000080000001, 0x8000000080008008,
	}

	if round < 0 || round >= 24 {
		panic("Keccak-f[1600] round must be between 0 and 23")
	}

	state[0] ^= roundConstants[round]
}

func keccakF1600(state *[25]uint64) {
	for round := 0; round < 24; round++ {
		theta(state)
		rho(state)
		pi(state)
		chi(state)
		iota(state, round)
	}
}
