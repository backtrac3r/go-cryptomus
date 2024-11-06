package tests

import (
	"net/http"
	"os"
	"testing"

	"github.com/backtrac3r/go-cryptomus"
)

var TestCryptomus *cryptomus.Cryptomus

func TestMain(m *testing.M) {
	httpClient := http.Client{}
	merchant := "replace with your merchant id"
	paymentAPIKey := "replace with your payment API key"
	payoutAPIKey := "replace with your payout API key"
	TestCryptomus = cryptomus.NewCryptomus(&httpClient, merchant, paymentAPIKey, payoutAPIKey)

	os.Exit(m.Run())
}
