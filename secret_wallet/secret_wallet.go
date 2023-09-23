package secretwallet

import (
	"bitcoin/formating"
	"bitcoin/uuid"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"golang.org/x/crypto/scrypt"
)

type CryptoParams struct {
	// Specifies the encryption algorithm used (e.g., "aes-128-ctr").
	Cipher string `json:"cipher"`
	// Additional parameters specific to the chosen cipher.
	CipherParams map[string]string `json:"cipherparams"`
	// Key Derivation Function used (e.g., "pbkdf2" or "scrypt").
	KDF string `json:"kdf"`
	// Additional parameters specific to the chosen KDF.
	KDFParams map[string]interface{} `json:"kdfparams"`
	// The encrypted data or ciphertext.
	Ciphertext string `json:"ciphertext"`
	// Message Authentication Code used for data integrity verification.
	MAC string `json:"mac"`
}

type WalletData struct {
	// Specifies the version of the wallet data structure (e.g., version 3).
	Version int `json:"version"`
	// Contains cryptographic parameters and information for data encryption and decryption.
	Crypto CryptoParams `json:"crypto"`
	// Represents a unique identifier or UUID associated with the wallet.
	ID string `json:"id"`
}

func (w *WalletData) Validate() error {
	if w.Version != 3 {
		return fmt.Errorf("library only supports version 3 of wallet files")
	}
	return nil
}

type KeyDerivator interface {
	// DeriveKey derives a cryptographic key from a given password.
	// It takes the password as input and returns the derived key as a byte slice.
	// It may also return an error if the key derivation process fails.
	DeriveKey(password []byte) ([]byte, error)
	// Name returns the name or identifier of the key derivation function (KDF).
	// This can be used to identify the specific KDF being used, such as "pbkdf2" or "scrypt."
	Name() string
	// Encode returns a map of parameters and their values that describe the key derivation function.
	// This information is typically used for encoding and decoding purposes when storing or exchanging data.
	Encode() map[string]interface{}
}

type PBDKDF2KeyDerivator struct {
	Iterations int    // The number of iterations performed during key derivation.
	Salt       []byte // A random salt value used as an additional input to the KDF.
	DkLen      int    // The length (in bytes) of the derived key.
}

// DeriveKey derives a cryptographic key from the provided password using PBKDF2.
// It takes the password as input and performs key derivation based on the PBKDF2 algorithm
// configured with the parameters stored in the PBDKDF2KeyDerivator struct.
// The method returns the derived key as a byte slice and may return an error if
// the key derivation process encounters any issues.
func (k *PBDKDF2KeyDerivator) DeriveKey(password []byte) ([]byte, error) {
	mac := sha256.New()
	mac.Write(password)
	mac.Write(k.Salt)
	dk := mac.Sum(nil)
	for i := 1; i < k.Iterations; i++ {
		mac.Reset()
		mac.Write(dk)
		dk = mac.Sum(nil)
	}
	if len(dk) > k.DkLen {
		dk = dk[:k.DkLen]
	}
	return dk, nil
}

// Name returns the name or identifier of the key derivation function (KDF).
// In this implementation, it always returns "pbkdf2" to indicate the use of the PBKDF2 algorithm.
// This information is useful for identifying the specific KDF being used within cryptographic operations.
func (k *PBDKDF2KeyDerivator) Name() string {
	return "pbkdf2"
}

// Encode returns a map of parameters and their values that describe the configuration
// of the PBKDF2 key derivation function.
// The returned map includes information such as the number of iterations, salt value,
// and desired key length.
// This information is typically used for encoding and decoding purposes when storing
// or exchanging data, allowing the configuration to be preserved and communicated.
func (k *PBDKDF2KeyDerivator) Encode() map[string]interface{} {
	return map[string]interface{}{
		"c":     k.Iterations,
		"dklen": k.DkLen,
		"prf":   "hmac-sha256",
		"salt":  formating.BytesToHex(k.Salt),
	}
}

// ScryptKeyDerivator represents the configuration parameters for the Scrypt key derivation function.
type ScryptKeyDerivator struct {
	// DkLen is the length (in bytes) of the derived key.
	DkLen int
	// N is the CPU/memory cost factor (parallelization factor) for Scrypt.
	N int
	// R is the block size parameter for Scrypt.
	R int
	// P is the parallelization parameter for Scrypt.
	P int
	// Salt is a random salt value used as an additional input to the KDF.
	Salt []byte
}

