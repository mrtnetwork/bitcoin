package scripts

import (
	"github.com/mrtnetwork/bitcoin/constant"
	"github.com/mrtnetwork/bitcoin/formating"
)

type ControlBlock struct {
	//  the internal public key object
	PublicXonly string
	// concatenated path (leafs/branches) hashes in bytes
	Scripts []byte
}

// NewControlBlock creates a new control block with the specified public key and scripts,
// and returns a pointer to the initialized ControlBlock object.
func NewControlBlock(public string, scripts []byte) *ControlBlock {
	return &ControlBlock{
		PublicXonly: public,
		Scripts:     scripts,
	}
}

// returns the control block as bytes
func (cb *ControlBlock) ToBytes() []byte {
	version := []byte{constant.LEAF_VERSION_TAPSCRIPT}

	pubKey := formating.HexToBytes(cb.PublicXonly)

	// If Scripts is nil, create an empty slice
	marklePath := cb.Scripts
	if marklePath == nil {
		marklePath = []byte{}
	}

	result := append(version, pubKey...)
	result = append(result, marklePath...)

	return result
}

// returns the control block as a hexadecimal string
func (cb *ControlBlock) ToHex() string {
	toBytes := cb.ToBytes()
	return formating.BytesToHex(toBytes)
}
