// / Base58 and base58Check implement
package base58

import (
	"bitcoin/digest"
	"bytes"
	"errors"
	"fmt"
	"math/big"
)

const (

	//In Bitcoin's Base58 encoding, the character set usually consists of the following 58 characters
	btcAlphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
)

var (
	bs58 = NewBase58(btcAlphabet)
)

type Base58 struct {
	alphabet    string
	alphabetMap map[byte]int
	base        *big.Int
	leader      string
}

func NewBase58(alphabet string) *Base58 {
	base58 := &Base58{
		alphabet:    alphabet,
		alphabetMap: make(map[byte]int),
		base:        new(big.Int).SetInt64(int64(len(alphabet))),
		leader:      string(alphabet[0]),
	}

	for i := 0; i < len(alphabet); i++ {
		base58.alphabetMap[alphabet[i]] = i
	}

	return base58
}

// The checksum is appended to the end of the binary data
func EncodeCheck(data []byte) string {
	hash := digest.DoubleHash(data)
	combined := append(data, hash[:4]...)

	return Encode(combined)
}

// convert binary data into a Base58-encoded string
func Encode(data []byte) string {
	if len(data) == 0 {
		return ""
	}

	digits := []int{0}

	for i := 0; i < len(data); i++ {
		carry := int(data[i])
		for j := 0; j < len(digits); j++ {
			carry += digits[j] << 8
			digits[j] = carry % int(bs58.base.Int64())
			carry = carry / int(bs58.base.Int64())
		}
		for carry > 0 {
			digits = append(digits, carry%int(bs58.base.Int64()))
			carry = carry / int(bs58.base.Int64())
		}
	}
	result := ""
	// Deal with leading zeros
	for k := 0; data[k] == 0 && k < len(data)-1; k++ {
		result += bs58.leader
	}

	for q := len(digits) - 1; q >= 0; q-- {
		result += string(bs58.alphabet[digits[q]])
	}
	return result
}

// convert a Base58-encoded string back into its original binary data or byte array
func Decode(s string) ([]byte, error) {
	if len(s) == 0 {
		return nil, errors.New("empty input")
	}

	bytes := []int{0}
	for i := 0; i < len(s); i++ {
		value, ok := bs58.alphabetMap[s[i]]
		if !ok {
			return nil, errors.New("Non-base58 character")
		}

		carry := value
		for j := 0; j < len(bytes); j++ {
			carry += bytes[j] * int(bs58.base.Int64())
			bytes[j] = carry & 0xff
			carry >>= 8
		}

		for carry > 0 {
			bytes = append(bytes, carry&0xff)
			carry >>= 8
		}
	}

	// Deal with leading zeros
	for k := 0; string(s[k]) == bs58.leader && k < len(s)-1; k++ {
		bytes = append(bytes, []int{0}...)
	}
	reversedBytes := make([]byte, len(bytes))
	for i, j := 0, len(bytes)-1; i <= j; i, j = i+1, j-1 {
		reversedBytes[i], reversedBytes[j] = byte(bytes[j]), byte(bytes[i])
	}

	return reversedBytes, nil
}

// decode and check for the presence of a checksum within the decoded data
func DecodeCheck(s string) ([]byte, error) {
	// Decode the Base58 string
	decoded, err := Decode(s)
	if err != nil {
		return nil, err
	}

	// Check if the decoded data has at least 4 bytes (for checksum)
	if len(decoded) < 4 {
		return nil, fmt.Errorf("invalid base58check: insufficient length")
	}

	// Split the decoded data into payload and checksum
	payload := decoded[:len(decoded)-4]
	checksum := decoded[len(decoded)-4:]

	// Calculate the checksum of the payload
	newChecksum := digest.DoubleHash(payload)[:4]

	// Compare the calculated checksum with the provided checksum
	if !bytes.Equal(checksum, newChecksum) {
		return nil, fmt.Errorf("invalid checksum")
	}

	return payload, nil
}
