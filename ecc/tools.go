package ecc

import (
	"bitcoin/formating"
	"crypto/elliptic"
	"errors"
	"math/big"
)

// Tweaks the public key with the specified tweak. Required to create the
// taproot public key from the internal key.
func TweakTaprootPoint(pub []byte, twek []byte) []byte {
	curve := P256k1()
	x := decodeBigInt(pub[:32])
	y := decodeBigInt(pub[32:])
	if y.Bit(0) == 1 {
		y.Sub(curve.Params().P, y)
	}
	qx, qy := curve.ScalarMult(curve.Params().Gx, curve.Params().Gy, twek)
	x, y = curve.Add(x, y, qx, qy)

	if y.Bit(0) == 1 { // Check if y is odd
		y = y.Sub(curve.Params().P, y)
	}
	r := padByteSliceTo32(encodeBigInt(x))
	s := padByteSliceTo32(encodeBigInt(y))
	combined := append(r, s...)
	return combined
}

// Tweaks the private key before signing with it. Check if public key's y
// is even and negate the private key before tweaking if it is not.
func TweakTaprootPrivate(secret []byte, tweak []byte) []byte {
	curve := P256k1()
	pub, _ := GenerateBitcoinPublicKey(secret)
	reEncode := ReEncodedFromForm(pub, false)
	point := decodeBigInt(reEncode[32:])
	tweakPoint := decodeBigInt(tweak)
	negatedKey := decodeBigInt(secret)
	if point.Bit(0) == 1 { // Check if y is odd
		negatedKey = negateSecretKey(secret)
	}
	var tw big.Int
	tw.Add(negatedKey, tweakPoint).Mod(&tw, curve.Params().N)
	return encodeBigInt(&tw)
}

func negateSecretKey(secret []byte) *big.Int {
	curve := P256k1()
	pub, _ := GenerateBitcoinPublicKey(secret)
	reEncode := ReEncodedFromForm(pub, false)
	point := decodeBigInt(reEncode[32:])
	nestedPoint := decodeBigInt(secret)
	if point.Bit(0) == 1 { // Check if y is odd
		// y = y.Sub(curve.Params().P, y)
		kexExpend := decodeBigInt(secret)
		nestedPoint = new(big.Int).Sub(curve.Params().N, kexExpend)
	}
	return nestedPoint
}
func xorBytes(a, b []byte) []byte {
	if len(a) != len(b) {
		panic("Input slices must have the same length")
	}

	result := make([]byte, len(a))

	for i := 0; i < len(a); i++ {
		result[i] = a[i] ^ b[i]
	}

	return result
}
func isPointCompressed(p []byte) bool {
	return p[0] != 0x04
}
func isOrderScalar(secret []byte) bool {
	if len(secret) != 32 {
		return false
	}
	var EC_GROUP_ORDER = formating.HexToBytes("fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141")
	return compare(secret, EC_GROUP_ORDER) < 0
}

func PointAddScalar(public []byte, tweak []byte, compress bool) ([]byte, error) {
	if !IsPoint(public) {
		return nil, errors.New("bad Point")
	}
	if !isOrderScalar(tweak) {
		return nil, errors.New("bad Tweek")
	}
	curve := P256k1()
	compressed := isPointCompressed(public)
	x, y := UnCompressedPoint(public)
	var ZERO32 = make([]uint8, 32)
	if compare(tweak, ZERO32) == 0 {
		if compressed {
			return MarshalCompressed(curve, x, y), nil
		}
		return elliptic.Marshal(curve, x, y), nil
	}

	qX, qY := curve.ScalarMult(curve.Params().Gx, curve.Params().Gy, tweak)
	uX, uY := curve.Add(x, y, qX, qY)
	isInfinity := curve.IsOnCurve(uX, uY) && uX.Cmp(big.NewInt(0)) == 0 && uY.Cmp(big.NewInt(0)) == 0
	if isInfinity {
		return nil, errors.New("invalid point")
	}
	if compressed {
		return MarshalCompressed(curve, x, y), nil
	}
	return elliptic.Marshal(curve, x, y), nil
}

func GenerateTweek(point []byte, tweak []byte) ([]byte, error) {
	curve := P256k1()
	if !IsValidBitcoinPrivateKey(point) {
		return nil, errors.New("bad Point")
	}
	if !isOrderScalar(tweak) {
		return nil, errors.New("bad Tweek")
	}
	d := decodeBigInt(point)
	t := decodeBigInt(tweak)

	dt := new(big.Int).Set(d)
	dt.Add(dt, t)
	dt.Mod(dt, curve.Params().N)
	p := encodeBigInt(dt)
	if !IsValidBitcoinPrivateKey(p) {
		return nil, errors.New("bad Private key")
	}
	return p, nil
}

