package ecc

import (
	"bitcoin/digest"
	"bitcoin/formating"
	"crypto/elliptic"
	"crypto/sha256"
	"errors"
	"math/big"
)

// SignMessage signs the given message using the provided private key
// and returns the digital signature as a byte slice.
func SingMessage(message []byte, privateKey []byte) []byte {
	curve := P256k1()
	n := curve.Params().N
	e := calculateE(n, message)
	pD := decodeBigInt(privateKey)
	kCalculator := NewRFC6979KCalculator(sha256.New(), n, pD, message, nil)
	var r *big.Int
	var s *big.Int
	for {
		var k *big.Int

		for {
			k = kCalculator.NextK()
			x, _ := ScalarBaseMultBigInt(curve, *k) // Assuming you have an appropriate function for scalar multiplication.
			r = new(big.Int).Mod(x, n)

			if r.Cmp(big.NewInt(0)) != 0 {
				break
			}
		}
		d := pD
		kinv := new(big.Int).ModInverse(k, n)
		dr := new(big.Int).Mul(d, r)
		edr := new(big.Int).Add(e, dr)
		s = new(big.Int).Mod(new(big.Int).Mul(kinv, edr), n)

		if s.Cmp(big.NewInt(0)) != 0 {
			break
		}
	}
	buffer := make([]byte, 64)
	copy(buffer[:32], encodeBigInt(r))

	nDiv2 := new(big.Int).Rsh(n, 1)
	// Compare sig.s to nDiv2
	compareResult := s.Cmp(nDiv2)

	// Set s based on the comparison result
	if compareResult > 0 {
		// s = n - sig.s
		s = new(big.Int).Sub(n, s)
	}
	copy(buffer[32:], encodeBigInt(s))
	return buffer

}

// SignDer signs the given message using the provided private key and
// entry points, returning the DER-encoded digital signature as a byte slice.
func SingDer(message []byte, privateKey []byte, entryPointes []byte) []byte {
	curve := P256k1()
	n := curve.Params().N
	e := calculateE(n, message)
	pD := decodeBigInt(privateKey)
	kCalculator := NewRFC6979KCalculator(sha256.New(), n, pD, message, entryPointes)
	var r *big.Int
	var s *big.Int
	for {
		var k *big.Int

		for {
			k = kCalculator.NextK()
			x, _ := ScalarBaseMultBigInt(curve, *k) // Assuming you have an appropriate function for scalar multiplication.
			r = new(big.Int).Mod(x, n)

			if r.Cmp(big.NewInt(0)) != 0 {
				break
			}
		}
		d := pD
		kinv := new(big.Int).ModInverse(k, n)
		dr := new(big.Int).Mul(d, r)
		edr := new(big.Int).Add(e, dr)
		s = new(big.Int).Mod(new(big.Int).Mul(kinv, edr), n)

		if s.Cmp(big.NewInt(0)) != 0 {
			break
		}
	}
	return ListBigIntToDER([]*big.Int{r, s})

}

// SignInput signs the given transaction digest using the provided private key,
// applying the specified signature hash type, and returns the resulting signature
// as a hexadecimal string.
func SingInput(privateKey []byte, message []byte, sigHash int) string {
	signature := SingDer(message, privateKey, nil)
	attempt := 1
	lengthR := int(signature[3])
	for lengthR == 33 {
		attemptBytes := formating.Bytes32FromInt(attempt)
		signature = SingDer(message, privateKey, attemptBytes)
		attempt++
		lengthR = int(signature[3])

		if attempt > 50 {
			panic("wrong !!!!! sign must implanet")
		}
	}
	derPrefix := signature[0]
	lengthTotal := signature[1]
	derTypeInt := signature[2]
	lengthR = int(signature[3])
	R := signature[4 : 4+lengthR]
	lengthS := int(signature[5+lengthR])
	S := signature[5+lengthR+1:]
	sAsBigint := formating.BytesToInt(S)

	var newS []byte

	if lengthS == 33 {
		newSAsBigint := new(big.Int).Sub(P256k1().Params().N, sAsBigint)
		newS = encodeBigInt(newSAsBigint)
		lengthS -= 1
		lengthTotal -= 1
	} else {
		newS = S
	}
	newSignature := append([]byte{derPrefix, byte(lengthTotal), byte(derTypeInt), byte(lengthR)}, R...)
	newSignature = append(newSignature, byte(derTypeInt), byte(lengthS))
	newSignature = append(newSignature, newS...)
	newSignature = append(newSignature, byte(sigHash))
	return formating.BytesToHex(newSignature)
}

