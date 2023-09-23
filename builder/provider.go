package builder

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

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
		fmt.Println("Error marshaling JSON:", err)
		return "", err
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
func (api *aPIConfig) GetUtxo(ownerDetals UtxoOwnerDetails) (UtxoWithOwnerList, error) {
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
				fmt.Println("Error decoding JSON:", err)
				return nil, err
			}
			return utxos.ToUtxoWithOwner(ownerDetals), nil
		}
	}
	var addressInfo blockCypherUtxo
	if err := json.Unmarshal(body, &addressInfo); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return nil, err
	}
	return addressInfo.ToUtxoWithOwner(ownerDetals), nil
}

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
		fmt.Println("Error marshaling JSON:", err)
		return "", err
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

	return string(body), nil
}
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
		fmt.Println("Error decoding JSON:", err)
		return nil, err
	}
	switch api.ApiType {
	case MempoolApi:
		{
			return NewBitcoinFeeRateFromMempool(fee), nil
		}
	}

	return NewBitcoinFeeRateFromBlockCyper(fee), nil
}
