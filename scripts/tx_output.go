package scripts

import (
	"encoding/binary"
	"math/big"

	"github.com/mrtnetwork/bitcoin/formating"
)

type TxOutput struct {
	// the value we want to send to this output in satoshis
	Amount *big.Int
	// the script that will lock this amount
	ScriptPubKey *Script
}

// NewTxOutput creates a new transaction output with the specified amount and script
// public key, and returns a pointer to the initialized TxOutput object.
func NewTxOutput(amount *big.Int, scriptPubKey *Script) *TxOutput {
	return &TxOutput{
		Amount:       amount,
		ScriptPubKey: scriptPubKey,
	}
}

// Create a copy of the object
func (txOutput *TxOutput) Copy() *TxOutput {

	return NewTxOutput(new(big.Int).Set(txOutput.Amount), &Script{txOutput.ScriptPubKey.Script})
}

// Serialize TxOutput to bytes
func (txOutput *TxOutput) ToBytes() []byte {

	amountBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(amountBytes, uint64(txOutput.Amount.Int64()))
	scriptBytes := txOutput.ScriptPubKey.ToBytes()
	scriptLengthVarint := formating.EncodeVarint(len(scriptBytes))
	data := append(amountBytes, append(scriptLengthVarint, scriptBytes...)...)
	return data
}

// raw The hexadecimal raw string of the Transaction
// The cursor of which the algorithm will start to read the data
// hasSegwit  Is the Tx Output segwit or not
func TxOutputFromRaw(outputBytes []byte, cursor int, hasSegwit bool) (*TxOutput, int, error) {

	// Parse TxOutput from raw bytes
	value := int64(binary.LittleEndian.Uint64(outputBytes[cursor : cursor+8]))
	cursor += 8

	vi, viSize := formating.ViToInt(outputBytes[cursor:])
	cursor += viSize

	lockScript := outputBytes[cursor : cursor+vi]
	cursor += vi

	scriptPubKey, err := ScriptFromRaw(lockScript, hasSegwit)
	if err != nil {
		return nil, cursor, err
	}
	return NewTxOutput(big.NewInt(value), scriptPubKey), cursor, nil
}
