package test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/mrtnetwork/bitcoin/address"
	"github.com/mrtnetwork/bitcoin/formating"
	"github.com/mrtnetwork/bitcoin/keypair"
	"github.com/mrtnetwork/bitcoin/scripts"
)

func TestPrivateKeys(t *testing.T) {
	keyWifc := "KwDiBf89QgGbjEhKnhXJuH7LrciVrZi3qYjgd9M7rFU73sVHnoWn"
	keyWif := "5HpHagT65TZzG1PH3CSu63k8DbpvD8s5ip4nEB3kEsreAnchuDf"
	keyBytes := []byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x01,
	}

	t.Run("test1", func(t *testing.T) {
		p, _ := keypair.NewECPrivateFromWIF(keyWifc)
		if !bytes.Equal(p.ToBytes(), keyBytes) {
			t.Errorf("Expected %v, but got %v", keyBytes, p.ToBytes())
		}
		if !strings.EqualFold(p.ToWIF(false, &address.MainnetNetwork), keyWif) {
			t.Errorf("Expected %v, but got %v", keyWif, p.ToWIF(false, &address.MainnetNetwork))
		}
	})
}
func TestSignAndVerify(t *testing.T) {
	message := "The test!"
	keyWifC, _ := keypair.NewECPrivateFromWIF("KwDiBf89QgGbjEhKnhXJuH7LrciVrZi3qYjgd9M7rFU73sVHnoWn")
	pub := keyWifC.GetPublic()
	deterministicSignature := "204890ee41df1aa9711d239c51fb73478802863ba925bb882090a26372ebc90f525f03de46806d25892b35dfeb814ed13fd8d7ea2d8868619830bb7d6d6fbf6db2"
	t.Run("getpublic", func(t *testing.T) {
		p := keypair.GetSignaturePublic(message, formating.HexToBytes(deterministicSignature))
		if p == nil || !strings.EqualFold(p.ToAddress().Show(address.MainnetNetwork), pub.ToAddress().Show(address.MainnetNetwork)) {
			t.Errorf("Expected %v, but got %v", pub.ToAddress().Show(address.MainnetNetwork), p.ToAddress().Show(address.MainnetNetwork))
		}

	})
	t.Run("test1", func(t *testing.T) {
		sign := keyWifC.SignMessage(message, true)
		verify := pub.Verify(message, sign)
		if !strings.EqualFold(sign, deterministicSignature) {
			t.Errorf("Expected %v, but got %v", deterministicSignature, sign)
		}
		if !verify {
			t.Errorf("Expected %v, but got %v", deterministicSignature, sign)
		}
	})
}
func TestPublicKeys(t *testing.T) {
	publicKeyHex := "0479be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8"
	unCompressedAddress := "1EHNa6Q4Jz2uvNExL497mE43ikXhwF6kZm"
	publicKeyBytes := []byte{
		121, 190, 102, 126, 249, 220, 187, 172, 85, 160, 98, 149, 206, 135, 11, 7,
		2, 155, 252, 219, 45, 206, 40, 217, 89, 242, 129, 91, 22, 248, 23, 152, 72,
		58, 218, 119, 38, 163, 196, 101, 93, 164, 251, 252, 14, 17, 8, 168, 253, 23,
		180, 72, 166, 133, 84, 25, 156, 71, 208, 143, 251, 16, 212, 184,
	}

	t.Run("test1", func(t *testing.T) {
		p, e := keypair.NewECPPublicFromHex(publicKeyHex)
		if e != nil || !bytes.Equal(p.ToUnCompressedBytes(false), publicKeyBytes) {
			t.Errorf("Expected %v, but got %v", publicKeyBytes, p.ToUnCompressedBytes(false))
		}
		if e != nil || !strings.EqualFold(p.ToAddress(false).Show(address.MainnetNetwork), unCompressedAddress) {
			t.Errorf("Expected %v, but got %v", unCompressedAddress, p.ToAddress(false).Show(address.MainnetNetwork))
		}
	})

}
func TestP2pkhAddresses(t *testing.T) {
	hash160 := "91b24bf9f5288532960ac687abb035127b1d28a5"
	hash160c := "751e76e8199196d454941c45d1b3a323f1433bd6"
	address1 := "1EHNa6Q4Jz2uvNExL497mE43ikXhwF6kZm"
	addressc := "1BgGZ9tcN4rm9KBzDn7KprQz87SZ26SAMH"

	p1, _ := address.P2PKHAddressFromHash160(hash160)
	p2, _ := address.P2PKHAddressFromHash160(hash160c)
	if !strings.EqualFold(p1.Show(address.MainnetNetwork), address1) {
		t.Errorf("Expected %v, but got %v", address1, p1.Show(address.MainnetNetwork))
	}
	if !strings.EqualFold(p2.Show(address.MainnetNetwork), addressc) {
		t.Errorf("Expected %v, but got %v", addressc, p2.Show(address.MainnetNetwork))
	}
	if !strings.EqualFold(p1.Program().Hash160, hash160) {
		t.Errorf("Expected %v, but got %v", hash160, p1.Program().Hash160)
	}
	if !strings.EqualFold(p2.Program().Hash160, hash160c) {
		t.Errorf("Expected %v, but got %v", hash160c, p2.Program().Hash160)
	}
	t.Run("testx", func(t *testing.T) {
		tx, err := address.P2PKHAddressFromAddress(address1)
		if err == nil {
			fmt.Println(tx.Program().Hash160 == hash160)
		} else {
			t.Errorf(err.Error())
		}

	})
}
func TestP2SHhAddresses(t *testing.T) {
	prive, _ := keypair.NewECPrivateFromWIF("cTALNpTpRbbxTCJ2A5Vq88UxT44w1PE2cYqiB3n4hRvzyCev1Wwo")
	pub := prive.GetPublic()
	p2shaddress := "2NDkr9uD2MSY5em3rsjkff8fLZcJzCfY3W1"
	t.Run("test_create", func(t *testing.T) {
		script := scripts.NewScript(pub.ToHex(), "OP_CHECKSIG")
		addr, _ := address.P2SHAddressFromScript(script, address.P2PKInP2SH)
		if !strings.EqualFold(addr.Show(address.TestnetNetwork), p2shaddress) {
			t.Errorf("Expected %v, but got %v", p2shaddress, addr.Show(address.TestnetNetwork))
		}
	})
	t.Run("p2sh_to_script", func(t *testing.T) {
		script := scripts.NewScript(pub.ToHex(), "OP_CHECKSIG")
		fromScript := script.ToP2shScriptPubKey().ToHex()
		addr, _ := address.P2SHAddressFromScript(script, address.P2PKInP2SH)
		fromP2shAddress := addr.Program().ToScriptPubKey()
		if !strings.EqualFold(addr.Show(address.TestnetNetwork), p2shaddress) {
			t.Errorf("Expected %v, but got %v", p2shaddress, addr.Show(address.TestnetNetwork))
		}
		if !strings.EqualFold(fromScript, fromP2shAddress.ToHex()) {
			t.Errorf("Expected %v, but got %v", fromScript, fromP2shAddress.ToHex())
		}
	})
}
func TestP2WPKHADDRESS(t *testing.T) {
	priv, _ := keypair.NewECPrivateFromWIF("cVdte9ei2xsVjmZSPtyucG43YZgNkmKTqhwiUA8M4Fc3LdPJxPmZ")
	pub := priv.GetPublic()
	correctP2wpkhAddress :=
		"tb1qxmt9xgewg6mxc4mvnzvrzu4f2v0gy782fydg0w"
	correctP2shP2wpkhAddress :=
		"2N8Z5t3GyPW1hSAEJZqQ1GUkZ9ofoGhgKPf"
	correctP2wshAddress :=
		"tb1qy4kdfavhluvnhpwcqmqrd8x0ge2ynnsl7mv2mdmdskx4g3fc6ckq8f44jg"
	correctP2shP2wshAddress :=
		"2NC2DBZd3WfEF9cZcpBRDYxCTGCVCfPUf7Q"
	t.Run("test1", func(t *testing.T) {
		addr, _ := address.P2WPKHAddresssFromProgram(pub.ToSegwitAddress().Program().Program)
		if !strings.EqualFold(correctP2wpkhAddress, addr.Show(address.TestnetNetwork)) {
			t.Errorf("Expected %v, but got %v", correctP2wpkhAddress, addr.Show(address.TestnetNetwork))
		}
	})
	t.Run("test2", func(t *testing.T) {
		addr, _ := keypair.NewECPrivateFromWIF("cTmyBsxMQ3vyh4J3jCKYn2Au7AhTKvqeYuxxkinsg6Rz3BBPrYKK")
		p2sh, _ := address.P2SHAddressFromScript(addr.GetPublic().ToSegwitAddress().Program().ToScriptPubKey(), address.P2WPKHInP2SH)
		if !strings.EqualFold(correctP2shP2wpkhAddress, p2sh.Show(address.TestnetNetwork)) {
			t.Errorf("Expected %v, but got %v", correctP2shP2wpkhAddress, p2sh.Show(address.TestnetNetwork))
		}
	})
	t.Run("test3", func(t *testing.T) {
		newPrivate, _ := keypair.NewECPrivateFromWIF("cNn8itYxAng4xR4eMtrPsrPpDpTdVNuw7Jb6kfhFYZ8DLSZBCg37")
		script := scripts.NewScript("OP_1", newPrivate.GetPublic().ToHex(), "OP_1", "OP_CHECKMULTISIG")
		p2wsh, _ := address.P2WSHAddresssFromScript(script)
		addr := p2wsh.Show(address.TestnetNetwork)
		if !strings.EqualFold(correctP2wshAddress, addr) {
			t.Errorf("Expected %v, but got %v", correctP2wshAddress, addr)
		}
	})
	t.Run("test4", func(t *testing.T) {
		newPrivate, _ := keypair.NewECPrivateFromWIF("cNn8itYxAng4xR4eMtrPsrPpDpTdVNuw7Jb6kfhFYZ8DLSZBCg37")
		script := scripts.NewScript("OP_1", newPrivate.GetPublic().ToHex(), "OP_1", "OP_CHECKMULTISIG")
		p2wsh, _ := address.P2WSHAddresssFromScript(script)
		p2sh, _ := address.P2SHAddressFromScript(p2wsh.Program().ToScriptPubKey(), address.P2WSHInP2SH)
		if !strings.EqualFold(correctP2shP2wshAddress, p2sh.Show(address.TestnetNetwork)) {
			t.Errorf("Expected %v, but got %v", correctP2shP2wshAddress, p2sh.Show(address.TestnetNetwork))
		}
	})

}
func TestP2trAddresses(t *testing.T) {
	privEven, _ := keypair.NewECPrivateFromWIF("cTLeemg1bCXXuRctid7PygEn7Svxj4zehjTcoayrbEYPsHQo248w")
	privOdd, _ := keypair.NewECPrivateFromWIF("cRPxBiKrJsH94FLugmiL4xnezMyoFqGcf4kdgNXGuypNERhMK6AT")

	correctEvenPk :=
		"0271fe85f75e97d22e74c2dd6425e843def8b662b928f24f724ae6a2fd0c4e0419"
	correctEvenTrAddr :=
		"tb1pk426x6qvmncj5vzhtp5f2pzhdu4qxsshszswga8ea6sycj9nulmsu7syz0"
	correctEvenTweakedPk :=
		"b555a3680cdcf12a305758689504576f2a03421780a0e474f9eea04c48b3e7f7"

	correctOddPk :=
		"03a957ff7ead882e4c95be2afa684ab0e84447149883aba60c067adc054472785b"
	correctOddTrAddr :=
		"tb1pdr8q4tuqqeglxxhkxl3trxt0dy5jrnaqvg0ddwu7plraxvntp8dqv8kvyq"
	correctOddTweakedPk :=
		"68ce0aaf800651f31af637e2b1996f692921cfa0621ed6bb9e0fc7d3326b09da"
	t.Run("t1", func(t *testing.T) {
		pub := privEven.GetPublic().ToHex()
		if !strings.EqualFold(pub, correctEvenPk) {
			t.Errorf("Expected %v, but got %v", correctEvenPk, pub)
		}

	})
	t.Run("t2", func(t *testing.T) {
		pub := privEven.GetPublic()
		addr := pub.ToTaprootAddress().Show(address.TestnetNetwork)

		if !strings.EqualFold(addr, correctEvenTrAddr) {
			t.Errorf("Expected %v, but got %v", correctEvenTrAddr, addr)
		}

	})
	t.Run("t3", func(t *testing.T) {
		pub := privEven.GetPublic()
		program := pub.ToTaprootAddress().Program().Program
		if !strings.EqualFold(program, correctEvenTweakedPk) {
			t.Errorf("Expected %v, but got %v", correctEvenTweakedPk, program)
		}

	})
	t.Run("t5", func(t *testing.T) {
		pub := privOdd.GetPublic().ToHex()
		if !strings.EqualFold(pub, correctOddPk) {
			t.Errorf("Expected %v, but got %v", correctOddPk, pub)
		}

	})
	t.Run("t6", func(t *testing.T) {
		pub := privOdd.GetPublic()
		addr := pub.ToTaprootAddress().Show(address.TestnetNetwork)
		if !strings.EqualFold(addr, correctOddTrAddr) {
			t.Errorf("Expected %v, but got %v", correctOddTrAddr, addr)
		}

	})
	t.Run("t7", func(t *testing.T) {
		pub := privOdd.GetPublic()
		program := pub.ToTaprootAddress().Program().Program
		if !strings.EqualFold(program, correctOddTweakedPk) {
			t.Errorf("Expected %v, but got %v", correctOddTweakedPk, program)
		}

	})
}
