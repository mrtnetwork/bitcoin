// Package address provides utilities for working with Bitcoin address types.
//
// Bitcoin uses different address types, such as Pay-to-Public-Key-Hash (P2PKH),
// Pay-to-Script-Hash (P2SH), and Segregated Witness (SegWit) addresses. This package
// offers functions for creating, parsing, and validating addresses of these types.
//
// Usage:
//  - Use the functions in this package to generate Bitcoin addresses.
//  - Parse Bitcoin addresses to extract information.
//  - Validate Bitcoin addresses to ensure they adhere to the correct format.
//	- Create address from scripts p2wsh, p2tr, tapleaf, ....
//
// Example:
//  package main
//
//  import (
//      "fmt"
//      "github.com/mrtnetwork/bitcoin/address"
//  )
//
//  func main() {
// P2PKH ADDRESS
// address in testnet:  myVMJgRi6arv4hLbeUcJYKUJWmFnpjtVme
// address in mainnet:  1JyQ1dLjHZRfHaryvudviQFyemf5vjbmUf
//	exampleAddr1, _ := address.P2PKHAddressFromAddress("myVMJgRi6arv4hLbeUcJYKUJWmFnpjtVme")
//	fmt.Println("address in testnet: ", exampleAddr1.Show(address.TestnetNetwork))
//	fmt.Println("address in mainnet: ", exampleAddr1.Show(address.MainnetNetwork))

// P2TR ADDRESS
// address in testnet:  tb1pyhmqwlcrws4dxcgalt4mrffgnys879vs59xf6sve4hazyvmhecxq3e6sc0
// address in mainnet:  bc1pyhmqwlcrws4dxcgalt4mrffgnys879vs59xf6sve4hazyvmhecxqx3vlzq
//	exampleAddr2, _ := address.P2TRAddressFromAddress("tb1pyhmqwlcrws4dxcgalt4mrffgnys879vs59xf6sve4hazyvmhecxq3e6sc0")
//	fmt.Println("address in testnet: ", exampleAddr2.Show(address.TestnetNetwork))
//	fmt.Println("address in mainnet: ", exampleAddr2.Show(address.MainnetNetwork))

// P2SH(P2PKH) ADDRESS
// address in testnet:  2N2yqygBJRvDzLzvPe91qKfSYnK5utGckJX
// address in mainnet:  3BRduwFGpTie9DHqy1PxhiTHZxsk7jpr9x
//	exampleAddr3, _ := address.P2SHAddressFromAddress("2N2yqygBJRvDzLzvPe91qKfSYnK5utGckJX", address.P2PKHInP2SH)
//	fmt.Println("address in testnet: ", exampleAddr3.Show(address.TestnetNetwork))
//	fmt.Println("address in mainnet: ", exampleAddr3.Show(address.MainnetNetwork))

// P2PKH ADDRESS
// address in testnet:  mzUzciYUGsNxLCaaHwou27F4RbnDTzKomV
// address in mainnet:  1Ky3KfTVTqwhZ66xaNqXCC2jZcBWZD6ppQ
//	exampleAddr4, _ := address.P2PKHAddressFromAddress("mzUzciYUGsNxLCaaHwou27F4RbnDTzKomV")
//	fmt.Println("address in testnet: ", exampleAddr4.Show(address.TestnetNetwork))
//	fmt.Println("address in mainnet: ", exampleAddr4.Show(address.MainnetNetwork))

// P2SH(P2PKH) ADDRESS
// address in testnet:  2MzibgEeJYCN8mjJsZTg79AH7au4PCkHXHo
// address in mainnet:  39APcViGvjrnZwgKtL4EXDHrNYrDT9YuCQ
//	exampleAddr5, _ := address.P2SHAddressFromAddress("2MzibgEeJYCN8mjJsZTg79AH7au4PCkHXHo", address.P2PKHInP2SH)
//	fmt.Println("address in testnet: ", exampleAddr5.Show(address.TestnetNetwork))
//	fmt.Println("address in mainnet: ", exampleAddr5.Show(address.MainnetNetwork))

// P2SH(P2WSH) ADDRESS
// address in testnet:  2N7bNV1WPwCVHfoqqRhvtmbAfktazfjHEW2
// address in mainnet:  3G3ARGaNKjywU2DHkaK29eBQYYNppD2xDE
//	exampleAddr6, _ := address.P2SHAddressFromAddress("2N7bNV1WPwCVHfoqqRhvtmbAfktazfjHEW2", address.P2WSHInP2SH)
//	fmt.Println("address in testnet: ", exampleAddr6.Show(address.TestnetNetwork))
//	fmt.Println("address in mainnet: ", exampleAddr6.Show(address.MainnetNetwork))