// VerifySchnorr verifies a Schnorr signature using the provided message, public key, and signature.
// It returns true if the signature is valid, and false otherwise.
func VerifySchnorr(message []byte, publicKey []byte, signature []byte) bool {
	curve := P256k1()
	if len(message) != 32 || len(publicKey) != 32 || len(signature) != 64 {
		return false
	}
	px, py, err := liftX(decodeBigInt(publicKey))
	if err != nil {
		return false
	}

	r := decodeBigInt(signature[:32])
	s := decodeBigInt(signature[32:])
	if r.Cmp(curve.Params().P) >= 0 || s.Cmp(curve.Params().N) >= 0 {
		return false
	}
	messageCombine := append(append(formating.CopyBytes(signature[:32]), publicKey...), message...)
	messageHash := digest.TaggedHash(messageCombine, "BIP0340/challenge")
	e := new(big.Int).Mod(decodeBigInt(messageHash), curve.Params().N)
	spX, spY := curve.ScalarMult(curve.Params().Gx, curve.Params().Gy, signature[32:])
	sub := new(big.Int).Sub(curve.Params().N, e)
	ePX, EPY := curve.ScalarMult(px, py, encodeBigInt(sub))

	rX, Ry := curve.Add(spX, spY, ePX, EPY)
	if Ry.Bit(0) == 1 || rX.Cmp(r) != 0 {
		return false
	}
	return true
}

// SchnorrSign generates a Schnorr signature for the given message using the secret key and auxiliary data.
// It returns the Schnorr signature as a byte slice.
func SchnorrSign(message []byte, secret []byte, aux []byte) []byte {
	if len(message) != 32 {
		panic("The message must be a 32-byte array.")
	}
	curve := P256k1()
	secretPoint := decodeBigInt(secret)
	one := big.NewInt(1)
	if !(one.Cmp(secretPoint) <= 0 && secretPoint.Cmp(new(big.Int).Sub(curve.Params().N, one)) <= 0) {
		panic("The secret key must be an integer in the range 1..n-1.")
	}
	if len(aux) != 32 {
		panic("aux_rand must be 32 bytes")

	}
	pX, pY := curve.ScalarMult(curve.Params().Gx, curve.Params().Gy, secret)

	if pY.Bit(0) == 1 {
		secretPoint = new(big.Int).Sub(curve.Params().N, secretPoint)
	}
	t := formating.XorBytes(encodeBigInt(secretPoint), digest.TaggedHash(aux, "BIP0340/aux"))

	combined := append(append(t, encodeBigInt(pX)...), message...)
	kHash := digest.TaggedHash(combined, "BIP0340/nonce")

	k0 := new(big.Int).Mod(decodeBigInt(kHash), curve.Params().N)
	if k0.Cmp(big.NewInt(0)) == 0 {
		panic("Failure. This happens only with negligible probability.")
	}
	rX, rY := curve.ScalarMult(curve.Params().Gx, curve.Params().Gy, encodeBigInt(k0))

	if rY.Bit(0) == 1 {
		k0 = new(big.Int).Sub(curve.Params().N, k0)
	}
	combinedMessageHash := append(append(encodeBigInt(rX), encodeBigInt(pX)...), message...)
	messageHash := digest.TaggedHash(combinedMessageHash, "BIP0340/challenge")
	e := new(big.Int).Mod(decodeBigInt(messageHash), curve.Params().N)
	temp := new(big.Int)
	temp.Mul(e, secretPoint)
	temp.Add(k0, temp)
	eKey := new(big.Int).Mod(temp, curve.Params().N)
	signature := append(encodeBigInt(rX), encodeBigInt(eKey)...)
	verify := VerifySchnorr(message, encodeBigInt(pX), signature)
	if !verify {
		panic("The created signature does not pass verification.")
	}
	return signature

}

