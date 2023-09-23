package keypair

import (
	"bitcoin/address"
	"bitcoin/base58"
	"bitcoin/constant"
	"bitcoin/digest"
	"bitcoin/ecc"
	"bitcoin/formating"
	"bytes"
	"encoding/hex"
	"fmt"
	"strings"
)

type ECPrivate struct {
	privateKey []byte
	publicKey  []byte
}

// NewECPrivate creates a new ECPrivate instance from a hexadecimal private key string.
// It returns the ECPrivate object or an error if parsing fails.
func NewECPrivate(privateHex string) (*ECPrivate, error) {
	privBytes, err := hex.DecodeString(privateHex)
	if err != nil {
		return nil, err
	}
	if !ecc.IsValidBitcoinPrivateKey(privBytes) {
		return nil, fmt.Errorf("invalid private key")
	}

	public, _ := ecc.GenerateBitcoinPublicKey(privBytes)
	return &ECPrivate{
		privateKey: privBytes,
		publicKey:  public,
	}, nil
}

// return private key as bytes
func (ecPriv *ECPrivate) ToBytes() []byte {
	return formating.CopyBytes(ecPriv.privateKey)
}

// return private key as hexadecimal string
func (ecPriv *ECPrivate) ToHex() string {
	return formating.BytesToHex(ecPriv.privateKey)
}

// return public key as hexadecimal string
func (ecPriv *ECPrivate) TopublicHex() string {
	return formating.BytesToHex(ecPriv.publicKey)
}

// return public key as bytes
func (ecPriv *ECPrivate) ToPublic() []byte {
	return formating.CopyBytes(ecPriv.publicKey)
}

// NewECPrivate creates a new ECPrivate instance from a private key bytes.
func NewECPrivateFromBytes(privBytes []byte) (*ECPrivate, error) {
	if !ecc.IsValidBitcoinPrivateKey(privBytes) {
		return nil, fmt.Errorf("invalid private key")
	}

	public, err := ecc.GenerateBitcoinPublicKey(privBytes)
	if err != nil {
		return nil, fmt.Errorf("invalid private key")
	}
	return &ECPrivate{
		privateKey: formating.CopyBytes(privBytes),
		publicKey:  public,
	}, nil
}

// NewECPrivateFromWIF creates an ECPrivate instance from a WIF string
// and returns a pointer to the initialized object.
func NewECPrivateFromWIF(wif string) (*ECPrivate, error) {
	b64, err := base58.Decode(wif)
	if err != nil {
		return nil, fmt.Errorf("invalid WIF length")
	}
	lengthWithoutChecksum := len(b64) - 4

	// Extract keyBytes
	keyBytes := b64[:lengthWithoutChecksum]

	// Extract checksum
	checksum := b64[lengthWithoutChecksum:]

	h := digest.DoubleHash(keyBytes)

	if !bytes.Equal(h[:4], checksum) {
		panic("invalid checksum")
	}
	keyBytes = keyBytes[1:]
	if len(keyBytes) > 32 {
		keyBytes = keyBytes[:len(keyBytes)-1]
	}

	return NewECPrivateFromBytes(keyBytes)
}

// ToWIF converts an ECPrivate key to its Wallet Import Format (WIF) representation.
func (ecPriv *ECPrivate) ToWIF(compressed bool, networkType address.NetworkInfo) string {
	var bytes []byte
	if compressed {
		bytes = append([]byte{networkType.WIF()}, ecPriv.privateKey...)
		bytes = append(bytes, 0x01)
	} else {
		bytes = append([]byte{networkType.WIF()}, ecPriv.privateKey...)
	}

	checksum := digest.DoubleHash(bytes)
	bytes = append(bytes, checksum[:4]...)
	return base58.Encode(bytes)
}

func (ecPriv *ECPrivate) GetPublic() *ECPublic {
	pub, e := NewECPPublicFromBytes(ecPriv.ToPublic())
	if e != nil {
		panic("invalid public key")
	}
	return pub
}

// SignTaprootTransaction signs a transaction digest using the ECPrivate key,
// specified signature hash type, and additional scripts. If 'tweak' is true,
// the private key may be tweaked before signing. Returns the resulting signature
// as a string.
func (ecPriv *ECPrivate) SignTaprootTransaction(txDigest []byte, sigHash int, scripts []interface{}, tweak bool) string {
	var keyBytes []byte
	if tweak {
		pub := ecPriv.GetPublic()
		tw, _ := pub.CalculateTweek(scripts)
		keyBytes = ecc.TweakTaprootPrivate(ecPriv.ToBytes(), tw)

	} else {
		keyBytes = ecPriv.ToBytes()
	}

	auxBytes := append(txDigest, keyBytes...)
	auxHash := digest.SingleHash(auxBytes)
	signature := ecc.SchnorrSign(txDigest, keyBytes, auxHash)
	if sigHash != constant.TAPROOT_SIGHASH_ALL {
		signature = append(signature, byte(sigHash))
	}
	return formating.BytesToHex(signature)
}

// sign transaction digest
func (ecPriv *ECPrivate) SingInput(txDigest []byte, sigHash ...interface{}) string {
	sig := constant.SIGHASH_ALL
	for _, opt := range sigHash {
		switch v := opt.(type) {
		case int:
			sig = v
		default:
			panic("invalid Tx Input argruments")
		}
	}
	return ecc.SingInput(ecPriv.ToBytes(), txDigest, sig)
}
func magicPrefix(message string) []byte {
	prefix := "\x18Bitcoin Signed Message:\n"
	size := formating.EncodeVarint(len(message))
	bytes := []byte(message)
	result := append([]byte(prefix), size...)
	result = append(result, bytes...)
	return result
}

func MagicMessage(message string) []byte {
	magic := magicPrefix(message)
	return digest.SingleHash(magic)
}

// signs the message's digest and returns the signature
func (ecPriv *ECPrivate) SignMessage(message string, compressed bool) string {
	m := digest.SingleHash(MagicMessage(message))

	signature := ecc.SingMessage(m, ecPriv.ToBytes())
	prefix := 27

	// Determine recid based on the prefix
	if compressed {
		prefix += 4
	}
	addr := ecPriv.GetPublic().ToAddress(compressed)
	for i := prefix; i < prefix+4; i++ {
		// Attempt to create a new byte slice containing the signature
		sig := make([]byte, 0)
		char := string(rune(i))
		charBytes := []byte(char)

		sig = append(sig, charBytes...)
		sig = append(sig, signature...)
		// Handle any potential errors (e.g., utf8 encoding errors)
		if pub := GetSignaturePublic(message, sig); pub != nil {
			if strings.EqualFold(pub.ToAddress(compressed).Program().Hash160, addr.Program().Hash160) {
				return formating.BytesToHex(sig)
			}
			continue
		}

	}
	panic("cannot validate message")
}
