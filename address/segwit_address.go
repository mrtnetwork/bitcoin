package address

import (
	"fmt"
	"github.com/mrtnetwork/bitcoin/bech32"
	"github.com/mrtnetwork/bitcoin/constant"
	"github.com/mrtnetwork/bitcoin/digest"
	"github.com/mrtnetwork/bitcoin/formating"
	"github.com/mrtnetwork/bitcoin/scripts"
)

// addressToSegwitProgram extracts the witness program (data part) from a SegWit (Bech32) address.
// The function takes a SegWit address string, an expected version (witness version), and optional network configurations
// for address validation. It uses the Bech32 decoder to decode the provided address and extract the witness program data.
// It then validates the version and, if specified, the network of the address. If all validation checks pass, the extracted
// witness program data is returned as a hexadecimal string.
//
// Parameters:
// - address: A string representing a SegWit (Bech32) address.
// - version: The expected witness version (e.g., 0 for P2WPKH, 1 for P2WSH).
// - network: Optional network configurations for address validation.
//
// Returns:
// - string: A hexadecimal string representing the witness program data.
// - error: An error if the address is invalid, the version is incorrect, the network is invalid, or other decoding issues occur.
func addressToSegwitProgram(address string, version int, network ...interface{}) (string, error) {
	// Determine the selected network based on optional network configurations.
	selectedNetwork := getNetworkParams(false, network...)

	// Decode the provided SegWit (Bech32) address using the Bech32 decoder.
	v, data, hrp, err := bech32.DecodeBech32(address)
	if err != nil {
		return "", err
	}

	// Validate that the extracted witness version matches the expected version.
	if v != version {
		return "", fmt.Errorf("invalid SegWit version")
	}

	// If a network is specified, validate that the address belongs to the expected network.
	if selectedNetwork != nil {
		if hrp != selectedNetwork.bech32 {
			return "", fmt.Errorf("invalid network address, address does not belong to %v", selectedNetwork.name)
		}
	}

	// Return the witness program data as a hexadecimal string.
	return formating.BytesToHex(data), nil
}

// ScriptToSegwitProgram converts a Bitcoin script to its hash equivalent as a witness program.
// The function takes a Bitcoin script and computes its hash (specifically, a single hash) to create a
// witness program. It then returns the resulting hash as a hexadecimal string.
//
// Parameters:
// - script: A Script object representing the Bitcoin script.
//
// Returns:
// - string: A hexadecimal string representing the hash-based witness program.
func ScriptToSegwitProgram(script *scripts.Script) string {
	// Convert the provided Bitcoin script to its byte representation.
	scriptBytes := script.ToBytes()

	// Compute a hash of the script to create a witness program.
	toHash160 := digest.SingleHash(scriptBytes)

	// Return the hash as a hexadecimal string representation.
	return formating.BytesToHex(toHash160)
}

// P2WSHAddresssFromProgram instantiates a P2WSH (Pay-to-Witness-Script-Hash) address object
// from a witness program hexadecimal string. The function takes the witness program as input,
// sets the appropriate version and type for a P2WSH address, and returns the resulting P2WSH address object.
//
// Parameters:
// - program: A hexadecimal string representing the witness program.
//
// Returns:
// - *P2WSHAddresss: A P2WSH address object created from the provided witness program.
// - error: An error if there are issues during object instantiation or if the witness program is invalid.
func P2WSHAddresssFromProgram(program string) (*P2WSHAddresss, error) {
	// Create a P2WSH address object with the specified witness program, version, and type.
	return &P2WSHAddresss{
		AddressProgram: SegwitAddress{
			Program:          program,
			Version:          constant.P2WSH_ADDRESS_V0,
			SegwitNumVersion: 0,
			Type:             P2WSH,
		},
	}, nil
}

