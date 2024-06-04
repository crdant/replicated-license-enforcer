package enforce

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/crdant/replicated-license-enforcer/pkg/client"
	"github.com/crdant/replicated-license-enforcer/pkg/events"

  "github.com/charmbracelet/log"
	backoff "github.com/cenkalti/backoff/v4"
)

func isValid(client client.ReplicatedClient) (bool, error) {
	expiration, err := client.GetExpirationDate()
	if err != nil {
		return false, err
	}

	if expiration.Before(time.Now()) {
    return false, nil
	}

	return true, nil
}

func check(sdkClient client.ReplicatedClient, k8sClient events.EventClient) error {
    valid, err := isValid(sdkClient)
    if err != nil {
      log.Error("checking license", "error", err)
      return err
    }
    log.Debug("Fetching license details and creating event")

    expiration, err := sdkClient.GetExpirationDate()
    if err != nil {
      log.Error("Error finding expiration date to incldue in kubernetes event", "error", err)
      return err
    }

    name, _ := sdkClient.GetAppName()
    slug, _ := sdkClient.GetAppSlug()

    k8sClient.CreateLicenseEvent(slug, expiration)
    if !valid {
      log.Infof("License for %s is expired", name)
      return errors.New(fmt.Sprintf("License for %s is expired", name))
    }

    log.Info("License is valid")
    return nil
}

func Check() error {
    endpoint := os.Getenv("REPLICATED_SDK_ENDPOINT")
    if endpoint == "" {
      // this is the default in the current build of the SDK
      endpoint = "http://replicated:3000"
    } 
    sdkClient := client.NewClient(endpoint)
    k8sClient, err := events.NewKubernetesEventClient()
    if err != nil {
      log.Error("Could not create Kuberenetes client", "error", err)
      return err
    }

    return check(sdkClient, k8sClient)
}

func Validate() error {
    err := backoff.Retry(Check, backoff.NewExponentialBackOff())
    if err != nil {
        log.Error("Error in license check, skipping current check", "error", err)
        return errors.New("Error in license check")
    }
    return nil
}
  
func Recheck() { 
  err := Validate()
  if err != nil {
      log.Error("Error in license check, skipping current check", "error", err)
  }
}

func Monitor() {

}
