package ecc

import (
	"encoding/hex"
	"fmt"
	"math/big"
)

// ListBigIntToDER converts a list of BigInt values into a DER-encoded byte slice.
// implements for Bitcoin transaction signature
func ListBigIntToDER(bigIntList []*big.Int) []byte {
	encodedIntegers := make([][]byte, len(bigIntList))

	for i, bi := range bigIntList {
		bytes := encodeInteger(bi)
		encodedIntegers[i] = bytes
	}

	length := 0
	for _, e := range encodedIntegers {
		length += len(e)
	}

	lengthBytes := encodeLength(length)
	contentBytes := make([]byte, 0)

	for _, e := range encodedIntegers {
		contentBytes = append(contentBytes, e...)
	}

	derBytes := make([]byte, 0)
	derBytes = append(derBytes, 0x30) // DER SEQUENCE tag
	derBytes = append(derBytes, lengthBytes...)
	derBytes = append(derBytes, contentBytes...)

	return derBytes
}

func encodeLength(length int) []byte {
	if length < 128 {
		return []byte{byte(length)}
	} else {
		lengthBytes := fmt.Sprintf("%X", length)
		if len(lengthBytes)%2 != 0 {
			lengthBytes = "0" + lengthBytes
		}
		lengthValue, _ := hex.DecodeString(lengthBytes)
		return append([]byte{byte(0x80 | (len(lengthValue)))}, lengthValue...)
	}
}

func encodeInteger(r *big.Int) []byte {
	if r.Sign() < 0 {
		// Negative numbers are not supported in this code
		panic("Negative numbers are not supported")
	}

	h := fmt.Sprintf("%X", r)
	if len(h)%2 != 0 {
		h = "0" + h
	}
	s, _ := hex.DecodeString(h)

	if s[0] <= 0x7F {
		return append([]byte{0x02}, append(encodeLength(len(s)), s...)...)
	} else {
		return append([]byte{0x02}, append(encodeLength(len(s)+1), append([]byte{0x00}, s...)...)...)
	}
}
