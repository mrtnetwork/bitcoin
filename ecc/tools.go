package ecc

import (
	"crypto/elliptic"
	"fmt"
	"github.com/mrtnetwork/bitcoin/formating"
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
	r := formating.PadByteSliceTo32(encodeBigInt(x))
	s := formating.PadByteSliceTo32(encodeBigInt(y))
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

// negateSecretKey negates a Bitcoin private key 'secret' and returns the resulting
// private key as a big.Int. The negation operation involves subtracting the private key
// from the curve order if the public key point derived from the private key has its
// least significant bit set to 1.
//
// The function first generates the corresponding Bitcoin public key from 'secret' and
// ensures it's uncompressed. It then extracts the x-coordinate of the public key as a big.Int.
//
// If the least significant bit of the x-coordinate is 1, indicating an odd value, the
// function performs the negation by subtracting 'secret' from the curve order and returns
// the negated private key as a big.Int.
//
// If the least significant bit is not 1, indicating an even value, the function returns
// 'secret' as is, since negation is not needed.
func negateSecretKey(secret []byte) *big.Int {
	curve := P256k1()

	// Generate the corresponding Bitcoin public key from the private key
	pub, _ := GenerateBitcoinPublicKey(secret)

	// Ensure the public key is uncompressed and extract the x-coordinate as a big.Int
	reEncode := ReEncodedFromForm(pub, false)
	point := decodeBigInt(reEncode[32:])

	nestedPoint := decodeBigInt(secret)

	// Check if the least significant bit of the x-coordinate is 1
	if point.Bit(0) == 1 {
		// Perform negation by subtracting 'secret' from the curve order
		kexExpend := decodeBigInt(secret)
		nestedPoint = new(big.Int).Sub(curve.Params().N, kexExpend)
	}

	return nestedPoint
}

// isOrderScalar checks if a byte slice 'secret' represents a scalar that is less than
// the elliptic curve group order. It is used to validate that the 'secret' is within
// the valid scalar range for the curve.
//
// The 'secret' input must be exactly 32 bytes in length. If it is not, the function
// returns false since it cannot be a valid scalar.
//
// The function compares the 'secret' with the elliptic curve group order and returns
// true if 'secret' is less than the group order; otherwise, it returns false.
func isOrderScalar(secret []byte) bool {
	if len(secret) != 32 {
		return false
	}

	var EC_GROUP_ORDER = formating.HexToBytes("fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141")

	return compare(secret, EC_GROUP_ORDER) < 0
}

// PointAddScalar computes a new elliptic curve point by adding a scalar 'tweak'
// to a given public key point 'public'. The resulting point is returned as a byte slice.
// The function also takes a 'compress' boolean parameter that determines whether the
// resulting point should be compressed or not.
//
// The 'public' input must represent a valid elliptic curve point, and the 'tweak' must be
// a valid scalar within the curve's order. If either the 'public' point or 'tweak' is invalid,
// the function returns an error.
//
// If 'tweak' is the zero scalar, the function returns the original 'public' point, optionally
// compressed as specified.
//
// If 'tweak' is not the zero scalar, the function computes the new point as follows:
// - Adds 'tweak' to the generator point (G) and obtains point (qX, qY).
// - Adds 'public' and (qX, qY) to obtain point (uX, uY).
// - Checks if (uX, uY) is a valid elliptic curve point.
//
// If 'compress' is true, the resulting point is compressed; otherwise, it's not compressed.
// The compressed or uncompressed point is returned as a byte slice.
//
// If any point computation results in an invalid point, the function returns an error.
func PointAddScalar(public []byte, tweak []byte, compress bool) ([]byte, error) {
	// Check if the 'public' point is a valid elliptic curve point
	if !IsPoint(public) {
		return nil, fmt.Errorf("bad Point")
	}

	// Check if 'tweak' is a valid scalar within the curve's order
	if !isOrderScalar(tweak) {
		return nil, fmt.Errorf("bad Tweak")
	}

	curve := P256k1()

	// Decode the 'public' point into (x, y) coordinates
	x, y := UnCompressedPoint(public)

	var ZERO32 = make([]uint8, 32)

	// If 'tweak' is the zero scalar, return the original 'public' point
	if compare(tweak, ZERO32) == 0 {
		if compress {
			return MarshalCompressed(curve, x, y), nil
		}
		return elliptic.Marshal(curve, x, y), nil
	}

	// Compute point (qX, qY) as (Gx, Gy) + 'tweak' * G
	qX, qY := curve.ScalarMult(curve.Params().Gx, curve.Params().Gy, tweak)

	// Compute point (uX, uY) as 'public' + (qX, qY)
	uX, uY := curve.Add(x, y, qX, qY)

	// Check if (uX, uY) is a valid elliptic curve point (not infinity)
	isInfinity := curve.IsOnCurve(uX, uY) && uX.Cmp(big.NewInt(0)) == 0 && uY.Cmp(big.NewInt(0)) == 0
	if isInfinity {
		return nil, fmt.Errorf("invalid point")
	}

	// Return the resulting point, optionally compressed as specified
	if compress {
		return MarshalCompressed(curve, uX, uY), nil
	}
	return elliptic.Marshal(curve, uX, uY), nil
}

// GenerateTweek generates a tweaked private key using a given point and tweak.
// It takes a byte slice 'point' representing a Bitcoin private key and a byte slice 'tweak'
// representing the tweak value. The function validates the input point and tweak for their
// correctness. It checks if the point is a valid Bitcoin private key and if the tweak is a
// valid scalar within the curve's order.
//
// If either the point or the tweak is invalid, it returns an error.
//
// If both the point and tweak are valid, the function computes the tweaked private key
// by adding the tweak to the point (in modular arithmetic), ensuring that the result stays
// within the curve's order. The resulting tweaked private key is returned as a byte slice.
//
// If the computed tweaked private key is invalid, it returns an error.
func GenerateTweek(point []byte, tweak []byte) ([]byte, error) {
	curve := P256k1()
	if !IsValidBitcoinPrivateKey(point) {
		return nil, fmt.Errorf("bad Point")
	}
	if !isOrderScalar(tweak) {
		return nil, fmt.Errorf("bad Tweek")
	}
	d := decodeBigInt(point)
	t := decodeBigInt(tweak)

	dt := new(big.Int).Set(d)
	dt.Add(dt, t)
	dt.Mod(dt, curve.Params().N)
	p := encodeBigInt(dt)
	if !IsValidBitcoinPrivateKey(p) {
		return nil, fmt.Errorf("bad Private key")
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
