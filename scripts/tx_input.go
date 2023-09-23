package scripts

import (
	"bitcoin/constant"
	"bitcoin/formating"
	"encoding/binary"
	"encoding/hex"
)

// TxInput represents a transaction input
type TxInput struct {
	// the transaction id as a hex string (little-endian as displayed by
	// tools)
	TxID string
	// the index of the UTXO that we want to spend
	TxIndex int
	// the script that satisfies the locking conditions (aka unlocking script)
	ScriptSig Script
	// the input sequence (for timelocks, RBF, etc.)
	Sequence []byte
}

// NewTxInput creates a new transaction input with the provided transaction ID and index,
// and additional optional arguments for the script and sequence. It returns a pointer
// to the initialized TxInput object.
func NewTxInput(txID string, txIndex int, options ...interface{}) *TxInput {
	script := Script{Script: []interface{}{}}
	sequance := constant.DEFAULT_TX_SEQUENCE
	for _, opt := range options {
		switch v := opt.(type) {
		case Script:
			script = v
		case []byte:
			sequance = v
		default:
			panic("invalid Tx Input argruments")
		}
	}
	return &TxInput{
		TxID:      txID,
		TxIndex:   txIndex,
		ScriptSig: script,
		Sequence:  sequance,
	}
}

// NewDefaultTxInput creates a new default transaction input with the provided transaction ID
// and index, and returns a pointer to the initialized TxInput object.
func NewDefaultTxInput(txID string, txIndex int) *TxInput {
	return &TxInput{
		TxID:      txID,
		TxIndex:   txIndex,
		ScriptSig: Script{},
		Sequence:  constant.DEFAULT_TX_SEQUENCE,
	}
}

// Copy creates a copy of the TxInput
func (ti *TxInput) Copy() *TxInput {
	return NewTxInput(ti.TxID, ti.TxIndex, Script{ti.ScriptSig.Script}, ti.Sequence)
}

// Serialize serializes the TxInput to bytes
func (ti *TxInput) ToBytes() []byte {
	txidBytes, _ := hex.DecodeString(ti.TxID)
	txidBytes = formating.ReverseBytes(txidBytes)
	txoutBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(txoutBytes, uint32(ti.TxIndex))

	scriptSigBytes := ti.ScriptSig.ToBytes()
	scriptSigLengthVarint := formating.EncodeVarint(len(scriptSigBytes))

	data := append(txidBytes, txoutBytes...)
	data = append(data, scriptSigLengthVarint...)
	data = append(data, scriptSigBytes...)
	data = append(data, ti.Sequence...)

	return data
}

// FromRaw parses a raw transaction input string into a TxInput
func TxInputFromRaw(raw string, cursor int, hasSegwit bool) (*TxInput, int) {
	txInputRaw, err := hex.DecodeString(raw)
	if err != nil {
		panic("invalid tx input hex")
	}

	if cursor+32 >= len(txInputRaw) {
		panic("Input transaction hash not found. Probably malformed raw transaction")
	}

	inpHash := make([]byte, 32)
	copy(inpHash, formating.ReverseBytes(txInputRaw[cursor:cursor+32]))
	cursor += 32

	if cursor+4 >= len(txInputRaw) {
		panic("Output number not found. Probably malformed raw transaction")
	}

	outputN := binary.LittleEndian.Uint32(formating.ReverseBytes(txInputRaw[cursor : cursor+4]))
	cursor += 4

	vi, viSize := formating.ViToInt(txInputRaw[cursor:])
	cursor += viSize

	if cursor+vi > len(txInputRaw) {
		panic("Unlocking script length exceeds available data. Probably malformed raw transaction")
	}

	unlockingScript := txInputRaw[cursor : cursor+vi]
	cursor += vi

	if cursor+4 > len(txInputRaw) {
		panic("Sequence number not found. Probably malformed raw transaction")
	}

	sequenceNumberData := txInputRaw[cursor : cursor+4]
	cursor += 4

	script := ScriptFromRaw(formating.BytesToHex(unlockingScript), hasSegwit)

	return NewTxInput(formating.BytesToHex(inpHash), int(outputN), *script, sequenceNumberData), cursor
}
