# XRP Dart Package

This package provides functionality to sign XRP transactions using two popular cryptographic algorithms, 
ED25519 and SECP256K1. It allows developers to create and sign XRP transactions securely.

## Features

### Transaction Types
The XRP Ledger supports various transaction types, each serving a different purpose:

- Payment Transactions: Standard transactions used to send XRP or Issue from one address to another.
- Escrow Transactions: These transactions lock up XRP until certain conditions are met, providing a trustless way to facilitate delayed payments.
- TrustSet Transactions: Used to create or modify trust lines, enabling users to hold and trade assets other than XRP on the ledger.
- OrderBook Transactions: Used to place and cancel offers on the decentralized exchange within the XRP Ledger.
- PaymentChannel Transactions: Allow for off-chain payments using payment channels.
- NFT: Mint NFTs, cancel them, create offers, and seamlessly accept NFT offers
- Issue: Issue custom assets
- Automated Market Maker: operations like bidding, creation, deletion, deposits, voting
- RegularKey: transactions to set or update account regular keys
- Offer: creation, cancel
- multi-signature transaction

### Addresses
- classAddress: They are straightforward representations of the recipient's address on the XRP Ledger
- xAddress: newer format, which provides additional features and interoperability. xAddresses include destination tags by default and are designed to simplify cross-network transactions and improve address interoperability

### Sign
- Sign XRP transactions with ED25519 and SECP256K1 algorithms.

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

### JSON-RPC Support
communicate with XRP nodes via the JSON-RPC protocol
It has been attempted to embed all the methods into RPC; however, currently, most of the data APIs are delivered in JSON format, and they have not been modeled.

## EXAMPLES
At least one example has been created for each transaction type, which you can find in the 'examples' folder.

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
Each type of transaction has its own class for creating transactions
Descriptions for some of these classes are provided below.

- Simple payment
  
  ```
    final transaction = Payment(
      destination: destination, // destination account
      account: ownerAddress, // Sender account
      amount: amount, // The amount sent can be in XRP or any other token.
      signingPubKey: ownerPublic); // Sender's public key

  ```
- NTF, mint, createOffer, acceptOffer
   
  ```
  // mint token
  final transaction = NFTokenMint(
      flags: NFTokenMintFlag.TF_TRANSFERABLE.value,
      account: ownerAddress,
      uri: "...", // that points to the data and/or metadata associated with the NFT
      signingPubKey: ownerPublic,
      memos: [memo], // Additional arbitrary information attached to this transaction
      nftokenTaxon: 1); // Indicates the taxon associated with this token

  // create offer
  final offer = NFTokenCreateOffer(
    amount: CurrencyAmount.xrp(BigInt.from(1000000)),
    flags: NFTokenCreateOfferFlag.TF_SELL_NFTOKEN.value,
    nftokenId: tokenId, /// Identifies the TokenID of the NFToken object that the offer references. 
    account: ownerAddress,
    signingPubKey: ownerPublic,
    memos: [memo],
  );
  
  // accept offer
  final offer = NFTokenAcceptOffer(
    nfTokenSellOffer: offerId,
    account: ownerAddress,
    signingPubKey: ownerPublic,
    memos: [memo],
  );

  ```
