package scripts

import (
	"bitcoin/formating"
)

// A list of the witness items required to satisfy the locking conditions of a segwit input (aka witness stack).
type TxWitnessInput struct {
	// the witness items (hex str) list
	Stack []string
}

// NewTxWitnessInput creates a new transaction witness input with the provided stack of strings
// and returns a pointer to the initialized TxWitnessInput object.
func NewTxWitnessInput(stack ...string) *TxWitnessInput {
	return &TxWitnessInput{
		Stack: append([]string{}, stack...),
	}
}

// Deep copy of TxWitnessInput
func (twi *TxWitnessInput) Copy() *TxWitnessInput {
	copiedStack := make([]string, len(twi.Stack))
	copy(copiedStack, twi.Stack)
	return &TxWitnessInput{
		Stack: copiedStack,
	}
}

// Converts to bytes
func (twi *TxWitnessInput) ToBytes() []byte {
	var stackBytes []byte

	for _, item := range twi.Stack {
		itemBytes := formating.PrependVarint(formating.HexToBytes(item))
		stackBytes = append(stackBytes, itemBytes...)
	}

	return stackBytes
}
