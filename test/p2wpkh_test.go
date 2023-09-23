package test

import (
	"bitcoin/constant"
	"bitcoin/keypair"
	"bitcoin/scripts"
	"math/big"
	"strings"
	"testing"
)

func TestP2WPKH(t *testing.T) {
	sk := keypair.NewECPrivateFromWIF("cTALNpTpRbbxTCJ2A5Vq88UxT44w1PE2cYqiB3n4hRvzyCev1Wwo")
	p2pkhAddr := sk.GetPublic().ToAddress()
	p2wpkhAddr := sk.GetPublic().ToSegwitAddress()

	txin1 := scripts.NewTxInput("5a7b3aaa66d6b7b7abcdc9f1d05db4eee94a700297a319e19454e143875e1078", 0)
	txout1 := scripts.NewTxOutput(big.NewInt(990000), p2wpkhAddr.Program().ToScriptPubKey())

	txinSpend := scripts.NewTxInput("b3ca1c4cc778380d1e5376a5517445104e46e97176e40741508a3b07a6483ad3", 0)
	txinSpendAmount := big.NewInt(990000)

	txout2 := scripts.NewTxOutput(big.NewInt(980000), p2pkhAddr.Program().ToScriptPubKey())

	p2pkhRedeemScript := scripts.Script{Script: []interface{}{"OP_DUP", "OP_HASH160", p2pkhAddr.Program().Hash160, "OP_EQUALVERIFY", "OP_CHECKSIG"}}

	txinSpendP2pkh := scripts.NewTxInput("1e2a5279c868d61fb2ff0b1c2b04aa3eff02cd74952a8b4e799532635a9132cc", 0)

	txinSpendP2wpkh := scripts.NewTxInput("fff39047310fbf04bdd0e0bc75dde4267ae4d25219d8ad95e0ca1cee907a60da", 0)
	txinSpendP2wpkhAmount := big.NewInt(950000)

	txout3 := scripts.NewTxOutput(big.NewInt(1940000), p2pkhAddr.Program().ToScriptPubKey())

	txin1Signone := scripts.NewTxInput("fb4c338a00a75d73f9a6bd203ed4bd8884edeb111fac25a7946d5df6562f1942", 0)

	txin1SignoneAmount := big.NewInt(1000000)

	txout1Signone := scripts.NewTxOutput(big.NewInt(800000), p2pkhAddr.Program().ToScriptPubKey())

	txout2Signone := scripts.NewTxOutput(big.NewInt(190000), p2pkhAddr.Program().ToScriptPubKey())

	txin1Sigsingle := scripts.NewTxInput("b04909d4b5239a56d676c1d9d722f325a86878c9aa535915aa0df97df47cedeb", 0)

	txin1SigsingleAmount := big.NewInt(1930000)

	txout1Sigsingle := scripts.NewTxOutput(big.NewInt(1000000), p2pkhAddr.Program().ToScriptPubKey())

	txout2Sigsingle := scripts.NewTxOutput(big.NewInt(920000), p2pkhAddr.Program().ToScriptPubKey())

	txin1SiganyonecanpayAll := scripts.NewTxInput("f67e97a2564dceed405e214843e3c954b47dd4f8b26ea48f82382f51f7626036", 0)

	txin1SiganyonecanpayAllAmount := big.NewInt(180000)

	txin2SiganyonecanpayAll := scripts.NewTxInput("f4afddb77cd11a79bed059463085382c50d60c7f9e4075d8469cfe60040f68eb", 0)

	txin2SiganyonecanpayAllAmount := big.NewInt(180000)

	txout1SiganyonecanpayAll := scripts.NewTxOutput(big.NewInt(180000), p2pkhAddr.Program().ToScriptPubKey())

	txout2SiganyonecanpayAll := scripts.NewTxOutput(big.NewInt(170000), p2pkhAddr.Program().ToScriptPubKey())

	txin1SiganyonecanpayNone := scripts.NewTxInput("d2ae5d4a3f390f108769139c9b5757846be6693b785c4e21eab777eec7289095", 0)

	txin1SiganyonecanpayNoneAmount := big.NewInt(900000)

	txin2SiganyonecanpayNone := scripts.NewTxInput("ee5062d426677372e6de96e2eb47d572af5deaaef3ef225f3179dfa1ece3f4f5", 0)

	txin2SiganyonecanpayNoneAmount := big.NewInt(700000)

	txout1SiganyonecanpayNone := scripts.NewTxOutput(big.NewInt(800000), p2pkhAddr.Program().ToScriptPubKey())

	txout2SiganyonecanpayNone := scripts.NewTxOutput(big.NewInt(700000), p2pkhAddr.Program().ToScriptPubKey())

	txin1SiganyonecanpaySingle := scripts.NewTxInput("c7bb5672266c8a5b64fe91e953a9e23e3206e3b1a2ddc8e5999b607b82485042", 0)

	txin1SiganyonecanpaySingleAmount := big.NewInt(1000000)

	txout1SiganyonecanpaySingle := scripts.NewTxOutput(big.NewInt(500000), p2pkhAddr.Program().ToScriptPubKey())

	txout2SiganyonecanpaySingle := scripts.NewTxOutput(big.NewInt(490000), p2pkhAddr.Program().ToScriptPubKey())

	createSendToP2wpkhResult := "020000000178105e8743e15494e119a39702704ae9eeb45dd0f1c9cdabb7b7d666aa3a7b5a000000006a4730440220415155963673e5582aadfdb8d53874c9764cfd56c28be8d5f2838fdab6365f9902207bf28f875e15ff53e81f3245feb07c6120df4a653feabba3b7bf274790ea1fd1012102d82c9860e36f15d7b72aa59e29347f951277c21cd4d34822acdeeadbcff8a546ffffffff01301b0f0000000000160014fd337ad3bf81e086d96a68e1f8d6a0a510f8c24a00000000"
	spendP2pkhResult := "02000000000101d33a48a6073b8a504107e47671e9464e10457451a576531e0d3878c74c1ccab30000000000ffffffff0120f40e00000000001976a914fd337ad3bf81e086d96a68e1f8d6a0a510f8c24a88ac0247304402201c7ec9b049daa99c78675810b5e36b0b61add3f84180eaeaa613f8525904bdc302204854830d463a4699b6d69e37c08b8d3c6158185d46499170cfcc24d4a9e9a37f012102d82c9860e36f15d7b72aa59e29347f951277c21cd4d34822acdeeadbcff8a54600000000"
	p2pkhAndP2wpkhToP2pkhResult := "02000000000102cc32915a633295794e8b2a9574cd02ff3eaa042b1c0bffb21fd668c879522a1e000000006a47304402200fe842622e656a6780093f60b0597a36a57481611543a2e9576f9e8f1b34edb8022008ba063961c600834760037be20f45bbe077541c533b3fd257eae8e08d0de3b3012102d82c9860e36f15d7b72aa59e29347f951277c21cd4d34822acdeeadbcff8a546ffffffffda607a90ee1ccae095add81952d2e47a26e4dd75bce0d0bd04bf0f314790f3ff0000000000ffffffff01209a1d00000000001976a914fd337ad3bf81e086d96a68e1f8d6a0a510f8c24a88ac00024730440220274bb5445294033a36c360c48cc5e441ba8cc2bc1554dcb7d367088ec40a0d0302202a36f6e03f969e1b0c582f006257eec8fa2ada8cd34fe41ae2aa90d6728999d1012102d82c9860e36f15d7b72aa59e29347f951277c21cd4d34822acdeeadbcff8a54600000000"
	testSignoneSendResult := "0200000000010142192f56f65d6d94a725ac1f11ebed8488bdd43e20bda6f9735da7008a334cfb0000000000ffffffff0200350c00000000001976a914fd337ad3bf81e086d96a68e1f8d6a0a510f8c24a88ac30e60200000000001976a914fd337ad3bf81e086d96a68e1f8d6a0a510f8c24a88ac0247304402202c47de56a42143ea94c15bdeee237104524a009e50d5359596f7c6f2208a280b022076d6be5dcab09f7645d1ee001c1af14f44420c0d0b16724d741d2a5c19816902022102d82c9860e36f15d7b72aa59e29347f951277c21cd4d34822acdeeadbcff8a54600000000"
	testSigsingleSendResult := "02000000000101ebed7cf47df90daa155953aac97868a825f322d7d9c176d6569a23b5d40949b00000000000ffffffff0240420f00000000001976a914fd337ad3bf81e086d96a68e1f8d6a0a510f8c24a88acc0090e00000000001976a914fd337ad3bf81e086d96a68e1f8d6a0a510f8c24a88ac0247304402205189808e5cd0d49a8211202ea1afd7d01c180892ddf054508c349c2aa5630ee202202cbe5efa11fdde964603f4b9112d5e9ac452fba2e8ad5b6cddffbc8f0043b59e032102d82c9860e36f15d7b72aa59e29347f951277c21cd4d34822acdeeadbcff8a54600000000"
	testSiganyonecanpayAllSendResult := "02000000000102366062f7512f38828fa46eb2f8d47db454c9e34348215e40edce4d56a2977ef60000000000ffffffffeb680f0460fe9c46d875409e7f0cd6502c3885304659d0be791ad17cb7ddaff40000000000ffffffff0220bf0200000000001976a914fd337ad3bf81e086d96a68e1f8d6a0a510f8c24a88ac10980200000000001976a914fd337ad3bf81e086d96a68e1f8d6a0a510f8c24a88ac024730440220046813b802c046c9cfa309e85d1f36b17f1eb1dfb3e8d3c4ae2f74915a3b1c1f02200c5631038bb8b6c7b5283892bb1279a40e7ac13d2392df0c7b36bde7444ec54c812102d82c9860e36f15d7b72aa59e29347f951277c21cd4d34822acdeeadbcff8a5460247304402206fb60dc79b5ca6c699d04ec96c4f196938332c2909fd17c04023ebcc7408f36e02202b071771a58c84e20b7bf1fcec05c0ef55c1100436a055bfcb2bf7ed1c0683a9012102d82c9860e36f15d7b72aa59e29347f951277c21cd4d34822acdeeadbcff8a54600000000"
	testSiganyonecanpayNoneSendResult := "02000000000102959028c7ee77b7ea214e5c783b69e66b8457579b9c136987100f393f4a5daed20000000000fffffffff5f4e3eca1df79315f22eff3aeea5daf72d547ebe296dee672736726d46250ee0000000000ffffffff0200350c00000000001976a914fd337ad3bf81e086d96a68e1f8d6a0a510f8c24a88ac60ae0a00000000001976a914fd337ad3bf81e086d96a68e1f8d6a0a510f8c24a88ac0247304402203bbcbd2003244e9ccde7f705d3017f3baa2cb2d47efb63ede7e39704eff3987702206932aa4b402de898ff2fd3b2182f344dc9051b4c326dacc07b1e59059042f3ad822102d82c9860e36f15d7b72aa59e29347f951277c21cd4d34822acdeeadbcff8a54602473044022052dd29ab8bb0814b13633691148feceded29466ff8a1812d6d51c6fa53c55b5402205f25b3ae0da860da29a6745b0b587aa3fc3e05bef3121d3693ca2e3f4c2c3195012102d82c9860e36f15d7b72aa59e29347f951277c21cd4d34822acdeeadbcff8a54600000000"
	testSiganyonecanpaySingleSendResult := "02000000000101425048827b609b99e5c8dda2b1e306323ee2a953e991fe645b8a6c267256bbc70000000000ffffffff0220a10700000000001976a914fd337ad3bf81e086d96a68e1f8d6a0a510f8c24a88ac107a0700000000001976a914fd337ad3bf81e086d96a68e1f8d6a0a510f8c24a88ac02473044022064b63a1da4181764a1e8246d353b72c420999c575807ec80329c64264fd5b19e022076ec4ba6c02eae7dc9340f8c76956d5efb7d0fbad03b1234297ebed8c38e43d8832102d82c9860e36f15d7b72aa59e29347f951277c21cd4d34822acdeeadbcff8a54600000000"

	t.Run("test_signed_send_to_p2wpkh", func(t *testing.T) {
		tx := scripts.NewBtcTransaction(
			[]scripts.TxInput{*txin1},
			[]scripts.TxOutput{*txout1},
			false,
		)
		digest := tx.GetTransactionDigest(0, p2pkhAddr.Program().ToScriptPubKey())
		sig := sk.SingInput(digest)
		tx.SetScriptSig(0, *scripts.NewScript(sig, sk.GetPublic().ToHex()))
		if !strings.EqualFold(tx.Serialize(), createSendToP2wpkhResult) {
			t.Errorf("Expected %v, but got %v", createSendToP2wpkhResult, tx.Serialize())
		}

	})
	t.Run("test_spend_p2wpkh", func(t *testing.T) {
		tx := scripts.NewBtcTransaction(
			[]scripts.TxInput{*txinSpend},
			[]scripts.TxOutput{*txout2},
			true,
		)
		digest := tx.GetTransactionSegwitDigit(0, p2pkhRedeemScript, txinSpendAmount)
		sig := sk.SingInput(digest)
		witness := scripts.TxWitnessInput{Stack: []string{sig, sk.GetPublic().ToHex()}}
		tx.Witnesses = append(tx.Witnesses, witness)

		if !strings.EqualFold(tx.Serialize(), spendP2pkhResult) {
			t.Errorf("Expected %v, but got %v", spendP2pkhResult, tx.Serialize())
		}

	})
	t.Run("test_p2pkh_and_p2wpkh_to_p2pkh", func(t *testing.T) {
		tx := scripts.NewBtcTransaction(
			[]scripts.TxInput{*txinSpendP2pkh, *txinSpendP2wpkh},
			[]scripts.TxOutput{*txout3},
			true,
		)
		digest := tx.GetTransactionDigest(0, p2pkhAddr.Program().ToScriptPubKey())
		sig := sk.SingInput(digest)
		tx.SetScriptSig(0, *scripts.NewScript(sig, sk.GetPublic().ToHex()))
		witness := scripts.TxWitnessInput{Stack: []string{}}
		tx.Witnesses = append(tx.Witnesses, witness)
		segwitDigest := tx.GetTransactionSegwitDigit(1, p2pkhRedeemScript, txinSpendP2wpkhAmount)
		sig2 := sk.SingInput(segwitDigest)
		witness2 := scripts.TxWitnessInput{Stack: []string{sig2, sk.GetPublic().ToHex()}}
		tx.Witnesses = append(tx.Witnesses, witness2)
		if !strings.EqualFold(tx.Serialize(), p2pkhAndP2wpkhToP2pkhResult) {
			t.Errorf("Expected %v, but got %v", p2pkhAndP2wpkhToP2pkhResult, tx.Serialize())
		}

	})

	t.Run("test_signone_send", func(t *testing.T) {
		tx := scripts.NewBtcTransaction(
			[]scripts.TxInput{*txin1Signone},
			[]scripts.TxOutput{*txout1Signone},
			true,
		)
		digest := tx.GetTransactionSegwitDigit(0, p2pkhRedeemScript, txin1SignoneAmount, constant.SIGHASH_NONE)
		sig := sk.SingInput(digest, constant.SIGHASH_NONE)
		witness := scripts.TxWitnessInput{Stack: []string{sig, sk.GetPublic().ToHex()}}
		tx.Witnesses = append(tx.Witnesses, witness)
		tx.Outputs = append(tx.Outputs, *txout2Signone)
		if !strings.EqualFold(tx.Serialize(), testSignoneSendResult) {
			t.Errorf("Expected %v, but got %v", testSignoneSendResult, tx.Serialize())
		}

	})
	t.Run("test_sigsingle_send", func(t *testing.T) {
		tx := scripts.NewBtcTransaction(
			[]scripts.TxInput{*txin1Sigsingle},
			[]scripts.TxOutput{*txout1Sigsingle},
			true,
		)
		digest := tx.GetTransactionSegwitDigit(0, p2pkhRedeemScript, txin1SigsingleAmount, constant.SIGHASH_SINGLE)
		sig := sk.SingInput(digest, constant.SIGHASH_SINGLE)
		witness := scripts.TxWitnessInput{Stack: []string{sig, sk.GetPublic().ToHex()}}
		tx.Witnesses = append(tx.Witnesses, witness)
		tx.Outputs = append(tx.Outputs, *txout2Sigsingle)
		if !strings.EqualFold(tx.Serialize(), testSigsingleSendResult) {
			t.Errorf("Expected %v, but got %v", testSigsingleSendResult, tx.Serialize())
		}

	})

	t.Run("test_siganyonecanpay_all_send", func(t *testing.T) {
		tx := scripts.NewBtcTransaction(
			[]scripts.TxInput{*txin1SiganyonecanpayAll},
			[]scripts.TxOutput{*txout1SiganyonecanpayAll, *txout2SiganyonecanpayAll},
			true,
		)
		digest := tx.GetTransactionSegwitDigit(0, p2pkhRedeemScript, txin1SiganyonecanpayAllAmount, constant.SIGHASH_ALL|constant.SIGHASH_ANYONECANPAY)
		sig := sk.SingInput(digest, constant.SIGHASH_ALL|constant.SIGHASH_ANYONECANPAY)
		witness := scripts.TxWitnessInput{Stack: []string{sig, sk.GetPublic().ToHex()}}
		tx.Witnesses = append(tx.Witnesses, witness)
		tx.Inputs = append(tx.Inputs, *txin2SiganyonecanpayAll)
		digit2 := tx.GetTransactionSegwitDigit(1, p2pkhRedeemScript, txin2SiganyonecanpayAllAmount)
		sig2 := sk.SingInput(digit2)
		witness2 := scripts.TxWitnessInput{Stack: []string{sig2, sk.GetPublic().ToHex()}}
		tx.Witnesses = append(tx.Witnesses, witness2)
		if !strings.EqualFold(tx.Serialize(), testSiganyonecanpayAllSendResult) {
			t.Errorf("Expected %v, but got %v", testSiganyonecanpayAllSendResult, tx.Serialize())
		}

	})
	t.Run("test_siganyonecanpay_none_send", func(t *testing.T) {

		tx := scripts.NewBtcTransaction(
			[]scripts.TxInput{*txin1SiganyonecanpayNone},
			[]scripts.TxOutput{*txout1SiganyonecanpayNone},
			true,
		)
		digest := tx.GetTransactionSegwitDigit(0, p2pkhRedeemScript, txin1SiganyonecanpayNoneAmount, constant.SIGHASH_NONE|constant.SIGHASH_ANYONECANPAY)
		sig := sk.SingInput(digest, constant.SIGHASH_NONE|constant.SIGHASH_ANYONECANPAY)
		witness := scripts.TxWitnessInput{Stack: []string{sig, sk.GetPublic().ToHex()}}
		tx.Witnesses = append(tx.Witnesses, witness)
		tx.Inputs = append(tx.Inputs, *txin2SiganyonecanpayNone)
		tx.Outputs = append(tx.Outputs, *txout2SiganyonecanpayNone)
		digit2 := tx.GetTransactionSegwitDigit(1, p2pkhRedeemScript, txin2SiganyonecanpayNoneAmount)
		sig2 := sk.SingInput(digit2)
		witness2 := scripts.TxWitnessInput{Stack: []string{sig2, sk.GetPublic().ToHex()}}
		tx.Witnesses = append(tx.Witnesses, witness2)
		if !strings.EqualFold(tx.Serialize(), testSiganyonecanpayNoneSendResult) {

			t.Errorf("Expected %v, but got %v", testSiganyonecanpayNoneSendResult, tx.Serialize())
		}

	})
	t.Run("test_siganyonecanpay_single_send", func(t *testing.T) {
		tx := scripts.NewBtcTransaction(
			[]scripts.TxInput{*txin1SiganyonecanpaySingle},
			[]scripts.TxOutput{*txout1SiganyonecanpaySingle},
			true,
		)
		digest := tx.GetTransactionSegwitDigit(0, p2pkhRedeemScript, txin1SiganyonecanpaySingleAmount, constant.SIGHASH_SINGLE|constant.SIGHASH_ANYONECANPAY)
		sig := sk.SingInput(digest, constant.SIGHASH_SINGLE|constant.SIGHASH_ANYONECANPAY)
		witness := scripts.TxWitnessInput{Stack: []string{sig, sk.GetPublic().ToHex()}}
		tx.Witnesses = append(tx.Witnesses, witness)
		tx.Outputs = append(tx.Outputs, *txout2SiganyonecanpaySingle)
		if !strings.EqualFold(tx.Serialize(), testSiganyonecanpaySingleSendResult) {
			t.Errorf("Expected %v, but got %v", testSiganyonecanpaySingleSendResult, tx.Serialize())
		}

	})
}
