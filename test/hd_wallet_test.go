package test

import (
	"bitcoin/address"
	hdwallet "bitcoin/hd_wallet"
	"fmt"
	"testing"
)

func TestHDWallet(t *testing.T) {

	semantic := address.P2PKH
	network := address.MainnetNetwork
	// generate random mnemonic 12,15,18,21,22 characters
	// mnemonic, _ := bip39.GenerateMnemonic(256)
	mnemonic := "spy often critic spawn produce volcano depart fire theory fog turn retire"

	// accsess to private and public keys
	masterWallet, _ := hdwallet.FromMnemonic(mnemonic, "")

	// wallet with path // does access to master wallet
	child, _ := hdwallet.DrivePath(masterWallet, "m/44'/0'/0'/0/1")

	// x prive from child wallet
	// xpriv := child.ToXPrivateKey(semantic, network)

	// xpub from child wallet
	xpub := child.ToXPublicKey(semantic, &network)

	// drive from public wallet
	publicWallet, _ := hdwallet.FromXPublicKey(xpub, false, &network)

	// access to publicKey
	publicKey := publicWallet.GetPublic()
	// priuvate
	// privateKey := p
	// p2pkh
	// 15q8rYdtyFQfY3DudmX2LBvugpRX1xRq9k
	p2pkh := publicKey.ToAddress().Show(network)
	fmt.Println("p2pkh: ", p2pkh)
	// p2sh(p2pkh)
	// 3NYTvWUCife9up1uVbeuvw36r4v9aog7kP
	p2sh := publicKey.ToP2PKHInP2SH().Show(network)
	fmt.Println("p2sh: ", p2sh)

	//p2wpk
	//bc1qxnasdyen3m9arnrfangmjfpxfuwsrkvlpszxa7
	p2wpk := publicKey.ToSegwitAddress().Show(network)
	fmt.Println("p2wpk: ", p2wpk)
	// p2wsh
	// bc1qczlcqww4t7wsqk4ee0hxqmwvhx64a8zeglfaw5k27yfev62aj99sgt5x9a
	p2wsh := publicKey.ToP2WSHAddress().Show(network)
	fmt.Println("p2wsh: ", p2wsh)

	// p2wshInP2sh
	// bc1qczlcqww4t7wsqk4ee0hxqmwvhx64a8zeglfaw5k27yfev62aj99sgt5x9a
	p2wshInP2sh := publicKey.ToP2WSHInP2SH().Show(network)
	fmt.Println("p2wshInP2sh: ", p2wshInP2sh)

	// s := publicWallet.ToXPublicKey(address.P2PKH, &network)

	// xpub6GC7MKTAjk8ch5DSiiLuzrNrnpu47pWMt5idjSH8dxa8s3zV6Zu5xd4LWXVtHSPjiCXzwvFcCL1Fy59X7JKVZwqhwKEoqmA9QrHDC1dK4sG

	chidFromPublicWallet, _ := hdwallet.DrivePath(publicWallet, "m/0/1")
	//xpub6KS5kbcU4a87Bxg1EkRdS1q2thTFyzkz7Ch1rRELyNwXco2yL6jTf7RXrjsc1vL44GSVwnLj15LhAGoSCBAdCGKi7gFfHYoux9MWkYuSy4d
	wif := chidFromPublicWallet.ToXPublicKey(address.P2PKH, &network)
	fmt.Println("wif: ", wif)
	// 1MYEc4wQ1HHtnWt7QQef5zQo8SEv71YZuG
	fmt.Println(chidFromPublicWallet.GetPublic().ToAddress().Show(network))
	// bc1qu99vghdkvt3tfd6fhwjdsv9t4twjmevqsqhv24
	fmt.Println(chidFromPublicWallet.GetPublic().ToSegwitAddress().Show(network))
	// bc1qt2sr2l8szs9yracsr8t3lk89gd5w3qynnyg8lx2s8w2zyfkvkausxwnlnz
	fmt.Println(chidFromPublicWallet.GetPublic().ToP2WSHAddress().Show(network))
	// 382DBHRPtJEMZuhBSgjU2AXN9qrRkf8GsA
	fmt.Println(chidFromPublicWallet.GetPublic().ToP2PKHInP2SH().Show(network))

	drivePrivate, _ := hdwallet.DrivePath(child, "m/0'/1'/1")
	// xpub6Mbgf8o6ZbmLpcAM52Ezy5aq7NAyRTrZnqkS6GR9QnjL4nSR3t9Q956JHwTuifMQN42pj9P51GdsepsTgy5uW8AtBRPpR41afSGJHnrPUac
	fmt.Println("wif: ", drivePrivate.ToXPrivateKey(address.P2PKH, &network))
	// xpub6Mbgf8o6ZbmLpcAM52Ezy5aq7NAyRTrZnqkS6GR9QnjL4nSR3t9Q956JHwTuifMQN42pj9P51GdsepsTgy5uW8AtBRPpR41afSGJHnrPUac
	fmt.Println("wif: ", drivePrivate.ToXPublicKey(address.P2PKH, &network))
	// 1MYEc4wQ1HHtnWt7QQef5zQo8SEv71YZuG
	fmt.Println(drivePrivate.GetPublic().ToAddress().Show(network))
	// bc1qu99vghdkvt3tfd6fhwjdsv9t4twjmevqsqhv24
	fmt.Println(drivePrivate.GetPublic().ToSegwitAddress().Show(network))
	// bc1qt2sr2l8szs9yracsr8t3lk89gd5w3qynnyg8lx2s8w2zyfkvkausxwnlnz
	fmt.Println(drivePrivate.GetPublic().ToP2WSHAddress().Show(network))
	// 382DBHRPtJEMZuhBSgjU2AXN9qrRkf8GsA
	fmt.Println(drivePrivate.GetPublic().ToP2PKHInP2SH().Show(network))

}
