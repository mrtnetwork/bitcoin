package provider

import (
	"math/big"
)

// BitcoinFeeRate represents a structure for storing different Bitcoin fee rate levels.
// It includes high, medium, and low fee rates, each as a pointer to a big.Int value.
//
// The big.Int pointers allow for precise representation of fee rates, where nil values
// indicate that the fee rate level is not available or unspecified.
type BitcoinFeeRate struct {
	High   *big.Int // High fee rate in satoshis per byte
	Medium *big.Int // Medium fee rate in satoshis per byte
	Low    *big.Int // Low fee rate in satoshis per byte
}

// parseMempoolFees takes a data interface and converts it to a big.Int representing
// mempool fees in satoshis per kilobyte (sat/KB). The function performs the conversion
// based on the type of the input data, which can be either a float64 (floating-point
// fee rate) or an int (integer fee rate in satoshis per byte).
func parseMempoolFees(data interface{}) *big.Int {
	const kb = 1024

	switch v := data.(type) {
	case float64:
		return new(big.Int).SetInt64(int64(v * kb))
	case int:
		return new(big.Int).SetInt64(int64(v * kb))
	default:
		return nil
	}
}

// NewBitcoinFeeRateFromMempool creates a BitcoinFeeRate structure from JSON data retrieved
// from a mempool API response. The function parses the JSON map and extracts fee rate
// information for high, medium, and low fee levels.
func NewBitcoinFeeRateFromMempool(json map[string]interface{}) *BitcoinFeeRate {
	return &BitcoinFeeRate{
		High:   parseMempoolFees(json["fastestFee"]),
		Medium: parseMempoolFees(json["halfHourFee"]),
		Low:    parseMempoolFees(json["minimumFee"]),
	}
}

// NewBitcoinFeeRateFromBlockCyper creates a BitcoinFeeRate structure from JSON data retrieved
// from a blockCypher API response. The function parses the JSON map and extracts fee rate
// information for high, medium, and low fee levels.
func NewBitcoinFeeRateFromBlockCyper(json map[string]interface{}) *BitcoinFeeRate {
	return &BitcoinFeeRate{
		High:   new(big.Int).SetInt64(int64(json["high_fee_per_kb"].(float64))),
		Medium: new(big.Int).SetInt64(int64(json["medium_fee_per_kb"].(float64))),
		Low:    new(big.Int).SetInt64(int64(json["low_fee_per_kb"].(float64))),
	}
}

func (b BitcoinFeeRate) String() string {
	return "high: " + b.High.String() + " medium: " + b.Medium.String() + " low: " + b.Low.String()
}

// GetEstimate calculates the estimated fee in satoshis for a given transaction size
// and fee rate (in satoshis per kilobyte) using the formula:
//
//	EstimatedFee = (TransactionSize * FeeRate) / 1024
//
// Parameters:
// - trSize: An integer representing the transaction size in bytes.
// - feeRate: A pointer to a big.Int representing the fee rate in satoshis per kilobyte.
//
// Returns:
// - *big.Int: A pointer to a big.Int containing the estimated fee in satoshis.
func (b BitcoinFeeRate) GetEstimate(trSize int, feeRate *big.Int) *big.Int {
	trSizeBigInt := new(big.Int).SetInt64(int64(trSize))
	return new(big.Int).Div(new(big.Int).Mul(trSizeBigInt, feeRate), big.NewInt(1024))
}
