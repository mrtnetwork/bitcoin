package scripts

import "bitcoin/constant"

type Sequence struct {
	// Specifies the type of sequence (TYPE_RELATIVE_TIMELOCK | TYPE_ABSOLUTE_TIMELOCK | TYPE_REPLACE_BY_FEE
	seqType int
	// The value of the block height or the 512 seconds increments
	value int
	// If type is TYPE_RELATIVE_TIMELOCK then this specifies its type (block height or 512 secs increments)
	isTypeBlock bool
}

// NewSequence creates a new Sequence object with the specified sequence type
func NewSequence(seqType, value int, isTypeBlock bool) *Sequence {
	if seqType == constant.TYPE_RELATIVE_TIMELOCK && (value < 1 || value > 0xffff) {
		panic("Sequence should be between 1 and 65535")
	}

	return &Sequence{seqType: seqType, value: value, isTypeBlock: isTypeBlock}
}

// Serializes the relative sequence as required in a transaction
func (s *Sequence) ForInputSequence() []byte {
	if s.seqType == constant.TYPE_ABSOLUTE_TIMELOCK {
		return constant.ABSOLUTE_TIMELOCK_SEQUENCE
	}

	if s.seqType == constant.TYPE_REPLACE_BY_FEE {
		return constant.REPLACE_BY_FEE_SEQUENCE
	}

	if s.seqType == constant.TYPE_RELATIVE_TIMELOCK {
		seq := 0
		if !s.isTypeBlock {
			seq |= 1 << 22
		}
		seq |= s.value
		return []byte{
			byte(seq & 0xFF),
			byte((seq >> 8) & 0xFF),
			byte((seq >> 16) & 0xFF),
			byte((seq >> 24) & 0xFF),
		}
	}

	panic("Invalid seqType")
}

// Returns the appropriate integer for a script; e.g. for relative timelocks
func (s *Sequence) ForScript() int {
	if s.seqType == constant.TYPE_REPLACE_BY_FEE {
		panic("RBF is not to be included in a script.")
	}

	scriptInteger := s.value
	if s.seqType == constant.TYPE_RELATIVE_TIMELOCK && !s.isTypeBlock {
		scriptInteger |= 1 << 22
	}
	return scriptInteger
}
