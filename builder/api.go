package builder

import (
	"bitcoin/address"
	"strings"
)

type APIType int

const (
	MempoolApi APIType = iota
	BlockCyperApi
)

type aPIConfig struct {
	URL             string
	FeeRate         string
	Transaction     string
	SendTransaction string
	ApiType         APIType
	Network         address.Network
}

const (
	blockCypherBaseURL = "https://api.blockcypher.com/v1/btc/test3"
	mempoolBaseURL     = "https://mempool.space/testnet/api"
	blockstreamBaseURL = "https://blockstream.info/testnet/api"
)
const (
	blockCypherMianBaseURL = "https://api.blockcypher.com/v1/btc/main"
	mempoolMainBaseURL     = "https://mempool.space/api"
	blockstreamMainBaseURL = "https://blockstream.info/api"
)

func createMempolApi(network address.NetworkInfo) *aPIConfig {
	baseUrl := mempoolMainBaseURL
	if !network.IsMainNet() {
		baseUrl = mempoolBaseURL
	}

	return &aPIConfig{
		URL:             baseUrl + "/address/###/utxo",
		FeeRate:         baseUrl + "/v1/fees/recommended",
		Transaction:     baseUrl + "/tx/###",
		SendTransaction: baseUrl + "/tx",
		ApiType:         MempoolApi,
		Network:         network.Network(),
	}
}
func createBlockCyperApi(network address.NetworkInfo) *aPIConfig {
	baseUrl := blockCypherMianBaseURL
	if !network.IsMainNet() {
		baseUrl = blockCypherBaseURL
	}

	return &aPIConfig{
		URL:             baseUrl + "/addrs/###/?unspentOnly=true&includeScript=true&limit=2000",
		FeeRate:         baseUrl,
		Transaction:     baseUrl + "/txs/###",
		SendTransaction: baseUrl + "/txs/push",
		ApiType:         BlockCyperApi,
		Network:         network.Network(),
	}
}

func SelectApi(apitype APIType, network address.NetworkInfo) *aPIConfig {
	switch apitype {
	case MempoolApi:
		{
			return createMempolApi(network)
		}
	default:
		{
			return createBlockCyperApi(network)
		}
	}
}
func (api *aPIConfig) GetUtxoUrl(address string) string {
	baseUrl := api.URL
	return strings.Replace(baseUrl, "###", address, -1)
}
func (api *aPIConfig) GetFeeApiUrl() string {
	return api.FeeRate
}

func (api *aPIConfig) GetNetwork() address.NetworkInfo {
	switch api.Network {
	case address.Mainnet:
		{
			return &address.MainnetNetwork
		}
	default:
		{
			return &address.TestnetNetwork
		}
	}
}
func (api *aPIConfig) GetSendTransactionUrl() string {
	return api.SendTransaction
}
