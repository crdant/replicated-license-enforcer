package main

import (
  "flag"
	"fmt"
	"os"
  "os/signal"
	"time"

	"github.com/crdant/replicated-license-enforcer/pkg/enforce"
	"github.com/crdant/replicated-license-enforcer/pkg/version"

  "github.com/charmbracelet/log"
  cron "github.com/robfig/cron/v3"
)

func main() {
  logLevel := os.Getenv("LOG_LEVEL")
  if logLevel == "" {
    logLevel = "info"
  }
  log.ParseLevel(logLevel)
  log.Infof("Version: %s, Build Time: %s, GitCommit: %s\n", version.Version, version.BuildTime, version.GitSHA)

  recheckInterval := flag.Duration("recheck", time.Duration(0), "Recheck license periodically to assure it's still valid")
  flag.Parse()

  err := enforce.Validate()
  if err != nil {
    log.Error("Error checking license validity", "error", err)
    os.Exit(1)
  }
  
  if *recheckInterval != time.Duration(0) {
    log.Infof("Rechecking license every %v", *recheckInterval)
    c := cron.New()
    _, err := c.AddFunc(fmt.Sprintf("@every %v", recheckInterval), enforce.Recheck)
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
