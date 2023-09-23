package uuid

import (
	"bitcoin/formating"
	"crypto/rand"
	"fmt"
	"strings"
)

// GenerateUUIDv4 generates a version 4 UUID.
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
func GenerateUUIDv4Bytes() ([]byte, error) {
	random, err := GenerateUUIDv4()
	if err != nil {
		return nil, err
	}
	toBytes := ToBuffer(random)
	return toBytes, nil
}
func ToBuffer(uuidString string) []byte {
	// Remove dashes and convert the hexadecimal string to bytes
	cleanUUIDString := removeDashes(uuidString)
	bytes := formating.HexToBytes(cleanUUIDString)
	return bytes
}

func removeDashes(uuidString string) string {
	return strings.ReplaceAll(uuidString, "-", "")
}
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
