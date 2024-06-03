package main

import (
	"errors"
  "flag"
	"fmt"
	"os"
  "os/signal"
	"time"

	"github.com/crdant/replicated-license-enforcer/pkg/client"
	"github.com/crdant/replicated-license-enforcer/pkg/events"
	"github.com/crdant/replicated-license-enforcer/pkg/version"

  "github.com/charmbracelet/log"
  cron "github.com/robfig/cron/v3"
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

func check() error {
    endpoint := os.Getenv("REPLICATED_SDK_ENDPOINT")
    if endpoint == "" {
      // this is the default in the current build of the SDK
      endpoint = "http://replicated:3000"
    } 
    sdkClient := client.NewClient(endpoint)

    valid, err := checkLicense(sdkClient)
    if err != nil {
      log.Error("checking license", "error", err)
      return err
    }
    log.Debug("Fetching license details and creating event")

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

    k8sClient.CreateLicenseEvent(slug, expiration)
    if !valid {
      log.Infof("License for %s is expired", name)
      return errors.New(fmt.Sprintf("License for %s is expired", name))
    }

    log.Info("License is valid")
    return nil
}

func validate() error {
    err := backoff.Retry(check, backoff.NewExponentialBackOff())
    if err != nil {
        log.Error("Error in license check, skipping current check", "error", err)
        return errors.New("Error in license check")
    }
    return nil
}
  
func recheck() { 
  err := validate()
  if err != nil {
      log.Error("Error in license check, skipping current check", "error", err)
  }
}

func main() {
  logLevel := os.Getenv("LOG_LEVEL")
  if logLevel == "" {
    logLevel = "info"
  }
  log.ParseLevel(logLevel)
  log.Infof("Version: %s, Build Time: %s, GitCommit: %s\n", version.Version, version.BuildTime, version.GitSHA)

  recheckInterval := flag.Duration("recheck", time.Duration(0), "Recheck license periodically to assure it's still valid")
  flag.Parse()

  err := validate()
  if err != nil {
    log.Error("Error checking license validity", "error", err)
    os.Exit(1)
  }
  
  if *recheckInterval != time.Duration(0) {
    log.Infof("Rechecking license every %v", *recheckInterval)
    c := cron.New()
    _, err := c.AddFunc(fmt.Sprintf("@every %v", recheckInterval), recheck)
    if err != nil {
      log.Error("Could not schedule periodic license check", "error", err)
      return
    }
    go c.Start()

    sig := make(chan os.Signal, 1)
    signal.Notify(sig, os.Interrupt, os.Kill)
    <-sig
  }
  os.Exit(0)
}
