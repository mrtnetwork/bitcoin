package test

import (
	"bitcoin/address"
	hdwallet "bitcoin/hd_wallet"
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
	child, _ := hdwallet.DrivePath(masterWallet, "m/44'/0'/0'/0/0/1")

	// x prive from child wallet
	// xpriv := child.ToXPrivateKey(semantic, network)

	// xpub from child wallet
	xpub := child.ToXPublicKey(semantic, &network)

	// drive from public wallet
	publicWallet, _ := hdwallet.FromXPublicKey(xpub, false, &network)

	// access to publicKey
	publicKey := publicWallet.GetPublic()

	// p2pkh
	publicKey.ToAddress().Show(network)
}
