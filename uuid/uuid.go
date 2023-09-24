package uuid

import (
	"bitcoin/formating"
	"crypto/rand"
	"fmt"
	"strings"
)

// GenerateUUIDv4 generates a version 4 UUID.
// It creates a random UUID following the version 4 format,
// where certain bits are set to specific values to indicate
// the version (4) and variant (2). The generated UUID is then
// formatted as a string with hyphens to match the typical UUID format.
// The resulting UUID string is returned along with a nil error.
// If there's an error during random byte generation, it raises a panic.
func GenerateUUIDv4() (string, error) {
	uuid := make([]byte, 16)
	_, err := rand.Read(uuid)
	if err != nil {
		panic("cannot generate random")
	}

	// Set the version (4) and variant (2) bits
	uuid[6] = (uuid[6] & 0x0f) | 0x40
	uuid[8] = (uuid[8] & 0x3f) | 0x80

	// Format the UUID as a string
	uuidString := fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
	return uuidString, nil
}

// GenerateUUIDv4Bytes generates a version 4 UUID as a byte slice.
// It first generates a version 4 UUID as a string using the GenerateUUIDv4 function.
// If there's an error during UUID generation, it is propagated as an error.
// Next, it converts the generated UUID string into a byte slice using the ToBuffer function.
// The resulting byte slice represents the UUID in binary form and is returned along with a nil error.
func GenerateUUIDv4Bytes() ([]byte, error) {
	random, err := GenerateUUIDv4()
	if err != nil {
		return nil, err
	}
	toBytes := ToBuffer(random)
	return toBytes, nil
}

// ToBuffer converts a hexadecimal UUID string to a byte slice.
// It takes a hexadecimal UUID string 'uuidString' as input and performs the following steps:
// 1. Removes dashes from the input string to create a clean hexadecimal representation.
// 2. Converts the clean hexadecimal string to a byte slice using the HexToBytes function.
// The resulting byte slice represents the UUID in binary form and is returned.
func ToBuffer(uuidString string) []byte {
	// Remove dashes and convert the hexadecimal string to bytes
	cleanUUIDString := removeDashes(uuidString)
	bytes := formating.HexToBytes(cleanUUIDString)
	return bytes
}

// removeDashes removes dashes from a string.
// It takes a string 'uuidString' as input and returns a new string
// where all dashes ('-') have been removed.
func removeDashes(uuidString string) string {
	return strings.ReplaceAll(uuidString, "-", "")
}

// FromBuffer converts a 16-byte buffer into a version 4 UUID string.
// It takes a 16-byte buffer 'buffer' as input, which is expected to represent
// a UUID in binary form. The function validates the input buffer's length to ensure
// it's 16 bytes, which is required for UUIDv4.
//
// If the buffer length is not 16 bytes, an error is returned.
//
// If the buffer length is valid, the function converts the binary UUID into a string
// by formatting its hexadecimal representation. Dashes are inserted at appropriate
// positions to create the standard UUIDv4 string format.
//
// The resulting UUIDv4 string is returned along with a nil error.
func FromBuffer(buffer []byte) (string, error) {
	if len(buffer) != 16 {
		return "", fmt.Errorf("invalid buffer length. UUIDv4 buffers must be 16 bytes long")
	}

	hexBytes := make([]string, 16)
	for _, byteVal := range buffer {
		hexBytes = append(hexBytes, fmt.Sprintf("%02x", byteVal))
	}

	// Join the hexadecimal values with dashes to form a UUIDv4 string
	uuidString := strings.Join(hexBytes, "")

	// Insert dashes at appropriate positions
	return fmt.Sprintf(
		"%s-%s-%s-%s-%s",
		uuidString[:8],
		uuidString[8:12],
		uuidString[12:16],
		uuidString[16:20],
		uuidString[20:],
	), nil
}
