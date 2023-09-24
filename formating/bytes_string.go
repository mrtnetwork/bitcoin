package formating

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

// strip0x removes the "0x" prefix from a hexadecimal string, if it is present.
// It takes a hex string as input and returns a modified string without the "0x" prefix.
// If the input string does not start with "0x", it returns the original string unchanged.
func strip0x(hex string) string {
	if strings.HasPrefix(hex, "0x") {
		return hex[2:]
	}
	return hex
}

// HexToBytes converts a hexadecimal string to a byte slice.
// It takes a hex string as input, removes the "0x" prefix (if present) using strip0x,
// and then decodes the remaining hexadecimal characters into bytes using the hex.DecodeString function.
// The resulting byte slice is returned.
// If there are any decoding errors, an empty byte slice is returned along with an error (ignored).
func HexToBytes(hexStr string) []byte {
	bytes, _ := hex.DecodeString(strip0x(hexStr))
	return bytes
}

// HexToBytesCatch converts a hexadecimal string to a byte slice while handling potential errors.
// It takes a hex string as input, removes the "0x" prefix (if present) using strip0x,
// and then attempts to decode the remaining hexadecimal characters into bytes using the hex.DecodeString function.
// If decoding is successful, it returns the resulting byte slice and a nil error.
// If there are any decoding errors, it returns a nil byte slice along with the encountered error.
func HexToBytesCatch(hexStr string) ([]byte, error) {
	bytes, err := hex.DecodeString(strip0x(hexStr))
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// BytesToHex converts a byte slice into a hexadecimal string representation.
// It takes a byte slice as input and encodes the bytes as a hexadecimal string
// using the hex.EncodeToString function. The resulting hexadecimal string is returned.
func BytesToHex(hexStr []byte) string {
	toHex := hex.EncodeToString(hexStr)
	return toHex
}

// IsLessThanBytes compares two byte slices lexicographically.
// It takes two byte slices, thashedA and thashedB, as input and compares them element by element.
// The function returns true if thashedA is lexicographically less than thashedB.
// If the two slices are equal up to the length of the shorter slice, it considers the shorter one as less.
// This function is useful for comparing byte representations of hashed values or keys.
func IsLessThanBytes(thashedA, thashedB []byte) bool {
	for i := 0; i < len(thashedA) && i < len(thashedB); i++ {
		if thashedA[i] < thashedB[i] {
			return true
		} else if thashedA[i] > thashedB[i] {
			return false
		}
	}
	return len(thashedA) < len(thashedB)
}

// ReverseBytes reverses the order of elements in a byte slice.
// It takes a byte slice 'data' as input and creates a new byte slice 'reversed'
// with the elements of 'data' reversed in order.
// The function uses a two-pointer approach to efficiently reverse the elements.
// The original 'data' slice remains unchanged, and the reversed slice is returned.
func ReverseBytes(data []byte) []byte {
	length := len(data)
	reversed := make([]byte, length)

	for i, j := 0, length-1; i < length; i, j = i+1, j-1 {
		reversed[i] = data[j]
	}

	return reversed
}

// CopyBytes creates a copy of a byte slice 'original' and returns the copied slice.
// It takes a byte slice 'original' as input and creates a new byte slice 'newBytes'
// with the same contents as 'original'. The copy is made using the built-in 'copy' function.
// This function is useful when you want to create a distinct copy of a byte slice to
// prevent modifications to the original slice from affecting the copied one.
func CopyBytes(original []byte) []byte {
	newBytes := make([]byte, len(original))

	copy(newBytes, original)
	return newBytes
}

// BytesToBinary converts a byte slice into a binary string representation.
// It takes a byte slice 'data' as input and creates a new string where each byte
// is represented as an 8-bit binary string. These binary strings are concatenated
// together to form the final binary representation of the input data.
// The resulting binary string is returned.
func BytesToBinary(data []byte) string {
	binaryStrings := make([]string, len(data))

	for _, b := range data {
		// Convert each byte to a binary string with 8 bits and pad with '0' on the left if needed
		binaryStrings = append(binaryStrings, fmt.Sprintf("%08b", b))
	}

	// Join the binary strings to form the final binary representation
	binaryRepresentation := fmt.Sprintf("%s", strings.Join(binaryStrings, ""))

	return binaryRepresentation

}

// BinaryToByte parses a binary string and returns the equivalent integer value.
// It takes a binary string as input and uses strconv.ParseInt to convert it to an integer.
// The base argument is set to 2, indicating binary representation.
// The function returns the integer value and ignores any parsing errors.
func BinaryToByte(binary string) (int, error) {
	result, err := strconv.ParseInt(binary, 2, 64)
	if err != nil {
		return -1, err
	}
	return int(result), nil
}

// FillRange sets all elements in a byte slice 'slice' to the specified 'value'.
// It takes a byte slice 'slice' and a byte 'value' as input. The function iterates
// over each element in the slice and assigns the 'value' to it, effectively filling
// the entire slice with the specified byte value.
func FillRange(slice []byte, value byte) {
	for i := range slice {
		slice[i] = value
	}
}

// Bytes32FromInt converts an integer 'x' into a 32-byte byte slice.
// It takes an integer 'x' as input and creates a 32-byte byte slice
// where the integer is encoded in little-endian order, with the least
// significant byte at the beginning of the slice and the most significant
// byte at the end. The resulting byte slice is returned.
func Bytes32FromInt(x int) []byte {
	result := make([]byte, 32)
	for i := 0; i < 32; i++ {
		result[32-i-1] = byte((x >> (8 * i)) & 0xFF)
	}
	return result
}

// BytesToInt converts a byte slice 'data' into a big.Int.
// It takes a byte slice 'data' as input and creates a new big.Int
// by setting its value to the bytes from the input slice.
// The resulting big.Int represents the integer value encoded by the bytes
// in little-endian order and is returned.
func BytesToInt(data []byte) *big.Int {
	return new(big.Int).SetBytes(data)
}

// PadByteSliceTo32 pads a byte slice 'data' to a length of 32 bytes.
// It takes a byte slice 'data' as input and ensures that it is exactly
// 32 bytes in length. If 'data' is already 32 bytes or longer, it is
// returned as is. If 'data' is shorter than 32 bytes, it is padded with
// zero bytes at the beginning to achieve the desired length.
// The resulting 32-byte byte slice is returned.
func PadByteSliceTo32(data []byte) []byte {
	if len(data) >= 32 {
		return data[:32]
	}

	paddedData := make([]byte, 32)
	copy(paddedData[32-len(data):], data)
	return paddedData
}

// xorBytes performs a bitwise XOR operation between two byte slices 'a' and 'b'.
// It takes two byte slices 'a' and 'b' as input and computes the bitwise XOR of
// each corresponding pair of bytes. The input slices must have the same length.
// The resulting byte slice contains the result of the XOR operation and is returned.
// If the input slices have different lengths, it raises a panic.
func XorBytes(a, b []byte) []byte {
	if len(a) != len(b) {
		panic("Input slices must have the same length")
	}

	result := make([]byte, len(a))

	for i := 0; i < len(a); i++ {
		result[i] = a[i] ^ b[i]
	}

	return result
}
