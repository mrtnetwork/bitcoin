package address

import (
	"bytes"
	"fmt"
	"github.com/mrtnetwork/bitcoin/base58"
	"github.com/mrtnetwork/bitcoin/digest"
	"github.com/mrtnetwork/bitcoin/ecc"
	"github.com/mrtnetwork/bitcoin/formating"
	"github.com/mrtnetwork/bitcoin/scripts"
	"math/big"
)

// IsValidAddress checks the validity of a Bitcoin address. It verifies whether the input address is
// well-formed and complies with the specified address type and network. The function performs the following checks:
// 1. Ensures that the address length is within a valid range (between 26 and 35 characters).
// 2. Decodes the base58-encoded address and verifies the checksum.
// 3. Validates the address type and network based on the decoded address prefix.
//
// Parameters:
// - address: A string representing the Bitcoin address to be validated.
// - addressType: An AddressType indicating the expected type of the address (e.g., P2PKH, P2PKHInP2SH, etc.).
// - network: A pointer to networkInfo representing the expected network configuration (can be nil for default networks).
//
// Returns:
// - bool: True if the input address is valid according to the specified criteria; otherwise, false.
func IsValidAddress(address string, addressType AddressType, network *networkInfo) bool {
	if len(address) < 26 || len(address) > 35 {
		return false
	}

	decode, err := base58.Decode(address)
	if err != nil {
		return false
	}
	networkPrefix := int(decode[0])
	data := decode[:len(decode)-4]
	checksum := decode[len(decode)-4:]
	hash := digest.DoubleHash(data)
	check := hash[:4]
	if !bytes.Equal(checksum, check) {
		return false
	}
	switch addressType {
	case P2PKH:
		if network != nil {
			return networkPrefix == network.p2PKHPrefix
		}
		return networkPrefix == TestnetNetwork.p2PKHPrefix || networkPrefix == MainnetNetwork.p2PKHPrefix

	case P2PKHInP2SH, P2PKInP2SH, P2WPKHInP2SH, P2WSHInP2SH:
		if network != nil {
			return networkPrefix == network.p2SHPrefix
		}
		return networkPrefix == TestnetNetwork.p2SHPrefix || networkPrefix == MainnetNetwork.p2SHPrefix

	default:
		{
			return true
		}
	}
}

// IsValidHash160 checks the validity of a hash160 string. It verifies whether the input string
// represents a valid hash160 by ensuring that it has a length of 40 characters and can be successfully
// parsed as a hexadecimal number. If the input hash160 is valid, the function returns true; otherwise,
// it returns false.
//
// Parameters:
// - hash160: A string representing the hash160 value to be validated.
//
// Returns:
// - bool: True if the input hash160 is valid; otherwise, false.
func IsValidHash160(hash160 string) bool {
	if len(hash160) != 40 {
		return false
	}
	_, err := new(big.Int).SetString(hash160, 16)
	return err
}

// fromHash160 instantiates an object from a hash160 hexadecimal string.
// The function takes a hash160 string as input and validates its format using the `IsValidHash160` function.
// If the provided hash160 is in a valid format, it is returned as is. Otherwise, an error is returned
// indicating that the hash160 is invalid.
//
// Parameters:
// - hash160: A hexadecimal string representing the hash160.
//
// Returns:
// - string: The input hash160 string if it is in a valid format.
// - error: An error if the input hash160 is not in a valid format.
func fromHash160(hash160 string) (string, error) {
	// Validate the format of the provided hash160 using the `IsValidHash160` function.
	if !IsValidHash160(hash160) {
		return "", fmt.Errorf("invalid hash160")
	}

	// Return the input hash160 as is, indicating that it is in a valid format.
	return hash160, nil
}

// fromAddress derives the hash160 representation from a Bitcoin address.
// The function takes a Bitcoin address string, an AddressType indicating the expected type of the address,
// and optional network configurations for address validation. It first determines the selected network based
// on the optional network configurations and validates the provided address using the `IsValidAddress` function.
// If the address is valid, it calls the `addressToHash160` function to convert the Bitcoin address to its hash160 form.
// Finally, it returns the hash160 as a hexadecimal string.
//
// Parameters:
// - address: A string representing a Bitcoin address.
// - addressType: An AddressType indicating the expected type of the address (e.g., P2PKH, P2SH, etc.).
// - network: Optional network configurations for address validation.
//
// Returns:
// - string: A hexadecimal string representing the hash160 of the provided Bitcoin address.
// - error: An error if the address is invalid or cannot be converted to hash160, or if there are issues during validation.
func fromAddress(address string, addressType AddressType, network ...interface{}) (string, error) {
	// Determine the selected network based on optional network configurations.
	selectedNetwork := getNetworkParams(false, network...)

	// Validate the provided Bitcoin address using the `IsValidAddress` function.
	if !IsValidAddress(address, addressType, selectedNetwork) {
		return "", fmt.Errorf("invalid address")
	}

	// Convert the valid Bitcoin address to its hash160 form using the `addressToHash160` function.
	toHash160, err := addressToHash160(address)
	if err != nil {
		return "", err
	}

	// Return the hash160 as a hexadecimal string representation.
	return toHash160, nil
}

