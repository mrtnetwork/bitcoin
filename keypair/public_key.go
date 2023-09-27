// keypair management for private and public keys.
package keypair

import (
	"fmt"
	"strings"

	"github.com/mrtnetwork/bitcoin/address"
	"github.com/mrtnetwork/bitcoin/constant"
	"github.com/mrtnetwork/bitcoin/digest"
	"github.com/mrtnetwork/bitcoin/ecc"
	"github.com/mrtnetwork/bitcoin/formating"
	"github.com/mrtnetwork/bitcoin/scripts"
)

type ECPublic struct {
	publicKey []byte
}

// NewECPPublicFromHex creates a new ECPublic instance from a hex-encoded public key
// and returns a pointer to the initialized object.
func NewECPPublicFromHex(publicHex string) (*ECPublic, error) {
	publicBytes := formating.HexToBytes(publicHex)
	if !ecc.IsPoint(publicBytes) {
		return nil, fmt.Errorf("invalid public key")
	}

	public := ecc.ReEncodedFromForm(publicBytes, false)
	return &ECPublic{
		publicKey: public,
	}, nil
}

// NewECPPublicFromBytes creates a new ECPublic instance from a byte slice
// containing a public key and returns a pointer to the initialized object.
func NewECPPublicFromBytes(publicBytes []byte) (*ECPublic, error) {
	if !ecc.IsPoint(publicBytes) {
		return nil, fmt.Errorf("invalid public key")
	}

	public := ecc.ReEncodedFromForm(publicBytes, false)
	return &ECPublic{
		publicKey: public,
	}, nil
}

// ToHex converts the ECPublic key to a hex-encoded string.
// If 'compressed' is true, the key is in compressed format.
func (ecPublic *ECPublic) ToHex(compressed ...interface{}) string {
	c := getCompressedArgs(compressed...)
	if c {
		p := ecc.ReEncodedFromForm(ecPublic.ToUnCompressedBytes(true), true)
		return formating.BytesToHex(p)
	}
	return formating.BytesToHex(ecPublic.ToUnCompressedBytes(true))
}

func getCompressedArgs(args ...interface{}) bool {
	argruments := formating.FlattenList(args)
	compressed := true
	for _, opt := range argruments {
		switch v := opt.(type) {
		case bool:
			compressed = v
		default:
			panic("invalid compressed argruments")
		}
	}
	return compressed
}

// toHash160 computes the RIPEMD160 hash of the ECPublic key.
// If 'compressed' is true, the key is in compressed format.
func (ecPublic *ECPublic) toHash160(compressed ...interface{}) []byte {
	toPublicHex := ecPublic.ToHex(compressed...)
	toBytes := formating.HexToBytes(toPublicHex)

	return digest.Hash160(toBytes)
}

// toHash160 computes the RIPEMD160 hash of the ECPublic key.
// If 'compressed' is true, the key is in compressed format.
func (ecPublic *ECPublic) ToHash160(compressed ...interface{}) string {
	hash160 := ecPublic.toHash160(compressed...)
	return formating.BytesToHex(hash160)
}

// ToAddress generates a P2PKH (Pay-to-Public-Key-Hash) address from the ECPublic key.
// If 'compressed' is true, the key is in compressed format.
func (ecPublic *ECPublic) ToAddress(compressed ...interface{}) *address.P2PKHAddress {
	h160 := ecPublic.ToHash160(compressed...)
	addr, _ := address.P2PKHAddressFromHash160(h160)
	return addr
}

// ToP2PKHInP2SH generates a P2PKH (Pay-to-Public-Key-Hash) address from the ECPublic key.
// If 'compressed' is true, the key is in compressed format.
func (ecPublic *ECPublic) ToP2PKHInP2SH(compressed ...interface{}) *address.P2SHAdress {
	addr := ecPublic.ToAddress(compressed...)
	p2sh, _ := address.P2SHAddressFromScript(addr.ToScriptPubKey(), address.P2PKHInP2SH)
	return p2sh
}

// ToSegwitAddress generates a P2WPKH (Pay-to-Witness-Public-Key-Hash) SegWit address
// from the ECPublic key. If 'compressed' is true, the key is in compressed format.
func (ecPublic *ECPublic) ToSegwitAddress(compressed ...interface{}) *address.P2WPKHAddresss {
	h160 := ecPublic.ToHash160(compressed...)
	addr, _ := address.P2WPKHAddresssFromProgram(h160)
	return addr
}