// DeriveKey takes a user-provided password and applies the Scrypt key derivation function
// to generate a derived key. The derived key is computed based on the Scrypt parameters,
// including N, R, P, and the salt value, which should be set prior to calling this method.
// The derived key is returned as a byte slice.
// If any error occurs during the key derivation process, an error is returned.
func (k *ScryptKeyDerivator) DeriveKey(password []byte) ([]byte, error) {
	// Implement Scrypt key derivation here
	return scrypt.Key(password, k.Salt, k.N, k.R, k.P, k.DkLen)
}

// Name returns the name or identifier for the Scrypt key derivation function.
// This method is used to indicate the specific key derivation function being employed,
// and it returns the string "scrypt" to represent Scrypt-based key derivation.
func (k *ScryptKeyDerivator) Name() string {
	return "scrypt"
}

// Encode converts the ScryptKeyDerivator's parameters into a map of strings and
// interfaces for serialization purposes. The resulting map includes the following
// key-value pairs:
// - "dklen": The length (in bytes) of the derived key.
// - "n": The CPU/memory cost factor (parallelization factor) for Scrypt.
// - "r": The block size parameter for Scrypt.
// - "p": The parallelization parameter for Scrypt.
// - "salt": The salt value, encoded as a hexadecimal string.
// This method is used to represent the Scrypt parameters in a format suitable for
// serialization, such as JSON.
func (k *ScryptKeyDerivator) Encode() map[string]interface{} {
	return map[string]interface{}{
		"dklen": k.DkLen,
		"n":     k.N,
		"r":     k.R,
		"p":     k.P,
		"salt":  formating.BytesToHex(k.Salt),
	}
}

// SecretWallet represents a wallet file used to securely store credentials,
// such as a private key associated with an bitcoin address. The private key
// within the wallet is encrypted with a secret password and can be decoded
// using the provided KeyDerivator, password, and initialization vector (IV).
// The ID field holds a unique identifier for the wallet.
type SecretWallet struct {
	// Credentials is a hexadecimal string representing the encrypted private key.
	Credentials string
	// Derivator is an interface for key derivation, specifying the method used
	// to derive the decryption key from the password.
	Derivator KeyDerivator
	// Password is a byte slice representing the user's secret password.

	Password []byte
	// IV (Initialization Vector) is a byte slice used in the AES decryption process
	// to ensure secure encryption and decryption of the private key.
	IV []byte
	// ID is a unique identifier associated with the wallet.
	ID []byte
}

// NewPBDKDF2KeyDerivator creates a new instance of PBDKDF2KeyDerivator with the specified
// key derivation parameters. It is used to derive keys using the PBKDF2 key derivation
// function, which employs HMAC with SHA-256 as the pseudo-random function (PRF).
//
// Parameters:
// - iterations: The number of iterations for key stretching.
// - salt: A byte slice representing the salt value for key derivation.
// - dkLen: The desired length (in bytes) of the derived key.
//
// Returns:
// - A pointer to a new PBDKDF2KeyDerivator instance initialized with the provided parameters.
//
// Example usage:
//
//	derivator := NewPBDKDF2KeyDerivator(10000, []byte("random_salt"), 32)
func NewPBDKDF2KeyDerivator(iterations int, salt []byte, dkLen int) *PBDKDF2KeyDerivator {
	return &PBDKDF2KeyDerivator{
		Iterations: iterations,
		Salt:       salt,
		DkLen:      dkLen,
	}
}

// NewScryptKeyDerivator creates a new instance of ScryptKeyDerivator with the specified
// key derivation parameters. It is used to derive keys using the scrypt key derivation
// function, which is designed to be memory-hard and CPU-intensive for enhanced security.
//
// Parameters:
// - dkLen: The desired length (in bytes) of the derived key.
// - n: The CPU/memory cost parameter, typically a power of two.
// - r: The block size parameter.
// - p: The parallelization parameter.
// - salt: A byte slice representing the salt value for key derivation.
//
// Returns:
// - A pointer to a new ScryptKeyDerivator instance initialized with the provided parameters.
//
// Example usage:
//
//	derivator := NewScryptKeyDerivator(32, 16384, 8, 1, []byte("random_salt"))
func NewScryptKeyDerivator(dkLen, n, r, p int, salt []byte) *ScryptKeyDerivator {
	return &ScryptKeyDerivator{
		DkLen: dkLen,
		N:     n,
		R:     r,
		P:     p,
		Salt:  salt,
	}
}

