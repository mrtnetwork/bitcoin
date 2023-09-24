package example

import (
	"github.com/mrtnetwork/bitcoin/address"
	"github.com/mrtnetwork/bitcoin/provider"

	"fmt"
	"github.com/mrtnetwork/bitcoin/constant"
	hdwallet "github.com/mrtnetwork/bitcoin/hd_wallet"
	"github.com/mrtnetwork/bitcoin/keypair"
	"math/big"
)

/*
To estimate the transaction cost, we must first create a mock transaction with the desired UTXO and the number of recipients,
similar to the examples of sending a transaction. This allows us to determine the transaction size.
After obtaining the transaction size, we can then create a real transaction with a new fee and send it.
*/
func TestExampleCalculateFee() {
	network := address.TestnetNetwork
	/*
		Avoid using Mempool to estimate costs in the Testnet network,
		as the transaction cost is not genuine. Using it may likely result
		in errors when attempting to send the transaction.

		Please note that Mempool displays the transaction cost in bytes,
		whereas BlockCypher displays the transaction cost in kilobytes.
		The GetEstimate method accurately calculates the cost based on the selected API.
	*/
	api := provider.SelectApi(provider.BlockCyperApi, &network)
	// i generate random mnemonic for test
	// mnemoic, _ := bip39.GenerateMnemonic(256)
	mnemonic := "spy often critic spawn produce volcano depart fire theory fog turn retire"

	// accsess to private and public keys
	masterWallet, _ := hdwallet.FromMnemonic(mnemonic, "")

	// wallet with path
	// i generate 4 HD wallet for this test and now i have access to private and pulic key of each wallet
	sp1, _ := hdwallet.DrivePath(masterWallet, "m/44'/0'/0'/0/0/1")
	sp2, _ := hdwallet.DrivePath(masterWallet, "m/44'/0'/0'/0/0/2")
	sp3, _ := hdwallet.DrivePath(masterWallet, "m/44'/0'/0'/0/0/3")
	sp4, _ := hdwallet.DrivePath(masterWallet, "m/44'/0'/0'/0/0/4")

	// access to private key `ECPrivate`
	private1, _ := sp1.GetPrivate()
	private2, _ := sp2.GetPrivate()
	private3, _ := sp3.GetPrivate()
	private4, _ := sp4.GetPrivate()
	// access to public key `ECPublic`
	public1 := sp1.GetPublic()
	public2 := sp2.GetPublic()
	public3 := sp3.GetPublic()
	public4 := sp4.GetPublic()

	// now we need some address for spending or receive let doint
	// For our test, I use public key to create addresses

	// P2PKH ADDRESS
	// myVMJgRi6arv4hLbeUcJYKUJWmFnpjtVme
	// equals to exampleAddr1 := address.P2PKHAddressFromAddress("myVMJgRi6arv4hLbeUcJYKUJWmFnpjtVme")
	exampleAddr1 := public1.ToAddress()

	// P2TR ADDRESS
	// tb1pyhmqwlcrws4dxcgalt4mrffgnys879vs59xf6sve4hazyvmhecxq3e6sc0
	// equals to exampleAddr2 := address.P2TRAddressFromAddress("tb1pyhmqwlcrws4dxcgalt4mrffgnys879vs59xf6sve4hazyvmhecxq3e6sc0")
	exampleAddr2 := public2.ToTaprootAddress()

	// P2SH(P2PKH) ADDRESS
	// 2N2yqygBJRvDzLzvPe91qKfSYnK5utGckJX
	// equals to exampleAddr3 := address.P2SHAddressFromAddress("2N2yqygBJRvDzLzvPe91qKfSYnK5utGckJX", address.P2PKHInP2SH)
	exampleAddr3 := public2.ToP2PKHInP2SH()

	// P2PKH ADDRESS
	// mzUzciYUGsNxLCaaHwou27F4RbnDTzKomV
	// equals to exampleAddr4 := address.P2PKHAddressFromAddress("mzUzciYUGsNxLCaaHwou27F4RbnDTzKomV")
	exampleAddr4 := public3.ToAddress()

	// P2SH(P2PKH) ADDRESS
	// 2MzibgEeJYCN8mjJsZTg79AH7au4PCkHXHo
	// equals to exampleAddr5 := address.P2SHAddressFromAddress("2MzibgEeJYCN8mjJsZTg79AH7au4PCkHXHo", address.P2PKHInP2SH)
	exampleAddr5 := public3.ToP2PKHInP2SH()

	// P2SH(P2WSH) ADDRESS
	// 2N7bNV1WPwCVHfoqqRhvtmbAfktazfjHEW2
	// equals to exampleAddr6 := address.P2SHAddressFromAddress("2N7bNV1WPwCVHfoqqRhvtmbAfktazfjHEW2", address.P2WSHInP2SH)
	exampleAddr6 := public3.ToP2WSHInP2SH()

	// P2SH(P2WPKH) ADDRESS
	// 2N38S8G9q6qyEjPWmicqxrVwjq4QiTbcyf4
	// equals to exampleAddr7 := address.P2SHAddressFromAddress("2N38S8G9q6qyEjPWmicqxrVwjq4QiTbcyf4", address.P2WPKHInP2SH)
	exampleAddr7 := public3.ToP2WPKHInP2SH()

	// P2SH(P2PK) ADDRESS
	// 2MugsNcgzLJ1HosnZyC2CfZVmgbMPK1XubR
	// equals to exampleAddr8 := address.P2SHAddressFromAddress("2MugsNcgzLJ1HosnZyC2CfZVmgbMPK1XubR", address.P2PKInP2SH)
	exampleAddr8 := public4.ToP2PKInP2SH()

	// P2WPKH ADDRESS
	// tb1q6q9halaazasd42gzsc2cvv5xls295w7kawkhxy
	// equals to exampleAddr9 := address.P2WPKHAddresssFromAddress("tb1q6q9halaazasd42gzsc2cvv5xls295w7kawkhxy")
	exampleAddr9 := public3.ToSegwitAddress()

	// P2WSH ADDRESS
	// tb1qf4qwtr5kp5q87dtp3ul3402vkzssxfv7f4aettjq2hcfhnt92dmq5xzs6n
	// equals to exampleAddr10 := address.P2WSHAddresssFromAddress("tb1qf4qwtr5kp5q87dtp3ul3402vkzssxfv7f4aettjq2hcfhnt92dmq5xzs6n")
	// created with 1-1 MultiSig script: ["OP_1", publicHex(Hex of compressed public key) , "OP_1", "OP_CHECKMULTISIG"]
	exampleAddr10 := public3.ToP2WSHAddress()

	// now we chose some address for spending from multiple address
	// i use some different address type for this
	spenders := []provider.UtxoOwnerDetails{
		{PublicKey: public1.ToHex(), Address: exampleAddr1}, // p2pkh address from public1
		{PublicKey: public2.ToHex(), Address: exampleAddr2}, // P2TRAddress address from public2
		{PublicKey: public3.ToHex(), Address: exampleAddr7}, // P2SH(P2WPKH) address from public3
		{PublicKey: public3.ToHex(), Address: exampleAddr9},
		{PublicKey: public3.ToHex(), Address: exampleAddr10},
		{PublicKey: public2.ToHex(), Address: exampleAddr3}, // P2SH(P2PKH) address public2
		{PublicKey: public4.ToHex(), Address: exampleAddr8}, // P2SH(P2PKH) address public2
		{PublicKey: public3.ToHex(), Address: exampleAddr4}, // p2pkh address from public1
	}
	// oh i dont have testnet btc for spending
	// i use https://coinfaucet.eu/en/btc-testnet/ web site for some faucet
	// i wait 10 min for confirmed coinfaucet transaction ...
	// waiting .... i think need more coffe :D

	// ok i got faucet
	// i need now to read spenders account UTXOS
	utxos := provider.UtxoWithOwnerList{}

	// i add some method for provider to read utxos from mempol or blockCypher
	// looping address to read Utxos
	for _, spender := range spenders {
		// read ech address utxo from mempol
		spenderUtxos, err := api.GetAccountUtxo(spender)

		// oh this address does not have any satoshi for spending
		if !spenderUtxos.CanSpending() {
			fmt.Println("address does not have any satoshi for spending: ", spender.Address.Show(network))
			continue
		}
		// oh something bad happen when reading Utxos
		if err != nil {
			fmt.Println("something bad happen when reading Utxos: ", err)
			return
		}
		// we append address utxos to utxos list
		utxos = append(utxos, spenderUtxos...)
	}
	// Well, now we calculate how much we can spend
	sumOfUtxo := utxos.SumOfUtxosValue()

	hasSatoshi := sumOfUtxo.Cmp(big.NewInt(0)) != 0

	if !hasSatoshi {
		// Are you kidding? We don't have btc to spend
		fmt.Println("Are you kidding? We don't have btc to spend")
		return
	}

	fmt.Println("sum of Utxos: ", *sumOfUtxo)
	// 1817320 sum of all utxos

	// We consider 50,000 satoshi for the cost
	// in next example i show you how to calculate fee
	FEE := big.NewInt(50000)

	// now we have 1,767,320 for spending let do it
	// we create 8 different output with  different address type like (pt2r,p2sh(p2wpkh),p2sh(p2wsh),p2sh(p2pkh),p2sh(p2pk),p2pkh,p2wph,p2wsh and etc..)
	// We consider the spendable amount for 10 outputs and divide by 10, each output 176,732
	output1 := provider.BitcoinOutputDetails{
		Address: exampleAddr4,
		Value:   big.NewInt(176732),
	}
	output2 := provider.BitcoinOutputDetails{
		Address: exampleAddr9,
		Value:   big.NewInt(176732),
	}
	output3 := provider.BitcoinOutputDetails{
		Address: exampleAddr10,
		Value:   big.NewInt(176732),
	}
	output4 := provider.BitcoinOutputDetails{
		Address: exampleAddr1,
		Value:   big.NewInt(176732),
	}
	output5 := provider.BitcoinOutputDetails{
		Address: exampleAddr3,
		Value:   big.NewInt(176732),
	}
	output6 := provider.BitcoinOutputDetails{
		Address: exampleAddr2,
		Value:   big.NewInt(176732),
	}
	output7 := provider.BitcoinOutputDetails{
		Address: exampleAddr7,
		Value:   big.NewInt(176732),
	}
	output8 := provider.BitcoinOutputDetails{
		Address: exampleAddr8,
		Value:   big.NewInt(176732),
	}
	output9 := provider.BitcoinOutputDetails{
		Address: exampleAddr5,
		Value:   big.NewInt(176732),
	}
	output10 := provider.BitcoinOutputDetails{
		Address: exampleAddr6,
		Value:   big.NewInt(176732),
	}

	// Well, now it is clear to whom we are going to pay the amount
	// Now let's create the transaction
	transactionBuilder := provider.NewBitcoinTransactionBuilder(
		// Now, we provide the UTXOs we want to spend.
		utxos,
		// We select transaction outputs
		[]provider.BitcoinOutputDetails{output1, output2, output3, output4, output5, output6, output7, output8, output9, output10},
		/*
			Transaction fee
			Ensure that you have accurately calculated the amounts.
			If the sum of the outputs, including the transaction fee,
			does not match the total amount of UTXOs,
			it will result in an error. Please double-check your calculations.
		*/
		FEE,
		// network (address.BitcoinNetwork ,ddress.TestnetNetwork)
		&network,

		// If you like the note write something else and leave it blank
		// I will put my GitHub address here
		"https://github.com/MohsenHaydari",
		/*
			RBF, or Replace-By-Fee, is a feature in Bitcoin that allows you to increase the fee of an unconfirmed
			transaction that you've broadcasted to the network.
			This feature is useful when you want to speed up a
			transaction that is taking longer than expected to get confirmed due to low transaction fees.
		*/
		true,
	)

	// now we use BuildTransaction to complete them
	// I considered a method parameter for this, to sign the transaction

	// utxo Utxo infos with owner details
	// trDigest transaction digest of current UTXO (must be sign with correct privateKey)

	// tweak: cheack is script path spending or tweaking the script.
	// If tweak is set to false, it implies that you are not using the script path spending feature of Taproot,
	// and you intend to sign the transaction using the actual script conditions.

	// sighash
	// Each input in a Bitcoin transaction can include a "sighash type."
	// This type is a flag that determines which parts of the transaction are covered by the digital signature.
	// Common sighash types include SIGHASH_ALL, SIGHASH_SINGLE, SIGHASH_ANYONECANPAY, etc.
	// This TransactionBuilder only works with SIGHASH_ALL and TAPROOT_SIGHASH_ALL for taproot input
	// If you want to use another sighash, you should create another TransactionBuilder
	transaction, err := transactionBuilder.BuildTransaction(func(trDigest []byte, utxo provider.UtxoWithOwner, multiSigPublicKey string) (string, error) {
		var key keypair.ECPrivate

		currentPublicKey := utxo.OwnerDetails.PublicKey
		if utxo.IsMultiSig() {
			currentPublicKey = multiSigPublicKey
		}
		// ok we have the public key of the current UTXO and we use some conditions to find private  key and sign transaction
		switch currentPublicKey {
		case public3.ToHex():
			{
				key = *private3
			}
		case public2.ToHex():
			{
				key = *private2
			}

		case public1.ToHex():
			{
				key = *private1
			}
		case public4.ToHex():
			{
				key = *private4
			}
		default:
			{
				return "", fmt.Errorf("cannot find private key")
			}
		}
		// Ok, now we have the private key, we need to check which method to use for signing
		// We check whether the UTX corresponds to the P2TR address or not.
		if utxo.Utxo.IsP2tr() {
			// yes is p2tr utxo and now we use SignTaprootTransaction(Schnorr sign)
			// for now this transaction builder support only tweak transaction
			return key.SignTaprootTransaction(
				trDigest, constant.TAPROOT_SIGHASH_ALL, []interface{}{}, true,
			), nil
		}
		// is seqwit(v0) or lagacy address we use  SingInput (ECDSA)
		return key.SingInput(trDigest, constant.SIGHASH_ALL), nil

	})

	if err != nil {
		fmt.Println("oh we have some error when build and sign transaction ", err)
		return
	}
	// ok everything is fine and we need a transaction output for broadcasting
	// We use the Serialize method to receive the transaction output
	_ = transaction.Serialize()

	// we check if transaction is segwit or not
	// When one of the input UTXO addresses is SegWit, the transaction is considered SegWit.
	isSegwitTr := transactionBuilder.HasSegwit()

	// transaction id
	transactionId := transaction.TxId()
	fmt.Println("transaction ID: ", transactionId)

	// transaction size
	var transactionSize int

	if isSegwitTr {
		transactionSize = transaction.GetVSize()
	} else {
		transactionSize = transaction.GetSize()
	}
	// 960 byte
	fmt.Println("transaction size: ", transactionSize)
	// Ok we now have the transaction size
	// Well, we use API to receive network fees
	networkFee, err := api.GetNetworkFee()
	if err != nil {
		fmt.Println("cannot read network fee: ", err)
		return
	}
	// We use the Get Estimate method to receive the transaction fee
	transactionFee := networkFee.GetEstimate(transactionSize, networkFee.Medium)
	fmt.Println("transaction fee: ", transactionFee)
	// 30952 statoshi
}