// ToP2PKAddress generates a P2PK (Pay-to-Public-Key) address from the ECPublic key.
// If 'compressed' is true, the key is in compressed format.
func (ecPublic *ECPublic) ToP2PKAddress(compressed ...interface{}) *address.P2PKAddress {
	redeem := ecPublic.ToHex(compressed...)
	addr, _ := address.P2PKAddressFromPublicKey(redeem)
	return addr
}

// ToTaprootAddress generates a Taproot address from the ECPublic key
// and an optional script. The 'script' parameter can be used to specify
// custom spending conditions.
func (ecPublic *ECPublic) ToTaprootAddress(script ...interface{}) *address.P2TRAddress {
	taproot, e := ecPublic.ToTapRotHex(script)
	if e != nil {
		panic("invalid taaproot program")
	}
	addr, e := address.P2TRAddressFromProgram(taproot)
	if e != nil {
		panic("invalid taaproot program")
	}
	return addr
}

// ToRedeemScript generates a redeem script from the ECPublic key.
// If 'compressed' is true, the key is in compressed format.
func (ecPublic *ECPublic) ToRedeemScript(compressed ...interface{}) *scripts.Script {
	redeem := ecPublic.ToHex(compressed...)
	return &scripts.Script{Script: []interface{}{redeem, "OP_CHECKSIG"}}
}

// ToP2PKInP2SH generates a P2SH (Pay-to-Script-Hash) address
// wrapping a P2PK (Pay-to-Public-Key) script derived from the ECPublic key.
// If 'compressed' is true, the key is in compressed format.
func (ecPublic *ECPublic) ToP2PKInP2SH(compressed ...interface{}) *address.P2SHAdress {
	p2sh, _ := address.P2SHAddressFromScript(ecPublic.ToRedeemScript(compressed...), address.P2PKInP2SH)
	return p2sh
}

// ToP2WPKHInP2SH generates a P2SH (Pay-to-Script-Hash) address
// wrapping a P2WPKH (Pay-to-Witness-Public-Key-Hash) script derived from the ECPublic key.
// If 'compressed' is true, the key is in compressed format.
func (ecPublic *ECPublic) ToP2WPKHInP2SH(compressed ...interface{}) *address.P2SHAdress {
	addr := ecPublic.ToSegwitAddress(compressed...)
	segwit, _ := address.P2SHAddressFromScript(addr.ToScriptPubKey(), address.P2WPKHInP2SH)
	return segwit
}

// ToP2WSHScript generates a P2WSH (Pay-to-Witness-Script-Hash) script
// derived from the ECPublic key. If 'compressed' is true, the key is in compressed format.
func (ecPublic *ECPublic) ToP2WSHScript(compressed ...interface{}) *scripts.Script {
	return scripts.NewScript("OP_1", ecPublic.ToHex(compressed...), "OP_1", "OP_CHECKMULTISIG")
}

// ToP2WSHAddress generates a P2WSH (Pay-to-Witness-Script-Hash) address
// from the ECPublic key. If 'compressed' is true, the key is in compressed format.
func (ecPublic *ECPublic) ToP2WSHAddress(compressed ...interface{}) *address.P2WSHAddresss {
	script := ecPublic.ToP2WSHScript(compressed...)
	p2wsh, _ := address.P2WSHAddresssFromScript(script)
	return p2wsh
}

// ToP2WSHInP2SH generates a P2SH (Pay-to-Script-Hash) address
// wrapping a P2WSH (Pay-to-Witness-Script-Hash) script derived from the ECPublic key.
// If 'compressed' is true, the key is in compressed format.
func (ecPublic *ECPublic) ToP2WSHInP2SH(compressed ...interface{}) *address.P2SHAdress {
	addr := ecPublic.ToP2WSHAddress(compressed...)
	p2sh, _ := address.P2SHAddressFromScript(addr.ToScriptPubKey(), address.P2WSHInP2SH)
	return p2sh
}

// ToUncompressedBytes returns the uncompressed byte representation of the ECPublic key.
// If 'withPrefix' is true, it includes the prefix byte (0x04) in the output.
func (ecPublic *ECPublic) ToUnCompressedBytes(withPrefix bool) []byte {
	if withPrefix {
		result := append([]byte{0x04}, ecPublic.publicKey...)
		return result
	}
	newBytes := make([]byte, len(ecPublic.publicKey))
	copy(newBytes, ecPublic.publicKey)
	return newBytes

}

// ToCompressedBytes returns the compressed byte representation of the ECPublic key.
func (ecPublic *ECPublic) ToCompressedBytes() []byte {
	comprossedBytes := ecc.ReEncodedFromForm(ecPublic.ToUnCompressedBytes(true), true)
	return comprossedBytes
}

