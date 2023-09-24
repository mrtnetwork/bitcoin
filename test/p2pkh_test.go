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

func TestP2PKH(t *testing.T) {
	txin := scripts.NewTxInput("fb48f4e23bf6ddf606714141ac78c3e921c8c0bebeb7c8abb2c799e9ff96ce6c", 0)
	addr, _ := address.P2PKHAddressFromAddress("n4bkvTyU1dVdzsrhWBqBw8fEMbHjJvtmJR")

	txout := scripts.NewTxOutput(
		big.NewInt(10000000),
		scripts.NewScript("OP_DUP", "OP_HASH160", addr.Program().Hash160, "OP_EQUALVERIFY", "OP_CHECKSIG"),
	)

	changeAddr, _ := address.P2PKHAddressFromAddress("mytmhndz4UbEMeoSZorXXrLpPfeoFUDzEp")

	changeTxout := scripts.NewTxOutput(
		big.NewInt(29000000),
		changeAddr.Program().ToScriptPubKey(),
	)

	changeLowSAddr, _ := address.P2PKHAddressFromAddress("mmYNBho9BWQB2dSniP1NJvnPoj5EVWw89w")

	changeLowSTxout := scripts.NewTxOutput(
		big.NewInt(29000000),
		changeLowSAddr.Program().ToScriptPubKey(),
	)

	sk, _ := keypair.NewECPrivateFromWIF("cRvyLwCPLU88jsyj94L7iJjQX5C2f8koG4G2gevN4BeSGcEvfKe9")

	fromAddr, _ := address.P2PKHAddressFromAddress("myPAE9HwPeKHh8FjKwBNBaHnemApo3dw6e")

	coreTxResult := "02000000016cce96ffe999c7b2abc8b7bebec0c821e9c378ac41417106f6ddf63be2f448fb0000000000ffffffff0280969800000000001976a914fd337ad3bf81e086d96a68e1f8d6a0a510f8c24a88ac4081ba01000000001976a914c992931350c9ba48538003706953831402ea34ea88ac00000000"
	coreTxSignedResult := "02000000016cce96ffe999c7b2abc8b7bebec0c821e9c378ac41417106f6ddf63be2f448fb000000006a473044022079dad1afef077fa36dcd3488708dd05ef37888ef550b45eb00cdb04ba3fc980e02207a19f6261e69b604a92e2bffdf6ddbed0c64f55d5003e9dfb58b874b07aef3d7012103a2fef1829e0742b89c218c51898d9e7cb9d51201ba2bf9d9e9214ebb6af32708ffffffff0280969800000000001976a914fd337ad3bf81e086d96a68e1f8d6a0a510f8c24a88ac4081ba01000000001976a914c992931350c9ba48538003706953831402ea34ea88ac00000000"
	coreTxSignedLowSSigallResult := "02000000016cce96ffe999c7b2abc8b7bebec0c821e9c378ac41417106f6ddf63be2f448fb000000006a473044022044ef433a24c6010a90af14f7739e7c60ce2c5bc3eab96eaee9fbccfdbb3e272202205372a617cb235d0a0ec2889dbfcadf15e10890500d184c8dda90794ecdf79492012103a2fef1829e0742b89c218c51898d9e7cb9d51201ba2bf9d9e9214ebb6af32708ffffffff0280969800000000001976a914fd337ad3bf81e086d96a68e1f8d6a0a510f8c24a88ac4081ba01000000001976a91442151d0c21442c2b038af0ad5ee64b9d6f4f4e4988ac00000000"
	coreTxSignedLowSSignoneResult := "02000000016cce96ffe999c7b2abc8b7bebec0c821e9c378ac41417106f6ddf63be2f448fb000000006a47304402201e4b7a2ed516485fdde697ba63f6670d43aa6f18d82f18bae12d5fd228363ac10220670602bec9df95d7ec4a619a2f44e0b8dcf522fdbe39530dd78d738c0ed0c430022103a2fef1829e0742b89c218c51898d9e7cb9d51201ba2bf9d9e9214ebb6af32708ffffffff0280969800000000001976a914fd337ad3bf81e086d96a68e1f8d6a0a510f8c24a88ac4081ba01000000001976a91442151d0c21442c2b038af0ad5ee64b9d6f4f4e4988ac00000000"
	coreTxSignedLowSSignoneTxid := "105933681b0ca37ae0c0af43ae6f111803c899232b7fd586584b532dbe21ae6f"

	sigTxin1 := scripts.NewTxInput("76464c2b9e2af4d63ef38a77964b3b77e629dddefc5cb9eb1a3645b1608b790f", 0)
	sigTxin2 := scripts.NewTxInput("76464c2b9e2af4d63ef38a77964b3b77e629dddefc5cb9eb1a3645b1608b790f", 1)

	sigFromAddr1, _ := address.P2PKHAddressFromAddress("n4bkvTyU1dVdzsrhWBqBw8fEMbHjJvtmJR")

	sigFromAddr2, _ := address.P2PKHAddressFromAddress("mmYNBho9BWQB2dSniP1NJvnPoj5EVWw89w")

	sigSk1, _ := keypair.NewECPrivateFromWIF("cTALNpTpRbbxTCJ2A5Vq88UxT44w1PE2cYqiB3n4hRvzyCev1Wwo")
	sigSk2, _ := keypair.NewECPrivateFromWIF("cVf3kGh6552jU2rLaKwXTKq5APHPoZqCP4GQzQirWGHFoHQ9rEVt")

	sigToAddr1, _ := address.P2PKHAddressFromAddress("myPAE9HwPeKHh8FjKwBNBaHnemApo3dw6e")

	sigTxout1 := scripts.NewTxOutput(
		big.NewInt(9000000),
		scripts.NewScript("OP_DUP", "OP_HASH160", sigToAddr1.Program().Hash160,
			"OP_EQUALVERIFY", "OP_CHECKSIG"),
	)

	sigToAddr2, _ := address.P2PKHAddressFromAddress("mmYNBho9BWQB2dSniP1NJvnPoj5EVWw89w")

	sigTxout2 := scripts.NewTxOutput(
		big.NewInt(900000),
		scripts.NewScript("OP_DUP", "OP_HASH160", sigToAddr2.Program().Hash160,
			"OP_EQUALVERIFY", "OP_CHECKSIG"),
	)

	sigSighashSingleResult := "02000000010f798b60b145361aebb95cfcdedd29e6773b4b96778af33ed6f42a9e2b4c4676000000006a47304402202cfd7077fe8adfc5a65fb3953fa3482cad1413c28b53f12941c1082898d4935102201d393772c47f0699592268febb5b4f64dabe260f440d5d0f96dae5bc2b53e11e032102d82c9860e36f15d7b72aa59e29347f951277c21cd4d34822acdeeadbcff8a546ffffffff0240548900000000001976a914c3f8e5b0f8455a2b02c29c4488a550278209b66988aca0bb0d00000000001976a91442151d0c21442c2b038af0ad5ee64b9d6f4f4e4988ac00000000"
	signSighashAll2in2outResult := "02000000020f798b60b145361aebb95cfcdedd29e6773b4b96778af33ed6f42a9e2b4c4676000000006a4730440220355c3cf50b1d320d4ddfbe1b407ddbe508f8e31a38cc5531dec3534e8cb2e565022037d4e8d7ba9dd1c788c0d8b5b99270d4c1d4087cdee7f139a71fea23dceeca33012102d82c9860e36f15d7b72aa59e29347f951277c21cd4d34822acdeeadbcff8a546ffffffff0f798b60b145361aebb95cfcdedd29e6773b4b96778af33ed6f42a9e2b4c4676010000006a47304402206b728374b8879fd7a10cbd4f347934d583f4301aa5d592211487732c235b85b6022030acdc07761f227c27010bd022df4b22eb9875c65a59e8e8a5722229bc7362f4012102364d6f04487a71b5966eae3e14a4dc6f00dbe8e55e61bedd0b880766bfe72b5dffffffff0240548900000000001976a914c3f8e5b0f8455a2b02c29c4488a550278209b66988aca0bb0d00000000001976a91442151d0c21442c2b038af0ad5ee64b9d6f4f4e4988ac00000000"
	signSighashNone2in2outResult := "02000000020f798b60b145361aebb95cfcdedd29e6773b4b96778af33ed6f42a9e2b4c4676000000006a47304402202a2804048b7f84f2dd7641ec05bbaf3da9ae0d2a9f9ad476d376adfd8bf5033302205170fee2ab7b955d72ae2beac3bae15679d75584c37d78d82b07df5402605bab022102d82c9860e36f15d7b72aa59e29347f951277c21cd4d34822acdeeadbcff8a546ffffffff0f798b60b145361aebb95cfcdedd29e6773b4b96778af33ed6f42a9e2b4c4676010000006a473044022021a82914b002bd02090fbdb37e2e739e9ba97367e74db5e1de834bbab9431a2f02203a11f49a3f6ac03b1550ee04f9d84deee2045bc038cb8c3e70869470126a064d022102364d6f04487a71b5966eae3e14a4dc6f00dbe8e55e61bedd0b880766bfe72b5dffffffff0240548900000000001976a914c3f8e5b0f8455a2b02c29c4488a550278209b66988aca0bb0d00000000001976a91442151d0c21442c2b038af0ad5ee64b9d6f4f4e4988ac00000000"

	signSighashAllSingleAnyone2in2outResult := "02000000020f798b60b145361aebb95cfcdedd29e6773b4b96778af33ed6f42a9e2b4c4676000000006a47304402205360315c439214dd1da10ea00a7531c0a211a865387531c358e586000bfb41b3022064a729e666b4d8ac7a09cb7205c8914c2eb634080597277baf946903d5438f49812102d82c9860e36f15d7b72aa59e29347f951277c21cd4d34822acdeeadbcff8a546ffffffff0f798b60b145361aebb95cfcdedd29e6773b4b96778af33ed6f42a9e2b4c4676010000006a473044022067943abe9fa7584ba9816fc9bf002b043f7f97e11de59155d66e0411a679ba2c02200a13462236fa520b80b4ed85c7ded363b4c9264eb7b2d9746200be48f2b6f4cb832102364d6f04487a71b5966eae3e14a4dc6f00dbe8e55e61bedd0b880766bfe72b5dffffffff0240548900000000001976a914c3f8e5b0f8455a2b02c29c4488a550278209b66988aca0bb0d00000000001976a91442151d0c21442c2b038af0ad5ee64b9d6f4f4e4988ac00000000"

	t.Run("test1", func(t *testing.T) {
		tx := scripts.NewBtcTransaction([]*scripts.TxInput{txin},
			[]*scripts.TxOutput{txout, changeTxout},
			false)
		if !strings.EqualFold(tx.Serialize(), coreTxResult) {
			t.Errorf("Expected %v, but got %v", coreTxResult, tx.Serialize())
		}

	})
	t.Run("test2", func(t *testing.T) {
		tx := scripts.NewBtcTransaction([]*scripts.TxInput{txin},
			[]*scripts.TxOutput{txout, changeTxout},
			false)

		digest := tx.GetTransactionDigest(0, scripts.NewScript("OP_DUP", "OP_HASH160", fromAddr.Program().Hash160, "OP_EQUALVERIFY", "OP_CHECKSIG"))
		sig := sk.SingInput(digest)
		tx.SetScriptSig(0, scripts.NewScript(sig, sk.GetPublic().ToHex()))
		if !strings.EqualFold(tx.Serialize(), coreTxSignedResult) {
			t.Errorf("Expected %v, but got %v", coreTxSignedResult, tx.Serialize())
		}

	})
	t.Run("test3", func(t *testing.T) {
		tx := scripts.NewBtcTransaction([]*scripts.TxInput{txin},
			[]*scripts.TxOutput{txout, changeLowSTxout},
			false)
		digest := tx.GetTransactionDigest(0, fromAddr.Program().ToScriptPubKey())
		sig := sk.SingInput(digest)
		tx.SetScriptSig(0, scripts.NewScript(sig, sk.GetPublic().ToHex()))
		if !strings.EqualFold(tx.Serialize(), coreTxSignedLowSSigallResult) {
			t.Errorf("Expected %v, but got %v", coreTxSignedLowSSigallResult, tx.Serialize())
		}

	})
	t.Run("test4", func(t *testing.T) {
		tx := scripts.NewBtcTransaction([]*scripts.TxInput{txin},
			[]*scripts.TxOutput{txout, changeLowSTxout},
			false)
		digest := tx.GetTransactionDigest(0, fromAddr.Program().ToScriptPubKey(), constant.SIGHASH_NONE)
		sig := sk.SingInput(digest, constant.SIGHASH_NONE)
		tx.SetScriptSig(0, scripts.NewScript(sig, sk.GetPublic().ToHex()))
		if !strings.EqualFold(tx.Serialize(), coreTxSignedLowSSignoneResult) {
			t.Errorf("Expected %v, but got %v", coreTxSignedLowSSignoneResult, tx.Serialize())
		}
		if !strings.EqualFold(tx.TxId(), coreTxSignedLowSSignoneTxid) {
			t.Errorf("Expected %v, but got %v", coreTxSignedLowSSignoneTxid, tx.Serialize())
		}

	})
	t.Run("test5", func(t *testing.T) {
		tx := scripts.NewBtcTransaction([]*scripts.TxInput{sigTxin1},
			[]*scripts.TxOutput{sigTxout1, sigTxout2},
			false)
		digest := tx.GetTransactionDigest(0, sigFromAddr1.Program().ToScriptPubKey(), constant.SIGHASH_SINGLE)
		sig := sigSk1.SingInput(digest, constant.SIGHASH_SINGLE)
		tx.SetScriptSig(0, scripts.NewScript(sig, sigSk1.GetPublic().ToHex()))
		if !strings.EqualFold(tx.Serialize(), sigSighashSingleResult) {
			t.Errorf("Expected %v, but got %v", sigSighashSingleResult, tx.Serialize())
		}

	})
	t.Run("test6", func(t *testing.T) {
		tx := scripts.NewBtcTransaction([]*scripts.TxInput{sigTxin1, sigTxin2},
			[]*scripts.TxOutput{sigTxout1, sigTxout2},
			false)
		digest := tx.GetTransactionDigest(0, sigFromAddr1.Program().ToScriptPubKey())
		sig := sigSk1.SingInput(digest)
		digest2 := tx.GetTransactionDigest(1, sigFromAddr2.Program().ToScriptPubKey())
		sig2 := sigSk2.SingInput(digest2)
		tx.SetScriptSig(0, scripts.NewScript(sig, sigSk1.GetPublic().ToHex()))
		tx.SetScriptSig(1, scripts.NewScript(sig2, sigSk2.GetPublic().ToHex()))
		if !strings.EqualFold(tx.Serialize(), signSighashAll2in2outResult) {
			t.Errorf("Expected %v, but got %v", signSighashAll2in2outResult, tx.Serialize())
		}

	})
	t.Run("test7", func(t *testing.T) {
		tx := scripts.NewBtcTransaction([]*scripts.TxInput{sigTxin1, sigTxin2},
			[]*scripts.TxOutput{sigTxout1, sigTxout2},
			false)
		digest := tx.GetTransactionDigest(0, sigFromAddr1.Program().ToScriptPubKey(), constant.SIGHASH_NONE)
		sig := sigSk1.SingInput(digest, constant.SIGHASH_NONE)
		digest2 := tx.GetTransactionDigest(1, sigFromAddr2.Program().ToScriptPubKey(), constant.SIGHASH_NONE)
		sig2 := sigSk2.SingInput(digest2, constant.SIGHASH_NONE)
		tx.SetScriptSig(0, scripts.NewScript(sig, sigSk1.GetPublic().ToHex()))
		tx.SetScriptSig(1, scripts.NewScript(sig2, sigSk2.GetPublic().ToHex()))
		if !strings.EqualFold(tx.Serialize(), signSighashNone2in2outResult) {
			t.Errorf("Expected %v, but got %v", signSighashNone2in2outResult, tx.Serialize())
		}

	})
	t.Run("test8", func(t *testing.T) {
		tx := scripts.NewBtcTransaction([]*scripts.TxInput{sigTxin1, sigTxin2},
			[]*scripts.TxOutput{sigTxout1, sigTxout2},
			false)
		digest := tx.GetTransactionDigest(0, sigFromAddr1.Program().ToScriptPubKey(), constant.SIGHASH_ALL|constant.SIGHASH_ANYONECANPAY)
		sig := sigSk1.SingInput(digest, constant.SIGHASH_ALL|constant.SIGHASH_ANYONECANPAY)
		digest2 := tx.GetTransactionDigest(1, sigFromAddr2.Program().ToScriptPubKey(), constant.SIGHASH_SINGLE|constant.SIGHASH_ANYONECANPAY)
		sig2 := sigSk2.SingInput(digest2, constant.SIGHASH_SINGLE|constant.SIGHASH_ANYONECANPAY)
		tx.SetScriptSig(0, scripts.NewScript(sig, sigSk1.GetPublic().ToHex()))
		tx.SetScriptSig(1, scripts.NewScript(sig2, sigSk2.GetPublic().ToHex()))
		if !strings.EqualFold(tx.Serialize(), signSighashAllSingleAnyone2in2outResult) {
			t.Errorf("Expected %v, but got %v", signSighashAllSingleAnyone2in2outResult, tx.Serialize())
		}

	})
}
