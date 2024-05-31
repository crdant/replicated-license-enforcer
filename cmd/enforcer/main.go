package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/crdant/replicated-license-enforcer/pkg/client"
	"github.com/crdant/replicated-license-enforcer/pkg/events"
	"github.com/crdant/replicated-license-enforcer/pkg/version"

	backoff "github.com/cenkalti/backoff/v4"
)

func checkLicense(client client.ExpirationClient) error {
	expiration, err := client.GetExpirationDate()
	if err != nil {
		return err
	}

	if expiration.Before(time.Now()) {
    client, err := events.NewKubernetesEventClient()
    if err != nil {
      return nil
    } 

    client.CreateExpiredEvent("slackernews", expiration)
		return errors.New("License is expired")
	}

	return nil
}

func main() {
	endpoint := os.Getenv("REPLICATED_SDK_ENDPOINT")
	sdkClient := client.NewClient(endpoint)
  
  fmt.Printf("Version: %s, Build Time: %s, GitCommit: %s\n", version.Version, version.BuildTime, version.GitSHA)

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
