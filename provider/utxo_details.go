package provider

import (
	"fmt"
	"github.com/mrtnetwork/bitcoin/address"
	"github.com/mrtnetwork/bitcoin/keypair"
	"math/big"
)

// UtxoOwnerDetails represents ownership details associated with a Bitcoin unspent transaction output (UTXO).
// It includes information such as the public key, Bitcoin address, and multi-signature address (if applicable)
// of the UTXO owner.
type UtxoOwnerDetails struct {
	// PublicKey is the public key associated with the UTXO owner.
	PublicKey string

	// Address is the Bitcoin address associated with the UTXO owner.
	Address address.BitcoinAddress

	// MultiSigAddress is a pointer to a MultiSignaturAddress instance representing a multi-signature address
	// associated with the UTXO owner. It may be nil if the UTXO owner is not using a multi-signature scheme.
	MultiSigAddress *MultiSignaturAddress
	// Utxo is a BitcoinUtxo instance representing the unspent transaction output.
	Utxo BitcoinUtxo
}

// UtxoWithOwner represents an unspent transaction output (UTXO) along with its associated owner details.
// It combines information about the UTXO itself (BitcoinUtxo) and the ownership details (UtxoOwnerDetails).
type UtxoWithOwner struct {
	// Utxo is a BitcoinUtxo instance representing the unspent transaction output.
	Utxo BitcoinUtxo

	// OwnerDetails is a UtxoOwnerDetails instance containing information about the UTXO owner.
	OwnerDetails UtxoOwnerDetails
}

// UtxoWithOwnerList is a slice of UtxoWithOwner instances, representing a list of Bitcoin UTXOs along with their
// respective ownership details. It is typically used to manage and process multiple UTXOs with ownership information.
type UtxoWithOwnerList []UtxoWithOwner

// BitcoinOutputDetails represents details about a Bitcoin transaction output, including
// the recipient address and the value of bitcoins sent to that address.
type BitcoinOutputDetails struct {
	// Address is a Bitcoin address representing the recipient of the transaction output.
	Address address.BitcoinAddress

	// Value is a pointer to a big.Int representing the amount of bitcoins sent to the recipient.
	Value *big.Int
}

// BitcoinUtxo represents an unspent transaction output (UTXO) on the Bitcoin blockchain.
// It includes details such as the transaction hash (TxHash), the amount of bitcoins (Value),
// the output index (Vout), the script type (ScriptType), and the block height at which the UTXO
// was confirmed (BlockHeight).
type BitcoinUtxo struct {
	// TxHash is the unique identifier of the transaction containing this UTXO.
	TxHash string

	// Value is a pointer to a big.Int representing the amount of bitcoins associated with this UTXO.
	Value *big.Int

	// Vout is the output index within the transaction that corresponds to this UTXO.
	Vout int

	// ScriptType specifies the type of Bitcoin script associated with this UTXO.
	ScriptType address.AddressType

	// BlockHeight represents the block height at which this UTXO was confirmed.
	BlockHeight int
}

// IsP2tr checks whether the BitcoinUtxo instance represents a Pay-to-Taproot (P2TR) UTXO
// based on its script type. It returns true if the script type is P2TR, indicating that the
// UTXO is associated with a Taproot address; otherwise, it returns false.
//
// Returns:
// - bool: True if the UTXO is of P2TR type, false otherwise.
func (utxo *BitcoinUtxo) IsP2tr() bool {
	return utxo.ScriptType == address.P2TR
}

// IsSegwit checks whether the BitcoinUtxo instance represents a Segregated Witness (SegWit) UTXO
// based on its script type. It returns true if the script type is any of the recognized SegWit types,
// which include P2WPKH, P2WSH, P2TR, or if it is a Pay-to-Witness-Public-Key-Hash (P2WPKH) wrapped
// in a Pay-to-Script-Hash (P2SH) script.
//
// Returns:
// - bool: True if the UTXO is of SegWit type, false otherwise.
func (utxo *BitcoinUtxo) IsSegwit() bool {
	return utxo.ScriptType == address.P2WPKH ||
		utxo.ScriptType == address.P2WSH ||
		utxo.ScriptType == address.P2TR ||
		utxo.IsP2shSegwit()
}

