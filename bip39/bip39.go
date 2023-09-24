// Package bip39_example demonstrates the usage of the bip39 package to work with BIP-39 mnemonic phrases,
// entropy generation, and mnemonic conversion. BIP-39 is a widely used standard for creating human-readable
// representations of cryptographic keys.
//
// This example script showcases how to generate a new BIP-39 mnemonic, derive its entropy, and convert it back
// to a mnemonic. It also includes error handling and helpful comments to guide you through the process.
//
// For more information on BIP-39, refer to the official specification: https://github.com/github.com/mrtnetwork/bitcoin/bips/blob/master/bip-0039.mediawiki

package bip39

import (
	"embed"
	"fmt"
	"github.com/mrtnetwork/bitcoin/digest"
	"github.com/mrtnetwork/bitcoin/formating"
	"regexp"
	"strings"
)

type Bip39Language string

// Bip39Language represents a string type that enumerates supported BIP-39 mnemonic languages.
const (
	English            Bip39Language = "english.txt"
	Spanish            Bip39Language = "spanish.txt"
	Portuguese         Bip39Language = "portuguese.txt"
	Korean             Bip39Language = "korean.txt"
	Japanese           Bip39Language = "japanese.txt"
	Italian            Bip39Language = "italian.txt"
	French             Bip39Language = "french.txt"
	Czech              Bip39Language = "czech.txt"
	ChineseTraditional Bip39Language = "chinese_traditional.txt"
	ChineseSimplified  Bip39Language = "chinese_simplified.txt"
)

type Bip39WordLength int

// Bip39WordLength represents an integer type that enumerates the supported word lengths for BIP-39 mnemonics.
const (
	Words12 Bip39WordLength = 128
	Words15 Bip39WordLength = 160
	Words18 Bip39WordLength = 192
	Words21 Bip39WordLength = 224
	Words24 Bip39WordLength = 256
)

//go:embed bip39_languages/*
var content embed.FS

// Bip39 represents a structure for working with BIP-39 mnemonics.
// It includes the language used for the mnemonics and the associated word list.
// The default language is English.
type Bip39 struct {
	Language Bip39Language // The language used for mnemonics (e.g., English, Spanish).
	wordList []string      // The list of words corresponding to the chosen language.
}

// ChangeLanguage changes the language used for BIP-39 mnemonics in the Bip39 instance.
// If the specified language is the same as the current language, it does nothing.
// Otherwise, it updates the language and resets the word list to an empty state.
func (bip39 *Bip39) ChangeLanguage(language Bip39Language) {
	// If the specified language is the same as the current language, no action is needed.
	if bip39.Language == language {
		return
	}

	// Update the language to the specified language.
	bip39.Language = language

	// Reset the word list to an empty state.
	bip39.wordList = make([]string, 0)
}

// LoadLanguages loads the word list for the specified BIP-39 mnemonic language.
// If the word list is already loaded (contains 2048 words), it does nothing.
// If the language is not specified, it defaults to English.
// It reads the word list from the corresponding file in the "bip39_languages" directory,
// removes newline characters, and stores the words in the Bip39 instance.
// Returns an error if any issues occur during file reading or if the word count is not 2048.
func (bip39 *Bip39) LoadLanguages() error {
	if len(bip39.wordList) == 2048 {
		return nil
	}
	if bip39.Language == "" {
		bip39.Language = English
	}

	bip39.wordList = make([]string, 0)
	files, err := content.ReadFile("bip39_languages/" + string(bip39.Language))
	if err != nil {
		return fmt.Errorf("error when reading file %v", err)
	}
	splitWords := strings.Split(string(files), "\n")

	// Create a slice to store the words without newline characters
	var bip39Words []string

	// Iterate through the words and remove leading/trailing whitespace
	for _, word := range splitWords {
		cleanedWord := strings.TrimSpace(word)
		if cleanedWord == "" {
			continue
		}
		bip39Words = append(bip39Words, cleanedWord)
	}
	if len(bip39Words) != 2048 {
		return fmt.Errorf("bip39 words list must be 2048 words but this file contains %v words", len(bip39Words))
	}
	bip39.wordList = append([]string{}, bip39Words...)
	return nil

}

// calculating the checksum for a mnemonic phrase
func deriveChecksumBits(entropy []byte) string {
	ent := len(entropy) * 8
	cs := ent / 32
	hash := digest.SingleHash(entropy)
	return formating.BytesToBinary(hash)[:cs]
}

// GenerateMnemonic creates a random BIP-39 mnemonic phrase of the specified word count.
// Returns the generated mnemonic and any error encountered during the process.
func (bip39 *Bip39) GenerateMnemonic(length Bip39WordLength) (string, error) {
	strength := int(length)
	if strength%32 != 0 {
		return "", fmt.Errorf("strength must be a multiple of 32")
	}
	entropy, err := digest.GenerateRandom(strength / 8)
	if err != nil {
		return "", err
	}
	lang := bip39.LoadLanguages()
	if lang != nil {
		return "", lang
	}
	return bip39.EntropyToMnemonic(entropy)
}

