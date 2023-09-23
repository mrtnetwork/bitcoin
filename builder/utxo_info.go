package builder

import (
	"bitcoin/address"
	"math/big"
)

type BitcoinUtxo struct {
	txHash      string
	value       *big.Int
	vout        int
	scriptType  address.AddressType
	blockHeight int
}

func NewBitcoinUtxo(txHash string, value *big.Int, vout int, scriptType address.AddressType) *BitcoinUtxo {
	return &BitcoinUtxo{
		txHash:      txHash,
		value:       value,
		vout:        vout,
		scriptType:  scriptType,
		blockHeight: 1,
	}
}

func NewBitcoinUtxoFromJSONBlockCypher(json map[string]interface{}, addressType address.AddressType) *BitcoinUtxo {
	txHash := json["tx_hash"].(string)
	vout := int(json["tx_output_n"].(float64))
	value := new(big.Int)
	valueStr := json["value"].(string)
	value.SetString(valueStr, 10)

	utxo := &BitcoinUtxo{
		txHash:      txHash,
		vout:        vout,
		value:       value,
		scriptType:  addressType,
		blockHeight: 0,
	}

	// scriptTypeStr := json["scriptType"].(string)
	// if scriptTypeStr != "" {
	// 	utxo.scriptType = BitcoinAddressTypeFromString(scriptTypeStr)
	// 	return utxo
	// }

	// scriptDataHex := json["script"].(string)
	// scriptType := scripts.GetScriptType(scriptDataHex)
	// if scriptType == nil {
	// 	utxo.scriptType = addressType
	// 	return utxo
	// }

	// if *scriptType == address.P2PK {
	// 	utxo.scriptType = address.P2PK
	// 	return utxo
	// }

	utxo.scriptType = addressType
	return utxo
}

func NewBitcoinUtxoFromMempool(json map[string]interface{}, scriptType address.AddressType) *BitcoinUtxo {
	txHash := json["txid"].(string)
	vout := int(json["vout"].(float64))
	value := new(big.Int)
	valueStr := json["value"].(string)
	value.SetString(valueStr, 10)

	return &BitcoinUtxo{
		txHash:      txHash,
		vout:        vout,
		value:       value,
		scriptType:  scriptType,
		blockHeight: int(json["status"].(map[string]interface{})["block_height"].(float64)),
	}
}

func (utxo *BitcoinUtxo) IsP2tr() bool {
	return utxo.scriptType == address.P2TR
}

func (utxo *BitcoinUtxo) IsSegwit() bool {
	return utxo.scriptType == address.P2WPKH ||
		utxo.scriptType == address.P2WSH ||
		utxo.scriptType == address.P2TR ||
		utxo.IsP2shSegwit()
}

func (utxo *BitcoinUtxo) IsP2shSegwit() bool {
	return utxo.scriptType == address.P2WPKHInP2SH ||
		utxo.scriptType == address.P2WSHInP2SH
}

func (utxo *BitcoinUtxo) IsSupported() bool {
	return utxo.scriptType != address.AddressType(-1)
}

func (utxo *BitcoinUtxo) CanUsedForTransaction() bool {
	return utxo.IsSupported() && utxo.Comfirmed()
}

func (utxo *BitcoinUtxo) Comfirmed() bool {
	return utxo.blockHeight > 0
}
