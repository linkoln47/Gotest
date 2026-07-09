package main

import (
	"crypto/aes"
	"fmt"
)

func main() {
	key := [32]byte{
		0x00, 0x01, 0x02, 0x03,
		0x04, 0x05, 0x06, 0x07,
		0x08, 0x09, 0x0a, 0x0b,
		0x0c, 0x0d, 0x0e, 0x0f,
		0x10, 0x11, 0x12, 0x13,
		0x14, 0x15, 0x16, 0x17,
		0x18, 0x19, 0x1a, 0x1b,
		0x1c, 0x1d, 0x1e, 0x1f,
	}

	plaintext := [16]byte{
		0x00, 0x11, 0x22, 0x33,
		0x44, 0x55, 0x66, 0x77,
		0x88, 0x99, 0xaa, 0xbb,
		0xcc, 0xdd, 0xee, 0xff,
	}

	expected := [16]byte{
		0x8e, 0xa2, 0xb7, 0xca,
		0x51, 0x67, 0x45, 0xbf,
		0xea, 0xfc, 0x49, 0x90,
		0x4b, 0x49, 0x60, 0x89,
	}

	cipher := New(key)

	ciphertext := cipher.EncryptBlock(plaintext)
	decrypted := cipher.DecryptBlock(ciphertext)

	fmt.Printf("key:        %x\n", key)
	fmt.Printf("plaintext:  %x\n", plaintext)
	fmt.Printf("ciphertext: %x\n", ciphertext)
	fmt.Printf("expected:   %x\n", expected)
	fmt.Printf("decrypted:  %x\n", decrypted)

	fmt.Println("matches test vector:", ciphertext == expected)
	fmt.Println("decrypt correct:    ", decrypted == plaintext)

	// Сравнение со стандартной библиотекой Go.
	stdCipher, _ := aes.NewCipher(key[:])

	var stdCiphertext [16]byte
	stdCipher.Encrypt(stdCiphertext[:], plaintext[:])

	fmt.Printf("crypto/aes: %x\n", stdCiphertext)
	fmt.Println("matches crypto/aes:", ciphertext == stdCiphertext)
}
