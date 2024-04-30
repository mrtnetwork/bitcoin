package scripts

import (
	"math/big"

	"github.com/mrtnetwork/bitcoin/constant"
	"github.com/mrtnetwork/bitcoin/digest"
	"github.com/mrtnetwork/bitcoin/formating"
)

// Define BtcTransaction struct
type BtcTransaction struct {
	//  A list of all the transaction inputs
	Inputs []*TxInput
	// A list of all the transaction outputs
	Outputs []*TxOutput
	// The transaction's locktime parameter
	Locktime []byte
	// The transaction version
	Version []byte
	// Specifies a tx that includes segwit inputs
	HasSegwit bool
	// The witness structure that corresponds to the inputs
	Witnesses []*TxWitnessInput
}

// NewBtcTransaction creates a new Bitcoin transaction with the specified inputs, outputs,
// and optional parameters such as locktime, version, witness inputs, and SegWit flag.
// It returns a pointer to the initialized BtcTransaction object.
func NewBtcTransaction(inputs []*TxInput, outputs []*TxOutput, hasSegwit bool, options ...interface{}) *BtcTransaction {
	var lock []byte
	var version []byte
	w := make([]*TxWitnessInput, 0)
	for _, opt := range options {
		switch v := opt.(type) {
		case []byte:
			if lock == nil {
				lock = v
			} else if version == nil {
				version = v
			}
		case []TxWitnessInput:
			for _, input := range v {
				w = append(w, &input)
			}
		case []*TxWitnessInput:
			w = append(w, v...)

		}
	}
	if lock == nil {
		lock = constant.DEFAULT_TX_LOCKTIME
	}
	if version == nil {
		version = constant.DEFAULT_TX_VERSION
	}

	transaction := &BtcTransaction{
		Inputs:    inputs,
		Outputs:   outputs,
		Witnesses: w,
		Locktime:  lock,
		Version:   version,
		HasSegwit: hasSegwit,
	}

	transaction.Witnesses = append(transaction.Witnesses, w...)

	return transaction
}

// Deep copy of Transaction
func (tx *BtcTransaction) Copy() *BtcTransaction {
	inputsCopy := make([]*TxInput, len(tx.Inputs))
	for i, input := range tx.Inputs {
		inputsCopy[i] = input.Copy()
	}

	outputsCopy := make([]*TxOutput, len(tx.Outputs))
	for i, output := range tx.Outputs {
		outputsCopy[i] = output.Copy()
	}

	witnessesCopy := make([]*TxWitnessInput, len(tx.Witnesses))
	for i, witness := range tx.Witnesses {
		witnessesCopy[i] = witness.Copy()
	}

	return &BtcTransaction{
		HasSegwit: tx.HasSegwit,
		Inputs:    inputsCopy,
		Outputs:   outputsCopy,
		Witnesses: witnessesCopy,
		Locktime:  tx.Locktime,
		Version:   tx.Version,
	}
}
func BtcTransactionFromRaw(raw string) (*BtcTransaction, error) {
	txBytes := formating.HexToBytes(raw)
	cursor := 4
	var flag []byte
	hasSegwit := false

	if txBytes[4] == 0 {
		flag = txBytes[5:6]
		if flag[0] == 1 {
			hasSegwit = true
		}
		cursor += 2
	}

	vi, viCursor := formating.ViToInt(txBytes[cursor:])
	cursor += viCursor

	inputs := make([]*TxInput, vi)
	for index := 0; index < len(inputs); index++ {
		inp, inpCursor, err := TxInputFromRaw(txBytes, cursor, hasSegwit)
		if err != nil {
			return nil, err
		}
		inputs[index] = inp
		cursor = inpCursor
	}
	viOut, viOutCursor := formating.ViToInt(txBytes[cursor:])
	cursor += viOutCursor

	outputs := make([]*TxOutput, viOut)
	for index := 0; index < len(outputs); index++ {
		out, outCursor, err := TxOutputFromRaw(txBytes, cursor, hasSegwit)
		if err != nil {
			return nil, err
		}
		outputs[index] = out
		cursor = outCursor
	}
	witnesses := make([]TxWitnessInput, len(inputs))
	if hasSegwit {
		for n := 0; n < len(inputs); n++ {
			wVi, wViCursor := formating.ViToInt(txBytes[cursor:])
			cursor += wViCursor
			witnessesTmp := make([]string, wVi)
			for m := 0; m < len(witnessesTmp); m++ {
				var witness []byte
				wtVi, wtViCursor := formating.ViToInt(txBytes[cursor:])
				if wtVi != 0 {
					witness = txBytes[cursor+wtViCursor : cursor+wtViCursor+wtVi]
				}
				cursor += wtViCursor + wtVi
				witnessesTmp[m] = formating.BytesToHex(witness)
			}
			witnesses[n] = TxWitnessInput{Stack: witnessesTmp}
		}
	}
	return NewBtcTransaction(inputs, outputs, hasSegwit, witnesses), nil

}