// P2SH(P2WPKH) ADDRESS
// address in testnet:  2N38S8G9q6qyEjPWmicqxrVwjq4QiTbcyf4
// address in mainnet:  3BaE4XDoVPTtXbtE3VE6EYxUciCYd5rspb
//	exampleAddr7, _ := address.P2SHAddressFromAddress("2N38S8G9q6qyEjPWmicqxrVwjq4QiTbcyf4", address.P2WPKHInP2SH)
//	fmt.Println("address in testnet: ", exampleAddr7.Show(address.TestnetNetwork))
//	fmt.Println("address in mainnet: ", exampleAddr7.Show(address.MainnetNetwork))

// P2SH(P2PK) ADDRESS
// address in testnet:  2MugsNcgzLJ1HosnZyC2CfZVmgbMPK1XubR
// address in mainnet:  348fJskxiqVwc6A2J4QL3cWWUF9DbWjTni
//	exampleAddr8, _ := address.P2SHAddressFromAddress("2MugsNcgzLJ1HosnZyC2CfZVmgbMPK1XubR", address.P2PKInP2SH)
//	fmt.Println("address in testnet: ", exampleAddr8.Show(address.TestnetNetwork))
//	fmt.Println("address in mainnet: ", exampleAddr8.Show(address.MainnetNetwork))

// P2WPKH ADDRESS
// address in testnet:  tb1q6q9halaazasd42gzsc2cvv5xls295w7kawkhxy
// address in mainnet:  bc1q6q9halaazasd42gzsc2cvv5xls295w7khgdyah
//	exampleAddr9, _ := address.P2WPKHAddresssFromAddress("tb1q6q9halaazasd42gzsc2cvv5xls295w7kawkhxy")
//	fmt.Println("address in testnet: ", exampleAddr9.Show(address.TestnetNetwork))
//	fmt.Println("address in mainnet: ", exampleAddr9.Show(address.MainnetNetwork))
//  }
//
// This package aims to simplify Bitcoin address handling and make it easier to work
// with different Bitcoin address formats.

package address

import "github.com/mrtnetwork/bitcoin/scripts"

type AddressType int

// address access interface in multiple structs
type BitcoinAddress interface {
	ToScriptPubKey() *scripts.Script
	Show(network ...interface{}) string
	GetType() AddressType
}

const (

	/*
		P2PKH (Pay-to-Public-Key-Hash)
		This is the most common address type in Bitcoin.
		FAddresses start with the number "1."
		It represents a single public key hashed with SHA-256 and then with RIPEMD-160.
		Transactions sent to a P2PKH address can only be spent by providing a valid signature corresponding to the public key.
	*/
	P2PKH AddressType = iota
	/*
		P2SH (Pay-to-Script-Hash):
		Addresses start with the number "3."
		It allows for more complex scripts (such as multisig or custom scripts) to be used as the recipient's address.
		The actual script that unlocks the funds is included in the transaction input.
	*/
	P2WPKH
	/*
		P2PK (Pay-to-Public-Key):
		Addresses start with the number "1."
		It involves sending funds directly to a recipient's public key, without hashing.
		It's not commonly used because it exposes the recipient's public key on the blockchain, reducing privacy.
	*/
	P2PK
	/*
		P2TR (Pay-to-Taproot):
		An advanced address type introduced with the Taproot upgrade.
		It allows for greater flexibility and privacy by enabling complex scripts and conditions.
		Addresses start with "bc1p."
	*/
	P2TR
	/*
		P2WSH (Pay-to-Witness-Script-Hash):
		Part of the SegWit upgrade.
		Addresses start with "bc1."
		Similar to P2SH but for SegWit-compatible scripts.
		Allows for more efficient use of block space and enhanced security.
	*/
	P2WSH
	/*
		P2WSHInP2SH (Pay-to-Witness-Script-Hash inside Pay-to-Script-Hash):
		Addresses start with the number "3."
		Combines the benefits of both P2SH and SegWit.
		It involves sending funds to a P2SH address, where the redeem script is a witness script (SegWit-compatible).
	*/
	P2WSHInP2SH
	/*
		P2WPKHInP2SH (Pay-to-Witness-Script-Hash inside Pay-to-Script-Hash):
		Addresses start with the number "3."
		Combines the benefits of both P2SH and SegWit.
		It involves sending funds to a P2SH address, where the redeem script is a witness script (SegWit-compatible).
	*/
	P2WPKHInP2SH
	P2PKInP2SH
	P2PKHInP2SH
)

type LegacyAddress struct {
	/*
	   the hash160 string representation of the address; hash160 represents
	   two consequtive hashes of the public key or the redeam script, first
	   a SHA-256 and then an RIPEMD-160
	*/
	Hash160 string
	Type    AddressType
}

