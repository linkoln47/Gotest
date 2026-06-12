package main

import (
	"bytes"
	"fmt"

	stdsha3 "crypto/sha3"
)

const (
	RateBytes    = 136
	shakeDomain  = 0x1F
	cshakeDomain = 0x04
)

func SHAKE256(input []byte, outLenBytes int) []byte {
	xof := NewKeccakXOF(RateBytes, shakeDomain)
	return xof.Sum(input, outLenBytes)
}

func CSHAKE256(input []byte, outLenBytes int, functionName, customization []byte) []byte {
	prefix := append(encodeString(functionName), encodeString(customization)...)
	prefix = bytepad(prefix, RateBytes)

	fullInput := append(prefix, input...)

	xof := NewKeccakXOF(RateBytes, cshakeDomain)
	return xof.Sum(fullInput, outLenBytes)
}

func lEncode(x uint64) []byte {
	n := 1
	for v := x; v > 0xFF; v >>= 8 {
		n++
	}

	digest := make([]byte, n+1)
	digest[0] = byte(n)

	for i := n; i >= 1; i-- {
		digest[i] = byte(x)
		x >>= 8
	}
	return digest
}

func encodeString(s []byte) []byte {
	digest := lEncode(uint64(len(s)) * 8)
	digest = append(digest, s...)
	return digest
}

func bytepad(x []byte, w int) []byte {
	digest := lEncode(uint64(w))
	digest = append(digest, x...)

	for len(digest)%w != 0 {
		digest = append(digest, 0x00)
	}

	return digest
}

func stdSHAKE256(input []byte, outLenBytes int) []byte {
	out := make([]byte, outLenBytes)

	h := stdsha3.NewSHAKE256()
	h.Write(input)
	h.Read(out)

	return out
}

func stdCSHAKE256(input []byte, outLenBytes int, functionName, customization []byte) []byte {
	out := make([]byte, outLenBytes)

	h := stdsha3.NewCSHAKE256(functionName, customization)
	h.Write(input)
	h.Read(out)

	return out
}

func main() {
	msg := []byte("abc")
	outLen := 32

	myHash1 := SHAKE256(msg, outLen)
	officialHash1 := stdSHAKE256(msg, outLen)

	myHash2 := CSHAKE256(msg, outLen, nil, []byte("test"))
	officialHash2 := stdCSHAKE256(msg, outLen, nil, []byte("test"))

	fmt.Println("=== SHAKE256 ===")
	fmt.Printf("%x\n", myHash1)
	fmt.Printf("%x\n", officialHash1)
	fmt.Println(bytes.Equal(myHash1, officialHash1))

	fmt.Println("=== cSHAKE256 ===")
	fmt.Printf("%x\n", myHash2)
	fmt.Printf("%x\n", officialHash2)
	fmt.Println(bytes.Equal(myHash2, officialHash2))

}
