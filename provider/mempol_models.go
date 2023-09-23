package provider

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

func (memplUtxos MempolUtxoList) ToUtxoWithOwner(owner UtxoOwnerDetails) UtxoWithOwnerList {
	utxos := make([]UtxoWithOwner, len(memplUtxos))
	for i := 0; i < len(memplUtxos); i++ {
		utxos[i] = UtxoWithOwner{
			Utxo: BitcoinUtxo{
				TxHash:      memplUtxos[i].Txid,
				Value:       &memplUtxos[i].Value,
				Vout:        memplUtxos[i].Vout,
				ScriptType:  owner.Address.GetType(),
				BlockHeight: 1,
			},
			OwnerDetails: owner,
		}
	}
	return utxos
}

// mempool transaction details
type MempoolPrevOut struct {
	ScriptPubKey        string `json:"scriptpubkey"`
	ScriptPubKeyAsm     string `json:"scriptpubkey_asm"`
	ScriptPubKeyType    string `json:"scriptpubkey_type"`
	ScriptPubKeyAddress string `json:"scriptpubkey_address"`
	Value               int    `json:"value"`
}

type MempoolVin struct {
	TxID         string         `json:"txid"`
	Vout         int            `json:"vout"`
	PrevOut      MempoolPrevOut `json:"prevout"`
	ScriptSig    string         `json:"scriptsig"`
	ScriptSigAsm string         `json:"scriptsig_asm"`
	Witness      []string       `json:"witness"`
	IsCoinbase   bool           `json:"is_coinbase"`
	Sequence     int            `json:"sequence"`
}

type MempoolVout struct {
	ScriptPubKey        string `json:"scriptpubkey"`
	ScriptPubKeyAsm     string `json:"scriptpubkey_asm"`
	ScriptPubKeyType    string `json:"scriptpubkey_type"`
	ScriptPubKeyAddress string `json:"scriptpubkey_address"`
	Value               int    `json:"value"`
}

type MempoolStatus struct {
	Confirmed   bool   `json:"confirmed"`
	BlockHeight int    `json:"block_height"`
	BlockHash   string `json:"block_hash"`
	BlockTime   int64  `json:"block_time"`
}

type MempoolTransaction struct {
	TxID     string        `json:"txid"`
	Version  int           `json:"version"`
	Locktime int           `json:"locktime"`
	Vin      []MempoolVin  `json:"vin"`
	Vout     []MempoolVout `json:"vout"`
	Size     int           `json:"size"`
	Weight   int           `json:"weight"`
	Fee      int           `json:"fee"`
	Status   MempoolStatus `json:"status"`
}

type MemoolTransactionList []MempoolTransaction
