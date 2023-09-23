package address

import (
	"bitcoin/base58"
	"bitcoin/digest"
	"bitcoin/ecc"
	"bitcoin/formating"
	"bitcoin/scripts"
	"bitcoin/tools"
	"errors"
)

// instantiates an object from a hash160 hex string
func fromHash160(hash160 string) (string, error) {
	if !tools.IsValidHash160(hash160) {
		return "", errors.New("invalid hash160")
	}
	return hash160, nil
}

// instantiates an object from address string encoding
func fromAddress(address string) (string, error) {
	if !tools.IsValidAddress(address) {
		return "", errors.New("invalid addres")
	}
	toHash160, err := addressToHash160(address)
	if err != nil {
		return "", err
	}
	return toHash160, nil
}

// instantiates an object from a redeem_script
func fromScript(script scripts.Script) string {
	toHash160 := scriptToHash160(script)
	return toHash160
}

func addressToHash160(address string) (string, error) {
	dec, err := base58.Decode(address)
	if err != nil {
		return "", errors.New("invalid addresss")
	}
	hash160 := dec[1 : len(dec)-4]
	return formating.BytesToHex(hash160), nil
}

func scriptToHash160(script scripts.Script) string {
	scriptBytes := script.ToBytes()
	toHash160 := digest.Hash160(scriptBytes)
	return formating.BytesToHex(toHash160)
}

// Creates an address object from a hash160 string
func P2SHAddressFromHash160(hash160 string) (P2SHAdress, error) {
	bip, err := fromHash160(hash160)
	if err != nil {
		return P2SHAdress{LegacyAddress{"", P2PKHInP2SH}}, err
	}
	return P2SHAdress{LegacyAddress{bip, P2PKHInP2SH}}, nil
}

// Creates an address object from an address string
// args [P2SHAddressTypeParam or AddressType ] default
/*
Choosing the type of address for transaction output is not important,
but for transaction input and spending from this address,
an error will be encountered if the transaction is entered incorrectly.
*/
func P2SHAddressFromAddress(address string, args ...interface{}) P2SHAdress {
	t := getP2shAddressParam(args...)
	bip, err := fromAddress(address)
	if err != nil {
		panic(err)
	}
	return P2SHAdress{LegacyAddress{bip, t}}
}

// Creates an address object from a Script object
func P2SHAddressFromScript(script scripts.Script, addressType AddressType) P2SHAdress {
	bip := fromScript(script)
	return P2SHAdress{LegacyAddress{bip, addressType}}
}

// Creates an address object from a hash160 string
func P2PKHAddressFromHash160(hash160 string) P2PKHAddress {
	bip, err := fromHash160(hash160)
	if err != nil {
		panic(err.Error())
	}
	return P2PKHAddress{LegacyAddress{bip, P2PKH}}
}

// Creates an address object from an address string
func P2PKHAddressFromAddress(address string) P2PKHAddress {
	bip, err := fromAddress(address)
	if err != nil {
		return P2PKHAddress{LegacyAddress{"", P2PKH}}
	}
	return P2PKHAddress{LegacyAddress{bip, P2PKH}}
}

// Creates an address object from a Script object
func P2PKHAddressFromScript(script scripts.Script) P2PKHAddress {
	bip := fromScript(script)
	return P2PKHAddress{LegacyAddress{bip, P2PKH}}
}

// Creates an address object from publick key
func P2PKAddressFromPublicKey(public string) P2PKAddress {
	toPublicBytes := formating.HexToBytes(public)
	if !ecc.IsPoint(toPublicBytes) {
		panic("invalid public key")
	}
	return P2PKAddress{LegacyAddress{Hash160: public, Type: P2PK}}

}

/*
returns the address's string encoding (Bech32)
You can use NetworkParams or NetworkInfo (TestnetNetwork, BitcoinNetwork) in arguments
to select the network, otherwise, the default network is used.
*/
func (s LegacyAddress) toAddress(network ...interface{}) string {
	networkType := getNetworkParams(network...)
	var tobytes []byte

	tobytes = formating.HexToBytes(s.Hash160)
	switch s.Type {
	case P2PKH:
		tobytes = append([]byte{byte(networkType.p2PKHPrefix)}, tobytes...)
	case P2PK:
		toHash160 := digest.Hash160(tobytes)
		tobytes = append([]byte{byte(networkType.p2PKHPrefix)}, toHash160...)
	default:
		tobytes = append([]byte{byte(networkType.p2SHPrefix)}, tobytes...)
	}
	hash := digest.DoubleHash(tobytes)
	addrBytes := append(tobytes, hash[:4]...)
	return base58.Encode(addrBytes)
}