func liftX(x *big.Int) (*big.Int, *big.Int, error) {
	curve := P256k1()
	prime := curve.Params().P
	temp := new(big.Int)
	temp.Exp(x, big.NewInt(3), prime)
	ySq := new(big.Int).Add(temp, big.NewInt(7))
	ySq.Mod(ySq, prime)
	// Calculate y = ySq^((prime+1)/4) % prime
	exponent := new(big.Int).Div(new(big.Int).Add(prime, big.NewInt(1)), big.NewInt(4))
	y := new(big.Int).Exp(ySq, exponent, prime)

	// Check if y^2 % prime == ySq, return null if not
	ySquared := new(big.Int).Exp(y, big.NewInt(2), prime)
	zero := big.NewInt(0)
	if ySquared.Cmp(ySq) != 0 {
		return zero, zero, errors.New("")
	}

	// Calculate result based on the least significant bit of y

	if y.Bit(0) == zero.Bit(0) {
		return x, y, nil
	} else {
		result := new(big.Int).Sub(prime, y)
		return x, result, nil
	}
}

func RecoverPublicKey(recId int, sig []byte, message []byte) []byte {
	r := decodeBigInt(sig[:32])
	s := decodeBigInt(sig[32:])
	i := new(big.Int).Div(big.NewInt(int64(recId)), big.NewInt(2))
	curve := P256k1()
	x := new(big.Int).Set(r)
	x.Mul(i, curve.Params().N)
	x.Add(x, r)
	if x.Cmp(curve.Params().P) >= 0 {
		return nil
	}
	Rx, Ry := decompressKey(x, (recId&1) == 1)
	if Rx == nil || Ry == nil {
		return nil
	}
	Px, Py := curve.ScalarMult(Rx, Ry, encodeBigInt(curve.Params().N))
	if Px == nil || Py == nil {
		return nil
	}
	e := decodeBigInt(message)
	eInv := new(big.Int).Neg(e)
	eInv.Mod(eInv, curve.Params().N)

	rInv := new(big.Int).ModInverse(r, curve.Params().N)

	srInv := new(big.Int).Mul(rInv, s)
	srInv.Mod(srInv, curve.Params().N)

	eInvrInv := new(big.Int).Mul(rInv, eInv)
	eInvrInv.Mod(eInvrInv, curve.Params().N)

	preQX, preQY := curve.ScalarMult(curve.Params().Gx, curve.Params().Gy, encodeBigInt(eInvrInv))
	if preQX == nil || preQY == nil {
		return nil
	}
	sX, sY := curve.ScalarMult(Rx, Ry, encodeBigInt(srInv))
	qX, qY := curve.Add(preQX, preQY, sX, sY)
	encode := elliptic.Marshal(curve, qX, qY)
	return encode

}

func x9IntegerToBytes(s *big.Int, qLength int) []byte {
	// Convert the big integer to a byte slice of qLength bytes
	bytes := s.Bytes()
	if len(bytes) > qLength {
		return bytes[len(bytes)-qLength:]
	} else if len(bytes) < qLength {
		tmp := make([]byte, qLength)
		copy(tmp[qLength-len(bytes):], bytes)
		return tmp
	}
	return bytes
}

func decompressKey(xBN *big.Int, yBit bool) (*big.Int, *big.Int) {
	// Define the elliptic curve parameters for secp256k1
	curve := P256k1()

	// Convert xBN to bytes with the required length
	qLength := ((curve.Params().BitSize + 7) / 8) + 1

	compEnc := x9IntegerToBytes(xBN, qLength)

	// Set the first byte based on the yBit (odd/even) and decode the point

	if yBit {
		compEnc[0] = 0x03 // Set the LSB for odd y
	} else {
		compEnc[0] = 0x02
	}
	return UnmarshalCompressed(curve, compEnc)
}
