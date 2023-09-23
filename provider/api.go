package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Tokenize func(url string) string

func TestMempoolAccept(txDigest []interface{}) (string, error) {
	url := "https://btc.getblock.io/786c97b8-f53f-427b-80f7-9af7bd5bdb84/testnet/"
	data := struct {
		Jsonrpc string        `json:"jsonrpc"`
		Method  string        `json:"method"`
		Id      string        `json:"id"`
		Params  []interface{} `json:"params"`
	}{
		Jsonrpc: "2.0",
		Method:  "testmempoolaccept",
		Id:      "123",
		Params:  txDigest,
	}

	// Marshal the struct into a JSON byte array
	payload, err := json.Marshal(data)
	if err != nil {

		return "", fmt.Errorf("error marshaling JSON: %v", err)
	}
	// Send an HTTP GET request.
	response, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	// Read the response body using io.Copy.
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// GetAccountUtxo retrieves a list of unspent transaction outputs (UTXOs) associated with a specific
// account
func (api *aPIConfig) GetAccountUtxo(ownerDetals UtxoOwnerDetails) (UtxoWithOwnerList, error) {
	url := api.GetUtxoUrl(ownerDetals.Address.Show(api.GetNetwork()))
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Read the response body using io.Copy.
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	switch api.ApiType {
	case MempoolApi:
		{
			var utxos MempolUtxoList
			if err := json.Unmarshal(body, &utxos); err != nil {
				return nil, fmt.Errorf("error decoding JSON: %v", err)
			}
			return utxos.ToUtxoWithOwner(ownerDetals), nil
		}
	}
	var addressInfo blockCypherUtxo
	if err := json.Unmarshal(body, &addressInfo); err != nil {
		return nil, fmt.Errorf("error decoding JSON: %v", err)
	}
	return addressInfo.ToUtxoWithOwner(ownerDetals), nil
}

// SendRawTransaction sends a raw transaction represented by its text digest to the blockchain
// network using the configured API. The method submits the transaction to the network for processing
// and returns the resulting transaction ID (TxID) if the submission is successful.
func (api *aPIConfig) SendRawTransaction(textDigest string) (string, error) {
	url := api.GetSendTransactionUrl()

	switch api.ApiType {
	case MempoolApi:
		{
			// Send an HTTP GET request.
			response, err := http.Post(url, "application/json", bytes.NewBuffer([]byte(textDigest)))
			if err != nil {
				return "", err
			}
			defer response.Body.Close()

			// Read the response body using io.Copy.
			body, err := io.ReadAll(response.Body)
			if err != nil {
				return "", err
			}

			return string(body), nil
		}

	}
	data := map[string]interface{}{
		"tx": textDigest,
	}

	// Marshal the map into JSON
	jsonData, err := json.Marshal(data)
	if err != nil {

		return "", fmt.Errorf("error marshaling JSON: %v", err)
	}
	// Send an HTTP GET request.
	response, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	// Read the response body using io.Copy.
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	var transaction BlocCyperTransaction
	if err := json.Unmarshal(body, &transaction); err != nil {
		return "", fmt.Errorf("error decoding JSON: %v", err)
	}

	return transaction.Hash, nil
}

// GetNetworkFee retrieves the current network fee rate information from the configured API. The method
// fetches and returns a BitcoinFeeRate structure containing fee rates for high, medium, and low priority
// transactions, which can be used to estimate transaction fees.
func (api *aPIConfig) GetNetworkFee() (*BitcoinFeeRate, error) {
	url := api.GetFeeApiUrl()
	// Send an HTTP GET request.
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Read the response body using io.Copy.
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var fee map[string]interface{}
	if err := json.Unmarshal(body, &fee); err != nil {
		return nil, fmt.Errorf("error decoding JSON: %v", err)
	}
	switch api.ApiType {
	case MempoolApi:
		{
			return NewBitcoinFeeRateFromMempool(fee), nil
		}
	}

	return NewBitcoinFeeRateFromBlockCyper(fee), nil
}

// GetTransaction retrieves information about a specific transaction identified by its unique transaction ID
// from the blockchain network using the configured API. The method fetches and returns information about
// the transaction, including details about its inputs, outputs, status, and more.
func (api *aPIConfig) GetTransaction(transactionId string) (interface{}, error) {
	url := api.GetTransactionUrl(transactionId)

	// Send an HTTP GET request.
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Read the response body using io.Copy.
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	switch api.ApiType {
	case MempoolApi:
		var transaction MempoolTransaction
		if err := json.Unmarshal(body, &transaction); err != nil {
			return nil, fmt.Errorf("error decoding JSON: %v", err)
		}
		return &transaction, nil
	}
	var transaction BlocCyperTransaction
	if err := json.Unmarshal(body, &transaction); err != nil {
		return nil, fmt.Errorf("error decoding JSON: %v", err)
	}
	return &transaction, nil
}

// GetAccountTransactions retrieves a list of transactions associated with a specific account or address
// from the blockchain network using the configured API. The method fetches and returns information about
// the account's transactions, including details such as transaction IDs, timestamps, inputs, outputs, and more.
//
// Parameters:
//   - api: A pointer to the aPIConfig instance representing the API configuration used to make the request.
//   - address: A string representing the account's or address's unique identifier.
//   - tokenize: A Tokenize function pointer that can be used to tokenize the API request URL, allowing for
//     customization of the request parameters.
func (api *aPIConfig) GetAccountTransactions(address string, tokenize Tokenize) (interface{}, error) {
	url := api.GetTransactionsUrl(address)
	url = tokenize(url)
	// Send an HTTP GET request.
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Read the response body using io.Copy.
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	switch api.ApiType {
	case MempoolApi:
		var transaction MemoolTransactionList
		if err := json.Unmarshal(body, &transaction); err != nil {
			return nil, fmt.Errorf("error decoding JSON: %v", err)
		}
		return &transaction, nil
	}
	var transaction BlockCypherAddressInfo
	if err := json.Unmarshal(body, &transaction); err != nil {
		return nil, fmt.Errorf("error decoding JSON: %v", err)
	}
	return &transaction.TXs, nil
}
