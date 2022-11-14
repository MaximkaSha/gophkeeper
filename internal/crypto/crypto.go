// Package crypto implements GOST "Kuznechik" ciphering.
package crypto

import (
	crypto "github.com/ddulesov/gogost/gost3412128"
)

// Crypto struct - ciphering class.
type Crypto struct {
	Crypto crypto.Cipher
}

// Constructor of Crypto.
// Key size 32 byte.
func NewCrypto(key []byte) *Crypto {
	return &Crypto{
		Crypto: *crypto.NewCipher(key),
	}
}

// Encrypt data.
func (c *Crypto) Encrypt(data []byte) []byte {
	dst := make([]byte, len(data))
	if a := len(data) % c.Crypto.BlockSize(); a != 0 {
		dst = append(dst, make([]byte, (c.Crypto.BlockSize()-a))...)
		data = append(data, make([]byte, (c.Crypto.BlockSize()-a))...)
	}
	for i := 0; i < len(data); i += 16 {
		c.Crypto.Encrypt(dst[i:c.Crypto.BlockSize()+i], data[i:c.Crypto.BlockSize()+i])
	}
	return dst

}

// Decrypt data.
func (c *Crypto) Decrypt(data []byte) []byte {
	dst := make([]byte, len(data))
	for i := 0; i < len(data); i += c.Crypto.BlockSize() {
		c.Crypto.Decrypt(dst[i:c.Crypto.BlockSize()+i], data[i:c.Crypto.BlockSize()+i])
	}
	var counter int64
	for i := len(dst) - 1; i > 0; i -= 1 {
		if dst[i] == 0x00 {
			counter++
		} else {
			break
		}
	}
	dst = dst[:len(dst)-int(counter)]
	return dst
}