// fromScript instantiates an object from a redeem_script (Script) and returns its hash160 representation.
// The function first calls the `scriptToHash160` function to compute the hash160 of the provided Script object.
// The resulting hash160 is then returned as a hexadecimal string.
//
// Parameters:
// - script: A Script object from which to derive the hash160 representation.
//
// Returns:
// - string: A hexadecimal string representing the hash160 of the provided Script.
func fromScript(script *scripts.Script) string {
	// Compute the hash160 of the provided Script object.
	toHash160 := scriptToHash160(script)

	// Return the hash160 as a hexadecimal string representation.
	return toHash160
}

// addressToHash160 converts a Bitcoin address to its hash160 representation.
// The function first decodes the provided Bitcoin address from Base58 encoding.
// It then extracts the hash160 bytes from the decoded address, excluding the
// address version byte and the checksum at the end. The extracted hash160 is
// converted to a hexadecimal string representation and returned.
//
// Parameters:
// - address: A string representing a Bitcoin address in Base58 encoding.
//
// Returns:
// - string: A hexadecimal string representing the hash160 of the Bitcoin address.
// - error: An error if the input address is invalid or cannot be decoded.
func addressToHash160(address string) (string, error) {
	dec, err := base58.Decode(address)
	if err != nil {
		return "", fmt.Errorf("invalid addresss")
	}
	hash160 := dec[1 : len(dec)-4]
	return formating.BytesToHex(hash160), nil
}

// scriptToHash160 converts a Script object to a hash160 string.
// The function first converts the Script object to its binary representation using `ToBytes()`.
// Then, it computes the hash160 of the script bytes using the `digest.Hash160` function.
// Finally, the computed hash160 is converted to a hexadecimal string representation and returned.
//
// Parameters:
// - script: A Script object to be converted to hash160.
//
// Returns:
// - string: A hexadecimal string representing the hash160 of the script.
func scriptToHash160(script *scripts.Script) string {
	scriptBytes := script.ToBytes()
	toHash160 := digest.Hash160(scriptBytes)
	return formating.BytesToHex(toHash160)
}

// P2SHAddressFromHash160 creates a P2SH (Pay-to-Script-Hash) address object from a hash160 string.
// The function uses the `fromHash160` function to convert the provided hash160 string to its canonical form
// and validate its correctness. If the hash160 is valid, the function constructs and returns a P2SHAddress
// object with the converted hash160, setting its address type to P2PKHInP2SH.
//
// Parameters:
// - hash160: A string representing the hash160 value from which to create the P2SH address.
//
// Returns:
// - *P2SHAddress: A P2SHAddress object representing the generated P2SH address.
// - error: An error if there are any issues during address creation, or if the input hash160 is invalid.
func P2SHAddressFromHash160(hash160 string) (*P2SHAdress, error) {
	bip, err := fromHash160(hash160)
	if err != nil {
		return nil, err
	}
	return &P2SHAdress{LegacyAddress{bip, P2PKHInP2SH}}, nil
}

// P2SHAddressFromAddress creates a P2SH (Pay-to-Script-Hash) address object from an address string.
// The function uses the `fromAddress` function to extract the relevant information from the provided
// address string. It also accepts optional `args` parameters, including an `AddressType` indicating the
// expected type of the address and network configurations for address validation. If provided, the function
// checks if the address matches the expected network and address type; if not, an error is returned. Finally,
// the function constructs and returns a P2SHAddress object with the extracted data, setting its address type
// based on the provided parameters.
//
// Parameters:
//   - address: A string representing the address from which to create the P2SH address.
//   - args: Optional arguments, including an `AddressType` (if specified) and network configurations for
//     address validation.
//
// Returns:
//   - *P2SHAddress: A P2SHAddress object representing the generated P2SH address.
//   - error: An error if there are any issues during address creation, if the address is invalid,
//     or if it does not match the expected network or address type (if provided).
func P2SHAddressFromAddress(address string, args ...interface{}) (*P2SHAdress, error) {
	addressType, e := getP2shAddressParam(args...)
	if e != nil {
		return nil, e
	}
	bip, err := fromAddress(address, addressType, args...)
	if err != nil {
		return nil, err
	}
	return &P2SHAdress{LegacyAddress{bip, addressType}}, nil
}

// P2SHAddressFromScript creates a P2SH (Pay-to-Script-Hash) address object from a Script object.
// The function first uses the `fromScript` function to extract the relevant information from the provided
// Script object. It also accepts an `addressType` parameter, indicating the expected type of the address
// (e.g., P2SHInP2PK, P2SHInP2PKH). Finally, the function constructs and returns a P2SHAddress object with
// the extracted data, setting its address type based on the provided `addressType`.
//
// Parameters:
// - script: A Script object representing the script from which to create the P2SH address.
// - addressType: An AddressType indicating the expected type of the address (e.g., P2SHInP2PK, P2SHInP2PKH).
//
// Returns:
// - *P2SHAddress: A P2SHAddress object representing the generated P2SH address.
// - error: An error if there are any issues during address creation or if the script is invalid.
func P2SHAddressFromScript(script *scripts.Script, addressType AddressType) (*P2SHAdress, error) {
	bip := fromScript(script)
	return &P2SHAdress{LegacyAddress{bip, addressType}}, nil
}