- Completely create, sign, and send transactions
  ```
  // create escrowCreate transaction
  final escrowCreate = EscrowCreate(
    account: ownerAddress,
    destination: destination,
    cancelAfterTime: cancelAfterOnDay,
    finishAfterTime: finishAfterOneHours,
    amount: BigInt.from(25000000),
    condition:
        "A0258020E488CD4C1AC9A7673CA2D2712B47049B87C308181BF3B89D6FBB74FC36836BB5810120",
    signingPubKey: ownerPublic,
    memos: [memo],
  );

  // It receives the transaction, the RPC class, and then fulfills the transaction requirements, including the fee amount, account sequence, and the last network ledger sequence.
  await autoFill(rpc, escrowCreate);
  
  // At this point, we need to sign the transaction with the sender's account.
  // We receive the transaction blob and sign it with the sender's private key.
  final sig = owner.sign(escrowCreate.toBlob());
  // After completing the signature, we add it to the transaction.
  escrowCreate.setSignature(sig);

  /// In the final step, we need to send the transaction to the network.
  /// We receive another transaction blob that already contains a signature. At this point, we no longer need to include a signature, and we must set the 'forSigning' variable to false.
  final trBlob = escrowCreate.toBlob(forSigning: false);

  // broadcasting transaction
  final result = await rpc.submit(trBlob)
  // transaction hash: result.txJson.hash ()
  // engine result: result.engineResult result.engineResult
  // engine result message: result.engineResultMessage
  
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
fmt.Println("encrypted wallet: ", encrypted)

// Decode the encrypted wallet using the password
decodedWallet, _ := secretwallet.DecodeSecretWallet(encrypted, myPassword)

// Create a new wallet from the decoded credentials, specifying SegWit support and the network
newWallet, _ := hdwallet.FromXPrivateKey(decodedWallet.Credentials, true, network)

```
### Node provider
```
// Select the network (testnet or mainnet)
network := address.TestnetNetwork

// Create an API instance (BlockCypherApi or MempoolApi).
// Currently, only a few critical methods have been implemented to retrieve unspent transactions,
// obtain network fees, receive transactions, and send transactions to the network.
api := provider.SelectApi(provider.MempoolApi, &network)

// ========================================================================================//

// Read the transaction ID (hash)
tr, e := api.GetTransaction("d4bad8e07d30ca4389ec8a203318aa523cc3e36c9730d0a6852a3801d086c5fe")
if e != nil {
	fmt.Println("Error:", e)
	return
}
if converted, ok := tr.(*provider.BlockCypherTransaction); ok {
	fmt.Println("It's a BlockCypher transaction struct")
	fmt.Println(converted.Hash)
	fmt.Println(converted.Inputs)
	fmt.Println(converted.Outputs)
	fmt.Println(converted.Confirmations)
} else if converted, ok := tr.(*provider.MempoolTransaction); ok {
	fmt.Println("Mempool transaction struct")
	fmt.Println(converted.TxID)
	fmt.Println(converted.Vout)
	fmt.Println(converted.Vin)
	fmt.Println(converted.Status.Confirmed)
}

// ========================================================================================//

addr, _ := address.P2WPKHAddressFromAddress("tb1q92nmnvhj04sqd4x7wjaewlt5jn8n3ngmplcymy")

// Read account UTXOs
utxos, e := api.GetAccountUtxo(provider.UtxoOwnerDetails{
	PublicKey: "",
	Address:   addr,
})
if e != nil {
	fmt.Println(e)
} else {
	fmt.Println("UTXOs: ", len(utxos))
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
	// Fees are in satoshis per kilobyte (PER KB)
	fmt.Println("Medium Fee: ", fee.Medium)
	fmt.Println("Low Fee: ", fee.Low)
	fmt.Println("High Fee: ", fee.High)
	
	// To calculate the transaction fee, you can use the EstimateFee method of the BitcoinFeeRate struct.
	// You'll need the transaction size (transaction.GetSize()) or virtual size (transaction.GetVSize()) for SegWit transactions.
	_ = fee.GetEstimate(500, fee.High)
}

// ========================================================================================//

// Send a transaction (replace "TRANSACTION DIGEST" with the actual transaction data)
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
	fmt.Println("It's a Mempool transactions struct")
	fmt.Println("Transactions: ", len(converted))
	for i := 0; i < len(converted); i++ {
		fmt.Println("Transaction ID: ", converted[i].TxID)
		fmt.Println("Status: ", converted[i].Status)
	}
} else if converted, ok := transaction.(provider.BlockCypherTransactionList); ok {
	fmt.Println("It's a BlockCypher transaction struct")
	for i := 0; i < len(converted); i++ {
		fmt.Println("Transaction Hash: ", converted[i].Hash)
		fmt.Println("Confirmations: ", converted[i].Confirmations)
	}
}

```

## Contributing

Contributions are welcome! Please follow these guidelines:
 - Fork the repository and create a new branch.
 - Make your changes and ensure tests pass.
 - Submit a pull request with a detailed description of your changes.

## Feature requests and bugs #

Please file feature requests and bugs in the issue tracker.