type SegwitAddress struct {
	/*
	   for segwit v0 this is the hash string representation of either the address;
	   it can be either a public key hash (P2WPKH) or the hash of the script (P2WSH)
	   for segwit v1 (aka taproot) this is the public key
	*/
	Program          string
	Version          string
	SegwitNumVersion int
	Type             AddressType
}

/*
	P2SH (Pay-to-Script-Hash):
	Addresses start with the number "3."
	It allows for more complex scripts (such as multisig or custom scripts) to be used as the recipient's address.
	The actual script that unlocks the funds is included in the transaction input.
*/
type P2SHAdress struct {
	AddressProgram LegacyAddress
}

/*
	P2PKH (Pay-to-Public-Key-Hash)
	This is the most common address type in Bitcoin.
	FAddresses start with the number "1."
	It represents a single public key hashed with SHA-256 and then with RIPEMD-160.
	Transactions sent to a P2PKH address can only be spent by providing a valid signature corresponding to the public key.
*/
type P2PKHAddress struct {
	AddressProgram LegacyAddress
}

type P2PKAddress struct {
	AddressProgram LegacyAddress
}

/*
	P2TR (Pay-to-Taproot):
	An advanced address type introduced with the Taproot upgrade.
	It allows for greater flexibility and privacy by enabling complex scripts and conditions.
	Addresses start with "bc1p."
*/
type P2TRAddress struct {
	AddressProgram SegwitAddress
}

/*
	P2WPKH (Pay-to-Witness-Public-Key-Hash):
	Also known as Bech32 addresses.
	Addresses start with "bc1."
	Part of the Segregated Witness (SegWit) upgrade.
	It represents a single public key hashed with SHA-256 and then with RIPEMD-160, but it's a SegWit-compatible format.
*/
type P2WPKHAddresss struct {
	AddressProgram SegwitAddress
}

/*
	P2WSH (Pay-to-Witness-Script-Hash):
	Part of the SegWit upgrade.
	Addresses start with "bc1."
	Similar to P2SH but for SegWit-compatible scripts.
	Allows for more efficient use of block space and enhanced security.
*/
type P2WSHAddresss struct {
	AddressProgram SegwitAddress
}

// Access to the address program LegacyAddress
func (s P2SHAdress) Program() LegacyAddress {
	return s.AddressProgram
}

// Access to the address program LegacyAddress
func (s P2PKAddress) Program() LegacyAddress {
	return s.AddressProgram
}

// Access to the address program LegacyAddress
func (s P2PKHAddress) Program() LegacyAddress {
	return s.AddressProgram
}

// returns the type of address
func (s P2SHAdress) GetType() AddressType {
	return s.AddressProgram.Type
}

// returns the type of address
func (s P2PKHAddress) GetType() AddressType {
	return s.AddressProgram.Type
}

// returns the type of address
func (s P2PKAddress) GetType() AddressType {
	return s.AddressProgram.Type
}

// Access to the address program SegwitAddress
func (s P2TRAddress) Program() SegwitAddress {
	return s.AddressProgram
}

// Access to the address program SegwitAddress
func (s P2WPKHAddresss) Program() SegwitAddress {
	return s.AddressProgram
}

// Access to the address program SegwitAddress
func (s P2WSHAddresss) Program() SegwitAddress {
	return s.AddressProgram
}

// returns the type of address
func (s P2TRAddress) GetType() AddressType {
	return s.AddressProgram.Type
}

// returns the type of address
func (s P2WPKHAddresss) GetType() AddressType {
	return s.AddressProgram.Type
}

// returns the type of address
func (s P2WSHAddresss) GetType() AddressType {
	return s.AddressProgram.Type
}

/*
 a "scriptPubKey" (short for "script public key")
  refers to a script that defines the conditions
  that must be satisfied in order to spend funds
  from a particular Bitcoin address. The script is
  associated with the output of a Bitcoin transaction
  and is used to lock the funds until the specified conditions are met.
*/
func (segwit SegwitAddress) toScriptPubKey() []string {
	switch segwit.Type {
	case P2WPKH:
		return []string{"OP_0", segwit.Program}
	case P2TR:
		return []string{"OP_1", segwit.Program}
	case P2WSH:
		return []string{"OP_0", segwit.Program}
	default:
		return []string{} // Default case, return an empty scriptPubKey
	}
}

/*
 a "scriptPubKey" (short for "script public key")
  refers to a script that defines the conditions
  that must be satisfied in order to spend funds
  from a particular Bitcoin address. The script is
  associated with the output of a Bitcoin transaction
  and is used to lock the funds until the specified conditions are met.
*/
func (b LegacyAddress) toScriptPubKey() []string {
	switch b.Type {
	case P2PKHInP2SH, P2PKInP2SH, P2WPKHInP2SH, P2WSHInP2SH:
		return []string{"OP_HASH160", b.Hash160, "OP_EQUAL"}
	case P2PKH:
		return []string{"OP_DUP", "OP_HASH160", b.Hash160, "OP_EQUALVERIFY", "OP_CHECKSIG"}
	case P2PK:
		return []string{b.Hash160, "OP_CHECKSIG"}
	default:
		return []string{} // Default case, return an empty scriptPubKey
	}
}

