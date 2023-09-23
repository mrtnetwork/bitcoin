# Bitcoin GO Package

A comprehensive Bitcoin library for GO that provides functionality to create, sign, and send Bitcoin transactions. This library supports a wide range of Bitcoin transaction types and features, making it suitable for various use cases.


## Features

- Create and sign Bitcoin transactions
- Addresses
  - Legacy Keys and Addresses (P2PK, P2PKH, P2SH)
  - Segwit Addresses (P2WPKH, P2SH-P2WPKH, P2WSH and P2SH-P2WSH, Taproot (segwit v1))
- Support for different transaction types:
  - Legacy transactions (P2PKH, P2SH(P2PK), P2SH(P2PKH) )
      - Transaction with P2PKH input and outputs
      - Create a P2PKH Transaction with different SIGHASHes
      - Create a P2SH Address
      - Create (spent) a P2SH Transaction
- Segwit Transactions
  - Transaction to pay to a P2WPKH, P2WSH, P2SH(segwit)
  - Spend from a P2SH(P2WPKH) nested segwit address
  - Spend from a P2SH(P2WSH) nested segwit address
- Timelock Transactions
  - Create a P2SH address with a relative timelock
  - Spend from a timelocked address
- Taproot (segwit v1) Transactions
  - Spend from a taproot address
  - Spend a multi input that contains both taproot and legacy UTXOs
  - Send to taproot address that contains a single script path spend
  - Spend taproot from key path (has single alternative script path spend)
  - Spend taproot from script path (has single alternative script path spend)
  - Send to taproot address that contains two scripts path spends
  - Send to taproot address that contains three scripts path spends
- Sign
  - sign message
  - sign transactions
  - Schnorr sign (segwit transactions)
  - support different `sighash`
  - get public key of signature

## Example
A large number of examples and tests have been prepared you can see them in the [test folder](https://github.com/MohsenHaydari/bitcoin/tree/main/test)

Finalizing the transaction with 15 different input types (spend: p2sh, p2wsh, p2wpkh, p2tr, p2sh, p2shInP2wsh, p2shInP2wpkh, p2shInP2pk, p2shInP2pkh, P2wsh(4-4 multi-sig), P2sh(4-7 multi-sig), and etc...) and 15 different output types within a single transaction. [mempol](https://mempool.space/testnet/tx/ffb96b60303eb8e76654d320204a2727dec57ca00cf947a50c2be40ff084e35e)


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

