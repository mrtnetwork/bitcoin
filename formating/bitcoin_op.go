package formating

import (
	"encoding/hex"
	"math/bits"
)

// Converts data to appropriate OP_PUSHDATA OP code including length
// 0x01-0x4b           -> just length plus data bytes
// 0x4c-0xff           -> OP_PUSHDATA1 plus 1-byte-length plus data bytes
// 0x0100-0xffff       -> OP_PUSHDATA2 plus 2-byte-length plus data bytes
// 0x010000-0xffffffff -> OP_PUSHDATA4 plus 4-byte-length plus data bytes
// Also note that according to standarardness rules (BIP-62) the minimum
// possible PUSHDATA operator must be used!
func OpPushData(hexData string) []byte {
	dataBytes, err := hex.DecodeString(hexData)
	if err != nil {
		panic("Invalid hex string")
	}

	dataLength := len(dataBytes)

	switch {
	case dataLength < 0x4c:
		return append([]byte{byte(dataLength)}, dataBytes...)
	case dataLength < 0xff:
		return append([]byte{0x4c, byte(dataLength)}, dataBytes...)
	case dataLength < 0xffff:
		return append([]byte{0x4d, byte(dataLength), byte(dataLength >> 8)}, dataBytes...)
	case dataLength < 0xffffffff:
		return append([]byte{0x4e, byte(dataLength), byte(dataLength >> 8), byte(dataLength >> 16), byte(dataLength >> 24)}, dataBytes...)
	default:
		panic("data too large. cannot push into script")
	}
}

// Converts integer to bytes; as signed little-endian integer
// Currently supports only positive integers
func PushInteger(integer int) []byte {
	if integer < 0 {
		panic("integer is currently required to be positive")
	}

	// Calculate the number of bytes required to represent the integer
	numberOfBytes := (bits.Len(uint(integer)) + 7) / 8

	// Convert to little-endian bytes
	integerBytes := make([]byte, numberOfBytes)
	for i := 0; i < numberOfBytes; i++ {
		integerBytes[i] = byte((integer >> (i * 8)) & 0xFF)
	}

	// If the last bit is set, add a sign byte to signify a positive integer
	if (integer & (1 << uint((numberOfBytes*8)-1))) != 0 {
		integerBytes = append(integerBytes, 0x00)
	}

	// Encode as a variable-length byte slice
	result := append([]byte{byte(len(integerBytes))}, integerBytes...)

	return result
}
