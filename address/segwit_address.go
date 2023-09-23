package address

import (
	"bitcoin/bech32"
	"bitcoin/constant"
	"bitcoin/digest"
	"bitcoin/formating"
	"bitcoin/scripts"
	"errors"
)

/*
Converts an address to it's hash equivalent
The size of the address determines between P2WPKH and P2WSH.
Then Bech32 decodes the address removing network prefix, checksum,
witness version.
*/
func addressToSegwitProgram(address string, version int) (string, error) {
	v, data, err := bech32.DecodeBech32(address)
	if err != nil {
		return "", err
	}
	if v != version {
		return "", errors.New("invalid segwit version")
	}
	return formating.BytesToHex(data), nil
}

// Converts a script to it's hash equivalent
func ScriptToSegwitProgram(script scripts.Script) string {
	scriptBytes := script.ToBytes()
	toHash160 := digest.SingleHash(scriptBytes)
	return formating.BytesToHex(toHash160)
}

// instantiates an object from a witness program hex string
func P2WSHAddresssFromProgram(program string) (P2WSHAddresss, error) {
	return P2WSHAddresss{
		AddressProgram: SegwitAddress{
			Program:          program,
			Version:          constant.P2WSH_ADDRESS_V0,
			SegwitNumVersion: 0,
			Type:             P2WSH,
		},
	}, nil
}

// instantiates an object from address string encoding
func P2WSHAddresssFromAddress(address string) (P2WSHAddresss, error) {
	program, err := addressToSegwitProgram(address, 0)
	if err != nil {
		return P2WSHAddresss{
			AddressProgram: SegwitAddress{
				Program:          "",
				Version:          constant.P2WSH_ADDRESS_V0,
				SegwitNumVersion: 0,
				Type:             P2WSH,
			},
		}, err
	}
	return P2WSHAddresss{
		AddressProgram: SegwitAddress{
			Program:          program,
			Version:          constant.P2WSH_ADDRESS_V0,
			SegwitNumVersion: 0,
			Type:             P2WSH,
		},
	}, nil
}

// instantiates an object from a witness_script
func P2WSHAddresssFromScript(script scripts.Script) P2WSHAddresss {
	program := ScriptToSegwitProgram(script)
	return P2WSHAddresss{
		AddressProgram: SegwitAddress{
			Program:          program,
			Version:          constant.P2WSH_ADDRESS_V0,
			SegwitNumVersion: 0,
			Type:             P2WSH,
		},
	}
}

// instantiates an object from a witness program hex string
func P2WPKHAddresssFromProgram(program string) P2WPKHAddresss {
	return P2WPKHAddresss{
		AddressProgram: SegwitAddress{
			Program:          program,
			Version:          constant.P2WPKH_ADDRESS_V0,
			SegwitNumVersion: 0,
			Type:             P2WPKH,
		},
	}
}

// instantiates an object from address string encoding
func P2WPKHAddresssFromAddress(address string) P2WPKHAddresss {
	program, err := addressToSegwitProgram(address, 0)
	if err != nil {
		panic("invalid p2wpkh address")
	}
	return P2WPKHAddresss{
		AddressProgram: SegwitAddress{
			Program:          program,
			Version:          constant.P2WPKH_ADDRESS_V0,
			SegwitNumVersion: 0,
			Type:             P2WPKH,
		},
	}
}

// instantiates an object from a witness_script
func P2WPKHAddresssFromScript(script scripts.Script) P2WPKHAddresss {
	program := ScriptToSegwitProgram(script)
	return P2WPKHAddresss{
		AddressProgram: SegwitAddress{
			Program:          program,
			Version:          constant.P2WPKH_ADDRESS_V0,
			SegwitNumVersion: 0,
			Type:             P2WPKH,
		},
	}
}

// instantiates an object from a witness program hex string
func P2TRAddressFromProgram(program string) P2TRAddress {
	return P2TRAddress{
		AddressProgram: SegwitAddress{
			Program:          program,
			Version:          constant.P2TR_ADDRESS_V1,
			SegwitNumVersion: 1,
			Type:             P2TR,
		},
	}
}

// instantiates an object from address string encoding
func P2TRAddressFromAddress(address string) P2TRAddress {
	program, err := addressToSegwitProgram(address, 1)
	if err != nil {
		panic("invalid segwit program")
	}
	return P2TRAddress{
		AddressProgram: SegwitAddress{
			Program:          program,
			Version:          constant.P2TR_ADDRESS_V1,
			SegwitNumVersion: 1,
			Type:             P2TR,
		},
	}
}

/*
returns the address's string encoding (Bech32)
You can use NetworkParams or NetworkInfo (TestnetNetwork, BitcoinNetwork) in arguments
to select the network, otherwise, the default network is used.
*/
func (s SegwitAddress) toAddress(network ...interface{}) string {
	networkType := getNetworkParams(network...)
	bytes := formating.HexToBytes(s.Program)
	sw, err := bech32.EncodeBech32(networkType.bech32, s.SegwitNumVersion, bytes)
	if err != nil {
		panic("invalid segwit program")
	}

	return sw
}
