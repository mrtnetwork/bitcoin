package formating

import (
	"encoding/binary"
	"math/big"
)

const largeUint64 = ^uint64(0)

// EncodeVarint encodes an integer 'i' as a variable-length byte representation.
// The function takes an integer 'i' as input and encodes it into a byte slice,
// using variable-length encoding based on the value of 'i'. The resulting byte
// slice represents 'i' with 1, 3, 5, or 9 bytes, depending on its magnitude.
//
// If 'i' is less than 253, it uses a single byte to represent it.
// If 'i' is between 253 and 65535 (0x10000), it uses 3 bytes.
// If 'i' is between 65536 (0x10000) and 4294967295 (0x100000000), it uses 5 bytes.
// If 'i' is less than 'largeUint64' (a predefined constant), it uses 9 bytes.
//
// If 'i' is too large to be represented, it raises a panic.
//
// The encoding is done using little-endian byte order.
// The resulting byte slice is returned.
func EncodeVarint(i int) []byte {
	if i < 253 {
		return []byte{byte(i)}
	} else if i < 0x10000 {
		bytes := make([]byte, 3)
		bytes[0] = 0xfd
		binary.LittleEndian.PutUint16(bytes[1:], uint16(i))
		return bytes
	} else if i < 0x100000000 {
		bytes := make([]byte, 5)
		bytes[0] = 0xfe
		binary.LittleEndian.PutUint32(bytes[1:], uint32(i))
		return bytes
	} else if uint64(i) < largeUint64 {
		bytes := make([]byte, 9)
		bytes[0] = 0xff
		binary.LittleEndian.PutUint64(bytes[1:], uint64(i))
		return bytes
	} else {
		panic("Integer is too large")
	}
}

// ViToInt decodes a variable-length integer representation from a byte slice.
// It takes a byte slice 'byteint' as input, which is expected to contain a variable-length
// encoded integer. The function parses this encoding and returns the decoded integer value
// and the number of bytes consumed from the 'byteint' slice.
//
// The encoding format is as follows:
// - If the first byte is less than 253, it represents the integer directly in 1 byte.
// - If the first byte is 253, the integer is represented in 2 bytes (16 bits).
// - If the first byte is 254, the integer is represented in 4 bytes (32 bits).
// - If the first byte is 255, the integer is represented in 8 bytes (64 bits).
//
// The decoded integer value is returned along with the number of bytes consumed
// from the input slice to represent that integer.
func ViToInt(byteint []byte) (int, int) {
	ni := int(byteint[0])
	size := 0

	if ni < 253 {
		return ni, 1
	}

	switch ni {
	case 253:
		size = 2
	case 254:
		size = 4
	default:
		size = 8
	}

	value := int(binary.LittleEndian.Uint64(byteint[1 : 1+size]))

	return value, size + 1
}

// PrependVarint encodes the length of a byte slice 'data' as a variable-length integer
// and prepends it to the original data. The resulting byte slice contains the encoded
// length followed by the original data.
//
// The function takes a byte slice 'data' as input, calculates the variable-length encoding
// of its length using the EncodeVarint function, and then prepends this encoded length to
// the 'data' slice. The resulting byte slice is returned.
func PrependVarint(data []byte) []byte {
	// Encode the length of the data as a variable-length integer
	n := EncodeVarint(len(data))

	// Prepend the encoded length to the original data
	result := append(n, data...)

	return result
}

// PackInt32LE packs a 32-bit integer 'value' into a little-endian byte slice.
// It takes an integer 'value' as input and creates a new 4-byte byte slice where
// the integer is packed in little-endian byte order using the binary.LittleEndian
// package. The resulting byte slice represents the 'value' in little-endian format
// and is returned.
func PackInt32LE(value int) []byte {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, uint32(value))
	return buf
}

// PackUint32LE packs a 32-bit unsigned integer 'value' into a little-endian byte slice.
// It takes a 32-bit unsigned integer 'value' as input and creates a new 4-byte byte slice
// where the integer is packed in little-endian byte order using the binary.LittleEndian
// package. The resulting byte slice represents the 'value' in little-endian format and is returned.
func PackUint32LE(value uint32) []byte {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, value)
	return buf
}

// PackBigIntToLittleEndian packs a big.Int 'value' into an 8-byte little-endian byte slice.
// It takes a big.Int 'value' as input and creates an 8-byte byte slice where the integer
// is packed in little-endian byte order. The resulting byte slice represents the 'value' in
// little-endian format and is returned.
func PackBigIntToLittleEndian(value *big.Int) []byte {
	amount := new(big.Int).Set(value)
	buffer := make([]byte, 8)
	for i := 0; i < 8; i++ {
		buffer[i] = byte(amount.Uint64() & 0xff)
		amount.Rsh(amount, 8)
	}
	return buffer
}

// IntFromBytes converts a byte slice 'bytes' into an integer using the specified byte order.
// It takes a byte slice 'bytes' and a byte order ('endian') as input and converts the byte
// slice into an integer, respecting the provided byte order.
//
// The function supports byte slices of lengths 1, 2, and 4 bytes and handles them accordingly:
// - For 1-byte slices, it returns the integer as an int8.
// - For 2-byte slices, it returns the integer as an int16 using the provided byte order.
// - For 4-byte slices, it returns the integer as an int32 using the provided byte order.
//
// If the input byte slice is empty, the function raises a panic.
// If the input byte slice length is not supported (other than 1, 2, or 4 bytes), it raises a panic.
func IntFromBytes(bytes []byte, endian binary.ByteOrder) int {
	if len(bytes) == 0 {
		panic("Input bytes should not be empty")
	}

	switch len(bytes) {
	case 1:
		return int(int8(bytes[0]))
	case 2:
		return int(endian.Uint16(bytes))
	case 4:
		return int(endian.Uint32(bytes))
	default:
		panic("unsupported byte lengt")
	}
}
