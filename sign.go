// sign.go
package cryptomus

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
)

// signRequest generates a signature for the request using the provided API key and request body.
// The signature is a hexadecimal MD5 hash of the base64-encoded request body concatenated with the API key.
func (c *Cryptomus) signRequest(apiKey string, reqBody []byte) (string, error) {
	if apiKey == "" {
		return "", errors.New("API key cannot be empty")
	}

	// Encode the request body using base64.
	data := base64.StdEncoding.EncodeToString(reqBody)

	// Compute the MD5 hash of the concatenated data and API key.
	hash := md5.Sum([]byte(data + apiKey))

	// Return the hexadecimal representation of the hash.
	return hex.EncodeToString(hash[:]), nil
}

// VerifySign verifies the signature of the incoming request.
// It checks whether the 'sign' field in the JSON body matches the expected signature.
// Parameters:
// - apiKey: The API key used for signing.
// - reqBody: The raw request body bytes.
// Returns:
// - error: Returns an error if the signature is invalid or if required fields are missing.
func (c *Cryptomus) VerifySign(apiKey string, reqBody []byte) error {
	// Unmarshal the request body into a generic map.
	var jsonBody map[string]interface{}
	err := json.Unmarshal(reqBody, &jsonBody)
	if err != nil {
		return fmt.Errorf("failed to unmarshal request body: %w", err)
	}

	// Extract the 'sign' field from the JSON body.
	reqSign, ok := jsonBody["sign"].(string)
	if !ok {
		return errors.New("missing signature field in request body")
	}

	// Remove the 'sign' field from the JSON body before generating the expected signature.
	delete(jsonBody, "sign")

	// Marshal the modified JSON body back to bytes.
	modifiedBody, err := json.Marshal(jsonBody)
	if err != nil {
		return fmt.Errorf("failed to marshal modified request body: %w", err)
	}

	// Generate the expected signature using the modified request body.
	expectedSign, err := c.signRequest(apiKey, modifiedBody)
	if err != nil {
		return fmt.Errorf("failed to generate expected signature: %w", err)
	}

	// Compare the expected signature with the one provided in the request.
	if reqSign != expectedSign {
		return errors.New("invalid signature")
	}

	return nil
}