/*
 a "scriptPubKey" (short for "script public key")
  refers to a script that defines the conditions
  that must be satisfied in order to spend funds
  from a particular Bitcoin address. The script is
  associated with the output of a Bitcoin transaction
  and is used to lock the funds until the specified conditions are met.
*/
func (segwit SegwitAddress) ToScriptPubKey() *scripts.Script {
	return scripts.ToScript(segwit.toScriptPubKey())
}

/*
 a "scriptPubKey" (short for "script public key")
  refers to a script that defines the conditions
  that must be satisfied in order to spend funds
  from a particular Bitcoin address. The script is
  associated with the output of a Bitcoin transaction
  and is used to lock the funds until the specified conditions are met.
*/
func (b LegacyAddress) ToScriptPubKey() *scripts.Script {
	return scripts.ToScript(b.toScriptPubKey())
}

/*
address string encoded in the Bech32 format.
*/
func (s P2TRAddress) Show(network ...interface{}) string {
	return s.AddressProgram.toAddress(network...)
}

/*
address string encoded in the Bech32 format.
*/
func (s P2WPKHAddresss) Show(network ...interface{}) string {
	return s.AddressProgram.toAddress(network...)
}

/*
address string encoded in the Bech32 format.
*/
func (s P2WSHAddresss) Show(network ...interface{}) string {
	return s.AddressProgram.toAddress(network...)
}

/*
The method calculates the address checksum and returns the Base58-encoded Bitcoin legacy address.
*/
func (s P2SHAdress) Show(network ...interface{}) string {
	return s.AddressProgram.toAddress(network...)
}

/*
The method calculates the address checksum and returns the Base58-encoded Bitcoin legacy address.
*/
func (s P2PKAddress) Show(network ...interface{}) string {
	return s.AddressProgram.toAddress(network...)
}

/*
The method calculates the address checksum and returns the Base58-encoded Bitcoin legacy address.
*/
func (s P2PKHAddress) Show(network ...interface{}) string {
	return s.AddressProgram.toAddress(network...)
}

///

/*
 a "scriptPubKey" (short for "script public key")
  refers to a script that defines the conditions
  that must be satisfied in order to spend funds
  from a particular Bitcoin address. The script is
  associated with the output of a Bitcoin transaction
  and is used to lock the funds until the specified conditions are met.
*/
func (s P2TRAddress) ToScriptPubKey() *scripts.Script {
	return s.AddressProgram.ToScriptPubKey()
}

/*
 a "scriptPubKey" (short for "script public key")
  refers to a script that defines the conditions
  that must be satisfied in order to spend funds
  from a particular Bitcoin address. The script is
  associated with the output of a Bitcoin transaction
  and is used to lock the funds until the specified conditions are met.
*/
func (s P2WPKHAddresss) ToScriptPubKey() *scripts.Script {
	return s.AddressProgram.ToScriptPubKey()
}

/*
 a "scriptPubKey" (short for "script public key")
  refers to a script that defines the conditions
  that must be satisfied in order to spend funds
  from a particular Bitcoin address. The script is
  associated with the output of a Bitcoin transaction
  and is used to lock the funds until the specified conditions are met.
*/
func (s P2WSHAddresss) ToScriptPubKey() *scripts.Script {
	return s.AddressProgram.ToScriptPubKey()
}

/*
 a "scriptPubKey" (short for "script public key")
  refers to a script that defines the conditions
  that must be satisfied in order to spend funds
  from a particular Bitcoin address. The script is
  associated with the output of a Bitcoin transaction
  and is used to lock the funds until the specified conditions are met.
*/
func (s P2SHAdress) ToScriptPubKey() *scripts.Script {
	return s.AddressProgram.ToScriptPubKey()
}

/*
 a "scriptPubKey" (short for "script public key")
  refers to a script that defines the conditions
  that must be satisfied in order to spend funds
  from a particular Bitcoin address. The script is
  associated with the output of a Bitcoin transactionf
  and is used to lock the funds until the specified conditions are met.
*/
func (s P2PKAddress) ToScriptPubKey() *scripts.Script {
	return s.AddressProgram.ToScriptPubKey()
}

/*
 a "scriptPubKey" (short for "script public key")
  refers to a script that defines the conditions
  that must be satisfied in order to spend funds
  from a particular Bitcoin address. The script is
  associated with the output of a Bitcoin transaction
  and is used to lock the funds until the specified conditions are met.
*/
func (s P2PKHAddress) ToScriptPubKey() *scripts.Script {
	return s.AddressProgram.ToScriptPubKey()
}

// types
