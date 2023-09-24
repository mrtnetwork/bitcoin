package test

import (
	"github.com/mrtnetwork/bitcoin/address"
	"github.com/mrtnetwork/bitcoin/constant"
	"github.com/mrtnetwork/bitcoin/keypair"
	"github.com/mrtnetwork/bitcoin/scripts"
	"math/big"
	"strings"
	"testing"
)

func TestP2SH(t *testing.T) {
	fromAddr, _ := address.P2PKHAddressFromAddress("n4bkvTyU1dVdzsrhWBqBw8fEMbHjJvtmJR")
	sk, _ := keypair.NewECPrivateFromWIF("cTALNpTpRbbxTCJ2A5Vq88UxT44w1PE2cYqiB3n4hRvzyCev1Wwo")
	p2pkSk, _ := keypair.NewECPrivateFromWIF("cRvyLwCPLU88jsyj94L7iJjQX5C2f8koG4G2gevN4BeSGcEvfKe9")
	p2pkRedeemScript := scripts.NewScript(p2pkSk.GetPublic().ToHex(),
		"OP_CHECKSIG")
	txout := scripts.NewTxOutput(
		big.NewInt(9000000),
		p2pkRedeemScript.ToP2shScriptPubKey(),
	)
	createP2shAndSendResult := "02000000010f798b60b145361aebb95cfcdedd29e6773b4b96778af33ed6f42a9e2b4c4676000000006a47304402206f4027d0a1720ea4cc68e1aa3cc2e0ca5996806971c0cd7d40d3aa4309d4761802206c5d9c0c26dec8edab91c1c3d64e46e4dd80d8da1787a9965ade2299b41c3803012102d82c9860e36f15d7b72aa59e29347f951277c21cd4d34822acdeeadbcff8a546ffffffff01405489000000000017a9142910fc0b1b7ab6c9789c5a67c22c5bcde5b903908700000000"

	txinSpend := scripts.NewTxInput("7db363d5a7fabb64ccce154e906588f1936f34481223ea8c1f2c935b0a0c945b", 0)

	toAddr := fromAddr

	txout2 := scripts.NewTxOutput(big.NewInt(8000000), toAddr.Program().ToScriptPubKey())

	spendP2shResult := "02000000015b940c0a5b932c1f8cea231248346f93f18865904e15cecc64bbfaa7d563b37d000000006c47304402204984c2089bf55d5e24851520ea43c431b0d79f90d464359899f27fb40a11fbd302201cc2099bfdc18c3a412afb2ef1625abad8a2c6b6ae0bf35887b787269a6f2d4d01232103a2fef1829e0742b89c218c51898d9e7cb9d51201ba2bf9d9e9214ebb6af32708acffffffff0100127a00000000001976a914fd337ad3bf81e086d96a68e1f8d6a0a510f8c24a88ac00000000"

	//
	skCsvP2pkh, _ := keypair.NewECPrivateFromWIF("cRvyLwCPLU88jsyj94L7iJjQX5C2f8koG4G2gevN4BeSGcEvfKe9")

	seq, _ := scripts.NewSequence(constant.TYPE_RELATIVE_TIMELOCK, 200, true)
	inputSeq, _ := seq.ForInputSequence()
	txinSeq := scripts.NewTxInput("f557c623e55f0affc696b742630770df2342c4aac395e0ed470923247bc51b95", 0, inputSeq)

	anotherAddr, _ := address.P2PKHAddressFromAddress("n4bkvTyU1dVdzsrhWBqBw8fEMbHjJvtmJR")

	spendP2shCsvP2pkhResult := "0200000001951bc57b24230947ede095c3aac44223df70076342b796c6ff0a5fe523c657f5000000008947304402205c2e23d8ad7825cf44b998045cb19b49cf6447cbc1cb76a254cda43f7939982002202d8f88ab6afd2e8e1d03f70e5edc2a277c713018225d5b18889c5ad8fd6677b4012103a2fef1829e0742b89c218c51898d9e7cb9d51201ba2bf9d9e9214ebb6af327081e02c800b27576a914c3f8e5b0f8455a2b02c29c4488a550278209b66988acc80000000100ab9041000000001976a914fd337ad3bf81e086d96a68e1f8d6a0a510f8c24a88ac00000000"

	t.Run("test1", func(t *testing.T) {
		txin := scripts.NewTxInput("76464c2b9e2af4d63ef38a77964b3b77e629dddefc5cb9eb1a3645b1608b790f", 0)
		tx := scripts.NewBtcTransaction([]*scripts.TxInput{txin},
			[]*scripts.TxOutput{txout},
			false)
		digest := tx.GetTransactionDigest(0, fromAddr.Program().ToScriptPubKey())
		sig := sk.SingInput(digest)
		tx.SetScriptSig(0, scripts.NewScript(sig, sk.GetPublic().ToHex()))
		if !strings.EqualFold(tx.Serialize(), createP2shAndSendResult) {
			t.Errorf("Expected %v, but got %v", createP2shAndSendResult, tx.Serialize())
		}

	})
	t.Run("test2", func(t *testing.T) {
		tx := scripts.NewBtcTransaction([]*scripts.TxInput{txinSpend},
			[]*scripts.TxOutput{txout2},
			false)
		digest := tx.GetTransactionDigest(0, p2pkRedeemScript)
		sig := p2pkSk.SingInput(digest)
		tx.SetScriptSig(0, scripts.NewScript(sig, p2pkRedeemScript.ToHex()))
		if !strings.EqualFold(tx.Serialize(), spendP2shResult) {
			t.Errorf("Expected %v, but got %v", spendP2shResult, tx.Serialize())
		}

	})
	t.Run("test3", func(t *testing.T) {
		scriptSequance, err := seq.ForScript()
		if err != nil {
			t.Errorf("invalid script sequance")
			return
		}
		redeemScript := scripts.NewScript(scriptSequance, "OP_CHECKSEQUENCEVERIFY", "OP_DROP",
			"OP_DUP", "OP_HASH160", skCsvP2pkh.GetPublic().ToHash160(),
			"OP_EQUALVERIFY", "OP_CHECKSIG")
		txout1 := scripts.NewTxOutput(big.NewInt(1100000000), anotherAddr.Program().ToScriptPubKey())
		tx := scripts.NewBtcTransaction([]*scripts.TxInput{txinSeq},
			[]*scripts.TxOutput{txout1},
			false)
		digest := tx.GetTransactionDigest(0, redeemScript)
		sig := p2pkSk.SingInput(digest)
		tx.SetScriptSig(0, scripts.NewScript(sig, skCsvP2pkh.GetPublic().ToHex(), redeemScript.ToHex()))
		if !strings.EqualFold(tx.Serialize(), spendP2shCsvP2pkhResult) {
			t.Errorf("Expected %v, but got %v", spendP2shCsvP2pkhResult, tx.Serialize())
		}

	})
}
