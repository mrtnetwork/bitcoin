package builder

import (
	"bitcoin/address"
	"bitcoin/scripts"
	"bitcoin/tools"
	"errors"
	"fmt"
)

type MultiSignatureSigner struct {
	// publicKey of signer
	PublicKey string
	// Participant's weight
	Weight int
}
type MultiSignaturAddress struct {
	Signers MultiSignaturAddressSigners
	// "m" is the minimum number of signatures required to spend the bitcoins.
	Threshold int
	// signers
	Address address.BitcoinAddress
	// N = sum of Signers weight
	// Script details ["OP_M", ...comperesedPublicKeys , "OP_N", "OP_CHECKMULTISIG"]
	ScriptDetails string
}
type MultiSignaturAddressSigners []MultiSignatureSigner

func CreateMultiSignatureAddress(threshold int, signers MultiSignaturAddressSigners, addressType address.AddressType) (*MultiSignaturAddress, error) {
	for i := 0; i < len(signers); i++ {
		isValid := tools.IsValidStruct(signers[i])
		if !isValid {
			return nil, errors.New("invalid MultiSignaturAddressOwnerDetals struct")
		}
	}
	if threshold > 16 || threshold < 1 {
		return nil, errors.New("the threshold should not be greater than 16 and less than 1")
	}
	sumWeight := signers.SumThreshHold()
	if threshold > 16 || threshold < 1 {
		return nil, errors.New("the total weight of the owners should not be more than 16")
	}
	if sumWeight < threshold {
		return nil, errors.New("the total weight of the signatories should reach the threshold")
	}
	multiSigScript := []interface{}{fmt.Sprint("OP_", threshold)}
	for i := 0; i < len(signers); i++ {
		for w := 0; w < signers[i].Weight; w++ {
			multiSigScript = append(multiSigScript, signers[i].PublicKey)
		}
	}
	multiSigScript = append(multiSigScript, []interface{}{fmt.Sprint("OP_", sumWeight), "OP_CHECKMULTISIG"}...)
	script := scripts.Script{Script: multiSigScript}
	var addr address.BitcoinAddress
	switch addressType {
	case address.P2WSH:
		{
			addr = address.P2WSHAddresssFromScript(script)
		}
	case address.P2WSHInP2SH:
		{
			p2wsh := address.P2WSHAddresssFromScript(script)
			addr = address.P2SHAddressFromScript(p2wsh.ToScriptPubKey(), address.P2WSHInP2SH)
		}
	default:
		{
			return nil, errors.New("addressType should be p2wsh or P2WSHInP2SH")
		}
	}
	return &MultiSignaturAddress{
		Signers:       signers,
		Threshold:     threshold,
		Address:       addr,
		ScriptDetails: script.ToHex(),
	}, nil
}
func (addr *MultiSignaturAddress) ShowScript() []string {
	sumWeight := addr.Signers.SumThreshHold()
	multiSigScript := []string{fmt.Sprint("OP_", addr.Threshold)}
	for i := 0; i < len(addr.Signers); i++ {
		for w := 0; w < addr.Signers[i].Weight; w++ {
			multiSigScript = append(multiSigScript, addr.Signers[i].PublicKey)
		}
	}
	multiSigScript = append(multiSigScript, []string{fmt.Sprint("OP_", sumWeight), "OP_CHECKMULTISIG"}...)
	return multiSigScript
}

func (owners MultiSignaturAddressSigners) SumThreshHold() int {
	sumWeight := int(0)
	for i := 0; i < len(owners); i++ {
		sumWeight += owners[i].Weight
	}
	return sumWeight

}
