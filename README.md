# BITCOIN GO Package
a comprehensive and versatile Go library for all your Bitcoin transaction needs. offers robust support for various Bitcoin transaction types, including spending transactions, Bitcoin address management, Bitcoin Schnorr signatures, BIP-39 mnemonic phrase generation, hierarchical deterministic (HD) wallet derivation, and Web3 Secret Storage Definition.

## Features

### Transaction Types
This comprehensive package provides robust support for a wide array of Bitcoin transaction types, encompassing the full spectrum of Bitcoin transaction capabilities. Whether you need to execute standard payments, facilitate complex multi-signature wallets, leverage Segregated Witness (SegWit) transactions for lower fees and enhanced scalability, or embrace the privacy and flexibility of Pay-to-Taproot (P2TR) transactions, this package has you covered. Additionally, it empowers users to engage in legacy transactions, create time-locked transactions, and harness the security of multisignature (multisig) transactions. With this package, you can seamlessly navigate the diverse landscape of Bitcoin transactions, ensuring your interactions with the Bitcoin network are secure, efficient, and tailored to your specific needs.

- P2PKH (Pay-to-Public-Key-Hash): The most common transaction type, it sends funds to a recipient's public key hash. Provides security and anonymity.

- P2SH (Pay-to-Script-Hash): Allows more complex scripts to be used, enhancing Bitcoin's capabilities by enabling features like multisignature wallets.

- P2WPKH (Pay-to-Witness-Public-Key-Hash): A Segregated Witness (SegWit) transaction type, it offers reduced fees and improved scalability while maintaining compatibility.

- P2WSH (Pay-to-Witness-Script-Hash): Another SegWit transaction, it extends the benefits of SegWit to more complex script scenarios, reducing transaction size and fees.

- P2TR (Pay-to-Taproot): An upgrade aiming to improve privacy and flexibility, allowing users to choose between various scripts and enhance transaction efficiency.

- Legacy Transactions: Refers to older transaction types used before SegWit, with higher fees and less scalability.

- Multisignature (Multisig) Transactions: Involves multiple signatures to authorize a Bitcoin transaction, commonly used for security purposes in shared wallets.

- SegWit Transactions: A collective term for P2WPKH and P2WSH transactions, leveraging segregated witness data to reduce transaction size and fees.

- Time-Locked Transactions: These transactions have a predetermined time or block height before they can be spent, adding security and functionality to Bitcoin smart contracts.

- Coinbase Transactions: The first transaction in each block, generating new Bitcoins as a block reward for miners. It includes the miner's payout address.

