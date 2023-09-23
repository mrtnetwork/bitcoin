// Reference implementation for Bech32/Bech32m and segwit addresses.
package bech32

import (
	"errors"
	"strings"
)

const charset = "qpzry9x8gf2tvdw0s3jn54khce6mua7l"

type Bech32Type int

const (
	Bech32 Bech32Type = iota + 1
	Bech32M
)

var generator = []int{0x3b6a57b2, 0x26508e6d, 0x1ea119fa, 0x3d4233dd, 0x2a1462b3}

// Internal function that computes the Bech32 checksum.
func bech32Polymod(values []byte) int {
	chk := 1
	for _, value := range values {
		top := chk >> 25
		chk = (chk&0x1ffffff)<<5 ^ int(value)
		for i := 0; i < 5; i++ {
			if (top>>i)&1 != 0 {
				chk ^= generator[i]
			}
		}
	}
	return chk
}

// Expand the HRP into values for checksum computation.
func bech32HrpExpand(hrp string) []byte {
	values := make([]byte, 0, len(hrp)*2+1)
	for i := 0; i < len(hrp); i++ {
		values = append(values, byte(hrp[i])>>5)
	}
	values = append(values, 0)
	for i := 0; i < len(hrp); i++ {
		values = append(values, byte(hrp[i])&31)
	}
	return values
}

// Verify a checksum given HRP and converted data characters.
func bech32VerifyChecksum(hrp string, data []byte) Bech32Type {
	combined := append(bech32HrpExpand(hrp), data...)
	c := bech32Polymod(combined)

	if c == int(Bech32) {
		return Bech32
	}
	if c == int(0x2bc830a3) {
		return Bech32M
	}
	return 0
}

// bech type to spec value
func getSpecValue(spec Bech32Type) int {
	if spec == Bech32M {
		return 0x2bc830a3
	}
	return 1
}

// Compute the checksum values given HRP and data."
func bech32CreateChecksum(hrp string, data []byte, spec Bech32Type) []byte {
	specValue := getSpecValue(spec)

	values := append(bech32HrpExpand(hrp), data...)
	polymod := bech32Polymod(append(values, 0, 0, 0, 0, 0, 0)) ^ int(specValue)
	checksum := make([]byte, 6)
	for i := 0; i < 6; i++ {
		checksum[i] = byte((polymod >> uint(5*(5-i))) & 31)
	}
	return checksum
}

// Compute a Bech32 string given HRP and data values.
func bech32Encode(hrp string, data []byte, spec Bech32Type) string {
	combined := append(data, bech32CreateChecksum(hrp, data, spec)...)
	encoded := hrp + "1"
	for _, d := range combined {
		encoded += string(charset[d])
	}
	return encoded
}

// Validate a Bech32/Bech32m string, and determine HRP and data.
func bech32Decode(bech string) (string, []byte, Bech32Type) {
	if len(bech) < 1 || len(bech) > 90 {
		return "", nil, 0
	}
	bech = toLowerCase(bech)
	pos := strings.LastIndex(bech, "1")
	if pos < 1 || pos+7 > len(bech) {
		return "", nil, 0
	}
	hrp := bech[:pos]
	data := make([]byte, len(bech)-pos-1)
	for i := pos + 1; i < len(bech); i++ {
		idx := strings.IndexRune(charset, rune(bech[i]))
		if idx == -1 {
			return "", nil, 0
		}
		data[i-pos-1] = byte(idx)
	}
	spec := bech32VerifyChecksum(hrp, data)
	if spec == 0 {
		return "", nil, 0
	}
	return hrp, data[:len(data)-6], spec
}

func toLowerCase(str string) string {
	var lower string
	for _, char := range str {
		if char >= 'A' && char <= 'Z' {
			lower += string(char + 32)
		} else {
			lower += string(char)
		}
	}
	return lower
}

// General power-of-2 base conversion.
func convertBits(data []byte, fromBits, toBits int, pad bool) []byte {
	acc := 0
	bits := 0
	var ret []byte
	maxv := (1 << toBits) - 1
	maxAcc := (1 << (fromBits + toBits - 1)) - 1

	for _, value := range data {
		acc = ((acc << fromBits) | int(value)) & maxAcc
		bits += fromBits
		for bits >= toBits {
			bits -= toBits
			ret = append(ret, byte((acc>>bits)&maxv))
		}
	}

	if pad {
		if bits > 0 {
			ret = append(ret, byte((acc<<(toBits-bits))&maxv))
		}
	} else if bits >= fromBits || ((acc<<(toBits-bits))&maxv) > 0 {
		return nil
	}

	return ret
}

// Decode a segwit address.
func DecodeBech32(address string) (int, []byte, error) {
	_, data, spec := bech32Decode(address)
	if data == nil {
		return 0, nil, errors.New("failed to decode Bech32")
	}

	bits := convertBits(data[1:], 5, 8, false)
	if bits == nil || len(bits) < 2 || len(bits) > 40 {
		return 0, nil, errors.New("invalid bits")
	}

	if data[0] > 16 {
		return 0, nil, errors.New("invalid data[0]")
	}

	if data[0] == 0 && len(bits) != 20 && len(bits) != 32 {
		return 0, nil, errors.New("invalid bits length")
	}

	if (data[0] == 0 && spec != Bech32) || (data[0] != 0 && spec != Bech32M) {
		return 0, nil, errors.New("invalid spec")
	}

	return int(data[0]), bits, nil
}

// Encode a segwit address.
func EncodeBech32(hrp string, version int, data []byte) (string, error) {
	var bech32Type Bech32Type
	if version == 0 {
		bech32Type = Bech32
	} else {
		bech32Type = Bech32M
	}

	bits := convertBits(data, 8, 5, true)
	if bits == nil {
		return "", errors.New("failed to convertBits")
	}

	combinedData := append([]byte{byte(version)}, bits...)
	encoded := bech32Encode(hrp, combinedData, bech32Type)
	if _, _, err := DecodeBech32(encoded); err != nil {
		return "", errors.New("failed to decodeBech32")
	}

	return encoded, nil
}