// NewSecretWallet creates a new instance of SecretWallet, which represents a wallet
// file used to securely store credentials like private keys. The private key is
// encrypted using the provided password and key derivation parameters.
//
// Parameters:
//   - credentials: The credentials to be stored in the wallet.
//   - password: The password used for encrypting the private key.
//   - scryptN: The CPU/memory cost parameter for the scrypt key derivation function,
//     typically a power of two.
//   - p: The parallelization parameter for the scrypt key derivation function.
func NewSecretWallet(credentials string, password string, scryptN int, p int) (*SecretWallet, error) {
	passwordBytes := []byte(password)
	salt, err := formating.GenerateRandom(32)
	if err != nil {
		return nil, err
	}

	derivator := NewScryptKeyDerivator(32, scryptN, 8, p, salt)
	u, err := uuid.GenerateUUIDv4Bytes()
	if err != nil {
		return nil, err
	}
	iv, err := formating.GenerateRandom(128 / 8)
	if err != nil {
		return nil, err
	}
	return &SecretWallet{
		Credentials: credentials,
		Derivator:   derivator,
		Password:    passwordBytes,
		IV:          iv,
		ID:          u,
	}, nil
}
func encodedwalletTobytes(encoded string) []byte {
	base68Bytes, err := base64.StdEncoding.DecodeString(encoded)
	if err == nil {
		return base68Bytes
	}
	return []byte(encoded)
}

// DecodeSecretWallet decodes a wallet from the given encoded JSON representation
// and uses the provided password to decrypt the private key stored in the wallet.
//
// Parameters:
// - encoded: The JSON-encoded representation of the wallet.
// - password: The password used for decrypting the private key.
func DecodeSecretWallet(encodedWallet string, password string) (*SecretWallet, error) {
	var data WalletData
	toBytes := encodedwalletTobytes(encodedWallet)
	if err := json.Unmarshal(toBytes, &data); err != nil {
		return nil, err
	}
	if err := data.Validate(); err != nil {
		return nil, err
	}

	crypto := data.Crypto
	kdf := crypto.KDF

	var derivator KeyDerivator

	switch kdf {
	case "pbkdf2":
		derParams := crypto.KDFParams
		if prf, ok := derParams["prf"].(string); !ok || prf != "hmac-sha256" {
			return nil, fmt.Errorf("invalid prf supplied with the pdf: was %v, expected hmac-sha256", prf)
		}

		salt := derParams["salt"].(string)
		saltBytes := formating.HexToBytes(salt)

		derivator = NewPBDKDF2KeyDerivator(
			derParams["c"].(int),
			saltBytes,
			derParams["dklen"].(int),
		)
	case "scrypt":
		derParams := crypto.KDFParams
		salt := derParams["salt"].(string)
		saltBytes := formating.HexToBytes(salt)

		derivator = NewScryptKeyDerivator(
			int(derParams["dklen"].(float64)),
			int(derParams["n"].(float64)),
			int(derParams["r"].(float64)),
			int(derParams["p"].(float64)),
			saltBytes,
		)
	default:
		return nil, fmt.Errorf("wallet file uses %s as key derivation function, which is not supported", kdf)
	}

	encodedPassword := []byte(password)
	derivedKey, err := derivator.DeriveKey(encodedPassword)
	if err != nil {
		return nil, err
	}

	aesKey := derivedKey[:16]
	encryptedPrivateKey := formating.HexToBytes(crypto.Ciphertext)

	derivedMac := generateMac(derivedKey, encryptedPrivateKey)
	if derivedMac != crypto.MAC {
		return nil, fmt.Errorf("could not unlock wallet file. You either supplied the wrong password or the file is corrupted")
	}

	if crypto.Cipher != "aes-128-ctr" {
		return nil, fmt.Errorf("wallet file uses %s as cipher, but only aes-128-ctr is supported", crypto.Cipher)
	}

	iv := formating.HexToBytes(crypto.CipherParams["iv"])

	decrypt, err := decryptPrivateKey(encryptedPrivateKey, aesKey, iv)
	if err != nil {
		return nil, fmt.Errorf("cannot decrypt data")
	}
	id := uuid.ToBuffer(data.ID)

	return &SecretWallet{
		Credentials: string(decrypt),
		Derivator:   derivator,
		Password:    encodedPassword,
		IV:          iv,
		ID:          id,
	}, nil
}