// Define a BtcTransaction method for converting it to bytes
func (tx *BtcTransaction) ToBytes(segwit bool) []byte {
	var data []byte

	// Add the version bytes to the data
	data = append(data, tx.Version...)

	if segwit {
		// If segwit is enabled, add the segwit marker and flag bytes
		data = append(data, 0x00, 0x01)
	}
	// Encode the input count and add it to the data
	txInCountBytes := formating.EncodeVarint(int(len(tx.Inputs)))
	data = append(data, txInCountBytes...)

	// Add serialized input data for each input
	for _, txIn := range tx.Inputs {
		data = append(data, txIn.ToBytes()...)
	}

	// Encode the output count and add it to the data
	txOutCountBytes := formating.EncodeVarint(int(len(tx.Outputs)))
	data = append(data, txOutCountBytes...)

	// Add serialized output data for each output
	for _, txOut := range tx.Outputs {
		tBytes := txOut.ToBytes()
		data = append(data, tBytes...)
	}

	if segwit {
		// If segwit is enabled, add witness data
		for _, witness := range tx.Witnesses {
			witnessCountBytes := formating.EncodeVarint(int(len(witness.Stack)))
			data = append(data, witnessCountBytes...)
			data = append(data, witness.ToBytes()...)
		}
	}

	// Add the locktime bytes to the data
	data = append(data, tx.Locktime...)

	return data
}
func getSigHashArgruments(defaultValue int, args ...interface{}) int {
	sigHash := defaultValue
	for _, opt := range args {
		switch v := opt.(type) {
		case int:
			sigHash = v
			return sigHash
		}
	}
	return sigHash
}

// Returns the transaction's digest for signing.
// https://en.bitcoin.it/wiki/OP_CHECKSIG
// |  SIGHASH types (see constants.py):
// |      SIGHASH_ALL - signs all inputs and outputs (default)
// |      SIGHASH_NONE - signs all of the inputs
// |      SIGHASH_SINGLE - signs all inputs but only txin_index output
// |      SIGHASH_ANYONECANPAY (only combined with one of the above)
// |      - with ALL - signs all outputs but only txin_index input
// |      - with NONE - signs only the txin_index input
// |      - with SINGLE - signs txin_index input and output
// txInIndex : The index of the input that we wish to sign
// script : The scriptPubKey of the UTXO that we want to spend
// sighash : The type of the signature hash to be created

