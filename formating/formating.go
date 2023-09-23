package formating

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"
	"reflect"
	"strings"
)

const largeUint64 = ^uint64(0)

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
func strip0x(hex string) string {
	if strings.HasPrefix(hex, "0x") {
		return hex[2:]
	}
	return hex
}

func HexToBytes(hexStr string) []byte {
	bytes, _ := hex.DecodeString(strip0x(hexStr))
	return bytes
}
func HexToBytesCatch(hexStr string) ([]byte, error) {
	bytes, err := hex.DecodeString(strip0x(hexStr))
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
func BytesToHex(hexStr []byte) string {
	toHex := hex.EncodeToString(hexStr)
	return toHex
}
func PrependVarint(data []byte) []byte {
	n := EncodeVarint(len(data))
	result := append(n, data...)
	return result
}

func GenerateRandom(size int) ([]byte, error) {
	randomBytes := make([]byte, size)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}
	return randomBytes, nil
}
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
func PackInt32LE(value int) []byte {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, uint32(value))
	return buf
}
func PackUint32LE(value uint32) []byte {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, value)
	return buf
}
func ReverseBytes(data []byte) []byte {
	length := len(data)
	reversed := make([]byte, length)

	for i, j := 0, length-1; i < length; i, j = i+1, j-1 {
		reversed[i] = data[j]
	}

	return reversed
}
func PackBigIntToLittleEndian(value *big.Int) []byte {
	amount := new(big.Int).Set(value)
	buffer := make([]byte, 8)
	for i := 0; i < 8; i++ {
		buffer[i] = byte(amount.Uint64() & 0xff)
		amount.Rsh(amount, 8)
	}
	return buffer
}

func CopyBytes(original []byte) []byte {
	newBytes := make([]byte, len(original))

	copy(newBytes, original)
	return newBytes
}
func FlattenList(input interface{}) []interface{} {
	var result []interface{}

	// Use reflection to check the type of input
	switch reflect.TypeOf(input).Kind() {
	case reflect.Slice:
		slice := reflect.ValueOf(input)
		for i := 0; i < slice.Len(); i++ {
			item := slice.Index(i).Interface()
			result = append(result, FlattenList(item)...)
		}
	default:
		result = append(result, input)
	}

	return result
}
func PrintBytes(prefix string, byteSlice []byte) {
	fmt.Printf("%v [", prefix)
	for i, b := range byteSlice {
		fmt.Printf("%d", int(b))
		if i < len(byteSlice)-1 {
			fmt.Print(", ")
		}
	}
	fmt.Println("]")
}
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

func BytesToBinary(data []byte) string {
	binaries := make([]string, len(data))
	for i, b := range data {
		binaries[i] = fmt.Sprintf("%08b", b)
	}
	return strings.Join(binaries, "")
}

func BinaryToByte(binary string) byte {
	value := 0
	for i, c := range binary {
		if c == '1' {
			value += 1 << uint(7-i)
		}
	}
	return byte(value)
}
func FillRange(slice []byte, value byte) {
	for i := range slice {
		slice[i] = value
	}
}
func ToInterfaceSlice(slice interface{}) []interface{} {
	v := reflect.ValueOf(slice)
	if v.Kind() != reflect.Slice {
		panic("Input is not a slice")
	}

	interfaceSlice := make([]interface{}, v.Len())
	for i := 0; i < v.Len(); i++ {
		interfaceSlice[i] = v.Index(i).Interface()
	}

	return interfaceSlice
}
