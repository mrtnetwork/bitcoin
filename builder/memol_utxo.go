package builder

import "math/big"

type MempolUtxo struct {
	Txid   string `json:"txid"`
	Vout   int    `json:"vout"`
	Status struct {
		Confirmed   bool   `json:"confirmed"`
		BlockHeight int    `json:"block_height"`
		BlockHash   string `json:"block_hash"`
		BlockTime   int    `json:"block_time"`
	} `json:"status"`
	Value big.Int `json:"value"`
}
type MempolUtxoList []MempolUtxo

// type MempolUtxoAddress struct {
// 	UtxoOwnerDetails UtxoOwnerDetails
// 	Utxos            []MempolUtxo
// }

func (memplUtxos MempolUtxoList) ToUtxoWithOwner(owner UtxoOwnerDetails) UtxoWithOwnerList {
	utxos := make([]UtxoWithOwner, len(memplUtxos))
	for i := 0; i < len(memplUtxos); i++ {
		utxos[i] = UtxoWithOwner{
			Utxo: BitcoinUtxo{
				txHash:      memplUtxos[i].Txid,
				value:       &memplUtxos[i].Value,
				vout:        memplUtxos[i].Vout,
				scriptType:  owner.Address.GetType(),
				blockHeight: 1,
			},
			OwnerDetails: owner,
		}
	}
	return utxos
}
