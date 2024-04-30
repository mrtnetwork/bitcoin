package test

import (
	"math/big"
	"strings"
	"testing"

	"github.com/mrtnetwork/bitcoin/address"
	"github.com/mrtnetwork/bitcoin/constant"
	"github.com/mrtnetwork/bitcoin/keypair"
	"github.com/mrtnetwork/bitcoin/scripts"
)

func TestCreateP2trWithSingleTapScript(t *testing.T) {
	toPriv1, _ := keypair.NewECPrivateFromWIF("cT33CWKwcV8afBs5NYzeSzeSoGETtAB8izjDjMEuGqyqPoF7fbQR")
	toPub1 := toPriv1.GetPublic()
	privkeyTrScript1, _ := keypair.NewECPrivateFromWIF("cSW2kQbqC9zkqagw8oTYKFTozKuZ214zd6CMTDs4V32cMfH3dgKa")
	pubkeyTrScript1 := privkeyTrScript1.GetPublic()
	trScriptP2pk1 := scripts.NewScript(pubkeyTrScript1.ToXOnlyHex(), "OP_CHECKSIG")
	toPriv2, _ := keypair.NewECPrivateFromWIF(
		"cNxX8M7XU8VNa5ofd8yk1eiZxaxNrQQyb7xNpwAmsrzEhcVwtCjs")
	toPub2 := toPriv2.GetPublic()
	toAddress2 := toPub2.ToTaprootAddress([]interface{}{})
	txIn2 := scripts.NewDefaultTxInput(
		"3d4c9d73c4c65772e645ff26493590ae4913d9c37125b72398222a553b73fa66",
		0)
	txOut2 := scripts.NewTxOutput(
		big.NewInt(3000),
		toAddress2.Program().ToScriptPubKey())
	fromPriv2, _ := keypair.NewECPrivateFromWIF(
		"cT33CWKwcV8afBs5NYzeSzeSoGETtAB8izjDjMEuGqyqPoF7fbQR")
	fromPub2 := fromPriv2.GetPublic()
	fromAddress2 := fromPub2.ToTaprootAddress([]interface{}{trScriptP2pk1})
	scriptPubKey2 := fromAddress2.Program().ToScriptPubKey()
	toTaprootScriptAddress1 :=
		"tb1p0fcjs5l5xqdyvde5u7ut7sr0gzaxp4yya8mv06d2ygkeu82l65xs6k4uqr"

	signedTx2 :=
		"0200000000010166fa733b552a229823b72571c3d91349ae90354926ff45e67257c6c4739d4c3d0000000000ffffffff01b80b000000000000225120d4213cd57207f22a9e905302007b99b84491534729bd5f4065bdcb42ed10fcd50140f1776ddef90a87b646a45ad4821b8dd33e01c5036cbe071a2e1e609ae0c0963685cb8749001944dbe686662dd7c95178c85c4f59c685b646ab27e34df766b7b100000000"
	signedTx3 := "0200000000010166fa733b552a229823b72571c3d91349ae90354926ff45e67257c6c4739d4c3d0000000000ffffffff01b80b000000000000225120d4213cd57207f22a9e905302007b99b84491534729bd5f4065bdcb42ed10fcd50340bf0a391574b56651923abdb256731059008a08b5a3406cd81ce10ef5e7f936c6b9f7915ec1054e2a480e4552fa177aed868dc8b28c6263476871b21584690ef8222013f523102815e9fbbe132ffb8329b0fef5a9e4836d216dce1824633287b0abc6ac21c01036a7ed8d24eac9057e114f22342ebf20c16d37f0d25cfd2c900bf401ec09c900000000"
	// 1-create address with single script spending path
	t.Run("address_with_script_path", func(t *testing.T) {
		addr := toPub1.ToTaprootAddress([]interface{}{trScriptP2pk1})
		inTestnet := addr.Show(address.TestnetNetwork)
		if !strings.EqualFold(inTestnet, toTaprootScriptAddress1) {
			t.Errorf("Expected %v, but got %v", toTaprootScriptAddress1, inTestnet)
		}
	})
	// 2-spend taproot from key path (has single tapleaf script for spending)
	t.Run("spend_key_path2", func(t *testing.T) {
		tx := scripts.NewBtcTransaction([]*scripts.TxInput{txIn2},
			[]*scripts.TxOutput{txOut2},
			true)
		digest := tx.GetTransactionTaprootDigest(0, []*scripts.Script{scriptPubKey2}, []*big.Int{big.NewInt(3500)}, 0, scripts.NewScript(), constant.TAPROOT_SIGHASH_ALL)
		signatur := fromPriv2.SignTaprootTransaction(digest, constant.TAPROOT_SIGHASH_ALL, []interface{}{trScriptP2pk1}, true)
		witness := scripts.NewTxWitnessInput(signatur)
		tx.Witnesses = append(tx.Witnesses, witness)
		digestHex := tx.Serialize()
		if !strings.EqualFold(digestHex, signedTx2) {
			t.Errorf("Expected %v, but got %v", signedTx2, digestHex)
		}
		fromRaw, _ := scripts.BtcTransactionFromRaw(tx.Serialize())
		if !strings.EqualFold(fromRaw.TxId(), tx.TxId()) {
			t.Errorf("Expected %v, but got %v", tx.TxId(), fromRaw.Serialize())
		}

	})
	// 3-spend taproot from script path (has single tapleaf script for spending)
	t.Run("spend_script_path2", func(t *testing.T) {
		tx := scripts.NewBtcTransaction([]*scripts.TxInput{txIn2},
			[]*scripts.TxOutput{txOut2},
			true)
		digest := tx.GetTransactionTaprootDigest(0, []*scripts.Script{scriptPubKey2}, []*big.Int{big.NewInt(3500)}, 1, trScriptP2pk1, constant.TAPROOT_SIGHASH_ALL)
		signatur := privkeyTrScript1.SignTaprootTransaction(digest, constant.TAPROOT_SIGHASH_ALL, []interface{}{trScriptP2pk1}, false)
		controllBlock := scripts.NewControlBlock(fromPub2.ToXOnlyHex(), nil)
		witness := scripts.NewTxWitnessInput(signatur, trScriptP2pk1.ToHex(), controllBlock.ToHex())
		tx.Witnesses = append(tx.Witnesses, witness)
		digestHex := tx.Serialize()
		if !strings.EqualFold(digestHex, signedTx3) {
			t.Errorf("Expected %v, but got %v", signedTx2, digestHex)
		}
		fromRaw, _ := scripts.BtcTransactionFromRaw(tx.Serialize())
		if !strings.EqualFold(fromRaw.TxId(), tx.TxId()) {
			t.Errorf("Expected %v, but got %v", tx.TxId(), fromRaw.Serialize())
		}

	})
}

