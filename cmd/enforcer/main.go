package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/crdant/replicated-license-enforcer/pkg/client"
	"github.com/crdant/replicated-license-enforcer/pkg/events"
	"github.com/crdant/replicated-license-enforcer/pkg/version"

  "github.com/charmbracelet/log"
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
  log.SetLevel(log.DebugLevel)
	endpoint := os.Getenv("REPLICATED_SDK_ENDPOINT")
	sdkClient := client.NewClient(endpoint)
  
  fmt.Printf("Version: %s, Build Time: %s, GitCommit: %s\n", version.Version, version.BuildTime, version.GitSHA)

	check := func() error {
		valid, err := checkLicense(sdkClient)
		if err != nil {
			log.Error("checking license", "error", err)
			return err
		}
    if !valid {
      log.Debug("License is expired, fetching details and creating event")
      k8sClient, err := events.NewKubernetesEventClient()
      if err != nil {
        log.Error("Could not create Kuberenetes client", "error", err)
        return err
      }
      expiration, err := sdkClient.GetExpirationDate()
      if err != nil {
        log.Error("Error finding expiration date to incldue in kubernetes event", "error", err)
        return err
      }

      name, _ := sdkClient.GetAppName()
      slug, _ := sdkClient.GetAppSlug()
      k8sClient.CreateExpiredEvent(slug, expiration)
      log.Infof("License for %s is expired", name)
      return errors.New(fmt.Sprintf("License for %s is expired", name))
    }
		return nil
	}

	err := backoff.Retry(check, backoff.NewExponentialBackOff())
	if err != nil {
		os.Exit(1)
	}

	log.Info("License is valid")
	os.Exit(0)
}
