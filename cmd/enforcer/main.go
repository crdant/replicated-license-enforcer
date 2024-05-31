package main

import (
    "os"
    "fmt"
    "time"
    "errors"

    "github.com/crdant/replicated-license-enforcer/pkg/client"

    backoff "github.com/cenkalti/backoff/v4"
)

func checkLicense(client client.ExpirationClient) error {
    expiration, err := client.GetExpirationDate()
    if err != nil {
        return err
    }

    if expiration.Before(time.Now()) {
        return errors.New("License is expired")
    }

    return nil
}

func main() {
    endpoint := os.Getenv("REPLICATED_SDK_ENDPOINT")
    sdkClient := client.NewClient(endpoint)

    check := func() error {
      err := checkLicense(sdkClient)
      if err != nil {
        fmt.Fprintln(os.Stderr, err)
        return err
      }
      return nil
    }

    err := backoff.Retry(check, backoff.NewExponentialBackOff())
    if err != nil {

        os.Exit(1) 
    }

    fmt.Fprintln(os.Stdout, "License is valid.")
    os.Exit(0)
}