// P2PKHAddressFromHash160 creates a P2PKH (Pay-to-Public-Key-Hash) address object from a hash160 string.
// The function uses the `fromHash160` function to convert the provided hash160 string to its canonical form
// and validate its correctness. If the hash160 is valid, the function constructs and returns a P2PKHAddress
// object with the converted hash160, setting its address type to P2PKH.
//
// Parameters:
// - hash160: A string representing the hash160 value from which to create the P2PKH address.
//
// Returns:
// - *P2PKHAddress: A P2PKHAddress object representing the generated P2PKH address.
// - error: An error if there are any issues during address creation, or if the input hash160 is invalid.
func P2PKHAddressFromHash160(hash160 string) (*P2PKHAddress, error) {
	bip, err := fromHash160(hash160)
	if err != nil {
		return nil, err
	}
	return &P2PKHAddress{LegacyAddress{bip, P2PKH}}, nil
}

// P2PKHAddressFromAddress creates a P2PKH (Pay-to-Public-Key-Hash) address object from an address string.
// The function first uses the `fromAddress` function to extract the relevant information from the provided
// address string. The `network` argument is optional and can be used to validate that the address belongs
// to a specific network configuration. If provided, the function checks if the address matches the expected
// network; if not, an error is returned. Finally, the function constructs and returns a P2PKHAddress object
// with the extracted data, setting its address type to P2PKH.
//
// Parameters:
// - address: A string representing the address from which to create the P2PKH address.
// - network: An optional network configuration for address validation (e.g., Mainnet, Testnet).
//
// Returns:
//   - *P2PKHAddress: A P2PKHAddress object representing the generated P2PKH address.
//   - error: An error if there are any issues during address creation, if the address is invalid,
//     or if the address does not match the expected network (if provided).
func P2PKHAddressFromAddress(address string, network ...interface{}) (*P2PKHAddress, error) {
	bip, err := fromAddress(address, P2PKH, network...)
	if err != nil {
		return nil, err
	}
	return &P2PKHAddress{LegacyAddress{bip, P2PKH}}, nil
}

// P2PKHAddressFromScript creates a P2PKH (Pay-to-Public-Key-Hash) address object from a Script object.
// The function extracts the relevant information from the provided Script object using the `fromScript`
// function and constructs a P2PKHAddress object with the derived data, setting its address type to P2PKH.
func P2PKHAddressFromScript(script *scripts.Script) (*P2PKHAddress, error) {
	bip := fromScript(script)
	return &P2PKHAddress{LegacyAddress{bip, P2PKH}}, nil
}

// P2PKAddressFromPublicKey creates a P2PK (Pay-to-Public-Key) address object from a public key string.
// The function first converts the provided public key string to a byte slice and checks if it represents
// a valid elliptic curve point.
//
// Parameters:
// - public: A string representing the public key to create the address from.
//
// Returns:
// - *P2PKAddress: A P2PKAddress object representing the generated P2PK address.
// - error: An error if the public key is invalid or if there are any issues during address creation.
func P2PKAddressFromPublicKey(public string) (*P2PKAddress, error) {
	toPublicBytes := formating.HexToBytes(public)
	if !ecc.IsPoint(toPublicBytes) {
		return nil, fmt.Errorf("invalid public key")
	}
	return &P2PKAddress{LegacyAddress{Hash160: public, Type: P2PK}}, nil

}

/*
toAddress generates a Bitcoin legacy address from the given hash160 and address type.
You can specify the desired Bitcoin network by passing network parameters.
Supported address types are P2PKH, P2PK, and P2SH.
The method calculates the address checksum and returns the Base58-encoded Bitcoin legacy address.
*/
func (s LegacyAddress) toAddress(network ...interface{}) string {
	networkType := getNetworkParams(true, network...)
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

/*
Find the P2SHAddress type in the function parameters
args must be AddressType or P2SHAddressTypeParam
Default P2PKInP2SH
Choosing the type of address for transaction output is not important,
but for transaction input and spending from this address,
an error will be encountered if the transaction is entered incorrectly.
*/
func getP2shAddressParam(args ...interface{}) (AddressType, error) {
	argruments := formating.FlattenList(args)
	defaultP2SHType := P2PKInP2SH
	for _, opt := range argruments {
		switch v := opt.(type) {
		case P2SHAddressTypeParam:
			if v.AddressType != nil {
				defaultP2SHType = *v.AddressType
				break
			}
		case AddressType:
			defaultP2SHType = v
		}
	}
	switch defaultP2SHType {
	case P2PKHInP2SH, P2WPKHInP2SH, P2WSHInP2SH, P2PKInP2SH:
		{
			break
		}
	default:
		return 0, fmt.Errorf("invalid p2sh address use one of P2PKHInP2SH,P2WPKHInP2SH,P2WSHInP2SH,P2PKInP2SH")
	}

	return defaultP2SHType, nil
}
