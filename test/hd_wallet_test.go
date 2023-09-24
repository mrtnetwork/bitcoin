package test

import (
	"fmt"
	"github.com/mrtnetwork/bitcoin/address"
	hdwallet "github.com/mrtnetwork/bitcoin/hd_wallet"
	"testing"
)

func TestHDWallet(t *testing.T) {
	mnemonic := "spy often critic spawn produce volcano depart fire theory fog turn retire"
	masterWalletXpub := "xpub661MyMwAqRbcFSz6LSmxGJiPKWXDcVf1P1LJuCNABxn8hJun3dbN5m2tpxjY1eX8Dpj7fGRzNfwCmepdbW6shNXb2aQ51esd3gqacTbTeYG"
	masterWalletXPrive := "xprv9s21ZrQH143K2xudEREwuAmemUgjD2wA1nQi6oxYddF9pWadW6H7XxiQyfjev243FZuhwqPkga7Lm4yYCaM8xHkQP979AnooaZEx1b3Wurz"
	childWalletXPUB := "xpub6GC7MKTAjk8ch5DSiiLuzrNrnpu47pWMt5idjSH8dxa8s3zV6Zu5xd4LWXVtHSPjiCXzwvFcCL1Fy59X7JKVZwqhwKEoqmA9QrHDC1dK4sG"
	publicWalletPublicKey := "027a6aea28f4a202d75a8e6db79276565994be98b2e91194b0c99ca017142d11c9"
	p2pkhPublicWallet := "15q8rYdtyFQfY3DudmX2LBvugpRX1xRq9k"
	p2shPublicWallet := "3NYTvWUCife9up1uVbeuvw36r4v9aog7kP"
	p2wpkhPublicWallet := "bc1qxnasdyen3m9arnrfangmjfpxfuwsrkvlpszxa7"
	p2wshPublicWallet := "bc1qczlcqww4t7wsqk4ee0hxqmwvhx64a8zeglfaw5k27yfev62aj99sgt5x9a"
	p2wshInP2shPublicWallet := "32oZLrQ6Au4jcKWA4G4HoyczzawcoM2e7d"
	newDriveFromPublicWalletXPub := "xpub6KS5kbcU4a87Bxg1EkRdS1q2thTFyzkz7Ch1rRELyNwXco2yL6jTf7RXrjsc1vL44GSVwnLj15LhAGoSCBAdCGKi7gFfHYoux9MWkYuSy4d"
	newDriveP2PKH := "1MYEc4wQ1HHtnWt7QQef5zQo8SEv71YZuG"
	newDriveP2WPKH := "bc1qu99vghdkvt3tfd6fhwjdsv9t4twjmevqsqhv24"
	newDriveP2WSH := "bc1qt2sr2l8szs9yracsr8t3lk89gd5w3qynnyg8lx2s8w2zyfkvkausxwnlnz"
	newDriveP2PKHInP2SH := "382DBHRPtJEMZuhBSgjU2AXN9qrRkf8GsA"

	newDriveFromChildPrivateKeyXpub := "xpub6Mbgf8o6ZbmLpcAM52Ezy5aq7NAyRTrZnqkS6GR9QnjL4nSR3t9Q956JHwTuifMQN42pj9P51GdsepsTgy5uW8AtBRPpR41afSGJHnrPUac"
	newDriveFromChildPrivateXPrv := "xprvA8cLFdGCjED3c85sxzhzbwe6ZLLV218iRcpqHt1XrTCMBz7GWLq9bGmpSgDnbSbze78NZV1YhpjzGZMQsyC3DVUJuwzXCqxV3zPkmaijLun"
	newDriveFromChildP2PKH := "1FazWtJ8wCebLs9LP1qEoLooQorCbUcZ1y"
	newDriveFromChildP2WPKH := "bc1qnllhl0llensmrwh4k58jl29jzwx73fzcq0cuzs"
	newDriveFromChildP2WSH := "bc1qjqa6gmaa8j5szrqpssvhlm3fnfva669gaqyf5p9cw4vt657e05sq8pxtad"
	newDriveFromChikdP2PKHInP2sh := "332T6VPmM6Yh7VAfyy6evqmEZcAUzGWniG"
	t.Run("wallet", func(t *testing.T) {

		semantic := address.P2PKH
		network := address.MainnetNetwork

		masterWallet, _ := hdwallet.FromMnemonic(mnemonic, "")
		if masterWallet.ToXPublicKey(address.P2PKH, &network) != masterWalletXpub {
			t.Errorf("Expected %v, but got %v", masterWalletXpub, masterWallet.ToXPublicKey(address.P2PKH, &network))
		}
		if masterWallet.ToXPrivateKey(address.P2PKH, &network) != masterWalletXPrive {
			t.Errorf("Expected %v, but got %v", masterWalletXPrive, masterWallet.ToXPrivateKey(address.P2PKH, &network))
		}
		child, _ := hdwallet.DrivePath(masterWallet, "m/44'/0'/0'/0/1")

		xpub := child.ToXPublicKey(semantic, &network)
		if childWalletXPUB != xpub {
			t.Errorf("Expected %v, but got %v", childWalletXPUB, xpub)
		}

		publicWallet, _ := hdwallet.FromXPublicKey(xpub, false, &network)

		publicKey := publicWallet.GetPublic()
		if publicKey.ToHex() != publicWalletPublicKey {
			t.Errorf("Expected %v, but got %v", publicWalletPublicKey, publicKey.ToHex())
		}

		p2pkh := publicKey.ToAddress().Show(network)
		if p2pkh != p2pkhPublicWallet {
			t.Errorf("Expected %v, but got %v", p2pkhPublicWallet, p2pkh)
		}

		p2sh := publicKey.ToP2PKHInP2SH().Show(network)
		if p2sh != p2shPublicWallet {
			t.Errorf("Expected %v, but got %v", p2shPublicWallet, p2sh)
		}

		p2wpkh := publicKey.ToSegwitAddress().Show(network)
		if p2wpkh != p2wpkhPublicWallet {
			t.Errorf("Expected %v, but got %v", p2wpkhPublicWallet, p2wpkh)
		}

		p2wsh := publicKey.ToP2WSHAddress().Show(network)
		fmt.Println("p2wsh: ", p2wsh)
		if p2wsh != p2wshPublicWallet {
			t.Errorf("Expected %v, but got %v", p2wshPublicWallet, p2wsh)
		}

		p2wshInP2sh := publicKey.ToP2WSHInP2SH().Show(network)
		if p2wshInP2sh != p2wshInP2shPublicWallet {
			t.Errorf("Expected %v, but got %v", p2wshInP2shPublicWallet, p2wshInP2sh)
		}

		chidFromPublicWallet, _ := hdwallet.DrivePath(publicWallet, "m/0/1")
		wif := chidFromPublicWallet.ToXPublicKey(address.P2PKH, &network)
		if wif != newDriveFromPublicWalletXPub {
			t.Errorf("Expected %v, but got %v", newDriveFromPublicWalletXPub, wif)
		}
		if newDriveP2PKH != chidFromPublicWallet.GetPublic().ToAddress().Show(network) {
			t.Errorf("Expected %v, but got %v", newDriveP2PKH, chidFromPublicWallet.GetPublic().ToAddress().Show(network))
		}
		if newDriveP2WPKH != chidFromPublicWallet.GetPublic().ToSegwitAddress().Show(network) {
			t.Errorf("Expected %v, but got %v", newDriveP2WPKH, chidFromPublicWallet.GetPublic().ToSegwitAddress().Show(network))
		}

		if newDriveP2WSH != chidFromPublicWallet.GetPublic().ToP2WSHAddress().Show(network) {
			t.Errorf("Expected %v, but got %v", newDriveP2WSH, chidFromPublicWallet.GetPublic().ToP2WSHAddress().Show(network))
		}
		if newDriveP2PKHInP2SH != chidFromPublicWallet.GetPublic().ToP2PKHInP2SH().Show(network) {
			t.Errorf("Expected %v, but got %v", newDriveP2PKHInP2SH, chidFromPublicWallet.GetPublic().ToP2PKHInP2SH().Show(network))
		}

		drivePrivate, _ := hdwallet.FromXPrivateKey(newDriveFromChildPrivateXPrv, false, &network)
		if newDriveFromChildPrivateXPrv != drivePrivate.ToXPrivateKey(address.P2PKH, &network) {
			t.Errorf("Expected %v, but got %v", newDriveFromChildPrivateXPrv, drivePrivate.ToXPrivateKey(address.P2PKH, &network))
		}
		if newDriveFromChildPrivateKeyXpub != drivePrivate.ToXPublicKey(address.P2PKH, &network) {
			t.Errorf("Expected %v, but got %v", newDriveFromChildPrivateKeyXpub, drivePrivate.ToXPublicKey(address.P2PKH, &network))
		}
		if newDriveFromChildP2PKH != drivePrivate.GetPublic().ToAddress().Show(network) {
			t.Errorf("Expected %v, but got %v", newDriveFromChildP2PKH, drivePrivate.GetPublic().ToAddress().Show(network))
		}
		if newDriveFromChildP2WPKH != drivePrivate.GetPublic().ToSegwitAddress().Show(network) {
			t.Errorf("Expected %v, but got %v", newDriveFromChildP2WPKH, drivePrivate.GetPublic().ToSegwitAddress().Show(network))
		}
		if newDriveFromChildP2WSH != drivePrivate.GetPublic().ToP2WSHAddress().Show(network) {
			t.Errorf("Expected %v, but got %v", newDriveFromChildP2WSH, drivePrivate.GetPublic().ToP2WSHAddress().Show(network))
		}
		if newDriveFromChikdP2PKHInP2sh != drivePrivate.GetPublic().ToP2PKHInP2SH().Show(network) {
			t.Errorf("Expected %v, but got %v", newDriveFromChikdP2PKHInP2sh, drivePrivate.GetPublic().ToP2PKHInP2SH().Show(network))
		}
	})

}
