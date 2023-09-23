package provider

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
				TxHash:      info.TxRefs[i].TxHash,
				Value:       &info.TxRefs[i].Value,
				Vout:        info.TxRefs[i].TxOutputN,
				ScriptType:  owner.Address.GetType(),
				BlockHeight: 1,
			},
			OwnerDetails: owner,
		}
	}
	return utxos
}

// blockcyper transaction details
type BlocCypherTransactionInput struct {
	PrevHash    string   `json:"prev_hash"`
	OutputIndex int      `json:"output_index"`
	OutputValue int      `json:"output_value"`
	Sequence    int      `json:"sequence"`
	Addresses   []string `json:"addresses"`
	ScriptType  string   `json:"script_type"`
	Age         int      `json:"age"`
	Witness     []string `json:"witness"`
}

type BlocCyperTransactionOutput struct {
	Value      int      `json:"value"`
	Script     string   `json:"script"`
	Addresses  []string `json:"addresses"`
	ScriptType string   `json:"script_type"`
	DataHex    string   `json:"data_hex"`
	DataString string   `json:"data_string"`
}

type BlocCyperTransaction struct {
	BlockHeight   int                          `json:"block_height"`
	BlockIndex    int                          `json:"block_index"`
	Hash          string                       `json:"hash"`
	Addresses     []string                     `json:"addresses"`
	Total         int                          `json:"total"`
	Fees          int                          `json:"fees"`
	Size          int                          `json:"size"`
	VSize         int                          `json:"vsize"`
	Preference    string                       `json:"preference"`
	RelayedBy     string                       `json:"relayed_by"`
	Received      time.Time                    `json:"received"`
	Ver           int                          `json:"ver"`
	DoubleSpend   bool                         `json:"double_spend"`
	VinSz         int                          `json:"vin_sz"`
	VoutSz        int                          `json:"vout_sz"`
	OptInRBF      bool                         `json:"opt_in_rbf"`
	DataProtocol  string                       `json:"data_protocol"`
	Confirmations int                          `json:"confirmations"`
	Inputs        []BlocCypherTransactionInput `json:"inputs"`
	Outputs       []BlocCyperTransactionOutput `json:"outputs"`
}
type BlockCypherAddressInfo struct {
	Address            string                     `json:"address"`
	TotalReceived      int64                      `json:"total_received"`
	TotalSent          int64                      `json:"total_sent"`
	Balance            int64                      `json:"balance"`
	UnconfirmedBalance int64                      `json:"unconfirmed_balance"`
	FinalBalance       int64                      `json:"final_balance"`
	NumTransactions    int                        `json:"n_tx"`
	UnconfirmedNumTx   int                        `json:"unconfirmed_n_tx"`
	FinalNumTx         int                        `json:"final_n_tx"`
	TXs                BlockCypherTransactionList `json:"txs"`
}

type BlockCypherTransactionList []BlocCyperTransaction
