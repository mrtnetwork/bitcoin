// RFC6979KCalculator provides a deterministic K-value generator based on RFC 6979
// for use in Bitcoin cryptographic operations.

package ecc

import (
	"bitcoin/formating"
	"crypto/hmac"
	"crypto/sha256"
	"hash"
	"math/big"
)

type RFC6979KCalculator struct {
	mac  hash.Hash
	K, V []byte
	n    *big.Int
}

// NewRFC6979KCalculator creates a new RFC6979KCalculator instance
// initialized with the provided parameters.
func NewRFC6979KCalculator(mac hash.Hash, n, d *big.Int, message []byte, entryPointes []byte) *RFC6979KCalculator {
	kCalculator := &RFC6979KCalculator{
		mac: mac,
		K:   make([]byte, mac.Size()),
		V:   make([]byte, mac.Size()),
		n:   new(big.Int).Set(n),
	}
	kCalculator.init(d, message, entryPointes)
	return kCalculator
}

func (k *RFC6979KCalculator) init(d *big.Int, message []byte, entryPointes []byte) {
	zeroBytes := []byte{0x00}
	oneBytes := []byte{0x01}
	formating.FillRange(k.K, 0x00)
	formating.FillRange(k.V, 0x01)
	x := make([]byte, (k.n.BitLen()+7)/8)
	dBytes := d.Bytes()
	copy(x[len(x)-len(dBytes):], dBytes)

	m := make([]byte, (k.n.BitLen()+7)/8)
	mInt := new(big.Int)
	mInt.SetBytes(message)
	if mInt.Cmp(k.n) > 0 {
		mInt.Sub(mInt, k.n)
	}
	mBytes := mInt.Bytes()
	copy(m[len(m)-len(mBytes):], mBytes)

	k.mac = hmac.New(sha256.New, k.K)
	k.mac.Write(k.V)
	k.mac.Write(zeroBytes)
	k.mac.Write(x)
	k.mac.Write(m)
	if entryPointes != nil {
		k.mac.Write(entryPointes)
	}
	k.mac.Sum(k.K[:0])
	k.mac = hmac.New(sha256.New, k.K)
	k.mac.Write(k.V)
	k.mac.Sum(k.V[:0])

	k.mac = hmac.New(sha256.New, k.K)
	k.mac.Write(k.V)
	k.mac.Write(oneBytes)
	k.mac.Write(x)
	k.mac.Write(m)
	if entryPointes != nil {
		k.mac.Write(entryPointes)
	}
	k.mac.Sum(k.K[:0])

	k.mac = hmac.New(sha256.New, k.K)
	k.mac.Write(k.V)
	k.mac.Sum(k.V[:0])
}

func (k *RFC6979KCalculator) NextK() *big.Int {
	t := make([]byte, (k.n.BitLen()+7)/8)

	for {
		tOff := 0
		for tOff < len(t) {
			k.mac.Reset()
			k.mac.Write(k.V)
			k.mac.Sum(k.V[:0])

			if len(t)-tOff < len(k.V) {
				copy(t[tOff:], k.V)
				tOff += len(t) - tOff
			} else {
				copy(t[tOff:], k.V)
				tOff += len(k.V)
			}
		}

		kInt := new(big.Int)
		kInt.SetBytes(t)

		if kInt.Cmp(big.NewInt(0)) == 0 || kInt.Cmp(k.n) >= 0 {
			k.mac.Reset()
			k.mac.Write(k.V)
			k.mac.Write([]byte{0x00})
			k.mac.Sum(k.K[:0])

			k.mac.Reset()
			k.mac.Write(k.K)
			k.mac.Write(k.V)
			k.mac.Sum(k.V[:0])
		} else {
			return kInt
		}
	}
}