// P2WSHAddresssFromAddress instantiates a P2WSH (Pay-to-Witness-Script-Hash) address object
// from a SegWit (Bech32) address string. The function takes the SegWit address as input and extracts
// the witness program from it. It then sets the appropriate version and type for a P2WSH address
// and returns the resulting P2WSH address object.
//
// Parameters:
// - address: A SegWit (Bech32) address string.
// - network: Optional network configurations for address validation.
//
// Returns:
//   - *P2WSHAddresss: A P2WSH address object created from the provided SegWit address.
//   - error: An error if there are issues during object instantiation, if the address is invalid,
//     or if the extracted witness program is invalid.
func P2WSHAddresssFromAddress(address string, network ...interface{}) (*P2WSHAddresss, error) {
	// Extract the witness program from the provided SegWit (Bech32) address.
	program, err := addressToSegwitProgram(address, 0, network...)
	if err != nil {
		return nil, err
	}

	// Create a P2WSH address object with the specified witness program, version, and type.
	return &P2WSHAddresss{
		AddressProgram: SegwitAddress{
			Program:          program,
			Version:          constant.P2WSH_ADDRESS_V0,
			SegwitNumVersion: 0,
			Type:             P2WSH,
		},
	}, nil
}

// P2WSHAddresssFromScript instantiates a P2WSH (Pay-to-Witness-Script-Hash) address object
// from a witness script. The function takes the witness script as input, computes the witness program
// by hashing the script, and sets the appropriate version and type for a P2WSH address. It then returns
// the resulting P2WSH address object.
//
// Parameters:
// - script: A Script object representing the witness script.
//
// Returns:
// - *P2WSHAddresss: A P2WSH address object created from the provided witness script.
// - error: An error if there are issues during object instantiation or if the witness script is invalid.
func P2WSHAddresssFromScript(script *scripts.Script) (*P2WSHAddresss, error) {
	// Compute the witness program by hashing the provided witness script.
	program := ScriptToSegwitProgram(script)

	// Create a P2WSH address object with the specified witness program, version, and type.
	return &P2WSHAddresss{
		AddressProgram: SegwitAddress{
			Program:          program,
			Version:          constant.P2WSH_ADDRESS_V0,
			SegwitNumVersion: 0,
			Type:             P2WSH,
		},
	}, nil
}

// P2WPKHAddresssFromProgram instantiates a P2WPKH (Pay-to-Witness-Public-Key-Hash) address object
// from a witness program hexadecimal string. The function takes the witness program as input, sets the
// appropriate version and type for a P2WPKH address, and returns the resulting P2WPKH address object.
//
// Parameters:
// - program: A hexadecimal string representing the witness program.
//
// Returns:
// - *P2WPKHAddresss: A P2WPKH address object created from the provided witness program.
// - error: An error if there are issues during object instantiation or if the witness program is invalid.
func P2WPKHAddresssFromProgram(program string) (*P2WPKHAddresss, error) {
	// Create a P2WPKH address object with the specified witness program, version, and type.
	return &P2WPKHAddresss{
		AddressProgram: SegwitAddress{
			Program:          program,
			Version:          constant.P2WPKH_ADDRESS_V0,
			SegwitNumVersion: 0,
			Type:             P2WPKH,
		},
	}, nil
}

// P2WPKHAddresssFromAddress instantiates a P2WPKH (Pay-to-Witness-Public-Key-Hash) address object
// from a SegWit (Bech32) address string. The function takes the SegWit address as input and extracts
// the witness program from it. It then sets the appropriate version and type for a P2WPKH address
// and returns the resulting P2WPKH address object.
//
// Parameters:
// - address: A SegWit (Bech32) address string.
// - network: Optional network configurations for address validation.
//
// Returns:
//   - *P2WPKHAddresss: A P2WPKH address object created from the provided SegWit address.
//   - error: An error if there are issues during object instantiation, if the address is invalid,
//     or if the extracted witness program is invalid.
func P2WPKHAddresssFromAddress(address string, network ...interface{}) (*P2WPKHAddresss, error) {
	// Extract the witness program from the provided SegWit (Bech32) address.
	program, err := addressToSegwitProgram(address, 0)
	if err != nil {
		return nil, err
	}

	// Create a P2WPKH address object with the specified witness program, version, and type.
	return &P2WPKHAddresss{
		AddressProgram: SegwitAddress{
			Program:          program,
			Version:          constant.P2WPKH_ADDRESS_V0,
			SegwitNumVersion: 0,
			Type:             P2WPKH,
		},
	}, nil
}

