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

// encodeLength encodes an integer 'length' as a variable-length prefix for a data structure.
// It is commonly used in various data encoding formats, such as ASN.1 DER or TLV encoding.
// The function takes an integer 'length' as input and encodes it into a byte slice,
// following the variable-length encoding rules.
//
// If 'length' is less than 128, it is encoded as a single byte with its value.
// If 'length' is 128 or greater, it is encoded as a multi-byte sequence. The first byte
// indicates the number of bytes used to encode the length, with the most significant bit set.
// The subsequent bytes represent the actual length value in big-endian format.
//
// The resulting encoded byte slice is returned.
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

// encodeInteger encodes a non-negative big.Int 'r' as an ASN.1 DER-encoded integer.
// The function first checks if 'r' is negative, and if so, it raises a panic since negative
// numbers are not supported in this code.
//
// If 'r' is non-negative, it converts it to a hexadecimal string 'h' and ensures that 'h' has
// an even number of hex digits. If not, a leading zero is added to 'h'.
//
// The function then checks the value of the first byte in the resulting byte slice 's':
//   - If 's[0]' is less than or equal to 0x7F, 's' is directly encoded as a positive integer,
//     and its length is determined using the 'encodeLength' function.
//   - If 's[0]' is greater than 0x7F, 's' is encoded as a positive integer with an additional
//     leading zero byte, and its length is adjusted accordingly.
//
// The resulting ASN.1 DER-encoded integer is returned as a byte slice.
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
