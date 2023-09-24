package hdwallet

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/mrtnetwork/bitcoin/address"
	"github.com/mrtnetwork/bitcoin/base58"
	"github.com/mrtnetwork/bitcoin/bip39"
	"github.com/mrtnetwork/bitcoin/digest"
	"github.com/mrtnetwork/bitcoin/ecc"
	"github.com/mrtnetwork/bitcoin/formating"
	"github.com/mrtnetwork/bitcoin/keypair"
	"regexp"
	"strconv"
	"strings"
)

// HdWallet represents an HD Wallet.
type HdWallet struct {
	// Represents the depth level in the hierarchical deterministic (HD) wallet structure
	depth int
	// Denotes the index of this wallet within its parent node in the HD wallet hierarchy
	index int
	// Contains the fingerprint or unique identifier of the parent wallet from which this wallet was derived
	fingerPrint []byte
	// Indicates whether this HD wallet is the root wallet in the hierarchy
	isRoot bool
	// Stores the private key associated with this HD wallet. (Note: Make sure it's appropriately secured.)
	privateKey *keypair.ECPrivate
	// Stores the public key corresponding to the private key of this wallet.
	publicKey *keypair.ECPublic
	// A flag that specifies whether this wallet was derived from an extended public key (xpub) rather than a private key
	fromXpub bool
	// Represents the chain code associated with this wallet, which is used in HD wallet key derivation
	chainCode []byte
}

const highBit = 0x80000000
const maxUint31 = 2147483647
const maxUint32 = 4294967295

// NewHdWalletFromPrivateKey creates an HdWallet instance from a private key.
func newHdWalletFromPrivateKey(privateKey *keypair.ECPrivate, chainCode []byte, depth, index int, fingerPrint []byte) *HdWallet {
	if len(fingerPrint) == 0 {
		fingerPrint = make([]byte, 4)
	}
	return &HdWallet{
		depth:       depth,
		index:       index,
		fingerPrint: fingerPrint,
		isRoot:      bytes.Equal(fingerPrint, make([]byte, 4)),
		privateKey:  privateKey,
		fromXpub:    false,
		chainCode:   chainCode,
		publicKey:   privateKey.GetPublic(),
	}
}

// NewHdWalletFromPublicKey creates an HdWallet instance from a public key.
func newHdWalletFromPublicKey(publicKey *keypair.ECPublic, chainCode []byte, depth, index int, fingerPrint []byte) *HdWallet {
	if len(fingerPrint) == 0 {
		fingerPrint = make([]byte, 4)
	}
	return &HdWallet{

		depth:       depth,
		index:       index,
		fingerPrint: fingerPrint,
		isRoot:      bytes.Equal(fingerPrint, make([]byte, 4)),
		publicKey:   publicKey,
		fromXpub:    true,
		chainCode:   chainCode,
	}
}

// FromMnemonic creates an HdWallet by generating its components from a given mnemonic phrase
// and an optional passphrase. It returns a pointer to the HdWallet and an error if there's any issue.
func FromMnemonic(mnemonic string, passphrase string) (*HdWallet, error) {
	// Perform the necessary operations to create an HdWallet from a mnemonic and passphrase
	seed := bip39.ToSeed(mnemonic, passphrase)
	if len(seed) < 16 {
		return nil, fmt.Errorf("seed should be at least 128 bits")
	}
	if len(seed) > 64 {
		return nil, fmt.Errorf("seed should be at most 512 bits")
	}

	hash := digest.HmacSHA512([]byte("Bitcoin seed"), seed)
	private, _ := keypair.NewECPrivateFromBytes(hash[:32])
	chainCode := hash[32:]

	wallet := newHdWalletFromPrivateKey(private, chainCode, 0, 0, nil)
	return wallet, nil
}

// AddDrive adds a new child HdWallet derived from the current wallet.
func (hd *HdWallet) addDrive(index int) (*HdWallet, error) {
	if uint32(index) > maxUint32 || index < 0 {
		return nil, fmt.Errorf("expected UInt32")
	}
	isHardened := uint32(index) >= highBit
	data := make([]byte, 37)

	if isHardened {
		if hd.fromXpub {
			return nil, fmt.Errorf("cannot use hardened path in public wallet")
		}
		data[0] = 0x00
		copy(data[1:], hd.privateKey.ToBytes())
		byteData := data[33:]
		byteData[0] = byte(index >> 24)
		byteData[1] = byte(index >> 16)
		byteData[2] = byte(index >> 8)
		byteData[3] = byte(index)
	} else {
		copy(data, hd.publicKey.ToCompressedBytes())
		byteData := data[33:]
		byteData[0] = byte(index >> 24)
		byteData[1] = byte(index >> 16)
		byteData[2] = byte(index >> 8)
		byteData[3] = byte(index)
	}

	masterKey := digest.HmacSHA512(hd.chainCode, data)
	key := masterKey[:32]
	chain := masterKey[32:]

	if !ecc.IsValidBitcoinPrivateKey(key) {
		return hd.addDrive(index + 1)
	}

	childDeph := hd.depth + 1
	childIndex := index

	if hd.fromXpub {
		newPoint, err := ecc.PointAddScalar(hd.publicKey.ToUnCompressedBytes(true), key, true)
		if err != nil {
			return hd.addDrive(index + 1)
		}
		finger := digest.Hash160(hd.publicKey.ToCompressedBytes())[:4]
		newPublicKey, e := keypair.NewECPPublicFromBytes(newPoint)
		if e != nil {
			return nil, e
		}
		return newHdWalletFromPublicKey(newPublicKey, chain, childDeph, childIndex, finger), nil
	}

	newPrivate, err := ecc.GenerateTweek(hd.privateKey.ToBytes(), key)

	if err != nil {
		return hd.addDrive(index + 1)
	}
	finger := digest.Hash160(hd.publicKey.ToCompressedBytes())[:4]
	newPrivateKey, e := keypair.NewECPrivateFromBytes(newPrivate)
	if e != nil {
		return nil, e
	}
	return newHdWalletFromPrivateKey(newPrivateKey, chain, childDeph, childIndex, finger), nil
}

