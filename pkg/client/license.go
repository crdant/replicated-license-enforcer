package client

import (
    "fmt"
    "time"
    "encoding/json"
    "crypto"
    "crypto/rsa"
    "crypto/x509"
    "encoding/base64"
    "encoding/pem"

    log "github.com/charmbracelet/log"
    license "github.com/replicatedhq/replicated-sdk/pkg/license/types"
)

// Return the expiration date for the license as a date field since it's
// a special case of the fields in the license
func (c *Client) GetExpirationDate() (time.Time, error) {
    expiresAt, err := c.GetLicenseField("expires_at")
    if err != nil {
      return time.Time{}, err
    }
    value := expiresAt.Value.(string)
    expirationDate, err := time.Parse(time.RFC3339, value)

    return expirationDate, nil 
}

// GetLicenseField fetches a field from the license by name, and returns it only
// if it's valid
func (c *Client) GetLicenseField(field string) (*license.LicenseField, error) {
    response, err := c.makeRequest("GET", fmt.Sprintf("/api/v1/license/fields/%s", field) , nil)
    if err != nil {
        log.Debug("Error calling Replicated SDK", "error", err)
        return nil, err
    }
    defer response.Body.Close()

    var licenseField license.LicenseField
    if err := json.NewDecoder(response.Body).Decode(&licenseField); err != nil {
        log.Debug("Error decoding API response", "error", err)
        return nil, err
    }
    if err := c.verifyLicenseField(&licenseField); err != nil {
        log.Debug("Error verifying license field", "error", err)
        return nil, err
    }
    return &licenseField, nil
}

func (c *Client) verifyLicenseField(field *license.LicenseField) (error) {
    log.Debug("Verifying license field", "field", field.Name, "type", field.ValueType, "value", field.Value)
    value := ""
    ok := true

    switch(field.ValueType) {
    case "String", "Text":
      value, ok = field.Value.(string)
      if !ok {
        return fmt.Errorf("%s value is not a valid int, license may have been tampered with", field.ValueType)
      }
    case "Integer":
      value = fmt.Sprintf("%d", int(field.Value.(float64)))
      log.Debug("String value", "value", value)
    case "Boolean":
      value = fmt.Sprintf("%t", field.Value.(bool))
    }
    signature := field.Signature.V1

    pubBlock, _ := pem.Decode([]byte(publicKeyPEM))
    publicKey, err := x509.ParsePKIXPublicKey(pubBlock.Bytes)
    if err != nil {
        return fmt.Errorf("parse public key PEM: %w", err)
    }

    var opts rsa.PSSOptions
    opts.SaltLength = rsa.PSSSaltLengthAuto

    newHash := crypto.MD5
    pssh := newHash.New()
    pssh.Write([]byte(value))
    hashed := pssh.Sum(nil)

    decodedSignature, err := base64.StdEncoding.DecodeString(signature)
    if err != nil {
        return fmt.Errorf("decode signature: %w", err)
    }

    if err := rsa.VerifyPSS(publicKey.(*rsa.PublicKey), newHash, hashed, decodedSignature, &opts); err != nil {
        return fmt.Errorf("verify PSS: %w", err)
    }

    return nil
}
