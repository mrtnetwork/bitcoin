package builder

import (
	"math/big"
)

type BitcoinFeeRate struct {
	High   *big.Int
	Medium *big.Int
	Low    *big.Int
}

func parseMempoolFees(data interface{}) *big.Int {
	const kb = 1024
	switch v := data.(type) {
	case float64:
		return new(big.Int).SetInt64(int64(v * kb))
	case int:
		return new(big.Int).SetInt64(int64(v * kb))
	default:
		return nil // Handle other types accordingly
	}
}

func NewBitcoinFeeRateFromMempool(json map[string]interface{}) *BitcoinFeeRate {
	return &BitcoinFeeRate{
		High:   parseMempoolFees(json["fastestFee"]),
		Medium: parseMempoolFees(json["halfHourFee"]),
		Low:    parseMempoolFees(json["minimumFee"]),
	}
}

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

func (b BitcoinFeeRate) GetEstimate(trSize int, feeRate *big.Int) *big.Int {
	trSizeBigInt := new(big.Int).SetInt64(int64(trSize))
	return new(big.Int).Div(new(big.Int).Mul(trSizeBigInt, feeRate), big.NewInt(1024))
}