func (tx *BtcTransaction) GetTransactionDigest(txInIndex int, script *Script, sighash ...interface{}) []byte {
	sig := getSigHashArgruments(constant.SIGHASH_ALL, sighash...)
	// Make a copy of the transaction
	txCopy := tx.Copy()

	// Set scriptSig for all inputs except the specified one
	for i := range txCopy.Inputs {
		txCopy.Inputs[i].ScriptSig = NewScript()
	}
	txCopy.Inputs[txInIndex].ScriptSig = script

	switch sig & 0x1f {
	case constant.SIGHASH_NONE:
		// Clear outputs for SIGHASH_NONE
		txCopy.Outputs = nil

		// Set empty sequence for all inputs except the specified one
		for i := range txCopy.Inputs {
			if i != txInIndex {
				txCopy.Inputs[i].Sequence = constant.EMPTY_TX_SEQUENCE
			}
		}

	case constant.SIGHASH_SINGLE:
		if txInIndex >= len(txCopy.Outputs) {
			panic("Transaction index is greater than the available outputs")
		}

		// Clear outputs except for the specified one
		txCopy.Outputs = txCopy.Outputs[:txInIndex+1]

		// Fill in with dummy outputs
		for i := 0; i < txInIndex; i++ {
			txCopy.Outputs[i] = &TxOutput{
				Amount:       big.NewInt(constant.NEGATIVE_SATOSHI),
				ScriptPubKey: NewScript(),
			}

		}

		// Set empty sequence for all inputs except the specified one
		for i := range txCopy.Inputs {
			if i != txInIndex {
				txCopy.Inputs[i].Sequence = constant.EMPTY_TX_SEQUENCE
			}
		}
	}

	if sig&constant.SIGHASH_ANYONECANPAY != 0 {
		input := txCopy.Inputs[txInIndex]
		txCopy.Inputs = []*TxInput{input}
	}

	// Serialize the transaction without SegWit data
	txForSign := txCopy.ToBytes(false)

	// Pack the sighash and append it to the serialized transaction
	packedData := formating.PackInt32LE(sig)
	txForSign = append(txForSign, packedData...)

	// Calculate the double hash of the serialized transaction
	return digest.DoubleHash(txForSign)
}

// Returns the segwit v0 transaction's digest for signing.
//https://github.com/github.com/mrtnetwork/bitcoin/bips/blob/master/bip-0143.mediawiki

// |  SIGHASH types (see constants.py):
// |      SIGHASH_ALL - signs all inputs and outputs (default)
// |      SIGHASH_NONE - signs all of the inputs
// |      SIGHASH_SINGLE - signs all inputs but only txin_index output
// |      SIGHASH_ANYONECANPAY (only combined with one of the above)
// |      - with ALL - signs all outputs but only txin_index input
// |      - with NONE - signs only the txin_index input
// |      - with SINGLE - signs txin_index input and output
// txin_index : The index of the input that we wish to sign
// script : The scriptCode (template) that corresponds to the segwit
// transaction output type that we want to spend
// amount : The amount of the UTXO to spend is included in the
// signature for segwit (in satoshis)
// sighash : The type of the signature hash to be created