// IsValidPath checks if a BIP32 path is valid.
func IsValidPath(path string) bool {
	pattern := `^(m\/)?(\d+'?\/)*\d+'?$`

	// Compile the regular expression pattern
	regex := regexp.MustCompile(pattern)

	// Check if the path matches the pattern
	return regex.MatchString(path)
}

// DrivePath derives an HdWallet instance from a BIP32 path.
func DrivePath(masterWallet *HdWallet, path string) (*HdWallet, error) {
	if !IsValidPath(path) {
		return nil, fmt.Errorf("invalid BIP32 Path")
	}

	splitPath := strings.Split(path, "/")
	if splitPath[0] == "m" || splitPath[0] == "M" {
		splitPath = splitPath[1:]
	}

	for _, indexStr := range splitPath {
		var index int
		var err error

		if strings.HasSuffix(indexStr, "'") {
			if masterWallet.fromXpub {
				return nil, fmt.Errorf("hardened derivation path is invalid for xpublic key")
			}
			indexStr = strings.TrimSuffix(indexStr, "'")
			index, err = strconv.Atoi(indexStr)

			if err != nil || index > int(maxUint31) || index < 0 {
				return nil, fmt.Errorf("wrong index")
			}
			newDrive, err := masterWallet.addDrive(index + highBit)
			if err != nil {
				return nil, err
			}
			masterWallet = newDrive
		} else {
			index, err = strconv.Atoi(indexStr)
			if err != nil {
				return nil, err
			}
			newDrive, err := masterWallet.addDrive(index)
			if err != nil {
				return nil, err
			}
			masterWallet = newDrive
		}
	}

	return masterWallet, nil
}

// isRootKey checks whether the provided extended private or public key is a root key
// for a given network. It returns a boolean indicating whether it's a root key and
// the fingerprint (unique identifier) associated with it.
func isRootKey(xPrivateKey string, network address.NetworkInfo, isPublicKey bool) (bool, []byte) {
	// Decode the Base58Check-encoded extended private key
	dec, err := base58.DecodeCheck(xPrivateKey)
	if err != nil {
		panic(err)
	}

	// Check if the decoded data has the expected length
	if len(dec) != 78 {
		panic("Invalid xPrivateKey")
	}

	// Extract the first 4 bytes (semantic) from the decoded data
	semantic := dec[:4]

	// Determine the version based on whether it's a public or private key
	var version address.AddressType
	if isPublicKey {
		// Use NetworkInfo.networkFromXPublicPrefix(semantic) for version lookup
		// Replace with your actual implementation
		version = address.NetworkFromXPublicPrefix(semantic)
	} else {
		// Use NetworkInfo.networkFromXPrivatePrefix(semantic) for version lookup
		// Replace with your actual implementation
		version = address.NetworkFromXPrivePrefix(semantic)
	}

	// Check if the version is valid
	if version == -1 {
		panic("Invalid network")
	}

	// Determine the expected network prefix based on whether it's a public or private key
	var networkPrefix string
	if isPublicKey {
		networkPrefix = network.ExtendPublic(version)
	} else {
		networkPrefix = network.ExtendPrivate(version)
	}

	// Convert the network prefix to bytes and append zeros
	prefix := formating.HexToBytes(networkPrefix + "000000000000000000")

	// Check if the prefix matches the first bytes of the decoded data
	return bytes.Equal(prefix, dec[:len(prefix)]), dec
}

// GetPrivate returns the private key associated with the HdWallet.
func (hd *HdWallet) GetPrivate() (*keypair.ECPrivate, error) {
	if hd.fromXpub {
		return nil, fmt.Errorf("cannot access private key from public wallet")
	}
	return hd.privateKey, nil
}

// GetPublic returns the public key associated with the HdWallet.
func (hd *HdWallet) GetPublic() *keypair.ECPublic {
	return hd.publicKey
}

