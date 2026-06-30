package main

import (
	"crypto/hmac"
	stdsha256 "crypto/sha256"
	"fmt"
	"mysha256"
)

func myHMACSHA256(key, msg []byte) []byte {
	const blockSize = mysha256.BlockSize

	k0 := make([]byte, blockSize)

	if len(key) > blockSize {
		keyHash := mysha256.Sum(key)
		copy(k0, keyHash[:])
	} else {
		copy(k0, key)
	}

	ipad := make([]byte, blockSize)
	opad := make([]byte, blockSize)

	for i := 0; i < blockSize; i++ {
		ipad[i] = k0[i] ^ 0x36
		opad[i] = k0[i] ^ 0x5c
	}

	innerMsg := append(ipad, msg...)
	innerSum := mysha256.Sum(innerMsg)

	outerMsg := append(opad, innerSum[:]...)
	outerSum := mysha256.Sum(outerMsg)

	return outerSum[:]
}

func stdHMACSHA256(key, msg []byte) []byte {
	mac := hmac.New(stdsha256.New, key)
	mac.Write(msg)

	return mac.Sum(nil)
}

func main() {
	key := []byte("secret-key")
	msg := []byte("hello")

	myMAC := myHMACSHA256(key, msg)
	stdMAC := stdHMACSHA256(key, msg)

	fmt.Printf("%x\n", myMAC)
	fmt.Printf("%x\n", stdMAC)

	fmt.Println(hmac.Equal(myMAC, stdMAC))
}