func TestCreateP2trWithTwoTapScripts(t *testing.T) {
	privkeyTrScriptA, _ := keypair.NewECPrivateFromWIF("cSW2kQbqC9zkqagw8oTYKFTozKuZ214zd6CMTDs4V32cMfH3dgKa")
	pubkeyTrScriptA := privkeyTrScriptA.GetPublic()
	trScriptP2pkA := scripts.NewScript(pubkeyTrScriptA.ToXOnlyHex(), "OP_CHECKSIG")

	privkeyTrScriptB, _ := keypair.NewECPrivateFromWIF("cSv48xapaqy7fPs8VvoSnxNBNA2jpjcuURRqUENu3WVq6Eh4U3JU")
	pubkeyTrScriptB := privkeyTrScriptB.GetPublic()
	trScriptP2pkB := scripts.NewScript(pubkeyTrScriptB.ToXOnlyHex(), "OP_CHECKSIG")

	fromPriv, _ := keypair.NewECPrivateFromWIF("cT33CWKwcV8afBs5NYzeSzeSoGETtAB8izjDjMEuGqyqPoF7fbQR")
	fromPub := fromPriv.GetPublic()
	fromAddress := fromPub.ToTaprootAddress(*trScriptP2pkA, *trScriptP2pkB)

	txIn := scripts.NewTxInput("808ec85db7b005f1292cea744b24e9d72ba4695e065e2d968ca17744b5c5c14d",
		0)

	toPriv, _ := keypair.NewECPrivateFromWIF("cNxX8M7XU8VNa5ofd8yk1eiZxaxNrQQyb7xNpwAmsrzEhcVwtCjs")
	toPub := toPriv.GetPublic()
	toAddress := toPub.ToTaprootAddress()

	txOut := scripts.NewTxOutput(big.NewInt(3000), toAddress.Program().ToScriptPubKey())
	scriptPubkey := fromAddress.Program().ToScriptPubKey()
	allUtxosScriptpubkeys := []*scripts.Script{scriptPubkey}
	signedTx3 :=
		"020000000001014dc1c5b54477a18c962d5e065e69a42bd7e9244b74ea2c29f105b0b75dc88e800000000000ffffffff01b80b000000000000225120d4213cd57207f22a9e905302007b99b84491534729bd5f4065bdcb42ed10fcd50340ab89d20fee5557e57b7cf85840721ef28d68e91fd162b2d520e553b71d604388ea7c4b2fcc4d946d5d3be3c12ef2d129ffb92594bc1f42cdaec8280d0c83ecc2222013f523102815e9fbbe132ffb8329b0fef5a9e4836d216dce1824633287b0abc6ac41c01036a7ed8d24eac9057e114f22342ebf20c16d37f0d25cfd2c900bf401ec09c9682f0e85d59cb20fd0e4503c035d609f127c786136f276d475e8321ec9e77e6c00000000"
	t.Run("test_spend_script_path_A_from_AB", func(t *testing.T) {
		tx := scripts.NewBtcTransaction([]*scripts.TxInput{txIn}, []*scripts.TxOutput{txOut}, true)
		txDigest := tx.GetTransactionTaprootDigest(0, allUtxosScriptpubkeys, []*big.Int{big.NewInt(3500)}, 1, trScriptP2pkA, constant.TAPROOT_SIGHASH_ALL)
		signature := privkeyTrScriptA.SignTaprootTransaction(txDigest, constant.TAPROOT_SIGHASH_ALL, []interface{}{trScriptP2pkA, trScriptP2pkB}, false)
		leafB := trScriptP2pkB.ToTapleafTaggedHash()
		controlBlock := scripts.NewControlBlock(fromPub.ToXOnlyHex(), leafB)
		witness := scripts.NewTxWitnessInput(signature, trScriptP2pkA.ToHex(), controlBlock.ToHex())
		tx.Witnesses = append(tx.Witnesses, witness)
		if !strings.EqualFold(tx.Serialize(), signedTx3) {
			t.Errorf("Expected %v, but got %v", signedTx3, tx.Serialize())
		}
		fromRaw, _ := scripts.BtcTransactionFromRaw(tx.Serialize())
		if !strings.EqualFold(fromRaw.TxId(), tx.TxId()) {
			t.Errorf("Expected %v, but got %v", tx.TxId(), fromRaw.Serialize())
		}

	})
}
func TestCreateP2trWithThreeTapScripts(t *testing.T) {
	privkeyTrScriptA, _ := keypair.NewECPrivateFromWIF("cSW2kQbqC9zkqagw8oTYKFTozKuZ214zd6CMTDs4V32cMfH3dgKa")
	pubkeyTrScriptA := privkeyTrScriptA.GetPublic()
	trScriptP2pkA := scripts.NewScript(pubkeyTrScriptA.ToXOnlyHex(), "OP_CHECKSIG")

	privkeyTrScriptB, _ := keypair.NewECPrivateFromWIF("cSv48xapaqy7fPs8VvoSnxNBNA2jpjcuURRqUENu3WVq6Eh4U3JU")
	pubkeyTrScriptB := privkeyTrScriptB.GetPublic()
	trScriptP2pkB := scripts.NewScript(pubkeyTrScriptB.ToXOnlyHex(), "OP_CHECKSIG")

	privkeyTrScriptC, _ := keypair.NewECPrivateFromWIF("cRkZPNnn3jdr64o3PDxNHG68eowDfuCdcyL6nVL4n3czvunuvryC")
	pubkeyTrScriptC := privkeyTrScriptC.GetPublic()
	trScriptP2pkC := scripts.NewScript(pubkeyTrScriptC.ToXOnlyHex(), "OP_CHECKSIG")

	fromPriv, _ := keypair.NewECPrivateFromWIF("cT33CWKwcV8afBs5NYzeSzeSoGETtAB8izjDjMEuGqyqPoF7fbQR")
	fromPub := fromPriv.GetPublic()
	fromAddress := fromPub.ToTaprootAddress([]interface{}{*trScriptP2pkA, *trScriptP2pkB}, *trScriptP2pkC)

	txIn := scripts.NewTxInput("9b8a01d0f333b2440d4d305d26641e14e0e1932ebc3c4f04387c0820fada87d3",
		0)

	toPriv, _ := keypair.NewECPrivateFromWIF("cNxX8M7XU8VNa5ofd8yk1eiZxaxNrQQyb7xNpwAmsrzEhcVwtCjs")
	toPub := toPriv.GetPublic()
	toAddress := toPub.ToTaprootAddress()

	txOut := scripts.NewTxOutput(big.NewInt(3000), toAddress.Program().ToScriptPubKey())
	scriptPubkey := fromAddress.Program().ToScriptPubKey()
	allUtxosScriptpubkeys := []*scripts.Script{scriptPubkey}
	signedTx3 :=
		"02000000000101d387dafa20087c38044f3cbc2e93e1e0141e64265d304d0d44b233f3d0018a9b0000000000ffffffff01b80b000000000000225120d4213cd57207f22a9e905302007b99b84491534729bd5f4065bdcb42ed10fcd50340644e392f5fd88d812bad30e73ff9900cdcf7f260ecbc862819542fd4683fa9879546613be4e2fc762203e45715df1a42c65497a63edce5f1dfe5caea5170273f2220e808f1396f12a253cf00efdf841e01c8376b616fb785c39595285c30f2817e71ac61c01036a7ed8d24eac9057e114f22342ebf20c16d37f0d25cfd2c900bf401ec09c9ed9f1b2b0090138e31e11a31c1aea790928b7ce89112a706e5caa703ff7e0ab928109f92c2781611bb5de791137cbd40a5482a4a23fd0ffe50ee4de9d5790dd100000000"
	t.Run("test_spend_script_path_A_from_AB", func(t *testing.T) {
		tx := scripts.NewBtcTransaction([]*scripts.TxInput{txIn}, []*scripts.TxOutput{txOut}, true)
		txDigest := tx.GetTransactionTaprootDigest(0, allUtxosScriptpubkeys, []*big.Int{big.NewInt(3500)}, 1, trScriptP2pkB, constant.TAPROOT_SIGHASH_ALL)
		signature := privkeyTrScriptB.SignTaprootTransaction(txDigest, constant.TAPROOT_SIGHASH_ALL, []interface{}{trScriptP2pkA, trScriptP2pkB}, false)
		leafA := trScriptP2pkA.ToTapleafTaggedHash()
		leafC := trScriptP2pkC.ToTapleafTaggedHash()
		controlBlock := scripts.NewControlBlock(fromPub.ToXOnlyHex(), append(leafA, leafC...))
		witness := scripts.NewTxWitnessInput(signature, trScriptP2pkB.ToHex(), controlBlock.ToHex())
		tx.Witnesses = append(tx.Witnesses, witness)
		if !strings.EqualFold(tx.Serialize(), signedTx3) {
			t.Errorf("Expected %v, but got %v", signedTx3, tx.Serialize())
		}
		fromRaw, _ := scripts.BtcTransactionFromRaw(tx.Serialize())
		if !strings.EqualFold(fromRaw.TxId(), tx.TxId()) {
			t.Errorf("Expected %v, but got %v", tx.TxId(), fromRaw.Serialize())
		}

	})
}
