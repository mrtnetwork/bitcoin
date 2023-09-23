package address

import (
	"bitcoin/formating"
	"fmt"
	"strconv"
)

type Network int

const (
	Mainnet Network = iota
	Testnet
)

// struct of network parameters in function arguments
type NetworkParams struct {
	Network *networkInfo
}

// struct of network parameters in function arguments
type P2SHAddressTypeParam struct {
	AddressType *AddressType
}

// NetworkInfo represents information about a Bitcoin network.
type networkInfo struct {
	bech32        string
	p2PKHPrefix   int
	p2SHPrefix    int
	wIF           int
	extendPrivate map[AddressType]string
	extendPublic  map[AddressType]string
	network       Network
	name          string
}

type NetworkInfo interface {
	Bech32() string
	P2PKHPrefix() int
	P2SHPrefix() int
	WIF() byte
	ExtendPrivate(AddressType) string
	ExtendPublic(AddressType) string
	IsMainNet() bool
	Network() Network
}

// BitcoinNetwork represents the Bitcoin network information.
var MainnetNetwork = networkInfo{
	network:     Mainnet,
	name:        "mainnet",
	bech32:      "bc",
	p2PKHPrefix: 0x00,
	p2SHPrefix:  0x05,
	wIF:         0x80,
	extendPrivate: map[AddressType]string{
		P2PKH:        "0x0488ade4",
		P2PKInP2SH:   "0x0488ade4",
		P2PKHInP2SH:  "0x0488ade4",
		P2WPKH:       "0x04b2430c",
		P2WPKHInP2SH: "0x049d7878",
		P2WSH:        "0x02aa7a99",
		P2WSHInP2SH:  "0x0295b005",
	},
	extendPublic: map[AddressType]string{
		P2PKH:        "0x0488b21e",
		P2PKInP2SH:   "0x0488b21e",
		P2PKHInP2SH:  "0x0488b21e",
		P2WPKH:       "0x04b24746",
		P2WPKHInP2SH: "0x049d7cb2",
		P2WSH:        "0x02aa7ed3",
		P2WSHInP2SH:  "0x0295b43f",
	},
}

// TestnetNetwork represents the Bitcoin testnet network information.
var TestnetNetwork = networkInfo{
	bech32:      "tb",
	p2PKHPrefix: 0x6f,
	name:        "testnet",
	p2SHPrefix:  0xc4,
	wIF:         0xef,
	extendPrivate: map[AddressType]string{
		P2PKH:        "0x04358394",
		P2PKInP2SH:   "0x04358394",
		P2PKHInP2SH:  "0x04358394",
		P2WPKH:       "0x045f18bc",
		P2WPKHInP2SH: "0x044a4e28",
		P2WSH:        "0x02575048",
		P2WSHInP2SH:  "0x024285b5",
	},
	extendPublic: map[AddressType]string{
		P2PKH:        "0x043587cf",
		P2PKInP2SH:   "0x043587cf",
		P2PKHInP2SH:  "0x043587cf",
		P2WPKH:       "0x045f1cf6",
		P2WPKHInP2SH: "0x044a5262",
		P2WSH:        "0x02575483",
		P2WSHInP2SH:  "0x024289ef",
	},
	network: Testnet,
}

// NetworkFromWIF returns the Bitcoin network information based on a WIF (Wallet Import Format) string.
func NetworkFromWIF(wif string) (networkInfo, error) {
	w, err := strconv.ParseInt(wif, 16, 64)
	if err != nil {
		return networkInfo{}, err
	}

	if MainnetNetwork.wIF == int(w) {
		return MainnetNetwork, nil
	} else if TestnetNetwork.wIF == int(w) {
		return TestnetNetwork, nil
	}

	return networkInfo{}, fmt.Errorf("WIF prefix not supported, only Bitcoin or Testnet accepted")
}

// NetworkFromXPrivePrefix returns the Bitcoin address type based on an extended private key prefix.
func NetworkFromXPrivePrefix(prefix []byte) AddressType {
	w := "0x" + formating.BytesToHex(prefix)
	for key, value := range TestnetNetwork.extendPrivate {
		if value == w {
			return key
		}
	}

	for key, value := range MainnetNetwork.extendPrivate {
		if value == w {
			return key
		}
	}

	return AddressType(-1)
}

// NetworkFromXPublicPrefix returns the Bitcoin address type based on an extended public key prefix.
func NetworkFromXPublicPrefix(prefix []byte) AddressType {
	w := "0x" + formating.BytesToHex(prefix)
	for key, value := range TestnetNetwork.extendPublic {
		if value == w {
			return key
		}
	}

	for key, value := range MainnetNetwork.extendPublic {
		if value == w {
			return key
		}
	}
	return AddressType(-1)
}

func (n *networkInfo) ExtendPublic(addressType AddressType) string {
	return n.extendPublic[addressType]
}
func (n *networkInfo) ExtendPrivate(addressType AddressType) string {
	return n.extendPrivate[addressType]
}

// access to wif version of network
func (n *networkInfo) WIF() byte {
	return byte(n.wIF)
}
func (n *networkInfo) Bech32() string {
	return n.bech32
}
func (n *networkInfo) Network() Network {
	return n.network
}
func (n *networkInfo) P2PKHPrefix() int {
	return n.p2PKHPrefix
}

func (n *networkInfo) P2SHPrefix() int {
	return n.p2SHPrefix
}
func (n *networkInfo) IsMainNet() bool {
	return n.network == Mainnet
}

var defaultNetwork = MainnetNetwork

// get default network of application
func DefaultNetwork() *networkInfo {
	return &defaultNetwork
}

// update the default network, the methods
// that require a network will use the default
// network if the desired parameter is not found.
func SetDefaultNetwork(network networkInfo) {
	defaultNetwork = network
}

// Find the NetworkInfo type in the function parameters
func getNetworkParams(useDefault bool, args ...interface{}) *networkInfo {
	argruments := formating.FlattenList(args)
	var currentNetwork *networkInfo
	for _, opt := range argruments {
		switch v := opt.(type) {
		case NetworkParams:
			if v.Network != nil {
				currentNetwork = v.Network
				break
			}
		case *NetworkParams:
			if v.Network != nil {
				currentNetwork = v.Network
				break
			}
		case *networkInfo:
			{
				currentNetwork = v
			}
		case networkInfo:
			{
				currentNetwork = &v
			}
		}

	}
	if currentNetwork == nil && useDefault {
		currentNetwork = &defaultNetwork
	}
	return currentNetwork
}
