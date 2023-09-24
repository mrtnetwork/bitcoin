package test

import (
	"fmt"
	"github.com/mrtnetwork/bitcoin/bip39"
	"github.com/mrtnetwork/bitcoin/formating"
	"testing"
)

func TestBip39(t *testing.T) {
	bip := bip39.Bip39{
		// Select the desired language, the default language is English
		// https://github.com/github.com/mrtnetwork/bitcoin/bips/tree/master/bip-0039
		// English
		// Spanish
		// Portuguese
		// Korean
		// Japanese
		// Italian
		// French
		// Czech
		// ChineseTraditional
		// ChineseSimplified
		Language: bip39.Japanese,
	}

	// Select the desired number of words. 12(Words12), 15(Words15), 18(Words18), 21(Words21) or 24(Words24) words
	mnemonic, err := bip.GenerateMnemonic(bip39.Words24)
	fmt.Println("mnemonic: ", mnemonic)
	if err != nil {
		t.Errorf(err.Error())
	}
	// passphrase: An optional passphrase used for seed derivation. Can be an empty string.
	toSeed := bip39.ToSeed(mnemonic, "PASSPHRASE")
	fmt.Println("seed: ", formating.BytesToHex(toSeed))

	toEntropy, err := bip.MnemonicToEntropy(mnemonic)
	if err != nil {
		t.Errorf(err.Error())
	}
	toMnemonicFromEntropy, err := bip.EntropyToMnemonic(toEntropy)
	if err != nil {
		t.Errorf(err.Error())
	}
	if toMnemonicFromEntropy != mnemonic {
		t.Errorf("Expected %v, but got %v", mnemonic, toMnemonicFromEntropy)
	}

	// Use the `bip.ChangeLanguage()` method to change the language or create new instance,
	// otherwise the word list will not change
	bip.ChangeLanguage(bip39.Italian)

}
