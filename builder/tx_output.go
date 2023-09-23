package builder

import (
	"bitcoin/address"
	"math/big"
)

type BitcoinOutputDetails struct {
	Address address.BitcoinAddress
	Value   *big.Int
}
