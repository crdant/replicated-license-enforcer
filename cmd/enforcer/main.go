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

func checkLicense(client client.ReplicatedClient) (bool, error) {
	expiration, err := client.GetExpirationDate()
	if err != nil {
		return false, err
	}

	if expiration.Before(time.Now()) {
    return false, nil
	}

	return true, nil
}

func main() {
	endpoint := os.Getenv("REPLICATED_SDK_ENDPOINT")
	sdkClient := client.NewClient(endpoint)
  
  fmt.Printf("Version: %s, Build Time: %s, GitCommit: %s", version.Version, version.BuildTime, version.GitSHA)

	check := func() error {
		valid, err := checkLicense(sdkClient)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return err
		}
    if !valid {
      k8sClient, err := events.NewKubernetesEventClient()
      if err != nil {
        fmt.Fprintln(os.Stderr, "License is expired.")
        fmt.Fprintln(os.Stderr, "Error sending license expired event to Kubernetes.") 
        return err
      }
      expiration, err := sdkClient.GetExpirationDate()
      if err != nil {
        fmt.Fprintln(os.Stderr, "License is expired.")
        fmt.Fprintln(os.Stderr, "Error finding expiration date to incldue in kubernetes event") 
        return err
      }

      name, _ := sdkClient.GetAppName()
      slug, _ := sdkClient.GetAppSlug()
      k8sClient.CreateExpiredEvent(slug, expiration)
      fmt.Fprintln(os.Stderr, fmt.Sprintf("License for %s is expired", name))
      return errors.New(fmt.Sprintf("License for %s is expired", name))
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
