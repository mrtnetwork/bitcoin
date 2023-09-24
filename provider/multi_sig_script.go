package provider

import (
	"fmt"
	"github.com/mrtnetwork/bitcoin/address"
	"github.com/mrtnetwork/bitcoin/keypair"
	"github.com/mrtnetwork/bitcoin/scripts"
)

// MultiSignatureSigner is an interface that defines methods required for representing
// signers in a multi-signature scheme. A multi-signature signer typically includes
// information about their public key and weight within the scheme.
type MultiSignatureSigner interface {
	// PublicKey returns the public key associated with the signer.
	PublicKey() string

	// Weight returns the weight or significance of the signer within the multi-signature scheme.
	// The weight is used to determine the number of signatures required for a valid transaction.
	Weight() int
}

type multiSignatureSigner struct {
	// PublicKey returns the public key associated with the signer.
	publicKey string
	// Weight returns the weight or significance of the signer within the multi-signature scheme.
	// The weight is used to determine the number of signatures required for a valid transaction.
	weight int
}

// PublicKey returns the public key associated with the signer.
func (signer *multiSignatureSigner) PublicKey() string {
	return signer.publicKey
}

// Weight returns the weight or significance of the signer within the multi-signature scheme.
// The weight is used to determine the number of signatures required for a valid transaction.
func (signer *multiSignatureSigner) Weight() int {
	return signer.weight
}

// CreateMultiSignatureSigner creates a new instance of a multi-signature signer with the
// specified public key and weight.
//
// Parameters:
//   - publicKey: A string representing the public key associated with the signer.
//   - weight: An integer indicating the weight or significance of the signer within the
//     multi-signature scheme. The weight is used to determine the number of
//     signatures required for a valid transaction.
//
// Returns:
//   - *multiSignatureSigner: A pointer to a multiSignatureSigner instance containing the
//     provided public key and weight.
//   - error: An error value if there was an issue creating the signer, or nil if successful.
func CreateMultiSignaturSigner(publicKey string, weight int) (*multiSignatureSigner, error) {
	_, e := keypair.NewECPPublicFromHex(publicKey)
	if e != nil {
		return nil, e
	}
	if weight > 16 || weight < 1 {
		return nil, fmt.Errorf("the weight of the owner's should not be more than 16 or less than 1")
	}
	return &multiSignatureSigner{
		publicKey: publicKey,
		weight:    weight,
	}, nil
}

// MultiSignaturAddress represents a multi-signature Bitcoin address configuration, including
// information about the required signers, threshold, the address itself,
// and the script details used for multi-signature transactions.
type MultiSignaturAddress struct {
	// Signers is a collection of signers participating in the multi-signature scheme.
	Signers MultiSignaturAddressSigners

	// Threshold is the minimum number of signatures required to spend the bitcoins associated
	// with this address.
	Threshold int

	// Address represents the Bitcoin address associated with this multi-signature configuration.
	Address address.BitcoinAddress

	// ScriptDetails provides details about the multi-signature script used in transactions,
	// including "OP_M", compressed public keys, "OP_N", and "OP_CHECKMULTISIG."
	ScriptDetails string
}

type MultiSignaturAddressSigners []MultiSignatureSigner

// CreateMultiSignatureAddress creates a new instance of a MultiSignaturAddress, representing
// a multi-signature Bitcoin address configuration. It allows you to specify the minimum number
// of required signatures (threshold), provide the collection of signers participating in the
// multi-signature scheme, and specify the address type.
//
// Parameters:
// - threshold: An integer indicating the minimum number of signatures required to spend the bitcoins.
// - signers: A collection of signers (MultiSignaturAddressSigners) participating in the multi-signature scheme.
// - addressType: An address.AddressType specifying the type of Bitcoin address to be associated with this configuration.
//
// Returns:
// - *MultiSignaturAddress: A pointer to a MultiSignaturAddress instance representing the created multi-signature address.
// - error: An error value if there was an issue creating the address, or nil if successful.
func CreateMultiSignatureAddress(threshold int, signers MultiSignaturAddressSigners, addressType address.AddressType) (*MultiSignaturAddress, error) {

	if threshold > 16 || threshold < 1 {
		return nil, fmt.Errorf("the threshold should not be greater than 16 and less than 1")
	}
	sumWeight := signers.SumThreshHold()
	if threshold > 16 || threshold < 1 {
		return nil, fmt.Errorf("the total weight of the owners should not be more than 16")
	}
	if sumWeight < threshold {
		return nil, fmt.Errorf("the total weight of the signatories should reach the threshold")
	}
	multiSigScript := []interface{}{fmt.Sprint("OP_", threshold)}
	for i := 0; i < len(signers); i++ {
		for w := 0; w < signers[i].Weight(); w++ {
			multiSigScript = append(multiSigScript, signers[i].PublicKey)
		}
	}
	multiSigScript = append(multiSigScript, []interface{}{fmt.Sprint("OP_", sumWeight), "OP_CHECKMULTISIG"}...)
	script := &scripts.Script{Script: multiSigScript}
	switch addressType {
	case address.P2WSH:
		{
			addr, err := address.P2WSHAddresssFromScript(script)
			if err != nil {
				return nil, err
			}
			return &MultiSignaturAddress{
				Signers:       signers,
				Threshold:     threshold,
				Address:       addr,
				ScriptDetails: script.ToHex(),
			}, nil

		}
	case address.P2WSHInP2SH:
		{
			p2wsh, err := address.P2WSHAddresssFromScript(script)
			if err != nil {
				return nil, err
			}
			addr, err := address.P2SHAddressFromScript(p2wsh.ToScriptPubKey(), address.P2WSHInP2SH)
			if err != nil {
				return nil, err
			}
			return &MultiSignaturAddress{
				Signers:       signers,
				Threshold:     threshold,
				Address:       addr,
				ScriptDetails: script.ToHex(),
			}, nil
		}
	default:
		{
			return nil, fmt.Errorf("addressType should be p2wsh or P2WSHInP2SH")
		}
	}

}

// ShowScript returns a representation of the multi-signature script details associated with
// the MultiSignaturAddress in a human-readable format. It provides a list of string elements
// describing each part of the script.
//
// Returns:
// - []string: A slice of strings representing the script details in a human-readable format.
func (addr *MultiSignaturAddress) ShowScript() []string {
	sumWeight := addr.Signers.SumThreshHold()
	multiSigScript := []string{fmt.Sprint("OP_", addr.Threshold)}
	for i := 0; i < len(addr.Signers); i++ {
		for w := 0; w < addr.Signers[i].Weight(); w++ {
			multiSigScript = append(multiSigScript, addr.Signers[i].PublicKey())
		}
	}
	multiSigScript = append(multiSigScript, []string{fmt.Sprint("OP_", sumWeight), "OP_CHECKMULTISIG"}...)
	return multiSigScript
}

func (owners MultiSignaturAddressSigners) SumThreshHold() int {
	sumWeight := int(0)
	for i := 0; i < len(owners); i++ {
		sumWeight += owners[i].Weight()
	}
	return sumWeight

}