// ToSeed generates a binary seed from a given mnemonic phrase and optional passphrase.
// It uses the mnemonic and passphrase (if provided) to derive a seed for cryptographic operations.
// Parameters:
//
//	mnemonic: The BIP-39 mnemonic phrase as a space-separated string.
//	passphrase: An optional passphrase used for seed derivation. Can be an empty string.
//
// Returns:
//
//	The binary seed as a byte slice.
func ToSeed(mnemonic, passphrase string) []byte {
	salt := "mnemonic" + passphrase
	return digest.PbkdfDeriveDigest(mnemonic, salt)
}

// EntropyToMnemonic converts binary entropy data into a human-readable mnemonic phrase.
// It takes the raw entropy and returns a space-separated list of words based on the BIP-39 standard.
// This mnemonic phrase is commonly used for creating and recovering cryptocurrency wallet keys.
// Parameters:
//
//	entropy: The binary entropy data to be converted into a mnemonic.
//	wordList: The list of words from which the mnemonic is generated (e.g., English word list).
//
// Returns:
//
//	The mnemonic phrase as a string.
func (bip39 *Bip39) EntropyToMnemonic(entropy []byte) (string, error) {
	entLen := len(entropy)
	if entLen < 16 || entLen > 32 || entLen%4 != 0 {
		return "", fmt.Errorf("invalid entropy")
	}

	entropyBits := formating.BytesToBinary(entropy)
	checksumBits := deriveChecksumBits(entropy)
	bits := entropyBits + checksumBits
	// Define the regular expression pattern to split the binary string into chunks of up to 11 characters
	re := regexp.MustCompile(".{1,11}")

	// Find all matches of the pattern in the binary string
	matches := re.FindAllString(bits, -1)

	var words []string
	for _, binary := range matches {
		byteValue, _ := formating.BinaryToByte(binary)
		if int(byteValue) < len(bip39.wordList) {
			words = append(words, bip39.wordList[byteValue])
		}
	}

	return strings.Join(words, " "), nil
}

// validate mnemonic
func (bip39 *Bip39) ValidateMnemonic(mnemonic string) bool {
	_, err := bip39.MnemonicToEntropy(mnemonic)
	return err == nil
}

// wordsToBinary converts a BIP-39 mnemonic word to its binary representation.
// It searches for the word in the Bip39 instance's word list and returns the binary representation.
// Returns an error if the word is not found in the word list.
func (bip39 *Bip39) wordsToBinary(word string) (string, error) {
	// Initialize the index to -1 (not found)
	index := -1

	// Search for the word in the word list
	for i, w := range bip39.wordList {
		if w == word {
			index = i
			break
		}
	}

	// If the word is not found in the word list, return an error
	if index == -1 {
		return "", fmt.Errorf("invalid mnemonic")
	}

	// Convert the found index to a binary representation with 11 bits
	binaryRepresentation := fmt.Sprintf("%011b", index)

	return binaryRepresentation, nil
}

// mapping algorithm to convert each word into its corresponding binary value
func (bip39 *Bip39) MnemonicToEntropy(mnemonic string) ([]byte, error) {
	loadLanguage := bip39.LoadLanguages()
	if loadLanguage != nil {
		return nil, loadLanguage
	}
	words := strings.Fields(mnemonic)
	if len(words)%3 != 0 {
		return nil, fmt.Errorf("invalid mnemonic")
	}

	bits := ""
	for _, word := range words {
		b, err := bip39.wordsToBinary(word)
		if err != nil {
			return nil, err
		}
		bits += b
	}
	dividerIndex := (len(bits) / 33) * 32
	entropyBits := bits[:dividerIndex]
	checksumBits := bits[dividerIndex:]

	var entropyBytes []byte
	for i := 0; i < len(entropyBits); i += 8 {
		end := i + 8
		if end > len(entropyBits) {
			end = len(entropyBits)
		}
		binaryByte := entropyBits[i:end]
		byteValue, err := formating.BinaryToByte(binaryByte)
		if err != nil {
			return nil, err
		}
		entropyBytes = append(entropyBytes, byte(byteValue))
	}

	if len(entropyBytes) < 16 || len(entropyBytes) > 32 || len(entropyBytes)%4 != 0 {
		return nil, fmt.Errorf("invalid entropy")
	}

	newChecksumBits := deriveChecksumBits(entropyBytes)
	if newChecksumBits != checksumBits {
		return nil, fmt.Errorf("invalid mnemonic checksum")
	}

	return entropyBytes, nil
}
