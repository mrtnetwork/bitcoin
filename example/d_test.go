// Examples, of how to use the package
package example

import (
	"testing"

	"github.com/mrtnetwork/bitcoin/address"
	"github.com/mrtnetwork/bitcoin/provider"

	"errors"
	"fmt"
	"math/big"

	"github.com/mrtnetwork/bitcoin/constant"
	hdwallet "github.com/mrtnetwork/bitcoin/hd_wallet"
	"github.com/mrtnetwork/bitcoin/keypair"
)

func TestD(t *testing.T) {
	network := address.TestnetNetwork
	api := provider.SelectApi(provider.BlockCyperApi, &network)
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

	signer1, _ := provider.CreateMultiSignaturSigner(
		// public key of signer
		public1.ToHex(),
		// siger weight
		2,
	)
	signer2, _ := provider.CreateMultiSignaturSigner(
		public2.ToHex(),
		2,
	)
	signer3, _ := provider.CreateMultiSignaturSigner(
		public3.ToHex(),
		1,
	)
	signer4, _ := provider.CreateMultiSignaturSigner(
		public4.ToHex(),
		1,
	)

	/*
		In general, this address requires 5 signatures to spend:
		2 signatures from signer1
		2 signatures from signer2
		and 1 signature from either signer 3 or signer 4.
		And the address script is as follows

		["OP_5", public1 ,public1 ,public2 ,public2 ,public3 ,public4, "OP_6", "OP_CHECKMULTISIG"]

		And the unlock script will be like this

		["", signer1Signataure, signer1Signataure, signer2Signatur, signer2Signatur, (signer3Signatur or signer4Signatur), ScriptInHex ]

	*/
	multiSigBuilder, err := provider.CreateMultiSignatureAddress(
		5, provider.MultiSignaturAddressSigners{
			signer1,
			signer2, signer3, signer4,
		}, address.P2WSHInP2SH, // P2SH(P2WSH)
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	/*
		In general, this address requires 5 signatures to spend:
		2 signatures from signer1
		2 signatures from signer2
		and 1 signature from either signer 3 or signer 4.
		And the address script is as follows

		["OP_5", public1 ,public1 ,public2 ,public2 ,public3 ,public4, "OP_6", "OP_CHECKMULTISIG"]

		And the unlock script will be like this

		["", signer1Signataure, signer1Signataure, signer2Signatur, signer2Signatur, (signer3Signatur or signer4Signatur), ScriptInHex ]
	*/
	multiSigBuilder2, err2 := provider.CreateMultiSignatureAddress(
		5, provider.MultiSignaturAddressSigners{
			signer1,
			signer2, signer3, signer4,
		}, address.P2WSH, // P2WSH
	)
	if err2 != nil {
		fmt.Println(err2)
		return
	}
	// P2SH(P2WSH) 5-6 multi-sig ADDRESS
	// 2MxVXBKFwvkWFeN4nij3n8s2GMeBeqF6cL4
	multiSigAddress := multiSigBuilder.Address

	// P2SH(P2WPKH) 5-6 multi-sig ADDRESS
	// tb1q4aw8qjc4eys27y8hnslzqexkgs920ewx8ssuxhwq0sc28vly0w0sv3mvu9
	multiSigAddress2 := multiSigBuilder2.Address

	// P2TR ADDRESS
	// tb1pyhmqwlcrws4dxcgalt4mrffgnys879vs59xf6sve4hazyvmhecxq3e6sc0
	// equals to exampleAddr1 := address.P2TRAddressFromAddress("tb1pyhmqwlcrws4dxcgalt4mrffgnys879vs59xf6sve4hazyvmhecxq3e6sc0")
	exampleAddr1 := public2.ToTaprootAddress()

	// P2SH(P2PK) ADDRESS
	// 2MugsNcgzLJ1HosnZyC2CfZVmgbMPK1XubR
	// equals to exampleAddr2 := address.P2SHAddressFromAddress("2MugsNcgzLJ1HosnZyC2CfZVmgbMPK1XubR", address.P2PKInP2SH)
	exampleAddr2 := public4.ToP2PKInP2SH()

	// P2WSH ADDRESS
	// tb1qf4qwtr5kp5q87dtp3ul3402vkzssxfv7f4aettjq2hcfhnt92dmq5xzs6n
	// equals to exampleAddr3 := address.P2WSHAddresssFromAddress("tb1qf4qwtr5kp5q87dtp3ul3402vkzssxfv7f4aettjq2hcfhnt92dmq5xzs6n")
	// created with 1-1 MultiSig script: ["OP_1", publicHex(Hex of compressed public key) , "OP_1", "OP_CHECKMULTISIG"]
	exampleAddr3 := public3.ToP2WSHAddress()

	// now we chose some address for spending from multiple address
	// i use some different address type for this
	spenders := []provider.UtxoOwnerDetails{
		{Address: multiSigAddress, MultiSigAddress: multiSigBuilder},
		{Address: multiSigAddress2, MultiSigAddress: multiSigBuilder2},
		{PublicKey: public2.ToHex(), Address: exampleAddr1},
	}

	// now we need to read spenders account UTXOS
	utxos := provider.UtxoWithOwnerList{}

	// i add some method for provider to read utxos from mempol or blockCypher
	// looping address to read Utxos
	for _, spender := range spenders {
		// read ech address utxo from mempol
		spenderUtxos, err := api.GetAccountUtxo(spender)
		// oh something bad happen when reading Utxos
		if err != nil {
			fmt.Println("something bad happen when reading Utxos: ", err)
			return
		}
		// oh this address does not have any satoshi for spending
		if !spenderUtxos.CanSpending() {
			fmt.Println("address does not have any satoshi for spending: ", spender.Address.Show(network))
			continue
		}
		fmt.Println("spending: ", spenderUtxos.SumOfUtxosValue(), spender.Address.Show(network), spender.Address.GetType())

		// we append address utxos to utxos list
		utxos = append(utxos, spenderUtxos...)

	}
	// Well, now we calculate how much we can spend
	sumOfUtxo := utxos.SumOfUtxosValue()
	fmt.Println("sum of utxos: ", sumOfUtxo)
	hasSatoshi := sumOfUtxo.Cmp(big.NewInt(0)) != 0

	if !hasSatoshi {
		// Are you kidding? We don't have btc to spend
		fmt.Println("Are you kidding? We don't have btc to spend")
		return
	}

	fmt.Println("sum of Utxos: ", *sumOfUtxo)
	// 656,928 sum of all utxos

	// We consider 50,000 satoshi for the cost
	// in next example i show you how to calculate fee
	FEE := big.NewInt(3000)

	// now we have 606,920 for spending let do it
	// we create 5 different output with  different address type
	// We consider the spendable amount for 5 outputs and divide by 5, each output 121,384

	output3 := provider.BitcoinOutputDetails{
		Address: exampleAddr3,
		Value:   big.NewInt(236768),
	}
	output4 := provider.BitcoinOutputDetails{
		Address: exampleAddr2,
		Value:   big.NewInt(1000),
	}
	output5 := provider.BitcoinOutputDetails{
		Address: exampleAddr1,
		Value:   big.NewInt(1000),
	}
	output6 := provider.BitcoinOutputDetails{
		Address: multiSigAddress,
		Value:   big.NewInt(1000),
	}
	output7 := provider.BitcoinOutputDetails{
		Address: multiSigAddress2,
		Value:   big.NewInt(1000),
	}

	// Well, now it is clear to whom we are going to pay the amount
	// Now let's create the transaction
	transactionBuilder := provider.NewBitcoinTransactionBuilder(
		// Now, we provide the UTXOs we want to spend.
		utxos,
		// We select transaction outputs
		[]provider.BitcoinOutputDetails{output3, output4, output5, output6, output7},
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
				return "", errors.New("cannot find private key")
			}
		}
		// Ok, now we have the private key, we need to check which method to use for signing
		// We check whether the UTX corresponds to the P2TR address or not.
		if utxo.Utxo.IsP2tr() {
			// yes is p2tr utxo and now we use SignTaprootTransaction(Schnorr sign)
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
	digest := transaction.Serialize()

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
	fmt.Println("transaction size: ", transactionSize)

	// now we send transaction to network
	trId, err := provider.TestMempoolAccept([]interface{}{[]string{digest}})

	if err != nil {
		fmt.Println("something bad happen when sending transaction: ", err)
		return
	}
	// Yes, we did :)  72b7244693960879bb07f9f96e87790a8b57bb2e91c8dfd79e6f9b8ee520adff
	// Now we check Mempol for what happened https://mempool.space/testnet/tx/72b7244693960879bb07f9f96e87790a8b57bb2e91c8dfd79e6f9b8ee520adff
	fmt.Println("Transaction ID: ", trId)
}