// decodeXKeys decodes an extended key (xKey) into its constituent parts. It returns a slice of byte slices
// representing the different components of the extended key, based on whether it's a public or private key.
func decodeXKeys(xKey []byte, isPublic bool) [][]byte {
	var parts [][]byte

	parts = append(parts, xKey[0:4])
	parts = append(parts, xKey[4:5])
	parts = append(parts, xKey[5:9])
	parts = append(parts, xKey[9:13])
	parts = append(parts, xKey[13:45])

	if isPublic {
		parts = append(parts, xKey[45:])
	} else {
		parts = append(parts, xKey[46:])
	}

	return parts
}

// ToXPublicKey generates an extended public key (xpub) string representation of the HdWallet
// using the specified semantic and network information.
func (hd *HdWallet) ToXPublicKey(semantic address.AddressType, network address.NetworkInfo) string {
	version := network.ExtendPublic(semantic)
	// Convert depth to a byte slice
	depthBytes := []byte{byte(hd.depth)}

	// Convert index to a 4-byte big-endian byte slice
	indexBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(indexBytes, uint32(hd.index))

	// Concatenate all the required byte slices
	result := bytes.Join([][]byte{
		formating.HexToBytes(version),
		depthBytes,
		formating.CopyBytes(hd.fingerPrint),
		indexBytes,
		formating.CopyBytes(hd.chainCode),
		hd.GetPublic().ToCompressedBytes(),
	}, nil)

	// Encode the result using base58 encoding with a checksum
	check := base58.EncodeCheck(result)

	return check
}

// ToXPrivateKey generates an extended private key (xpriv) string representation of the HdWallet
// using the specified semantic and network information.
func (hd *HdWallet) ToXPrivateKey(semantic address.AddressType, network address.NetworkInfo) string {
	if hd.fromXpub {
		panic("connot access private from publicKey wallet")
	}
	version := formating.HexToBytes(network.ExtendPrivate(semantic))
	// Convert depth to a byte slice
	depthBytes := []byte{byte(hd.depth)}
	//0488ade4
	// Convert index to a 4-byte big-endian byte slice
	indexBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(indexBytes, uint32(hd.index))
	privateBytes := append([]byte{0}, hd.privateKey.ToBytes()...)
	// Concatenate all the required byte slices
	result := bytes.Join([][]byte{
		version,
		depthBytes,
		formating.CopyBytes(hd.fingerPrint),
		indexBytes,
		formating.CopyBytes(hd.chainCode),
		privateBytes,
	}, nil)

	// Encode the result using base58 encoding with a checksum
	check := base58.EncodeCheck(result)

	return check
}

// FromXPrivateKey creates an HdWallet from an extended private key (xpriv) string. The 'forRootKey' parameter
// indicates whether the provided key is a root key. The 'network' parameter specifies the network information.
// It returns a pointer to the HdWallet and an error if there's any issue.
func FromXPrivateKey(xPrivateKey string, forRootKey bool, network address.NetworkInfo) (*HdWallet, error) {
	// Check if it's a root key
	isRootKey, xKeyBytes := isRootKey(xPrivateKey, network, false)
	// Verify if it matches the expected root key status
	if forRootKey && !isRootKey {
		return nil, fmt.Errorf("invalid rootXPrivateKey")
	} else if !forRootKey && isRootKey {
		return nil, fmt.Errorf("invalid xPrivateKey")
	}

	// Decode the xKey bytes
	xKeyParts := decodeXKeys(xKeyBytes, false)
	chainCode := xKeyParts[4]
	privateKey, e := keypair.NewECPrivateFromBytes(xKeyParts[5])
	if e != nil {
		return nil, e
	}
	index := formating.IntFromBytes(xKeyParts[3], binary.BigEndian)
	depth := formating.IntFromBytes(xKeyParts[1], binary.BigEndian)
	fingerprint := xKeyParts[2]
	return newHdWalletFromPrivateKey(
		privateKey, chainCode, depth, int(index), fingerprint,
	), nil
}

// FromXPublicKey creates an HdWallet from an extended public key (xpub) string. The 'forceRootKey' parameter
// indicates whether to treat the provided key as a root key. The 'network' parameter specifies the network information.
// It returns a pointer to the HdWallet and an error if there's any issue.
func FromXPublicKey(xPublicKey string, forceRootKey bool, network address.NetworkInfo) (*HdWallet, error) {
	// Check if it's a root key (public key)
	isRootKey, xKeyBytes := isRootKey(xPublicKey, network, true)

	// Verify if it matches the expected root key status
	if forceRootKey && !isRootKey {
		return nil, fmt.Errorf("invalid rootPublicKey")
	} else if !forceRootKey && isRootKey {
		return nil, fmt.Errorf("invalid publicKey")
	}

	// Decode the xKey bytes (public key)
	xKeyParts := decodeXKeys(xKeyBytes, true)

	// Extract the necessary parts
	chainCode := xKeyParts[4]
	publicKey, e := keypair.NewECPPublicFromBytes(xKeyParts[5])
	if e != nil {
		return nil, e
	}
	index := formating.IntFromBytes(xKeyParts[3], binary.BigEndian)
	depth := formating.IntFromBytes(xKeyParts[1], binary.BigEndian)
	fingerprint := xKeyParts[2]

	return newHdWalletFromPublicKey(
		publicKey,
		chainCode,
		depth,
		int(index),
		fingerprint,
	), nil
}
