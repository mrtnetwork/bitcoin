// access BlockCypher and Mempool APIs for fetching UTXos, transaction data, network fees, and sending transactions in the Bitcoin network
package provider

import (
	"strings"

	"github.com/mrtnetwork/bitcoin/address"
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
	Transactions    string
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
		Transactions:    baseUrl + "/address/###/txs",
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
		Transactions:    baseUrl + "/addrs/###/full?limit=200",
		ApiType:         BlockCyperApi,
		Network:         network.Network(),
	}
}

// SelectApi chooses and returns the appropriate API configuration based on the given APIType
// and network information.
//
// Parameters:
// - apitype: The APIType representing the desired API.
// - network: The address.NetworkInfo providing network-specific details.
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

// Return UTXO url contains address
func (api *aPIConfig) GetUtxoUrl(address string) string {
	baseUrl := api.URL
	return strings.Replace(baseUrl, "###", address, -1)
}

// Return Network Fee api url contains address
func (api *aPIConfig) GetFeeApiUrl() string {
	return api.FeeRate
}

// get current network of api
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

// get send transaction url
func (api *aPIConfig) GetSendTransactionUrl() string {
	return api.SendTransaction
}

// get transaction url contains hash
func (api *aPIConfig) GetTransactionUrl(transactionId string) string {
	baseUrl := api.Transaction
	return strings.Replace(baseUrl, "###", transactionId, -1)
}

// get account transaction url contains address
func (api *aPIConfig) GetTransactionsUrl(address string) string {
	baseUrl := api.Transactions
	return strings.Replace(baseUrl, "###", address, -1)
}
