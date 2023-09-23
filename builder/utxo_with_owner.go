package builder

import (
	"bitcoin/address"
	"bitcoin/keypair"
	"math/big"
)

type UtxoOwnerDetails struct {
	PublicKey       string
	Address         address.BitcoinAddress
	MultiSigAddress *MultiSignaturAddress
}

type UtxoWithOwner struct {
	Utxo         BitcoinUtxo
	OwnerDetails UtxoOwnerDetails
	// MultiSigAddress *MultiSignaturAddress
}
type UtxoWithOwnerList []UtxoWithOwner

func NewUtxoWithOwner(utxo BitcoinUtxo, owner UtxoOwnerDetails) *UtxoWithOwner {
	return &UtxoWithOwner{
		Utxo:         utxo,
		OwnerDetails: owner,
		// MultiSigAddress: owner.MultiSigAddress,
	}
}

func (utxo *UtxoWithOwner) Public() keypair.ECPublic {
	return *keypair.NewECPPublicFromHex(utxo.OwnerDetails.PublicKey)
}
func (utxo *UtxoWithOwner) IsMultiSig() bool {

	return utxo.OwnerDetails.MultiSigAddress != nil
}

func (utxos UtxoWithOwnerList) SumOfUtxosValue() *big.Int {
	sum := big.NewInt(0)
	for _, utxo := range utxos {
		sum.Add(sum, utxo.Utxo.value)
	}
	return sum
}

func (utxos UtxoWithOwnerList) CanSpending() bool {
	value := utxos.SumOfUtxosValue()
	return value.Cmp(big.NewInt(0)) != 0
}
