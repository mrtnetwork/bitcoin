package example

import (
	"fmt"
	"github.com/mrtnetwork/bitcoin/address"
	"github.com/mrtnetwork/bitcoin/provider"
)

func ExampleNetworkProvider() {
	// select network testnet or mainnet
	network := address.TestnetNetwork

	// create api (BlockCyperApi or MempoolApi)
	// Currently, only a few critical methods have been implemented to retrieve unspent transactions,
	// obtain network fees, receive transactions, and send transactions to the network.
	api := provider.SelectApi(provider.MempoolApi, &network)

	// ========================================================================================//

	// Read Transaction id(hash)
	tr, e := api.GetTransaction("d4bad8e07d30ca4389ec8a203318aa523cc3e36c9730d0a6852a3801d086c5fe")
	if e != nil {
		fmt.Println("error: ", e)
		return
	}
	if converted, ok := tr.(*provider.BlocCyperTransaction); ok {
		fmt.Println("is blockcypher transaction struct")
		fmt.Println(converted.Hash)
		fmt.Println(converted.Inputs)
		fmt.Println(converted.Outputs)
		fmt.Println(converted.Confirmations)
	} else if converted, ok := tr.(*provider.MempoolTransaction); ok {
		fmt.Println("Memool transaction struct")
		fmt.Println(converted.TxID)
		fmt.Println(converted.Vout)
		fmt.Println(converted.Vin)
		fmt.Println(converted.Status.Confirmed)
	}

	// ========================================================================================//

	addr, _ := address.P2WPKHAddresssFromAddress("tb1q92nmnvhj04sqd4x7wjaewlt5jn8n3ngmplcymy")

	// Read accounts UTXOS
	utxos, e := api.GetAccountUtxo(provider.UtxoOwnerDetails{
		PublicKey: "",
		Address:   addr,
	})
	if e != nil {
		fmt.Println(e)
	} else {
		fmt.Println("utxos ", len(utxos))
		for i := 0; i < len(utxos); i++ {
			fmt.Println("TxID: ", utxos[i].Utxo.TxHash)
			fmt.Println("Value: ", utxos[i].Utxo.Value)
			fmt.Println("ScriptType: ", utxos[i].Utxo.ScriptType)
			fmt.Println("Vout: ", utxos[i].Utxo.Vout)
			fmt.Println("BlockHeight: ", utxos[i].Utxo.BlockHeight)
		}
	}

	// ========================================================================================//

	// Network fee
	fee, e := api.GetNetworkFee()
	if e != nil {
		fmt.Println(e)
	} else {
		// PER KB
		fmt.Println("MEDIUM: ", fee.Medium)
		fmt.Println("LOW: ", fee.Low)
		fmt.Println("LOW: ", fee.High)
		// for calculate transaction fee u can use estimate fee method of BitcoinFeeRate struct
		// transaction size transaction.GetSize() or transaction.GetVSize() for segwit transaction
		_ = fee.GetEstimate(500, fee.High)
	}

	// ========================================================================================//

	//  Send transaction
	_, _ = api.SendRawTransaction("TRANSACTION DIGEST")

	// ========================================================================================//

	// Read account transactions
	transaction, _ := api.GetAccountTransactions(addr.Show(network), func(url string) string {
		/*
			You have the option to modify the address before making the request,
			such as adding parameters like a limit or page number. For more information,
			please consult the Mempool API or BlockCypher documentation.
			You have the option to modify the address before making the request,
			such as adding parameters like a limit or page number. For more information, please consult the Mempool API or BlockCypher documentation.

		*/
		return url
	})
	if converted, ok := transaction.(provider.MemoolTransactionList); ok {
		fmt.Println("is mempool transactions struct")
		fmt.Println("transactions: ", len(converted))
		for i := 0; i < len(converted); i++ {
			fmt.Println("transactions: ", converted[i].TxID)
			fmt.Println("transactions: ", converted[i].Status)
		}
	} else if converted, ok := transaction.(provider.BlockCypherTransactionList); ok {
		fmt.Println("is blockCypher transaction struct")
		for i := 0; i < len(converted); i++ {
			fmt.Println("transactions: ", converted[i].Hash)
			fmt.Println("transactions: ", converted[i].Confirmations)
		}
	}

}