// P2WPKHAddresssFromScript instantiates a P2WPKH (Pay-to-Witness-Public-Key-Hash) address object
// from a witness script. The function takes the witness script as input, derives the witness program
// from it, and sets the appropriate version and type for a P2WPKH address. It then returns the
// resulting P2WPKH address object.
//
// Parameters:
// - script: A witness script.
//
// Returns:
// - *P2WPKHAddresss: A P2WPKH address object created from the provided witness script.
// - error: An error if there are issues during object instantiation or if the witness script is invalid.
func P2WPKHAddresssFromScript(script *scripts.Script) (*P2WPKHAddresss, error) {
	// Derive the witness program from the provided witness script.
	program := ScriptToSegwitProgram(script)

	// Create a P2WPKH address object with the specified witness program, version, and type.
	return &P2WPKHAddresss{
		AddressProgram: SegwitAddress{
			Program:          program,
			Version:          constant.P2WPKH_ADDRESS_V0,
			SegwitNumVersion: 0,
			Type:             P2WPKH,
		},
	}, nil
}

// P2TRAddressFromProgram instantiates a P2TR (Pay-to-Taproot) address object from a Taproot program.
// The function takes the Taproot program as input and sets the appropriate version and type for a P2TR address.
// It then returns the resulting P2TR address object.
//
// Parameters:
// - program: A Taproot program in hexadecimal string format.
//
// Returns:
// - *P2TRAddress: A P2TR address object created from the provided Taproot program.
// - error: An error if there are issues during object instantiation or if the Taproot program is invalid.
func P2TRAddressFromProgram(program string) (*P2TRAddress, error) {
	// Create a P2TR address object with the specified program, version, and type.
	return &P2TRAddress{
		AddressProgram: SegwitAddress{
			Program:          program,
			Version:          constant.P2TR_ADDRESS_V1,
			SegwitNumVersion: 1,
			Type:             P2TR,
		},
	}, nil
}

// P2TRAddressFromAddress instantiates a P2TR (Pay-to-Taproot) address object from a P2TR address string encoding.
// The function takes the P2TR address string and an optional network parameter as input and sets the appropriate
// version and type for a P2TR address. It then returns the resulting P2TR address object.
//
// Parameters:
// - address: A P2TR address string encoding.
// - network: An optional network parameter for address validation (e.g., Mainnet or Testnet).
//
// Returns:
// - *P2TRAddress: A P2TR address object created from the provided P2TR address string.
// - error: An error if there are issues during object instantiation or if the address is invalid.
func P2TRAddressFromAddress(address string, network ...interface{}) (*P2TRAddress, error) {
	// Derive the Taproot program from the provided P2TR address string encoding.
	program, err := addressToSegwitProgram(address, 1, network...)
	if err != nil {
		return nil, err
	}

	// Create a P2TR address object with the specified program, version, and type.
	return &P2TRAddress{
		AddressProgram: SegwitAddress{
			Program:          program,
			Version:          constant.P2TR_ADDRESS_V1,
			SegwitNumVersion: 1,
			Type:             P2TR,
		},
	}, nil
}

// toAddress converts a SegwitAddress object to a P2TR (Pay-to-Taproot) or P2WSH (Pay-to-Witness-Script-Hash) address string.
// The function encodes the SegwitAddress program in the specified network's Bech32 format, taking into account the
// SegwitNumVersion and network type.
//
// Parameters:
//   - network: An optional network parameter for address encoding (e.g., Mainnet or Testnet). If not provided, the default
//     network of the SegwitAddress object is used.
//
// Returns:
// - string: A P2TR or P2WSH or P2WPKH address string encoded in the Bech32 format.
func (s SegwitAddress) toAddress(network ...interface{}) string {
	// Determine the network type to be used for address encoding.
	networkType := getNetworkParams(true, network...)

	// Convert the SegwitAddress program from hexadecimal to bytes.
	bytes := formating.HexToBytes(s.Program)

	// Encode the program in the Bech32 format, including the SegwitNumVersion and network type.
	sw, _ := bech32.EncodeBech32(networkType.bech32, s.SegwitNumVersion, bytes)

	return sw
}
