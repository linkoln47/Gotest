package main

import (
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

func padMasg(msg []byte) []byte {
	originalBitLen := uint64(len(msg)) * 8 // Create slice to save msg in bits

	padded := make([]byte, len(msg))
	copy(padded, msg)
	padded = append(padded, 0x80) // Append 0x80 (10000000 in binary) to msg

	for len(padded)%64 != 56 { // Append 0x00 until length of msg is 56 mod 64
		padded = append(padded, 0x00)
	}
	for i := 7; i >= 0; i-- { // 56 bytes of data + 8 bytes for orig msg
		padded = append(padded, byte(originalBitLen>>(uint(i)*8)))
	}
	return padded
}

func lRotate(x uint32, n uint) uint32 { // x << n shifts to the left
	return (x << n) | (x >> (32 - n)) // x >> (32 - n) adding missed bits to the right
}

func blockToWords(block []byte) [80]uint32 { // Create 80 words (32 bits per word) from each block
	var w [80]uint32 // sha1 works with uint32, not bytes

	for i := 0; i < 16; i++ { // first 16 words of 4 bytes
		w[i] = binary.BigEndian.Uint32(block[i*4 : i*4+4])
	}

	for i := 16; i < 80; i++ { // expands 16 to 80 words by XORing and left rotating by 1 bit
		w[i] = lRotate(w[i-3]^w[i-8]^w[i-14]^w[i-16], 1) // new word depends on 4 prev words
	}

	return w
}

func prcBlock(block []byte, h0, h1, h2, h3, h4 uint32) (uint32, uint32, uint32, uint32, uint32) {
	w := blockToWords(block) // create 80 words from block
	/*
	   Creating 5 var of state.
	   a...e are temp reg for each round
	   h0...h4 are state of hash
	*/
	a := h0
	b := h1
	c := h2
	d := h3
	e := h4

	for i := 0; i < 80; i++ {
		var f uint32 // logic func
		var k uint32 // constant

		if i <= 19 { // if b = 1, take c, if b = 0, take d
			f = (b & c) | (^b & d)
			k = 0x5A827999
		} else if i <= 39 { // XOR of b, c, d
			f = (b ^ c ^ d)
			k = 0x6ED9EBA1
		} else if i <= 59 { // majority func: choses most common bit
			f = (b & c) | (b & d) | (c & d)
			k = 0x8F1BBCDC
		} else {
			f = (b ^ c ^ d)
			k = 0xCA62C1D6
		}
		// mixing current state, logic func, old reg, round const, msg slice
		temp := lRotate(a, 5) + f + e + k + w[i] // + is addition mod 2^32
		// update state for next round
		e = d
		d = c
		c = lRotate(b, 30)
		b = a
		a = temp
	}
	// feed-forward
	h0 += a
	h1 += b
	h2 += c
	h3 += d
	h4 += e

	return h0, h1, h2, h3, h4
}

func sha1Hash(msg []byte) [20]byte {
	padded := padMasg(msg)

	h0 := uint32(0x67452301)
	h1 := uint32(0xEFCDAB89)
	h2 := uint32(0x98BADCFE)
	h3 := uint32(0x10325476)
	h4 := uint32(0xC3D2E1F0)

	for i := 0; i < len(padded); i += 64 {
		block := padded[i : i+64]
		h0, h1, h2, h3, h4 = prcBlock(block, h0, h1, h2, h3, h4)
	}
	var digest [20]byte

	binary.BigEndian.PutUint32(digest[0:4], h0)
	binary.BigEndian.PutUint32(digest[4:8], h1)
	binary.BigEndian.PutUint32(digest[8:12], h2)
	binary.BigEndian.PutUint32(digest[12:16], h3)
	binary.BigEndian.PutUint32(digest[16:20], h4)

	return digest
}

func main() {
	msg := []byte("abc")
	digest := sha1Hash(msg)
	fmt.Println(hex.EncodeToString(digest[:]))
	digest = sha1.Sum(msg)
	fmt.Println(hex.EncodeToString(digest[:]))
}
