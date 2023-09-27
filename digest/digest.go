// Implementation of Bitcoin-specific hash functions, including SHA-256, RIPEMD-160 (Hash160), and TaggedHash.
package digest

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"

	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/ripemd160"
)

// DoubleHash computes the SHA-256 hash of data twice, returning the resulting hash.
func DoubleHash(buffer []byte) []byte {
	firstHash := sha256.Sum256(buffer)
	secondHash := sha256.Sum256(firstHash[:])
	return secondHash[:]
}

// SingleHash computes the SHA-256 hash of data.
func SingleHash(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}

// Hash160 computes the RIPEMD160 hash of the SHA-256 hash of data.
func Hash160(buffer []byte) []byte {
	toSh256 := SingleHash(buffer)
	hasher := ripemd160.New()
	hasher.Write(toSh256)
	return hasher.Sum(nil)
}

// TaggedHash computes a tagged hash by prepending a tag byte to the input data
func TaggedHash(data []byte, tag string) []byte {
	tagHash := SingleHash([]byte(tag))
	concat := append(append(tagHash[:], tagHash[:]...), data...)
	hash := SingleHash(concat)
	result := make([]byte, len(hash))
	copy(result, hash[:])

	return result
}

// HmacSHA512 computes the HMAC-SHA-512 hash of the given key and data.
func HmacSHA512(key, data []byte) []byte {
	h := hmac.New(sha512.New, key)
	h.Write(data)
	hash := h.Sum(nil)
	return hash
}

// PbkdfDeriveDigest generates a derived key using the PBKDF2 algorithm with HMAC-SHA-512.
func PbkdfDeriveDigest(mnemonic string, salt string) []byte {
	saltBytes := []byte(salt)
	key := pbkdf2.Key([]byte(mnemonic), saltBytes, 2048, 64, sha512.New)
	return key
}

// GenerateRandom creates a byte slice of the specified 'size' containing random data.
// It generates 'size' random bytes using the crypto/rand package's rand.Read function
// and returns the resulting byte slice. If any error occurs during the random data
// generation, it is returned as an error along with a nil byte slice.
//
// The function is suitable for generating cryptographically secure random data.
func GenerateRandom(size int) ([]byte, error) {
	randomBytes := make([]byte, size)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}
	return randomBytes, nil
}