### Create Transaction
Using this package, you can create a Bitcoin transaction in two ways: either through the `BtcTransaction` struct or the `BitcoinTransactionBuilder` struct
- BtcTransaction: To use the `BtcTransaction` struct, you should have a general understanding of how Bitcoin transactions work, including knowledge of UTXOs, scripts, various types of scripts, Bitcoin addresses, signatures, and more. We created examples and tests to enhance your understanding. An example of this transaction type is explained below, and you can also find numerous examples in the [`test`](https://github.com/MohsenHaydari/bitcoin/tree/main/test) folder.

- BitcoinTransactionBuilder: Even with limited prior knowledge, you can utilize this class to send various types of transactions. Below, I've provided an example in which a transaction features 8 distinct input addresses with different types and private keys, as well as 10 different output addresses. Furthermore, additional examples have been prepared, which you can find in the [`example`](https://github.com/MohsenHaydari/bitcoin/tree/main/example) folder.

### Addresses
- P2PKH A P2PKH (Pay-to-Public-Key-Hash) address in Bitcoin represents ownership of a cryptocurrency wallet by encoding a hashed public key
  
- P2WPKH: A P2WPKH (Pay-to-Witness-Public-Key-Hash) address in Bitcoin is a Segregated Witness (SegWit) address that enables more efficient and secure transactions by segregating witness data, enhancing network scalability and security.
  
- P2WSH: A P2WSH (Pay-to-Witness-Script-Hash) address in Bitcoin is a Segregated Witness (SegWit) address that allows users to spend bitcoins based on the conditions specified in a witness script, offering improved security and flexibility for complex transaction types.
  
- P2TR: A P2TR (Pay-to-Taproot) address in Bitcoin is a type of address that allows users to send and receive bitcoins using the Taproot smart contract, offering enhanced privacy and scalability features.
  
- P2SH: A P2SH (Pay-to-Script-Hash) address in Bitcoin is an address type that enables the use of more complex scripting, often associated with multi-signature transactions or other advanced smart contract functionality, enhancing flexibility and security.
  
- P2SH(SEGWIT): A P2SH (Pay-to-Script-Hash) Segregated Witness (SegWit) address in Bitcoin combines the benefits of P2SH and SegWit technologies, allowing for enhanced transaction security, reduced fees, and improved scalability.

### Sign
- Sign message: ECDSA Signature Algorithm
  
- Sign Segwit(v0) and legacy transaction: ECDSA Signature Algorithm
  
- Sign Taproot transaction
  
  - Script Path and TapTweak: Taproot allows for multiple script paths (smart contract conditions) to be included in a single transaction. The "taptweak" ensures that the correct 	 
    script path is used when spending. This enhances privacy by making it difficult to determine the spending conditions from the transaction.
    
  - Schnorr Signatures: While ECDSA is still used for Taproot, it also provides support for Schnorr signatures. Schnorr signatures offer benefits such as smaller signature sizes and 	 
    signature aggregation, contributing to improved scalability and privacy.
    
  - Schnorr-Musig: Taproot can leverage Schnorr-Musig, a technique for securely aggregating multiple signatures into a single signature. This feature enables collaborative spending and 
    enhances privacy.

### BIP-39
- Generate BIP39 mnemonics, providing a secure and standardized way to manage keys and seed phrases

### HD Wallet
- Implement hierarchical deterministic (HD) wallet derivation

### Web3 Secret Storage Definition
- JSON Format: Private keys are stored in a JSON (JavaScript Object Notation) format, making it easy to work with in various programming languages.
  
- Encryption: The private key is encrypted using the user's chosen password. This ensures that even if the JSON file is compromised, an attacker cannot access the private key without the password.
  
- Key Derivation: The user's password is typically used to derive an encryption key using a key derivation function (KDF). This derived key is then used to encrypt the private key.
  
- Scrypt Algorithm: The Scrypt algorithm is commonly used for key derivation, as it is computationally intensive and resistant to brute-force attacks.
  
- Checksum: A checksum is often included in the JSON file to help detect errors in the password.
  
- Initialization Vector (IV): An IV is typically used to add an extra layer of security to the encryption process.
  
- Versioning: The JSON file may include a version field to indicate which version of the encryption and storage format is being used.
  
- Metadata: Additional metadata, such as the address associated with the private key, may be included in the JSON file.

### Node Provider
We have added two APIs (Mempool and BlockCypher) to the plugin for network access. You can easily use these two APIs to obtain information such as unspent transactions (UTXO), network fees, sending transactions, receiving transaction information, and retrieving account transactions.

## EXAMPLES

### Key and addresses
  - Private key
    ```
    // Create an EC private key instance from a WIF (Wallet Import Format) encoded string.
    privateKey, _ := keypair.NewECPrivateFromWIF("cT33CWKwcV8afBs5NYzeSzeSoGETtAB8izjDjMEuGqyqPoF7fbQR")

    // Retrieve the corresponding public key from the private key.
    publicKey := privateKey.GetPublic()

    // Sign an input using the private key.
    signSegwitV0OrLagacy := privateKey.SignInput()

    // Sign a Taproot transaction using the private key.
    signSegwitV1TapprotTransaction := privateKey.SignTaprootTransaction()

    // Convert the private key to a WIF (Wallet Import Format) encoded string.
    // The boolean argument specifies whether to use the compressed format.
    toWif := privateKey.ToWIF(true, &address.MainnetNetwork)

    // Convert the private key to its hexadecimal representation.
    toHex := privateKey.ToHex()
    ```
- Public key
  ```
	// Create an instance of an EC public key from a hexadecimal representation.
	publicKey, _ := keypair.NewECPPublicFromHex()

	// Generate a Pay-to-Public-Key-Hash (P2PKH) address from the public key.
	p2pkh := publicKey.ToAddress()

	// Generate a Pay-to-Witness-Public-Key-Hash (P2WPKH) Segregated Witness (SegWit) address from the public key.
	p2wpkh := publicKey.ToSegwitAddress()

	// Generate a Pay-to-Witness-Script-Hash (P2WSH) Segregated Witness (SegWit) address from the public key.
	p2wsh := publicKey.ToP2WSHAddress()

	// Generate a Taproot address from the public key.
	p2tr := publicKey.ToTaprootAddress()

	// Generate a Pay-to-Public-Key-Hash (P2PKH) inside Pay-to-Script-Hash (P2SH) address from the public key.
	p2pkhInP2sh := publicKey.ToP2PKHInP2SH()

	// Generate a Pay-to-Witness-Public-Key-Hash (P2WPKH) inside Pay-to-Script-Hash (P2SH) address from the public key.
	p2wpkhInP2sh := publicKey.ToP2WPKHInP2SH()

	// Generate a Pay-to-Witness-Script-Hash (P2WSH) inside Pay-to-Script-Hash (P2SH) address from the public key.
	p2wshInP2sh := publicKey.ToP2WSHInP2SH()

	// Generate a Pay-to-Public-Key (P2PK) inside Pay-to-Script-Hash (P2SH) address from the public key.
	p2pkInP2sh := publicKey.ToP2PKInP2SH()

	// Get the compressed bytes representation of the public key.
	compressedBytes := publicKey.ToCompressedBytes()

	// Get the uncompressed bytes representation of the public key.
	unCompressedBytes := publicKey.ToUnCompressedBytes(true)

	// extracts and returns the x-coordinate (first 32 bytes) of the ECPublic key
	// as a hexadecimal string.
	onlyX := publicKey.ToXOnlyHex()

	// CalculateTweak computes and returns the TapTweak value based on the ECPublic key
	// and an optional script. It uses the key's x-coordinate and the Merkle root of the script
	// (if provided) to calculate the tweak.
	tweak, _ := publicKey.CalculateTweak(scripts.Script{})
	
	// computes and returns the Taproot commitment point's x-coordinate
	// derived from the ECPublic key and an optional script, represented as a hexadecimal string.
	taproot, _ := publicKey.ToTapRotHex([]interface{}{})

	// Verify verifies a signature against a message
	verify := publicKey.Verify()
  ```
- Addresses
  ```
  // The `...FromAddress` methods only check the address and return an error if it does not belong to Bitcoin.
  // If you also want to verify that the address belongs to a specific network,
  // please select the desired network using the parameters.

  // Generate a Pay-to-Public-Key-Hash (P2PKH) address from the public key.
  p2pkh, _ := address.P2PKHAddressFromAddress("1Q5odQtVCc4PDmP5ncrp7DSuVbh2ML4Gnb", network)

  // Generate a Pay-to-Witness-Public-Key-Hash (P2WPKH) Segregated Witness (SegWit) address from the public key.
  p2wpkh, _ := address.P2WPKHAddresssFromAddress("bc1ql5eh45als8sgdkt2drsl344q55g03sj2u9enzz")

  // Generate a Pay-to-Witness-Script-Hash (P2WSH) Segregated Witness (SegWit) address from the public key.
  p2wsh, _ := address.P2WSHAddresssFromAddress("bc1qf90kcg2ktg0wm983cyvhy0jsrj2fmqz26ugf5jz3uw68mtnr8ljsnf8pqe")

  // Generate a Taproot address from the public key.
  p2tr, _ := address.P2TRAddressFromAddress("bc1pmelvn3xz2n3dmcsvk2k99na7kc55ry77zmhg4z39upry05myjthq37f6jk")

  // Generate a Pay-to-Public-Key-Hash (P2PKH) inside Pay-to-Script-Hash (P2SH) address from the public key.
  p2pkhInP2sh, _ := address.P2SHAddressFromAddress("3HDtvvRMu3yKGFXYFSubTspbhbLagpdKJ7")

  // Generate a Pay-to-Witness-Public-Key-Hash (P2WPKH) inside Pay-to-Script-Hash (P2SH) address from the public key.
  p2wpkhInP2sh, _ := address.P2SHAddressFromAddress("36Dq32LRMW8EJyD3T2usHaxeMBmUpsXhq2")

  // Generate a Pay-to-Witness-Script-Hash (P2WSH) inside Pay-to-Script-Hash (P2SH) address from the public key.
  p2wshInP2sh, _ := address.P2SHAddressFromAddress("3PPL49fMytbEKJsjjPnkfWh3iWzrZxQZAg")

  // Generate a Pay-to-Public-Key (P2PK) inside Pay-to-Script-Hash (P2SH) address from the public key.
  p2pkInP2sh, _ := address.P2SHAddressFromAddress("3NCe6AGzjz2jSyRKCc8o3Bg5MG6pUM92bg")

  // You can create any type of Bitcoin address with scripts.
  // Create an address with scripts for P2WSH multisig 3-of-5.
  newScript := scripts.NewScript("OP_3", publicKey.ToHex(true), publicKey.ToHex(true), publicKey.ToHex(true), publicKey.ToHex(true), publicKey.ToHex(true), "OP_5", "OP_CHECKMULTISIG")

  // Generate a P2WSH 3-of-5 address.
  p2wsh3of5Address, _ := address.P2WSHAddresssFromScript(newScript)

  // Generate a P2SH 3-of-5 address from the P2WSH address.
  p2sh3Of5, _ := address.P2SHAddressFromScript(p2wsh3of5Address.ToScriptPubKey(), address.P2WPKHInP2SH)

  // The method calculates the address checksum and returns the Base58-encoded
  // Bitcoin legacy address or the Bech32 format for SegWit addresses.
  p2sh3Of5.Show(network)

  // Return the scriptPubKey that corresponds to this address.
  p2sh3Of5.ToScriptPubKey()

  // Access the legacy or SegWit program of the address.
  p2sh3Of5.Program()
  ```
  
### Transaction
- With TransactionBuilder
  ```
  // spending from 3 private key with 8 different address to 10 address
  // network
  network := address.TestnetNetwork
  // create node provider
  api := provider.SelectApi(provider.MempoolApi, &network)

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
  exampleAddr1 := public1.ToAddress()

  // P2TR ADDRESS
  exampleAddr2 := public2.ToTaprootAddress()

  // P2SH(P2PKH) ADDRESS
  exampleAddr3 := public2.ToP2PKHInP2SH()

  // P2PKH ADDRESS
  exampleAddr4 := public3.ToAddress()

  // P2SH(P2PKH) ADDRESS
  exampleAddr5 := public3.ToP2PKHInP2SH()

  // P2SH(P2WSH) ADDRESS
  exampleAddr6 := public3.ToP2WSHInP2SH()

  // P2SH(P2WPKH) ADDRESS
  exampleAddr7 := public3.ToP2WPKHInP2SH()

  // P2SH(P2PK) ADDRESS
  exampleAddr8 := public4.ToP2PKInP2SH()

  // P2WPKH ADDRESS
  exampleAddr9 := public3.ToSegwitAddress()

  // P2WSH ADDRESS
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
  
  utxos := provider.UtxoWithOwnerList{}

  // i add some method for provider to read utxos from mempol or blockCypher
  // looping address to read Utxos
  for _, spender := range spenders {
		// read ech address utxo from mempol
		spenderUtxos, err := api.GetAccountUtxo(spender)

		// oh this address does not have any satoshi for spending
		if !spenderUtxos.CanSpending() {
			continue
		}
		// oh something bad happen when reading Utxos
		if err != nil {
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
  	return
  }

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

  	// Transaction fee
  	// Ensure that you have accurately calculated the amounts.
  	// If the sum of the outputs, including the transaction fee,
  	// does not match the total amount of UTXOs,
  	// it will result in an error. Please double-check your calculations.

  	FEE,
  	// network (address.BitcoinNetwork ,ddress.TestnetNetwork)
  	&network,

  	// If you like the note write something else and leave it blank
  	// I will put my GitHub address here
  	"https://github.com/MohsenHaydari",
  
  	// RBF, or Replace-By-Fee, is a feature in Bitcoin that allows you to increase the fee of an unconfirmed
  	// transaction that you've broadcasted to the network.
  	// This feature is useful when you want to speed up a
  	// transaction that is taking longer than expected to get confirmed due to low transaction fees.
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
  
  // transaction size
  var transactionSize int

  if isSegwitTr {
  	transactionSize = transaction.GetVSize()
  } else {
  	transactionSize = transaction.GetSize()
  }
    
  // now we send transaction to network
  trId, err := api.SendRawTransaction(digest)

  if err != nil {
  	return
  }
  // Yes, we did :)  5015a7748d8d6df47358902b6cdc6d77ef839945c479924f4592fd89315ac0e0
  // Now we check Mempol for what happened https://mempool.space/testnet/tx/5015a7748d8d6df47358902b6cdc6d77ef839945c479924f4592fd89315ac0e0

  ```
- With BtcTransaction
  - Spend P2TR UTXO
    ```
    // Private key of the UTXO owner
    privateKey, _ := keypair.NewECPrivateFromWIF("cRvyLwCPLU88jsyj94L7iJjQX5C2f8koG4G2gevN4BeSGcEvfKe9")

    // Address we want to spend from
    fromAddr := privateKey.GetPublic().ToTaprootAddress()

    // Create an input
    // Insert transaction ID and index of UTXO
    sigTxin1 := scripts.NewTxInput("3d4c9d73c4c65772e645ff26493590ae4913d9c37125b72398222a553b73fa66", 0)

    // Address we want to send funds to
    addr, _ := address.P2PKHAddressFromAddress("n4bkvTyU1dVdzsrhWBqBw8fEMbHjJvtmJR")

    // Create an output: Send 3000 to `n4bkvTyU1dVdzsrhWBqBw8fEMbHjJvtmJR`
    txout := scripts.NewTxOutput(
		big.NewInt(3000),
		addr.ToScriptPubKey(),
	)

    // Create a transaction
    tx := scripts.NewBtcTransaction(
		[]*scripts.TxInput{sigTxin1},
		[]*scripts.TxOutput{txout},
		true, // The transaction contains one or more segwit UTXOs
	)

    // Get the transaction digest for signing input at index 0
    // Arguments:
    // - Index 0
    // - The scriptPubkeys that correspond to all the inputs/UTXOs
    // - The amounts that correspond to all the inputs/UTXOs
    // - Ext_flag: Extension mechanism, default is 0; 1 is for script spending (BIP342)
    // - Script: The script that we are spending when ext_flag is 1
    digest := tx.GetTransactionTaprootDigest(0, []*scripts.Script{fromAddr.ToScriptPubKey()}, []*big.Int{big.NewInt(3500)}, 0, scripts.NewScript(), constant.TAPROOT_SIGHASH_ALL)

    // Sign the transaction
    // Arguments:
    // - Transaction digest related to the index
    // - Signature Hash Type (TAPROOT_SIGHASH_ALL)
    // - Script path (tapleafs)
    // - Tweak: Note that we don't use tapleafs script in this transaction, so the tweak should be set to False
    sig := privateKey.SignTaprootTransaction(digest, constant.TAPROOT_SIGHASH_ALL, []interface{}{}, false)

    // Create a witness signature and set it to the transaction at the current index
    witness := scripts.NewTxWitnessInput(sig)
    tx.Witnesses = append(tx.Witnesses, witness)

    // Transaction ID
    tx.TxId()

    // In this case, the transaction is segwit, and we must use GetVSize for transaction size
    tx.GetVSize()

    // Transaction digest ready for broadcast
    tx.Serialize()
    
    ```
  - Spend P2PKH UTXO
    ```
    // private key of UTXO owner
    privateKey, _ := keypair.NewECPrivateFromWIF("...")
    // address we want to spend
    fromAddr := privateKey.GetPublic().ToP2PKAddress()

    // create input
    // insert transaction ID and index of UTXO
    sigTxin1 := scripts.NewTxInput("3d4c9d73c4c65772e645ff26493590ae4913d9c37125b72398222a553b73fa66", 0)

    // address we want to send
    addr, _ := address.P2PKHAddressFromAddress("n4bkvTyU1dVdzsrhWBqBw8fEMbHjJvtmJR")

    // outputs 1 send 10000000 to `n4bkvTyU1dVdzsrhWBqBw8fEMbHjJvtmJR`
    txout := scripts.NewTxOutput(
		big.NewInt(3000),
		addr.ToScriptPubKey())

    // We will return the spent transaction balance to the account holder
    changeTxout := scripts.NewTxOutput(
		big.NewInt(29000000),
		fromAddr.ToScriptPubKey(),
    )

    // create transaction
    tx := scripts.NewBtcTransaction(
		[]*scripts.TxInput{sigTxin1},
		[]*scripts.TxOutput{txout, changeTxout},
		false)

    // get transaction digest for signing input at index 0
    digest := tx.GetTransactionDigest(0, fromAddr.Program().ToScriptPubKey())

    // sign the transaction
    sig := privateKey.SingInput(digest)
    
    // set unlock script to transaction
    tx.SetScriptSig(0, scripts.NewScript(sig, privateKey.GetPublic().ToHex()))
    // transaction id
    tx.TxId()
    // transaction size
    tx.GetSize()
    // transaction digest ready for broadcast
    tx.Serialize()
    ```
  

### BIP39
```
// Create a new Bip39 instance with the desired language.
// The default language is English. You can choose from various supported languages.
// Language options include: English, Spanish, Portuguese, Korean, Japanese, Italian, French, Czech, ChineseTraditional, and ChineseSimplified.
bip := bip39.Bip39{
	Language: bip39.Japanese, // Set the language to Japanese
}

// Generate a mnemonic phrase with the specified number of words.
// You can choose from 12, 15, 18, 21, or 24 words (e.g., bip39.Words24).
mnemonic, err := bip.GenerateMnemonic(bip39.Words24)

// Derive a seed from the generated mnemonic phrase.
// You can optionally provide a passphrase, which can be an empty string.
toSeed := bip39.ToSeed(mnemonic, "PASSPHRASE")

// Convert the mnemonic phrase back to entropy.
toEntropy, err := bip.MnemonicToEntropy(mnemonic)

// Convert the entropy back to a mnemonic phrase.
toMnemonicFromEntropy, err := bip.EntropyToMnemonic(toEntropy)

// Change the language of the Bip39 instance to Italian.
// You can use the `bip.ChangeLanguage()` method to switch between supported languages.
bip.ChangeLanguage(bip39.Italian)

```
### HD Wallet
```
// Create a pointer to the Mainnet network
network := &address.MainnetNetwork

// Initialize a Bip39 instance with the Japanese language
bip39Instance := bip39.Bip39{
	Language: bip39.Japanese,
}

// Generate a 24-word mnemonic using the Bip39 instance
mn, _ := bip39Instance.GenerateMnemonic(bip39.Words24)

// Create a master wallet from the generated mnemonic and a passphrase
masterWallet, _ := hdwallet.FromMnemonic(mn, "MYPASSPHRASE")

// Derive an extended private key (xPrive) for the specified address type and network
xPrive := masterWallet.ToXPrivateKey(address.P2PKH, network)

// Derive an extended public key (xPub) for the specified address type and network
xPub := masterWallet.ToXPublicKey(address.P2PKH, network)

// Create a master public wallet from the derived xPub, specifying SegWit support and the network
masterPublicWallet, _ := hdwallet.FromXPublicKey(xPub, true, network)

// Derive a child wallet (driveFromPublicWallet) from the master public wallet using a specific derivation path
driveFromPublicWallet, _ := hdwallet.DrivePath(masterPublicWallet, "m/0/1")

// Get the public key of the child wallet (driveFromPublicWallet)
publicKeyOfMasterPublicWallet := driveFromPublicWallet.GetPublic()

// Derive a child wallet (drivePrivateWallet) from the master public wallet using a different derivation path
drivePrivateWallet, _ := hdwallet.DrivePath(masterPublicWallet, "m/44'/0'/0'/1")

// Derive an extended private key from the child wallet (drivePrivateWallet) for the specified address type and network
_ = drivePrivateWallet.ToXPrivateKey(address.P2PKH, network)

// Derive an extended public key from the child wallet (drivePrivateWallet) for the specified address type and network
_ = drivePrivateWallet.ToXPublicKey(address.P2PKH, network)

// Get the public key of the child wallet (drivePrivateWallet)
publicKey := drivePrivateWallet.GetPublic()

// Get the private key of the child wallet (drivePrivateWallet)
privateKey, _ := drivePrivateWallet.GetPrivate()

```
### Web3 Secret Storage Definition
```
// Create a pointer to the Mainnet network
network := &address.MainnetNetwork

// Initialize a Bip39 instance with the Japanese language
bip39Instance := bip39.Bip39{
	Language: bip39.Japanese,
}

// Generate a 24-word mnemonic using the Bip39 instance
mn, _ := bip39Instance.GenerateMnemonic(bip39.Words24)

// Define your password
myPassword := "myPassword"

// Create a master wallet from the generated mnemonic and a passphrase
materWallet, _ := hdwallet.FromMnemonic(mn, "MYPASSPHRASE")

// Derive an extended private key (xPrive) for the specified address type and network
xPrive := materWallet.ToXPrivateKey(address.P2PKH, network)

// Define Scrypt parameters (N and p)
scryptN := 8192
p := 1

// Create an encrypted wallet using the xPrive key, password, Scrypt parameters, and p
encryptWallet, _ := secretwallet.NewSecretWallet(xPrive, myPassword, scryptN, p)

// Convert the encrypted wallet to a JSON representation
encrypted, _ := encryptWallet.ToJSON()

// Decode the encrypted wallet using the password
decodedWallet, _ := secretwallet.DecodeSecretWallet(encrypted, myPassword)

// Create a new wallet from the decoded credentials, specifying SegWit support and the network
newWallet, _ := hdwallet.FromXPrivateKey(decodedWallet.Credentials, true, network)

```
### Node provider
```
// select network testnet or mainnet
network := address.TestnetNetwork

// create api (BlockCyperApi or MempoolApi)
// Currently, only a few critical methods have been implemented to retrieve unspent transactions,
// obtain network fees, receive transactions, and send transactions to the network.
api := provider.SelectApi(provider.MempoolApi, &network)

// Read Transaction id(hash)
tr, e := api.GetTransaction("d4bad8e07d30ca4389ec8a203318aa523cc3e36c9730d0a6852a3801d086c5fe")

// Read accounts UTXOS
addr, _ := address.P2WPKHAddresssFromAddress("tb1q92nmnvhj04sqd4x7wjaewlt5jn8n3ngmplcymy")
utxos, e := api.GetAccountUtxo(provider.UtxoOwnerDetails{
	PublicKey: "",
	Address:   addr,
})

// Network fee
fee, e := api.GetNetworkFee()

//  Send transaction
_, _ = api.SendRawTransaction("TRANSACTION DIGEST")

// Read account transactions
transaction, _ := api.GetAccountTransactions(addr.Show(network), func(url string) string {
	return url
})

```

## Contributing

Contributions are welcome! Please follow these guidelines:
 - Fork the repository and create a new branch.
 - Make your changes and ensure tests pass.
 - Submit a pull request with a detailed description of your changes.

## Feature requests and bugs #

Please file feature requests and bugs in the issue tracker.


## Support

If you find this repository useful, show us some love, and give us a star!
Small Bitcoin donations to the following address are also welcome:

bc1q2y5rcf6uwmg57zsmqy86xsc0pt0rw7qujn2a2h

18Q7w7aQYSAQonWTr2Uj5Z7VT7dNTZFwqk