func (tx *BtcTransaction) GetTransactionSegwitDigit(txInIndex int, script *Script, amount *big.Int, sigshash ...interface{}) []byte {
	sig := getSigHashArgruments(constant.SIGHASH_ALL, sigshash...)
	// Make a copy of the transaction
	txCopy := tx.Copy()

	// Initialize hashPrevouts, hashSequence, and hashOutputs
	hashPrevouts := make([]byte, 32)
	hashSequence := make([]byte, 32)
	hashOutputs := make([]byte, 32)

	// Extract basicSigHashType and flags
	basicSigHashType := sig & 0x1F
	anyoneCanPay := (sig & 0xF0) == constant.SIGHASH_ANYONECANPAY
	signAll := (basicSigHashType != constant.SIGHASH_SINGLE) && (basicSigHashType != constant.SIGHASH_NONE)
	if !anyoneCanPay {
		hashPrevouts = []byte{}
		for _, txin := range txCopy.Inputs {
			txidBytes := formating.ReverseBytes(formating.HexToBytes(txin.TxID))
			txoutIndexBytes := formating.PackUint32LE(uint32(txin.TxIndex))
			hashPrevouts = append(hashPrevouts, txidBytes...)
			hashPrevouts = append(hashPrevouts, txoutIndexBytes...)
		}
		hashPrevouts = digest.DoubleHash(hashPrevouts)

	}

	if !anyoneCanPay && signAll {

		// Calculate hashSequence
		hashSequence = []byte{}
		for _, txin := range txCopy.Inputs {
			hashSequence = append(hashSequence, txin.Sequence...)
		}
		hashSequence = digest.DoubleHash(hashSequence)
	}

	if signAll {
		// Calculate hashOutputs
		hashOutputs = []byte{}
		for _, txout := range txCopy.Outputs {
			amountBytes := formating.PackBigIntToLittleEndian(txout.Amount)
			scriptBytes := txout.ScriptPubKey.ToBytes()
			outputLength := byte(len(scriptBytes))

			// Concatenate amountBytes, outputLength, and scriptBytes to hashOutputs
			hashOutputs = append(hashOutputs, amountBytes...)
			hashOutputs = append(hashOutputs, outputLength)
			hashOutputs = append(hashOutputs, scriptBytes...)
		}
		hashOutputs = digest.DoubleHash(hashOutputs)

	} else if basicSigHashType == constant.SIGHASH_SINGLE && txInIndex < len(txCopy.Outputs) {
		out := txCopy.Outputs[txInIndex]
		packedAmount := formating.PackBigIntToLittleEndian(out.Amount)
		scriptBytes := out.ScriptPubKey.ToBytes()
		lenScriptBytes := []byte{byte(len(scriptBytes))}
		hashOutputs = append(packedAmount, lenScriptBytes...)
		hashOutputs = append(hashOutputs, scriptBytes...)
		hashOutputs = digest.DoubleHash(hashOutputs)

	}

	// Create a byte slice to assemble the transaction data
	var txForSigning []byte

	// Add version, hashPrevouts, and hashSequence to txForSigning
	txForSigning = append(txForSigning, txCopy.Version...)
	txForSigning = append(txForSigning, hashPrevouts...)
	txForSigning = append(txForSigning, hashSequence...)
	// Add the relevant input 	data for the specified input
	txIn := txCopy.Inputs[txInIndex]
	txidBytes := formating.ReverseBytes(formating.HexToBytes(txIn.TxID))
	txoutIndexBytes := formating.PackUint32LE(uint32(txIn.TxIndex))
	txForSigning = append(txForSigning, append(txidBytes, txoutIndexBytes...)...)
	txForSigning = append(txForSigning, byte(len(script.ToBytes())))
	txForSigning = append(txForSigning, script.ToBytes()...)
	packedAmount := formating.PackBigIntToLittleEndian(amount)
	txForSigning = append(txForSigning, packedAmount...)
	txForSigning = append(txForSigning, txIn.Sequence...)
	txForSigning = append(txForSigning, hashOutputs...)
	txForSigning = append(txForSigning, txCopy.Locktime...)
	packedSighash := formating.PackInt32LE(sig)
	txForSigning = append(txForSigning, packedSighash...)

	// Calculate the double hash of txForSigning
	return digest.DoubleHash(txForSigning)
}

// Returns the segwit v1 (taproot) transaction's digest for signing.
// https://github.com/github.com/mrtnetwork/bitcoin/bips/blob/master/bip-0341.mediawiki
// Also consult Bitcoin Core code at: https://github.com/github.com/mrtnetwork/bitcoin/github.com/mrtnetwork/bitcoin/blob/29c36f070618ea5148cd4b2da3732ee4d37af66b/src/script/interpreter.cpp#L1478
// And: https://github.com/github.com/mrtnetwork/bitcoin/github.com/mrtnetwork/bitcoin/blob/b5f33ac1f82aea290b4653af36ac2ad1bf1cce7b/test/functional/test_framework/script.py