// IsP2shSegwit checks whether the BitcoinUtxo instance represents a Pay-to-Script-Hash (P2SH)
// Segregated Witness (SegWit) UTXO based on its script type. It returns true if the script type
// indicates that the UTXO is either a Pay-to-Witness-Public-Key-Hash-in-Pay-to-Script-Hash (P2WPKH-in-P2SH)
// or a Pay-to-Witness-Script-Hash-in-Pay-to-Script-Hash (P2WSH-in-P2SH).
//
// Returns:
// - bool: True if the UTXO is of P2SH SegWit type, false otherwise.
func (utxo *BitcoinUtxo) IsP2shSegwit() bool {
	return utxo.ScriptType == address.P2WPKHInP2SH ||
		utxo.ScriptType == address.P2WSHInP2SH
}

// NewUtxoWithOwner creates a new instance of UtxoWithOwner by combining a BitcoinUtxo representing an
// unspent transaction output (UTXO) and UtxoOwnerDetails containing ownership information. It returns a
// pointer to the created UtxoWithOwner instance.
//
// Parameters:
// - utxo: A BitcoinUtxo instance representing the unspent transaction output.
// - owner: UtxoOwnerDetails containing ownership information associated with the UTXO.
//
// Returns:
// - *UtxoWithOwner: A pointer to the UtxoWithOwner instance combining the UTXO and ownership details.
func NewUtxoWithOwner(utxo BitcoinUtxo, owner UtxoOwnerDetails) *UtxoWithOwner {
	return &UtxoWithOwner{
		Utxo:         utxo,
		OwnerDetails: owner,
	}
}

// Public retrieves the public key associated with the UTXO owner. If the UTXO owner is using a multi-signature
// address, it returns an error indicating that public key access is not available for multi-signature addresses.
// Otherwise, it returns the ECDSA public key as an *keypair.ECPublic.
//
// Returns:
// - *keypair.ECPublic: A pointer to the ECDSA public key of the UTXO owner.
// - error: An error if public key access is not available for multi-signature addresses or other issues.
func (utxo *UtxoWithOwner) Public() (*keypair.ECPublic, error) {
	if utxo.IsMultiSig() {
		return nil, fmt.Errorf("cannot access public in multisig address; use owner's public keys")
	}
	return keypair.NewECPPublicFromHex(utxo.OwnerDetails.PublicKey)
}

// IsMultiSig checks whether the UTXO owner is using a multi-signature address based on the presence of a
// MultiSignaturAddress instance in the ownership details. It returns true if the UTXO owner is using a multi-
// signature address, indicating that the UTXO is associated with a multi-signature scheme.
//
// Returns:
// - bool: True if the UTXO owner is using a multi-signature address, false otherwise.
func (utxo *UtxoWithOwner) IsMultiSig() bool {
	return utxo.OwnerDetails.MultiSigAddress != nil
}

// SumOfUtxosValue calculates and returns the total value of all UTXOs in the UtxoWithOwnerList. It iterates
// through each UTXO in the list and adds their values to compute the sum of UTXO values.
//
// Returns:
// - *big.Int: A pointer to a big.Int representing the total sum of UTXO values.
func (utxos UtxoWithOwnerList) SumOfUtxosValue() *big.Int {
	sum := big.NewInt(0)
	for _, utxo := range utxos {
		sum.Add(sum, utxo.Utxo.Value)
	}
	return sum
}

// CanSpending checks whether there are UTXOs in the UtxoWithOwnerList that can be spent. It calculates
// the total value of UTXOs in the list and returns true if the total value is greater than zero
// that there are spendable UTXOs available.
func (utxos UtxoWithOwnerList) CanSpending() bool {
	value := utxos.SumOfUtxosValue()
	return value.Cmp(big.NewInt(0)) != 0
}
