// Implementation of various Bitcoin address types, including P2PKH, P2SH, and SegWit
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