// IsValidBitcoinPrivateKey checks if the given bytes represent a valid Bitcoin private key.
func IsValidBitcoinPrivateKey(privateKeyBytes []byte) bool {
	// Check if the private key is 32 bytes in length
	if len(privateKeyBytes) != 32 {
		return false
	}
	curve := P256k1()
	// Convert the private key bytes to a big integer
	privateKey := new(big.Int).SetBytes(privateKeyBytes)

	// Check if the private key is within the valid range [1, n-1]
	if privateKey.Cmp(big.NewInt(1)) < 0 || privateKey.Cmp(curve.Params().N) >= 0 {
		return false
	}

	// Optionally, check that the private key is not zero
	if privateKey.Cmp(big.NewInt(0)) == 0 {
		return false
	}

	return true
}

// convert bytes to bigint, then compare
func compare(a, b []uint8) int {
	aa := decodeBigInt(a)
	bb := decodeBigInt(b)

	return aa.Cmp(bb)
}

// IsPoint checks if the input bytes represent a valid Bitcoin public key point.
func IsPoint(p []uint8) bool {
	if len(p) < 33 {
		return false
	}

	t := p[0]
	x := p[1:33]
	var ZERO32 = make([]uint8, 32)
	var EC_P = formating.HexToBytes("fffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f")
	if compare(x, ZERO32) == 0 {
		return false
	}
	if compare(x, EC_P) == 1 {
		return false
	}

	xx, yy := UnCompressedPoint(p)
	if xx == nil || yy == nil {
		return false
	}

	if (t == 0x02 || t == 0x03) && len(p) == 33 {
		return true
	}

	y := p[33:]
	if compare(y, ZERO32) == 0 {

		return false
	}
	if compare(y, EC_P) == 1 {
		return false
	}

	if t == 0x04 && len(p) == 65 {
		return true
	}

	return false
}

// decodeBigInt decodes a little-endian byte slice into a BigInt value.
func decodeBigInt(bytes []byte) *big.Int {
	result := new(big.Int)
	for i := 0; i < len(bytes); i++ {
		bigByte := new(big.Int).SetInt64(int64(bytes[len(bytes)-i-1]))
		shift := uint(8 * i)
		bigByte.Lsh(bigByte, shift)
		result.Add(result, bigByte)
	}
	return result
}

// encodeBigInt encodes a BigInt value into a little-endian byte slice.
func encodeBigInt(number *big.Int) []byte {
	var needsPaddingByte int
	var rawSize int

	if number.Cmp(big.NewInt(0)) > 0 {
		rawSize = (number.BitLen() + 7) >> 3
		mask := new(big.Int).Lsh(big.NewInt(1), uint((rawSize-1)*8))
		needsPaddingByte = 0
		if new(big.Int).And(number, mask).Cmp(mask) != 0 {
			needsPaddingByte = 1
		}

		if rawSize < 32 {
			needsPaddingByte = 1
		}
	} else {
		needsPaddingByte = 0
		rawSize = (number.BitLen() + 8) >> 3
	}

	size := rawSize
	if rawSize < 32 {
		size += needsPaddingByte
	}

	result := make([]byte, size)
	for i := 0; i < size; i++ {
		result[size-i-1] = byte(number.Uint64())
		number = new(big.Int).Rsh(number, 8)
	}
	return result
}
func bytes32FromInt(x int) []byte {
	result := make([]byte, 32)
	for i := 0; i < 32; i++ {
		result[32-i-1] = byte((x >> (8 * i)) & 0xFF)
	}
	return result
}

func bytesToInt(data []byte) *big.Int {
	return new(big.Int).SetBytes(data)
}
func padByteSliceTo32(data []byte) []byte {
	if len(data) >= 32 {
		return data[:32]
	}

	paddedData := make([]byte, 32)
	copy(paddedData[32-len(data):], data)
	return paddedData
}

// calculateE computes the integer 'e' from the provided message and modulus 'n'.
// It returns the resulting 'e' as a big integer.
func calculateE(n *big.Int, message []byte) *big.Int {
	log2n := n.BitLen()
	messageBitLength := len(message) * 8

	if log2n >= messageBitLength {
		return new(big.Int).SetBytes(message)
	} else {
		trunc := new(big.Int).SetBytes(message)
		shift := uint(messageBitLength - log2n)
		trunc.Rsh(trunc, shift)
		return trunc
	}
}