//	  SIGHASH types (see constants.py):
//		 TAPROOT_SIGHASH_ALL - signs all inputs and outputs (default)
//	      SIGHASH_ALL - signs all inputs and outputs
//	      SIGHASH_NONE - signs all of the inputs
//	      SIGHASH_SINGLE - signs all inputs but only txin_index output
//	      SIGHASH_ANYONECANPAY (only combined with one of the above)
//	      - with ALL - signs all outputs but only txin_index input
//	      - with NONE - signs only the txin_index input
//	      - with SINGLE - signs txin_index input and output
//
// txin_index : The index of the input that we wish to sign
// script_pubkeys : The scriptPubkeys that correspond to all the inputs/UTXOs
// amounts : The amounts that correspond to all the inputs/UTXOs
// ext_flag : Extension mechanism, default is 0; 1 is for script spending (BIP342)
// script : The script that we are spending (ext_flag=1)
// leaf_ver : The script version, LEAF_VERSION_TAPSCRIPT for the default tapscript
// sighash : The type of the signature hash to be created
func (tx *BtcTransaction) GetTransactionTaprootDigest(txIndex int, scriptPubKeys []*Script, amounts []*big.Int, extFlags int, script *Script, sighash int) []byte {
	newTx := tx.Copy()
	sighashNone := (sighash & 0x03) == constant.SIGHASH_NONE
	sighashSingle := (sighash & 0x03) == constant.SIGHASH_SINGLE
	anyoneCanPay := (sighash & 0x80) == constant.SIGHASH_ANYONECANPAY
	txForSign := make([]byte, 0)
	txForSign = append(txForSign, byte(0))
	txForSign = append(txForSign, byte(sighash))
	txForSign = append(txForSign, tx.Version...)
	txForSign = append(txForSign, tx.Locktime...)
	hashPrevouts := make([]byte, 0)
	hashAmounts := make([]byte, 0)
	hashScriptPubkeys := make([]byte, 0)
	hashSequences := make([]byte, 0)
	hashOutputs := make([]byte, 0)

	if !anyoneCanPay {
		for _, txin := range newTx.Inputs {
			txidBytes := formating.ReverseBytes(formating.HexToBytes(txin.TxID))
			txoutIndexBytes := formating.PackUint32LE(uint32(txin.TxIndex))
			hashPrevouts = append(hashPrevouts, txidBytes...)
			hashPrevouts = append(hashPrevouts, txoutIndexBytes...)
		}
		hashPrevouts = digest.SingleHash(hashPrevouts)
		txForSign = append(txForSign, hashPrevouts...)

		for _, i := range amounts {
			bytes := formating.PackBigIntToLittleEndian(i)
			hashAmounts = append(hashAmounts, bytes...)
		}
		hashAmounts = digest.SingleHash(hashAmounts)
		txForSign = append(txForSign, hashAmounts...)

		for _, s := range scriptPubKeys {
			h := s.ToHex() // must checked
			scriptLen := len(h) / 2
			scriptBytes := formating.HexToBytes(h)
			lenBytes := []byte{byte(scriptLen)}
			hashScriptPubkeys = append(hashScriptPubkeys, lenBytes...)
			hashScriptPubkeys = append(hashScriptPubkeys, scriptBytes...)
		}
		hashScriptPubkeys = digest.SingleHash(hashScriptPubkeys)
		txForSign = append(txForSign, hashScriptPubkeys...)

		for _, txIn := range newTx.Inputs {
			hashSequences = append(hashSequences, txIn.Sequence...)
		}
		hashSequences = digest.SingleHash(hashSequences)
		txForSign = append(txForSign, hashSequences...)
	}

	if !(sighashNone || sighashSingle) {
		for _, txOut := range newTx.Outputs {
			packedAmount := formating.PackBigIntToLittleEndian(txOut.Amount)
			scriptBytes := txOut.ScriptPubKey.ToBytes()
			lenScriptBytes := []byte{byte(len(scriptBytes))}
			hashOutputs = append(hashOutputs, packedAmount...)
			hashOutputs = append(hashOutputs, lenScriptBytes...)
			hashOutputs = append(hashOutputs, scriptBytes...)
		}
		hashOutputs = digest.SingleHash(hashOutputs)
		txForSign = append(txForSign, hashOutputs...)
	}

	spendType := extFlags*2 + 0
	txForSign = append(txForSign, byte(spendType))

	if anyoneCanPay {
		txin := newTx.Inputs[txIndex]
		txidBytes := formating.ReverseBytes(formating.HexToBytes(txin.TxID))
		txoutIndexBytes := formating.PackUint32LE(uint32(txin.TxIndex))
		result := append(txidBytes, txoutIndexBytes...)
		txForSign = append(txForSign, result...)
		txForSign = append(txForSign, formating.PackBigIntToLittleEndian(amounts[txIndex])...)
		sPubKey := scriptPubKeys[txIndex].ToHex()
		sLength := len(sPubKey) / 2
		txForSign = append(txForSign, byte(sLength))
		txForSign = append(txForSign, formating.HexToBytes(sPubKey)...)
		txForSign = append(txForSign, newTx.Inputs[txIndex].Sequence...)
	} else {
		index := txIndex
		bytes := make([]byte, 4)
		for i := 0; i < 4; i++ {
			bytes[i] = byte(index & 0xFF)
			index >>= 8
		}
		txForSign = append(txForSign, bytes...)
	}

	if sighashSingle {
		txOut := newTx.Outputs[txIndex]
		packedAmount := formating.PackBigIntToLittleEndian(txOut.Amount)
		sBytes := txOut.ScriptPubKey.ToBytes()
		lenScriptBytes := []byte{byte(len(sBytes))}
		hashOut := append(packedAmount, lenScriptBytes...)
		hashOut = append(hashOut, sBytes...)
		txForSign = append(txForSign, digest.SingleHash(hashOut)...)
	}

	if extFlags == 1 {

		leafVarBytes := append([]byte{byte(constant.LEAF_VERSION_TAPSCRIPT)}, formating.PrependVarint(script.ToBytes())...)
		txForSign = append(txForSign, digest.TaggedHash(leafVarBytes, "TapLeaf")...)
		txForSign = append(txForSign, 0)
		txForSign = append(txForSign, []byte{0xFF, 0xFF, 0xFF, 0xFF}...)
	}

	return digest.TaggedHash(txForSign, "TapSighash")
}

