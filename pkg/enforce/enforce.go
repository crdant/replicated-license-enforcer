package enforce

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/crdant/replicated-license-enforcer/pkg/client"
	"github.com/crdant/replicated-license-enforcer/pkg/events"

  "github.com/charmbracelet/log"
  cron "github.com/robfig/cron/v3"
	backoff "github.com/cenkalti/backoff/v4"
)

type Enforcer struct {
    sdkClient client.ReplicatedClient ;
    eventClient events.EventClient ;
    scheduler *cron.Cron ;
}

func DefaultEnforcer() *Enforcer {
    endpoint := os.Getenv("REPLICATED_SDK_ENDPOINT")
    if endpoint == "" {
      // this is the default in the current build of the SDK
      endpoint = "http://replicated:3000"
    } 

    sdkClient := client.NewClient(endpoint)
    eventClient, err := events.NewKubernetesEventClient()
    if err != nil {
      log.Error("Error creating Kubernetes event client", "error", err)
      return nil
    }
    return NewEnforcer(sdkClient, eventClient)
}

func NewEnforcer(sdkClient client.ReplicatedClient, eventClient events.EventClient) *Enforcer {
    return &Enforcer{sdkClient: sdkClient, eventClient: eventClient, scheduler: cron.New()}
}

func (e *Enforcer) isValid(client client.ReplicatedClient) (bool, error) {
	expiration, err := client.GetExpirationDate()
	if err != nil {
		return false, err
	}

	if expiration.Before(time.Now()) {
    return false, nil
	}

	return true, nil
}

func (e *Enforcer) Check() error {
    valid, err := e.isValid(e.sdkClient)
    if err != nil {
      log.Error("checking license", "error", err)
      return err
    }
    log.Debug("Fetching license details and creating event")

    expiration, err := e.sdkClient.GetExpirationDate()
    if err != nil {
      log.Error("Error finding expiration date to incldue in kubernetes event", "error", err)
      return err
    }

    name, _ := e.sdkClient.GetAppName()
    slug, _ := e.sdkClient.GetAppSlug()

    e.eventClient.CreateLicenseEvent(slug, expiration)
    if !valid {
      log.Infof("License for %s is expired", name)
      return errors.New(fmt.Sprintf("License for %s is expired", name))
    }

    log.Info("License is valid")
    return nil
}

func (e *Enforcer) Validate() error {
    err := backoff.Retry(e.Check, backoff.NewExponentialBackOff())
    if err != nil {
        log.Error("Error in license check, skipping current check", "error", err)
        return errors.New("Error in license check")
    }
    return nil
}
  
func (e *Enforcer) Recheck() { 
  err := e.Validate()
  if err != nil {
      log.Error("Error in license check, skipping current check", "error", err)
  }
}

func (e *Enforcer) Monitor(interval time.Duration) {
	_, err := e.scheduler.AddFunc(fmt.Sprintf("@every %v", interval), e.Recheck)
	if err != nil {
		log.Error("Could not schedule periodic license check", "error", err)
		return
	}
  defer e.Stop()  
	go e.scheduler.Start()
}

func (e *Enforcer) Stop() {
	e.scheduler.Stop()
}