// ToJSON serializes the SecretWallet into a JSON-encoded string. This method
// represents the wallet's data, including encryption details and credentials, in a
// JSON format that can be stored or transmitted.
func (w *SecretWallet) ToJSON() (string, error) {
	ciphertextBytes, err := w.encryptPrivateKey()
	if err != nil {
		return "", err
	}
	drive, err := w.Derivator.DeriveKey(w.Password)
	if err != nil {
		return "", err
	}
	u, err := w.uuid()
	if err != nil {
		return "", err
	}
	data := map[string]interface{}{
		"crypto": map[string]interface{}{
			"cipher":       "aes-128-ctr",
			"cipherparams": map[string]interface{}{"iv": formating.BytesToHex(w.IV)},
			"ciphertext":   formating.BytesToHex(ciphertextBytes),
			"kdf":          w.Derivator.Name(),
			"kdfparams":    w.Derivator.Encode(),
			"mac":          generateMac(drive, ciphertextBytes),
		},
		"id":      u,
		"version": 3,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

// ToBase68 encodes the credentials of the SecretWallet instance to a Base68-encoded string.
// This method represents the wallet's data, including encryption details and credentials,
// in a Base68-encoded format that can be stored or transmitted.
func (w *SecretWallet) ToBase68() (string, error) {
	toJson, err := w.ToJSON()
	if err != nil {
		return "", err
	}
	jsonBytes := []byte(toJson)
	encodedString := base64.StdEncoding.EncodeToString(jsonBytes)
	return encodedString, nil
}

// uuid returns the UUID (Universally Unique Identifier) associated with this SecretWallet.
// The UUID is a unique identifier assigned to the wallet and is used to distinguish it
// from other wallets. It is typically represented as a string in UUID format.
func (w *SecretWallet) uuid() (string, error) {
	u, err := uuid.FromBuffer(w.ID)
	return u, err
}

// encryptPrivateKey encrypts the private key associated with the SecretWallet using AES-CTR
// encryption with the provided password and initialization vector (IV).
func (w *SecretWallet) encryptPrivateKey() ([]byte, error) {
	derived, err := w.Derivator.DeriveKey(w.Password)
	if err != nil {
		return nil, err
	}
	aesKey := derived[:16]
	aesCipher, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}

	ctr := cipher.NewCTR(aesCipher, w.IV)
	ciphertext := make([]byte, len(w.Credentials))
	ctr.XORKeyStream(ciphertext, []byte(w.Credentials))
	return ciphertext, nil
}

// decryptPrivateKey decrypts a message using AES-CTR decryption with the provided AES key and initialization vector (IV).
// Parameters:
// - message: The message to be decrypted as a byte slice.
// - aesKey: The AES encryption key used for decryption.
// - iv: The initialization vector (IV) used for AES-CTR decryption.
func decryptPrivateKey(message []byte, aesKey []byte, iv []byte) ([]byte, error) {

	aesCipher, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}
	ctr := cipher.NewCTR(aesCipher, iv)
	decryptedData := make([]byte, len(message))
	ctr.XORKeyStream(decryptedData, message)
	return decryptedData, nil
}

// generateMac computes a Message Authentication Code (MAC) using SHA-256 for the given derived key (dk)
// and ciphertext. The MAC is generated by hashing the concatenation of the last 16 bytes of the derived key
// and the ciphertext with SHA-256.
func generateMac(dk, ciphertext []byte) string {
	macBody := append(dk[16:32], ciphertext...)
	mac := sha256.Sum256(macBody)
	return formating.BytesToHex(mac[:])
}