// Converts object to hexadecimal string
func (tx *BtcTransaction) ToHex() string {
	toBytes := tx.ToBytes(tx.HasSegwit)
	return formating.BytesToHex(toBytes)
}

// Converts object to hexadecimal string
func (tx *BtcTransaction) Serialize() string {
	return tx.ToHex()
}

// Hashes the serialized (bytes) tx to get a unique id
func (tx *BtcTransaction) TxId() string {
	toBytes := tx.ToBytes(false)
	h := formating.ReverseBytes(digest.DoubleHash(toBytes))
	return formating.BytesToHex(h)
}

// Gets the size of the transaction
func (tx *BtcTransaction) GetSize() int {
	toBytes := tx.ToBytes(tx.HasSegwit)
	return len(toBytes)
}

// Gets the virtual size of the transaction.
// For non-segwit txs this is identical to get_size(). For segwit txs the
// marker and witnesses length needs to be reduced to 1/4 of its original
// length. Thus it is substructed from size and then it is divided by 4
// before added back to size to produce vsize (always rounded up).
// https://en.bitcoin.it/wiki/Weight_units
func (tx *BtcTransaction) GetVSize() int {
	if !tx.HasSegwit {
		return tx.GetSize()
	}
	markerSize := 2
	witSize := 0
	data := make([]byte, 0)

	for _, w := range tx.Witnesses {
		countBytes := []byte{byte(len(w.Stack))}
		data = append(data, countBytes...)
		data = append(data, w.ToBytes()...)
	}

	witSize = len(data)
	size := tx.GetSize() - (markerSize + witSize)
	vSize := float64(size + (markerSize+witSize)/4)
	return int(vSize)
}

// Hashes the serialized (bytes) tx including segwit marker and witnesses"
func (tx *BtcTransaction) GetHash() string {
	toBytes := tx.ToBytes(tx.HasSegwit)
	toHash := digest.DoubleHash(toBytes)
	revers := formating.ReverseBytes(toHash)
	return formating.BytesToHex(revers)
}

// Hashes the serialized (bytes) tx including segwit marker and witnesses
func (tx *BtcTransaction) GetWTXID() string {
	return tx.GetHash()
}

// SetScriptSig sets the script signature for the transaction input at the specified index.
func (tx *BtcTransaction) SetScriptSig(index int, script *Script) {
	tx.Inputs[index].ScriptSig = script
}
