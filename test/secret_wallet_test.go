package test

import (
	secretwallet "bitcoin/secret_wallet"
	"math/rand"
	"testing"
	"time"
)

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	randomString := make([]byte, length)
	for i := range randomString {
		randomString[i] = charset[seededRand.Intn(len(charset))]
	}

	return string(randomString)
}
func TestSecretWallet(t *testing.T) {
	for i := 0; i < 100; i++ {
		credentials := generateRandomString(500)
		password := generateRandomString(32)
		scryptN := 8192
		p := 1

		wallet, err := secretwallet.NewSecretWallet(credentials, password, scryptN, p)
		if err != nil {
			t.Error(err)
			return
		}

		jsonData, err := wallet.ToBase68()
		if err != nil {
			t.Error(err)
			return
		}

		decodedWallet, err := secretwallet.DecodeSecretWallet(jsonData, password)
		if err != nil {
			t.Error(err)
			return
		}
		if decodedWallet.Credentials != credentials {
			t.Error("invalid wallet")
			return
		}
	}

}
