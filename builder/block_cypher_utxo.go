package builder

import (
	"math/big"
	"time"
)

type TxRef struct {
	TxHash        string    `json:"tx_hash"`
	BlockHeight   int       `json:"block_height"`
	TxInputN      int       `json:"tx_input_n"`
	TxOutputN     int       `json:"tx_output_n"`
	Value         big.Int   `json:"value"`
	RefBalance    int       `json:"ref_balance"`
	Spent         bool      `json:"spent"`
	Confirmations int       `json:"confirmations"`
	Confirmed     time.Time `json:"confirmed"`
	DoubleSpend   bool      `json:"double_spend"`
	Script        string    `json:"script"`
}

type blockCypherUtxo struct {
	// UtxoOwnerDetails   UtxoOwnerDetails
	Address            string  `json:"address"`
	TotalReceived      int     `json:"total_received"`
	TotalSent          int     `json:"total_sent"`
	Balance            int     `json:"balance"`
	UnconfirmedBalance int     `json:"unconfirmed_balance"`
	FinalBalance       int     `json:"final_balance"`
	NTx                int     `json:"n_tx"`
	UnconfirmedNTx     int     `json:"unconfirmed_n_tx"`
	FinalNTx           int     `json:"final_n_tx"`
	TxRefs             []TxRef `json:"txrefs"`
	TxURL              string  `json:"tx_url"`
}

func (info *blockCypherUtxo) ToUtxoWithOwner(owner UtxoOwnerDetails) UtxoWithOwnerList {
	utxos := make([]UtxoWithOwner, len(info.TxRefs))
	for i := 0; i < len(info.TxRefs); i++ {
		utxos[i] = UtxoWithOwner{
			Utxo: BitcoinUtxo{
				txHash:      info.TxRefs[i].TxHash,
				value:       &info.TxRefs[i].Value,
				vout:        info.TxRefs[i].TxOutputN,
				scriptType:  owner.Address.GetType(),
				blockHeight: 1,
			},
			OwnerDetails: owner,
		}
	}
	return utxos
}

// func (utxos *UtxoWithOwnerList) CanSpending() bool {
// 	value := utxos.SumOfUtxosValue()
// 	return value.Cmp(big.NewInt(0)) != 0
// }