// ToXOnlyHex extracts and returns the x-coordinate (first 32 bytes) of the ECPublic key
// as a hexadecimal string.
func (ecPublic *ECPublic) ToXOnlyHex() string {
	return formating.BytesToHex(ecPublic.publicKey[:32])
}

// ToTapRootHex computes and returns the Taproot commitment point's x-coordinate
// derived from the ECPublic key and an optional script, represented as a hexadecimal string.
func (ecPublic *ECPublic) ToTapRotHex(script []interface{}) (string, error) {
	publicBytes := ecPublic.ToUnCompressedBytes(false)
	tweak, e := ecPublic.CalculateTweek(script)
	if e != nil {
		return "", e
	}
	point := ecc.TweakTaprootPoint(publicBytes, tweak)
	return formating.BytesToHex(point[:32]), nil

}

// tapleafTaggedHash computes and returns the tagged hash of a script for Taproot,
// using the specified script. It prepends a version byte and then tags the hash with "TapLeaf".
func tapleafTaggedHash(script *scripts.Script) []byte {
	scriptBytes := formating.PrependVarint(script.ToBytes())
	part := append([]byte{constant.LEAF_VERSION_TAPSCRIPT}, scriptBytes...)
	return digest.TaggedHash(part, "TapLeaf")
}

// tapBranchTaggedHash computes and returns the tagged hash of two byte slices
// for Taproot, where 'a' and 'b' are the input byte slices. It ensures that 'a' and 'b'
// are sorted and concatenated before tagging the hash with "TapBranch".
func tapBranchTaggedHash(a, b []byte) []byte {
	var part []byte

	if formating.IsLessThanBytes(a, b) {
		part = append(a, b...)
	} else {
		part = append(b, a...)
	}

	return digest.TaggedHash(part, "TapBranch")
}

// getTagHashedMerkleRoot computes and returns the tagged hashed Merkle root for Taproot
// based on the provided argument. It handles different argument types, including scripts
// and lists of scripts.
func getTagHashedMerkleRoot(args interface{}) ([]byte, error) {
	switch val := args.(type) {
	case scripts.Script:
		return tapleafTaggedHash(&val), nil
	case *scripts.Script:
		return tapleafTaggedHash(val), nil
	case []interface{}:
		if len(val) == 0 {
			return nil, nil
		} else if len(val) == 1 {
			return getTagHashedMerkleRoot(val[0])
		} else if len(val) == 2 {
			left, e := getTagHashedMerkleRoot(val[0])
			if e != nil {
				return nil, e
			}
			right, e := getTagHashedMerkleRoot(val[1])
			if e != nil {
				return nil, e
			}
			return tapBranchTaggedHash(left, right), nil
		} else {
			return nil, fmt.Errorf("list cannot have more than 2 branches")

		}
	default:
		return nil, fmt.Errorf("unsupported argument tyxpe")
	}
}

// CalculateTweak computes and returns the TapTweak value based on the ECPublic key
// and an optional script. It uses the key's x-coordinate and the Merkle root of the script
// (if provided) to calculate the tweak.
func (ecPublic *ECPublic) CalculateTweek(script interface{}) ([]byte, error) {

	keyX := formating.CopyBytes(ecPublic.publicKey[:32])

	if script == nil {
		tweek := digest.TaggedHash(keyX, "TapTweak")
		return tweek, nil
	}

	merkleRoot, e := getTagHashedMerkleRoot(script)
	if e != nil {
		return nil, e
	}
	tweek := digest.TaggedHash(append(keyX, merkleRoot...), "TapTweak")
	return tweek, nil
}

// Verify verifies a signature against a message using the ECPublic key.
// It compares the recovered public key from the signature to the provided ECPublic key.
func (ecPublic *ECPublic) Verify(message string, signature string) bool {
	pub := GetSignaturePublic(message, formating.HexToBytes(signature))
	if pub != nil {
		return strings.EqualFold(pub.ToHex(), ecPublic.ToHex())
	}
	return false
}

// GetSignaturePublic extracts and returns the public key associated with a signature
// for the given message. If the extraction is successful, it returns an ECPublic key;
// otherwise, it returns nil.
func GetSignaturePublic(message string, signature []byte) *ECPublic {
	m := digest.SingleHash(MagicMessage(message))
	prefix := int(signature[0])

	var recid int

	// Determine recid based on the prefix
	if prefix >= 31 {
		recid = prefix - 31
	} else {
		recid = prefix - 27
	}

	rec := ecc.RecoverPublicKey(recid, formating.CopyBytes(signature[1:]), m)

	if rec != nil {
		pub, _ := NewECPPublicFromBytes(rec)
		return pub
	}
	return nil
}
